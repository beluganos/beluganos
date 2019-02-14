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

package gonslib

import (
	"fabricflow/fibc/api"
	"sync"

	"github.com/beluganos/go-opennsl/opennsl"
	// log "github.com/sirupsen/logrus"
)

//
// IDMapL2Station has l2 station entries.
// key: some ID
// vak: L2 Station ID
//
type IDMapL2Station struct {
	sync.Map
}

//
// Register adds entry.
//
func (m *IDMapL2Station) Register(gID uint32, l2stationID opennsl.L2StationID) bool {
	_, ok := m.Map.LoadOrStore(gID, l2stationID)
	return ok
}

//
// Unregister remove entry.
//
func (m *IDMapL2Station) Unregister(gID uint32) {
	m.Map.Delete(gID)
}

//
// Get returns entry by group-id.
//
func (m *IDMapL2Station) Get(gID uint32) (opennsl.L2StationID, bool) {
	if v, ok := m.Map.Load(gID); ok {
		return v.(opennsl.L2StationID), ok
	}
	return 0, false
}

//
// Traverse enumerates all entries.
//
func (m *IDMapL2Station) Traverse(f func(uint32, opennsl.L2StationID) bool) {
	m.Map.Range(func(key, value interface{}) bool {
		return f(key.(uint32), value.(opennsl.L2StationID))
	})
}

//
// IDMapL3Iface has l3 interface entries.
// key: L2 Interface Group ID
// val: L3 Interface ID
//
type IDMapL3Iface struct {
	sync.Map
}

//
// Register adds entry(portID, vlanID, ifaceID).
//
func (m *IDMapL3Iface) Register(portID uint32, vid uint16, ifaceID opennsl.L3IfaceID) bool {
	gid := fibcapi.NewL2InterfaceGroupID(portID, vid)
	_, ok := m.Map.LoadOrStore(gid, ifaceID)
	return !ok
}

//
// Unregister removes entry.
//
func (m *IDMapL3Iface) Unregister(portID uint32, vid uint16) {
	gid := fibcapi.NewL2InterfaceGroupID(portID, vid)
	m.Map.Delete(gid)
}

//
// Get returns entry by portID
//
func (m *IDMapL3Iface) Get(portID uint32, vid uint16) (opennsl.L3IfaceID, bool) {
	gid := fibcapi.NewL2InterfaceGroupID(portID, vid)
	if v, ok := m.Map.Load(gid); ok {
		return v.(opennsl.L3IfaceID), ok
	}
	return 0, false
}

//
// Traverse enumerates all entries.
//
func (m *IDMapL3Iface) Traverse(f func(uint32, opennsl.L3IfaceID) bool) {
	m.Map.Range(func(key, value interface{}) bool {
		return f(key.(uint32), value.(opennsl.L3IfaceID))
	})
}

//
// IDMapL3Egress has l3 egress entries.
// key: L3 Unicast Group ID
// val: L3 Egress ID
//
type IDMapL3Egress struct {
	sync.Map
}

//
// Register adds entry(niID, egressId)
//
func (m *IDMapL3Egress) Register(neID uint32, egressID opennsl.L3EgressID) bool {
	gid := fibcapi.NewL3UnicastGroupID(neID)
	_, ok := m.Map.LoadOrStore(gid, egressID)
	return !ok
}

//
// Unregister removes entry.
//
func (m *IDMapL3Egress) Unregister(neID uint32) {
	gid := fibcapi.NewL3UnicastGroupID(neID)
	m.Map.Delete(gid)
}

//
// Get entry by neID.
//
func (m *IDMapL3Egress) Get(neID uint32) (opennsl.L3EgressID, bool) {
	gid := fibcapi.NewL3UnicastGroupID(neID)
	if v, ok := m.Map.Load(gid); ok {
		return v.(opennsl.L3EgressID), ok
	}
	return 0, false
}

//
// Traverse enumerates all entries.
//
func (m *IDMapL3Egress) Traverse(f func(uint32, opennsl.L3EgressID) bool) {
	m.Map.Range(func(key, value interface{}) bool {
		return f(key.(uint32), value.(opennsl.L3EgressID))
	})
}

//
// IDMaps has sub-maps.
//
type IDMaps struct {
	L2Stations IDMapL2Station
	L3Ifaces   IDMapL3Iface
	L3Egress   IDMapL3Egress
}

//
// NewIDMaps returns new instance.
//
func NewIDMaps() *IDMaps {
	return &IDMaps{}
}
