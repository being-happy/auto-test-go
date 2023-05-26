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

package enum

const (
	Resp_OK                           int = 200
	Resp_Forbid                       int = 403
	Resp_Error                        int = 500
	Resp_Invalid                      int = 400
	ProtocolTypeHttp_DoRequest            = "ProtocolTypeHttp_DoRequest"
	LuaFuncType_DoHttpRequest             = "LuaFuncType_DoHttpRequest"
	LuaFuncType_DoSqlExecute              = "LuaFuncType_DoSqlExecute"
	LuaFuncType_AssertUserCase            = "LuaFuncType_AssertUserCase"
	LuaFuncType_DoBaseUserCaseExecute     = "LuaFuncType_DoBaseExecute"
	LuaFuncName_DoHttpRequest             = "do_http_request_execute"
	LuaFuncName_DoSqlExecute              = "do_sql_execute"
	LuaFuncName_AssertUserCase            = "assert_user_case_execute"
	LuaFuncName_DoBaseExecute             = "do_base_execute"
	ScriptType_LuaScript                  = "ScriptType_LuaScript"
	ScriptType_SqlScript                  = "ScriptType_SqlScript"
	ScriptType_HttpCall                   = "ScriptType_HttpCall"
	DirectorType_UserCase                 = "DirectorType_UserCase"
	CommandType_Execute                   = "Execute"
	CommandType_Stop                      = "Stop"
	DirectorType_ScriptDebugger           = "DirectorType_ScriptDebugger"
	DirectorType_BatchUserCase            = "DirectorType_BatchUserCase"
	DirectorType_ScenariorCase            = "DirectorType_ScenariorCase"
	LoopType_Data                         = "LoopType_Data"
	LoopType_Normal                       = "LoopType_Normal"
	TextAssert_ResponseCode               = "TextAssert_ResponseCode"
	TextAssert_ResponseData               = "TextAssert_ResponseData"
	TextAssert_ResponseHeaders            = "TextAssert_ResponseHeaders"
	OperationType_Contains                = "OperationType_Contains"
	OperationType_NoContains              = "OperationType_NoContains"
	OperationType_Equals                  = "OperationType_Equals"
	OperationType_StartWith               = "OperationType_StartWith"
	OperationType_EndWith                 = "OperationType_EndWith"
)
