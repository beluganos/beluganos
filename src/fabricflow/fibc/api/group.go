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
	"github.com/golang/protobuf/proto"
	"net"
)

func (g *GroupMod) Type() uint16 {
	return uint16(FFM_GROUP_MOD)
}

func (g *GroupMod) Bytes() ([]byte, error) {
	return proto.Marshal(g)
}

//
// L2 Interface Group
//
func NewL2InterfaceGroup(portId uint32, vlanId uint16, vlanTranslation bool) *L2InterfaceGroup {
	return &L2InterfaceGroup{
		PortId:          portId,
		VlanVid:         uint32(vlanId),
		VlanTranslation: vlanTranslation,
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

//
// L2 Rewrite Group Id
//

//
// L3 Unicast Group
//
func NewL3UnicastGroup(neId, portId uint32, vlanVid uint16, ethDst, ethSrc net.HardwareAddr) *L3UnicastGroup {
	return &L3UnicastGroup{
		NeId:    neId,
		PortId:  portId,
		VlanVid: uint32(vlanVid),
		EthDst:  ethDst.String(),
		EthSrc:  ethSrc.String(),
	}
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
// L2 Multicast Group Id
//

//
// L2 Flood Group Id
//

//
// L3 Interface Group Id
//

//
// L3 Multicast Group Id
//

//
// L3 ECMP Group Id
//

//
// L2 Overlay Group Id
//

//
// MPLS Interface Group Id
//
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
