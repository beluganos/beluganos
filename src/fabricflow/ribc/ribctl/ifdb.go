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

package ribctl

import (
	fibcapi "fabricflow/fibc/api"
	"fmt"
	"gonla/nlamsg"
	"strings"
	"sync"
)

type IfDBField uint32

const (
	IfDBFieldNone IfDBField = 0
	IfDBFieldAny  IfDBField = 1 << iota
	IfDBFieldStatus
	IfDBFieldMaster
	IfDBFieldLinkType
	IfDBFieldLinkID
)

var ifdbField_names = map[IfDBField]string{
	IfDBFieldAny:      "any",
	IfDBFieldStatus:   "status",
	IfDBFieldMaster:   "master",
	IfDBFieldLinkType: "linkYype",
	IfDBFieldLinkID:   "linkId",
}

func (v IfDBField) String() string {
	if v.IsNull() {
		return "none"
	}

	names := []string{}
	for field, name := range ifdbField_names {
		if (field & v) != 0 {
			names = append(names, name)
		}
	}

	return strings.Join(names, "|")
}

func (v IfDBField) IsNull() bool {
	return v == IfDBFieldNone
}

func (v IfDBField) Has(fields ...IfDBField) bool {
	for _, field := range fields {
		if (v & field) == field {
			return true
		}
	}

	return false
}

func NewIfDBKey(nid uint8, ifindex int) string {
	return fmt.Sprintf("%d@%d", nid, ifindex)
}

type IfDBEntry struct {
	NId         uint8
	LnId        uint16
	Index       int
	MasterIndex int
	PortStatus  fibcapi.PortStatus_Status
	Associated  bool
	LinkType    fibcapi.LinkType_Type
}

func NewIfDBEntryFromLink(link *nlamsg.Link) *IfDBEntry {
	return &IfDBEntry{
		NId:         link.NId,
		LnId:        link.LnId,
		Index:       link.Attrs().Index,
		MasterIndex: link.Attrs().MasterIndex,
		PortStatus:  NewPortStatus(link),
		Associated:  false,
		LinkType:    LinkTypeFromLink(link),
	}
}

func (e *IfDBEntry) CopyTo(dst *IfDBEntry) {
	*dst = *e
}

func (e *IfDBEntry) PortId() uint32 {
	return PortId(e.NId, e.LnId)
}

func (e *IfDBEntry) Key() string {
	return NewIfDBKey(e.NId, e.Index)
}

func (e *IfDBEntry) Update(link *nlamsg.Link) IfDBField {
	fields := IfDBFieldNone

	if lnId := link.LnId; lnId != e.LnId {
		e.LnId = lnId
		fields |= IfDBFieldLinkID
	}

	if master := link.Attrs().MasterIndex; master != e.MasterIndex {
		e.MasterIndex = master
		fields |= IfDBFieldMaster
	}

	if pstatus := NewPortStatus(link); pstatus != e.PortStatus {
		e.PortStatus = pstatus
		fields |= IfDBFieldStatus
	}

	if linkType := LinkTypeFromLink(link); linkType != e.LinkType {
		e.LinkType = linkType
		fields |= IfDBFieldLinkType
	}

	return fields
}

func (e *IfDBEntry) String() string {
	return fmt.Sprintf("nid:%d index:%d lnId:%d port:%08x %s %t", e.NId, e.Index, e.LnId, e.PortId(), e.LinkType, e.Associated)
}

type IfDB struct {
	lock    sync.RWMutex
	entries map[string]*IfDBEntry
	portIds map[uint32]*IfDBEntry
}

func NewIfDB() *IfDB {
	return &IfDB{
		entries: map[string]*IfDBEntry{},
		portIds: map[uint32]*IfDBEntry{},
	}
}

func (db *IfDB) clear() {
	db.entries = map[string]*IfDBEntry{}
	db.portIds = map[uint32]*IfDBEntry{}
}

func (db *IfDB) find(key string) *IfDBEntry {
	if e, ok := db.entries[key]; ok {
		return e
	}
	return nil
}

func (db *IfDB) findByPortId(portId uint32) *IfDBEntry {
	if e, ok := db.portIds[portId]; ok {
		return e
	}
	return nil
}

func (db *IfDB) add(e *IfDBEntry) {
	db.entries[e.Key()] = e
	db.portIds[e.PortId()] = e
}

func (db *IfDB) del(key string) *IfDBEntry {
	if e := db.find(key); e != nil {
		delete(db.entries, key)
		delete(db.portIds, e.PortId())
		return e
	}

	return nil
}

func (db *IfDB) delByPortId(portId uint32) *IfDBEntry {
	if e := db.findByPortId(portId); e != nil {
		return db.del(e.Key())
	}

	return nil
}

func (db *IfDB) Clear() {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.clear()
}

func (db *IfDB) Set(e *IfDBEntry) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.add(e)
}

func (db *IfDB) Delete(portId uint32) *IfDBEntry {
	db.lock.Lock()
	defer db.lock.Unlock()

	return db.delByPortId(portId)
}

func (db *IfDB) SelectBy(ifentry *IfDBEntry, nid uint8, ifindex int) bool {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if e := db.find(NewIfDBKey(nid, ifindex)); e != nil {
		e.CopyTo(ifentry)
		return true
	}

	return false
}

func (db *IfDB) Select(ifentry *IfDBEntry, portId uint32) bool {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if e := db.findByPortId(portId); e != nil {
		e.CopyTo(ifentry)
		return true
	}

	return false
}

func (db *IfDB) ListSlaves(masterIndex int, f func(*IfDBEntry)) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	for _, e := range db.entries {
		if e.MasterIndex == masterIndex {
			f(e)
		}
	}
}

func (db *IfDB) Update(portId uint32, f func(e *IfDBEntry) IfDBField) IfDBField {
	db.lock.Lock()
	defer db.lock.Unlock()

	if e := db.findByPortId(portId); e != nil {
		return f(e)
	}

	return IfDBFieldNone
}

func (db *IfDB) Associated(nid uint8, ifindex int) bool {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if e := db.find(NewIfDBKey(nid, ifindex)); e != nil {
		return e.Associated
	}

	return false
}
