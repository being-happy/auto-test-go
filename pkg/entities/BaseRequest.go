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

type BaseRequest struct {
	Name       string            `json:"name"`
	Headers    map[string]string `json:"headers"`
	Protocol   string            `json:"protocol"`
	Method     string            `json:"method"`
	Body       string            `json:"body"`
	Url        string            `json:"url"`
	Timeout    int               `json:"timeout"`
	Parameters map[string]string `json:"parameters"`
}

func (r *BaseRequest) Validation() bool {
	if r.Url == "" || r.Method == "" {
		return false
	}
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	if r.Timeout == 0 {
		r.Timeout = 30
	}
	if r.Timeout > 60 {
		r.Timeout = 60
	}
	if r.Headers["Content-Type"] == "" {
		r.Headers["Content-Type"] = "application/json"
	}
	return true
}
