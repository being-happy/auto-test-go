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
	"auto-test-go/pkg"
	"auto-test-go/pkg/command"
	"auto-test-go/pkg/director"
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
	"auto-test-go/pkg/script"
	"auto-test-go/pkg/util"
	"encoding/json"
	"os"
)

func prepare() {
	util.Init()
	var handler script.IScriptHandler

	handler, err := script.NewLuaScriptBaseHandler()
	if err != nil {
		panic(err)
	}
	script.CaseRegister.Register(enum.LuaFuncType_DoBaseExecute, enum.ScriptType_LuaScript, handler)

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
}

func main() {
	prepare()
	scenariorTest()
	os.Setenv("NACOS_ADDRESS", "127.0.0.1")
	os.Setenv("NACOS_PORT", "80")
	os.Setenv("GROUP_NAME", "03a1c325-7c9b-41bf-b4f6-404a0cf22d5a")
	//util.NacosHelper{}.RegisterServiceInstance()
	//	batchUserCaseTest()
	//scriptDebuggerTest()
}

func scriptDebuggerTest() {
	execCommand := command.SingleScriptExecuteCommand{
		Name: "scriptLady",
		Id:   "1",
		Parameters: []entities.CaseParameter{
			{
				Name:  "id",
				Value: "1",
			}, {
				Name:  "name",
				Value: "zhangSan",
			}, {
				Name:  "age",
				Value: "15",
			}, {
				Name:  "host",
				Value: "127.0.0.1",
			},
		},
		BaseScript: entities.BaseScript{
			ScriptType: enum.ScriptType_LuaScript,
			Script: entities.LuaScript{
				FuncType: enum.LuaFuncType_DoBaseExecute,
				Script:   "ctx.name='wangWu'",
			},
		},
	}

	bytes, _ := json.Marshal(execCommand)
	str := string(bytes)
	util.Logger.Warn(str)
	factory := director.BaseDirectorFactory{}.Create(enum.DirectorType_ScriptDebugger)
	factory.Action(&execCommand, false)
}

