// -*- coding; utf-8 -*-

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

package nladbm

import (
	"net"

	"gonla/nlamsg"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestLinkTable_Insert(t *testing.T) {
	tbl := NewLinkTable().(*linkTable)
	ln1 := &netlink.Device{}
	ln1.Attrs().Index = 1
	link0 := nlamsg.NewLink(ln1, 0, 0)
	link1 := nlamsg.NewLink(ln1, 1, 0)

	if old := tbl.Insert(link0); old != nil {
		t.Errorf("linkTable Insert link0 error. %v", old)
	}

	if old := tbl.Insert(link1); old != nil {
		t.Errorf("linkTable Insert link1 error. %v", old)
	}

	if l := len(tbl.Links); l != 2 {
		t.Errorf("linkTable Insert size unmatch. %d", l)
	}
}

func TestLinkIptun(t *testing.T) {
	nid := uint8(0)
	tbl := NewLinkTable().(*linkTable)

	remotes := []string{
		"2001:2001::1",
		"2001:2001::2",
		"2001:2001::3",
		"2001:2002::1",
	}

	check := map[string]struct{}{}

	for index, remote := range remotes {
		check[remote] = struct{}{}
		iptun := &netlink.Iptun{}
		iptun.Attrs().Index = index + 1
		iptun.Remote = net.ParseIP(remote)
		link := nlamsg.NewLink(iptun, nid, uint16(index+11))
		tbl.Insert(link)
	}

	_, route, _ := net.ParseCIDR("2001:2001::/32")
	tbl.WalkTunByRemote(nid, route, func(iptun *nlamsg.Iptun) error {
		remote := iptun.Remote().String()
		if _, ok := check[remote]; !ok {
			t.Errorf("WalkTunByRemote unmatch. iptun=%s", iptun)
		}
		delete(check, remote)
		return nil
	})
	if v := len(check); v != 1 {
		t.Errorf("WalkTunByRemote unmatch. check=%v", check)
	}

	_, route, _ = net.ParseCIDR("2001:2002::/32")
	tbl.WalkTunByRemote(nid, route, func(iptun *nlamsg.Iptun) error {
		remote := iptun.Remote().String()
		if _, ok := check[remote]; !ok {
			t.Errorf("WalkTunByRemote unmatch. iptun=%s", iptun)
		}
		delete(check, remote)
		return nil
	})
	if v := len(check); v != 0 {
		t.Errorf("WalkTunByRemote unmatch. check=%v", check)
	}
}
