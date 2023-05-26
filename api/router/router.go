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

package router

import (
	controllers "auto-test-go/api/controller"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/user-case/", &controllers.UserCaseController{}, "post:Receive;delete:Delete;get:GetUserCase"),
		beego.NSRouter("/user-case/command/receive-single-sync", &controllers.UserCaseController{}, "post:ReceiveSingleSync"),
		beego.NSRouter("/script-debugger/command/receive-single-sync", &controllers.ScriptDebuggerController{}, "post:ReceiveSingleSync"),
		beego.NSRouter("/user-case/command/receive-batch-sync", &controllers.UserCaseController{}, "post:ReceiveBatchSync"),
		beego.NSRouter("/scenario-case/command/receive-single-sync", &controllers.ScenarioController{}, "post:ReceiveSingelSync"),
		beego.NSRouter("/scenario-case/command/query-single-async/:taskId", &controllers.ScenarioController{}, "get:QueryCommand"),
		beego.NSRouter("/user-case/command/query-single-async/:taskId", &controllers.UserCaseController{}, "get:QueryCommand"),
		beego.NSRouter("/scenario-case/command/receive-single-async", &controllers.ScenarioController{}, "post:ReceiveSingelAsync"),
		beego.NSRouter("/user-case/command/receive-single-async", &controllers.UserCaseController{}, "post:ReceiveSingleAsync"),
		beego.NSRouter("/statistics/query", &controllers.ScenarioController{}, "get:QueryStatus"),
	)
	beego.AddNamespace(ns)
}
