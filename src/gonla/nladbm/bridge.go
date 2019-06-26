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
	"fmt"
	"gonla/nlalib"
	"gonla/nlamsg"
	"sync"
)

type BridgeVlanInfoKey struct {
	NId   uint8
	Index int // netlink.LinkAttrs.Index
	Vid   uint16
}

func NewBridgeVlanInfoKey(nid uint8, index int, vid uint16) *BridgeVlanInfoKey {
	return &BridgeVlanInfoKey{
		NId:   nid,
		Index: index,
		Vid:   vid,
	}
}

func BridgeVlanInfoToKey(br *nlamsg.BridgeVlanInfo) *BridgeVlanInfoKey {
	return NewBridgeVlanInfoKey(br.NId, br.Index, br.Vid)
}

func (k *BridgeVlanInfoKey) String() string {
	return fmt.Sprintf("nid:%d index:%d vid:%d", k.NId, k.Index, k.Vid)
}

type BridgeVlanInfoTable interface {
	Count() uint32
	Insert(*nlamsg.BridgeVlanInfo) *nlamsg.BridgeVlanInfo
	Select(*BridgeVlanInfoKey) *nlamsg.BridgeVlanInfo
	ListByIndex(int, func(*nlamsg.BridgeVlanInfo) bool) bool
	Update(*BridgeVlanInfoKey, func(*nlamsg.BridgeVlanInfo) error) error
	Delete(*BridgeVlanInfoKey) *nlamsg.BridgeVlanInfo
	Walk(func(*nlamsg.BridgeVlanInfo) error) error
	WalkFree(func(*nlamsg.BridgeVlanInfo) error) error
	GCInit()
	GCList(int) []BridgeVlanInfoKey
}

func NewBridgeVlanInfoTable() BridgeVlanInfoTable {
	return newBridgeVlanInfoTable()
}

type bridgeVlanInfoTable struct {
	Mutex    sync.RWMutex
	bridges  map[BridgeVlanInfoKey]*nlamsg.BridgeVlanInfo
	gc       *BridgeVlanInfoGC
	indexIdx *BridgeVlanInfoIfIndexIndex
	counter  *nlalib.Counters32
}

func newBridgeVlanInfoTable() *bridgeVlanInfoTable {
	return &bridgeVlanInfoTable{
		bridges:  map[BridgeVlanInfoKey]*nlamsg.BridgeVlanInfo{},
		gc:       NewBridgeVlanInfoGC(),
		indexIdx: NewBridgeVlanInfoIfIndexIndex(),
		counter:  nlalib.NewCounters32(),
	}
}

func (t *bridgeVlanInfoTable) find(key *BridgeVlanInfoKey) *nlamsg.BridgeVlanInfo {
	br, _ := t.bridges[*key]
	return br
}

func (t *bridgeVlanInfoTable) Count() uint32 {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return uint32(len(t.bridges))
}

func (t *bridgeVlanInfoTable) GCInit() {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	t.gc.Init()
	for k, _ := range t.bridges {
		t.gc.Reset(k)
	}
}

func (t *bridgeVlanInfoTable) GCList(index int) []BridgeVlanInfoKey {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if index < 0 {
		return t.gc.List()
	}

	if e, ok := t.indexIdx.Select(index); ok {
		return t.gc.ListBy(e.Keys())
	}

	return []BridgeVlanInfoKey{}
}

func (t *bridgeVlanInfoTable) Insert(br *nlamsg.BridgeVlanInfo) (old *nlamsg.BridgeVlanInfo) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := BridgeVlanInfoToKey(br)
	if old = t.find(key); old == nil {
		br.BrId = t.counter.Next(br.NId)
	} else {
		br.BrId = old.BrId
	}

	t.bridges[*key] = br.Copy()
	t.gc.Using(*key)

	if old != nil {
		t.indexIdx.Delete(br)
	}
	t.indexIdx.Insert(br)

	return
}

func (t *bridgeVlanInfoTable) Select(key *BridgeVlanInfoKey) *nlamsg.BridgeVlanInfo {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *bridgeVlanInfoTable) ListByIndex(index int, f func(*nlamsg.BridgeVlanInfo) bool) bool {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	e, ok := t.indexIdx.Select(index)
	if !ok {
		return false
	}

	for _, key := range e.Keys() {
		if brvlan := t.find(&key); brvlan != nil {
			if next := f(brvlan); !next {
				return true
			}
		}
	}

	return true
}

func (t *bridgeVlanInfoTable) Update(key *BridgeVlanInfoKey, f func(*nlamsg.BridgeVlanInfo) error) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if br := t.find(key); br != nil {
		return f(br)
	}

	return fmt.Errorf("BridgeVlan not found. %s", key)
}

