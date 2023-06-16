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
	"auto-test-go/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"
)

type ScenarioCaseExecute struct {
}

func (u ScenarioCaseExecute) DoWork(scenariorCase *entities.ScenarioCase, ctx *entities.ScenarioContext, flow *entities.SenarioFlowDesign) {
	err := u.ScenarioCaseValidation(ctx.Self)
	if err != nil {
		ctx.Self.AddLogs(err.Error())
		util.Logger.Error(err.Error())
		return
	}

	var log string
	// execute before script
	if scenariorCase.DependFunctions != nil {
		copyDenpendFunctions(scenariorCase.DependFunctions, scenariorCase.PreScripts)
		copyDenpendFunctions(scenariorCase.DependFunctions, scenariorCase.AfterScripts)
	}

	if scenariorCase.PreScripts != nil && len(scenariorCase.PreScripts) > 0 {
		sort.Sort(scenariorCase.PreScripts)
		scriptStr, _ := json.Marshal(scenariorCase.PreScripts)
		log = fmt.Sprintf("[ScenarioCaseExecute] User case pre script ordered: %s", scriptStr)
		ctx.Self.AddLogs(log)
		util.Logger.Info(log)
		UserCaseExecute{}.execScripts(&scenariorCase.PreScripts, ctx.Self)
	}

	executeFlow(scenariorCase, flow.Flows, ctx.Self, ctx)
	//execute after script
	if scenariorCase.AfterScripts != nil && len(scenariorCase.AfterScripts) > 0 {
		sort.Sort(scenariorCase.AfterScripts)
		scriptStr, _ := json.Marshal(scenariorCase.AfterScripts)
		log = fmt.Sprintf("[ScenarioCaseExecute] User case after script ordered: %s", scriptStr)
		ctx.Self.AddLogs(log)
		util.Logger.Info(log)
		UserCaseExecute{}.execScripts(&scenariorCase.AfterScripts, ctx.Self)
	}

	if ctx.Self.GetStatus() != entities.Failed {
		ctx.Self.SetStatus(entities.Success)
	}
	return
}

func (ScenarioCaseExecute) ScenarioCaseValidation(execCtx *entities.ExecContext) error {
	if execCtx.CaseId == "" || execCtx.Name == "" || execCtx.TaskId == "" {
		return errors.New("[ScenarioCaseValidation] Scenario case context is not valid, lost parameters")
	}
	return nil
}

func copyDenpendFunctions(dependFunctions []string, scripts entities.BaseScripts) {
	for _, s := range scripts {
		s.Script.DependFunctions = dependFunctions
	}
}

