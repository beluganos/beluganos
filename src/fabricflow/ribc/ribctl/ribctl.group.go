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

//
// L2 Interface Group
//
func NewL2InterfaceGroup(link *nlamsg.Link) *fibcapi.L2InterfaceGroup {
	hwaddr := func() net.HardwareAddr {
		if link.Iptun() != nil {
			return fibcapi.HardwareAddrDummy
		}
		return link.Attrs().HardwareAddr
	}()
	return fibcapi.NewL2InterfaceGroup(
		NewPortId(link),
		link.VlanId(),
		false, // vlanTranslation
		hwaddr,
		link.Attrs().MTU,
		link.NId,
	)
}

func (r *RIBController) SendL2InterfaceGroup(cmd fibcapi.GroupMod_Cmd, link *nlamsg.Link) error {
	g := NewL2InterfaceGroup(link)
	return r.fib.Send(g.ToMod(cmd, r.reId), 0)
}

//
// L3 Unicast Group
//
func NewL3UnicastGroup(link, phyLink *nlamsg.Link, neigh *nlamsg.Neigh) *fibcapi.L3UnicastGroup {
	vid := func() uint16 {
		if link == nil {
			return 0
		}
		return link.VlanId()
	}()
	srcMAC := func() net.HardwareAddr {
		if phyLink == nil {
			return fibcapi.HardwareAddrNone
		}
		return phyLink.Attrs().HardwareAddr
	}()

	g := fibcapi.NewL3UnicastGroup(
		NewNeighId(neigh),
		NewPortId(link),
		NewPortId(phyLink),
		vid,
		neigh.HardwareAddr,
		srcMAC,
	)

	if iptun := neigh.GetIptun(); iptun != nil {
		tunType, _ := fibcapi.ParseTunnelTypeFromNative(iptun.TunType)
		g.SetTunnel(tunType, neigh.IP, iptun.SrcIP)
	}

	return g
}

func (r *RIBController) SendL3UnicastGroup(cmd fibcapi.GroupMod_Cmd, neigh *nlamsg.Neigh) error {
	hwtype, _ := fibcapi.ParseHardwareAddrType(neigh.HardwareAddr)

	if hwtype == fibcapi.HWADDR_TYPE_NONE {
		if cmd != fibcapi.GroupMod_DELETE {
			return nil
		}
	}

	if ok := hwtype.Has(fibcapi.HWADDR_TYPE_MULTICAST); ok {
		return nil
	}

	link, err := r.nla.GetLink_GroupMod(cmd, neigh.NId, neigh.LinkIndex)
	if err != nil {
		return err
	}

	phyLink, err := func() (*nlamsg.Link, error) {
		if neigh.IsTunnelRemote() {
			return r.nla.GetLink_GroupMod(cmd, neigh.NId, neigh.PhyLink)
		}
		return link, nil
	}()
	if err != nil {
		return err
	}

	g := NewL3UnicastGroup(link, phyLink, neigh)
	return r.fib.Send(g.ToMod(cmd, r.reId), 0)
}

//
// MPLS Interface Group
//
func NewMPLSInterfaceGroup(link *nlamsg.Link, neigh *nlamsg.Neigh) *fibcapi.MPLSInterfaceGroup {
	if link == nil {
		// when DELNEIGH, only neigh-id needed,
		return fibcapi.NewMPLSInterfaceGroup(
			NewNeighId(neigh),
			NewPortId(nil),
			0,
			neigh.HardwareAddr,
			neigh.HardwareAddr, // dummy
		)
	} else {
		return fibcapi.NewMPLSInterfaceGroup(
			NewNeighId(neigh),
			NewPortId(link),
			link.VlanId(),
			neigh.HardwareAddr,
			link.Attrs().HardwareAddr,
		)
	}
}

