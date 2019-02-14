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

package nladbm

import (
	"gonla/nlalib"
	"gonla/nlamsg"
	"net"
	"sync"
)

//
// Key
//
type RouteKey struct {
	// note: do not use AdId field.
	NId  uint8
	Addr string // ip/mask (net.IPNet.String())
}

func NewRouteKey(nid uint8, addr *net.IPNet) *RouteKey {
	return &RouteKey{
		NId:  nid,
		Addr: addr.String(),
	}
}

func RouteToKey(r *nlamsg.Route) *RouteKey {
	return NewRouteKey(r.NId, r.GetDst())
}

//
// Table interface
//
type RouteTable interface {
	Insert(*nlamsg.Route) *nlamsg.Route
	Select(*RouteKey) *nlamsg.Route
	Delete(*RouteKey) *nlamsg.Route
	Walk(f func(*nlamsg.Route) error) error
	WalkFree(f func(*nlamsg.Route) error) error
	WalkByGw(uint8, net.IP, func(*nlamsg.Route) error) error
	WalkByGwFree(uint8, net.IP, func(*nlamsg.Route) error) error
	WalkByLink(uint8, int, func(*nlamsg.Route) error) error
	WalkByLinkFree(uint8, int, func(*nlamsg.Route) error) error

	RegisterTunRemote(uint8, *net.IPNet)
	SelectByTunRemote(uint8, net.IP) *nlamsg.Route
}

func NewRouteTable() RouteTable {
	return &routeTable{
		Routes:  make(map[RouteKey]*nlamsg.Route),
		Counter: nlalib.NewCounters32(),
		GwIdx:   NewRouteGwIndex(),
		LinkIdx: NewRouteLinkIndex(),

		TunRemote: NewIptunPeerTable(),
	}
}

//
// Table
//
type routeTable struct {
	Mutex   sync.RWMutex
	Routes  map[RouteKey]*nlamsg.Route
	Counter *nlalib.Counters32
	GwIdx   *RouteGwIndex
	LinkIdx *RouteLinkIndex

	TunRemote *IptunPeerTable
}

func (t *routeTable) find(key *RouteKey) *nlamsg.Route {
	n, _ := t.Routes[*key]
	return n
}

func (t *routeTable) Insert(r *nlamsg.Route) (old *nlamsg.Route) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := RouteToKey(r)
	if old = t.find(key); old == nil {
		r.RtId = t.Counter.Next(r.NId)
	} else {
		r.RtId = old.RtId
	}

	t.Routes[*key] = r.Copy()

	if old != nil {
		t.GwIdx.Delete(old)
		t.LinkIdx.Delete(old)
	}
	t.GwIdx.Insert(r)
	t.LinkIdx.Insert(r)

	return
}

func (t *routeTable) Select(key *RouteKey) *nlamsg.Route {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *routeTable) Walk(f func(*nlamsg.Route) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *routeTable) WalkFree(f func(*nlamsg.Route) error) error {
	for _, n := range t.Routes {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (t *routeTable) Delete(key *RouteKey) (old *nlamsg.Route) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.Routes, *key)
		t.GwIdx.Delete(old)
		t.LinkIdx.Delete(old)
	}

	return
}

func (t *routeTable) WalkByGw(nid uint8, ip net.IP, f func(*nlamsg.Route) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkByGwFree(nid, ip, f)
}

