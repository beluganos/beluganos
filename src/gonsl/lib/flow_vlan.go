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
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"

	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

//
// FIBCVLANFlowMod process FlowMod(VLAN)
//
func (s *Server) FIBCVLANFlowMod(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.VLANFlow) {
	log.Debugf("Server: FlowMod(VLAN): %v %v %v", hdr, mod, flow)

	vid := func() opennsl.Vlan {
		if vid := opennsl.Vlan(flow.Match.Vid); vid != opennsl.VLAN_ID_NONE {
			return vid
		}
		return opennsl.VlanDefaultMustGet(s.Unit())
	}()

	port := opennsl.Port(flow.Match.InPort)
	pvid := s.vlanPorts.ConvVID(port, vid)
	pbmp := opennsl.NewPBmp()
	pbmp.Add(port)

	switch mod.Cmd {
	case fibcapi.FlowMod_ADD:
		if _, err := pvid.Create(s.Unit()); err != nil {
			log.Errorf("Server: FlowMod(VLAN): vlan create error. vid:%d/%d", vid, pvid)
			return
		}

		ubmp := opennsl.NewPBmp()
		if untag := (vid == opennsl.VlanDefaultMustGet(s.Unit())); untag {
			ubmp.Add(port)
		}

		if err := pvid.PortAdd(s.Unit(), pbmp, ubmp); err != nil {
			log.Errorf("Server: FlowMod(VLAN): Port add error. vid:%d/%d port:%d %s", vid, pvid, port, err)
		}

	case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
		if _, err := pvid.PortRemove(s.Unit(), pbmp); err != nil {
			log.Errorf("Server: FlowMod(VLAN): Port remove error. vid:%d/%d port:%d %s", vid, pvid, port, err)
		}

	default:
		log.Warnf("Server: FlowMod(VLAN): Invalid cmd. %s", mod.Cmd)
	}
}
