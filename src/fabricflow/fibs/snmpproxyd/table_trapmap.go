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
	"sync"
)

type TrapMapEntry struct {
	Ifname  string
	Ifindex int
	PortId  uint
}

func NewTrapMapEntry(ifname string, ifindex int, portId uint) *TrapMapEntry {
	return &TrapMapEntry{
		Ifname:  ifname,
		Ifindex: ifindex,
		PortId:  portId,
	}
}

func (e *TrapMapEntry) String() string {
	return fmt.Sprintf("name:'%s', index:%d, port:%d", e.Ifname, e.Ifindex, e.PortId)
}

func (e *TrapMapEntry) Key() string {
	return e.Ifname
}

type TrapMapTable struct {
	mutex   sync.Mutex
	entries map[string]*TrapMapEntry
}

func (t *TrapMapTable) WriteTo(w io.Writer) (sum int64, err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, val := range t.entries {
		var n int
		n, err = fmt.Fprintf(w, "TrapMap %s\n", val)
		sum += int64(n)
		if err != nil {
			return
		}
	}
	return
}

func NewTrapMapTable() *TrapMapTable {
	return &TrapMapTable{
		entries: map[string]*TrapMapEntry{},
	}
}

func (t *TrapMapTable) add(e *TrapMapEntry) bool {
	key := e.Key()
	if _, ok := t.find(key); ok {
		return false
	}

	t.entries[key] = e
	return true
}

func (t *TrapMapTable) Add(e *TrapMapEntry) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.add(e)
}

func (t *TrapMapTable) del(key string) (e *TrapMapEntry, ok bool) {
	if e, ok = t.find(key); ok {
		delete(t.entries, key)
	}
	return
}

func (t *TrapMapTable) find(key string) (e *TrapMapEntry, ok bool) {
	e, ok = t.entries[key]
	return
}

func (t *TrapMapTable) findByIfindex(ifindex int) (*TrapMapEntry, bool) {
	for _, e := range t.entries {
		if ifindex == e.Ifindex {
			return e, true
		}
	}
	return nil, false
}

func (t *TrapMapTable) FindByIfindex(ifindex int) (e *TrapMapEntry, ok bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	e, ok = t.findByIfindex(ifindex)
	return
}

func (t *TrapMapTable) UpdateByIfindex(ifindex int, f func(*TrapMapEntry)) (ok bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var e *TrapMapEntry
	if e, ok = t.findByIfindex(ifindex); ok {
		f(e)
	}
	return
}

func (t *TrapMapTable) UpdateByIfname(ifname string, f func(*TrapMapEntry)) (ok bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var e *TrapMapEntry
	if e, ok = t.find(ifname); ok {
		f(e)
	}
	return
}
