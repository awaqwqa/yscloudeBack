package cluster

import (
	"context"
	"encoding/json"
	"fastbuilder-core/utils/sync_wrapper"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"yscloudeBack/source/app/utils"
)

type RespCb func(msg map[string]interface{})
type StreamCb func(msg string) error

type ClusterRequester struct {
	conn *websocket.Conn
	// 连接是否正常
	ConnectStatus   bool
	cbs             sync_wrapper.SyncMap[RespCb]
	streamListeners sync_wrapper.SyncMap[StreamCb]
}

func NewClusterRequester() *ClusterRequester {
	r := &ClusterRequester{
		cbs:             sync_wrapper.SyncMap[RespCb]{},
		streamListeners: sync_wrapper.SyncMap[StreamCb]{},
		ConnectStatus:   false,
	}

	return r
}
func (cr *ClusterRequester) Init(host string) error {
	conn, _, err := websocket.DefaultDialer.Dial(host, http.Header{})
	if err != nil {
		utils.Error(err.Error())
		cr.ConnectStatus = false
		return err
	}
	cr.conn = conn
	cr.ConnectStatus = true
	return nil
}

type InstanceDetail struct {
	InstanceID   string   `json:"instance_id"`
	Name         string   `json:"name"`
	Cmd          string   `json:"cmd"`
	Args         []string `json:"args"`
	Status       string   `json:"status"`
	StatusDetail string   `json:"status_detail"`
}

func (r *ClusterRequester) List(cb func([]InstanceDetail)) {
	r.sendRequest("list", map[string]interface{}{}, func(msg map[string]interface{}) {
		marshal, _ := json.Marshal(msg["instances"])
		instances := make([]InstanceDetail, 0)
		json.Unmarshal(marshal, &instances)
		cb(instances)
	})
}

func toString(data interface{}) string {
	if s, ok := data.(string); ok {
		return s
	} else {
		return fmt.Sprintf("%v", s)
	}
}

// taskName cmd为builder_wrapper地址 args 为需要携带的参数
// args主要是:--convert-dir 转化地址 --file为文件地址 --user-token为本次导入的token --server为服务器code --pos [x,y,z]为导入的地址
func (r *ClusterRequester) Run(name, cmd string, args []string, cb func(instanceID string, err string)) {
	r.sendRequest("run", map[string]interface{}{
		"name": name,
		"cmd":  cmd,
		"args": args,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["instance_id"]), toString(msg["err"]))
	})
}

func (r *ClusterRequester) Rm(instanceID string, cb func(err string)) {
	r.sendRequest("rm", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["err"]))
	})
}

func (r *ClusterRequester) Status(instanceID string, cb func(status, detail, name string, err string)) {
	r.sendRequest("status", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["status"]), toString(msg["status_detail"]), toString(msg["name"]), toString(msg["err"]))
	})
}

func (r *ClusterRequester) Journal(instanceID string, lines int, cb func(journal string, err string)) {
	r.sendRequest("journal", map[string]interface{}{
		"instance_id": instanceID,
		"lines":       lines,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["journal"]), toString(msg["err"]))
	})
}

func (r *ClusterRequester) Kill(instanceID string, cb func(err string)) {
	r.sendRequest("kill_instance", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["err"]))
	})
}

func (r *ClusterRequester) stopStream(instanceID string, streamID string) {
	r.sendRequest("stop_stream", map[string]interface{}{
		"instance_id": instanceID,
		"stream_id":   streamID,
	}, func(msg map[string]interface{}) {
		if err := toString(msg["err"]); err == "" {
			fmt.Printf("stream %v of %v stopped\n", streamID, instanceID)
		} else {
			fmt.Printf("stream %v of %v stop fail: err=%v\n", streamID, instanceID)
		}
	})
}

func (r *ClusterRequester) startStream(instanceID string, resultCb func(streamID string, stopFn func(), err string), streamCB func(msg string) error) {
	r.sendRequest("start_stream", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		streamID := toString(msg["stream"])
		if err := toString(msg["err"]); err == "" {
			r.streamListeners.Set(streamID, streamCB)
			resultCb(streamID, func() {
				r.stopStream(instanceID, streamID)
			}, "")
		} else {
			resultCb("", nil, err)
		}
	})
}

func (r *ClusterRequester) sendRequest(action string, req map[string]interface{}, cb RespCb) {
	opID := uuid.New().String()
	req["op_id"] = opID
	req["action"] = action
	if cb != nil {
		r.cbs.Set(opID, cb)
	}
	err := r.conn.WriteJSON(req)
	if err != nil {
		utils.Error(err.Error())
	}
}

func (r *ClusterRequester) InitReadLoop(ctx context.Context) (err error) {
	var data []byte

	for {
		var msg map[string]interface{}
		_, data, err = r.conn.ReadMessage()
		if err != nil {
			utils.Error(err.Error())
			r.ConnectStatus = false
			return
		}
		if ctx.Err() != nil {
			r.ConnectStatus = false
			return err
		}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			utils.Error(err.Error())
			r.ConnectStatus = false
			return err
		}
		opID, found := msg["op_id"]
		if !found {
			utils.Error("msg no op_id: %v\n", msg)
			return
		}
		// 出现stream_id 这里可以暂时理解这个id就是在开始文件导入的时候才有的
		if streamID, found := msg["stream_id"]; found {
			l, found := r.streamListeners.Get(streamID.(string))
			if !found {
				utils.Error("cant find streamId %v function", streamID.(string))
				return
			}
			// 发送消息
			if err := l(toString(msg["msg"])); err != nil {
				utils.Error(err.Error())
				r.streamListeners.Delete(streamID.(string))
			}
			return
		}
		// 查看是否有对应id
		if l, found := r.cbs.Get(opID.(string)); found {
			// 这里调用函数
			// 然后删除
			l(msg)
			r.cbs.Delete(opID.(string))
		}

	}
}
