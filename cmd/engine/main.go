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

package main

import (
	_ "auto-test-go/api/router"
	"auto-test-go/pkg"
	"auto-test-go/pkg/enum"
	"auto-test-go/pkg/script"
	"auto-test-go/pkg/util"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	prepare()
	exitHandler()
	beego.Run()
}

func prepare() {
	util.NacosHelper{}.RegisterServiceInstance()
	util.Init()
	var handler script.IScriptHandler

	handler, err := script.NewLuaScriptBaseHandler()
	if err != nil {
		panic(err)
	}
	script.CaseRegister.Register(enum.LuaFuncType_DoBaseUserCaseExecute, enum.ScriptType_LuaScript, handler)

	handler, err = script.NewLuaScriptDoHttpCallHandler()
	if err != nil {
		panic(err)
	}
	script.CaseRegister.Register(enum.LuaFuncType_DoHttpRequest, enum.ScriptType_LuaScript, handler)

	handler, err = script.NewLuaScriptSqlHandler()
	if err != nil {
		panic(err)
	}
	script.CaseRegister.Register(enum.LuaFuncType_DoSqlExecute, enum.ScriptType_LuaScript, handler)

	handler, err = script.NewLuaScriptAssertHandler()
	if err != nil {
		panic(err)
	}
	script.CaseRegister.Register(enum.LuaFuncType_AssertUserCase, enum.ScriptType_LuaScript, handler)

	handler, err = script.NewHttpProtocolExecuteHandler()
	if err != nil {
		panic(err)
	}
	script.CaseRegister.Register(enum.ProtocolTypeHttp_DoRequest, enum.ScriptType_HttpCall, handler)
	err = pkg.TaskDispatch.Init()
	if err != nil {
		panic(err)
	}

	util.Nacos_Helper.RegisterServiceInstance()
}

func exitHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
				fmt.Println("got signal and try to exit: ", s)
				fmt.Println("try do some clear jobs")
				util.Nacos_Helper.DeregisterInstance()
				fmt.Println("run done")
				os.Exit(0)
			case syscall.SIGUSR1:
				fmt.Println("usr1: ", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2: ", s)
			default:
				fmt.Println("other: ", s)
			}
		}
	}()
}
