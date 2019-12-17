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

type TunnelInitiator struct {
}

func NewTunnelInitiator() *TunnelInitiator {
	return &TunnelInitiator{}
}

func (e *TunnelInitiator) Name() string {
	return "tun-init"
}

func (e *TunnelInitiator) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetTunnelInitiators(context.Background(), api.NewGetTunnelInitiatorsRequest())
	if err != nil {
		return err
	}

	for _, tun := range reply.Tunnels {
		fmt.Fprintf(w, "Tun(Init): flags:%08x id:%d type:%s ifaceId:%d mac:%s-%s ip:%s-%s port:%d-%d ttl:%d mtu:%d vlan:%d\n",
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

	return nil
}

type TunnelTerminator struct {
}

func NewTunnelTerminator() *TunnelTerminator {
	return &TunnelTerminator{}
}

func (e *TunnelTerminator) Name() string {
	return "tun-term"
}

func (e *TunnelTerminator) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetTunnelTerminators(context.Background(), api.NewGetTunnelTerminatorsRequest())
	if err != nil {
		return err
	}

	for _, tun := range reply.Tunnels {
		fmt.Fprintf(w, "Tun(Term): flags:%08x id:%d type:%s rport:%d ip:%s-%s port:%d-%d vlan:%d vrf:%d\n",
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
	return nil
}
