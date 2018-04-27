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
	"gonla/nlamsg"
	"net"
	"sync"
)

//
// Key
//
type MplsKey struct {
	NId    uint8
	LLabel uint32
}

func NewMplsKey(nid uint8, llabel uint32) *MplsKey {
	return &MplsKey{
		NId:    nid,
		LLabel: llabel,
	}
}

func MplsToKey(r *nlamsg.Route) *MplsKey {
	if r.Route.MPLSDst == nil {
		return nil
	}
	return NewMplsKey(r.NId, uint32(*r.Route.MPLSDst))
}

//
// Table interface
//
type MplsTable interface {
	Insert(*nlamsg.Route) *nlamsg.Route
	Select(*MplsKey) *nlamsg.Route
	Delete(*MplsKey) *nlamsg.Route
	Walk(f func(*nlamsg.Route) error) error
	WalkFree(f func(*nlamsg.Route) error) error
	WalkByGw(uint8, net.IP, func(*nlamsg.Route) error) error
	WalkByGwFree(uint8, net.IP, func(*nlamsg.Route) error) error
}

func NewMplsTable() MplsTable {
	return &mplsTable{
		Mplss: make(map[MplsKey]*nlamsg.Route),
		GwIdx: NewMplsGwIndex(),
	}
}

//
// Table
//
type mplsTable struct {
	Mutex sync.RWMutex
	Mplss map[MplsKey]*nlamsg.Route
	GwIdx *MplsGwIndex
}

func (t *mplsTable) find(key *MplsKey) *nlamsg.Route {
	n, _ := t.Mplss[*key]
	return n
}

func (t *mplsTable) Insert(r *nlamsg.Route) (old *nlamsg.Route) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := MplsToKey(r)
	if old = t.find(key); old == nil {
		t.Mplss[*key] = r.Copy()
		t.GwIdx.Insert(r)
	}

	return
}

func (t *mplsTable) Select(key *MplsKey) *nlamsg.Route {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *mplsTable) Walk(f func(*nlamsg.Route) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *mplsTable) WalkFree(f func(*nlamsg.Route) error) error {
	for _, n := range t.Mplss {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (t *mplsTable) Delete(key *MplsKey) (old *nlamsg.Route) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.Mplss, *key)
		t.GwIdx.Delete(old)
	}

	return
}

func (t *mplsTable) WalkByGw(nid uint8, ip net.IP, f func(*nlamsg.Route) error) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	return t.WalkByGwFree(nid, ip, f)
}

func (t *mplsTable) WalkByGwFree(nid uint8, ip net.IP, f func(*nlamsg.Route) error) error {
	if e, ok := t.GwIdx.Select(nid, ip); ok {
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

//
// GW Index entry
//
type MplsGwIndexEntry struct {
	Keys map[MplsKey]*MplsKey
}

func NewMplsGwIndexEntry() *MplsGwIndexEntry {
	return &MplsGwIndexEntry{
		Keys: make(map[MplsKey]*MplsKey),
	}
}

func (m *MplsGwIndexEntry) Insert(key *MplsKey) {
	m.Keys[*key] = key
}

func (m *MplsGwIndexEntry) Delete(key *MplsKey) {
	delete(m.Keys, *key)
}

func (m *MplsGwIndexEntry) Len() int {
	return len(m.Keys)
}

//
// GW Index Table
//
type MplsGwIndex struct {
	Entry map[NeighKey]*MplsGwIndexEntry
}

func NewMplsGwIndex() *MplsGwIndex {
	return &MplsGwIndex{
		Entry: make(map[NeighKey]*MplsGwIndexEntry),
	}
}

func (m *MplsGwIndex) Insert(route *nlamsg.Route) {
	gw := route.GetGw()
	if gw == nil {
		return
	}

	key := NewNeighKey(route.NId, gw)
	e, ok := m.Entry[*key]
	if !ok {
		e = NewMplsGwIndexEntry()
		m.Entry[*key] = e
	}

	e.Insert(MplsToKey(route))
}

func (m *MplsGwIndex) Delete(route *nlamsg.Route) {
	gw := route.GetGw()
	if gw == nil {
		return
	}

	key := NewNeighKey(route.NId, gw)
	e, ok := m.Entry[*key]
	if !ok {
		return
	}

	e.Delete(MplsToKey(route))

	if e.Len() == 0 {
		delete(m.Entry, *key)
	}
}

func (m *MplsGwIndex) Select(nid uint8, ip net.IP) (e *MplsGwIndexEntry, ok bool) {
	e, ok = m.Entry[*NewNeighKey(nid, ip)]
	return
}
