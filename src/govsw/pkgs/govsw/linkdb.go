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

package govsw

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type LinkDB struct {
	links map[int]*Link

	stripWPkt uint16
	stripRPkt uint16

	mutex sync.RWMutex
	log   *log.Entry
}

func NewLinkDB() *LinkDB {
	return &LinkDB{
		links: make(map[int]*Link),
		log:   log.WithField("module", "linkdb"),
	}
}

func (db *LinkDB) find(ifindex int) (*Link, bool) {
	link, ok := db.links[ifindex]
	return link, ok
}

func (db *LinkDB) findByName(ifname string) (*Link, bool) {
	for _, link := range db.links {
		if name := link.Name(); name == ifname {
			return link, true
		}
	}

	return nil, false
}

func (db *LinkDB) SetStripWPkt(v uint16) {
	db.stripWPkt = v
}

func (db *LinkDB) SetStripRPkt(v uint16) {
	db.stripRPkt = v
}

func (db *LinkDB) GetByName(ifname string, f func(*Link) error) error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if link, ok := db.findByName(ifname); ok {
		return f(link)
	}

	return fmt.Errorf("Link not found. %s", ifname)
}

func (db *LinkDB) Get(ifindex int, f func(*Link) error) error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if link, ok := db.find(ifindex); ok {
		return f(link)
	}

	return fmt.Errorf("Link not found. %d", ifindex)
}

func (db *LinkDB) GetOrAdd(ln netlink.Link, f func(*Link)) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	ifindex := ln.Attrs().Index
	link, ok := db.find(ifindex)
	if !ok {
		link = NewLink(ln)
		link.SetStripWPkt(db.stripWPkt)
		link.SetStripRPkt(db.stripRPkt)
		db.links[ifindex] = link
	}

	f(link)
}

func (db *LinkDB) Add(link *Link) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	ifindex := link.Index()
	db.links[ifindex] = link
}

func (db *LinkDB) Delete(ifindex int) (*Link, bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if link, ok := db.find(ifindex); ok {
		delete(db.links, ifindex)
		return link, true
	}

	return nil, false
}

func (db *LinkDB) Range(f func(int, *Link)) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for index, link := range db.links {
		f(index, link)
	}
}
