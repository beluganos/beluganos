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

package nladbm

import (
	"gonla/nlalib"
	"gonla/nlamsg"
	"net"
	"sync"
)

//
// Key
//
type VpnKey struct {
	NId uint8
	Dst string // ip/mask (net.IPNet.String())
	Gw  string // ip (net.IP.String())
}

func NewVpnKey(nid uint8, dst *net.IPNet, gw net.IP) *VpnKey {
	return &VpnKey{
		NId: nid,
		Dst: dst.String(),
		Gw:  gw.String(),
	}
}

func VpnToKey(v *nlamsg.Vpn) *VpnKey {
	return NewVpnKey(v.NId, v.Vpn.GetIPNet(), v.Vpn.NetGw())
}

//
// Table interface
//
type VpnTable interface {
	Insert(*nlamsg.Vpn) *nlamsg.Vpn
	Select(*VpnKey) *nlamsg.Vpn
	Delete(*VpnKey) *nlamsg.Vpn
	Walk(f func(*nlamsg.Vpn) error) error
	WalkFree(f func(*nlamsg.Vpn) error) error
	WalkByGw(gw net.IP, f func(*nlamsg.Vpn) error) error
	WalkByGwFree(gw net.IP, f func(*nlamsg.Vpn) error) error
	WalkByVpnGw(vpnGw net.IP, f func(*nlamsg.Vpn) error) error
	WalkByVpnGwFree(vpnGw net.IP, f func(*nlamsg.Vpn) error) error
}

func NewVpnTable() VpnTable {
	return newVpnTable()
}

//
// Table
//
type vpnTable struct {
	Mutex    sync.RWMutex
	Vpns     map[VpnKey]*nlamsg.Vpn
	Counter  *nlalib.Counter32
	GwIdx    *VpnGwIndex
	VpnGwIdx *VpnVpnGwIndex
}

func newVpnTable() *vpnTable {
	return &vpnTable{
		Vpns:     make(map[VpnKey]*nlamsg.Vpn),
		Counter:  &nlalib.Counter32{},
		GwIdx:    NewVpnGwIndex(),
		VpnGwIdx: NewVpnVpnGwIndex(),
	}
}

func (t *vpnTable) find(key *VpnKey) *nlamsg.Vpn {
	n, _ := t.Vpns[*key]
	return n
}

func (t *vpnTable) Insert(v *nlamsg.Vpn) (old *nlamsg.Vpn) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := VpnToKey(v)
	if old = t.find(key); old == nil {
		v.VpnId = t.Counter.Next()
	} else {
		v.VpnId = old.VpnId
	}

	t.Vpns[*key] = v.Copy()

	if old != nil {
		t.GwIdx.Delete(old)
		t.VpnGwIdx.Delete(old)
	}
	t.GwIdx.Insert(v)
	t.VpnGwIdx.Insert(v)

	return
}

func (t *vpnTable) Select(key *VpnKey) *nlamsg.Vpn {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *vpnTable) Delete(key *VpnKey) (old *nlamsg.Vpn) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.Vpns, *key)
		t.GwIdx.Delete(old)
		t.VpnGwIdx.Delete(old)
	}

	return
}

func (t *vpnTable) Walk(f func(*nlamsg.Vpn) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkFree(f)
}

func (t *vpnTable) WalkFree(f func(*nlamsg.Vpn) error) error {
	for _, n := range t.Vpns {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (t *vpnTable) WalkByGw(gw net.IP, f func(*nlamsg.Vpn) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkByGwFree(gw, f)
}

func (t *vpnTable) WalkByGwFree(gw net.IP, f func(*nlamsg.Vpn) error) error {
	e, ok := t.GwIdx.Select(gw)
	if ok {
		for _, key := range e.Keys {
			if vpn := t.find(key); vpn != nil {
				if err := f(vpn); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *vpnTable) WalkByVpnGw(vpnGw net.IP, f func(*nlamsg.Vpn) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkByVpnGwFree(vpnGw, f)
}

func (t *vpnTable) WalkByVpnGwFree(vpnGw net.IP, f func(*nlamsg.Vpn) error) error {
	e, ok := t.VpnGwIdx.Select(vpnGw)
	if ok {
		for _, key := range e.Keys {
			if vpn := t.find(key); vpn != nil {
				if err := f(vpn); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//
// Index Entry
//
type VpnGwIndexEntry struct {
	Keys map[VpnKey]*VpnKey
}

func NewVpnGwIndexEntry() *VpnGwIndexEntry {
	return &VpnGwIndexEntry{
		Keys: make(map[VpnKey]*VpnKey),
	}
}

func (e *VpnGwIndexEntry) Insert(k *VpnKey) {
	e.Keys[*k] = k
}

func (e *VpnGwIndexEntry) Delete(k *VpnKey) {
	delete(e.Keys, *k)
}

func (e *VpnGwIndexEntry) Len() int {
	return len(e.Keys)
}

//
// Index(Gw)
//
type VpnGwIndex struct {
	Entry map[string]*VpnGwIndexEntry // key: net.IP.String()
}

func NewVpnGwIndex() *VpnGwIndex {
	return &VpnGwIndex{
		Entry: make(map[string]*VpnGwIndexEntry),
	}
}

func (i *VpnGwIndex) Insert(v *nlamsg.Vpn) {
	dkey := v.Vpn.NetGw().String()
	e, ok := i.Entry[dkey]
	if !ok {
		e = NewVpnGwIndexEntry()
		i.Entry[dkey] = e
	}

	e.Insert(VpnToKey(v))
}

func (i *VpnGwIndex) Delete(v *nlamsg.Vpn) {
	dkey := v.Vpn.NetGw().String()
	e, ok := i.Entry[dkey]
	if !ok {
		return
	}

	e.Delete(VpnToKey(v))

	if e.Len() == 0 {
		delete(i.Entry, dkey)
	}
}

func (i *VpnGwIndex) Select(gw net.IP) (e *VpnGwIndexEntry, ok bool) {
	e, ok = i.Entry[gw.String()]
	return
}

//
// Index(VpnGw)
//
type VpnVpnGwIndex struct {
	Entry map[string]*VpnGwIndexEntry // key: net.IP.String()
}

func NewVpnVpnGwIndex() *VpnVpnGwIndex {
	return &VpnVpnGwIndex{
		Entry: make(map[string]*VpnGwIndexEntry),
	}
}

func (i *VpnVpnGwIndex) Insert(v *nlamsg.Vpn) {
	dkey := v.Vpn.NetVpnGw().String()
	e, ok := i.Entry[dkey]
	if !ok {
		e = NewVpnGwIndexEntry()
		i.Entry[dkey] = e
	}

	e.Insert(VpnToKey(v))
}

func (i *VpnVpnGwIndex) Delete(v *nlamsg.Vpn) {
	dkey := v.Vpn.NetVpnGw().String()
	e, ok := i.Entry[dkey]
	if !ok {
		return
	}

	e.Delete(VpnToKey(v))

	if e.Len() == 0 {
		delete(i.Entry, dkey)
	}
}

func (i *VpnVpnGwIndex) Select(vpnGw net.IP) (e *VpnGwIndexEntry, ok bool) {
	e, ok = i.Entry[vpnGw.String()]
	return
}
