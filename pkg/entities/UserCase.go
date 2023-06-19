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

type UserCase struct {
	Name            string          `json:"name"`
	Id              int64           `json:"id"`
	Request         BaseRequest     `json:"request"`
	Parameters      []CaseParameter `json:"parameters"`
	PreScripts      BaseScripts     `json:"preScripts"`
	AfterScripts    BaseScripts     `json:"afterScripts"`
	Assert          LuaScript       `json:"assert"`
	IsSkipError     bool            `json:"isSkipError"`
	TextAsserts     []TextAssert    `json:"textAsserts"`
	DependFunctions []string        `json:"dependFunctions"`
}

type TextAssert struct {
	ResponseType string `json:"responseType"`
	Operation    string `json:"operation"`
	Data         string `json:"data"`
}
