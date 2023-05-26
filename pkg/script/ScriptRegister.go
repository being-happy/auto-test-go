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
	"auto-test-go/pkg/util"
	"context"
	"errors"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
	"os"
	"strings"
	"time"
)

type CaseScriptHandleRegister interface {
	Register(funcType string, scriptType string, handler IScriptHandler)
	Trigger(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error
}

type UserCaseScriptHandleRegister struct {
	routes map[string]IScriptHandler
}

var CaseRegister = NewCaseScriptHandleRegister()

func NewCaseScriptHandleRegister() CaseScriptHandleRegister {
	return &UserCaseScriptHandleRegister{}
}

func (register *UserCaseScriptHandleRegister) Register(funcType string, scriptType string, handler IScriptHandler) {
	if register.routes == nil {
		register.routes = make(map[string]IScriptHandler)
	}

	if funcType != "" && scriptType != "" {
		route := fmt.Sprintf("%s:%s", funcType, scriptType)
		if register.routes[route] != nil {
			util.Logger.Warn("[CaseScriptHandleRegister] Script handler route: %s is regist, no need regist again!", route)
			return
		}
		register.routes[route] = handler
	}
}

func (register *UserCaseScriptHandleRegister) Trigger(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	route := fmt.Sprintf("%s:%s", funcCtx.FuncType, funcCtx.ScriptType)
	scriptHandler := register.routes[route]
	if funcCtx == nil || funcCtx.FuncBody == "" {
		log := fmt.Sprintf("[CaseScriptHandleRegister] User case: %d-%s not find function body", execCtx.Id, execCtx.Name)
		execCtx.AddLogs(log)
		util.Logger.Warn(CombineLogInfo(log, execCtx))
		return errors.New(log)
	}
	if scriptHandler == nil {
		log := fmt.Sprintf("[CaseScriptHandleRegister] Can not find script engine for user case: %s, route name: %s", execCtx.Name, route)
		execCtx.AddLogs(log)
		util.Logger.Warn(log)
		return errors.New(log)
	}
	err := scriptHandler.BuildScript(execCtx, funcCtx)
	if err != nil {
		return err
	}
	return scriptHandler.Execute(execCtx, funcCtx)
}

func loadScript(fileName string) (fileBody string, err error) {
	basePath := os.Getenv("LUA_PATH")
	if basePath == "" {
		basePath = "."
	}

	file, err := os.Open(fmt.Sprintf("%s/lua/%s.lua", basePath, fileName))
	if err != nil {
		return "", errors.New(fmt.Sprintf("[CaseScriptHandleRegister] Open file err :%s", err.Error()))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprintf("[CaseScriptHandleRegister] Close file error :%s", err.Error()))
		}
	}(file)
	var buf [128]byte
	var content []byte
	for {
		n, err := file.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", errors.New(fmt.Sprintf("[CaseScriptHandleRegister] Read file err: %s", err.Error()))
		}
		content = append(content, buf[:n]...)
	}

	util.Logger.Info("[CaseScriptHandleRegister] Load lua script from disk success,script name: %s", fileName)
	return string(content), nil
}

func buildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext, function string, functionName string) error {
	//注意两个问题:(1)同一个 JIT 如果函数名一样，函数体不一样，不同 case 的执行结果会产生偏差。
	//(2) 过多的脚本或者函数的加载会导致, JIT 无效函数的加载内容过多，占用内存。
	// 使用函数名+ case ID 确定函数名
	var log string
	if execCtx.CaseId != "" && execCtx.Name != "" && execCtx.TaskId != "" {
		funcName := fmt.Sprintf("%s_%s_%s", functionName, execCtx.CaseId, execCtx.TaskId)
		funcScript := strings.Replace(function, "@functionName", funcName, -1)
		funcScript = strings.Replace(funcScript, "@funcBody", funcCtx.FuncBody, -1)
		funcCtx.FuncBody = funcScript
		log = fmt.Sprintf("[CaseScriptHandleRegister] Current user case script build success,function name: %s,scripts: %s", funcName, funcCtx.FuncBody)
		util.Logger.Info(log)
		execCtx.AddLogs(log)
		funcCtx.FuncName = funcName
		return nil
	}

	log = fmt.Sprintf("[CaseScriptHandleRegister] Function context is not valid, lost parameters!")
	execCtx.AddLogs(log)
	return errors.New(log)
}

