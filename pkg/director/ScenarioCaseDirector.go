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
	"auto-test-go/pkg/exec"
	"auto-test-go/pkg/util"
	"encoding/json"
	"errors"
	"strconv"
)

type ScenarioCaseDirector struct {
}

func NewScenarioCaseDirector() ScenarioCaseDirector {
	return ScenarioCaseDirector{}
}

func (u ScenarioCaseDirector) Action(command c.BaseCommand, async bool) (interface{}, error) {
	// specially no min execute unit has no executer,such as senarior and loop and so on.
	caseCommand := command.(*c.ScenarioCaseExecuteCommand)
	if caseCommand == nil {
		if command != nil {
			commandStr, err := json.Marshal(&command)
			if err != nil {
				util.Logger.Error("[ScenarioCaseDirector] Script command convert fail, error: " + err.Error())
			} else {
				util.Logger.Error("[ScenarioCaseDirector] Script command convert fail, error: " + string(commandStr))
			}
			return nil, err
		}
		util.Logger.Warn("[ScenarioCaseDirector] Script command convert fail, error: command is empty!")
		return nil, errors.New("script debugger command is nil")
	}

	str, _ := json.Marshal(caseCommand)
	util.Logger.Info("[ScenarioCaseDirector] Begin execute script case, task id: %d ,name: %s, data:%s", caseCommand.Id, caseCommand.Name, str)
	senarioCase := caseCommand.ScenarioCase
	if senarioCase.Design == "" {
		return nil, errors.New("senario design can not be empty")
	}

	design, err := exec.SenarioXmlResolver{}.ResolveDesign(senarioCase.Design)
	if err != nil {
		return nil, err
	}

	ctx := &entities.ScenarioContext{}
	ctx.Self = &entities.ExecContext{
		CaseId: strconv.FormatInt(senarioCase.Id, 10),
		Name:   senarioCase.Name,
		TaskId: caseCommand.Id,
		Id:     0,
	}

	ctx.Self.Variables = make(map[string]entities.VarValue)
	if senarioCase.Parameters != nil && len(senarioCase.Parameters) > 0 {
		for _, v := range senarioCase.Parameters {
			value := entities.VarValue{Value: v.Value, Type: v.PType}
			ctx.Self.Variables[v.Name] = value
		}
	}

	ctx.ExecIds = []string{}
	ctx.Self.Variables["inner_log"] = entities.VarValue{Value: "", Type: ""}
	ctx.Self.ParentId = 0
	ctx.Self.Reset()
	if async {
		err := db.BoltDbManager.RefreshScenarioContext(ctx, false)
		if err != nil {
			ctx.Self.AddLogs("Current task save error:" + err.Error())
		}
	}
	exec.ScenarioCaseExecute{}.DoWork(&senarioCase, ctx, design)
	return ctx, nil
}
