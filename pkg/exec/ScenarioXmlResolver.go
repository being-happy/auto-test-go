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

package exec

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/util"
	"github.com/beevik/etree"
	"sort"
	"strings"
)

type SenarioXmlResolver struct {
}

func (s SenarioXmlResolver) ResolveDesign(design string) (*entities.SenarioFlowDesign, error) {
	doc := etree.NewDocument()
	design = strings.Replace(design, "&", "", -1)
	err := doc.ReadFromBytes([]byte(design))
	util.Logger.Info("[SenariorXmlResolver] Begin to resolve senarior xml data: %s", design)
	if err != nil {
		return nil, err
	}

	root := doc.SelectElement("scenario-flow")
	flowDesign := entities.SenarioFlowDesign{}
	flowDesign.Flows = s.commonResolve(root)

	return &flowDesign, err
}

func (s SenarioXmlResolver) resolveLoop(loop *etree.Element) entities.LoopCaseDesign {
	loopCaseDesign := entities.LoopCaseDesign{
		Ltype:     loop.SelectAttrValue("loop-type", ""),
		DataId:    loop.SelectAttrValue("data-id", ""),
		LoopCount: loop.SelectAttrValue("loop-count", ""),
		Order:     loop.SelectAttrValue("order", "1"),
		Id:        loop.SelectAttrValue("id", "1"),
	}

	loopCaseDesign.Flows = s.commonResolve(loop)
	return loopCaseDesign
}

func (s SenarioXmlResolver) commonResolve(root *etree.Element) entities.Flows {
	flows := entities.Flows{}
	for _, element := range root.SelectElements("user-case") {
		var userCase = entities.UserCaseUnitDesign{
			Id:    element.SelectAttrValue("case-id", "0"),
			Name:  element.SelectAttrValue("case-name", ""),
			Order: element.SelectAttrValue("order", "0"),
		}
		flows = append(flows, userCase)
	}

	for _, element := range root.SelectElements("time-wait-case") {
		timeWait := entities.TimeWaitUnitDesign{
			WaitTime: element.SelectAttrValue("wait-time", "0"),
			Order:    element.SelectAttrValue("order", "0"),
		}
		flows = append(flows, timeWait)
	}

	for _, element := range root.SelectElements("script-case") {
		scriptDesign := entities.ScriptUnitDesign{
			Id:    element.SelectAttrValue("script-id", "0"),
			Name:  element.SelectAttrValue("script-name", ""),
			Stype: element.SelectAttrValue("script-type", "ScriptType_LuaScript"),
			Order: element.SelectAttrValue("order", "0"),
		}
		flows = append(flows, scriptDesign)
	}

	childLoop := root.SelectElements("loop-script-case")
	if len(childLoop) >= 0 {
		for _, element := range childLoop {
			flows = append(flows, s.resolveLoop(element))
		}
	}

	for _, element := range root.SelectElements("condition") {
		conditoinDesign := entities.ConditoinUnitDesign{
			Id:            element.SelectAttrValue("condition-id", "0"),
			Name:          element.SelectAttrValue("condition-name", ""),
			Expr:          element.SelectAttrValue("expr", ""),
			Order:         element.SelectAttrValue("order", "0"),
			CorrectBranch: s.commonResolve(element.SelectElement("correct-condition")),
			ErrorBranch:   s.commonResolve(element.SelectElement("deny-condition")),
		}
		flows = append(flows, conditoinDesign)
	}

	sort.Sort(flows)
	return flows
}
