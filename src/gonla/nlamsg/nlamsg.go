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

package nlamsg

import (
	"fmt"
	"gonla/nlamsg/nlalink"
	"syscall"
)

//
// Netlink type to RTMGRP
//
func NlMsgGroupFromType(nlmsgType uint16) uint16 {
	switch nlmsgType {
	case syscall.RTM_NEWLINK, syscall.RTM_DELLINK, syscall.RTM_SETLINK:
		return nlalink.RTMGRP_LINK
	case syscall.RTM_NEWADDR, syscall.RTM_DELADDR, nlalink.RTM_SETADDR:
		return nlalink.RTMGRP_ADDR
	case syscall.RTM_NEWNEIGH, syscall.RTM_DELNEIGH, nlalink.RTM_SETNEIGH:
		return nlalink.RTMGRP_NEIGH
	case syscall.RTM_NEWROUTE, syscall.RTM_DELROUTE, nlalink.RTM_SETROUTE:
		return nlalink.RTMGRP_ROUTE
	case nlalink.RTM_NEWNODE, nlalink.RTM_DELNODE, nlalink.RTM_SETNODE:
		return nlalink.RTMGRP_NODE
	case nlalink.RTM_NEWVPN, nlalink.RTM_DELVPN, nlalink.RTM_SETVPN:
		return nlalink.RTMGRP_VPN
	case nlalink.RTM_NEWBRIDGE, nlalink.RTM_DELBRIDGE, nlalink.RTM_SETBRIDGE:
		return nlalink.RTMGRP_BRIDGE
	default:
		return nlalink.RTMGRP_UNSPEC
	}
}

type NlMsgSrc int32

const (
	SRC_NOP NlMsgSrc = 0
	SRC_KNL NlMsgSrc = 1
	SRC_API NlMsgSrc = 2
)

var NlMsgSrc_name = map[NlMsgSrc]string{
	SRC_NOP: "NOP",
	SRC_KNL: "KNL",
	SRC_API: "API",
}

func (n NlMsgSrc) String() string {
	if s, ok := NlMsgSrc_name[n]; ok {
		return s
	}
	return fmt.Sprintf("NlMsgSrc(%d)", n)
}

//
// NetlinkMessage
//
type NetlinkMessage struct {
	syscall.NetlinkMessage
	NId uint8
	Src NlMsgSrc
}

func (m *NetlinkMessage) Type() uint16 {
	return m.Header.Type
}

func (m *NetlinkMessage) Group() uint16 {
	return NlMsgGroupFromType(m.Header.Type)
}

func (m *NetlinkMessage) Len() uint32 {
	return m.Header.Len
}

func (m *NetlinkMessage) String() string {
	return fmt.Sprintf("%s NId:%d Src:%s", NlMsgHdrStr(&m.Header), m.NId, m.Src)
}

func NewNetlinkMessage(msg *syscall.NetlinkMessage, nid uint8, src NlMsgSrc) *NetlinkMessage {
	return &NetlinkMessage{
		NetlinkMessage: *msg,
		NId:            nid,
		Src:            src,
	}
}

type NetlinkMessageUnion struct {
	Header syscall.NlMsghdr
	Msg    interface{}
	NId    uint8
	Src    NlMsgSrc
}

func NewNetlinkMessageUnion(header *syscall.NlMsghdr, msg interface{}, nid uint8, src NlMsgSrc) *NetlinkMessageUnion {
	return &NetlinkMessageUnion{
		Header: *header,
		Msg:    msg, // netlink.Link, nlamsg.Node, ...
		NId:    nid,
		Src:    src,
	}
}

func (m *NetlinkMessageUnion) String() string {
	return fmt.Sprintf("%s NId:%d Src:%s %v", NlMsgHdrStr(&m.Header), m.NId, m.Src, m.Msg)
}

func (m *NetlinkMessageUnion) Type() uint16 {
	return m.Header.Type
}

func (m *NetlinkMessageUnion) Group() uint16 {
	return NlMsgGroupFromType(m.Header.Type)
}

func (n *NetlinkMessageUnion) GetLink() *Link {
	if v, ok := n.Msg.(*Link); ok {
		return v
	}
	return nil
}

func (n *NetlinkMessageUnion) GetAddr() *Addr {
	if v, ok := n.Msg.(*Addr); ok {
		return v
	}
	return nil
}

func (n *NetlinkMessageUnion) GetNeigh() *Neigh {
	if v, ok := n.Msg.(*Neigh); ok {
		return v
	}
	return nil
}

func (n *NetlinkMessageUnion) GetRoute() *Route {
	if v, ok := n.Msg.(*Route); ok {
		return v
	}
	return nil
}

func (n *NetlinkMessageUnion) GetNode() *Node {
	if v, ok := n.Msg.(*Node); ok {
		return v
	}
	return nil
}

func (n *NetlinkMessageUnion) GetVpn() *Vpn {
	if v, ok := n.Msg.(*Vpn); ok {
		return v
	}
	return nil
}

func (n *NetlinkMessageUnion) GetBridgeVlanInfo() *BridgeVlanInfo {
	if v, ok := n.Msg.(*BridgeVlanInfo); ok {
		return v
	}
	return nil
}
