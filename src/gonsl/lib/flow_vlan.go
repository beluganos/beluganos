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

	vlan := opennsl.Vlan(flow.Match.Vid)
	if vlan == opennsl.VLAN_ID_NONE || vlan == opennsl.VlanDefaultMustGet(s.Unit) {
		log.Debugf("Server: FlowMod(VLAN): Skip. vid=%d", vlan)
		return
	}

	if _, err := vlan.Create(s.Unit); err != nil {
		log.Errorf("Server: FlowMod(VLAN): Create error. %d %s", vlan, err)
		return
	}

	port := opennsl.Port(flow.Match.InPort)
	pbmp := opennsl.NewPBmp()
	pbmp.Add(port)

	switch mod.Cmd {
	case fibcapi.FlowMod_ADD:
		ubmp := opennsl.NewPBmp()
		if err := vlan.PortAdd(s.Unit, pbmp, ubmp); err != nil {
			log.Errorf("Server: FlowMod(VLAN): Port add error. %d %d %s", vlan, port, err)
		}

	case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
		if err := vlan.PortRemove(s.Unit, pbmp); err != nil {
			log.Errorf("Server: FlowMod(VLAN): Port remove error. %d %d %s", vlan, port, err)
		}

	default:
		log.Warnf("Server: FlowMod(VLAN): Invalid cmd. %s", mod.Cmd)
	}
}