func (r *RIBController) SendMPLSInterfaceGroup(cmd fibcapi.GroupMod_Cmd, neigh *nlamsg.Neigh) error {
	hwtype, _ := fibcapi.ParseHardwareAddrType(neigh.HardwareAddr)

	if hwtype == fibcapi.HWADDR_TYPE_NONE {
		if cmd != fibcapi.GroupMod_DELETE {
			return nil
		}
	}

	if ok := hwtype.Has(fibcapi.HWADDR_TYPE_MULTICAST); ok {
		return nil
	}

	if neigh.IsTunnelRemote() {
		// if neigh is tunnel remote peer, no group set.
		return nil
	}

	link, err := r.nla.GetLink_GroupMod(cmd, neigh.NId, neigh.LinkIndex)
	if err != nil {
		return err
	}

	g := NewMPLSInterfaceGroup(link, neigh)
	return r.fib.Send(g.ToMod(cmd, r.reId), 0)
}

//
// MPLS L3 VPN Group
//
func NewMPLSLabelGroupL3VPN(enId, label, neId, nextEnId uint32) *fibcapi.MPLSLabelGroup {
	return fibcapi.NewMPLSLabelGroup(
		fibcapi.GroupMod_MPLS_L3_VPN,
		enId,
		label,
		neId,
		nextEnId,
	)
}

//
// MPLS Tunnel1 Group
//
func NewMPLSLabelGroupTun1(enId, label, neId uint32) *fibcapi.MPLSLabelGroup {
	return fibcapi.NewMPLSLabelGroup(
		fibcapi.GroupMod_MPLS_TUNNEL1,
		enId,
		label,
		neId,
		0,
	)
}

//
// MPLS Label Groups for encap mpls single label
//
func (r *RIBController) SendMPLSLabelGroupMPLS(cmd fibcapi.GroupMod_Cmd, route *nlamsg.Route) error {
	labels := route.GetMPLSEncap().Labels
	if len(labels) != 1 {
		return nil
	}

	neigh, err := r.nla.GetNeigh_GroupMod(cmd, route.NId, route.GetGw())
	if err != nil {
		return err
	}

	neId := NewNeighId(neigh)
	g := NewMPLSLabelGroupL3VPN(route.EnIds[0], uint32(labels[0]), neId, 0)
	if err := r.fib.Send(g.ToMod(cmd, r.reId), 0); err != nil {
		return err
	}

	g = NewMPLSLabelGroupTun1(route.EnIds[0], uint32(labels[0]), neId)
	if err := r.fib.Send(g.ToMod(cmd, r.reId), 0); err != nil {
		return err
	}

	return nil
}

//
// MPLS Label Groups for encap mpls double label
//
func (r *RIBController) SendMPLSLabelGroupVPN(cmd fibcapi.GroupMod_Cmd, route *nlamsg.Route) error {
	labels := route.GetMPLSEncap().Labels
	if len(labels) < 2 {
		return nil
	}

	g := NewMPLSLabelGroupL3VPN(route.EnIds[1], uint32(labels[1]), 0, route.EnIds[0])
	if err := r.fib.Send(g.ToMod(cmd, r.reId), 0); err != nil {
		return err
	}

	return nil
}

//
// MPLS Swap Group
//
func NewMPLSLabelGroupSwap(neigh *nlamsg.Neigh, route *nlamsg.Route) *fibcapi.MPLSLabelGroup {
	return fibcapi.NewMPLSLabelGroup(
		fibcapi.GroupMod_MPLS_SWAP,
		uint32(*route.MPLSDst),
		uint32(route.GetMPLSNewDst().Labels[0]),
		NewNeighId(neigh),
		0,
	)
}

func (r *RIBController) SendMPLSLabelGroupSwap(cmd fibcapi.GroupMod_Cmd, route *nlamsg.Route) error {
	neigh, err := r.nla.GetNeigh_GroupMod(cmd, route.NId, route.GetGw())
	if err != nil {
		return err
	}

	g := NewMPLSLabelGroupSwap(neigh, route)
	return r.fib.Send(g.ToMod(cmd, r.reId), 0)
}
