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
)

type ScriptDebuggerExecute struct {
}

func (u ScriptDebuggerExecute) DoWork(baseScripts *entities.BaseScript, ctx *entities.ExecContext) {
	var log string
	if ctx.Name == "" {
		ctx.Name = baseScripts.ScriptType
	}

	if baseScripts.ScriptType == "" {
		log = fmt.Sprintf("[ScriptDebuggerExecute] script type can not nil")
		ctx.AddLogs(log)
		util.Logger.Error(script.CombineLogInfo(log, ctx))
		return
	}

	if baseScripts.IsSkipError {
		ctx.IgnoreStop = true
	} else {
		ctx.IgnoreStop = false
	}

	u.execScript(baseScripts.ScriptType, &baseScripts.Script, ctx)
	//if ctx.GetStatus() != entities.Failed {
	//	ctx.SetStatus(entities.Success)
	//}
}

func (u ScriptDebuggerExecute) execScript(scriptType string, luaScript *entities.LuaScript, ctx *entities.ExecContext) {
	if ctx.Stop() {
		return
	}

	if luaScript.Script == "" || luaScript.FuncType == "" {
		log := fmt.Sprintf("[ScriptDebuggerExecute] LuaScript field Script or  FuncType can not be nil!, function type: %s , script type: %s", luaScript.FuncType, scriptType)
		ctx.AddLogs(log)
		util.Logger.Error(script.CombineLogInfo(log, ctx))
		ctx.SetStop()
		return
	}

	funcCtx := entities.FuncContext{
		FuncBody:     luaScript.Script,
		FuncType:     luaScript.FuncType,
		ScriptType:   enum.ScriptType_LuaScript,
		CallFunction: luaScript.CallFunction,
	}
	var err error
	switch scriptType {
	case enum.ScriptType_LuaScript:
		err = script.CaseRegister.Trigger(ctx, &funcCtx)
		break
	case enum.ScriptType_SqlScript:
		if !luaScript.ValidScript() {
			log := "[ScriptDebuggerExecute] SqlScript parameter has nil value!"
			ctx.AddLogs(log)
			util.Logger.Warn(script.CombineLogInfo(log, ctx))
			break
		}
		luaScript.CopyToFuncContext(&funcCtx)
		err = script.CaseRegister.Trigger(ctx, &funcCtx)
		break
	case enum.ScriptType_HttpCall:
		funcCtx.Request = ctx.CurrentRequest
		funcCtx.ScriptType = enum.ScriptType_HttpCall
		err = script.CaseRegister.Trigger(ctx, &funcCtx)
		break
	}

	if err != nil {
		ctx.AddLogs(err.Error())
		util.Logger.Error(err.Error())
		ctx.SetStop()
	}

	if ctx.Stop() {
		u.stop(ctx)
	}
}

func (ScriptDebuggerExecute) stop(ctx *entities.ExecContext) {
	log := "[UserCaseExecute] Use case execute stopped by script!"
	ctx.AddLogs(log)
	util.Logger.Info(script.CombineLogInfo(log, ctx))
}
