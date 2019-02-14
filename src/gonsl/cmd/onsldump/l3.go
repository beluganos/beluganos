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

func printL3Egress(l3egr *api.L3Egress) {
	fmt.Printf("L3Egress: flags:%08x/%08x l3egrId:%d ifaceId:%d mac:%s vid:%d port:%d\n",
		l3egr.GetFlags(),
		l3egr.GetFlags2(),
		l3egr.GetEgressId(),
		l3egr.GetIfaceId(),
		l3egr.GetMac(),
		l3egr.GetVid(),
		l3egr.GetPort(),
	)
}

func dumpL3Egresses(client api.GoNSLApiClient) {
	l3egrs, err := client.GetL3Egresses(context.Background(), api.NewGetL3EgressesRequest())
	if err != nil {
		log.Errorf("GetL3Egresses error. %s", err)
		return
	}

	for _, l3egr := range l3egrs.Egresses {
		printL3Egress(l3egr)
	}
}

func printL3Host(l3host *api.L3Host) {
	fmt.Printf("L3Host: Flags:%08x l3egrId:%d ip:%s ip6:%s mac:%s vrf:%d\n",
		l3host.GetFlags(),
		l3host.GetEgressId(),
		l3host.GetIpAddr(),
		l3host.GetIp6Addr(),
		l3host.GetMac(),
		l3host.GetVrf(),
	)
}

func dumpL3Hosts(client api.GoNSLApiClient) {
	l3hosts, err := client.GetL3Hosts(context.Background(), api.NewGetL3HostsRequest())
	if err != nil {
		log.Errorf("GetL3Hosts error. %s", err)
		return
	}

	for _, l3host := range l3hosts.Hosts {
		printL3Host(l3host)
	}
}

func printL3Route(l3route *api.L3Route) {
	fmt.Printf("L3Route: Flags:%08x l3egrId:%d ip:%s ip6:%s vrf:%d\n",
		l3route.GetFlags(),
		l3route.GetEgressId(),
		l3route.GetIpAddr(),
		l3route.GetIp6Addr(),
		l3route.GetVrf(),
	)
}

func dumpL3Routes(client api.GoNSLApiClient) {
	l3routes, err := client.GetL3Routes(context.Background(), api.NewGetL3RoutesRequest())
	if err != nil {
		log.Errorf("GetL3Routes error. %s", err)
		return
	}

	for _, l3route := range l3routes.Routes {
		printL3Route(l3route)
	}
}

func printL3Iface(iface *api.L3Iface) {
	fmt.Printf("L3Iface: flags:%08x, ifaceId:%d mac:%s vid:%d vrf:%d mtu:%d:%d ttl:%d\n",
		iface.GetFlags(),
		iface.GetIfaceId(),
		iface.GetMac(),
		iface.GetVid(),
		iface.GetVrf(),
		iface.GetMtu(),
		iface.GetMtuFwd(),
		iface.GetTtl(),
	)
}

func dumpL3Ifaces(client api.GoNSLApiClient) {
	reply, err := client.GetL3Ifaces(context.Background(), api.NewGetL3IfacesRequest())
	if err != nil {
		log.Errorf("GetL3Routes error. %s", err)
		return
	}

	for _, iface := range reply.Ifaces {
		printL3Iface(iface)
	}
}
