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
	"auto-test-go/pkg/enum"
	"github.com/cjoudrey/gluahttp"
	json "github.com/layeh/gopher-json"
	mysql "github.com/tengattack/gluasql/mysql"
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

type LuaScriptWrapper struct {
	BaseScripHandler
	FileNames map[string]string
	isInit    bool
}

var luaScriptWrapper = &LuaScriptWrapper{}

func (l *LuaScriptWrapper) init() {
	/*l.FileNames = map[string]string{"json": "./lua/json.lua"}*/ /*
		for moduleName, pathName := range l.FileNames {
			proto, err := CompileLua(pathName)
			if err != nil {
				util.Logger.Errorf("[LuaScriptDecorator] Load common file error : %s, file path: %s", err.Error(), pathName)
			}

			if l.functionProto == nil {
				l.functionProto = map[string]*lua.FunctionProto{}
			}
			l.functionProto[moduleName] = proto
		}*/
}

func (l *LuaScriptWrapper) wrap(state *lua.LState, funcName string) {
	switch funcName {
	case enum.LuaFuncName_DoSqlExecute:
		state.PreloadModule("json", json.Loader)
		state.PreloadModule("mysql", mysql.Loader)
		state.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
		break
	default:
		state.PreloadModule("json", json.Loader)
		state.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	}
}