func userCaseTest() {
	execCommand := command.SingleUserCaseExecuteCommand{
		Name:        "test command",
		Id:          "1",
		CommandType: enum.CommandType_Execute,
		UserCase: entities.UserCase{
			Name: "preside",
			Id:   1,
			Parameters: []entities.CaseParameter{
				{
					Name:  "id",
					Value: "1",
				}, {
					Name:  "name",
					Value: "zhangSan",
				}, {
					Name:  "age",
					Value: "15",
				}, {
					Name:  "host",
					Value: "127.0.0.1",
				},
			},
			PreScripts: entities.BaseScripts([]entities.BaseScript{
				{
					ScriptType: enum.ScriptType_LuaScript,
					Order:      2,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoBaseExecute,
						Script: "print(ctx.name) \n " +
							" print(ctx.message)",
					},
				}, {
					ScriptType: enum.ScriptType_LuaScript,
					Order:      1,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoBaseExecute,
						Script:   "ctx.name='wangWu'",
					},
				},
				{
					ScriptType: enum.ScriptType_LuaScript,
					Order:      0,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoHttpRequest,
						Script: "  local response, error_message = http.get(\"https://127.0.0.1/mixmicro-demo/mock/latbox-test\")\n  " +
							"if error_message  then\n    " +
							"   print(\"http request call fail:\" .. error_message)\n    end\n  " +
							"  if response.body ~=nil and response.status_code==200 then\n    " +
							"   local my_body = json.decode(response.body)\n      " +
							"   ctx.message = my_body.message\n   " +
							" else\n     " +
							"   print(\"response code is :\" .. response.status_code .. \", url: https://127.0.0.1/mixmicro-demo/mock/latbox-test\")\n  " +
							" end",
					},
				}, {
					ScriptType: enum.ScriptType_SqlScript,
					Order:      4,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoSqlExecute,
						DbName:   "ops-bate",
						Host:     "127.0.0.1",
						Port:     "3306",
						UserName: "userName",
						Password: "password",
						Script: "  local resp, err = c:query(\"SELECT * FROM service_config WHERE  service_name='mixmicro-ops-api'\")\n " +
							"  if err then\n     " +
							" add_log(ctx, err)\n   " +
							"else\n     " +
							"    for i = 1, #resp do\n   " +
							" ctx.rows = json.encode(resp[i])" +
							"   add_log(ctx,'sql: ' .. json.encode(resp[i]))  \n " +
							"  end \n" +
							"end",
					},
				}, {
					ScriptType: enum.ScriptType_SqlScript,
					Order:      5,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoSqlExecute,
						DbName:   "ops-bata",
						Host:     "127.0.0.1",
						Port:     "3306",
						UserName: "userName",
						Password: "password",
						Script: "  local resp, err = c:query(\"update service_config set project_id =1  WHERE  id = 191\")\n " +
							"  if err then\n     " +
							" add_log(ctx, err)\n   " +
							"else\n     " +
							"    for i = 1, #resp do\n   " +
							" ctx.rows = json.encode(resp[i])" +
							"   add_log(ctx,'sql: ' .. json.encode(resp[i]))  \n " +
							"  end \n" +
							"end",
					},
				},
			}),
			Assert: entities.LuaScript{
				FuncType: enum.LuaFuncType_AssertUserCase,
				Script: "add_log(ctx,'http response code: ' .. code) \n " +
					"add_log(ctx,'http response data: ' .. resp) \n" +
					"ctx.stop=false",
			},
			AfterScripts: entities.BaseScripts([]entities.BaseScript{
				{
					ScriptType: enum.ScriptType_LuaScript,
					Order:      1,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoBaseExecute,
						Script:   "add_log(ctx,'http response by after script:')",
					},
				},
			}),
			Request: entities.BaseRequest{
				Url:     "https://@host/mixmicro-demo/mock/latbox-test",
				Method:  "GET",
				Timeout: 30,
			},
		},
	}
	str, _ := json.Marshal(execCommand)
	util.Logger.Warn(string(str))
	factory := director.BaseDirectorFactory{}.Create(enum.DirectorType_UserCase)
	_, err := factory.Action(&execCommand, false)
	if err != nil {
		return
	}
}

