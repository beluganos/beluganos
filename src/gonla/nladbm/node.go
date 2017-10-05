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
	"gonla/nlamsg"
	"sync"
)

//
// Key
//
type NodeKey struct {
	NId uint8
}

func NewNodeKey(nid uint8) *NodeKey {
	return &NodeKey{
		NId: nid,
	}
}

func NodeToKey(n *nlamsg.Node) *NodeKey {
	return NewNodeKey(n.NId)
}

//
// Table interface
//
type NodeTable interface {
	Insert(*nlamsg.Node) *nlamsg.Node
	Select(*NodeKey) *nlamsg.Node
	Delete(*NodeKey) *nlamsg.Node
	Walk(f func(*nlamsg.Node) error) error
	WalkUnsafe(f func(*nlamsg.Node) error) error
}

func NewNodeTable() NodeTable {
	return &nodeTable{
		Nodes: make(map[NodeKey]*nlamsg.Node),
	}
}

//
// Table
//
type nodeTable struct {
	Mutex sync.RWMutex
	Nodes map[NodeKey]*nlamsg.Node
}

func (t *nodeTable) find(key *NodeKey) *nlamsg.Node {
	n, _ := t.Nodes[*key]
	return n
}

func (t *nodeTable) Insert(n *nlamsg.Node) (old *nlamsg.Node) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	key := NodeToKey(n)
	if old = t.find(key); old == nil {
		n.Open()
		t.Nodes[*key] = n.Copy()
	}

	return
}

func (t *nodeTable) Select(key *NodeKey) *nlamsg.Node {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.find(key)
}

func (t *nodeTable) Walk(f func(*nlamsg.Node) error) error {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()

	return t.WalkUnsafe(f)
}

func (t *nodeTable) WalkUnsafe(f func(*nlamsg.Node) error) error {
	for _, n := range t.Nodes {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (t *nodeTable) Delete(key *NodeKey) (old *nlamsg.Node) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if old = t.find(key); old != nil {
		delete(t.Nodes, *key)
		old.Close()
	}

	return
}
