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

package pkg

import (
	"auto-test-go/pkg/command"
	"auto-test-go/pkg/db"
	"auto-test-go/pkg/director"
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
	"auto-test-go/pkg/util"
	"errors"
	"os"
	"strconv"
)

type ExecuteResult struct {
	ctx interface{}
	err error
}

var MaxTaskCount = os.Getenv("MAX_TASK_COUNT")

type ExecuteCommand struct {
	Command interface{}
	TaskId  string
}

var TaskDispatch = DefaultTaskDispatch{}

type DefaultTaskDispatch struct {
	commandChan  chan *ExecuteCommand
	resultChan   chan *ExecuteResult
	counter      int
	maxTaskCount int
}

func (d *DefaultTaskDispatch) Init() error {
	d.commandChan = make(chan *ExecuteCommand, 128)
	d.resultChan = make(chan *ExecuteResult, 128)
	err := db.BoltDbManager.Init()
	if err != nil {
		return err
	}

	go d.createExecutePool()
	go d.resultHandler()
	if MaxTaskCount == "" {
		d.maxTaskCount = 200
	} else {
		count, err := strconv.Atoi(MaxTaskCount)
		if err != nil {
			d.maxTaskCount = 200
		} else {
			d.maxTaskCount = count
		}
	}
	return err
}

func (d *DefaultTaskDispatch) createExecutePool() {
	for {
		cmd, ok := <-d.commandChan
		if !ok {
			break
		}

		go func(cmd *ExecuteCommand, resultChan chan *ExecuteResult) {
			var factory director.Director
			isScenario := false
			switch cmd.Command.(type) {
			case command.ScenarioCaseExecuteCommand:
				isScenario = true
				//转换成指针传递,解决通道直接传递指针带来的gc逃逸问题
				obj := cmd.Command.(command.ScenarioCaseExecuteCommand)
				cmd.Command = &obj
				factory = director.BaseDirectorFactory{}.Create(enum.DirectorType_ScenariorCase)
				break
			case command.SingleScriptExecuteCommand:
				obj := cmd.Command.(command.SingleScriptExecuteCommand)
				cmd.Command = &obj
				factory = director.BaseDirectorFactory{}.Create(enum.DirectorType_ScriptDebugger)
				break
			case command.SingleUserCaseExecuteCommand:
				obj := cmd.Command.(command.SingleUserCaseExecuteCommand)
				cmd.Command = &obj
				factory = director.BaseDirectorFactory{}.Create(enum.DirectorType_UserCase)
				break
			case command.UserCaseBatchExecuteCommand:
				obj := cmd.Command.(command.UserCaseBatchExecuteCommand)
				cmd.Command = &obj
				factory = director.BaseDirectorFactory{}.Create(enum.DirectorType_BatchUserCase)
				break
			}

			ctx, err := factory.Action(cmd.Command, true)
			if err != nil && ctx == nil {
				if isScenario {
					result := &entities.ScenarioContext{}
					result.Self = &entities.ExecContext{}
					result.Self.Reset()
					result.Self.TaskId = cmd.TaskId
					result.Self.SetStatus(entities.Failed)
					result.Self.AddLogs("[DefaultTaskDispatch] current command execute error: " + err.Error())
					ctx = result
				} else {
					result := &entities.ExecContext{}
					result.Reset()
					result.TaskId = cmd.TaskId
					result.SetStatus(entities.Failed)
					result.AddLogs("[DefaultTaskDispatch] current command execute error: " + err.Error())
					ctx = result
				}
			}
			resultChan <- &ExecuteResult{ctx: ctx, err: err}
		}(cmd, d.resultChan)
		d.counter++
	}
}

func (d *DefaultTaskDispatch) ReceiveCommand(cmd *ExecuteCommand) error {
	if cap(d.commandChan) == len(d.commandChan) || d.counter > d.maxTaskCount {
		return errors.New("command is full, can not handler")
	} else {
		d.commandChan <- cmd
	}
	return nil
}

func (d *DefaultTaskDispatch) resultHandler() {
	for {
		result, ok := <-d.resultChan
		d.counter--
		if !ok {
			break
		}

		if result.err != nil {
			util.Logger.Error("[DefaultTaskDispatch] current command execute error: " + result.err.Error())
		}

		switch result.ctx.(type) {
		case *entities.ExecContext:
			err := db.BoltDbManager.RefreshUserContext(result.ctx.(*entities.ExecContext), true)
			if err != nil {
				util.Logger.Error("[DefaultTaskDispatch] current user command result save error: " + err.Error())
			}
			break
		case *entities.ScenarioContext:
			err := db.BoltDbManager.RefreshScenarioContext(result.ctx.(*entities.ScenarioContext), true)
			if err != nil {
				util.Logger.Error("[DefaultTaskDispatch] current scenario result save error: " + err.Error())
			}
			break
		}
	}
}
