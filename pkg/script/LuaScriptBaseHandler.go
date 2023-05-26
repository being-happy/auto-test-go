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

package script

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
	"auto-test-go/pkg/util"
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

type BaseScripHandler struct {
	FuncType   string
	ScriptType string
	Name       string
}

type IScriptHandler interface {
	BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error
	Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error
	Init() (err error)
}

type LuaScriptBaseHandler struct {
	BaseScripHandler
	function string
}

func NewLuaScriptBaseHandler() (handler *LuaScriptBaseHandler, err error) {
	handler = &LuaScriptBaseHandler{}
	handler.Name = enum.LuaFuncName_DoBaseExecute
	handler.FuncType = enum.LuaFuncType_DoBaseUserCaseExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	err = handler.Init()
	return handler, err
}

func (l *LuaScriptBaseHandler) Init() (err error) {
	l.function, err = loadScript(enum.LuaFuncName_DoBaseExecute)
	return err
}

func (l LuaScriptBaseHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return buildScript(execCtx, funcCtx, l.function, enum.LuaFuncName_DoBaseExecute)
}

func (l LuaScriptBaseHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) (err error) {
	return scriptExecute(enum.LuaFuncName_DoBaseExecute, execCtx, funcCtx)
}

func prepareLuaCtx(execCtx *entities.ExecContext) *lua.LTable {
	table := lua.LTable{}
	if execCtx.Variables != nil {
		for k, v := range execCtx.Variables {
			table.RawSet(lua.LString(k), lua.LString(v.Value))
		}
	}
	return &table
}

func refreshLuaCtx(ret lua.LValue, execCtx *entities.ExecContext, jIT *lua.LState) {
	if ret.Type().String() == "table" {
		table := ret.(*lua.LTable)
		execCtx.Refresh(jIT, table)
		log := fmt.Sprintf("[LuaScriptHandler] Context is refresh after script executed, context vars: %s", execCtx.GetStringVariables())
		execCtx.AddLogs(log)
		util.Logger.Info(CombineLogInfo(log, execCtx))
	} else {
		log := fmt.Sprintf("[LuaScriptHandler] Can not convert lua response data to ctx, please check script is correct.")
		execCtx.AddLogs(log)
		util.Logger.Warn(CombineLogInfo(log, execCtx))
	}
}

func CombineLogInfo(log string, execCtx *entities.ExecContext) string {
	return fmt.Sprintf("%s, id: %d, name: %s, task id: %s", log, execCtx.Id, execCtx.Name, execCtx.TaskId)
}