func executeFlow(scenariorCase *entities.ScenarioCase, flows entities.Flows, parentContext *entities.ExecContext, scenariorContext *entities.ScenarioContext) {
	//上下文的id 由父节点的id 开始一次向下递增
	for _, flow := range flows {
		if parentContext.Stop() {
			break
		}

		switch flow.(type) {
		case entities.LoopCaseDesign:
			loopDesign := flow.(entities.LoopCaseDesign)
			LoopExecute{}.doWork(scenariorCase, loopDesign, scenariorContext, parentContext)
			break
		case entities.UserCaseUnitDesign:
			caseDesign := flow.(entities.UserCaseUnitDesign)
			userCase := scenariorCase.UserCases[caseDesign.Id]
			if userCase == nil {
				log := fmt.Sprintf("[ScenarioCaseExecute] Can not find user case id: %s ,in ScenariorCase.", caseDesign.Id)
				parentContext.AddLogs(log)
				util.Logger.Error(log)
				parentContext.SetStop()
				break
			}

			if scenariorCase.DependFunctions != nil {
				//复制通用函数到用例中
				userCase.DependFunctions = scenariorCase.DependFunctions
			}

			execute := UserCaseExecute{}
			userCaseContext := execute.CreateContext(userCase)
			userCaseContext.Merge(parentContext)
			userCaseContext.TaskId = scenariorContext.Self.TaskId
			scenariorContext.Counter++
			addTrace(scenariorContext, scenariorContext.Counter, parentContext, userCaseContext)
			execute.DoWork(userCase, userCaseContext)
			//如果当前节点执行失败，且没有设置为跳过，则会影响父级节点继续执行，所以在每一个case 执行之前必须先判定父级节点的状态
			if !userCase.IsSkipError && userCaseContext.GetStatus() == entities.Failed {
				log := fmt.Sprintf("[ScenarioCaseExecute] User case id: %s , execute fail.", caseDesign.Id)
				parentContext.AddLogs(log)
				util.Logger.Error(log)
				parentContext.SetStop()
				break
			}
			//当前case 执行结束后，将结果合并到父级节点，带下一个节点使用
			parentContext.Merge(userCaseContext)
			break
		case entities.TimeWaitUnitDesign:
			timeContext := parentContext.Copy()
			timeDesign := flow.(entities.TimeWaitUnitDesign)
			timeContext.AddLogs(fmt.Sprintf("[ScenarioCaseExecute] Begin Time Wait on time: %s", timeDesign.WaitTime))
			timeContext.TaskId = scenariorContext.Self.TaskId
			timeContext.Name = "TimeContext"
			timeContext.Reset()
			count, _ := strconv.Atoi(timeDesign.WaitTime)
			<-time.After(time.Duration(count) * time.Second)
			timeContext.AddLogs(fmt.Sprintf("[ScenarioCaseExecute] End Time Wait on time: %s", timeDesign.WaitTime))
			timeContext.SetStatus(entities.Success)
			scenariorContext.Counter++
			addTrace(scenariorContext, scenariorContext.Counter, parentContext, timeContext)

		case entities.ScriptUnitDesign:
			scriptDesign := flow.(entities.ScriptUnitDesign)
			script := scenariorCase.Scripts[scriptDesign.Id]
			if script == nil {
				log := fmt.Sprintf("[ScenarioCaseExecute] Can not find script case id: %s ,in ScenariorCase.", scriptDesign.Id)
				parentContext.AddLogs(log)
				util.Logger.Error(log)
				parentContext.SetStop()
				break
			}

			if scenariorCase.DependFunctions != nil {
				//复制通用函数到脚本中
				script.Script.DependFunctions = scenariorCase.DependFunctions
			}
			scriptContext := parentContext.Copy()
			scriptContext.Name = "ScriptContext"
			scriptExecute := ScriptDebuggerExecute{}
			scriptContext.CaseId = scriptDesign.Id
			scenariorContext.Counter++
			scriptContext.TaskId = scenariorContext.Self.TaskId
			addTrace(scenariorContext, scenariorContext.Counter, parentContext, scriptContext)
			scriptExecute.DoWork(script, scriptContext)
			if !script.IsSkipError && scriptContext.GetStatus() == entities.Failed {
				log := fmt.Sprintf("[ScenarioCaseExecute] User case id: %s , execute fail.", scriptDesign.Id)
				parentContext.AddLogs(log)
				util.Logger.Error(log)
				parentContext.SetStop()
				break
			} else {
				scriptContext.SetStatus(entities.Success)
			}
			parentContext.Merge(scriptContext)
			break
		}
	}
}

func addTrace(scenariorContext *entities.ScenarioContext, beforeId int, parentContext *entities.ExecContext, currentContext *entities.ExecContext) {
	if scenariorContext.ContextsTrace == nil {
		scenariorContext.ContextsTrace = map[int][]*entities.ExecContext{}
	}

	contexts := scenariorContext.ContextsTrace[parentContext.Id]
	if contexts == nil {
		contexts = []*entities.ExecContext{}
	}

	currentContext.Id = beforeId
	currentContext.ParentId = parentContext.Id
	scenariorContext.ContextsTrace[parentContext.Id] = append(contexts, currentContext)
}
