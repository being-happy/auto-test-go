// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
