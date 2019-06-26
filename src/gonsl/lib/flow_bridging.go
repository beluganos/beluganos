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

package gonslib

import (
	fibcapi "fabricflow/fibc/api"
	fibcnet "fabricflow/fibc/net"
	"net"

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

//
// FIBCBridgingFlowMod process FlowMod (Bridging)
//
func (s *Server) FIBCBridgingFlowMod(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.BridgingFlow) {
	log.Debugf("Server: FlowMod(Bridge): %v %v %v", hdr, mod, flow)

	ethDst, err := net.ParseMAC(flow.Match.EthDst)
	if err != nil {
		log.Errorf("Server: FlowMod(Bridge): Bad eth_dst. %s", err)
		return
	}

	if name := flow.Action.Name; name != fibcapi.BridgingFlow_Action_OUTPUT {
		log.Warnf("Server: FlowMod(Bridge): Bad action. %d %s", name, err)
		return
	}

	port, portType := fibcapi.ParseDPPortId(flow.Action.Value)
	switch portType {
	case fibcapi.LinkType_BRIDGE, fibcapi.LinkType_BOND:
		log.Debugf("Server: FlowMod(Bridge): %d %s skip.", port, portType)
		return
	}

	vid := opennsl.Vlan(flow.Match.VlanVid)

	l2addr := opennsl.NewL2Addr(ethDst, vid)
	if port != 0 {
		l2addr.SetPort(opennsl.Port(port))
	}

	switch mod.Cmd {
	case fibcapi.FlowMod_ADD, fibcapi.FlowMod_MODIFY:
		log.Debugf("Server: FlowMod(Bridge): l2addr add. %s", l2addr)
		if err := l2addr.Add(s.Unit()); err != nil {
			log.Errorf("Server: FlowMod(Bridge): l2addr add error. %s", err)
		}

	case fibcapi.FlowMod_DELETE:
		log.Debugf("Server: FlowMod(Bridge): l2addr delete. %s", l2addr)
		if err := l2addr.Delete(s.Unit()); err != nil {
			log.Errorf("Server: FlowMod(Bridge): l2addr delete error. %s", err)
		}

	default:
		log.Errorf("Server: FlowMod(Bridge): Bad  flow cmd. %d", mod.Cmd)
	}
}
