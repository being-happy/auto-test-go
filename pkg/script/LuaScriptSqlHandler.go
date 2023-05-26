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
	"errors"
	"fmt"
	"strings"
)

type LuaScriptSqlHandler struct {
	BaseScripHandler
	function string
}

func NewLuaScriptSqlHandler() (*LuaScriptSqlHandler, error) {
	handler := LuaScriptSqlHandler{}
	handler.Name = enum.LuaFuncType_DoSqlExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	handler.FuncType = enum.LuaFuncType_DoSqlExecute
	err := handler.Init()
	return &handler, err
}

func (l *LuaScriptSqlHandler) Init() error {
	body, err := loadScript(enum.LuaFuncName_DoSqlExecute)
	l.function = body
	return err
}

func (l *LuaScriptSqlHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	if funcCtx.FuncBody == "" || funcCtx.UserName == "" || funcCtx.Password == "" || funcCtx.Host == "" || funcCtx.Port == "" {
		log := fmt.Sprintf("[LuaScriptSqlHandler] One or more of userName,password,host,port is nil!")
		execCtx.AddLogs(log)
		return errors.New(log)
	}

	funcScript := strings.Replace(l.function, "@host", funcCtx.Host, -1)
	funcScript = strings.Replace(funcScript, "@port", funcCtx.Port, -1)
	funcScript = strings.Replace(funcScript, "@dbName", funcCtx.DbName, -1)
	funcScript = strings.Replace(funcScript, "@userName", funcCtx.UserName, -1)
	funcScript = strings.Replace(funcScript, "@password", funcCtx.Password, -1)
	return buildScript(execCtx, funcCtx, funcScript, enum.LuaFuncName_DoSqlExecute)
}

func (l *LuaScriptSqlHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return scriptExecute(enum.LuaFuncName_DoSqlExecute, execCtx, funcCtx)
}
