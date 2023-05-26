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

import (
	"auto-test-go/pkg/util"
	"strconv"
)

type SenarioFlowDesign struct {
	Flows Flows
}

type Order interface {
	GetOrder() int
}

type Flows []Order

func (b Flows) Len() int {
	return len(b)
}

func (b Flows) Less(i, j int) bool {
	return b[i].GetOrder() < b[j].GetOrder()
}

func (b Flows) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type UserCaseUnitDesign struct {
	Id    string `xml:"case-id,attr"`
	Name  string `xml:"case-name,attr"`
	Order string `xml:"order,attr"`
}

func (u UserCaseUnitDesign) GetOrder() int {
	v, err := strconv.Atoi(u.Order)
	if err != nil {
		util.Logger.Error("[UserCaseUnitDesign] convert string to int error: %s", err.Error())
	}
	return v
}

type ScriptUnitDesign struct {
	Id    string `xml:"script-id,attr"`
	Stype string `xml:"script-type,attr"`
	Order string `xml:"order,attr"`
}

func (u ScriptUnitDesign) GetOrder() int {
	v, err := strconv.Atoi(u.Order)
	if err != nil {
		util.Logger.Error("[UserCaseUnitDesign] convert string to int error: %s", err.Error())
	}
	return v
}

type LoopCaseDesign struct {
	Ltype     string `xml:"loop-type,attr"`
	DataId    string `xml:"data-id,attr"`
	LoopCount string `xml:"loop-count,attr"`
	Flows     Flows
	Order     string `xml:"order,attr"`
	Id        string `xml:"loop-id,attr"`
}

func (u LoopCaseDesign) GetOrder() int {
	if u.Order == "" {
		return 0
	}

	v, err := strconv.Atoi(u.Order)
	if err != nil {
		util.Logger.Error("[UserCaseUnitDesign] convert string to int error: %s", err.Error())
	}
	return v
}

type TimeWaitUnitDesign struct {
	WaitTime string `xml:"wait-time,attr"`
	Order    string `xml:"order,attr"`
}

func (u TimeWaitUnitDesign) GetOrder() int {
	v, err := strconv.Atoi(u.Order)
	if err != nil {
		util.Logger.Error("[UserCaseUnitDesign] convert string to int error: %s", err.Error())
	}
	return v
}
