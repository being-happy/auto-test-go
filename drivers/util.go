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

package drivers

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"reflect"
	"strconv"
	"strings"
)

func GetValue(l *lua.LState, n int) interface{} {
	return GetArbitraryValue(l, l.Get(n))
}

func GetArbitraryValue(l *lua.LState, v lua.LValue) interface{} {
	switch t := v.Type(); t {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return lua.LVAsBool(v)
	case lua.LTNumber:
		f := lua.LVAsNumber(v)
		if float64(f) == float64(int(f)) {
			return int(f)
		}
		return float64(f)
	case lua.LTString:
		return lua.LVAsString(v)
	case lua.LTTable:
		m := map[string]interface{}{}
		tb := v.(*lua.LTable)
		arrSize := 0
		tb.ForEach(func(k, val lua.LValue) {
			key := GetArbitraryValue(l, k)
			if keyi, ok := key.(int); ok {
				if arrSize >= 0 && arrSize < keyi {
					arrSize = keyi
				}
				key = strconv.Itoa(keyi)
			} else {
				arrSize = -1
			}
			m[key.(string)] = GetArbitraryValue(l, val)
		})

		if arrSize >= 0 {
			ms := make([]interface{}, arrSize)
			for i := 0; i < arrSize; i++ {
				ms[i] = m[strconv.Itoa(i+1)]
			}
			return ms
		}

		return m
	default:
		panic(fmt.Sprintf("unknown lua type: %s", t))
	}
}

func ToArbitraryValue(l *lua.LState, i interface{}) lua.LValue {
	if i == nil {
		return lua.LNil
	}

	switch ii := i.(type) {
	case bool:
		return lua.LBool(ii)
	case int:
		return lua.LNumber(ii)
	case int8:
		return lua.LNumber(ii)
	case int16:
		return lua.LNumber(ii)
	case int32:
		return lua.LNumber(ii)
	case int64:
		return lua.LNumber(ii)
	case uint:
		return lua.LNumber(ii)
	case uint8:
		return lua.LNumber(ii)
	case uint16:
		return lua.LNumber(ii)
	case uint32:
		return lua.LNumber(ii)
	case uint64:
		return lua.LNumber(ii)
	case float64:
		return lua.LNumber(ii)
	case float32:
		return lua.LNumber(ii)
	case string:
		return lua.LString(ii)
	case []byte:
		return lua.LString(ii)
	default:
		v := reflect.ValueOf(i)
		switch v.Kind() {
		case reflect.Ptr:
			return ToArbitraryValue(l, v.Elem().Interface())

		case reflect.Struct:
			return ToTableFromStruct(l, v)

		case reflect.Map:
			return ToTableFromMap(l, v)

		case reflect.Slice:
			return ToTableFromSlice(l, v)

		default:
			panic(fmt.Sprintf("unknown type being pushed onto lua stack: %T %+v", i, i))
		}
	}
}

func ToTableFromStruct(l *lua.LState, v reflect.Value) lua.LValue {
	tb := l.NewTable()
	return toTableFromStructInner(l, tb, v)
}

func toTableFromStructInner(l *lua.LState, tb *lua.LTable, v reflect.Value) lua.LValue {
	t := v.Type()
	for j := 0; j < v.NumField(); j++ {
		var inline bool
		name := t.Field(j).Name
		if tag := t.Field(j).Tag.Get("luautil"); tag != "" {
			tagParts := strings.Split(tag, ",")
			if tagParts[0] == "-" {
				continue
			} else if tagParts[0] != "" {
				name = tagParts[0]
			}
			if len(tagParts) > 1 && tagParts[1] == "inline" {
				inline = true
			}
		}
		if inline {
			toTableFromStructInner(l, tb, v.Field(j))
		} else {
			tb.RawSetString(name, ToArbitraryValue(l, v.Field(j).Interface()))
		}
	}
	return tb
}

func ToTableFromMap(l *lua.LState, v reflect.Value) lua.LValue {
	tb := &lua.LTable{}
	for _, k := range v.MapKeys() {
		tb.RawSet(ToArbitraryValue(l, k.Interface()),
			ToArbitraryValue(l, v.MapIndex(k).Interface()))
	}
	return tb
}

func ToTableFromSlice(l *lua.LState, v reflect.Value) lua.LValue {
	tb := &lua.LTable{}
	for j := 0; j < v.Len(); j++ {
		tb.RawSet(ToArbitraryValue(l, j+1), // because lua is 1-indexed
			ToArbitraryValue(l, v.Index(j).Interface()))
	}
	return tb
}
