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
	"auto-test-go/pkg/util"
	"encoding/json"
	"fmt"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
	"net/url"
	"strings"
	"time"
)

type HttpProtocolExecuteHandler struct {
	BaseScripHandler
}

func NewHttpProtocolExecuteHandler() *HttpProtocolExecuteHandler {
	handler := HttpProtocolExecuteHandler{}
	handler.Name = enum.ProtocolTypeHttp_DoRequest
	handler.ScriptType = enum.ScriptType_HttpCall
	handler.FuncType = enum.ProtocolTypeHttp_DoRequest
	err := handler.Init()
	if err != nil {
		panic(err)
	}
	return &handler
}

func (l *HttpProtocolExecuteHandler) Init() error {
	return nil
}

func (l *HttpProtocolExecuteHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	request := funcCtx.Request
	reqUrl, _ := url.QueryUnescape(request.Url)
	if !strings.Contains(reqUrl, "?") {
		//解析path上变量
		reqUrl = l.replaceUrlParameters(reqUrl, execCtx.Variables, false) + "?"
	} else {
		//针对已有URL进行编码
		arrys := strings.Split(reqUrl, "?")
		if len(arrys) > 1 {
			//解析path上变量
			url := l.replaceUrlParameters(arrys[0], execCtx.Variables, false) + "?"
			//解析&符号 例如: 解析前: a=1&b=2&c=3 解析后： [0]:a=1;[1]b=2;[3]c=3
			kvs := strings.Split(arrys[1], "&")
			if len(kvs) > 0 {
				//解析=符号
				for i := 0; i < len(kvs); i++ {
					//解析=符号 例如：解析前 a=1 解析后: [0]:a;[1]:1
					kvArray := strings.Split(kvs[i], "=")
					url = url + kvArray[0] + "=" + l.replaceUrlParameters(kvArray[1], execCtx.Variables, true) + "&"
				}
			}
			reqUrl = url
		}
	}

	// url 的变量必须单独替换，否则会引发 URL 编码问题
	if request.Parameters != nil && len(request.Parameters) > 0 {
		for k, v := range request.Parameters {
			reqUrl = reqUrl + k + "=" + l.replaceUrlParameters(v, execCtx.Variables, true) + "&"
		}
	}

	request.Url = reqUrl
	if len(execCtx.Variables) > 0 {
		fastTemplate := FastTemplate{}
		//优先替换所有的花括号
		vars := fastTemplate.convertVar(execCtx.Variables)
		for hk, hv := range request.Headers {
			//优先根据变量进行匹配,若变量不匹配则根据名称匹配，若名称相同则替换
			request.Headers[hk] = fastTemplate.template(hv, vars)
		}
		//替换请求体当中的变量
		request.Body = fastTemplate.template(request.Body, vars)
	}

	//解码 url 忽略错误
	//if strings.Contains(baseRequest.Url, "?") {
	//	arrys := strings.Split(baseRequest.Url, "?")
	//	if len(arrys) > 1 {
	//		urlParameter, _ := url.QueryUnescape(arrys[1])
	//		urlParameter = url.QueryEscape(urlParameter)
	//		baseRequest.Url = arrys[0] + "?" + urlParameter
	//	}
	//}
	funcCtx.Request = request
	execCtx.CurrentRequest = request
	return nil
}

func (l *HttpProtocolExecuteHandler) replaceUrlParameters(dist string, variables map[string]entities.VarValue, encode bool) string {
	if dist == "" {
		return dist
	}

	if len(variables) > 0 {
		fastTemplate := FastTemplate{}
		dist = fastTemplate.template(dist, fastTemplate.convertVar(variables))
	}

	if encode {
		return url.QueryEscape(dist)
	}
	return dist
}

func (l *HttpProtocolExecuteHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	str, _ := json.Marshal(funcCtx.Request)
	log := fmt.Sprintf("[HttpProtocolExecuteHandler] Begin to execute http request, base request info is: %s", string(str))
	execCtx.AddLogs(log)
	util.Logger.Info(CombineLogInfo(log, execCtx))
	baseRequest := funcCtx.Request
	cli := gentleman.New()
	cli.URL(baseRequest.Url)
	req := cli.Request()
	cli.Use(timeout.Request(time.Duration(baseRequest.Timeout) * time.Second))
	req.Method(baseRequest.Method)
	req.SetHeaders(baseRequest.Headers)
	req.BodyString(baseRequest.Body)
	resp, err := req.Send()
	if err != nil {
		log = fmt.Sprintf("[HttpProtocolExecuteHandler] Execute http request error: %s, url: %s", err.Error(), baseRequest.Url)
		execCtx.AddLogs(log)
		util.Logger.Error(CombineLogInfo(log, execCtx))
		execCtx.RespBody = ""
		execCtx.RespCode = 600
		return err
	}

	execCtx.RespCode = resp.StatusCode
	execCtx.RespBody = resp.String()
	log = fmt.Sprintf("[HttpProtocolExecuteHandler] Execute http request success, url: %s , code: %d, body: %s.", baseRequest.Url, execCtx.RespCode, execCtx.RespBody)
	execCtx.AddLogs(log)
	util.Logger.Info(CombineLogInfo(log, execCtx))
	return err
}

//func (l *HttpProtocolExecuteHandler) template(source string, keyWord string, value string) string {
//	arg0 := "{@" + keyWord + "}"
//	arg1 := "@" + keyWord
//	if strings.Contains(source, arg0) {
//		source = strings.Replace(source, arg0, value, -1)
//	}
//
//	//兼容老的替换方式
//	if strings.Contains(source, arg1) {
//		source = strings.Replace(source, arg1, value, -1)
//	}
//	return source
//}
