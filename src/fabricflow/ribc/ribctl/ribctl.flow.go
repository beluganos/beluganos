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

package ribctl

import (
	"fabricflow/fibc/api"
	"gonla/nlamsg"
	"net"
)

func checkIPNet(ip *net.IPNet) bool {
	if ip == nil {
		return false
	}
	return checkIP(ip.IP)
}

func checkIP(ip net.IP) bool {
	switch {
	case ip == nil:
		return false
	case ip.IsInterfaceLocalMulticast():
		return false
	case ip.IsLinkLocalMulticast():
		return false
	case ip.IsLoopback():
		return false
	case ip.IsMulticast():
		return false
	case ip.IsUnspecified():
		return false
	default:
		return true
	}
}

//
// VLAN Flow()
//
func NewVLANFilterFlow(link *nlamsg.Link) *fibcapi.VLANFlow {
	vid := link.VlanId()
	if vid == fibcapi.OFPVID_NONE {
		vid = fibcapi.OFPVID_UNTAGGED
	}

	m := fibcapi.NewVLANFlowMatch(
		NewPortId(link),
		uint32(vid),
		fibcapi.OFPVID_ABSENT,
	)
	a := []*fibcapi.VLANFlow_Action{}
	if link.NId != 0 {
		a = append(a, fibcapi.NewVLANFlowAction("SET_VRF", uint32(link.NId)))
	}

	return fibcapi.NewVLANFlow(m, a, uint32(fibcapi.FlowMod_TERM_MAC))
}

func NewVLANUntagFlow(link *nlamsg.Link) *fibcapi.VLANFlow {
	m := fibcapi.NewVLANFlowMatch(
		NewPortId(link),
		fibcapi.OFPVID_NONE,
		fibcapi.OFPVID_ABSENT,
	)
	a := []*fibcapi.VLANFlow_Action{
		fibcapi.NewVLANFlowAction("PUSH_VLAN", fibcapi.OFPVID_UNTAGGED),
	}
	if link.NId != 0 {
		a = append(a, fibcapi.NewVLANFlowAction("SET_VRF", uint32(link.NId)))
	}

	return fibcapi.NewVLANFlow(m, a, uint32(fibcapi.FlowMod_TERM_MAC))
}

func (r *RIBController) SendVLANFlow(cmd fibcapi.FlowMod_Cmd, link *nlamsg.Link) error {
	f := NewVLANFilterFlow(link)
	if err := r.fib.Send(f.ToMod(cmd, r.reId), 0); err != nil {
		return err
	}

	if link.VlanId() == fibcapi.OFPVID_NONE {
		f := NewVLANUntagFlow(link)
		return r.fib.Send(f.ToMod(cmd, r.reId), 0)
	}

	return nil
}

//
// Term MAC floe
//
func NewTermMACFlowIPv4(link *nlamsg.Link) *fibcapi.TerminationMacFlow {
	m := fibcapi.NewTermMACMatch(
		NewPortId(link),
		fibcapi.ETHTYPE_IPV4,
		link.Attrs().HardwareAddr.String(),
		link.VlanId(),
	)
	a := []*fibcapi.TerminationMacFlow_Action{}
	return fibcapi.NewTermMACFlow(m, a, uint32(fibcapi.FlowMod_UNICAST_ROUTING))
}

func NewTermMACFlowIPv6(link *nlamsg.Link) *fibcapi.TerminationMacFlow {
	m := fibcapi.NewTermMACMatch(
		NewPortId(link),
		fibcapi.ETHTYPE_IPV6,
		link.Attrs().HardwareAddr.String(),
		link.VlanId(),
	)
	a := []*fibcapi.TerminationMacFlow_Action{}
	return fibcapi.NewTermMACFlow(m, a, uint32(fibcapi.FlowMod_UNICAST_ROUTING))
}

func NewTermMACFlowMPLS(link *nlamsg.Link) *fibcapi.TerminationMacFlow {
	m := fibcapi.NewTermMACMatch(
		NewPortId(link),
		fibcapi.ETHTYPE_MPLS,
		link.Attrs().HardwareAddr.String(),
		link.VlanId(),
	)
	a := []*fibcapi.TerminationMacFlow_Action{}
	return fibcapi.NewTermMACFlow(m, a, uint32(fibcapi.FlowMod_MPLS1))
}

func (r *RIBController) SendTermMACFlow(cmd fibcapi.FlowMod_Cmd, link *nlamsg.Link) error {

	flows := []*fibcapi.TerminationMacFlow{
		NewTermMACFlowIPv4(link),
		NewTermMACFlowIPv6(link),
		NewTermMACFlowMPLS(link),
	}
	for _, f := range flows {
		if err := r.fib.Send(f.ToMod(cmd, r.reId), 0); err != nil {
			return err
		}
	}

	return nil
}

//
// MPLS Flow (POP single label for VRF)
//
func NewMPLSFlowVRF(label uint32, nid uint8) *fibcapi.MPLSFlow {
	m := fibcapi.NewMPLSMatch(label, true)
	a := []*fibcapi.MPLSFlow_Action{
		fibcapi.NewMPLSAction("SET_VRF", uint32(nid)),
	}
	t := uint32(fibcapi.FlowMod_MPLS_L3_TYPE)
	return fibcapi.NewMPLSFlow(m, a, t, fibcapi.GroupMod_UNSPEC, 0)
}