func batchUserCaseTest() {
	execCommand := command.UserCaseBatchExecuteCommand{
		Name:        "test command",
		Id:          "1",
		CommandType: enum.CommandType_Execute,
		UserCases: []entities.UserCase{
			entities.UserCase{
				Name: "preside",
				Id:   1,
				Parameters: []entities.CaseParameter{
					{
						Name:  "id",
						Value: "1",
					}, {
						Name:  "name",
						Value: "zhangSan",
					}, {
						Name:  "age",
						Value: "15",
					}, {
						Name:  "host",
						Value: "127.0.0.1",
					},
				},
				PreScripts: entities.BaseScripts([]entities.BaseScript{
					{
						ScriptType: enum.ScriptType_LuaScript,
						Order:      2,
						Script: entities.LuaScript{
							FuncType: enum.LuaFuncType_DoBaseExecute,
							Script: "print(ctx.name) \n " +
								" print(ctx.message)",
						},
					}, {
						ScriptType: enum.ScriptType_LuaScript,
						Order:      1,
						Script: entities.LuaScript{
							FuncType: enum.LuaFuncType_DoBaseExecute,
							Script:   "ctx.name='wangWu'",
						},
					},
					{
						ScriptType: enum.ScriptType_LuaScript,
						Order:      0,
						Script: entities.LuaScript{
							FuncType: enum.LuaFuncType_DoHttpRequest,
							Script: "  local response, error_message = http.get(\"https://127.0.0.1/mixmicro-demo/mock/latbox-test\")\n  " +
								"if error_message  then\n    " +
								"   print(\"http request call fail:\" .. error_message)\n    end\n  " +
								"  if response.body ~=nil and response.status_code==200 then\n    " +
								"   local my_body = json.decode(response.body)\n      " +
								"   ctx.message = my_body.message\n   " +
								" else\n     " +
								"   print(\"response code is :\" .. response.status_code .. \", url: https://127.0.0.1/mixmicro-demo/mock/latbox-test\")\n  " +
								" end",
						},
					}, {
						ScriptType: enum.ScriptType_SqlScript,
						Order:      4,
						Script: entities.LuaScript{
							FuncType: enum.LuaFuncType_DoSqlExecute,
							DbName:   "ops-bata",
							Host:     "127.0.0.1",
							Port:     "3306",
							UserName: "userName",
							Password: "password",
							Script: "  local resp, err = c:query(\"SELECT * FROM service_config WHERE  service_name='mixmicro-ops-api'\")\n " +
								"  if err then\n     " +
								" add_log(ctx, err)\n   " +
								"else\n     " +
								"    for i = 1, #resp do\n   " +
								" ctx.rows = json.encode(resp[i])" +
								"   add_log(ctx,'sql: ' .. json.encode(resp[i]))  \n " +
								"  end \n" +
								"end",
						},
					}, {
						ScriptType: enum.ScriptType_SqlScript,
						Order:      5,
						Script: entities.LuaScript{
							FuncType: enum.LuaFuncType_DoSqlExecute,
							DbName:   "ops-bata",
							Host:     "127.0.0.1",
							Port:     "3306",
							UserName: "userName",
							Password: "password",
							Script: "  local resp, err = c:query(\"update service_config set project_id =1  WHERE  id = 191\")\n " +
								"  if err then\n     " +
								" add_log(ctx, err)\n   " +
								"else\n     " +
								"    for i = 1, #resp do\n   " +
								" ctx.rows = json.encode(resp[i])" +
								"   add_log(ctx,'sql: ' .. json.encode(resp[i]))  \n " +
								"  end \n" +
								"end",
						},
					},
				}),
				//Assert: entities.LuaScript{
				//	FuncType: enum.LuaFuncType_AssertUserCase,
				//	Script: "add_log(ctx,'http response code: ' .. code) \n " +
				//		"add_log(ctx,'http response data: ' .. resp) \n" +
				//		"ctx.stop=false",
				//},
				TextAsserts: []entities.TextAssert{
					{
						ResponseType: enum.TextAssert_ResponseCode,
						Data:         "200",
						Operation:    enum.OperationType_Contains,
					},
				},
				AfterScripts: entities.BaseScripts([]entities.BaseScript{
					{
						ScriptType: enum.ScriptType_LuaScript,
						Order:      1,
						Script: entities.LuaScript{
							FuncType: enum.LuaFuncType_DoBaseExecute,
							Script:   "add_log(ctx,'http response by after script:')",
						},
					},
				}),
				Request: entities.BaseRequest{
					Url:     "https://@host/mixmicro-demo/mock/latbox-test",
					Method:  "GET",
					Timeout: 30,
				},
			},
		},
	}

	byte1, _ := json.Marshal(execCommand)
	str := string(byte1)
	util.Logger.Warn(str)
	factory := director.BaseDirectorFactory{}.Create(enum.DirectorType_BatchUserCase)
	_, err := factory.Action(&execCommand, false)
	if err != nil {
		return
	}
}

