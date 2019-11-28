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

package ribsdbm

import (
	"fabricflow/ribs/api/ribsapi"
	"sync"
)

const (
	// RTany is RT for all.
	RTany = "*"
	// RTmic is RT for mic.
	RTmic = ""
)

//
// RicEntry is ric monitoring messages.
//
type RicEntry struct {
	NId    uint8
	Rt     string
	Stream ribsapi.RIBSCoreApi_MonitorRibServer
}

//
// Key returns key for RibTable.
//
func (e *RicEntry) Key() string {
	return e.Rt
}

//
// RicTable is ric table.
//
type RicTable struct {
	Mutex   sync.RWMutex
	Entries map[string]*RicEntry // key: RT
}

//
// NewRicTable returns new RicTable
//
func NewRicTable() *RicTable {
	return &RicTable{
		Entries: make(map[string]*RicEntry),
	}
}

//
// Add registers ric entry.
//
func (t *RicTable) Add(c *RicEntry) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	t.Entries[c.Rt] = c
}

//
// Delete removes ric entry.
//
func (t *RicTable) Delete(rt string) *RicEntry {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if entry, ok := t.Entries[rt]; ok {
		delete(t.Entries, rt)
		return entry
	}

	return nil
}

//
// Select retuens ric entry.
//
func (t *RicTable) Select(rt string, f func(*RicEntry)) bool {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	if e, ok := t.Entries[rt]; ok {
		f(e)
		return true
	}
	return false
}

//
// Update updates ric entry.
//
func (t *RicTable) Update(rt string, f func(*RicEntry)) bool {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if e, ok := t.Entries[rt]; ok {
		f(e)
		return true
	}
	return false
}

//
// Range returns all ric entry.
//
func (t *RicTable) Range(f func(string, *RicEntry) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	for key, entry := range t.Entries {
		if err := f(key, entry); err != nil {
			return err
		}
	}
	return nil
}
