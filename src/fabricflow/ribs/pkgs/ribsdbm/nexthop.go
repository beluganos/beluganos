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
	"fmt"
	"net"
	"sync"
)

func isRTmic(rt string) bool {
	return rt == RTmic
}

//
// NewNexthopKey returns ket for NexthopTable.
//
func NewNexthopKey(ip net.IP, rt string) string {
	if isRTmic(rt) {
		return fmt.Sprintf("%s", ip)
	}

	return fmt.Sprintf("%s@%s", ip, rt)
}

//
// Nexthop is nexthop entry.
//
type Nexthop struct {
	Rt       string
	Addr     net.IP
	SourceID net.IP
}

//
// IsMic detects rt is mic or not.
//
func (e *Nexthop) IsMic() bool {
	return isRTmic(e.Rt)
}

//
// NewNexthop returns new Nexthop
//
func NewNexthop(addr net.IP, rt string, sourceID net.IP) *Nexthop {
	return &Nexthop{
		Rt:       rt,
		Addr:     addr,
		SourceID: sourceID,
	}
}

//
// Key returns key for Nexthop table.
//
func (e *Nexthop) Key() string {
	return NewNexthopKey(e.Addr, e.Rt)
}

//
// NexthopTable is used to detect infinit loop.
//
type NexthopTable struct {
	Mutex   sync.RWMutex
	Entries map[string]*Nexthop // Key: <RT>_<IP> or <IP>
}

//
// NewNexthopTable returns new NexthopTable.
//
func NewNexthopTable() *NexthopTable {
	return &NexthopTable{
		Entries: make(map[string]*Nexthop),
	}
}

func (t *NexthopTable) find(key string) *Nexthop {
	if e, ok := t.Entries[key]; ok {
		return e
	}
	return nil
}

//
// Add register nexthop entry.
//
func (t *NexthopTable) Add(ip net.IP, rt string, srcID net.IP) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	ent := NewNexthop(ip, rt, srcID)
	t.Entries[ent.Key()] = ent
}

//
// Select returns nexthop entry.
//
func (t *NexthopTable) Select(ip net.IP, rt string, f func(*Nexthop)) bool {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	if e := t.find(NewNexthopKey(ip, rt)); e != nil {
		f(e)
		return true
	}

	return false
}

//
// Update updates nexthop entry.
//
func (t *NexthopTable) Update(ip net.IP, rt string, f func(*Nexthop)) bool {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	if e := t.find(NewNexthopKey(ip, rt)); e != nil {
		f(e)
		return true
	}

	return false
}

//
// Delete removes nexthop entry
//
func (t *NexthopTable) Delete(ip net.IP, rt string) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	delete(t.Entries, NewNexthopKey(ip, rt))
}

func (t *NexthopTable) DeleteByRT(rt string, f func(*Nexthop) bool) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	delKeys := []string{}

FOR_LOOP:
	for _, entry := range t.Entries {
		if entry.Rt != rt {
			continue FOR_LOOP
		}

		if ok := f(entry); ok {
			delKeys = append(delKeys, entry.Key())
		}
	}

	for _, key := range delKeys {
		delete(t.Entries, key)
	}

}

//
// Range returns all nexthop entry.
//
func (t *NexthopTable) Range(f func(string, *Nexthop)) {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	for key, entry := range t.Entries {
		f(key, entry)
	}
}
