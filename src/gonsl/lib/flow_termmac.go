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

package gonslib

import (
	fibcapi "fabricflow/fibc/api"
	fibcnet "fabricflow/fibc/net"

	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

//
// FIBCTerminationMacFlowMod process FlowMod(Termination MAC)
//
func (s *Server) FIBCTerminationMacFlowMod(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.TerminationMacFlow) {
	s.log.Debugf("FlowMod(TermMAC): %v", hdr)
	fibcapi.LogFlowMod(s.log, log.DebugLevel, mod)

	port, portType := fibcapi.ParseDPPortId(flow.Match.InPort)
	switch portType {
	case fibcapi.LinkType_BRIDGE, fibcapi.LinkType_BOND:
		s.log.Debugf("FlowMod(TermMAC): %d %s skip", port, portType)
		return
	}

	mac, mask, err := fibcapi.ParseMaskedMAC(flow.Match.EthDst)
	if err != nil {
		s.log.Errorf("FlowMod(TermMAC): Invalid MAC. %s", err)
		return
	}

	if mask.String() != fibcapi.HWADDR_EXACT_MASK {
		s.log.Debugf("FlowMod(TermMAC): MAC is masked. %s '%s'", flow.Match.EthDst, mask)
		return
	}

	if ethType := flow.Match.EthType; ethType != fibcapi.ETHTYPE_IPV4 && ethType != fibcapi.ETHTYPE_IPV6 {
		s.log.Debugf("FlowMod(TermMAC): Not IPv4/6. %d %s", ethType, flow.Match.EthDst)
		return
	}

	vlan := func() opennsl.Vlan {
		if vid := opennsl.Vlan(flow.Match.VlanVid); vid != opennsl.VLAN_ID_NONE {
			return vid
		}
		return opennsl.VlanDefaultMustGet(s.Unit())
	}()

	l2addr := opennsl.NewL2Addr(mac, vlan)
	l2addr.SetFlags(opennsl.L2_L3LOOKUP | opennsl.L2_STATIC)
	l2addr.SetPort(opennsl.Port(port))

	switch mod.Cmd {
	case fibcapi.FlowMod_ADD:
		s.log.Debugf("FlowMod(TermMAC): L2Addr add. %s/%s port:%d vid:%d", mac, mask, flow.Match.InPort, vlan)
		if err := l2addr.Add(s.Unit()); err != nil {
			s.log.Errorf("FlowMod(TermMAC): L2 Addr Add error. %v %s", l2addr, err)
		}

	case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
		s.log.Debugf("FlowMod(TermMAC): L2Addr del. %s/%s port:%d vid:%d", mac, mask, flow.Match.InPort, vlan)
		if err := l2addr.Delete(s.Unit()); err != nil {
			s.log.Errorf("FlowMod(TermMAC): L2 Addr Delete error. %v %s", l2addr, err)
		}

	default:
		s.log.Warnf("FlowMod(TermMAC): Ignored. %s %s", mod.Cmd, flow)
	}
}
