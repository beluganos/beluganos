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
	"github.com/vishvananda/netlink"
	"net"
)

//
// netlink.Route
//
func CopyRoute(src *netlink.Route) *netlink.Route {
	dst := *src
	return &dst
}

//
// Route
//
type Route struct {
	*netlink.Route
	RtId  uint32 // auto increment
	NId   uint8
	VpnGw net.IP
	EnIds []uint32
}

func (r *Route) Copy() *Route {
	return &Route{
		Route: CopyRoute(r.Route),
		RtId:  r.RtId,
		NId:   r.NId,
		VpnGw: r.VpnGw,
		EnIds: r.EnIds,
	}
}

func (r *Route) String() string {
	return fmt.Sprintf("%s RtId: %d NId: %d VpnGw: %s EnId: %v", r.Route, r.RtId, r.NId, r.VpnGw, r.EnIds)
}

func (r *Route) GetDst() *net.IPNet {
	if r.Dst == nil {
		return &net.IPNet{}
	}
	return r.Dst
}

func (r *Route) MultiPathIndex() int {
	return len(r.MultiPath) - 1
}

func (r *Route) GetLinkIndex() int {
	if i := r.MultiPathIndex(); i >= 0 {
		return r.MultiPath[i].LinkIndex
	}
	return r.LinkIndex
}

func (r *Route) GetGw() net.IP {
	if i := r.MultiPathIndex(); i >= 0 {
		return r.MultiPath[i].Gw
	}
	return r.Gw
}

func (r *Route) GetEncap() netlink.Encap {
	if i := r.MultiPathIndex(); i >= 0 {
		return r.MultiPath[i].Encap
	}
	return r.Encap
}

func (r *Route) GetMPLSEncap() *netlink.MPLSEncap {
	if en := r.GetEncap(); en != nil {
		if men, ok := en.(*netlink.MPLSEncap); ok {
			return men
		}
	}

	return nil
}

func (r *Route) GetMPLSNewDst() *netlink.MPLSDestination {
	d := func() netlink.Destination {
		if i := r.MultiPathIndex(); i >= 0 {
			return r.MultiPath[i].NewDst
		}
		return r.NewDst
	}()
	if dst, ok := d.(*netlink.MPLSDestination); ok {
		return dst
	}
	return nil
}

func NewRoute(route *netlink.Route, nid uint8, id uint32, vpnGw net.IP, enIds []uint32) *Route {
	return &Route{
		RtId:  id,
		Route: route,
		NId:   nid,
		VpnGw: vpnGw,
		EnIds: enIds,
	}
}

func RouteDeserialize(nlmsg *NetlinkMessage) (*Route, error) {
	route, err := netlink.RouteDeserialize(nlmsg.Data)
	if err != nil {
		return nil, err
	}

	return NewRoute(route, nlmsg.NId, 0, nil, []uint32{}), nil
}