func (r *RIBController) SendMPLSFlowVRF(cmd fibcapi.FlowMod_Cmd, nid uint8) error {
	label := NewVRFLabel(r.label, nid)
	f := NewMPLSFlowVRF(label, nid)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}

//
// MPLS Flow (POP single label)
//
func NewMPLSFlowPop1(route *nlamsg.Route) *fibcapi.MPLSFlow {
	return NewMPLSFlowVRF(uint32(*route.MPLSDst), route.NId)
}

func (r *RIBController) SendMPLSFlowPop1(cmd fibcapi.FlowMod_Cmd, route *nlamsg.Route) error {
	f := NewMPLSFlowPop1(route)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}

//
// MPLS Flow (POP double label)
//
func NewMPLSFlowPop2(neigh *nlamsg.Neigh, route *nlamsg.Route) *fibcapi.MPLSFlow {
	m := fibcapi.NewMPLSMatch(uint32(*route.MPLSDst), false)
	a := []*fibcapi.MPLSFlow_Action{
		fibcapi.NewMPLSAction("POP_LABEL", fibcapi.ETHTYPE_MPLS),
	}
	t := uint32(fibcapi.FlowMod_MPLS_TYPE)
	return fibcapi.NewMPLSFlow(m, a, t, fibcapi.GroupMod_MPLS_INTERFACE, NewNeighId(neigh))
}

func (r *RIBController) SendMPLSFlowPop2(cmd fibcapi.FlowMod_Cmd, route *nlamsg.Route) error {
	neigh, err := r.nla.GetNeigh_FlowMod(cmd, route.NId, route.GetGw())
	if err != nil {
		return err
	}

	f := NewMPLSFlowPop2(neigh, route)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}

//
// MPLS Flow (SWAP)
//
func NewMPLSFlowSwap(route *nlamsg.Route, bos bool) *fibcapi.MPLSFlow {
	m := fibcapi.NewMPLSMatch(uint32(*route.MPLSDst), bos)
	a := []*fibcapi.MPLSFlow_Action{}
	t := uint32(fibcapi.FlowMod_MPLS_TYPE)
	return fibcapi.NewMPLSFlow(m, a, t, fibcapi.GroupMod_MPLS_SWAP, uint32(*route.MPLSDst))
}

func (r *RIBController) SendMPLSFlowSwap(cmd fibcapi.FlowMod_Cmd, route *nlamsg.Route, bos bool) error {
	f := NewMPLSFlowSwap(route, bos)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}

//
// Unicast Routing (for Neighbor)
//
func NewUnicastRoutingFlowNeigh(neigh *nlamsg.Neigh) *fibcapi.UnicastRoutingFlow {
	m := fibcapi.NewUnicastRoutingMatchNeigh(neigh.IP, neigh.NId)
	return fibcapi.NewUnicastRoutingFlow(m, nil, fibcapi.GroupMod_L3_UNICAST, NewNeighId(neigh))
}

func (r *RIBController) SendUnicastRoutingFlowNeigh(cmd fibcapi.FlowMod_Cmd, neigh *nlamsg.Neigh) error {
	if !checkIP(neigh.IP) {
		return nil
	}

	f := NewUnicastRoutingFlowNeigh(neigh)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}

//
// Unicast Routing (for Route)
//
func NewUnicastRoutingFlow(neigh *nlamsg.Neigh, route *nlamsg.Route) *fibcapi.UnicastRoutingFlow {
	m := fibcapi.NewUnicastRoutingMatchRoute(route.GetDst(), route.NId)
	return fibcapi.NewUnicastRoutingFlow(m, nil, fibcapi.GroupMod_L3_UNICAST, NewNeighId(neigh))
}

func (r *RIBController) SendUnicastRoutingFlow(cmd fibcapi.FlowMod_Cmd, route *nlamsg.Route) error {
	if !checkIPNet(route.GetDst()) {
		return nil
	}

	neigh, err := r.nla.GetNeigh_FlowMod(cmd, route.NId, route.GetGw())
	if err != nil {
		return err
	}

	f := NewUnicastRoutingFlow(neigh, route)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}

//
// Unicast Routing (for MPLS)
//
func NewUnicastRoutingFlowMPLS(route *nlamsg.Route) *fibcapi.UnicastRoutingFlow {
	enId := route.EnIds[len(route.EnIds)-1]
	m := fibcapi.NewUnicastRoutingMatchRoute(route.GetDst(), route.NId)
	return fibcapi.NewUnicastRoutingFlow(m, nil, fibcapi.GroupMod_MPLS_L3_VPN, enId)
}

func (r *RIBController) SendUnicastRoutingFlowMPLS(cmd fibcapi.FlowMod_Cmd, route *nlamsg.Route) error {
	f := NewUnicastRoutingFlowMPLS(route)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}

//
// PolicyACL (match ip_dst and send controller)
//
func NewACLFlowByAddr(addr *nlamsg.Addr) *fibcapi.PolicyACLFlow {
	return fibcapi.NewPolicyACLFlowByAddr(addr.Family, addr.IPNet.IP, addr.NId)
}

func (r *RIBController) SendACLFlowByAddr(cmd fibcapi.FlowMod_Cmd, addr *nlamsg.Addr) error {
	if !checkIP(addr.IP) {
		return nil
	}

	f := NewACLFlowByAddr(addr)
	return r.fib.Send(f.ToMod(cmd, r.reId), 0)
}
