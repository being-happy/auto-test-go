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

package util

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"os"
	"strconv"
)

var Nacos_Helper = NacosHelper{}

type NacosHelper struct {
	client    naming_client.INamingClient
	parameter vo.RegisterInstanceParam
}

func (n NacosHelper) RegisterServiceInstance() {
	nacosAdr := os.Getenv("NACOS_ADDRESS")
	nacosPort := os.Getenv("NACOS_PORT")
	namespace := os.Getenv("NAMESPACE")
	if nacosAdr == "" || nacosPort == "" {
		Logger.Warn("[NacosHelper] Nacos config is invalid.")
		return
	}

	intNum, _ := strconv.Atoi(nacosPort)
	sc := []constant.ServerConfig{{
		IpAddr:      nacosAdr,
		Port:        uint64(intNum),
		ContextPath: "/nacos",
	},
	}

	cc := constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// create naming client
	client, err := clients.CreateNamingClient(
		map[string]interface{}{
			"serverConfigs": sc,
			"clientConfig":  cc,
		})

	if err != nil {
		panic(err)
	}

	n.client = client
	n.parameter = vo.RegisterInstanceParam{
		Ip:          getCurrentIp(),
		Port:        8090,
		ServiceName: "auto-test-engine",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	}

	serviceClient_RegisterServiceInstance(n.client, n.parameter)
}

func (n NacosHelper) DeregisterInstance() {
	param := vo.DeregisterInstanceParam{
		Ip:          n.parameter.Ip,
		Port:        n.parameter.Port,
		ServiceName: n.parameter.ServiceName,
		Ephemeral:   true,
	}

	success, err := n.client.DeregisterInstance(param)
	if !success || err != nil {
		panic("Dregister Service Instance failed!" + err.Error())
	}

	fmt.Printf("Dregister Service Instance,param:%+v,result:%+v \n\n", param, success)
}

func serviceClient_RegisterServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	success, err := client.RegisterInstance(param)
	if !success || err != nil {
		panic("RegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func getCurrentIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic("net.Interfaces failed, err: " + err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}
