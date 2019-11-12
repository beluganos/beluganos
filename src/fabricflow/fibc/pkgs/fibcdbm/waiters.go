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
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// WaitTimeout is timeout of waiter.
var WaitTimeout = errors.New("timeout")

//
// Waiter is used to sync request and async response
//
type Waiter interface {
	Set(interface{})
	SetError(error)
	Wait(time.Duration) error
}

//
// SimpleWaiter is simple waiter.
//
type SimpleWaiter struct {
	once sync.Once
	done chan struct{}
	err  error
}

//
// NewSimpleWaiter returns new Simple waiter
//
func NewSimpleWaiter() *SimpleWaiter {
	return &SimpleWaiter{
		done: make(chan struct{}),
	}
}

//
// Set sets wait signal.
//
func (w *SimpleWaiter) Set(v interface{}) {
	w.Close()
}

//
// Set error
//
func (w *SimpleWaiter) SetError(err error) {
	w.err = err
	w.Close()
}

//
// Close closes chan.
//
func (w *SimpleWaiter) Close() {
	w.once.Do(func() {
		close(w.done)
	})
}

//
// Wait blocks until Set is called or timeout.
//
func (w *SimpleWaiter) Wait(d time.Duration) error {
	t := time.NewTimer(d)
	defer func() {
		t.Stop()
		w.Close()
	}()

	select {
	case <-t.C:
		return WaitTimeout

	case <-w.done:
		return w.err
	}
}

//
// WaiterTable is table of waiters.
//
type WaiterTable struct {
	mutex   sync.RWMutex
	xid     uint32
	waiters map[uint32]Waiter
}

//
// NewWaiterTable returns new WaiterTable
//
func NewWaiterTable() *WaiterTable {
	return &WaiterTable{
		waiters: map[uint32]Waiter{},
	}
}

func (t *WaiterTable) nextXid() uint32 {
	return atomic.AddUint32(&t.xid, 1)
}

func (t *WaiterTable) find(xid uint32) (Waiter, bool) {
	w, ok := t.waiters[xid]
	return w, ok
}

//
// Register add waiter.
//
func (t *WaiterTable) Register(w Waiter) (xid uint32) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	xid = t.nextXid()
	t.waiters[xid] = w

	return
}

//
// Unregister deletes waiter.
//
func (t *WaiterTable) Unregister(xid uint32) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, ok := t.find(xid); ok {
		delete(t.waiters, xid)
	}
}

//
// Select return waiter.
//
func (t *WaiterTable) Select(xid uint32, f func(Waiter)) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if w, ok := t.find(xid); ok {
		f(w)
		return true
	}

	return false
}

//
// Range returns all waiters.
//
func (t *WaiterTable) Range(f func(uint32, Waiter) bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for xid, w := range t.waiters {
		if ok := f(xid, w); !ok {
			return
		}
	}
}
