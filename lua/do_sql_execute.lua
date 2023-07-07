--
-- Licensed to the Apache Software Foundation (ASF) under one or more
-- contributor license agreements.  See the NOTICE file distributed with
-- this work for additional information regarding copyright ownership.
-- The ASF licenses this file to You under the Apache License, Version 2.0
-- (the "License"); you may not use this file except in compliance with
-- the License.  You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.
--

local json = require("json")
local http = require("http")

function add_log(ctx, log)
     if not ctx then
         print(log)
        end

        if not ctx.inner_log then
             ctx.inner_log = ''
        end

        if log then
             ctx.inner_log = ctx.inner_log .. log .. '\n'
        end
 end

@commonFunctions

 function inner_function_@functionName(ctx)
     mysql = require('mysql')
     c = mysql.new()
     add_log(ctx ,'msql host: ' .. '@host'  .. ', port: ' .. '@port' .. ', dbName: ' .. '@dbName' .. ', userName: ' .. '@userName' .. ', password:' ..  '@password')
     ok, err = c:connect({ host = '@host', port = @port, database = '@dbName', user = '@userName', password = '@password'})
     if err then
        add_log(ctx ,'mysql connect error: ' .. err)
     end
     if ok then
       @funcBody
     end
     c:close()
 end

 function @functionName(ctx)
    if type(ctx) ~= "table" then
        ctx.add_log('input ctx is not a table, can not execute function')
        return
    end

   -- add_log(ctx,'script input ctx is:' .. json.encode(ctx))
    inner_function_@functionName(ctx)
    return ctx
end
