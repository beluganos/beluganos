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

package nlaapi

import (
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"syscall"
)

//
// NlMsghdr
//
func (h *NlMsghdr) ToNative() *syscall.NlMsghdr {
	return &syscall.NlMsghdr{
		Len:   h.Len,
		Type:  uint16(h.Type),
		Flags: uint16(h.Flags),
		Seq:   h.Seq,
		Pid:   h.Pid,
	}
}

func NewNlMsghdrFromNative(h *syscall.NlMsghdr) *NlMsghdr {
	return &NlMsghdr{
		Len:   h.Len,
		Type:  uint32(h.Type),
		Flags: uint32(h.Flags),
		Seq:   h.Seq,
		Pid:   h.Pid,
	}
}

//
// NlMsgSrc
//
var NlMsgSrc_to_native = map[NlMsgSrc]nlamsg.NlMsgSrc{
	NlMsgSrc_NOP: nlamsg.SRC_NOP,
	NlMsgSrc_KNL: nlamsg.SRC_KNL,
	NlMsgSrc_API: nlamsg.SRC_API,
}

var NlMsgSrc_from_native = map[nlamsg.NlMsgSrc]NlMsgSrc{
	nlamsg.SRC_NOP: NlMsgSrc_NOP,
	nlamsg.SRC_KNL: NlMsgSrc_KNL,
	nlamsg.SRC_API: NlMsgSrc_API,
}

func (n NlMsgSrc) ToNative() nlamsg.NlMsgSrc {
	return NlMsgSrc_to_native[n]
}

func NewNlMsgSrcFromNative(n nlamsg.NlMsgSrc) NlMsgSrc {
	return NlMsgSrc_from_native[n]
}

//
// NetlinkMessage
//
func (n *NetlinkMessage) Type() uint16 {
	return uint16(n.Header.Type)
}

func (n *NetlinkMessage) ToNative() *nlamsg.NetlinkMessage {
	msg := &syscall.NetlinkMessage{
		Header: *n.GetHeader().ToNative(),
		Data:   n.Data,
	}
	m := nlamsg.NewNetlinkMessage(msg, uint8(n.NId), n.Src.ToNative())
	return m
}

func NewNetlinkMessageFromNative(n *nlamsg.NetlinkMessage) *NetlinkMessage {
	return &NetlinkMessage{
		Header: NewNlMsghdrFromNative(&n.Header),
		Data:   n.Data,
		NId:    uint32(n.NId),
		Src:    NewNlMsgSrcFromNative(n.Src),
	}
}

//
// NetlinkMessageUnion
//
func (n *NlMsgUni) ToNative(g uint16) interface{} {
	switch g {
	case nlalink.RTMGRP_LINK:
		return n.GetLink().ToNative()
	case nlalink.RTMGRP_ADDR:
		return n.GetAddr().ToNative()
	case nlalink.RTMGRP_NEIGH:
		return n.GetNeigh().ToNative()
	case nlalink.RTMGRP_ROUTE:
		return n.GetRoute().ToNative()
	case nlalink.RTMGRP_NODE:
		return n.GetNode().ToNative()
	case nlalink.RTMGRP_VPN:
		return n.GetVpn().ToNative()
	default:
		return nil
	}
}

func NewNlMsgUni(n interface{}, g uint16) *NlMsgUni {
	var msg isNlMsgUni_Msg

	switch g {
	case nlalink.RTMGRP_LINK:
		msg = &NlMsgUni_Link{Link: n.(*Link)}

	case nlalink.RTMGRP_ADDR:
		msg = &NlMsgUni_Addr{Addr: n.(*Addr)}

	case nlalink.RTMGRP_NEIGH:
		msg = &NlMsgUni_Neigh{Neigh: n.(*Neigh)}

	case nlalink.RTMGRP_ROUTE:
		msg = &NlMsgUni_Route{Route: n.(*Route)}

	case nlalink.RTMGRP_NODE:
		msg = &NlMsgUni_Node{Node: n.(*Node)}

	case nlalink.RTMGRP_VPN:
		msg = &NlMsgUni_Vpn{Vpn: n.(*Vpn)}

	default:
		return nil
	}
	return &NlMsgUni{
		Msg: msg,
	}
}

func NewNetlinkMessageUnion(nid uint8, t uint16, m interface{}) *NetlinkMessageUnion {
	g := nlamsg.NlMsgGroupFromType(t)
	return &NetlinkMessageUnion{
		Header: &NlMsghdr{Type: uint32(t)},
		Msg:    NewNlMsgUni(m, g),
		NId:    uint32(nid),
	}
}

func NewNlMsgUniFromNative(n interface{}, g uint16) *NlMsgUni {
	switch g {
	case nlalink.RTMGRP_LINK:
		return NewNlMsgUni(NewLinkFromNative(n.(*nlamsg.Link)), g)

	case nlalink.RTMGRP_ADDR:
		return NewNlMsgUni(NewAddrFromNative(n.(*nlamsg.Addr)), g)

	case nlalink.RTMGRP_NEIGH:
		return NewNlMsgUni(NewNeighFromNative(n.(*nlamsg.Neigh)), g)

	case nlalink.RTMGRP_ROUTE:
		return NewNlMsgUni(NewRouteFromNative(n.(*nlamsg.Route)), g)

	case nlalink.RTMGRP_NODE:
		return NewNlMsgUni(NewNodeFromNative(n.(*nlamsg.Node)), g)

	case nlalink.RTMGRP_VPN:
		return NewNlMsgUni(NewVpnFromNative(n.(*nlamsg.Vpn)), g)

	default:
		return nil
	}
}

func (n *NetlinkMessageUnion) Type() uint16 {
	return uint16(n.Header.Type)
}

func (n *NetlinkMessageUnion) Group() uint16 {
	return nlamsg.NlMsgGroupFromType(uint16(n.Header.Type))
}

func (n *NetlinkMessageUnion) ToNative() *nlamsg.NetlinkMessageUnion {
	return &nlamsg.NetlinkMessageUnion{
		Header: *n.GetHeader().ToNative(),
		Msg:    n.GetMsg().ToNative(n.Group()),
		NId:    uint8(n.NId),
		Src:    n.Src.ToNative(),
	}
}

func NewNetlinkMessageUnionFromNative(n *nlamsg.NetlinkMessageUnion) *NetlinkMessageUnion {
	return &NetlinkMessageUnion{
		Header: NewNlMsghdrFromNative(&n.Header),
		Msg:    NewNlMsgUniFromNative(n.Msg, n.Group()),
		NId:    uint32(n.NId),
		Src:    NewNlMsgSrcFromNative(n.Src),
	}
}