func (t *bridgeVlanInfoTable) Delete(key *BridgeVlanInfoKey) (old *nlamsg.BridgeVlanInfo) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.bridges, *key)
		t.gc.Delete(*key)
		t.indexIdx.Delete(old)
	}

	return
}

func (t *bridgeVlanInfoTable) Walk(f func(*nlamsg.BridgeVlanInfo) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *bridgeVlanInfoTable) WalkFree(f func(*nlamsg.BridgeVlanInfo) error) error {
	for _, bridge := range t.bridges {
		if err := f(bridge); err != nil {
			return err
		}
	}
	return nil
}

//
// Index Entry (Ifindex)
//
type BridgeVlanInfoIfIndexIndexEntry struct {
	keys map[BridgeVlanInfoKey]*BridgeVlanInfoKey
}

func NewBridgeVlanInfoIfIndexIndexEntry() *BridgeVlanInfoIfIndexIndexEntry {
	return &BridgeVlanInfoIfIndexIndexEntry{
		keys: map[BridgeVlanInfoKey]*BridgeVlanInfoKey{},
	}
}

func (e *BridgeVlanInfoIfIndexIndexEntry) Insert(key *BridgeVlanInfoKey) {
	e.keys[*key] = key
}

func (e *BridgeVlanInfoIfIndexIndexEntry) Delete(key *BridgeVlanInfoKey) {
	delete(e.keys, *key)
}

func (e *BridgeVlanInfoIfIndexIndexEntry) Len() int {
	return len(e.keys)
}

func (e *BridgeVlanInfoIfIndexIndexEntry) Keys() []BridgeVlanInfoKey {
	keys := []BridgeVlanInfoKey{}
	for key, _ := range e.keys {
		keys = append(keys, key)
	}

	return keys
}

//
// Index Table (Ifindex)
//
type BridgeVlanInfoIfIndexIndex struct {
	entries map[int]*BridgeVlanInfoIfIndexIndexEntry // key: Ifindex
}

func NewBridgeVlanInfoIfIndexIndex() *BridgeVlanInfoIfIndexIndex {
	return &BridgeVlanInfoIfIndexIndex{
		entries: map[int]*BridgeVlanInfoIfIndexIndexEntry{},
	}
}

func (i *BridgeVlanInfoIfIndexIndex) Insert(b *nlamsg.BridgeVlanInfo) {
	k := b.Index
	e, ok := i.entries[k]
	if !ok {
		e = NewBridgeVlanInfoIfIndexIndexEntry()
		i.entries[k] = e
	}

	e.Insert(BridgeVlanInfoToKey(b))
}

func (i *BridgeVlanInfoIfIndexIndex) Delete(b *nlamsg.BridgeVlanInfo) {
	k := b.Index
	e, ok := i.entries[k]
	if !ok {
		return
	}

	e.Delete(BridgeVlanInfoToKey(b))

	if e.Len() == 0 {
		delete(i.entries, k)
	}
}

func (i *BridgeVlanInfoIfIndexIndex) Select(index int) (e *BridgeVlanInfoIfIndexIndexEntry, ok bool) {
	e, ok = i.entries[index]
	return
}

//
// GC Table
//
type BridgeVlanInfoGC struct {
	gc map[BridgeVlanInfoKey]bool
}

func NewBridgeVlanInfoGC() *BridgeVlanInfoGC {
	return &BridgeVlanInfoGC{
		gc: map[BridgeVlanInfoKey]bool{},
	}
}

func (g *BridgeVlanInfoGC) Init() {
	g.gc = map[BridgeVlanInfoKey]bool{}
}

func (g *BridgeVlanInfoGC) List() []BridgeVlanInfoKey {
	keys := []BridgeVlanInfoKey{}
	for key, leaked := range g.gc {
		if leaked {
			keys = append(keys, key)
		}
	}

	return keys
}

func (g *BridgeVlanInfoGC) ListBy(keys []BridgeVlanInfoKey) []BridgeVlanInfoKey {
	ks := []BridgeVlanInfoKey{}
	for _, key := range keys {
		if leaked, ok := g.gc[key]; ok && leaked {
			ks = append(ks, key)
		}
	}

	return ks
}

func (g *BridgeVlanInfoGC) Reset(key BridgeVlanInfoKey) {
	g.gc[key] = true
}

func (g *BridgeVlanInfoGC) Using(key BridgeVlanInfoKey) {
	g.gc[key] = false
}

func (g *BridgeVlanInfoGC) Delete(key BridgeVlanInfoKey) {
	delete(g.gc, key)
}
