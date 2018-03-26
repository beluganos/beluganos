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
	"github.com/vishvananda/netlink/nl"
	"gonla/nladbm"
	"gonla/nlamsg"
	"net"
)

//
// NexthopInfo
//
func (n *NexthopInfo) NetGw() net.IP {
	return net.IP(n.Gw)
}

func (n *NexthopInfo) ToNative() *netlink.NexthopInfo {
	return &netlink.NexthopInfo{
		LinkIndex: int(n.LinkIndex),
		Hops:      int(n.Hops),
		Gw:        n.NetGw(),
		Flags:     int(n.Flags),
		NewDst:    n.NewDst.ToNative(),
		Encap:     n.Encap.ToNative(),
	}
}

func NewNexthopInfoToNatives(ns []*NexthopInfo) []*netlink.NexthopInfo {
	ret := make([]*netlink.NexthopInfo, len(ns))
	for i, n := range ns {
		ret[i] = n.ToNative()
	}
	return ret
}

func NewNexthopInfoFromNative(n *netlink.NexthopInfo) *NexthopInfo {
	return &NexthopInfo{
		LinkIndex: int32(n.LinkIndex),
		Hops:      int32(n.Hops),
		Gw:        n.Gw,
		Flags:     int32(n.Flags),
		NewDst:    NewDestinationFromNative(n.NewDst),
		Encap:     NewEncapFromNative(n.Encap),
	}
}

func NewNexthopInfoFromNatives(ns []*netlink.NexthopInfo) []*NexthopInfo {
	ret := make([]*NexthopInfo, len(ns))
	for i, n := range ns {
		ret[i] = NewNexthopInfoFromNative(n)
	}
	return ret
}

func NewMPLSLabelToNatives(labels []uint32) []int {
	n := make([]int, len(labels))
	for i, label := range labels {
		n[i] = int(label)
	}
	return n
}

func NewMPLSLabelFromNatives(labels []int) []uint32 {
	n := make([]uint32, len(labels))
	for i, label := range labels {
		n[i] = uint32(label)
	}
	return n
}

//
// MPLSDestination
//
func (d *MPLSDestination) ToNative() *netlink.MPLSDestination {
	if d == nil {
		return nil
	}

	return &netlink.MPLSDestination{
		Labels: NewMPLSLabelToNatives(d.Labels),
	}
}

func NewMPLSDestinationFromNative(n *netlink.MPLSDestination) *MPLSDestination {
	if n == nil {
		return nil
	}
	return &MPLSDestination{
		Labels: NewMPLSLabelFromNatives(n.Labels),
	}
}

//
// Destination
//
func (d *Destination) ToNative() netlink.Destination {
	if d == nil {
		return nil
	}

	switch d.Family {
	case nl.FAMILY_MPLS:
		dest := d.Dest.(*Destination_Mpls)
		return dest.Mpls.ToNative()
	default:
		return nil
	}
}

func NewDestinationFromNative(n netlink.Destination) *Destination {
	if n == nil {
		return nil
	}

	switch n.Family() {
	case nl.FAMILY_MPLS:
		mpls := n.(*netlink.MPLSDestination)
		return &Destination{
			Family: nl.FAMILY_MPLS,
			Dest: &Destination_Mpls{
				Mpls: NewMPLSDestinationFromNative(mpls),
			},
		}
	default:
		return nil
	}
}

//
// MPLSEncap
//
func (e *MPLSEncap) ToNative() *netlink.MPLSEncap {
	if e == nil {
		return nil
	}
	return &netlink.MPLSEncap{
		Labels: NewMPLSLabelToNatives(e.Labels),
	}
}

func NewMPLSEncapFromNatives(n *netlink.MPLSEncap) *MPLSEncap {
	if n == nil {
		return nil
	}
	return &MPLSEncap{
		Labels: NewMPLSLabelFromNatives(n.Labels),
	}
}

//
// Encap
//
func (e *Encap) ToNative() netlink.Encap {
	if e == nil {
		return nil
	}

	switch e.Type {
	case nl.LWTUNNEL_ENCAP_MPLS:
		encap := e.Encap.(*Encap_Mpls)
		return encap.Mpls.ToNative()
	default:
		return nil
	}
}

