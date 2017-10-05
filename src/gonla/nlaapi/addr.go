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
	"github.com/vishvananda/netlink"
	"gonla/nladbm"
	"gonla/nlamsg"
	"net"
)

//
// Addr
//
func (a *Addr) GetIPNet() *net.IPNet {
	return BytesToIPNet(a.Ip, a.IpMask)
}

func (a *Addr) NetPeer() *net.IPNet {
	return BytesToIPNet(a.Peer, a.PeerMask)
}

func (a *Addr) NetBroadcast() net.IP {
	return net.IP(a.Broadcast)
}

func (a *Addr) ToNetlink() *netlink.Addr {
	return &netlink.Addr{
		IPNet:     a.GetIPNet(),
		Label:     a.Label,
		Flags:     int(a.Flags),
		Scope:     int(a.Scope),
		Peer:      a.NetPeer(),
		Broadcast: a.NetBroadcast(),
	}
}

func (a *Addr) ToNative() *nlamsg.Addr {
	return &nlamsg.Addr{
		Addr:   a.ToNetlink(),
		AdId:   a.AdId,
		Family: a.Family,
		Index:  a.Index,
		NId:    uint8(a.NId),
	}
}

func NewAddrFromNative(n *nlamsg.Addr) *Addr {
	ip, ipMask := IPNetToBytes(n.IPNet)
	peer, peerMask := IPNetToBytes(n.Peer)
	return &Addr{
		Ip:        ip,
		IpMask:    ipMask,
		Label:     n.Label,
		Flags:     int32(n.Flags),
		Scope:     int32(n.Scope),
		Peer:      peer,
		PeerMask:  peerMask,
		Broadcast: n.Broadcast,
		Family:    n.Family,
		Index:     n.Index,
		NId:       uint32(n.NId),
		AdId:      n.AdId,
	}
}

//
// Addr (Key)
//
func (k *AddrKey) ToNative() *nladbm.AddrKey {
	return &nladbm.AddrKey{
		NId:  uint8(k.NId),
		Addr: k.Addr,
	}
}

func NewAddrKeyFromNative(n *nladbm.AddrKey) *AddrKey {
	return &AddrKey{
		NId:  uint32(n.NId),
		Addr: n.Addr,
	}
}
