// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/osrg/gobgp/pkg/packet/bgp"
	"github.com/vishvananda/netlink"
)

//
// TunnelRoute
//
type TunnelRoute struct {
	Prefix     *net.IPNet
	Nexthop    net.IP
	Family     uint16         // bgp.AFI_IP ot bgp.API_IP6
	TunnelType bgp.TunnelType // bgp.TUNNEL_TPPE_XXX
}

func (r *TunnelRoute) String() string {
	return fmt.Sprintf("%s nexthop %s %d %s", r.Prefix, r.Nexthop, r.Family, r.TunnelType)
}

func (r *TunnelRoute) WriteTo(w io.Writer) (sum int64, err error) {
	var n int
	n, err = fmt.Fprintf(w, "%s\n", r)
	sum += int64(n)
	return
}

//
// TunnelEntry is entry of TunnelTable.
//
type TunnelEntry struct {
	Id     uint32
	Type   bgp.TunnelType // ipip, ip6tnl
	remote net.IP
	local  net.IP
	attrs  *netlink.LinkAttrs
	Routes map[string]*TunnelRoute // key: prefix
}

func NewTunnelEntry(link netlink.Link, id uint32) *TunnelEntry {
	tunType, ok := func() (bgp.TunnelType, bool) {
		switch link.Type() {
		case "ipip":
			return bgp.TUNNEL_TYPE_IP_IN_IP, true

		case "ip6tnl":
			return bgp.TUNNEL_TYPE_IPV6, true

		default:
			return 0, false
		}
	}()
	if !ok {
		return nil
	}

	tun := link.(*netlink.Iptun)
	return &TunnelEntry{
		Id:     id,
		Type:   tunType,
		remote: tun.Remote,
		local:  tun.Local,
		attrs:  tun.Attrs(),
		Routes: map[string]*TunnelRoute{},
	}
}

func (c *TunnelEntry) Ifname() string {
	return c.attrs.Name
}

func (c *TunnelEntry) Ifindex() int {
	return c.attrs.Index
}

func (c *TunnelEntry) Remote() string {
	return c.remote.String()
}

func (c *TunnelEntry) Local() string {
	return c.local.String()
}

func (e *TunnelEntry) String() string {
	return fmt.Sprintf("%s %d %s %s %s %d", e.Ifname(), e.Id, e.Type, e.remote, e.local, len(e.Routes))
}

func (e *TunnelEntry) AddRoute(route *TunnelRoute) {
	e.Routes[route.Prefix.String()] = route
}

func (e *TunnelEntry) DelRoute(route *TunnelRoute) int {
	delete(e.Routes, route.Prefix.String())
	return len(e.Routes)
}

func (c *TunnelEntry) WriteTo(w io.Writer) (sum int64, err error) {
	var n int

	n, err = fmt.Fprintf(w, "%s\n", c)
	sum += int64(n)
	if err != nil {
		return
	}

	for _, route := range c.Routes {
		var nn int64
		nn, err = route.WriteTo(w)
		sum += nn
		if err != nil {
			return
		}
	}

	return
}

//
// TunnelTable is tunnel device table.
//
type TunnelTable struct {
	mutex    sync.Mutex
	factory  *TunnelFactory
	remotes  map[string]*TunnelEntry // key: tunnel remote
	ifnames  map[string]*TunnelEntry // key: ifname
	prefixes map[string]*TunnelEntry // key: prefix
}

func NewTunnelTable(device string) *TunnelTable {
	t := &TunnelTable{factory: NewTunnelFactory(device)}
	t.Reset()
	return t
}

func (t *TunnelTable) Reset() {
	t.remotes = map[string]*TunnelEntry{}
	t.ifnames = map[string]*TunnelEntry{}
	t.prefixes = map[string]*TunnelEntry{}
}

func (t *TunnelTable) findByRemote(remote string) (e *TunnelEntry, ok bool) {
	e, ok = t.remotes[remote]
	return
}

func (t *TunnelTable) FindByRemote(remote string) (*TunnelEntry, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.findByRemote(remote)
}

func (t *TunnelTable) findByIfname(ifname string) (e *TunnelEntry, ok bool) {
	e, ok = t.ifnames[ifname]
	return
}

