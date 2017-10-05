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

package nlalib

import (
	"github.com/vishvananda/netlink"
	"testing"
)

func TestGetNetlinkLinks(t *testing.T) {
	links, err := GetNetlinkLinks()

	if err != nil {
		t.Errorf("nlalib.GetNetlinkLinks error. %s", err)
	}

	for _, b := range links {
		link, err := netlink.LinkDeserialize(nil, b)
		if err != nil {
			t.Errorf("netlink.LinkDeserialize error. %s %v", err, b)
		}
		t.Logf("netlink.LinkDeserialize %v", link)
	}
}

func TestGetNetlinkAddrs(t *testing.T) {
	addrs, err := GetNetlinkAddrs()

	if err != nil {
		t.Errorf("nlalib.GetNetlinkAddrs error. %s", err)
	}

	for _, b := range addrs {
		addr, _, _, err := netlink.AddrDeserialize(b)
		if err != nil {
			t.Errorf("netlink.AddrDeserialize error. %s %v", err, b)
		}
		t.Logf("netlink.AddrDeserialize %v", addr)
	}
}

func TestGetNetlinkNeighs(t *testing.T) {
	neighs, err := GetNetlinkNeighs()

	if err != nil {
		t.Errorf("nlalib.GetNetlinkNeighs error. %s", err)
	}

	for _, b := range neighs {
		neigh, err := netlink.NeighDeserialize(b)
		if err != nil {
			t.Errorf("netlink.NeighDeserialize error. %s %v", err, b)
		}
		t.Logf("netlink.NeighDeserialize %v", neigh)
	}
}

func TestGetNetlinkRoutes(t *testing.T) {
	routes, err := GetNetlinkRoutes()

	if err != nil {
		t.Errorf("nlalib.GetNetlinkRoutes error. %s", err)
	}

	for _, b := range routes {
		route, err := netlink.RouteDeserialize(b)
		if err != nil {
			t.Errorf("netlink.RouteDeserialize error. %s %v", err, b)
		}
		t.Logf("netlink.RouteDeserialize %v", route)
	}
}
