# ClusterRequester 模块

- run 符合传入行动指针 然后参数以map形式传入 并且传入回调函数

  ```go
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
  
  ```

  

- sendRequests

  - 传入指令`run` 之类的 分配一个opID以及向上下文中传入对应的回调方法

  ```go
  func (r *ClusterRequester) sendRequest(action string, req map[string]interface{}, cb RespCb) {
  	opID := uuid.New().String()
  	req["op_id"] = opID
  	req["action"] = action
  	if cb != nil {
  		utils.Error("启动1")
  		r.cbs.Set(opID, cb)
  	}
  	err := r.conn.WriteJSON(req)
  	if err != nil {
  		utils.Error(err.Error())
  	}
  }
  
  ```

  

- InitReadLoop函数

  - 从通道中接受信息
  - 然后判断steam_id

  - 分两种情况 一种是有stream的和无的 都会获取他们的callback 来调用

  ```go
  
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
  			// 执行call back
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
  
  ```

  