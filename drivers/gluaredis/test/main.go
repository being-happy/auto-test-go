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
	"auto-test-go/drivers/gluaredis"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("redis", gluaredis.Loader)
	script := fmt.Sprintf(` 
           redis = require("redis");
           local c = redis.new();
           local ok, err =c:connect({ host = "%s" , port = "%s",database = "%d", password = "%s" });
           if err then
               print(err);
           end

           local ok, err = c:set_key("name","zhangshan");
           if err then
				print(err);
           end

		   local value, err = c:get_key("name");
           print(value);
           local success, err = c:del_key("name");
           if err then
              print(err);
           end

           local ok = c:close();
           if ok then
               print('close success');
           end
`, "127.0.0.1", "6379", 11, "123456")
	err := L.DoString(script)
	if err != nil {
		fmt.Println(err.Error())
	}
}
