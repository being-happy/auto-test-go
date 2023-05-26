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
	"auto-test-go/pkg/enum"
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
)

type ScenarioController struct {
	beego.Controller
}

func (s ScenarioController) ReceiveSingelSync() {
	var execCommand command.ScenarioCaseExecuteCommand
	res := entities.Response{}
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &execCommand)
	if err != nil {
		res.Code = enum.Resp_Error
		res.Message = err.Error()
	} else {
		factory := director.BaseDirectorFactory{}.Create(enum.DirectorType_ScenariorCase)
		db.BoltDbManager.Statistics.AddDoneScenarioCaseCount()
		ctx, err := factory.Action(&execCommand, false)
		db.BoltDbManager.Statistics.AddDoneScenarioCaseCount()
		if err != nil {
			res.Code = enum.Resp_Forbid
			res.Message = err.Error()
		} else {
			res.Code = enum.Resp_OK
			res.Data = ctx
		}
	}

	s.Data["json"] = res
	s.ServeJSON()
}

func (s ScenarioController) ReceiveSingelAsync() {
	var execCommand command.ScenarioCaseExecuteCommand
	res := entities.Response{}
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &execCommand)
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

	s.Data["json"] = res
	s.ServeJSON()
}

func (s ScenarioController) QueryCommand() {
	taskId := s.Ctx.Input.Param(":taskId")
	parentId := s.GetString("parentId", "")
	ctx, err := db.BoltDbManager.QueryScenarioTask(taskId)
	res := entities.Response{}
	if err != nil {
		res.Code = enum.Resp_Error
		res.Message = err.Error()
	} else {
		if parentId != "" && ctx != nil {
			parseInt, err := strconv.Atoi(parentId)
			if err != nil {
				res.Code = enum.Resp_Error
				res.Message = "parentId is inValid"
			} else {
				res.Code = enum.Resp_OK
				res.Data = ctx.Filter(parseInt)
			}
		} else {
			res.Code = enum.Resp_OK
			res.Data = ctx
		}
	}

	s.Data["json"] = res
	s.ServeJSON()
}

func (s ScenarioController) QueryStatus() {
	res := entities.Response{
		Code: enum.Resp_OK,
		Data: db.BoltDbManager.Statistics,
	}

	s.Data["json"] = res
	s.ServeJSON()
}
