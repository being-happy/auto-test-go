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

type ScenarioCase struct {
	Id              int64                  `json:"id"`
	Name            string                 `json:"name"`
	Parameters      []CaseParameter        `json:"parameters"`
	Design          string                 `json:"design"`
	UserCases       map[string]*UserCase   `json:"userCases"`
	Baggages        map[string]*Baggage    `json:"baggages"`
	Scripts         map[string]*BaseScript `json:"scripts"`
	PreScripts      BaseScripts            `json:"preScripts"`
	AfterScripts    BaseScripts            `json:"afterScripts"`
	DependFunctions []string               `json:"dependFunctions"`
}
