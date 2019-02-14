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
)

func TestIptunPeerTable1(t *testing.T) {
	nid := uint8(0)
	var (
		route *net.IPNet
		peer  *IptunPeer
	)

	tbl := NewIptunPeerTable()

	_, route, _ = net.ParseCIDR("2001:2001::/32")
	tbl.Insert(NewIptunPeer(nid, route))

	peer = tbl.SelectByIP(nid, net.ParseIP("2001:2001:1::254"))
	if peer == nil {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}
	if v := peer.Dst.String(); v != "2001:2001::/32" {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}

	peer = tbl.SelectByIP(nid, net.ParseIP("2001:2002:1::254"))
	if peer != nil {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}
}

func TestIptunPeerTableN(t *testing.T) {
	nid := uint8(0)
	var (
		peer *IptunPeer
	)
	routes := []string{
		"2001:2001::/32",
		"2001:2001:1::/64",
		"2001:2001:1::254/128",
	}
	tbl := NewIptunPeerTable()

	for _, r := range routes {
		_, route, _ := net.ParseCIDR(r)
		tbl.Insert(NewIptunPeer(nid, route))
	}

	peer = tbl.SelectByIP(nid, net.ParseIP("2001:2001::1"))
	if peer == nil {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}
	if v := peer.Dst.String(); v != "2001:2001::/32" {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}

	peer = tbl.SelectByIP(nid, net.ParseIP("2001:2001:1::1"))
	if peer == nil {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}
	if v := peer.Dst.String(); v != "2001:2001:1::/64" {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}

	peer = tbl.SelectByIP(nid, net.ParseIP("2001:2001:1::254"))
	if peer == nil {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}
	if v := peer.Dst.String(); v != "2001:2001:1::254/128" {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}

	peer = tbl.SelectByIP(nid, net.ParseIP("2001:2002:1::254"))
	if peer != nil {
		t.Errorf("IptunPeerTable.SelectByIP unmatch. peer=%s", peer)
	}
}
