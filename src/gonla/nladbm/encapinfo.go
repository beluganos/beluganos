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
// Encap Key
//
type EncapInfoKey struct {
	Dst string // ip/mask (net.IPNet.String())
	Vrf uint32 // mpls:0, vrf:label
}

func NewEncapInfoKey(dst *net.IPNet, vrf uint32) *EncapInfoKey {
	return &EncapInfoKey{
		Dst: dst.String(),
		Vrf: vrf,
	}
}

func EncapInfoToKey(e *nlamsg.EncapInfo) *EncapInfoKey {
	return NewEncapInfoKey(e.Dst, e.Vrf)
}

//
// Table Interface
//
type EncapInfoTable interface {
	EncapId(*net.IPNet, uint32) uint32
	Insert(*nlamsg.EncapInfo) *nlamsg.EncapInfo
	Select(*EncapInfoKey) *nlamsg.EncapInfo
	Delete(*EncapInfoKey) *nlamsg.EncapInfo
	Walk(func(*nlamsg.EncapInfo) error) error
	WalkFree(func(*nlamsg.EncapInfo) error) error
}

func NewEncapInfoTable() EncapInfoTable {
	return newEncapInfoTable()
}

//
// Table
//
type encapTnfoTable struct {
	mutex   sync.RWMutex
	infos   map[EncapInfoKey]*nlamsg.EncapInfo
	counter *nlalib.Counter32
}

func newEncapInfoTable() *encapTnfoTable {
	return &encapTnfoTable{
		infos:   make(map[EncapInfoKey]*nlamsg.EncapInfo),
		counter: &nlalib.Counter32{},
	}
}

func (t *encapTnfoTable) find(key *EncapInfoKey) *nlamsg.EncapInfo {
	if e, ok := t.infos[*key]; ok {
		return e
	}
	return nil
}

func (t *encapTnfoTable) EncapId(dst *net.IPNet, vrf uint32) uint32 {
	ei := nlamsg.NewEncapInfo(dst, vrf, 0)
	t.Insert(ei)
	return ei.EnId
}

func (t *encapTnfoTable) Insert(e *nlamsg.EncapInfo) (old *nlamsg.EncapInfo) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	key := EncapInfoToKey(e)
	if old = t.find(key); old != nil {
		e.EnId = old.EnId
	} else {
		e.EnId = t.counter.Next()
	}

	t.infos[*key] = e.Copy()

	return
}

func (t *encapTnfoTable) Select(key *EncapInfoKey) *nlamsg.EncapInfo {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.find(key)
}

func (t *encapTnfoTable) Delete(key *EncapInfoKey) (old *nlamsg.EncapInfo) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if old := t.find(key); old != nil {
		delete(t.infos, *key)
	}

	return
}

func (t *encapTnfoTable) Walk(f func(*nlamsg.EncapInfo) error) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *encapTnfoTable) WalkFree(f func(*nlamsg.EncapInfo) error) error {
	for _, e := range t.infos {
		if err := f(e); err != nil {
			return err
		}
	}
	return nil
}
