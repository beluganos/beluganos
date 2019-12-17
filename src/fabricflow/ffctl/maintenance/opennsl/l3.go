// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package opennsl

import (
	"context"
	"fmt"
	api "gonsl/api"
	"io"
)

type L3Egress struct {
}

func NewL3Egress() *L3Egress {
	return &L3Egress{}
}

func (e *L3Egress) Name() string {
	return "l3-egress"
}

func (e *L3Egress) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetL3Egresses(context.Background(), api.NewGetL3EgressesRequest())
	if err != nil {
		return err
	}

	for _, l3egr := range reply.Egresses {
		fmt.Fprintf(w, "L3Egress: flags:%08x/%08x l3egrId:%d ifaceId:%d mac:%s vid:%d port:%d\n",
			l3egr.GetFlags(),
			l3egr.GetFlags2(),
			l3egr.GetEgressId(),
			l3egr.GetIfaceId(),
			l3egr.GetMac(),
			l3egr.GetVid(),
			l3egr.GetPort(),
		)
	}

	return nil
}

type L3Host struct {
}

func NewL3Host() *L3Host {
	return &L3Host{}
}

func (e *L3Host) Name() string {
	return "l3-host"
}

func (e *L3Host) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetL3Hosts(context.Background(), api.NewGetL3HostsRequest())
	if err != nil {
		return err
	}

	for _, l3host := range reply.Hosts {
		fmt.Fprintf(w, "L3Host: Flags:%08x l3egrId:%d ip:%s ip6:%s mac:%s vrf:%d\n",
			l3host.GetFlags(),
			l3host.GetEgressId(),
			l3host.GetIpAddr(),
			l3host.GetIp6Addr(),
			l3host.GetMac(),
			l3host.GetVrf(),
		)
	}

	return nil
}

type L3Route struct {
}

func NewL3Route() *L3Route {
	return &L3Route{}
}

func (e *L3Route) Name() string {
	return "l3-route"
}

func (e *L3Route) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetL3Routes(context.Background(), api.NewGetL3RoutesRequest())
	if err != nil {
		return err
	}

	for _, l3route := range reply.Routes {
		fmt.Fprintf(w, "L3Route: Flags:%08x l3egrId:%d ip:%s ip6:%s vrf:%d\n",
			l3route.GetFlags(),
			l3route.GetEgressId(),
			l3route.GetIpAddr(),
			l3route.GetIp6Addr(),
			l3route.GetVrf(),
		)
	}

	return nil
}

type L3Iface struct {
}

func NewL3Iface() *L3Iface {
	return &L3Iface{}
}

func (e *L3Iface) Name() string {
	return "l3-iface"
}

func (e *L3Iface) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetL3Ifaces(context.Background(), api.NewGetL3IfacesRequest())
	if err != nil {
		return err
	}

	for _, iface := range reply.Ifaces {
		fmt.Fprintf(w, "L3Iface: flags:%08x, ifaceId:%d mac:%s vid:%d vrf:%d mtu:%d:%d ttl:%d\n",
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

	return nil
}
