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
type AddrKey struct {
	// note: do not use AdId field.
	NId  uint8
	Addr string // ip/mask (net.IPNet.String())
}

func NewAddrKey(nid uint8, addr *net.IPNet) *AddrKey {
	return &AddrKey{
		NId:  nid,
		Addr: addr.String(),
	}
}

func AddrToKey(a *nlamsg.Addr) *AddrKey {
	return NewAddrKey(a.NId, a.IPNet)
}

//
// Table interface
//
type AddrTable interface {
	Insert(*nlamsg.Addr) *nlamsg.Addr
	Select(*AddrKey) *nlamsg.Addr
	Delete(*AddrKey) *nlamsg.Addr
	Walk(f func(*nlamsg.Addr) error) error
	WalkFree(f func(*nlamsg.Addr) error) error
}

func NewAddrTable() AddrTable {
	return &addrTable{
		Addrs:   make(map[AddrKey]*nlamsg.Addr),
		Counter: nlalib.NewCounters32(),
	}
}

//
// Table
//
type addrTable struct {
	Mutex   sync.RWMutex
	Addrs   map[AddrKey]*nlamsg.Addr
	Counter *nlalib.Counters32
}

func (t *addrTable) find(key *AddrKey) *nlamsg.Addr {
	n, _ := t.Addrs[*key]
	return n
}

func (t *addrTable) Insert(a *nlamsg.Addr) (old *nlamsg.Addr) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := AddrToKey(a)
	if old = t.find(key); old == nil {
		a.AdId = t.Counter.Next(a.NId)
		t.Addrs[*key] = a.Copy()
	} else {
		a.AdId = old.AdId
	}

	return
}

func (t *addrTable) Select(key *AddrKey) *nlamsg.Addr {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *addrTable) Walk(f func(*nlamsg.Addr) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *addrTable) WalkFree(f func(*nlamsg.Addr) error) error {
	for _, n := range t.Addrs {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (t *addrTable) Delete(key *AddrKey) (old *nlamsg.Addr) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.Addrs, *key)
	}

	return
}
