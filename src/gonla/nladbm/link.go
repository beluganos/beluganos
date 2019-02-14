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
type LinkKey struct {
	// note: do not use LnId field.
	NId   uint8
	Index int // netlink.LinkAttrs.Index
}

func NewLinkKey(nid uint8, index int) *LinkKey {
	return &LinkKey{
		NId:   nid,
		Index: index,
	}
}

func LinkToKey(ln *nlamsg.Link) *LinkKey {
	return NewLinkKey(ln.NId, ln.Attrs().Index)
}

//
// Table interface
//
type LinkTable interface {
	Insert(*nlamsg.Link) *nlamsg.Link
	Select(*LinkKey) *nlamsg.Link
	Delete(*LinkKey) *nlamsg.Link
	Walk(f func(*nlamsg.Link) error) error
	WalkFree(f func(*nlamsg.Link) error) error
	SelectTun(*IptunKey) *nlamsg.Iptun
	UpdateTun(*IptunKey, func(*nlamsg.Iptun) error) error
	WalkTun(func(*nlamsg.Iptun) error) error

	WalkTunByRemote(uint8, *net.IPNet, func(*nlamsg.Iptun) error) error
}

func NewLinkTable() LinkTable {
	return &linkTable{
		Links:   make(map[LinkKey]*nlamsg.Link),
		Counter: nlalib.NewCounters16(),
		TunIdx:  NewIptunTable(),
	}
}

//
// Table
//
type linkTable struct {
	Mutex   sync.RWMutex
	Links   map[LinkKey]*nlamsg.Link
	Counter *nlalib.Counters16
	TunIdx  IptunTable
}

func (t *linkTable) find(key *LinkKey) *nlamsg.Link {
	n, _ := t.Links[*key]
	return n
}

func (t *linkTable) Insert(ln *nlamsg.Link) (old *nlamsg.Link) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := LinkToKey(ln)
	if old = t.find(key); old == nil {
		ln.LnId = t.Counter.Next(ln.NId)
	} else {
		ln.LnId = old.LnId
	}

	t.Links[*key] = ln.Copy()

	iptun := nlamsg.NewIptun(ln)
	if old != nil {
		t.TunIdx.Delete(IptunToKey(iptun))
	}
	t.TunIdx.Insert(iptun)

	return
}

func (t *linkTable) Select(key *LinkKey) *nlamsg.Link {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *linkTable) Walk(f func(*nlamsg.Link) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *linkTable) WalkFree(f func(*nlamsg.Link) error) error {
	for _, n := range t.Links {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (t *linkTable) Delete(key *LinkKey) (old *nlamsg.Link) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.Links, *key)
		t.TunIdx.Delete(IptunToKey(nlamsg.NewIptun(old)))
	}

	return
}

func (t *linkTable) SelectTun(key *IptunKey) (e *nlamsg.Iptun) {
	return t.TunIdx.Select(key)
}

func (t *linkTable) UpdateTun(key *IptunKey, f func(*nlamsg.Iptun) error) error {
	return t.TunIdx.Update(key, f)
}

func (t *linkTable) WalkTun(f func(*nlamsg.Iptun) error) error {
	return t.TunIdx.Walk(f)
}

func (t *linkTable) WalkTunByRemote(nid uint8, route *net.IPNet, f func(*nlamsg.Iptun) error) error {
	return t.TunIdx.Walk(func(iptun *nlamsg.Iptun) error {
		if (iptun.NId == nid) && route.Contains(iptun.Remote()) {
			if err := f(iptun); err != nil {
				return err
			}
		}
		return nil
	})
}
