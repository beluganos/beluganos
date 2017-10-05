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

package ribsyn

import (
	"fabricflow/ribs/ribsmsg"
	"fmt"
	"net"
	"sync"
)

//
// Ric Table (used by RicController)
//
type RicTable struct {
	Mutex   sync.RWMutex
	Entries map[string]*RicEntry // key: RT
}

func NewRicTable() *RicTable {
	return &RicTable{
		Entries: make(map[string]*RicEntry),
	}
}

func (t *RicTable) Add(c *RicEntry) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	t.Entries[c.Rt] = c
}

func (t *RicTable) Delete(rt string) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if entry, ok := t.Entries[rt]; ok {
		delete(t.Entries, rt)
		entry.Close()
	}
}

func (t *RicTable) FindByAddr(addr string) *RicEntry {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	for _, c := range t.Entries {
		if c.Addr == addr {
			return c
		}
	}
	return nil
}

func (t *RicTable) FindByRt(rt string) *RicEntry {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	if c, ok := t.Entries[rt]; ok {
		return c
	}
	return nil
}

func (t *RicTable) Walk(f func(string, *RicEntry) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	for key, entry := range t.Entries {
		if err := f(key, entry); err != nil {
			return err
		}
	}
	return nil
}

func NewNexthopKeyMic(ip net.IP) string {
	return fmt.Sprintf("%s", ip)
}

func NewNexthopKeyRic(ip net.IP, rt string) string {
	return fmt.Sprintf("%s@%s", ip, rt)
}

func NewNexthopKey(e *ribsmsg.Nexthop) string {
	if e.IsMic() {
		return NewNexthopKeyMic(e.Addr)
	} else {
		return NewNexthopKeyRic(e.Addr, e.Rt)
	}
}

//
// Nexthop Table (for detect loop)
//
type NexthopTable struct {
	Mutex   sync.RWMutex
	Entries map[string]*ribsmsg.Nexthop // Key: <RT>_<IP> or <IP>
}

func NewNexthopTable() *NexthopTable {
	return &NexthopTable{
		Entries: make(map[string]*ribsmsg.Nexthop),
	}
}

func (t *NexthopTable) AddMic(ip net.IP) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	ent := ribsmsg.NewNexthop("", ip, nil)
	t.Entries[NewNexthopKey(ent)] = ent
}

func (t *NexthopTable) AddRic(ip net.IP, rt string, srcId net.IP) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	ent := ribsmsg.NewNexthop(rt, ip, srcId)
	t.Entries[NewNexthopKey(ent)] = ent
}

func (t *NexthopTable) find(key string) *ribsmsg.Nexthop {
	if e, ok := t.Entries[key]; ok {
		return e
	}
	return nil
}

func (t *NexthopTable) FindMic(ip net.IP) *ribsmsg.Nexthop {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(NewNexthopKeyMic(ip))
}

func (t *NexthopTable) FindRic(ip net.IP, rt string) *ribsmsg.Nexthop {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(NewNexthopKeyRic(ip, rt))
}

func (t *NexthopTable) DelMic(ip net.IP) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	delete(t.Entries, NewNexthopKeyMic(ip))
}

func (t *NexthopTable) DelRic(ip net.IP, rt string) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	delete(t.Entries, NewNexthopKeyRic(ip, rt))
}

func (t *NexthopTable) Walk(f func(*ribsmsg.Nexthop) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	for _, entry := range t.Entries {
		if err := f(entry); err != nil {
			return err
		}
	}
	return nil
}