func scenariorTest() {
	factory := director.BaseDirectorFactory{}.Create(enum.DirectorType_ScenariorCase)
	var execCommand = command.ScenarioCaseExecuteCommand{
		Name:        "test command",
		Id:          "1",
		CommandType: enum.CommandType_Execute,
		ScenarioCase: entities.ScenarioCase{
			Name: "preside",
			Id:   1,
			Parameters: []entities.CaseParameter{
				{
					Name:  "id",
					Value: "1",
				}, {
					Name:  "name",
					Value: "zhangSan",
				}, {
					Name:  "age",
					Value: "15",
				}, {
					Name:  "host",
					Value: "127.0.0.1",
				},
			},
			PreScripts: entities.BaseScripts([]entities.BaseScript{
				{
					ScriptType: enum.ScriptType_LuaScript,
					Order:      2,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoBaseExecute,
						Script: "print(ctx.name) \n " +
							" print(ctx.message)",
					},
				}}),
			AfterScripts: entities.BaseScripts([]entities.BaseScript{
				{
					ScriptType: enum.ScriptType_LuaScript,
					Order:      1,
					Script: entities.LuaScript{
						FuncType: enum.LuaFuncType_DoBaseExecute,
						Script:   "add_log(ctx,'http response by after script:')",
					},
				},
			}),
			UserCases: map[string]*entities.UserCase{"1": {
				Name: "preside",
				Id:   1,
				Parameters: []entities.CaseParameter{
					{
						Name:  "id",
						Value: "1",
					}, {
						Name:  "name",
						Value: "zhangSan",
					}, {
						Name:  "age",
						Value: "15",
					}, {
						Name:  "host",
						Value: "127.0.0.1",
					},
				},
				Assert: entities.LuaScript{
					FuncType: enum.LuaFuncType_AssertUserCase,
					Script: "add_log(ctx,'http response code: ' .. code) \n " +
						"add_log(ctx,'http response data: ' .. resp) \n" +
						"ctx.stop=false",
				},
				Request: entities.BaseRequest{
					Url:     "https://@host/mixmicro-demo/mock/latbox-test",
					Method:  "GET",
					Timeout: 30,
				},
			},
				"2": {
					Name: "preside",
					Id:   2,
					Parameters: []entities.CaseParameter{
						{
							Name:  "id",
							Value: "1",
						}, {
							Name:  "name",
							Value: "zhangSan",
						}, {
							Name:  "age",
							Value: "15",
						}, {
							Name:  "host",
							Value: "127.0.0.1",
						},
					},
					Assert: entities.LuaScript{
						FuncType: enum.LuaFuncType_AssertUserCase,
						Script: "add_log(ctx,'http response code: ' .. code) \n " +
							"add_log(ctx,'http response data: ' .. resp) \n" +
							"ctx.stop=false",
					},
					Request: entities.BaseRequest{
						Url:     "https://@host/mixmicro-demo/mock/latbox-test",
						Method:  "GET",
						Timeout: 30,
					},
				}},
			Scripts: map[string]*entities.BaseScript{"1": &entities.BaseScript{
				ScriptType: enum.ScriptType_LuaScript,
				Script: entities.LuaScript{
					FuncType: enum.LuaFuncType_DoBaseExecute,
					Script:   "ctx.name='wangWu'",
				},
			}},
			Baggages: map[string]*entities.Baggage{
				"1": &entities.Baggage{
					Data: []string{"['name1','age1','sex']", "['name2','age2','sex1']"},
				},
			},
			Design: "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<scenario-flow>\n   " +
				" <user-case case-id=\"1\" case-name=\"GET http://127.0.0.1/mixmicro-ops-api/v1/workbench/buildStatByGroupMonth?status==&month=6 [55700]\" order=\"1\" />\n   " +
				" <user-case case-id=\"2\" case-name=\"helloword\" order=\"2\" />\n  " +
				"  <time-wait-case wait-time=\"1\" order =\"3\" />\n   " +
				" <script-case script-id =\"1\" script-type =\"ScriptType_LuaScript\" order=\"4\"/>\n   " +
				" <loop-script-case loop-type=\"LoopType_Data\" data-id=\"1\" loop-id=\"2\"> \n       " +
				" <user-case case-id=\"1\" case-name=\"helloword\" order=\"1\" />\n    " +
				"    <user-case case-id=\"2\" case-name=\"helloword\" order=\"2\" />\n     " +
				"   <time-wait-case wait-time=\"1\" order =\"3\" />\n       " +
				" <loop-script-case loop-type=\"LoopType_Data\" data-id=\"1\" loop-count=\"\" order=\"\" loop-id=\"3\"> \n        " +
				"    <user-case case-id=\"1\" case-name=\"helloword\" order=\"1\" />\n         " +
				"   <user-case case-id=\"2\" case-name=\"helloword\" order=\"2\" />\n          " +
				"  <time-wait-case wait-time=\"1\" order =\"3\"/>\n      " +
				"  </loop-script-case>\n   " +
				" </loop-script-case>\n</scenario-flow>\n",
		},
	}
	str, _ := json.Marshal(execCommand)
	util.Logger.Warn(string(str))
	ctx, err := factory.Action(&execCommand, false)
	if err != nil {
		return
	}
	if ctx != nil {

	}
}
