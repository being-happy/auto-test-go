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

package gluaredis

import (
	"auto-test-go/drivers"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	lua "github.com/yuin/gopher-lua"
	"time"
)

const (
	CLIENT_TYPENAME = "redis{client}"
)

// Client mysql
type Client struct {
	DB      *redis.Client
	Timeout time.Duration
	Ctx     context.Context
}

var clientMethods = map[string]lua.LGFunction{
	"connect":     clientConnectMethod,
	"close":       clientCloseMethod,
	"get_key":     clientGetStrMethod,
	"set_key":     clientSetStrMethod,
	"del_key":     clientDelMethod,
	"set_timeout": clientSetTimeoutMethod,
}

var ErrConnectionString = errors.New("options or connection string excepted")

// parseConnectionString parse options or connection string to golang redis driverName and dsn
func parseConnectionString(cs interface{}, timeout time.Duration) (dsn *redis.Options, err error) {
	if cs == nil {
		return nil, ErrConnectionString
	} else if options, ok := cs.(map[string]interface{}); ok {
		host, _ := options["host"].(string)
		if host == "" {
			host = "127.0.0.1"
		}

		port, _ := options["port"].(int)
		if port == 0 {
			port = 6379
		}
		database, _ := options["database"].(int)
		password, _ := options["password"].(string)

		if timeout == 0 {
			timeout = 10 * time.Second
		}

		dsn = &redis.Options{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Password:     password,
			DB:           database,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
			DialTimeout:  timeout,
		}
		return
	}
	return nil, ErrConnectionString
}

func clientConnectMethod(L *lua.LState) int {
	if err := recover(); err != nil {
		L.Push(lua.LBool(false))
		L.ArgError(2, "connect error")
	}

	client := checkClient(L)
	cs := drivers.GetValue(L, 2)

	dsn, err := parseConnectionString(cs, client.Timeout)
	if err != nil {
		L.ArgError(2, err.Error())
		return 0
	}

	client.DB = redis.NewClient(dsn)
	L.Push(lua.LBool(true))
	return 1
}

func clientGetStrMethod(L *lua.LState) int {
	client := checkClient(L)
	key := L.ToString(2)
	if client.DB == nil {
		return 0
	}

	if key == "" {
		L.ArgError(2, "key string required")
		return 0
	}

	resp, err := client.DB.Get(client.Ctx, key).Result()
	if err != nil {
		L.ArgError(2, err.Error())
		return 0
	}
	L.Push(lua.LString(resp))
	return 1
}

func clientSetStrMethod(L *lua.LState) int {
	client := checkClient(L)
	key := L.ToString(2)
	value := L.ToString(3)
	expiration := L.ToInt(4)
	if client.DB == nil {
		return 0
	}

	if key == "" || value == "" {
		L.ArgError(2, "key|value string required")
		return 0
	}

	var outTime time.Duration = 0
	if expiration != 0 {
		outTime = time.Duration(expiration) * time.Millisecond
	}

	err := client.DB.Set(client.Ctx, key, value, outTime).Err()
	if err != nil {
		L.ArgError(2, err.Error())
	}

	L.Push(lua.LBool(true))
	return 1
}

func clientDelMethod(L *lua.LState) int {
	client := checkClient(L)
	key := L.ToString(2)
	if client.DB == nil {
		return 0
	}

	if key == "" {
		L.ArgError(2, "key string required")
		return 0
	}

	err := client.DB.Del(client.Ctx, key).Err()
	if err != nil {
		L.ArgError(2, err.Error())
	}
	L.Push(lua.LBool(true))
	return 1
}

func checkClient(L *lua.LState) *Client {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Client); ok {
		return v
	}
	L.ArgError(1, "client expected")
	return nil
}

func clientCloseMethod(L *lua.LState) int {
	client := checkClient(L)
	if client.DB == nil {
		L.Push(lua.LBool(true))
		return 1
	}

	err := client.DB.Close()
	// always clean
	client.DB = nil
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

func clientSetTimeoutMethod(L *lua.LState) int {
	client := checkClient(L)
	timeout := L.ToInt64(2) // timeout (in ms)

	client.Timeout = time.Millisecond * time.Duration(timeout)
	return 0
}
