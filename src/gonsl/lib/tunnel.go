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
	"net"

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

func tunnelinitiatorSetSrcMAC(unit int, ifaceId opennsl.L3IfaceID, hwaddr net.HardwareAddr, vid opennsl.Vlan) error {
	l3iface, err := opennsl.L3IfaceGet(unit, ifaceId)
	if err != nil {
		return err
	}

	l3iface.SetFlags(opennsl.L3_REPLACE | opennsl.L3_WITH_ID)
	l3iface.SetMAC(hwaddr)
	l3iface.SetVID(vid)

	log.Debugf("GroupMod(L3-UC) initiator mac replaced. %d %s vid:%d", ifaceId, hwaddr, vid)
	return l3iface.Create(unit)
}

func tunnelInitiatorAdd(unit int, group *fibcapi.L3UnicastGroup, ifaceId opennsl.L3IfaceID, pvid opennsl.Vlan) {

	tun := &opennsl.TunnelInitiator{}
	tun.Init()
	tun.SetTTL(64)
	tun.SetVID(pvid)
	tun.SetL3IfaceID(ifaceId)

	switch group.TunType {
	case fibcapi.TunnelType_IPIP:
		tun.SetType(opennsl.TunnelTypeIPIP4encap)
		tun.SetDstIP4(group.GetTunRemoteIP())
		tun.SetSrcIP4(group.GetTunLocalIP())

		log.Debugf("GroupMod(L3-UC) tunnel initiator add. type=%s iface=%d vid=%d dst=%s src=%s",
			tun.Type(), tun.L3IfaceID(), tun.VID(), tun.DstIP4(), tun.SrcIP4())

	case fibcapi.TunnelType_IPV6:
		tun.SetType(opennsl.TunnelTypeIPIP6encap)
		tun.SetDstIP6(group.GetTunRemoteIP())
		tun.SetSrcIP6(group.GetTunLocalIP())

		log.Debugf("GroupMod(L3-UC) tunnel initiator add. type=%s iface=%d vid=%d dst=%s src=%s",
			tun.Type(), tun.L3IfaceID(), tun.VID(), tun.DstIP6(), tun.SrcIP6())

	case fibcapi.TunnelType_NOP:
		log.Debugf("GroupMod(L3-UC) tunnel initiator add. type=%s iface=%d not tunnel.", group.TunType, ifaceId)
		return

	default:
		log.Errorf("GroupMod(L3-UC) tunnel initiator add. invalid tunnel. type=%s iface=%d", group.TunType, ifaceId)
		return
	}

	if err := tunnelinitiatorSetSrcMAC(unit, ifaceId, group.GetEthSrcHwAddr(), pvid); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel initiator set-mac errror. iface:%d %s", ifaceId, err)
		return
	}

	iface := opennsl.NewL3Iface()
	iface.SetIfaceID(ifaceId)
	if err := tun.Create(unit, iface); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel initiator  crate error. iface:%d %s", ifaceId, err)
	}
}

func tunnelInitiatorDelete(unit int, group *fibcapi.L3UnicastGroup, ifaceId opennsl.L3IfaceID) {

	switch group.TunType {
	case fibcapi.TunnelType_IPIP:
	case fibcapi.TunnelType_IPV6:
	default:
		log.Debugf("GroupMod(L3-UC) tunnel initiator del. type=%s iface=%d not tunnel.",
			group.TunType, ifaceId)
		return
	}

	log.Debugf("GroupMod(L3-UC) tunnel initiator del. type=%s iface=%d", group.TunType, ifaceId)

	iface := opennsl.NewL3Iface()
	iface.SetIfaceID(ifaceId)

	if err := iface.TunnelInitiatorClear(unit); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel initiator clear error. iface:%d %s", ifaceId, err)
	}

	if err := tunnelinitiatorSetSrcMAC(unit, ifaceId, fibcapi.HardwareAddrDummy, 1); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel initiator L3-IF clear error. iface:%d %s", ifaceId, err)
	}
}