func scriptExecute(funcName string, execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	jIT, err := GetLuaJITWithRetry(execCtx, funcName)
	if err != nil {
		return err
	}
	defer DestroyedLuaPool(jIT)
	log := fmt.Sprintf("[CaseScriptHandleRegister] Begin to execute Function name: %s, script: %s, parameters: %s", funcCtx.FuncName, funcCtx.FuncBody, execCtx.GetStringVariables())
	execCtx.AddLogs(log)
	util.Logger.Info(CombineLogInfo(log, execCtx))
	err = jIT.DoString(funcCtx.FuncBody)
	if err != nil {
		log = fmt.Sprintf("[CaseScriptHandleRegister] Load function error, compile fail, err: %s!, function name: %s, function script: %s", err.Error(), funcCtx.FuncName, funcCtx.FuncBody)
		execCtx.LastErrorScript = funcCtx.FuncBody
		execCtx.AddLogs(log)
		util.Logger.Error(log)
		execCtx.Stop()
		return err
	}

	table := prepareLuaCtx(execCtx)
	timeCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	jIT.SetContext(timeCtx)
	if err := jIT.CallByParam(lua.P{
		Fn:      jIT.GetGlobal(funcCtx.FuncName),
		NRet:    1,
		Protect: true,
	}, table); err != nil {
		log = fmt.Sprintf("[CaseScriptHandleRegister] Execute fail, error: %s, function name: %s, script: %s  , parameters: %s", err.Error(), funcCtx.FuncName, funcCtx.FuncBody, execCtx.GetStringVariables())
		execCtx.LastErrorScript = funcCtx.FuncBody
		execCtx.AddLogs(log)
		util.Logger.Warn(CombineLogInfo(log, execCtx))
		execCtx.SetStop()
		return err
	}
	ret := jIT.Get(-1)
	refreshLuaCtx(ret, execCtx, jIT)
	jIT.Pop(1)
	return err
}

func scriptExecuteAfterResp(funcName string, execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	jIT, err := GetLuaJITWithRetry(execCtx, funcName)
	if err != nil {
		return err
	}
	defer DestroyedLuaPool(jIT)
	log := fmt.Sprintf("[CaseScriptHandleRegister] Begin to execute Function name: %s, script: %s, parameters: %s", funcCtx.FuncName, funcCtx.FuncBody, execCtx.GetStringVariables())
	execCtx.AddLogs(log)
	util.Logger.Info(CombineLogInfo(log, execCtx))
	err = jIT.DoString(funcCtx.FuncBody)
	if err != nil {
		execCtx.LastErrorScript = funcCtx.FuncBody
		log = fmt.Sprintf("[CaseScriptHandleRegister] Load function error, compile fail, err: %s!, function name: %s, function script: %s", err.Error(), funcCtx.FuncName, funcCtx.FuncBody)
		execCtx.AddLogs(log)
		util.Logger.Error(log)
		execCtx.SetStop()
		return err
	}

	table := prepareLuaCtx(execCtx)
	statusCode := lua.LNumber(execCtx.RespCode)
	resp := lua.LString(execCtx.RespBody)
	timeCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	jIT.SetContext(timeCtx)
	if err := jIT.CallByParam(lua.P{
		Fn:      jIT.GetGlobal(funcCtx.FuncName),
		NRet:    1,
		Protect: true,
	}, table, resp, statusCode); err != nil {
		execCtx.LastErrorScript = funcCtx.FuncBody
		log = fmt.Sprintf("[CaseScriptHandleRegister] Execute fail, error: %s, function name: %s, script: %s  , parameters: %s", err.Error(), funcCtx.FuncName, funcCtx.FuncBody, execCtx.GetStringVariables())
		execCtx.AddLogs(log)
		util.Logger.Warn(CombineLogInfo(log, execCtx))
		execCtx.SetStop()
		return err
	}
	ret := jIT.Get(-1)
	refreshLuaCtx(ret, execCtx, jIT)
	jIT.Pop(1)
	return err
}
