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

package ofproto

import (
	"fmt"
)

type NameConv func(uint32) string

var groupId_converts = map[uint32]NameConv{
	0:  ConvL2Interface,
	1:  ConvL2Rewrite,
	2:  ConvL3Unicast,
	3:  ConvL2Multicast,
	4:  ConvL2Flood,
	5:  ConvL3Interface,
	6:  ConvL3Multicast,
	7:  ConvL3Ecmp,
	8:  ConvL2DCOverlay,
	9:  ConvMPLSLabel,
	10: ConvMPLSForward,
	11: ConvUnfilteredInterface,
}

func ConvGroupId(gID uint32) string {
	gType := (gID & 0xf0000000) >> 28
	if conv, ok := groupId_converts[gType]; ok {
		return conv(gID)
	}
	return fmt.Sprintf("(%08x)", gID)
}

func ConvL2Interface(gID uint32) string {
	vid := (gID & 0x0fff0000) >> 16
	port := gID & 0xffff
	return fmt.Sprintf("L2_IFACE,vid=%d,port=%d", vid, port)
}

func ConvL2Rewrite(gID uint32) string {
	id := gID & 0x0fffffff
	return fmt.Sprintf("L2_REWRT,id=0x%x", id)
}

func ConvL3Unicast(gID uint32) string {
	vrf := (gID & 0x00ff0000) >> 16
	neid := gID & 0x0000ffff
	return fmt.Sprintf("L3_U.C.,vrf=%d,nei=%d", vrf, neid)
}

func ConvL2Multicast(gID uint32) string {
	id := gID & 0x0fffffff
	return fmt.Sprintf("L2_M.C.,id=0x%x", id)
}

func ConvL2Flood(gID uint32) string {
	vid := (gID & 0x0fff0000) >> 16
	id := gID & 0x0000ffff
	return fmt.Sprintf("L2_FLOD,vid=%d,id=0x%x", vid, id)
}

func ConvL3Interface(gID uint32) string {
	id := gID & 0x0fffffff
	return fmt.Sprintf("L3_IFACE,id=0x%x", id)
}

func ConvL3Multicast(gID uint32) string {
	vid := (gID & 0x0fff0000) >> 16
	index := gID & 0x0000ffff
	return fmt.Sprintf("L3_M.C.,vid=%d,idx=0x%x", vid, index)
}

func ConvL3Ecmp(gID uint32) string {
	id := gID & 0x0fffffff
	return fmt.Sprintf("L3_ECMP,id=0x%x", id)
}

func ConvL2DCOverlay(gID uint32) string {
	id := gID & 0x0fffffff
	return fmt.Sprintf("L2_D.C.,id=0x%x", id)
}

func ConvMPLSLabel(gID uint32) string {
	subType := (gID & 0x0f000000) >> 24
	if conv, ok := mplsGroup_convets[subType]; ok {
		return conv(gID)
	}
	return fmt.Sprintf("MPLS_LBL,0xx%08x", gID)
}

func ConvMPLSForward(gID uint32) string {
	return ""
}

func ConvUnfilteredInterface(gID uint32) string {
	port := gID & 0x0000ffff
	return fmt.Sprintf("L2_UF-IF,port=%d", port)
}

var mplsGroup_convets = map[uint32]NameConv{
	0: ConvMPLSInterface, // "MPLS_IFACE",
	1: ConvMPLSL2VPN,     // "MPLS_L2VPN",
	2: ConvMPLSL3VPN,     // "MPLS_L3VPN",
	3: ConvMPLSTun1,      // "MPLS_TUN1",
	4: ConvMPLSTun2,      // "MPLS_TUN2",
	5: ConvMPLSSwap,      // "MPLS_SWAP",
}

func ConvMPLSInterface(gID uint32) string {
	vrf := (gID & 0x00ff0000) >> 16
	neid := gID & 0x0000ffff
	return fmt.Sprintf("MPLS_IFACE,vrf=%d,nei=%d", vrf, neid)
}

func ConvMPLSL2VPN(gID uint32) string {
	index := gID & 0x00ffffff
	return fmt.Sprintf("MPLS_L2VPN,idx=0x%08x", index)
}

func ConvMPLSL3VPN(gID uint32) string {
	index := (gID & 0x00f00000) >> 20
	enc := gID & 0x000fffff
	return fmt.Sprintf("MPLS_L3VPN,i=%d,enc=%d", index, enc)
}

func ConvMPLSTun1(gID uint32) string {
	index := (gID & 0x00f00000) >> 20
	enc := gID & 0x000fffff
	return fmt.Sprintf("MPLS_TUN1,i=%d,enc=%d", index, enc)
}

func ConvMPLSTun2(gID uint32) string {
	index := (gID & 0x00f00000) >> 20
	enc := gID & 0x000fffff
	return fmt.Sprintf("MPLS_TUN1,i=%d,enc=%d", index, enc)
}

func ConvMPLSSwap(gID uint32) string {
	index := (gID & 0x00f00000) >> 20
	label := gID & 0x000fffff
	return fmt.Sprintf("MPLS_SWAP,i=%d,label=%d", index, label)
}
