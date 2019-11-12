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
	"sort"
	"sync"
	"sync/atomic"
)

//
// StatsEntry is entry of stats.
//
type StatsEntry struct {
	cnt uint64
}

//
// NewStatsEntry returns new StatsEntry
//
func NewStatsEntry() *StatsEntry {
	return &StatsEntry{}
}

//
// Inc add count +1.
//
func (e *StatsEntry) Inc() {
	atomic.AddUint64(&e.cnt, 1)
}

//
// Add adds count +n
//
func (e *StatsEntry) Add(n uint64) {
	atomic.AddUint64(&e.cnt, n)
}

//
// Get returns count.
//
func (e *StatsEntry) Get() uint64 {
	return atomic.LoadUint64(&e.cnt)
}

//
// StatsGroup is set of StatsEntry.
//
type StatsGroup struct {
	entries map[string]*StatsEntry
	defent  *StatsEntry
	keys    []string
}

//
// NewStatsGroup returns new StatsGroup
//
func NewStatsGroup() *StatsGroup {
	return &StatsGroup{
		entries: map[string]*StatsEntry{},
		defent:  NewStatsEntry(),
		keys:    []string{},
	}
}

//
// Register add entry by name.
//
func (g *StatsGroup) Register(name string) {
	if _, ok := g.entries[name]; ok {
		return
	}

	e := NewStatsEntry()
	g.entries[name] = e
	g.keys = append(g.keys, name)
	sort.Strings(g.keys)
}

//
// RegisterList add entries by names.
//
func (g *StatsGroup) RegisterList(names []string) {
	for _, name := range names {
		g.Register(name)
	}
}

//
// Inc add count +1 by name.
//
func (g *StatsGroup) Inc(name string) {
	if e, ok := g.entries[name]; ok {
		e.Inc()
	} else {
		g.defent.Inc()
	}
}

//
// Add add count +n by name.
//
func (g *StatsGroup) Add(name string, n uint64) {
	if e, ok := g.entries[name]; ok {
		e.Add(n)
	} else {
		g.defent.Add(n)
	}
}

//
// Range retuens all count.
//
func (g *StatsGroup) Range(f func(string, uint64)) {
	for _, name := range g.keys {
		f(name, g.entries[name].Get())
	}

	if n := g.defent.Get(); n != 0 {
		f("*", n)
	}
}

//
// StatsTable is table of StatsGroup
//
type StatsTable struct {
	mutex  sync.RWMutex
	groups map[string]*StatsGroup
	defgrp *StatsGroup
	keys   []string
}

//
// NewStatsTable returns new StatsTable.
//
func NewStatsTable() *StatsTable {
	return &StatsTable{
		groups: map[string]*StatsGroup{},
		defgrp: NewStatsGroup(),
		keys:   []string{},
	}
}

//
// Register add stats group by name.
//
func (t *StatsTable) Register(name string) *StatsGroup {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if g, ok := t.groups[name]; ok {
		return g
	}

	g := NewStatsGroup()
	t.groups[name] = g
	t.keys = append(t.keys, name)
	sort.Strings(t.keys)

	return g
}

//
// Select returns StatsGroup.
//
func (t *StatsTable) Select(name string) *StatsGroup {
	if g, ok := t.groups[name]; ok {
		return g
	}

	return t.defgrp
}

//
// Range returns all counts.
//
func (t *StatsTable) Range(f func(string, string, uint64)) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, gname := range t.keys {
		g := t.groups[gname]
		g.Range(func(name string, n uint64) {
			f(gname, name, n)
		})
	}

	t.defgrp.Range(func(name string, n uint64) {
		f("*", name, n)
	})
}
