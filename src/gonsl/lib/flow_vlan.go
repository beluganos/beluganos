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
// FIBCVLANFlowMod process FlowMod(VLAN)
//
func (s *Server) FIBCVLANFlowMod(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.VLANFlow) {
	if flags, ok := flow.BridgeVlanInfoFlags(); ok {
		s.fibcVLANFlowModBrVlan(hdr, mod, flow, flags)
	} else {
		s.fibcVLANFlowMod(hdr, mod, flow)
	}
}

func (s *Server) fibcVLANFlowMod(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.VLANFlow) {
	s.log.Debugf("FlowMod(VLAN): %v", hdr)
	fibcapi.LogFlowMod(s.log, log.DebugLevel, mod)

	port, portType := fibcapi.ParseDPPortId(flow.Match.InPort)
	if portType.IsVirtual() {
		s.log.Debugf("FlowMod(VLAN): %d %s skip.", port, portType)
		return
	}

	vlan := NewL3Vlan(s.Unit(), opennsl.Vlan(flow.Match.Vid))
	vlan.Pbmp.Add(opennsl.Port(port))
	vlan.Vlan = s.vlanPorts.ConvVID(opennsl.Port(port), vlan.Vid)

	switch mod.Cmd {
	case fibcapi.FlowMod_ADD:
		s.log.Infof("FlowMod(VLAN): ADD Port. %s", vlan)
		if err := vlan.Create(s.Unit()); err != nil {
			s.log.Errorf("FlowMod(VLAN): ADD Port error. %s", err)
		}

	case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
		s.log.Infof("FlowMod(VLAN): DEL port. %s", vlan)
		if err := vlan.Delete(s.Unit()); err != nil {
			s.log.Errorf("FlowMod(VLAN): DEL Port error. %s", err)
		}

	default:
		s.log.Warnf("FlowMod(VLAN): Invalid cmd. %s %s", mod.Cmd, vlan)
	}
}

func (s *Server) fibcVLANFlowModBrVlan(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.VLANFlow, flags fibcapi.BridgeVlanInfo_Flags) {
	s.log.Debugf("FlowMod(BrVLAN): %v", hdr)
	fibcapi.LogFlowMod(s.log, log.DebugLevel, mod)

	port, _ := fibcapi.ParseDPPortId(flow.Match.InPort)
	vid := opennsl.Vlan(flow.Match.Vid)
	vlan := NewBrVlan(s.Unit(), vid)
	vlan.Pbmp.Add(opennsl.Port(port))

	if (flags & fibcapi.BridgeVlanInfo_PVID) != 0 {
		// Drop tagged packet.
		vlan.StrictlyUntagged = true
	}

	if (flags & fibcapi.BridgeVlanInfo_UNTAGGED) != 0 {
		// Egress packets are untagged.
		vlan.UntagBmp.Add(opennsl.Port(port))
	}

	switch mod.Cmd {
	case fibcapi.FlowMod_ADD:
		s.log.Infof("FlowMod(BrVLAN): ADD Port. %s", vlan)
		if err := vlan.Create(s.Unit()); err != nil {
			s.log.Errorf("FlowMod(BrVLAN): ADD Port error. %s", err)
			return
		}

		if err := s.notifyL2Addrs(opennsl.Port(port), vid); err != nil {
			s.log.Errorf("FlowMod(BrVLAN): Notify L2Addrs error. %s", err)
		}

	case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
		s.log.Infof("FlowMod(BrVLAN): DEL Port. %s", vlan)
		vlan.Delete(s.Unit())

	default:
		s.log.Errorf("FlowMod(BrVLAN): Invalid cmd. %s %s", mod.Cmd, vlan)
	}
}

func (s *Server) notifyL2Addrs(portId opennsl.Port, vid opennsl.Vlan) error {
	l2addrs := []*L2addrmonEntry{}
	if err := opennsl.L2Traverse(s.Unit(), func(unit int, l2addr *opennsl.L2Addr) opennsl.OpenNSLError {
		if l2addr.Port() == portId || l2addr.VID() == vid {
			e := NewL2addrmonEntry(l2addr, opennsl.L2_CALLBACK_ADD)
			l2addrs = append(l2addrs, e)
		}

		return opennsl.E_NONE

	}); err != nil {
		return err
	}

	s.l2addrCh <- l2addrs

	return nil
}
