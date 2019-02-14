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

const (
	IFMAP_NAME_OIDMAP = "oidmap"
	IFMAP_NAME_SHIFT  = "shift"
)

type IfMapEntry struct {
	Name string
	Min  uint
	Max  uint
}

func NewIfMapOidMap(min, max uint) *IfMapEntry {
	return &IfMapEntry{
		Name: IFMAP_NAME_OIDMAP,
		Min:  min,
		Max:  max,
	}
}

func NewIfMapShift(min, max uint) *IfMapEntry {
	return &IfMapEntry{
		Name: IFMAP_NAME_SHIFT,
		Min:  min,
		Max:  max,
	}
}

func (e *IfMapEntry) String() string {
	return fmt.Sprintf("%s min:%d, max:%v", e.Name, e.Min, e.Max)
}

func (e *IfMapEntry) Key() string {
	return e.Name
}

type IfMapTable struct {
	mutex   sync.Mutex
	entries map[string]*IfMapEntry
}

func NewIfMapTable() *IfMapTable {
	return &IfMapTable{
		entries: map[string]*IfMapEntry{},
	}
}

func (t *IfMapTable) WriteTo(w io.Writer) (sum int64, err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, v := range t.entries {
		var n int
		n, err = fmt.Fprintf(w, "IfMap %s\n", v)
		sum += int64(n)
		if err != nil {
			return
		}
	}
	return
}

func (t *IfMapTable) add(e *IfMapEntry) {
	t.entries[e.Key()] = e
}

func (t *IfMapTable) Add(e *IfMapEntry) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.add(e)
}

func (t *IfMapTable) get(name string) (e *IfMapEntry, ok bool) {
	e, ok = t.entries[name]
	return
}

func (t *IfMapTable) Get(name string) (e *IfMapEntry, ok bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.get(name)
}

func (t *IfMapTable) match(name string, v uint) (e *IfMapEntry, ok bool) {
	if e, ok = t.entries[name]; !ok {
		return
	}

	if v < e.Min || v > e.Max {
		return e, false
	}

	return
}

func (t *IfMapTable) Match(name string, v uint) (*IfMapEntry, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.match(name, v)
}
