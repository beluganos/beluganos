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
)

func DispatchLink(nlmsg *NetlinkMessage, link *Link, app interface{}) {
	if h, ok := app.(NetlinkLinkHandler); ok {
		h.NetlinkLink(nlmsg, link)
	}
}

func DispatchToLink(nlmsg *NetlinkMessage, app interface{}) error {
	if h, ok := app.(NetlinkLinkHandler); ok {
		link, err := LinkDeserialize(nlmsg)
		if err != nil {
			return err
		}

		h.NetlinkLink(nlmsg, link)
	}
	return nil
}

func DispatchAddr(nlmsg *NetlinkMessage, addr *Addr, app interface{}) {
	if h, ok := app.(NetlinkAddrHandler); ok {
		h.NetlinkAddr(nlmsg, addr)
	}
}

func DispatchToAddr(nlmsg *NetlinkMessage, app interface{}) error {
	if h, ok := app.(NetlinkAddrHandler); ok {
		addr, err := AddrDeserialize(nlmsg)
		if err != nil {
			return err
		}

		h.NetlinkAddr(nlmsg, addr)
	}
	return nil
}

func DispatchNeigh(nlmsg *NetlinkMessage, neigh *Neigh, app interface{}) {
	if h, ok := app.(NetlinkNeighHandler); ok {
		h.NetlinkNeigh(nlmsg, neigh)
	}
}

func DispatchToNeigh(nlmsg *NetlinkMessage, app interface{}) error {
	if h, ok := app.(NetlinkNeighHandler); ok {
		neigh, err := NeighDeserialize(nlmsg)
		if err != nil {
			return err
		}

		h.NetlinkNeigh(nlmsg, neigh)
	}
	return nil
}

func DispatchRoute(nlmsg *NetlinkMessage, route *Route, app interface{}) {
	if h, ok := app.(NetlinkRouteHandler); ok {
		h.NetlinkRoute(nlmsg, route)
	}
}

func DispatchToRoute(nlmsg *NetlinkMessage, app interface{}) error {
	if h, ok := app.(NetlinkRouteHandler); ok {
		route, err := RouteDeserialize(nlmsg)
		if err != nil {
			return err
		}

		h.NetlinkRoute(nlmsg, route)
	}
	return nil
}

func DispatchNode(nlmsg *NetlinkMessage, node *Node, app interface{}) {
	if h, ok := app.(NetlinkNodeHandler); ok {
		h.NetlinkNode(nlmsg, node)
	}
}

func DispatchToNode(nlmsg *NetlinkMessage, app interface{}) error {
	if h, ok := app.(NetlinkNodeHandler); ok {
		node, err := NodeDeserialize(nlmsg)
		if err != nil {
			return err
		}

		h.NetlinkNode(nlmsg, node)
	}
	return nil
}

func DispatchVpn(nlmsg *NetlinkMessage, vpn *Vpn, app interface{}) {
	if h, ok := app.(NetlinkVpnHandler); ok {
		h.NetlinkVpn(nlmsg, vpn)
	}
}

func DispatchToVpn(nlmsg *NetlinkMessage, app interface{}) error {
	if h, ok := app.(NetlinkVpnHandler); ok {
		vpn, err := VpnDeserialize(nlmsg)
		if err != nil {
			return err
		}

		h.NetlinkVpn(nlmsg, vpn)
	}
	return nil
}

func DispatchNetlinkMessage(nlmsg *NetlinkMessage, app interface{}) {
	if h, ok := app.(NetlinkMessageHandler); ok {
		h.NetlinkMessage(nlmsg)
	}
}

func Dispatch(nlmsg *NetlinkMessage, app interface{}) error {

	DispatchNetlinkMessage(nlmsg, app)

	switch nlmsg.Group() {
	case nlalink.RTMGRP_LINK:
		return DispatchToLink(nlmsg, app)

	case nlalink.RTMGRP_ADDR:
		return DispatchToAddr(nlmsg, app)

	case nlalink.RTMGRP_NEIGH:
		return DispatchToNeigh(nlmsg, app)

	case nlalink.RTMGRP_ROUTE:
		return DispatchToRoute(nlmsg, app)

	case nlalink.RTMGRP_NODE:
		return DispatchToNode(nlmsg, app)

	case nlalink.RTMGRP_VPN:
		return DispatchToVpn(nlmsg, app)

	default:
		return fmt.Errorf("Dispatcher: unsupported nlmsg. %d", nlmsg.Type())
	}
}

func DispatchUnion(nlmsg *NetlinkMessageUnion, app interface{}) error {

	m := &NetlinkMessage{}
	m.Header = nlmsg.Header
	m.Data = []byte{}
	m.NId = nlmsg.NId
	m.Src = nlmsg.Src

	switch nlmsg.Group() {
	case nlalink.RTMGRP_LINK:
		DispatchLink(m, nlmsg.GetLink(), app)

	case nlalink.RTMGRP_ADDR:
		DispatchAddr(m, nlmsg.GetAddr(), app)

	case nlalink.RTMGRP_NEIGH:
		DispatchNeigh(m, nlmsg.GetNeigh(), app)

	case nlalink.RTMGRP_ROUTE:
		DispatchRoute(m, nlmsg.GetRoute(), app)

	case nlalink.RTMGRP_NODE:
		DispatchNode(m, nlmsg.GetNode(), app)

	case nlalink.RTMGRP_VPN:
		DispatchVpn(m, nlmsg.GetVpn(), app)

	default:
		return fmt.Errorf("Dispatcher: unsupported nlmsg. %d", nlmsg.Type())
	}

	return nil
}
