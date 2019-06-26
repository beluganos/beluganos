// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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

package fibcapi

import (
	"net"

	"github.com/golang/protobuf/proto"
)

func (g *GroupMod) Type() uint16 {
	return uint16(FFM_GROUP_MOD)
}

func (g *GroupMod) Bytes() ([]byte, error) {
	return proto.Marshal(g)
}

func NewGroupModFromBytes(data []byte) (*GroupMod, error) {
	group_mod := &GroupMod{}
	if err := proto.Unmarshal(data, group_mod); err != nil {
		return nil, err
	}

	return group_mod, nil
}

//
// L2 Interface Group
//
func NewL2InterfaceGroupID(portId uint32, vlanId uint16) uint32 {
	return (AdjustVlanVID(vlanId) << 16) + (portId & 0xffff)
}

func NewL2InterfaceGroup(portId uint32, vlanId uint16, vlanTranslation bool, hwAddr net.HardwareAddr, mtu int, vrf uint8, master uint32) *L2InterfaceGroup {
	return &L2InterfaceGroup{
		PortId:          portId,
		VlanVid:         uint32(vlanId),
		VlanTranslation: vlanTranslation,
		HwAddr:          hwAddr.String(),
		Mtu:             uint32(mtu),
		Vrf:             uint32(vrf),
		Master:          master,
	}
}

func (g *L2InterfaceGroup) ToMod(cmd GroupMod_Cmd, reId string) *GroupMod {
	return &GroupMod{
		Cmd:   cmd,
		GType: GroupMod_L2_INTERFACE,
		ReId:  reId,
		Entry: &GroupMod_L2Iface{L2Iface: g},
	}
}

func (g *L2InterfaceGroup) GetAdjustedVlanVid() uint16 {
	return AdjustVlanVID16(uint16(g.VlanVid))
}

//
// L2 Rewrite Group
//
func NewL2RewriteGroupID(neighId uint32) uint32 {
	return 0x10000000 + (neighId & 0x0fffffff)
}

//
// L3 Unicast Group
//
func NewL3UnicastGroupID(neighId uint32) uint32 {
	return 0x20000000 + (neighId & 0x0fffffff)
}

func NewL3UnicastGroup(neId, portId, phyPortId uint32, vlanVid uint16, ethDst, ethSrc net.HardwareAddr) *L3UnicastGroup {
	return &L3UnicastGroup{
		NeId:      neId,
		PortId:    portId,
		VlanVid:   uint32(vlanVid),
		EthDst:    ethDst.String(),
		EthSrc:    ethSrc.String(),
		PhyPortId: phyPortId,
		TunType:   TunnelType_NOP,
		TunLocal:  "",
		TunRemote: "",
	}
}

func (g *L3UnicastGroup) SetTunnel(tunType TunnelType_Type, remote, local net.IP) {
	g.TunType = tunType
	g.TunLocal = local.String()
	g.TunRemote = remote.String()
}

func (g *L3UnicastGroup) GetAdjustedVlanVid() uint16 {
	return AdjustVlanVID16(uint16(g.VlanVid))
}

func (g *L3UnicastGroup) GetEthDstHwAddr() net.HardwareAddr {
	if hwaddr, err := net.ParseMAC(g.EthDst); err == nil {
		return hwaddr
	}
	return net.HardwareAddr{}
}

func (g *L3UnicastGroup) GetEthSrcHwAddr() net.HardwareAddr {
	if hwaddr, err := net.ParseMAC(g.EthSrc); err == nil {
		return hwaddr
	}
	return net.HardwareAddr{}
}

func (g *L3UnicastGroup) GetTunLocalIP() net.IP {
	if ip := net.ParseIP(g.TunLocal); ip != nil {
		return ip
	}
	return net.IP{}
}

func (g *L3UnicastGroup) GetTunRemoteIP() net.IP {
	if ip := net.ParseIP(g.TunRemote); ip != nil {
		return ip
	}
	return net.IP{}
}

func (g *L3UnicastGroup) ToMod(cmd GroupMod_Cmd, reId string) *GroupMod {
	return &GroupMod{
		Cmd:   cmd,
		GType: GroupMod_L3_UNICAST,
		ReId:  reId,
		Entry: &GroupMod_L3Unicast{L3Unicast: g},
	}
}

