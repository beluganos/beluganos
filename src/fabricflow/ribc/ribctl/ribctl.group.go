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
	fibcapi "fabricflow/fibc/api"
	"fmt"
	"gonla/nlamsg"
	"net"
)

//
// L2 Interface Group
//
func NewL2InterfaceGroup(link *nlamsg.Link, master *IfDBEntry) *fibcapi.L2InterfaceGroup {
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
		master.PortId(),
	)
}

func (r *RIBController) SendL2InterfaceGroup(cmd fibcapi.GroupMod_Cmd, link *nlamsg.Link, master *IfDBEntry) error {
	g := NewL2InterfaceGroup(link, master)
	return r.fib.Send(g.ToMod(cmd, r.reId), 0)
}

//
// L3 Unicast Group
//
func NewL3UnicastGroup(link, phyLink *IfDBEntry, neigh *nlamsg.Neigh) *fibcapi.L3UnicastGroup {
	g := fibcapi.NewL3UnicastGroup(
		NewNeighId(neigh),
		link.PortId(),
		phyLink.PortId(),
		phyLink.Vid,
		neigh.HardwareAddr,
		phyLink.HardwareAddr,
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

	var (
		link    IfDBEntry
		phyLink IfDBEntry
	)

	if ok := r.ifdb.SelectBy(&link, neigh.NId, neigh.LinkIndex); !ok {
		return fmt.Errorf("link not found. nid:%d ifindex:%d", neigh.NId, neigh.LinkIndex)
	}

	if neigh.IsTunnelRemote() && cmd != fibcapi.GroupMod_DELETE {
		if ok := r.ifdb.SelectBy(&phyLink, neigh.NId, neigh.PhyLink); !ok {
			return fmt.Errorf("phy link not found. nid:%d ifindex:%d", neigh.NId, neigh.PhyLink)
		}
	} else {
		phyLink = link
	}

	g := NewL3UnicastGroup(&link, &phyLink, neigh)
	return r.fib.Send(g.ToMod(cmd, r.reId), 0)
}

//
// MPLS Interface Group
//
func NewMPLSInterfaceGroup(link *IfDBEntry, neigh *nlamsg.Neigh) *fibcapi.MPLSInterfaceGroup {
	return fibcapi.NewMPLSInterfaceGroup(
		NewNeighId(neigh),
		link.PortId(),
		link.Vid,
		neigh.HardwareAddr,
		link.HardwareAddr,
	)
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

	var link IfDBEntry
	if ok := r.ifdb.SelectBy(&link, neigh.NId, neigh.LinkIndex); !ok {
		return fmt.Errorf("link not found. nid:%d ifindex:%d", neigh.NId, neigh.LinkIndex)
	}

	g := NewMPLSInterfaceGroup(&link, neigh)
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
