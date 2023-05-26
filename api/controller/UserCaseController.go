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

package controller

import (
	"auto-test-go/pkg"
	"auto-test-go/pkg/command"
	"auto-test-go/pkg/db"
	"auto-test-go/pkg/director"
	"auto-test-go/pkg/entities"
	enum "auto-test-go/pkg/enum"
	"encoding/json"
	"github.com/astaxie/beego"
)

type UserCaseController struct {
	beego.Controller
}

//@router
func (o UserCaseController) Prepare() {
	o.EnableXSRF = false
}

func (o UserCaseController) Receive() {
}

func (o UserCaseController) GetUserCase() {
	caseName := o.GetString("caseName")
	res := entities.Response{}
	if caseName == "" {
		res.Code = enum.Resp_Invalid
		res.Message = "case name is empty!"
		o.Data["json"] = res
		o.ServeJSON()
		return
	}

	res.Data = entities.BaseRequest{Name: caseName}
	o.Data["json"] = res
	o.ServeJSON()
}

func (o UserCaseController) ReceiveSingleSync() {
	var execCommand command.SingleUserCaseExecuteCommand
	res := entities.Response{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &execCommand)
	if err != nil {
		res.Code = enum.Resp_Error
		res.Message = err.Error()
	} else {
		factory := director.BaseDirectorFactory{}.Create(enum.DirectorType_UserCase)
		db.BoltDbManager.Statistics.AddDoUserCaseCount()
		ctx, err := factory.Action(&execCommand, false)
		db.BoltDbManager.Statistics.AddDoneUserCaseCount()
		if err != nil {
			res.Code = enum.Resp_Invalid
			res.Message = err.Error()
		} else {
			res.Code = enum.Resp_OK
			res.Data = ctx
		}
	}

	o.Data["json"] = res
	o.ServeJSON()
}

func (o UserCaseController) ReceiveSingleAsync() {
	var execCommand command.SingleUserCaseExecuteCommand
	res := entities.Response{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &execCommand)
	if err != nil {
		res.Code = enum.Resp_Error
		res.Message = err.Error()
	} else {
		err = pkg.TaskDispatch.ReceiveCommand(&pkg.ExecuteCommand{Command: execCommand, TaskId: execCommand.Id})
		if err != nil {
			res.Code = enum.Resp_Error
			res.Message = err.Error()
		} else {
			res.Code = enum.Resp_OK
			res.Message = "resceive command success."
		}
	}

	o.Data["json"] = res
	o.ServeJSON()
}

func (o UserCaseController) ReceiveBatchSync() {
	var execCommand command.UserCaseBatchExecuteCommand
	res := entities.Response{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &execCommand)
	if err != nil {
		res.Code = enum.Resp_Error
		res.Message = err.Error()
	} else {
		switch execCommand.CommandType {
		case enum.CommandType_Execute:
			factory := director.BaseDirectorFactory{}.Create(enum.DirectorType_BatchUserCase)
			ctx, err := factory.Action(&execCommand, false)
			if err != nil {
				res.Code = enum.Resp_Invalid
				res.Message = err.Error()
			} else {
				res.Code = enum.Resp_OK
				res.Data = ctx
			}
			break
		}
	}

	o.Data["json"] = res
	o.ServeJSON()
}

func (o UserCaseController) QueryCommand() {
	taskId := o.Ctx.Input.Param(":taskId")
	ctx, err := db.BoltDbManager.QueryUserTask(taskId)
	res := entities.Response{}
	if err != nil {
		res.Code = enum.Resp_Error
		res.Message = err.Error()
	} else {
		res.Code = enum.Resp_OK
		res.Data = ctx
	}

	o.Data["json"] = res
	o.ServeJSON()
}
