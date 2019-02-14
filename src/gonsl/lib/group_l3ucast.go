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
	"net"

	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

//
// FIBCL3UnicastGroupMod process GroupMod(L3 Unicst)
//
func (s *Server) FIBCL3UnicastGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.L3UnicastGroup) {
	log.Debugf("Server: GroupMod(L3-UC): %v %v %v", hdr, mod, group)

	hwaddr, err := net.ParseMAC(group.EthDst)
	if err != nil {
		log.Errorf("Server: L3-UC group: Invalid MAC. '%s'", group.EthDst)
		return
	}

	vid := fibcapi.AdjustVlanVID16(uint16(group.VlanVid))
	port := group.PortId
	neid := group.NeId

	switch mod.Cmd {
	case fibcapi.GroupMod_ADD, fibcapi.GroupMod_MODIFY:
		if _, ok := s.idmaps.L3Egress.Get(neid); ok {
			log.Errorf("Server: L3-UC group: L3-UC(neid:%d) already exists. ", neid)
			return
		}

		ifaceID, ok := s.idmaps.L3Ifaces.Get(port, vid)
		if !ok {
			log.Errorf("Server: L3-UC group: L2-IF(port:%d, vid:%d) not found. ", port, vid)
			return
		}

		flags := opennsl.L3_NONE
		if mod.Cmd == fibcapi.GroupMod_MODIFY {
			flags = opennsl.L3_REPLACE
		}

		l3egr := opennsl.NewL3Egress()
		l3egr.SetIfaceID(opennsl.L3IfaceID(ifaceID))
		l3egr.SetPort(opennsl.Port(port))
		l3egr.SetMAC(hwaddr)
		l3egr.SetVID(opennsl.Vlan(vid))

		l3egrID, err := l3egr.Create(s.Unit, flags, 0)
		if err != nil {
			log.Errorf("Server: L3-UC group: L3 Egress create error. %s", err)
			return
		}

		s.idmaps.L3Egress.Register(neid, l3egrID)

	case fibcapi.GroupMod_DELETE:
		l3egrID, ok := s.idmaps.L3Egress.Get(neid)
		if !ok {
			log.Errorf("Server: L3-UC group: L3-UC(%08x) not found. ", neid)
			return
		}

		s.idmaps.L3Egress.Unregister(neid)

		if err := l3egrID.Destroy(s.Unit); err != nil {
			log.Errorf("Server: L3-UC group: L3 Egress delete error. %s", err)
		}

	default:
		log.Errorf("Server: L3-UC group: Invalid Cmd. %d", mod.Cmd)
	}
}
