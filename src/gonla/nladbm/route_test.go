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

package nladbm

import (
	"net"
	"testing"

	"gonla/nlamsg"

	"github.com/vishvananda/netlink"
)

func TestTouteTableIptunRemote(t *testing.T) {
	nid := uint8(0)
	cfgRemotes := []string{
		"2001:2001::/32",
		"2001:2001::/64",
		"2001:2002::/32",
		"2001:2002:1::/64",
	}

	var (
		route *nlamsg.Route
	)
	tbl := NewRouteTable().(*routeTable)

	for _, r := range cfgRemotes {
		_, remote, _ := net.ParseCIDR(r)
		tbl.RegisterTunRemote(nid, remote)
	}

	routes := []string{
		"2001:2001::/32",
		"2001:2001::/64",
		"2001:2001::254/128",
		"2001:2002::/32",
		"2001:2002:1::/64",
	}

	for index, r := range routes {
		_, dst, _ := net.ParseCIDR(r)
		route := &netlink.Route{
			LinkIndex: index + 1,
			Dst:       dst,
		}
		rtid := uint32(index + 11)
		tbl.Insert(nlamsg.NewRoute(route, nid, rtid, nil, []uint32{}))
	}

	route = tbl.SelectByTunRemote(nid, net.ParseIP("2001:2001::253"))
	if route != nil {
		t.Errorf("routeTable.SelectByTunRemote unmatch. route=%s", route)
	}

	route = tbl.SelectByTunRemote(nid, net.ParseIP("2001:2001::254"))
	if route == nil {
		t.Errorf("routeTable.SelectByTunRemote unmatch. route=%s", route)
	}

	route = tbl.SelectByTunRemote(nid, net.ParseIP("2001:2002::254"))
	if route != nil {
		t.Errorf("routeTable.SelectByTunRemote unmatch. route=%s", route)
	}

	route = tbl.SelectByTunRemote(nid, net.ParseIP("2001:2002:1::254"))
	if route != nil {
		t.Errorf("routeTable.SelectByTunRemote unmatch. route=%s", route)
	}
}