func (t *routeTable) WalkByGwFree(nid uint8, ip net.IP, f func(*nlamsg.Route) error) error {
	e, ok := t.GwIdx.Select(nid, ip)
	if ok {
		for _, key := range e.Keys {
			if route := t.find(key); route != nil {
				if err := f(route); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *routeTable) WalkByLink(nid uint8, ifindex int, f func(*nlamsg.Route) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkByLinkFree(nid, ifindex, f)
}

func (t *routeTable) WalkByLinkFree(nid uint8, ifindex int, f func(*nlamsg.Route) error) error {
	e, ok := t.LinkIdx.Select(nid, ifindex)
	if ok {
		for _, key := range e.Keys {
			if route := t.find(key); route != nil {
				if err := f(route); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *routeTable) RegisterTunRemote(nid uint8, remote *net.IPNet) {
	peer := NewIptunPeer(nid, remote)
	t.TunRemote.Insert(peer)
}

func (t *routeTable) SelectByTunRemote(nid uint8, remote net.IP) (route *nlamsg.Route) {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	// first: select by remote/128(ipv6) or remote/32(ipv4)
	dst := nlalib.NewIPNetFromIP(remote)
	if route = t.find(NewRouteKey(nid, dst)); route != nil {
		return
	}

	// second: select by network contains remote.
	if tun := t.TunRemote.SelectByIP(nid, remote); tun != nil {
		route = t.find(NewRouteKey(nid, tun.Dst))
	}

	return
}

//
// GW Index Entry
//
type RouteGwIndexEntry struct {
	Keys map[RouteKey]*RouteKey
}

func NewRouteGwIndexEntry() *RouteGwIndexEntry {
	return &RouteGwIndexEntry{
		Keys: make(map[RouteKey]*RouteKey),
	}
}

func (r *RouteGwIndexEntry) Insert(key *RouteKey) {
	r.Keys[*key] = key
}

func (r *RouteGwIndexEntry) Delete(key *RouteKey) {
	delete(r.Keys, *key)
}

func (r *RouteGwIndexEntry) Len() int {
	return len(r.Keys)
}

//
// Link Index Table
//
type RouteLinkIndex struct {
	Entry map[LinkKey]*RouteGwIndexEntry
}

func NewRouteLinkIndex() *RouteLinkIndex {
	return &RouteLinkIndex{
		Entry: map[LinkKey]*RouteGwIndexEntry{},
	}
}

func (r *RouteLinkIndex) Insert(route *nlamsg.Route) {
	ifindex := route.LinkIndex
	if ifindex <= 0 {
		return
	}

	key := NewLinkKey(route.NId, ifindex)
	e, ok := r.Entry[*key]
	if !ok {
		e = NewRouteGwIndexEntry()
		r.Entry[*key] = e
	}

	e.Insert(RouteToKey(route))
}

func (r *RouteLinkIndex) Delete(route *nlamsg.Route) {
	ifindex := route.LinkIndex
	if ifindex <= 0 {
		return
	}

	key := NewLinkKey(route.NId, ifindex)
	e, ok := r.Entry[*key]
	if !ok {
		return
	}

	e.Delete(RouteToKey(route))

	if e.Len() == 0 {
		delete(r.Entry, *key)
	}
}

func (r *RouteLinkIndex) Select(nid uint8, ifindex int) (e *RouteGwIndexEntry, ok bool) {
	e, ok = r.Entry[*NewLinkKey(nid, ifindex)]
	return
}

//
// GW Index Table
//
type RouteGwIndex struct {
	Entry map[NeighKey]*RouteGwIndexEntry
}

func NewRouteGwIndex() *RouteGwIndex {
	return &RouteGwIndex{
		Entry: make(map[NeighKey]*RouteGwIndexEntry),
	}
}

func (r *RouteGwIndex) Insert(route *nlamsg.Route) {
	gw := route.GetGw()
	if gw == nil {
		return
	}

	key := NewNeighKey(route.NId, gw)
	e, ok := r.Entry[*key]
	if !ok {
		e = NewRouteGwIndexEntry()
		r.Entry[*key] = e
	}

	e.Insert(RouteToKey(route))
}

func (r *RouteGwIndex) Delete(route *nlamsg.Route) {
	gw := route.GetGw()
	if gw == nil {
		return
	}

	key := NewNeighKey(route.NId, gw)
	e, ok := r.Entry[*key]
	if !ok {
		return
	}

	e.Delete(RouteToKey(route))

	if e.Len() == 0 {
		delete(r.Entry, *key)
	}
}

func (r *RouteGwIndex) Select(nid uint8, ip net.IP) (e *RouteGwIndexEntry, ok bool) {
	e, ok = r.Entry[*NewNeighKey(nid, ip)]
	return
}
