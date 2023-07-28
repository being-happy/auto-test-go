package observer

import (
	"auto-test-go/contract"
	"github.com/astaxie/beego/logs"
)

var RunrObserver = RunCommandObserver{}

type RunCommandObserver struct {
	commpleteTaskChans map[string]chan *contract.StreamCaseCommandSubResponse
}

func (observer *RunCommandObserver) PubChannel(clientId string, streamChan chan *contract.StreamCaseCommandSubResponse) {
	if observer.commpleteTaskChans == nil {
		observer.commpleteTaskChans = make(map[string]chan *contract.StreamCaseCommandSubResponse)
	}

	ch := observer.commpleteTaskChans[clientId]
	//对于已经存在的客户端，直接剔除上次客户端
	if ch != nil {
		observer.RemoveChannelById(clientId)
	}

	observer.commpleteTaskChans[clientId] = streamChan
}

func (observer *RunCommandObserver) Notify(response *contract.StreamCaseCommandSubResponse) {
	if observer.commpleteTaskChans == nil || len(observer.commpleteTaskChans) == 0 {
		return
	}

	for _, channel := range observer.commpleteTaskChans {
		channel <- response
	}
}

func (observer *RunCommandObserver) RemoveChannelById(clientId string) {
	ch := observer.commpleteTaskChans[clientId]
	if ch != nil {
		close(ch)
	}

	delete(observer.commpleteTaskChans, clientId)
	logs.Warning("[RunCommandObserver] grpc client:" + clientId + ", channel closed.")
}

func (observer *RunCommandObserver) ExistClientId(clientId string) bool {
	return observer.commpleteTaskChans[clientId] == nil
}
