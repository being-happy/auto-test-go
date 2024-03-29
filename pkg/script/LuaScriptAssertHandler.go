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
)

type LuaScriptAssertHandler struct {
	BaseScripHandler
	function string
}

func NewLuaScriptAssertHandler() *LuaScriptAssertHandler {
	handler := LuaScriptAssertHandler{}
	handler.Name = enum.LuaFuncName_AssertUserCase
	handler.ScriptType = enum.ScriptType_LuaScript
	handler.FuncType = enum.LuaFuncType_AssertUserCase
	err := handler.Init()
	if err != nil {
		panic(err)
	}
	return &handler
}

func (l *LuaScriptAssertHandler) Init() error {
	body, err := loadScript(enum.LuaFuncName_AssertUserCase)
	l.function = body
	return err
}

func (l *LuaScriptAssertHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return buildScript(execCtx, funcCtx, l.function, enum.LuaFuncName_AssertUserCase)
}

func (l *LuaScriptAssertHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return scriptExecuteAfterResp(enum.LuaFuncName_AssertUserCase, execCtx, funcCtx)
}
