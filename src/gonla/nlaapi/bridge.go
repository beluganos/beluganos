// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package nlaapi

import (
	"gonla/nladbm"
	"gonla/nlamsg"

	"github.com/vishvananda/netlink/nl"
)

//
// BridgeVlanInfo
//
func NewBridgeVlanInfo(nid uint8, vid uint16, ifindex int, masterIndex int) *BridgeVlanInfo {
	return &BridgeVlanInfo{
		NId:         uint32(nid),
		BrId:        0,
		Flags:       0,
		Vid:         uint32(vid),
		Index:       int32(ifindex),
		Name:        "",
		MasterIndex: int32(masterIndex),
		Mtu:         0,
	}
}

func NewBridgeVlanInfoFromNative(v *nlamsg.BridgeVlanInfo) *BridgeVlanInfo {
	return &BridgeVlanInfo{
		NId:         uint32(v.NId),
		BrId:        v.BrId,
		Flags:       BridgeVlanInfo_Flags(v.Flags),
		Vid:         uint32(v.Vid),
		Index:       int32(v.Index),
		Name:        v.Name,
		MasterIndex: int32(v.MasterIndex),
		Mtu:         uint32(v.Mtu),
	}
}

func (v *BridgeVlanInfo) ToNative() *nlamsg.BridgeVlanInfo {
	return &nlamsg.BridgeVlanInfo{
		NId:  uint8(v.NId),
		BrId: v.BrId,
		BridgeVlanInfo: nl.BridgeVlanInfo{
			Flags: uint16(v.Flags),
			Vid:   uint16(v.Vid),
		},
		Index:       int(v.Index),
		Name:        v.Name,
		MasterIndex: int(v.MasterIndex),
		Mtu:         uint16(v.Mtu),
	}
}

func (v *BridgeVlanInfo) PortType() BridgeVlanInfo_PortType {
	if v.MasterIndex == v.Index {
		return BridgeVlanInfo_MASTER_PORT
	}

	if (v.Flags&BridgeVlanInfo_PVID) != 0 && (v.Flags&BridgeVlanInfo_UNTAGGED) != 0 {
		return BridgeVlanInfo_ACCESS_PORT
	}

	if (v.Flags&BridgeVlanInfo_PVID) == 0 && (v.Flags&BridgeVlanInfo_UNTAGGED) == 0 {
		return BridgeVlanInfo_TRUNK_PORT
	}

	return BridgeVlanInfo_NONE_PORT
}

//
// BridgeVlanInfo (Key)
//
func (k *BridgeVlanInfoKey) ToNative() *nladbm.BridgeVlanInfoKey {
	return &nladbm.BridgeVlanInfoKey{
		NId:   uint8(k.NId),
		Index: int(k.Index),
		Vid:   uint16(k.Vid),
	}
}

func NewBridgeVlanInfoKeyFromNative(n *nladbm.BridgeVlanInfoKey) *BridgeVlanInfoKey {
	return &BridgeVlanInfoKey{
		NId:   uint32(n.NId),
		Index: int32(n.Index),
		Vid:   uint32(n.Vid),
	}
}

func NewGetBridgeVlanInfosRequest(nid uint8) *GetBridgeVlanInfosRequest {
	return &GetBridgeVlanInfosRequest{
		NId: uint32(nid),
	}
}
