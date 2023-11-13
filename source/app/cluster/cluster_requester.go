package main

import (
	"encoding/json"
	"fastbuilder-core/utils/sync_wrapper"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// 声明响应回调函数类型
type RespCb func(msg map[string]interface{})

// 声明流回调函数类型
type StreamCb func(msg string) error

// ClusterRequester 表示与集群通信的客户端
type ClusterRequester struct {
	conn            *websocket.Conn
	cbs             sync_wrapper.SyncMap[RespCb]   // 存储响应回调函数的map，带同步锁
	streamListeners sync_wrapper.SyncMap[StreamCb] // 存储流回调函数的map，带同步锁
}

// NewClusterRequester 创建一个新的 ClusterRequester 实例
func NewClusterRequester(conn *websocket.Conn) *ClusterRequester {
	r := &ClusterRequester{
		conn:            conn,
		cbs:             sync_wrapper.SyncMap[RespCb]{},
		streamListeners: sync_wrapper.SyncMap[StreamCb]{},
	}

	return r
}

// InstanceDetail 表示集群中一个实例的详细信息
type InstanceDetail struct {
	InstanceID   string   `json:"instance_id"`
	Name         string   `json:"name"`
	Cmd          string   `json:"cmd"`
	Args         []string `json:"args"`
	Status       string   `json:"status"`
	StatusDetail string   `json:"status_detail"`
}

// List 发送请求以列出集群中的实例
func (r *ClusterRequester) List(cb func([]InstanceDetail)) {
	r.sendRequest("list", map[string]interface{}{}, func(msg map[string]interface{}) {
		marshal, _ := json.Marshal(msg["instances"])
		instances := make([]InstanceDetail, 0)
		json.Unmarshal(marshal, &instances)
		cb(instances)
	})
}

// toString 将数据转换为字符串
func toString(data interface{}) string {
	if s, ok := data.(string); ok {
		return s
	} else {
		return fmt.Sprintf("%v", s)
	}
}

// Run 发送请求以在集群中启动一个新实例
func (r *ClusterRequester) Run(name, cmd string, args []string, cb func(instanceID string, err string)) {
	r.sendRequest("run", map[string]interface{}{
		"name": name,
		"cmd":  cmd,
		"args": args,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["instance_id"]), toString(msg["err"]))
	})
}

// Rm 发送请求以从集群中删除一个实例
func (r *ClusterRequester) Rm(instanceID string, cb func(err string)) {
	r.sendRequest("rm", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["err"]))
	})
}

// Status 发送请求以获取集群中实例的状态
func (r *ClusterRequester) Status(instanceID string, cb func(status, detail, name string, err string)) {
	r.sendRequest("status", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["status"]), toString(msg["status_detail"]), toString(msg["name"]), toString(msg["err"]))
	})
}

// Journal 发送请求以获取集群中实例的日志
func (r *ClusterRequester) Journal(instanceID string, lines int, cb func(journal string, err string)) {
	r.sendRequest("journal", map[string]interface{}{
		"instance_id": instanceID,
		"lines":       lines,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["journal"]), toString(msg["err"]))
	})
}

// Kill 发送请求以终止集群中的实例
func (r *ClusterRequester) Kill(instanceID string, cb func(err string)) {
	r.sendRequest("kill_instance", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		cb(toString(msg["err"]))
	})
}

// stopStream 发送请求以停止集群中实例的流
func (r *ClusterRequester) stopStream(instanceID string, streamID string) {
	r.sendRequest("stop_stream", map[string]interface{}{
		"instance_id": instanceID,
		"stream_id":   streamID,
	}, func(msg map[string]interface{}) {
		if err := toString(msg["err"]); err != "" {
			fmt.Printf("停止流失败：%s\n", err)
		}
	})
}

// Stream 发送请求以监听集群中实例的流
func (r *ClusterRequester) Stream(instanceID string, cb StreamCb) error {
	var streamID string

	// 发送请求以开始流
	r.sendRequest("stream", map[string]interface{}{
		"instance_id": instanceID,
	}, func(msg map[string]interface{}) {
		if err := toString(msg["err"]); err != "" { // 如果返回的消息中有错误信息，则打印错误日志并返回
			fmt.Printf("启动流失败：%s\n", err)
			return
		}

		if id, ok := msg["id"].(string); ok { // 如果返回的消息中有流ID，则将其记录下来
			streamID = id
		} else {
			fmt.Println("Invalid stream ID received")
			return
		}

		// 将回调函数添加到 streamListeners 中
		r.streamListeners.Store(cb, struct{}{})
	})

	// 开始循环读取流数据
	for {
		select {
		case <-r.conn.Context().Done(): // 如果连接已关闭，则停止监听流并返回错误
			r.stopStream(instanceID, streamID)
			return fmt.Errorf("连接已关闭")
		default:
			_, message, err := r.conn.ReadMessage() // 读取流数据
			if err != nil {                         // 如果读取出错，则停止监听流并返回错误
				r.stopStream(instanceID, streamID)
				return err
			}

			var msg map[string]interface{}
			err = json.Unmarshal(message, &msg) // 解析流数据为 JSON 格式
			if err != nil {                     // 如果解析出错，则打印错误日志并继续循环
				fmt.Printf("解析JSON失败：%s\n", err.Error())
				continue
			}

			if id, ok := msg["id"].(string); ok && id == streamID { // 如果流ID匹配，则调用回调函数处理流数据
				r.streamListeners.Range(func(key, value interface{}) bool {
					cb := key.(StreamCb)
					err := cb(toString(msg["data"]))
					if err != nil { // 如果回调函数返回错误，则从 streamListeners 中删除该回调函数
						r.streamListeners.Delete(cb)
					}
					return true
				})
			}
		}
	}
}

// sendRequest 发送请求并处理响应
func (r *ClusterRequester) sendRequest(method string, data map[string]interface{}, cb RespCb) {
	messageID := uuid.New().String() // 生成消息ID

	// 构造请求消息
	request := map[string]interface{}{
		"id":     messageID,
		"method": method,
		"data":   data,
	}

	// 将响应回调函数添加到 cbs 中
	r.cbs.Store(messageID, cb)

	// 发送请求消息
	r.conn.WriteJSON(request)

	// 循环读取响应消息
	for {
		select {
		case <-r.conn.Context().Done(): // 如果连接已关闭，则从 cbs 中删除响应回调函数并返回
			r.cbs.Delete(messageID)
			return
		default:
			_, message, err := r.conn.ReadMessage() // 读取响应消息
			if err != nil {                         // 如果读取出错，则从 cbs 中删除响应回调函数并返回
				r.cbs.Delete(messageID)
				return
			}

			var msg map[string]interface{}
			err = json.Unmarshal(message, &msg) // 解析响应消息为 JSON 格式
			if err != nil {                     // 如果解析出错，则打印错误日志并继续循环
				fmt.Printf("解析JSON失败：%s\n", err.Error())
				continue
			}

			if id, ok := msg["id"].(string); ok && id == messageID { // 如果消息ID匹配，则调用响应回调函数处理响应
				r.cbs.Range(func(key, value interface{}) bool {
					cb := key.(RespCb)
					cb(msg)
					r.cbs.Delete(key) // 从 cbs 中删除响应回调函数
					return true
				})
				break
			}
		}
	}
}

func main() {}
