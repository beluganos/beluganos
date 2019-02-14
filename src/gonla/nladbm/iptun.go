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

package nladbm

import (
	"fmt"
	"gonla/nlalib"
	"gonla/nlamsg"
	"net"
	"sync"
)

//
// Key
//
type IptunKey struct {
	NId    uint8
	Remote string // ipv4 or ipv6 (net.IP.String())
}

func NewIptunKey(nid uint8, remote net.IP) *IptunKey {
	return &IptunKey{
		NId:    nid,
		Remote: remote.String(),
	}
}

func IptunToKey(iptun *nlamsg.Iptun) *IptunKey {
	if tun := iptun.Iptun(); tun != nil {
		return NewIptunKey(iptun.NId, tun.Remote)
	}

	return nil
}

func (k *IptunKey) String() string {
	return fmt.Sprintf("nid:%d, remote:'%s'", k.NId, k.Remote)
}

//
// Table Interface
//
type IptunTable interface {
	Insert(*nlamsg.Iptun) *nlamsg.Iptun
	Select(*IptunKey) *nlamsg.Iptun
	Delete(*IptunKey) *nlamsg.Iptun
	Update(*IptunKey, func(*nlamsg.Iptun) error) error
	Walk(func(*nlamsg.Iptun) error) error
	WalkFree(func(*nlamsg.Iptun) error) error
}

func NewIptunTable() IptunTable {
	return newIptunTable()
}

//
// Table
//
type iptunTable struct {
	mutex   sync.RWMutex
	iptuns  map[IptunKey]*nlamsg.Iptun
	counter *nlalib.Counters16
}

func newIptunTable() *iptunTable {
	return &iptunTable{
		iptuns:  map[IptunKey]*nlamsg.Iptun{},
		counter: nlalib.NewCounters16(),
	}
}

func (t *iptunTable) find(key *IptunKey) (iptun *nlamsg.Iptun) {
	if key == nil {
		return
	}

	iptun, _ = t.iptuns[*key]
	return
}

func (t *iptunTable) Insert(iptun *nlamsg.Iptun) (old *nlamsg.Iptun) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	key := IptunToKey(iptun)
	if key == nil {
		return
	}

	if old = t.find(key); old == nil {
		iptun.TnlId = t.counter.Next(iptun.NId)
	} else {
		iptun.TnlId = old.TnlId
	}

	t.iptuns[*key] = iptun.Copy()

	return
}

func (t *iptunTable) Select(key *IptunKey) *nlamsg.Iptun {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.find(key)
}

func (t *iptunTable) Delete(key *IptunKey) (tun *nlamsg.Iptun) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if tun = t.find(key); tun != nil {
		delete(t.iptuns, *key)
	}

	return
}

func (t *iptunTable) Update(key *IptunKey, f func(*nlamsg.Iptun) error) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	tun := t.find(key)
	if tun == nil {
		return fmt.Errorf("Tunnel not found. %s", key)
	}

	return f(tun)
}

func (t *iptunTable) Walk(f func(*nlamsg.Iptun) error) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.WalkFree(f)
}

func (t *iptunTable) WalkFree(f func(*nlamsg.Iptun) error) error {
	for _, link := range t.iptuns {
		if err := f(link); err != nil {
			return err
		}
	}
	return nil
}

//
// IptunPeer
//
type IptunPeer struct {
	NId uint8
	Dst *net.IPNet
}

func NewIptunPeerKey(nid uint8, dst *net.IPNet) string {
	return fmt.Sprintf("%d_%s", nid, dst)
}

func (r *IptunPeer) Key() string {
	return NewIptunPeerKey(r.NId, r.Dst)
}

func (r *IptunPeer) Contains(nid uint8, ip net.IP) bool {
	return (r.NId == nid) && r.Dst.Contains(ip)
}

func (r *IptunPeer) String() string {
	return fmt.Sprintf("nid:%d %s", r.NId, r.Dst)
}

func NewIptunPeer(nid uint8, dst *net.IPNet) *IptunPeer {
	return &IptunPeer{
		NId: nid,
		Dst: &net.IPNet{
			IP:   dst.IP,
			Mask: dst.Mask,
		},
	}
}

func (r *IptunPeer) PrefixLen() int {
	ones, _ := r.Dst.Mask.Size()
	return ones
}

//
// IptunPeerTable
//
type IptunPeerTable struct {
	entries map[string]*IptunPeer
}

func NewIptunPeerTable() *IptunPeerTable {
	return &IptunPeerTable{
		entries: map[string]*IptunPeer{},
	}
}

func (t *IptunPeerTable) find(route *IptunPeer) *IptunPeer {
	if e, ok := t.entries[route.Key()]; ok {
		return e
	}
	return nil
}

func (t *IptunPeerTable) Insert(route *IptunPeer) {
	t.entries[route.Key()] = route
}

func (t *IptunPeerTable) SelectByIP(nid uint8, ip net.IP) *IptunPeer {
	var curLen int
	var entry *IptunPeer
	for _, e := range t.entries {
		if e.Contains(nid, ip) {
			if plen := e.PrefixLen(); plen >= curLen {
				entry = e
				curLen = plen
			}
		}
	}
	return entry
}

func (t *IptunPeerTable) Delete(route *IptunPeer) *IptunPeer {
	if e := t.find(route); e != nil {
		delete(t.entries, route.Key())
		return e
	}
	return nil
}

func (t *IptunPeerTable) Range(f func(*IptunPeer)) {
	for _, e := range t.entries {
		f(e)
	}
}