func newTunnelTerminator4(dst, src net.IP, port opennsl.Port, tunType opennsl.TunnelType) *opennsl.TunnelTerminator {

	mask := net.CIDRMask(32, 32)
	tun := opennsl.NewTunnelTerminator(tunType)
	tun.SetDstIP4(dst)
	tun.SetDstIPMask4(mask)
	tun.SetSrcIP4(src)
	tun.SetSrcIPMask4(mask)
	tun.PBmp().Add(port)
	// tun.SetVID(vid)

	log.Debugf("GroupMod(L3-UC) tunnel terminator. type=%s port=%d vid=%d dst=%s src=%s",
		tun.Type(), port, tun.VID(), tun.DstIPNet4(), tun.SrcIPNet4())

	return tun
}

func newTunnelTerminator6(dst, src net.IP, port opennsl.Port, tunType opennsl.TunnelType) *opennsl.TunnelTerminator {

	mask := net.CIDRMask(128, 128)
	tun := opennsl.NewTunnelTerminator(tunType)
	tun.SetDstIP6(dst)
	tun.SetDstIPMask6(mask)
	tun.SetSrcIP6(src)
	tun.SetSrcIPMask6(mask)
	tun.PBmp().Add(port)
	// tun.SetVID(vid)

	log.Debugf("GroupMod(L3-UC) tunnel terminator. type=%s port=%d vid=%d dst=%s src=%s",
		tun.Type(), port, tun.VID(), tun.DstIPNet6(), tun.SrcIPNet6())

	return tun
}

func newTunnelTerminators(group *fibcapi.L3UnicastGroup) (to4Tun *opennsl.TunnelTerminator, to6Tun *opennsl.TunnelTerminator) {

	dst := group.GetTunLocalIP()
	src := group.GetTunRemoteIP()
	port := opennsl.Port(group.PhyPortId)

	switch group.TunType {
	case fibcapi.TunnelType_IPIP:
		to4Tun = newTunnelTerminator4(dst, src, port, opennsl.TunnelTypeIPIP4toIP4)
		to6Tun = newTunnelTerminator4(dst, src, port, opennsl.TunnelTypeIPIP4toIP6)

	case fibcapi.TunnelType_IPV6:
		to4Tun = newTunnelTerminator6(dst, src, port, opennsl.TunnelTypeIPIP6toIP4)
		to6Tun = newTunnelTerminator6(dst, src, port, opennsl.TunnelTypeIPIP6toIP6)

	case fibcapi.TunnelType_NOP:
		log.Debugf("GroupMod(L3-UC) tunnel terminator init. type=%s port=%d not tunnel.", group.TunType, port)
		return

	default:
		log.Warnf("GroupMod(L3-UC) tunnel terminator init. invalid tunnel. type=%s port=%d", group.TunType, port)
		return
	}

	return
}

func tunnelTerminatorAdd(unit int, group *fibcapi.L3UnicastGroup) {
	log.Debugf("GroupMod(L3-UC) tunnel terminator add.")

	to4Tun, to6Tun := newTunnelTerminators(group)
	if to4Tun == nil {
		return
	}

	if err := to4Tun.Create(unit); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel terminator(to ipv6) create error. %v %s", to4Tun, err)
	}
	if err := to6Tun.Create(unit); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel terminator(to ipv6) create error. %v %s", to6Tun, err)
	}
}

func tunnelTerminatorDelete(unit int, group *fibcapi.L3UnicastGroup) {
	log.Debugf("GroupMod(L3-UC) tunnel terminator delete.")

	to4Tun, to6Tun := newTunnelTerminators(group)
	if to4Tun == nil {
		return
	}

	if err := to4Tun.Delete(unit); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel terminator(to ipv4) delete error. %v %s", to4Tun, err)
	}
	if err := to6Tun.Delete(unit); err != nil {
		log.Errorf("GroupMod(L3-UC) tunnel terminator(to ipv6) delete error. %v %s", to6Tun, err)
	}
}
