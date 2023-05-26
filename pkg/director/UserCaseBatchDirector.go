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
	"auto-test-go/pkg/entities"
	exec2 "auto-test-go/pkg/exec"
	"auto-test-go/pkg/util"
	"encoding/json"
	"errors"
	"strconv"
)

type UserCaseBatchDirector struct {
}

func NewUserCaseBatchDirector() *UserCaseBatchDirector {
	return &UserCaseBatchDirector{}
}

func (u UserCaseBatchDirector) Action(command c.BaseCommand, async bool) (interface{}, error) {
	caseCommand := command.(*c.UserCaseBatchExecuteCommand)
	if caseCommand == nil {
		if command != nil {
			commandStr, err := json.Marshal(&command)
			if err != nil {
				util.Logger.Warn("[UserCaseBatchDirector] Use case command convert fail, error: " + err.Error())
			} else {
				util.Logger.Warn("[UserCaseBatchDirector] Use case command convert fail, error: " + string(commandStr))
			}
			return nil, err
		}
		util.Logger.Warn("[UserCaseBatchDirector] Use case command convert fail, error: command is empty!")
		return nil, errors.New("caseCommand is nil")
	}

	str, _ := json.Marshal(caseCommand)
	util.Logger.Info("[UserCaseBatchDirector] Begin execute user case, task id: %d ,name: %s, data:%s", caseCommand.Id, caseCommand.Name, str)
	exec := exec2.UserCaseExecute{}
	ctxs := make([]*entities.ExecContext, len(caseCommand.UserCases))
	var index = 0
	for _, userCase := range caseCommand.UserCases {
		ctx := &entities.ExecContext{
			CaseId: strconv.FormatInt(userCase.Id, 10),
			Name:   userCase.Name,
			TaskId: caseCommand.Id,
		}

		ctx.Variables = make(map[string]entities.VarValue)
		if userCase.Parameters != nil && len(userCase.Parameters) > 0 {
			for _, v := range userCase.Parameters {
				value := entities.VarValue{Value: v.Value, Type: v.PType}
				ctx.Variables[v.Name] = value
			}
		}

		ctx.Variables["inner_log"] = entities.VarValue{Value: "", Type: ""}
		if async {
			err := db.BoltDbManager.RefreshUserContext(ctx, false)
			if err != nil {
				ctx.AddLogs("Current task save error:" + err.Error())
			}
		}
		exec.DoWork(&userCase, ctx)
		ctxs[index] = ctx
		index++
	}
	return ctxs, nil
}
