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
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type UserCaseExecute struct {
}

func (u UserCaseExecute) DoWork(userCase *entities.UserCase, ctx *entities.ExecContext) {
	err := u.UserCaseValidation(ctx, userCase)
	if err != nil {
		ctx.AddLogs(err.Error())
		util.Logger.Error(err.Error())
		return
	}

	var log string
	// execute before script
	if userCase.PreScripts != nil && len(userCase.PreScripts) > 0 {
		sort.Sort(userCase.PreScripts)
		scriptStr, _ := json.Marshal(userCase.PreScripts)
		log = fmt.Sprintf("[UserCaseExecute] User case pre script ordered: %s", scriptStr)
		ctx.AddLogs(log)
		util.Logger.Info(log)
		u.execScripts(&userCase.PreScripts, ctx)
	}
	httpScript := entities.LuaScript{
		FuncType: enum.ProtocolTypeHttp_DoRequest,
		Script:   "local",
	}

	ctx.CurrentRequest = userCase.Request
	// http request also as a special type to execute
	ScriptDebuggerExecute{}.execScript(enum.ScriptType_HttpCall, &httpScript, ctx)
	//execute assert
	if userCase.TextAsserts != nil && len(userCase.TextAsserts) > 0 {
		for _, assert := range userCase.TextAsserts {
			switch assert.ResponseType {
			case enum.TextAssert_ResponseCode:
				ctx.AssertSuccess = u.textOperation(assert.Operation, assert.Data, strconv.Itoa(ctx.RespCode))
			case enum.TextAssert_ResponseData:
				ctx.AssertSuccess = u.textOperation(assert.Operation, assert.Data, ctx.RespBody)
			}
		}
	} else {
		ScriptDebuggerExecute{}.execScript(enum.ScriptType_LuaScript, &userCase.Assert, ctx)
	}

	//execute after script
	if userCase.AfterScripts != nil && len(userCase.AfterScripts) > 0 {
		sort.Sort(userCase.AfterScripts)
		scriptStr, _ := json.Marshal(userCase.AfterScripts)
		log = fmt.Sprintf("[UserCaseExecute] User case after script ordered: %s", scriptStr)
		ctx.AddLogs(log)
		util.Logger.Info(log)
		u.execScripts(&userCase.AfterScripts, ctx)
	}

	if ctx.GetStatus() != entities.Failed {
		ctx.SetStatus(entities.Success)
	}
	return
}

func (u UserCaseExecute) textOperation(operation string, source string, data string) bool {
	data = strings.Trim(data, " ")
	source = strings.Trim(source, " ")
	switch operation {
	case enum.OperationType_Contains:
		return strings.Contains(data, source)
	case enum.OperationType_Equals:
		return data == source
	case enum.OperationType_NoContains:
		return !strings.Contains(data, source)
	case enum.OperationType_StartWith:
		return strings.HasPrefix(data, source)
	case enum.OperationType_EndWith:
		return strings.HasSuffix(data, source)
	default:
		return false
	}
}

func (u UserCaseExecute) CreateContext(userCase *entities.UserCase) *entities.ExecContext {
	ctx := &entities.ExecContext{
		CaseId: strconv.FormatInt(userCase.Id, 10),
		Name:   userCase.Name,
	}

	ctx.Reset()
	ctx.Variables = make(map[string]entities.VarValue)
	if userCase.Parameters != nil && len(userCase.Parameters) > 0 {
		for _, v := range userCase.Parameters {
			value := entities.VarValue{Value: v.Value, Type: v.PType}
			ctx.Variables[v.Name] = value
		}
	}
	ctx.Variables["inner_log"] = entities.VarValue{Value: "", Type: ""}
	return ctx
}

func (u UserCaseExecute) execScripts(baseScripts *entities.BaseScripts, ctx *entities.ExecContext) {
	if ctx.Stop() {
		return
	}

	for _, baseScript := range *baseScripts {
		if ctx.Stop() {
			return
		}
		ScriptDebuggerExecute{}.DoWork(&baseScript, ctx)
	}
}

func (u UserCaseExecute) UserCaseValidation(execCtx *entities.ExecContext, userCase *entities.UserCase) error {
	if execCtx.CaseId == "" || execCtx.Name == "" || execCtx.TaskId == "" {
		return errors.New("[UserCaseExecute] Use case context is not valid, lost parameters")
	}

	if !userCase.Request.Validation() {
		return errors.New("[UserCaseExecute] HttpRequest lost parameters")
	}
	return nil
}
