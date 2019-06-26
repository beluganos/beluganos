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

package gonslib

import (
	"fmt"
	"sync"
	"time"

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

func (s *Server) L2AddrInit() {
	agingTime := s.dpCfg.GetL2AgingTimer()
	if agingTime > 0 {
		if err := opennsl.L2AddrAgeTimerSet(s.Unit(), agingTime); err != nil {
			log.Warnf("Server: L2AddrInit AgeTimeSet error. %s", err)
		}

		log.Infof("Server: L2AddrInit AgeTime %d", agingTime)
	}

	log.Infof("Server: L2AddrInit ok.")
}

type L2addrmonEntry struct {
	L2Addr *opennsl.L2Addr
	Oper   opennsl.L2CallbackOper
}

func NewL2addrmonEntry(src *opennsl.L2Addr, oper opennsl.L2CallbackOper) *L2addrmonEntry {
	l2addr := *src
	return &L2addrmonEntry{
		L2Addr: &l2addr,
		Oper:   oper,
	}
}

func (e *L2addrmonEntry) SetAdd() bool {
	switch e.Oper {
	case opennsl.L2_CALLBACK_ADD:
		return true

	case opennsl.L2_CALLBACK_DELETE:
		e.Oper = opennsl.L2_CALLBACK_NONE
		return false

	default:
		// oper is NONE or other.
		e.Oper = opennsl.L2_CALLBACK_ADD
		return true
	}
}

func (e *L2addrmonEntry) SetDel() bool {
	switch e.Oper {
	case opennsl.L2_CALLBACK_ADD:
		e.Oper = opennsl.L2_CALLBACK_NONE
		return false

	case opennsl.L2_CALLBACK_DELETE:
		return true

	default:
		// oper is NONE or other.
		e.Oper = opennsl.L2_CALLBACK_DELETE
		return true
	}
}

func (e *L2addrmonEntry) Key() string {
	return NewL2addrmonKeyFromL2Addr(e.L2Addr)
}

func (e *L2addrmonEntry) String() string {
	return fmt.Sprintf("%s %s", e.Oper, e.L2Addr)
}

func NewL2addrmonKeyFromL2Addr(l2addr *opennsl.L2Addr) string {
	return fmt.Sprintf("%s_%d_%d", l2addr.MAC(), l2addr.VID(), l2addr.Port())
}

type L2addrmonTable struct {
	entries  map[string]*L2addrmonEntry
	mutex    sync.Mutex
	maxEntry uint32
}

func NewL2addrmonTable(maxEntry uint32) *L2addrmonTable {
	return &L2addrmonTable{
		entries:  map[string]*L2addrmonEntry{},
		maxEntry: maxEntry,
	}
}

func (t *L2addrmonTable) find(key string) *L2addrmonEntry {
	e, _ := t.entries[key]
	return e
}

func (t *L2addrmonTable) del(key string) {
	delete(t.entries, key)
}

func (t *L2addrmonTable) count() uint32 {
	return uint32(len(t.entries))
}

func (t *L2addrmonTable) Add(newEntry *L2addrmonEntry) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	key := newEntry.Key()
	if e := t.find(key); e != nil {
		if hold := e.SetAdd(); hold {
			e.L2Addr = newEntry.L2Addr
		} else {
			t.del(key)
		}
	} else {
		newEntry.SetAdd()
		t.entries[key] = newEntry
	}
}

func (t *L2addrmonTable) Del(delEntry *L2addrmonEntry) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	key := delEntry.Key()
	if e := t.find(key); e != nil {
		if hold := e.SetDel(); hold {
			e.L2Addr = delEntry.L2Addr
		} else {
			t.del(key)
		}
	} else {
		delEntry.SetDel()
		t.entries[key] = delEntry
	}
}

func (t *L2addrmonTable) Put(e *L2addrmonEntry) {
	switch e.Oper {
	case opennsl.L2_CALLBACK_ADD:
		t.Add(e)

	case opennsl.L2_CALLBACK_DELETE:
		t.Del(e)

	default:
		// do nothind
	}
}

func (t *L2addrmonTable) Reset(force bool) []*L2addrmonEntry {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var entries []*L2addrmonEntry
	if n := t.count(); force || n >= t.maxEntry {
		var index int

		entries = make([]*L2addrmonEntry, n)
		for _, e := range t.entries {
			entries[index] = e
			index++
		}

		t.entries = map[string]*L2addrmonEntry{}
	}

	return entries
}

func (s *Server) L2AddrMonStart(done <-chan struct{}) {
	s.L2AddrInit()

	sweepTime := s.dpCfg.GetL2SweepTime()
	go s.L2AddrMonServe(s.Unit(), sweepTime, done)
}

func (s *Server) L2AddrMonServe(unit int, sweepTime time.Duration, done <-chan struct{}) {
	ch := make(chan *L2addrmonEntry)
	defer close(ch)

	if err := opennsl.L2AddrRegister(unit, func(unitCb int, l2addr *opennsl.L2Addr, oper opennsl.L2CallbackOper) {
		ch <- NewL2addrmonEntry(l2addr, oper)
	}); err != nil {
		log.Errorf("L2AddrMon: L2AddrRegister error. %s", err)
		return
	}

	defer opennsl.L2AddrUnregister(unit)

	log.Infof("L2AddrMon: Started.")

	tbl := NewL2addrmonTable(s.dpCfg.GetL2NotifyLimit())
	doNotify := func(force bool) {
		if entries := tbl.Reset(force); len(entries) > 0 {
			log.Debugf("Server: notify #%d force:%t", len(entries), force)
			s.l2addrCh <- entries
		}
	}

	tick := time.NewTicker(sweepTime)
	defer tick.Stop()

FOR_LABEL:
	for {
		select {
		case e := <-ch:
			log.Debugf("L2AddrMon: entry %s", e)
			tbl.Put(e)
			doNotify(false)

		case <-tick.C:
			// log.Debugf("L2AddrMon: all")
			doNotify(true)

		case <-done:
			log.Infof("L2AddrMon: Exit.")
			break FOR_LABEL
		}
	}
}
