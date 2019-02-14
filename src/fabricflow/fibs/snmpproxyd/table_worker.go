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

//
// WorkerTableEntry is entry for WorkerTable
//
type WorkerTableEntry struct {
	worker *ProxyWorker
	ttl    int
}

const (
	SERVER_WORKER_TTL_DEFAULT = 5
)

//
// NewWorkerTableEntry creates new WorkerTableEntry instance.
//
func NewWorkerTableEntry(worker *ProxyWorker) *WorkerTableEntry {
	return &WorkerTableEntry{
		worker: worker,
		ttl:    SERVER_WORKER_TTL_DEFAULT,
	}
}

//
// Reset resets ttl.
//
func (e *WorkerTableEntry) Reset() {
	e.ttl = SERVER_WORKER_TTL_DEFAULT
}

//
// Check check and decrement ttl.
//
func (e *WorkerTableEntry) CheckAlive() bool {
	if e.ttl == 0 {
		return false
	}
	e.ttl--
	return true
}

//
// WorkerTable is workers table.
//
type WorkerTable struct {
	mutex   sync.Mutex
	entries map[string]*WorkerTableEntry
}

func NewWorkerTable() *WorkerTable {
	return &WorkerTable{
		entries: map[string]*WorkerTableEntry{},
	}
}

func (t *WorkerTable) WriteTo(w io.Writer) (sum int64, err error) {
	for key, val := range t.entries {
		var n int
		n, err = fmt.Fprintf(w, "%s:%d\n", key, val.ttl)
		sum += int64(n)
		if err != nil {
			return
		}
	}
	return
}

func (t *WorkerTable) Put(addr string, w *ProxyWorker) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.entries[addr] = NewWorkerTableEntry(w)
}

func (t *WorkerTable) Find(addr string) (*ProxyWorker, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if e, ok := t.entries[addr]; ok {
		e.Reset()
		return e.worker, true
	}
	return nil, false
}

func (t *WorkerTable) CheckAlive(f func(string, *ProxyWorker)) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for addr, entry := range t.entries {
		if ok := entry.CheckAlive(); !ok {
			delete(t.entries, addr)
			f(addr, entry.worker)
		}
	}
}
