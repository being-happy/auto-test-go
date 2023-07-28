package service

import (
	"auto-test-go/contract"
	"auto-test-go/pkg/observer"
	"fmt"
	"github.com/astaxie/beego/logs"
)

type CaseCommandStreamService struct {
}

func (*CaseCommandStreamService) Listen(request *contract.CaseCommandSubRequest, server contract.StreamCaseCommandListener_ListenServer) error {
	ch := make(chan *contract.StreamCaseCommandSubResponse, 50)
	observer.RunrObserver.PubChannel(request.GetCommandId(), ch)
	select {
	case response, finish := <-ch:
		if !finish {
			err := server.SendMsg(response)
			if err != nil {
				logs.Error(fmt.Sprintf("[CaseCommandStreamService] send message to client error: %s, client id: %s", err.Error(), request.GetClientId()))
				observer.RunrObserver.RemoveChannelById(request.GetCommandId())
				return err
			}
		} else {
			break
		}
	}

	return nil
}
