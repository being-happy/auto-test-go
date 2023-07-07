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

type LuaScript struct {
	Script          string   `json:"luaScript"`
	FuncType        string   `json:"type"`
	Host            string   `json:"host"`
	Port            string   `json:"port"`
	UserName        string   `json:"userName"`
	Password        string   `json:"password"`
	DbName          string   `json:"dbName"`
	CallFunction    string   `json:"callFunction"`
	DependFunctions []string `json:"dependFunctions"`
}

type SqlAuth struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
}

func (sqlScript *LuaScript) ValidScript() bool {
	if sqlScript.UserName == "" || sqlScript.DbName == "" || sqlScript.Script == "" || sqlScript.Password == "" || sqlScript.Host == "" || sqlScript.Port == "" {
		return false
	}
	return true
}

func (sqlScript *LuaScript) CopySqlToFuncContext(ctx *FuncContext) {
	ctx.UserName = sqlScript.UserName
	ctx.Password = sqlScript.Password
	ctx.Host = sqlScript.Host
	ctx.Port = sqlScript.Port
	ctx.DbName = sqlScript.DbName
}

type HttpCall struct {
	IsVoid bool   `json:"isVoid"`
	Name   string `json:"name"`
	BaseRequest
	LuaScript LuaScript `json:"luaScript"`
}

type CaseParameter struct {
	Name   string     `json:"name"`
	Value  string     `json:"value"`
	PType  string     `json:"pType"`
	Script BaseScript `json:"script"`
}

type BaseScript struct {
	ScriptType  string    `json:"scriptType"`
	Script      LuaScript `json:"script"`
	Order       int32
	IsSkipError bool `json:"isSkipError"`
}

type BaseScripts []BaseScript

func (b BaseScripts) Len() int {
	return len(b)
}

func (b BaseScripts) Less(i, j int) bool {
	return b[i].Order < b[j].Order
}

func (b BaseScripts) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
