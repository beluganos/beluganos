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
	"fmt"
	"net"

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

func (s *Server) fibcL2InterfaceGroupTrunkAdd(port uint32, vid uint16) error {
	log.Debugf("Server: L2-IF group: add Trunk. port:%x vid:%d", port, vid)

	if _, ok := s.idmaps.Trunks.Get(port, vid); ok {
		log.Warnf("Server: L2-IF group: Trunk already exists. port:%x vid:%d", port, vid)
		return nil
	}

	trunk, err := opennsl.TrunkCreate(s.Unit(), opennsl.TRUNK_FLAG_NONE)
	if err != nil {
		log.Errorf("Server: L2-IF group: create Trunk error. %s", err)
		return err
	}

	s.idmaps.Trunks.Register(port, vid, trunk)

	log.Debugf("Server: L2-IF group: add Trunk. id:%x port:%x vid:%d", trunk, port, vid)

	return nil
}

func (s *Server) fibcL2InterfaceGroupTrunkDel(port uint32, vid uint16) error {
	log.Debugf("Server: L2-IF group: del Trunk. port:%x vid:%d", port, vid)

	trunk, ok := s.idmaps.Trunks.Get(port, vid)
	if !ok {
		log.Warnf("Server: L2-IF group: Trunk not exists. port:%x vid:%d", port, vid)
		return nil
	}

	if err := trunk.Destroy(s.Unit()); err != nil {
		log.Errorf("Server: L2-IF group: Trunk destroy error. %d %e", trunk, err)
		return err
	}

	s.idmaps.Trunks.Unregister(port, vid)

	log.Debugf("Server: L2-IF group: del Trunk. id:%x port:%x vid:%d", trunk, port, vid)

	return nil
}

func (s *Server) fibcL2InterfaceGroupTrunkMemberAdd(master, slave uint32, vid uint16) error {
	log.Debugf("Server: L2-IF group: add Trunk member. master:%x slave:%x vid:%d", master, slave, vid)

	trunk, ok := s.idmaps.Trunks.Get(master, vid)
	if !ok {
		log.Errorf("Server: L2-IF group: Trunk not exists. master:%x vid:%d", master, vid)
		return fmt.Errorf("Trunk not exists. master:%x vid:%d", master, vid)
	}

	port, _ := fibcapi.ParseDPPortId(slave)
	gport, err := opennsl.Port(port).GPortGet(s.Unit())
	if err != nil {
		log.Errorf("server: L2-IF group: Trunk member gport error. port:%d vid:%d %s",
			port, vid, err)
		return err
	}

	member := opennsl.NewTrunkMember()
	member.SetGPort(gport)

	if err := trunk.MemberAdd(s.Unit(), member); err != nil {
		log.Errorf("server: L2-IF group: Trunk member add error. id:%x port:%d vid:%d %s",
			trunk, port, vid, err)
		return err
	}

	log.Debugf("Server: L2-IF group: add Trunk member. id:%x port:%d vid:%d",
		trunk, port, vid)

	return nil
}

func (s *Server) fibcL2InterfaceGroupTrunkMemberDel(master, slave uint32, vid uint16) error {
	log.Debugf("Server: L2-IF group: del Trunk member. master:%x slave:%x vid:%d", master, slave, vid)

	trunk, ok := s.idmaps.Trunks.Get(master, vid)
	if !ok {
		log.Warnf("Server: L2-IF group: Trunk not exists. master:%x vid:%d", master, vid)
		return nil
	}

	port, _ := fibcapi.ParseDPPortId(slave)
	gport, err := opennsl.Port(port).GPortGet(s.Unit())
	if err != nil {
		log.Errorf("server: L2-IF group: Trunk member gport error. port:%d vid:%d %s",
			port, vid, err)
		return err
	}

	member := opennsl.NewTrunkMember()
	member.SetGPort(gport)

	if err := trunk.MemberDelete(s.Unit(), member); err != nil {
		log.Errorf("server: L2-IF group: Trunk member del error. id:%x port:%d vid:%d %s",
			trunk, port, vid, err)
		return err
	}

	log.Debugf("Server: L2-IF group: del Trunk member. id:%x port:%d vid:%d",
		trunk, port, vid)

	return nil
}

