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
	"strconv"
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

func NewLuaScriptBaseHandler() (handler *LuaScriptBaseHandler) {
	handler = &LuaScriptBaseHandler{}
	handler.Name = enum.LuaFuncName_DoBaseExecute
	handler.FuncType = enum.LuaFuncType_DoBaseExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	err := handler.Init()
	if err != nil {
		panic(err)
	}
	return handler
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
	// "number", "", "function", "userdata", "thread", "table"
	if execCtx.Variables != nil {
		for k, v := range execCtx.Variables {
			if v.Object != nil {
				switch v.Object.(type) {
				case lua.LBool:
					table.RawSet(lua.LString(k), v.Object.(lua.LBool))
					break
				case lua.LNumber:
					table.RawSet(lua.LString(k), v.Object.(lua.LNumber))
					break
				case lua.LString:
					table.RawSet(lua.LString(k), v.Object.(lua.LString))
				case *lua.LTable:
					subTable := v.Object.(*lua.LTable)
					table.RawSet(lua.LString(k), subTable)
				default:
					table.RawSet(lua.LString(k), lua.LString(v.Value))
				}
			} else if v.Type != "" {
				//外部指令传递的类型
				switch v.Type {
				case "boolean":
					boolean, err := strconv.ParseBool(v.Value)
					if err != nil {
						table.RawSet(lua.LString(k), lua.LBool(boolean))
					} else {
						table.RawSet(lua.LString(k), lua.LString(v.Value))
					}
					break
				case "number":
					float, err := strconv.ParseFloat(v.Value, 32)
					if err != nil {
						table.RawSet(lua.LString(k), lua.LNumber(float))
					}
					break
				default:
					table.RawSet(lua.LString(k), lua.LString(v.Value))
				}
			} else {
				table.RawSet(lua.LString(k), lua.LString(v.Value))
			}
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
