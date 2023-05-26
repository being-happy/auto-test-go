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

type LuaScriptHttpCallHandler struct {
	BaseScripHandler
	function string
}

func NewLuaScriptDoHttpCallHandler() (*LuaScriptHttpCallHandler, error) {
	handler := LuaScriptHttpCallHandler{}
	handler.Name = enum.LuaFuncType_DoHttpRequest
	handler.ScriptType = enum.ScriptType_LuaScript
	handler.FuncType = enum.LuaFuncType_DoHttpRequest
	err := handler.Init()
	return &handler, err
}

func (l *LuaScriptHttpCallHandler) Init() error {
	body, err := loadScript(enum.LuaFuncName_DoHttpRequest)
	l.function = body
	return err
}

func (l *LuaScriptHttpCallHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return buildScript(execCtx, funcCtx, l.function, enum.LuaFuncName_DoHttpRequest)
}

func (l *LuaScriptHttpCallHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return scriptExecute(enum.LuaFuncName_DoHttpRequest, execCtx, funcCtx)
}
