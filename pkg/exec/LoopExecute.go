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
	"auto-test-go/pkg/util"
	"fmt"
	"strconv"
)

type LoopExecute struct {
}

func (u LoopExecute) doWork(scenarioCase *entities.ScenarioCase, loopDesign entities.LoopCaseDesign, scenariorContext *entities.ScenarioContext, parentContext *entities.ExecContext) {
	currentCtx := parentContext.Copy()
	currentCtx.Reset()
	currentCtx.ParentId = parentContext.Id
	scenariorContext.Counter++
	currentCtx.TaskId = scenariorContext.Self.TaskId
	currentCtx.Name = "LoopContext"
	addTrace(scenariorContext, scenariorContext.Counter, parentContext, currentCtx)
	scenariorContext.ExecIds = append(scenariorContext.ExecIds, loopDesign.Id)
	switch loopDesign.Ltype {
	case enum.LoopType_Data:
		data := scenarioCase.Baggages[loopDesign.DataId]
		if data == nil || data.Data == nil {
			log := fmt.Sprintf("[LoopExecute] Can not find data id: %s for loop: %s", loopDesign.DataId, loopDesign.Id)
			currentCtx.AddLogs(log)
			util.Logger.Error(log)
			currentCtx.Stop()
			break
		}
		//loop each
		for _, row := range data.Data {
			if currentCtx.Stop() {
				break
			}
			currentCtx.Variables["CURRENT_ROW"] = entities.VarValue{
				Value: row,
			}

			log := fmt.Sprintf("[LoopExecute] Current row:%s", row)
			currentCtx.AddLogs(log)
			//将值带入循环每一个步骤
			executeFlow(scenarioCase, loopDesign.Flows, currentCtx, scenariorContext)
			waterContext := currentCtx.Copy()
			waterContext.Name = "WaterContext"
			scenariorContext.Counter++
			addTrace(scenariorContext, scenariorContext.Counter, currentCtx, waterContext)
		}
		break
	case enum.LoopType_Normal:
		i := 0
		max, _ := strconv.Atoi(loopDesign.LoopCount)
		for i < max {
			if currentCtx.Stop() {
				break
			}
			log := fmt.Sprintf("[LoopExecute] Current row index:%d", i)
			currentCtx.AddLogs(log)
			i++
			currentCtx.Variables["CURRENT_INDEX"] = entities.VarValue{
				Value: string(rune(i)),
			}

			executeFlow(scenarioCase, loopDesign.Flows, currentCtx, scenariorContext)
			waterContext := currentCtx.Copy()
			waterContext.Name = "WaterContext"
			scenariorContext.Counter++
			addTrace(scenariorContext, scenariorContext.Counter, currentCtx, waterContext)
		}
		break
	}

	if currentCtx.GetStatus() != entities.Failed {
		currentCtx.SetStatus(entities.Success)
	}
	parentContext.Merge(currentCtx)
}
