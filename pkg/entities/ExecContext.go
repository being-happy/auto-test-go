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

package entities

import (
	"auto-test-go/pkg/util"
	"encoding/json"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"strings"
	"time"
)

type VarValue struct {
	Value  string      `json:"value"`
	Object interface{} `json:"object"`
	Type   string      `json:"type"`
}

const (
	Success = 1
	Failed  = 3
	Started = 0
	Pending = -1
)

type ExecContext struct {
	ParentId        int                 `json:"parentId"`
	Variables       map[string]VarValue `json:"variables"`
	Name            string              `json:"name"`
	CaseId          string              `json:"caseId"`
	Id              int                 `json:"id"`
	TaskId          string              `json:"taskId"`
	Logs            []string            `json:"logs"`
	stop            bool
	RespBody        string      `json:"respBody"`
	RespCode        int         `json:"respCode"`
	LastErrorScript string      `json:"lastErrorScript"`
	CurrentRequest  BaseRequest `json:"currentRequest"`
	IgnoreStop      bool        `json:"ignoreStop"`
	Status          int         `json:"status"`
	Start           time.Time   `json:"start"`
	Duration        int64       `json:"duration"`
	AssertSuccess   bool        `json:"assertSuccess"`
}

func (ctx *ExecContext) GetStringVariables() string {
	if s, err := json.Marshal(ctx.Variables); err != nil {
		util.Logger.Warn("[ExecContext] Marshal variables error :" + err.Error())
		return ""
	} else {
		return string(s)
	}
}

func (ctx *ExecContext) Close() {
	ctx.Logs = nil
}

func (ctx *ExecContext) AddLogs(log string) {
	if ctx.Logs == nil {
		ctx.Logs = []string{}
	}
	log = time.Now().Format("2006-01-02 15:04:05") + " " + log
	ctx.Logs = append(ctx.Logs, log)
}

func (ctx *ExecContext) ReadLogs() *[]string {
	return &ctx.Logs
}

func (ctx *ExecContext) Refresh(jIT *lua.LState, table *lua.LTable) {
	jIT.ForEach(table, func(propertyName lua.LValue, propertyValue lua.LValue) {
		if strings.ToLower(propertyName.String()) == "inner_log" {
			ctx.AddLogs("[ScriptInnerLog]" + propertyValue.String())
			return
		}

		if strings.ToLower(propertyName.String()) == "assert" {
			if strings.ToLower(propertyValue.String()) == "false" {
				//断言失败，将状态设置为失败，并且继续执行后置脚本
				ctx.AssertSuccess = false
			} else {
				ctx.AssertSuccess = true
			}
		}

		ctx.Variables[propertyName.String()] = VarValue{Value: propertyValue.String(), Type: propertyValue.Type().String(), Object: propertyValue}
		if propertyValue.Type().String() == "table" {
			goType := map[string]interface{}{}
			luaData := propertyValue.(*lua.LTable)
			err := gluamapper.Map(luaData, &goType)
			if err != nil {
				ctx.AddLogs("change table type to go type error: " + err.Error())
				return
			}

			value, err := json.Marshal(goType)
			if err != nil {
				return
			}

			ctx.Variables[propertyName.String()] = VarValue{Value: string(value), Type: propertyValue.Type().String(), Object: propertyValue}
		}
	})
}
func (ctx *ExecContext) Stop() bool {
	return ctx.stop
}

func (ctx *ExecContext) DeepCopy() *ExecContext {
	str, err := json.Marshal(ctx)
	if err != nil {
		return nil
	}

	copyCtx := &ExecContext{}
	json.Unmarshal(str, copyCtx)
	return copyCtx
}

func (ctx *ExecContext) SetStop() {
	if ctx.IgnoreStop {
		ctx.stop = false
	} else {
		ctx.stop = true
		ctx.SetStatus(Failed)
	}
}

func (ctx *ExecContext) Merge(ctx2 *ExecContext) {
	if len(ctx2.Variables) > 0 {
		for k, v := range ctx2.Variables {
			//特殊变量不与合并
			if strings.Contains(k, "assert") {
				continue
			}
			ctx.Variables[k] = v
		}
	}
	//ctx.Id = ctx2.Id + 1
}

func (ctx *ExecContext) Copy() *ExecContext {
	return &ExecContext{
		Variables: ctx.Variables,
		Start:     time.Now(),
	}
}

func (ctx *ExecContext) Reset() {
	ctx.stop = false
	ctx.SetStatus(Started)
	ctx.IgnoreStop = false
	ctx.Start = time.Now()
	ctx.AssertSuccess = true
}
func (ctx *ExecContext) SetStatus(status int) {
	ctx.Duration = time.Since(ctx.Start).Milliseconds()
	ctx.Status = status
}
func (ctx *ExecContext) GetStatus() int {
	return ctx.Status
}

type FuncContext struct {
	FuncBody   string
	FuncType   string
	ScriptType string
	FuncName   string
	SqlAuth
	Request BaseRequest
}
