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

package nlamsg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gonla/nlalib"
	"syscall"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

const BRVLAN_INFO_IFNAME_SIZE = 16

type BridgeVlanPortType uint8

const (
	BRIDGE_VLAN_PORT_NONE   = 0
	BRIDGE_VLAN_PORT_ACCESS = 1
	BRIDGE_VLAN_PORT_TRUNK  = 2
	BRIDGE_VLAN_PORT_MASTER = 3
)

var bridgeVlanPortType_names = map[BridgeVlanPortType]string{
	BRIDGE_VLAN_PORT_NONE:   "none",
	BRIDGE_VLAN_PORT_ACCESS: "access",
	BRIDGE_VLAN_PORT_TRUNK:  "trunk",
	BRIDGE_VLAN_PORT_MASTER: "master",
}

var bridgeVlanPortType_values = map[string]BridgeVlanPortType{
	"none":   BRIDGE_VLAN_PORT_NONE,
	"access": BRIDGE_VLAN_PORT_ACCESS,
	"trunk":  BRIDGE_VLAN_PORT_TRUNK,
	"master": BRIDGE_VLAN_PORT_MASTER,
}

func (t BridgeVlanPortType) String() string {
	if s, ok := bridgeVlanPortType_names[t]; ok {
		return s
	}
	return fmt.Sprintf("BridgeVlanPortType(%d)", t)
}

func ParseBridgeVlanPortType(s string) (BridgeVlanPortType, error) {
	if v, ok := bridgeVlanPortType_values[s]; ok {
		return v, nil
	}
	return BRIDGE_VLAN_PORT_NONE, fmt.Errorf("Invalid BridgeVlanPortType. %s", s)
}

type BridgeVlanInfo struct {
	nl.BridgeVlanInfo
	Index       int
	Name        string
	MasterIndex int
	Mtu         uint16
	BrId        uint32
	NId         uint8
}

func NewBridgeVlanInfo(bridgeVlan nl.BridgeVlanInfo, index int, name string, master int, nid uint8, id uint32) *BridgeVlanInfo {
	return &BridgeVlanInfo{
		BridgeVlanInfo: bridgeVlan,
		Index:          index,
		Name:           name,
		MasterIndex:    master,
		Mtu:            0,
		BrId:           id,
		NId:            nid,
	}
}

func NewBridgeVlanInfoFromNetlink(nid uint8, brvlan *nl.BridgeVlanInfo, link netlink.Link) *BridgeVlanInfo {
	return NewBridgeVlanInfo(
		*brvlan,
		link.Attrs().Index,
		link.Attrs().Name,
		link.Attrs().MasterIndex,
		nid,
		0, // BrId
	)
}

func (b *BridgeVlanInfo) Copy() *BridgeVlanInfo {
	return &BridgeVlanInfo{
		BridgeVlanInfo: b.BridgeVlanInfo,
		Index:          b.Index,
		Name:           b.Name,
		MasterIndex:    b.MasterIndex,
		Mtu:            b.Mtu,
		BrId:           b.BrId,
		NId:            b.NId,
	}
}

func (b *BridgeVlanInfo) Equals(other *BridgeVlanInfo) bool {
	if b == nil || other == nil {
		return false
	}
	return (*b) == (*other)
}

func (b *BridgeVlanInfo) PortType() BridgeVlanPortType {
	if m := b.MasterIndex; m == 0 || m == b.Index {
		return BRIDGE_VLAN_PORT_MASTER
	}

	if (b.Flags&nl.BRIDGE_VLAN_INFO_PVID) != 0 && (b.Flags&nl.BRIDGE_VLAN_INFO_UNTAGGED) != 0 {
		return BRIDGE_VLAN_PORT_ACCESS
	}

	if (b.Flags&nl.BRIDGE_VLAN_INFO_PVID) == 0 && (b.Flags&nl.BRIDGE_VLAN_INFO_UNTAGGED) == 0 {
		return BRIDGE_VLAN_PORT_TRUNK
	}

	return BRIDGE_VLAN_PORT_NONE
}

func (b *BridgeVlanInfo) String() string {
	flags := nlalib.StringBridgeVlanInfoFlags(b.Flags)
	return fmt.Sprintf("{BridgeVlanInfo: Index:%d, Name:'%s' Master:%d Vlan:%d Mtu:%d flags:%s type:%s, BrId:%d NId:%d}", b.Index, b.Name, b.MasterIndex, b.Vid, b.Mtu, flags, b.PortType(), b.BrId, b.NId)
}

type bridgeVlanInfo struct {
	nl.BridgeVlanInfo                               // [4]
	Index             int32                         // [4]
	Name              [BRVLAN_INFO_IFNAME_SIZE]byte // [16]
	MasterIndex       int32                         // [4]
	Mtu               uint32                        // [4]
}

func (b *BridgeVlanInfo) Bytes() ([]byte, error) {
	data := bridgeVlanInfo{}

	data.BridgeVlanInfo = b.BridgeVlanInfo
	data.Index = int32(b.Index)
	copy(data.Name[:], b.Name)
	data.MasterIndex = int32(b.MasterIndex)
	data.Mtu = uint32(b.Mtu)

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, &data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b *BridgeVlanInfo) ToNetlinkMessage(msgType uint16) (*syscall.NetlinkMessage, error) {
	data, err := b.Bytes()
	if err != nil {
		return nil, err
	}
	return nlalib.NewNetlinkMessage(msgType, data), nil
}

func ParseBridgeVlanInfo(b []byte) (*BridgeVlanInfo, error) {
	data := bridgeVlanInfo{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.BigEndian, &data); err != nil {
		return nil, err
	}

	br := BridgeVlanInfo{}

	name := func() string {
		if index := bytes.IndexByte(data.Name[:], 0x00); index >= 0 {
			return string(data.Name[:index])
		}
		return string(data.Name[:])
	}()

	br.BridgeVlanInfo = data.BridgeVlanInfo
	br.Index = int(data.Index)
	br.Name = name
	br.MasterIndex = int(data.MasterIndex)
	br.Mtu = uint16(data.Mtu)

	return &br, nil
}

func BridgeVlanInfoDeserialize(nlmsg *NetlinkMessage) (*BridgeVlanInfo, error) {
	brvlan, err := ParseBridgeVlanInfo(nlmsg.Data)
	if err != nil {
		return nil, err
	}
	brvlan.NId = nlmsg.NId

	return brvlan, nil
}

//
// BridgeVlanInfoSerialize coverts BridgeVlanInfo struct to netlink message.
// msgType is RTM_NEWBRIDGE, RTM_DELBRIDGE, RTM_SETBRIDGE (gonla/nlamsg/nlalink package)
//
func BridgeVlanInfoSerialize(brvlan *BridgeVlanInfo, msgType uint16) (*NetlinkMessage, error) {
	nlmsg, err := brvlan.ToNetlinkMessage(msgType)
	if err != nil {
		return nil, err
	}
	return NewNetlinkMessage(nlmsg, brvlan.NId, SRC_NOP), nil
}
