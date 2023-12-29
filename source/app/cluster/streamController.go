package cluster

import (
	"fastbuilder-core/utils/sync_wrapper"
	"fmt"
	"yscloudeBack/source/app/utils"
)

type InstanceStream struct {
	listeners sync_wrapper.SyncMap[func(msg string) error]
	streamID  string
	stopFn    func()
}

type StreamController struct {
	client    *ClusterRequester
	instances sync_wrapper.SyncMap[*InstanceStream]
}

func NewStreamController(clinet *ClusterRequester) *StreamController {
	return &StreamController{client: clinet, instances: sync_wrapper.SyncMap[*InstanceStream]{}}
}

func (c *StreamController) AttachListener(instanceID string, cb func(msg string) error) (stopFn func(), err error) {
	instance, found := c.instances.Get(instanceID)
	if !found {
		choker := make(chan string, 1)
		c.client.startStream(instanceID, func(streamID string, stopFn func(), err string) {
			defer func() {
				choker <- err
				close(choker)
			}()
			if err != "" {
				return
			}

			instance = &InstanceStream{
				listeners: sync_wrapper.SyncMap[func(msg string) error]{},
				streamID:  streamID,
				stopFn:    stopFn,
			}
			c.instances.Set(instanceID, instance)

		}, func(msg string) error {
			failureListeners := make([]string, 0)
			oks := 0
			instance.listeners.Iter(func(k string, v func(msg string) error) (continueInter bool) {
				if v(msg) != nil {
					failureListeners = append(failureListeners, k)
					return true
				}
				oks++
				return true
			})
			for _, k := range failureListeners {
				instance.listeners.Delete(k)
			}
			if oks == 0 {
				instance.stopFn()
				c.instances.Delete(instanceID)
			}
			return nil
		})
		errS := <-choker
		if errS != "" {
			return nil, fmt.Errorf(errS)
		}
	}
	listenerID, err := utils.GenerateRandomKey()
	if err != nil {
		return nil, err
	}
	instance.listeners.Set(listenerID, cb)
	return func() {
		instance.listeners.Delete(listenerID)
	}, nil
}
