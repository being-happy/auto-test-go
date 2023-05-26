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
	"bufio"
	"errors"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"os"
	"sync/atomic"
	"time"
)

var luaPoolMax int32 = 1000
var luaPoolCount int32

func newLuaScriptJIT() (luaState *lua.LState, err error) {
	if luaPoolCount > luaPoolMax {
		return nil, errors.New("[LuaScriptFacade] lua pool is over max limit, please try latter")
	}

	luaPool := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  1024 * 20,
	})

	atomic.AddInt32(&luaPoolCount, 1)
	return luaPool, nil
}

// CompileLua reads the pssed lua file from disk and compiles it.
func CompileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			util.Logger.Warn("[LuaScriptFacade] Close file path error, path: %s, error: %s!", filePath, err.Error())
		}
	}(file)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func DoCompiledFile(l *lua.LState, path string) error {
	err := l.DoFile(path)
	return err
}

func GetLuaJITWithRetry(execCtx *entities.ExecContext, funcName string) (*lua.LState, error) {
	jIT, err := newLuaScriptJIT()
	if err != nil {
		execCtx.AddLogs(err.Error())
		timer := time.NewTimer(5 * time.Second)
		i := 0
		for {
			i++
			<-timer.C
			jIT, err = newLuaScriptJIT()
			if err != nil {
				execCtx.AddLogs(fmt.Sprintf("[LuaScriptFacade] Retry %d times, but execute script also error,info :%s", i, err.Error()))
			} else {
				break
			}
			if i > 3 {
				return jIT, err
			}
		}
	}
	luaScriptWrapper.wrap(jIT, funcName)
	return jIT, err
}

func DestroyedLuaPool(state *lua.LState) {
	state.Close()
	if luaPoolCount <= 0 {
		return
	}
	atomic.AddInt32(&luaPoolCount, -1)
}
