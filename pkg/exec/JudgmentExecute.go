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

package exec

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
	"auto-test-go/pkg/script"
	"auto-test-go/pkg/util"
	"fmt"
	"strconv"
)

type JudgmentExecute struct {
}

func (u JudgmentExecute) doWork(scenarioCase *entities.ScenarioCase, conditoinUnitDesign entities.ConditoinUnitDesign, scenariorContext *entities.ScenarioContext, parentContext *entities.ExecContext) {
	judgmentContext := parentContext.Copy()
	judgmentContext.Reset()
	judgmentContext.AssertSuccess = true
	if conditoinUnitDesign.Name == "" {
		judgmentContext.Name = "JudgmentContext"
	} else {
		judgmentContext.Name = conditoinUnitDesign.Name
	}

	judgmentContext.CaseId = conditoinUnitDesign.Id
	judgmentContext.TaskId = scenariorContext.Self.TaskId
	scenariorContext.Counter++
	addTrace(scenariorContext, scenariorContext.Counter, parentContext, judgmentContext)
	if conditoinUnitDesign.Expr == "" {
		log := fmt.Sprintf("[JudgmentExecute] expr can not be empty: %s", conditoinUnitDesign.Id)
		judgmentContext.AddLogs(log)
		util.Logger.Error(log)
		judgmentContext.SetStop()
		parentContext.SetStop()
		return
	}

	template := script.FastTemplate{}
	currentLuaScript := conditoinUnitDesign.Expr
	if len(parentContext.Variables) > 0 {
		vars := template.ConvertVar(parentContext.Variables)
		currentLuaScript = template.Template(currentLuaScript, vars)
	}

	judgmentScript := u.createBaseScript(currentLuaScript)
	if len(parentContext.Variables) > 0 {
		vars := template.ConvertVar(parentContext.Variables)
		judgmentScript.Script.Script = template.Template(currentLuaScript, vars)
	}

	scriptExecute := ScriptDebuggerExecute{}
	scriptExecute.DoWork(&judgmentScript, judgmentContext)
	if judgmentContext.GetStatus() == entities.Failed {
		log := fmt.Sprintf("[JudgmentExecute] condition id: %s , execute fail.", conditoinUnitDesign.Id)
		judgmentContext.AddLogs(log)
		util.Logger.Error(log)
		judgmentContext.SetStop()
		parentContext.SetStop()
		return
	}

	judgmentContext.SetStatus(entities.Success)
	result := judgmentContext.Variables["return_value"].Value
	boolValue, err := strconv.ParseBool(result)
	if err != nil {
		log := fmt.Sprintf("[JudgmentExecute] parse bool error,condition id: %s ,error:%s ,in ScenariorCase.", conditoinUnitDesign.Id, err)
		judgmentContext.AddLogs(log)
		util.Logger.Error(log)
		judgmentContext.SetStop()
		parentContext.SetStop()
		return
	}

	scenariorContext.ExecIds = append(scenariorContext.ExecIds, conditoinUnitDesign.Id)
	log := fmt.Sprintf("[JudgmentExecute] expression:%s , result:%v ,", conditoinUnitDesign.Expr, boolValue)
	judgmentContext.AddLogs(log)
	if boolValue {
		executeFlow(scenarioCase, conditoinUnitDesign.CorrectBranch, judgmentContext, scenariorContext)
	} else {
		executeFlow(scenarioCase, conditoinUnitDesign.ErrorBranch, judgmentContext, scenariorContext)
	}

	if judgmentContext.GetStatus() != entities.Failed {
		judgmentContext.SetStatus(entities.Success)
	} else {
		parentContext.SetStop()
	}
}

func (u JudgmentExecute) createBaseScript(currentLuaScript string) entities.BaseScript {
	judgmentScript := entities.BaseScript{
		ScriptType: enum.ScriptType_LuaScript,
		Script: entities.LuaScript{
			FuncType: enum.LuaFuncType_DoJudgmentExecute,
			Script:   currentLuaScript,
		},
	}
	return judgmentScript
}