func (t *TunnelTable) FindByIfname(ifname string) (*TunnelEntry, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.findByIfname(ifname)
}

func (t *TunnelTable) findByPrefix(prefix string) (e *TunnelEntry, ok bool) {
	e, ok = t.prefixes[prefix]
	return
}

func (t *TunnelTable) FindByPrefix(prefix string) (*TunnelEntry, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.findByPrefix(prefix)
}

func (t *TunnelTable) put(tun *TunnelEntry) error {
	if _, ok := t.findByRemote(tun.Remote()); ok {
		return fmt.Errorf("%s already exists.", tun.Remote())
	}
	if _, ok := t.findByIfname(tun.Ifname()); ok {
		return fmt.Errorf("%s already exists.", tun.Ifname())
	}

	t.remotes[tun.Remote()] = tun
	t.ifnames[tun.Ifname()] = tun

	return nil
}

func (t *TunnelTable) Put(e *TunnelEntry) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.put(e)
}

func (t *TunnelTable) pop(remote string) (tun *TunnelEntry, ok bool) {
	if tun, ok = t.findByRemote(remote); ok {
		delete(t.remotes, remote)
		delete(t.ifnames, tun.Ifname())
	}

	return
}

func (t *TunnelTable) Pop(remote string) (tun *TunnelEntry, ok bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.pop(remote)
}

func (t *TunnelTable) addRoute(route *TunnelRoute, tun *TunnelEntry) {
	prefix := route.Prefix.String()
	if _, ok := t.findByPrefix(prefix); !ok {
		t.prefixes[prefix] = tun
		tun.AddRoute(route)
	}
}

func (t *TunnelTable) AddRoute(route *TunnelRoute, tun *TunnelEntry) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.addRoute(route, tun)
}

func (t *TunnelTable) delRoute(route *TunnelRoute, tun *TunnelEntry) int {
	delete(t.prefixes, route.Prefix.String())
	return tun.DelRoute(route)
}

func (t *TunnelTable) DelRoute(route *TunnelRoute, tun *TunnelEntry) int {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.delRoute(route, tun)
}

func (t *TunnelTable) Load() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.factory.LinkList(func(link netlink.Link, id uint32) {
		if tun := NewTunnelEntry(link, id); tun != nil {
			t.put(tun)
			t.factory.RouteList(link, func(afi uint16, route *netlink.Route) {
				rt := &TunnelRoute{
					Prefix:     route.Dst,
					Nexthop:    route.Gw,
					Family:     afi,
					TunnelType: tun.Type,
				}
				t.addRoute(rt, tun)
			})
		}
	})
}

func (t *TunnelTable) NewIfName() (string, uint32) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for {
		counter := t.factory.NextCounter()
		ifname := t.factory.NewIfName(counter)
		if _, ok := t.findByIfname(ifname); !ok {
			return ifname, counter
		}
	}
}

func (t *TunnelTable) RangeByRemote(f func(string, *TunnelEntry)) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for remote, e := range t.remotes {
		f(remote, e)
	}
}

func (t *TunnelTable) RangeByIfname(f func(string, *TunnelEntry)) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for ifname, e := range t.ifnames {
		f(ifname, e)
	}
}

func (t *TunnelTable) RangeByPrefix(f func(string, *TunnelEntry)) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for prefix, e := range t.prefixes {
		f(prefix, e)
	}
}

func (t *TunnelTable) Range(f func(string, string, *TunnelEntry)) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for remote, e := range t.remotes {
		f("remote", remote, e)
	}

	for ifname, e := range t.ifnames {
		f("ifname", ifname, e)
	}

	for prefix, e := range t.prefixes {
		f("prefix", prefix, e)
	}
}

func (t *TunnelTable) WriteTo(w io.Writer) (sum int64, err error) {
	var n int
	var nn int64

	t.Range(func(name string, key string, tun *TunnelEntry) {
		n, err = fmt.Fprintf(w, "%s: %s\n", name, key)
		sum += int64(n)
		if err != nil {
			return
		}
		nn, err = tun.WriteTo(w)
		sum += nn
		if err != nil {
			return
		}
	})

	return
}
