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
type NeighKey struct {
	// do not use NeId
	NId  uint8
	Addr string // ip (net.IP.String())
}

func NewNeighKey(nid uint8, ip net.IP) *NeighKey {
	return &NeighKey{
		NId:  nid,
		Addr: ip.String(),
	}
}

func NeighToKey(n *nlamsg.Neigh) *NeighKey {
	return NewNeighKey(n.NId, n.IP)
}

//
// Table interface
//
type NeighTable interface {
	Insert(*nlamsg.Neigh) *nlamsg.Neigh
	Select(*NeighKey) *nlamsg.Neigh
	Delete(*NeighKey) *nlamsg.Neigh
	Walk(f func(*nlamsg.Neigh) error) error
	WalkFree(f func(*nlamsg.Neigh) error) error
}

func NewNeighTable() NeighTable {
	return &neighTable{
		Neighs:  make(map[NeighKey]*nlamsg.Neigh),
		Counter: nlalib.NewCounters16(),
	}
}

//
// Table
//
type neighTable struct {
	Mutex   sync.RWMutex
	Neighs  map[NeighKey]*nlamsg.Neigh
	Counter *nlalib.Counters16
}

func (t *neighTable) find(key *NeighKey) *nlamsg.Neigh {
	n, _ := t.Neighs[*key]
	return n
}

func (t *neighTable) Insert(n *nlamsg.Neigh) (old *nlamsg.Neigh) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := NeighToKey(n)
	if old = t.find(key); old == nil {
		n.NeId = t.Counter.Next(n.NId)
		t.Neighs[*key] = n.Copy()
	} else {
		n.NeId = old.NeId
	}

	return
}

func (t *neighTable) Select(key *NeighKey) *nlamsg.Neigh {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *neighTable) Walk(f func(*nlamsg.Neigh) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *neighTable) WalkFree(f func(*nlamsg.Neigh) error) error {
	for _, n := range t.Neighs {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (t *neighTable) Delete(key *NeighKey) (old *nlamsg.Neigh) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.Neighs, *key)
	}

	return
}
