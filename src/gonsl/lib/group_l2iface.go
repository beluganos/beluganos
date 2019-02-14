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
// FIBCL2InterfaceGroupMod process GroupMod(L2 Interface)
//
func (s *Server) FIBCL2InterfaceGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.L2InterfaceGroup) {
	log.Debugf("Server: GroupMod(L2-IF): %v %v %v", hdr, mod, group)

	hwaddr, err := net.ParseMAC(group.HwAddr)
	if err != nil {
		log.Errorf("Server: L2-IF group: Invalid MAC. '%s'", group.HwAddr)
		return
	}

	vid := fibcapi.AdjustVlanVID16(uint16(group.VlanVid))
	port := group.PortId

	switch mod.Cmd {
	case fibcapi.GroupMod_ADD, fibcapi.GroupMod_MODIFY:
		if _, ok := s.idmaps.L3Ifaces.Get(port, vid); ok {
			log.Errorf("Server: L2-IF group: L2-IF(port:%d vid:%d) already exists. ", port, vid)
			return
		}

		l3iface := opennsl.NewL3Iface()
		l3iface.SetMAC(hwaddr)
		l3iface.SetVID(opennsl.Vlan(vid))
		if vrf := opennsl.Vrf(group.Vrf); vrf != 0 {
			l3iface.SetVRF(vrf)
		}

		if err := l3iface.Create(s.Unit); err != nil {
			log.Errorf("Server: L2-IF group: L3Iface create error. %s", err)
			return
		}

		s.idmaps.L3Ifaces.Register(port, vid, l3iface.IfaceID())

	case fibcapi.GroupMod_DELETE:
		ifaceID, ok := s.idmaps.L3Ifaces.Get(port, vid)
		if !ok {
			log.Errorf("Server: L2-IF group: L2-IF(port:%d, vid:%d) not found. ", port, vid)
			return
		}

		s.idmaps.L3Ifaces.Unregister(port, vid)

		l3iface := opennsl.NewL3Iface()
		l3iface.SetIfaceID(ifaceID)
		if err := l3iface.Delete(s.Unit); err != nil {
			log.Errorf("Server: L2-IF group: L3Iface delete error. %s", err)
		}

	default:
		log.Errorf("Server: L2-IF group: Invalid Cmd. %d", mod.Cmd)
	}
}
