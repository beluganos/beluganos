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
}

func NewRouteTable() RouteTable {
	return &routeTable{
		Routes:  make(map[RouteKey]*nlamsg.Route),
		Counter: nlalib.NewCounters32(),
	}
}

//
// Table
//
type routeTable struct {
	Mutex   sync.RWMutex
	Routes  map[RouteKey]*nlamsg.Route
	Counter *nlalib.Counters32
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
	}

	return
}
