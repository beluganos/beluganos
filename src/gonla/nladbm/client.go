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

type Client chan *nlamsg.NetlinkMessageUnion

type ClientTable interface {
	New() Client
	Delete(Client)
	Send(*nlamsg.NetlinkMessageUnion)
}

func NewClientTable() ClientTable {
	return &clientTable{
		clients: make(map[Client]struct{}),
	}
}

//
// Table
//
type clientTable struct {
	clients map[Client]struct{}
	Mutex   sync.Mutex
}

func (t *clientTable) New() Client {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	c := make(chan *nlamsg.NetlinkMessageUnion)
	t.clients[c] = struct{}{}
	return c
}

func (t *clientTable) Delete(c Client) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	delete(t.clients, c)
	close(c)
}

func (t *clientTable) Send(msg *nlamsg.NetlinkMessageUnion) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	t.WalkFree(func(c Client) {
		c <- msg
	})
}

func (t *clientTable) WalkFree(f func(Client)) {
	for c, _ := range t.clients {
		f(c)
	}
}