//
// L2 Multicast Group
//
func NewL2MulticastGroupID(mcId uint16, vlanId uint16) uint32 {
	return 0x30000000 + ((AdjustVlanVID(vlanId) << 16) & 0x0fff0000) + ((uint32)(mcId) & 0xffff)
}

//
// L2 Flood Group
//
func NewL2FloodGroupID(floodId uint16, vlanId uint16) uint32 {
	return 0x40000000 + ((AdjustVlanVID(vlanId) << 16) & 0x0fff0000) + (uint32(floodId) & 0xffff)
}

//
// L3 Interface Group Id
//
func NewL3InterfaceGroupID(neighId uint32) uint32 {
	return 0x50000000 + (neighId & 0x0fffffff)
}

//
// L3 Multicast Group
//
func NewL3MulticastGroupID(mcId uint16, vlanId uint16) uint32 {
	return 0x60000000 + ((AdjustVlanVID(vlanId) << 16) & 0x0fff0000) + (uint32(mcId) & 0xffff)
}

//
// L3 ECMP Group
//
func NewL3EcmpGroupId(ecmpId uint32) uint32 {
	return 0x70000000 + (ecmpId & 0x0fffffff)
}

//
// L2 Overlay Group
//
func NewOverlayGroupID(tunnelId uint16, subType uint16, index uint16) uint32 {
	return 0x80000000 +
		(((uint32)(tunnelId) << 12) & 0x0ffff000) +
		(((uint32)(subType) << 10) & 0x0800) +
		((uint32)(index) & 0x07ff)
}

//
// MPLS Interface Group
//
func NewMPLSInterfaceGroupID(neighId uint32) uint32 {
	return 0x90000000 + (neighId & 0x00ffffff)
}

func NewMPLSInterfaceGroup(neId, portId uint32, vlanVid uint16, ethDst, ethSrc net.HardwareAddr) *MPLSInterfaceGroup {
	return &MPLSInterfaceGroup{
		NeId:    neId,
		PortId:  portId,
		VlanVid: uint32(vlanVid),
		EthDst:  ethDst.String(),
		EthSrc:  ethSrc.String(),
	}
}

func (g *MPLSInterfaceGroup) ToMod(cmd GroupMod_Cmd, reId string) *GroupMod {
	return &GroupMod{
		Cmd:   cmd,
		GType: GroupMod_MPLS_INTERFACE,
		ReId:  reId,
		Entry: &GroupMod_MplsIface{MplsIface: g},
	}
}

//
// MPLS Label Group
//
var MPLSLabelGroup_subtype = map[GroupMod_GType]uint32{
	GroupMod_MPLS_L2_VPN:  1,
	GroupMod_MPLS_L3_VPN:  2,
	GroupMod_MPLS_TUNNEL1: 3,
	GroupMod_MPLS_TUNNEL2: 4,
	GroupMod_MPLS_SWAP:    5,
}

func NewMPLSLabelGroupID(subType uint32, label uint32) uint32 {
	return 0x90000000 + ((subType << 24) & 0x0f000000) + (label & 0x00ffffff)
}

func NewMPLSLabelGroup(gType GroupMod_GType, dstId, newLabel, neId, newDstId uint32) *MPLSLabelGroup {
	return &MPLSLabelGroup{
		GType:    gType,
		DstId:    dstId,
		NewLabel: newLabel,
		NeId:     neId,
		NewDstId: newDstId,
	}
}

func (g *MPLSLabelGroup) ToMod(cmd GroupMod_Cmd, reId string) *GroupMod {
	return &GroupMod{
		Cmd:   cmd,
		GType: g.GType,
		ReId:  reId,
		Entry: &GroupMod_MplsLabel{MplsLabel: g},
	}
}

//
// MPLS FastFailover Group
//
func NewMPLSFastFailoverGroupID(index uint32) uint32 {
	return 0xa6000000 + (index & 0x00ffffff)
}

//
// MPLS ECMP Group
//
func NewMPLSEcmpGroupID(index uint32) uint32 {
	return 0xa8000000 + (index & 0x00ffffff)
}

//
// Unfiltered Interface Group
//
func NewUnfilteredInterfaceGroupID(portId uint32) uint32 {
	return 0xb0000000 + (portId & 0xffff)
}
