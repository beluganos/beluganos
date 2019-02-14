// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	api "gonsl/api"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func printTunnelInitiator(tun *api.TunnelInitiator) {
	fmt.Printf("Tun(Init): flags:%08x id:%d type:%s ifaceId:%d mac:%s-%s ip:%s-%s port:%d-%d ttl:%d mtu:%d vlan:%d\n",
		tun.GetFlags(),
		tun.GetTunnelId(),
		tun.GetTunnelType(),
		tun.GetL3IfaceId(),
		tun.GetDstMac(), tun.GetSrcMac(),
		tun.GetDstIp(), tun.GetSrcIp(),
		tun.GetDstPort(), tun.GetSrcPort(),
		tun.GetTtl(),
		tun.GetMtu(),
		tun.GetVlan(),
	)
}

func dumpTunnelInitiators(client api.GoNSLApiClient) {
	reply, err := client.GetTunnelInitiators(context.Background(), api.NewGetTunnelInitiatorsRequest())
	if err != nil {
		log.Errorf("GetTunnelInitiators error. %s", err)
		return
	}

	for _, tun := range reply.Tunnels {
		printTunnelInitiator(tun)
	}
}

func printTunnelTerminator(tun *api.TunnelTerminator) {
	fmt.Printf("Tun(Term): flags:%08x id:%d type:%s rport:%d ip:%s-%s port:%d-%d vlan:%d vrf:%d\n",
		tun.GetFlags(),
		tun.GetTunnelId(),
		tun.GetTunnelType(),
		tun.GetRemotePort(),
		tun.GetDstIp(), tun.GetSrcIp(),
		tun.GetDstPort(), tun.GetSrcPort(),
		tun.GetVlan(),
		tun.GetVrf(),
	)
}

func dumpTunnelTerminators(client api.GoNSLApiClient) {
	reply, err := client.GetTunnelTerminators(context.Background(), api.NewGetTunnelTerminatorsRequest())
	if err != nil {
		log.Errorf("GetTunnelTerminators error. %s", err)
		return
	}

	for _, tun := range reply.Tunnels {
		printTunnelTerminator(tun)
	}
}
