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

package fibcdbm

import (
	"sync"
)

//
// DPEntry is interface for entry of DPSet.
//
type DPEntry interface {
	EntryID() string
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	Remote() string
>>>>>>> develop
=======
	Remote() string
>>>>>>> develop
=======
	Remote() string
>>>>>>> develop
	Start(<-chan struct{})
	Stop()
}

//
// DPSet is set of DPEntry.
//
type DPSet struct {
	mutex   sync.RWMutex
	entries map[string]DPEntry
}

//
// NewDPSet returns new DPSet
//
func NewDPSet() *DPSet {
	return &DPSet{
		entries: map[string]DPEntry{},
	}
}

func (t *DPSet) find(id string) DPEntry {
	if e, ok := t.entries[id]; ok {
		return e
	}

	return nil
}

//
// Add adds entry if not exist.
//
func (t *DPSet) Add(e DPEntry) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	eid := e.EntryID()
	if old := t.find(eid); old != nil {
		return false
	}

	t.entries[eid] = e

	return true
}

//
// Delete removes entry of id key.
//
func (t *DPSet) Delete(id string) DPEntry {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if e := t.find(id); e != nil {
		delete(t.entries, id)
		return e
	}

	return nil
}

//
// Select selects entry and call func.
//
func (t *DPSet) Select(id string, f func(DPEntry)) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if e := t.find(id); e != nil {
		f(e)
		return true
	}

	return false
}

//
// Range call f for each entries.
//
func (t *DPSet) Range(f func(DPEntry)) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, e := range t.entries {
		f(e)
	}
}
