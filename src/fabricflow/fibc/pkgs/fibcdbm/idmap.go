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
	"fmt"
	"sync"
)

//
// IDEntry is entry of IDMap
//
type IDEntry struct {
	DpID uint64
	ReID string

	version uint64
}

//
// NewIDEntry returns new IDEntry.
//
func NewIDEntry(dpID uint64, reID string) *IDEntry {
	return &IDEntry{
		DpID: dpID,
		ReID: reID,
	}
}

//
// String is stringer.
//
func (e *IDEntry) String() string {
	return fmt.Sprintf("dpid:%d, reid:'%s'", e.DpID, e.ReID)
}

//
// IDMap is map of dp_id and re_id.
//
type IDMap struct {
	mutex   sync.RWMutex
	entries map[uint64]*IDEntry
	reIDKey *IDMapReIDKey

	version uint64
}

//
// NewIDMap returns new IDMap.
//
func NewIDMap() *IDMap {
	return &IDMap{
		entries: map[uint64]*IDEntry{},
		reIDKey: NewIDMapReIDKey(),
	}
}

func (m *IDMap) find(dpID uint64) *IDEntry {
	if e, ok := m.entries[dpID]; ok {
		return e
	}

	return nil
}

func (m *IDMap) findByReID(reID string) *IDEntry {
	if dpID, ok := m.reIDKey.Select(reID); ok {
		return m.find(dpID)
	}

	return nil
}

//
// VerUp updates version.
//
func (m *IDMap) VerUp() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	m.version++
}

//
// SelectByDpID select entry by dpId.
//
func (m *IDMap) SelectByDpID(dpID uint64) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if e := m.find(dpID); e != nil {
		return e.ReID, true
	}

	return "", false
}

//
// SelectByReID select entry by reId.
//
func (m *IDMap) SelectByReID(reID string) (uint64, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if e := m.findByReID(reID); e != nil {
		return e.DpID, true
	}

	return 0, false
}

//
// Register add entry.
//
func (m *IDMap) Register(e *IDEntry) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	dpID := e.DpID
	if old := m.find(dpID); old != nil {
		if old.ReID != e.ReID {
			return fmt.Errorf("Invalid id. %d %s", e.DpID, e.ReID)
		}

		e.version = m.version
		return nil
	}

	e.version = m.version
	m.entries[dpID] = e
	m.reIDKey.Add(e.ReID, e.DpID)

	return nil
}

//
// SelectOrRegister returns IDEntry or add if not exist.
//
func (m *IDMap) SelectOrRegister(e *IDEntry, f func(*IDEntry) bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	dpID := e.DpID
	old := m.find(dpID)
	if old != nil {
		if old.ReID != e.ReID {
			return fmt.Errorf("Invalid id. %d %s", e.DpID, e.ReID)
		}
	} else {
		old = e
		m.entries[dpID] = e
		m.reIDKey.Add(e.ReID, e.DpID)
	}

	if ok := f(old); ok {
		old.version = m.version
	}

	return nil
}

//
// GC free unused entries.
//
func (m *IDMap) GC(f func(*IDEntry) bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delEntries := []*IDEntry{}
	for _, e := range m.entries {
		if e.version != m.version {
			if ok := f(e); ok {
				e.version = m.version
			} else {
				delEntries = append(delEntries, e)
			}
		}
	}

	for _, e := range delEntries {
		delete(m.entries, e.DpID)
		m.reIDKey.Delete(e.ReID)
	}
}

//
// UnregisterByDpID remove entry by dpId.
//
func (m *IDMap) UnregisterByDpID(dpID uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if e := m.find(dpID); e != nil {
		delete(m.entries, dpID)
		m.reIDKey.Delete(e.ReID)
	}
}

//
// UnregisterByReID remove entry by reId,
//
func (m *IDMap) UnregisterByReID(reID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if e := m.findByReID(reID); e != nil {
		delete(m.entries, e.DpID)
		m.reIDKey.Delete(reID)
	}
}

//
// Range enumerate all entry.
//
func (m *IDMap) Range(f func(*IDEntry)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, e := range m.entries {
		f(e)
	}
}

//
// IDMapReIDKey is reId key table.
//
type IDMapReIDKey struct {
	entries map[string]uint64
}

//
// NewIDMapReIDKey returns new IDMapReIdKey.
//
func NewIDMapReIDKey() *IDMapReIDKey {
	return &IDMapReIDKey{
		entries: map[string]uint64{},
	}
}

//
// Select returns dpId and exist/not-exist.
//
func (k *IDMapReIDKey) Select(reID string) (uint64, bool) {
	dpID, ok := k.entries[reID]
	return dpID, ok
}

//
// Add add entry
//
func (k *IDMapReIDKey) Add(reID string, dpID uint64) {
	if _, ok := k.Select(reID); !ok {
		k.entries[reID] = dpID
	}
}

//
// Delete remove entry.
//
func (k *IDMapReIDKey) Delete(reID string) {
	if _, ok := k.Select(reID); ok {
		delete(k.entries, reID)
	}
}