func NewEncapFromNative(n netlink.Encap) *Encap {
	if n == nil {
		return nil
	}

	switch n.Type() {
	case nl.LWTUNNEL_ENCAP_MPLS:
		mpls := n.(*netlink.MPLSEncap)
		return &Encap{
			Type: nl.LWTUNNEL_ENCAP_MPLS,
			Encap: &Encap_Mpls{
				Mpls: NewMPLSEncapFromNatives(mpls),
			},
		}
	default:
		return nil
	}
}

//
// Route
//
func (r *Route) NetScope() netlink.Scope {
	return netlink.Scope(r.Scope)
}

func (r *Route) NetDst() *net.IPNet {
	return BytesToIPNet(r.Dst, r.DstMask)
}

func (r *Route) NetSrc() net.IP {
	return net.IP(r.Src)
}

func (r *Route) NetGw() net.IP {
	return net.IP(r.Gw)
}

func (r *Route) NetVpnGw() net.IP {
	return net.IP(r.VpnGw)
}

func (r *Route) ToNetlink() *netlink.Route {
	mplsDst := func(v int32) *int {
		if v == -1 {
			return nil
		}
		dst := int(v)
		return &dst
	}
	return &netlink.Route{
		LinkIndex:  int(r.LinkIndex),
		ILinkIndex: int(r.ILinkIndex),
		Scope:      r.NetScope(),
		Dst:        r.NetDst(),
		Src:        r.NetSrc(),
		Gw:         r.NetGw(),
		MultiPath:  NewNexthopInfoToNatives(r.MultiPath),
		Protocol:   int(r.Protocol),
		Priority:   int(r.Priority),
		Table:      int(r.Table),
		Type:       int(r.Type),
		Tos:        int(r.Tos),
		Flags:      int(r.Flags),
		MPLSDst:    mplsDst(r.MplsDst),
		NewDst:     r.NewDst.ToNative(),
		Encap:      r.Encap.ToNative(),
	}
}

func (r *Route) ToNative() *nlamsg.Route {
	return &nlamsg.Route{
		Route: r.ToNetlink(),
		NId:   uint8(r.NId),
		RtId:  r.RtId,
		VpnGw: r.NetVpnGw(),
		EnIds: r.EnIds,
	}
}

func NewRouteFromNative(r *nlamsg.Route) *Route {
	dst, dstMask := IPNetToBytes(r.Dst)
	mplsDst := func(v *int) int32 {
		if v == nil {
			return -1
		}
		return int32(*v)
	}
	return &Route{
		LinkIndex:  int32(r.LinkIndex),
		ILinkIndex: int32(r.ILinkIndex),
		Scope:      int32(r.Scope),
		Dst:        dst,
		DstMask:    dstMask,
		Src:        r.Src,
		Gw:         r.Gw,
		MultiPath:  NewNexthopInfoFromNatives(r.MultiPath),
		Protocol:   int32(r.Protocol),
		Priority:   int32(r.Priority),
		Table:      int32(r.Table),
		Type:       int32(r.Type),
		Tos:        int32(r.Tos),
		Flags:      int32(r.Flags),
		MplsDst:    mplsDst(r.MPLSDst),
		NewDst:     NewDestinationFromNative(r.NewDst),
		Encap:      NewEncapFromNative(r.Encap),
		NId:        uint32(r.NId),
		RtId:       r.RtId,
		VpnGw:      r.VpnGw,
		EnIds:      r.EnIds,
	}
}

//
// Route (Key)
//
func (k *RouteKey) ToNative() *nladbm.RouteKey {
	return &nladbm.RouteKey{
		NId:  uint8(k.NId),
		Addr: k.Addr,
	}
}

func NewRouteKeyFromNative(n *nladbm.RouteKey) *RouteKey {
	return &RouteKey{
		NId:  uint32(n.NId),
		Addr: n.Addr,
	}
}

//
// Routes
//
func NewGetRoutesRequest(nid uint8) *GetRoutesRequest {
	return &GetRoutesRequest{
		NId: uint32(nid),
	}
}