func (s *Server) fibcL2InterfaceGroupAdd(port uint32, vid uint16, mac string, vrf uint32) error {
	if _, ok := s.idmaps.L3Ifaces.Get(port, vid); ok {
		log.Warnf("Server: L2-IF group: L2-IF already exists. port:%x vid:%d", port, vid)
		return nil
	}

	hwaddr, err := net.ParseMAC(mac)
	if err != nil {
		log.Errorf("Server: L2-IF group: Invalid MAC. '%s'", mac)
		return err
	}

	pvid := s.vlanPorts.ConvVID(opennsl.Port(port), opennsl.Vlan(vid))

	l3iface := opennsl.NewL3Iface()
	l3iface.SetMAC(hwaddr)
	l3iface.SetVID(opennsl.Vlan(pvid))
	if vrf != 0 {
		l3iface.SetVRF(opennsl.Vrf(vrf))
	}

	if err := l3iface.Create(s.Unit()); err != nil {
		log.Errorf("Server: L2-IF group: L3Iface create error. %s", err)
		return err
	}

	s.idmaps.L3Ifaces.Register(port, vid, l3iface.IfaceID())

	log.Debugf("Server: L2-IF group: add L3Iface. id:%x port:%x vid:%d/%d.",
		l3iface.IfaceID(), port, vid, pvid)

	return nil
}

func (s *Server) fibcL2InterfaceGroupDel(port uint32, vid uint16) error {
	ifaceID, ok := s.idmaps.L3Ifaces.Get(port, vid)
	if !ok {
		log.Warnf("Server: L2-IF group: L2-IF not found. port:%x vid:%d", port, vid)
		return nil
	}

	s.idmaps.L3Ifaces.Unregister(port, vid)

	l3iface := opennsl.NewL3Iface()
	l3iface.SetIfaceID(ifaceID)
	if err := l3iface.Delete(s.Unit()); err != nil {
		log.Errorf("Server: L2-IF group: L3Iface delete error. port:%x vid:%d %s", port, vid, err)
		return err
	}

	log.Debugf("Server: L2-IF group: del L3Iface. id:%x port:%x vid:%d.",
		ifaceID, port, vid)

	return nil
}

func (s *Server) fibcL2InterfaceGroupUpd(port uint32, vid uint16, mac string, vrf uint32) error {
	ifaceID, ok := s.idmaps.L3Ifaces.Get(port, vid)
	if !ok {
		log.Warnf("Server: L2-IF group: L2-IF(port:%x vid:%d) not exists. ", port, vid)
		return fmt.Errorf("L2-IF(port:%x vid:%d)  not exist.", port, vid)
	}

	hwaddr, err := net.ParseMAC(mac)
	if err != nil {
		log.Errorf("Server: L2-IF group: Invalid MAC. '%s'", mac)
		return err
	}

	pvid := s.vlanPorts.ConvVID(opennsl.Port(port), opennsl.Vlan(vid))

	l3iface := opennsl.NewL3Iface()
	l3iface.SetFlags(opennsl.L3_WITH_ID | opennsl.L3_REPLACE)
	l3iface.SetIfaceID(ifaceID)
	l3iface.SetMAC(hwaddr)
	l3iface.SetVID(opennsl.Vlan(pvid))
	if vrf != 0 {
		l3iface.SetVRF(opennsl.Vrf(vrf))
	}

	if err := l3iface.Create(s.Unit()); err != nil {
		log.Errorf("Server: L2-IF group: L3Iface create error. %s", err)
		return err
	}

	log.Debugf("Server: L2-IF group: upd L3Iface. id:%x port:%x vid:%d/%d.",
		ifaceID, port, vid, pvid)

	return nil
}

//
// FIBCL2InterfaceGroupMod process GroupMod(L2 Interface)
//
func (s *Server) FIBCL2InterfaceGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.L2InterfaceGroup) {
	log.Debugf("Server: L2-IF group: %v %v %v", hdr, mod, group)

	vid := group.GetAdjustedVlanVid()
	_, portType := fibcapi.ParseDPPortId(group.PortId)
	_, masterType := fibcapi.ParseDPPortId(group.Master)

	log.Debugf("Server: L2-IF group: port:%d %s master:%d %s vid:%d",
		group.PortId, portType, group.Master, masterType, vid)

	switch mod.Cmd {
	case fibcapi.GroupMod_ADD:
		if err := s.fibcL2InterfaceGroupAdd(group.PortId, vid, group.HwAddr, group.Vrf); err != nil {
			return
		}

		if portType == fibcapi.LinkType_BOND {
			s.fibcL2InterfaceGroupTrunkAdd(group.PortId, vid)
		} else if masterType == fibcapi.LinkType_BOND {
			s.fibcL2InterfaceGroupTrunkMemberAdd(group.Master, group.PortId, vid)
		}

	case fibcapi.GroupMod_DELETE:
		if portType == fibcapi.LinkType_BOND {
			s.fibcL2InterfaceGroupTrunkDel(group.PortId, vid)
		} else if masterType == fibcapi.LinkType_BOND {
			s.fibcL2InterfaceGroupTrunkMemberDel(group.Master, group.PortId, vid)
		}

		s.fibcL2InterfaceGroupDel(group.PortId, vid)

	case fibcapi.GroupMod_MODIFY:
		s.fibcL2InterfaceGroupUpd(group.PortId, vid, group.HwAddr, group.Vrf)

	default:
		log.Errorf("Server: L2-IF group: Invalid command. %d", mod.Cmd)
	}
}
