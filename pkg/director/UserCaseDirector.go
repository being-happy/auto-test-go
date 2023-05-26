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

package director

import (
	c "auto-test-go/pkg/command"
	"auto-test-go/pkg/db"
	exec2 "auto-test-go/pkg/exec"
	"auto-test-go/pkg/util"
	"encoding/json"
	"errors"
)

type UserCaseDirector struct {
}

func NewUserCaseDirector() UserCaseDirector {
	return UserCaseDirector{}
}

func (u UserCaseDirector) Action(command c.BaseCommand, async bool) (interface{}, error) {
	caseCommand := command.(*c.SingleUserCaseExecuteCommand)
	if caseCommand == nil {
		if command != nil {
			commandStr, err := json.Marshal(&command)
			if err != nil {
				util.Logger.Warn("[UserCaseDirector] Use case command convert fail, error: " + err.Error())
			} else {
				util.Logger.Warn("[UserCaseDirector] Use case command convert fail, error: " + string(commandStr))
			}
			return nil, err
		}
		util.Logger.Warn("[UserCaseDirector] Use case command convert fail, error: command is empty!")
		return nil, errors.New("caseCommand is nil")
	}

	str, _ := json.Marshal(caseCommand)
	util.Logger.Info("[UserCaseDirector] Begin execute user case, task id: %d ,name: %s, data:%s", caseCommand.Id, caseCommand.Name, str)
	exec := exec2.UserCaseExecute{}
	userCase := caseCommand.UserCase
	ctx := exec.CreateContext(&userCase)
	ctx.TaskId = caseCommand.Id
	if async {
		err := db.BoltDbManager.RefreshUserContext(ctx, false)
		if err != nil {
			ctx.AddLogs("Current task save error:" + err.Error())
		}
	}
	exec.DoWork(&userCase, ctx)
	return ctx, nil
}
