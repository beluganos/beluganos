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
	"fmt"
	"sync"

	"github.com/beluganos/go-opennsl/opennsl"
)

//
// IDMaps has sub-maps.
//
type IDMaps struct {
	L2Stations *L2StationIDMap
	L3Ifaces   *L3IfaceIDMap
	L3Egress   *L3EgressIDMap
	Trunks     *TrunkIDMap
}

//
// NewIDMaps returns new instance.
//
func NewIDMaps() *IDMaps {
	return &IDMaps{
		L2Stations: NewL2StationIDMap(),
		L3Ifaces:   NewL3IfaceIDMap(),
		L3Egress:   NewL3EgressIDMap(),
		Trunks:     NewTrunkIDMap(),
	}
}

//
// L2StationIDKey is key of L2StationIDMap
//
type L2StationIDKey uint32

//
// NewL2StationIDKey returns new L2StationIDKey
//
func NewL2StationIDKey(id uint32) L2StationIDKey {
	return L2StationIDKey(id)
}

//
// String returns string.
//
func (k L2StationIDKey) String() string {
	return fmt.Sprintf("0x%08x", uint32(k))
}

//
// L2StationIDMap has id and l2station id.
//
type L2StationIDMap struct {
	sync.Map
}

//
// NewL2StationIDMap returns new L2StationIDMap
//
func NewL2StationIDMap() *L2StationIDMap {
	return &L2StationIDMap{}
}

//
// Register registers id and l3egrId.
//
func (m *L2StationIDMap) Register(id uint32, l2stationId opennsl.L2StationID) bool {
	_, ok := m.Map.LoadOrStore(NewL2StationIDKey(id), l2stationId)
	return !ok
}

//
// Unregister removes id.
//
func (m *L2StationIDMap) Unregister(id uint32) {
	m.Map.Delete(NewL2StationIDKey(id))
}

//
// Get returns L2StationID by id.
//
func (m *L2StationIDMap) Get(id uint32) (opennsl.L2StationID, bool) {
	if v, ok := m.Map.Load(NewL2StationIDKey(id)); ok {
		return v.(opennsl.L2StationID), true
	}
	return 0, false
}

//
// Traverse enumerates all entries.
//
func (m *L2StationIDMap) Traverse(f func(L2StationIDKey, opennsl.L2StationID) bool) {
	m.Map.Range(func(key, value interface{}) bool {
		return f(key.(L2StationIDKey), value.(opennsl.L2StationID))
	})
}

//
// L3IfaceIDKey is key of L3IfaceIDMap.
//
type L3IfaceIDKey struct {
	Port uint32
	Vid  uint16
}

//
// NewL3IfaceIDKey returns new L3IfaceIDKey
//
func NewL3IfaceIDKey(port uint32, vid uint16) L3IfaceIDKey {
	return L3IfaceIDKey{
		Port: port,
		Vid:  vid,
	}
}

//
// String returns string.
//
func (k *L3IfaceIDKey) String() string {
	return fmt.Sprintf("port:%d vid:%d", k.Port, k.Vid)
}

//
// L3IfaceIDMap has id/vid(key) and l3-interface-id(value).
//
type L3IfaceIDMap struct {
	sync.Map
}

//
// NewL3IfaceIDMap returns new NewL3IfaceIDMap
//
func NewL3IfaceIDMap() *L3IfaceIDMap {
	return &L3IfaceIDMap{}
}

//
// Register registers id and l3ifaceId.
//
func (m *L3IfaceIDMap) Register(port uint32, vid uint16, l3ifaceId opennsl.L3IfaceID) bool {
	_, ok := m.Map.LoadOrStore(NewL3IfaceIDKey(port, vid), l3ifaceId)
	return !ok
}

//
// Unregister removes id.
//
func (m *L3IfaceIDMap) Unregister(port uint32, vid uint16) {
	m.Map.Delete(NewL3IfaceIDKey(port, vid))
}

//
// Get returns L3IfaceID by id.
//

func (m *L3IfaceIDMap) Get(port uint32, vid uint16) (opennsl.L3IfaceID, bool) {
	if v, ok := m.Map.Load(NewL3IfaceIDKey(port, vid)); ok {
		return v.(opennsl.L3IfaceID), true
	}
	return 0, false
}

//
// Traverse enumerates all entries.
//
func (m *L3IfaceIDMap) Traverse(f func(L3IfaceIDKey, opennsl.L3IfaceID) bool) {
	m.Map.Range(func(key, value interface{}) bool {
		return f(key.(L3IfaceIDKey), value.(opennsl.L3IfaceID))
	})
}

//
// L3EgressIDKey is key of L3EgressIDMap
//
type L3EgressIDKey uint32

//
// NewL3EgressIDKey returns new L3EgressIDKey
//
func NewL3EgressIDKey(neighId uint32) L3EgressIDKey {
	return L3EgressIDKey(neighId)
}

//
// String returns string.
//
func (k L3EgressIDKey) String() string {
	return fmt.Sprintf("0x%08x", uint32(k))
}

//
// L3EgressIDMap has id(key) and l3-egress-id(value).
//
type L3EgressIDMap struct {
	sync.Map
}

//
// NewL3EgressIDMap returns new L3EgressIDMap
//
func NewL3EgressIDMap() *L3EgressIDMap {
	return &L3EgressIDMap{}
}

//
// Register registers id and l3egrId.
//
func (m *L3EgressIDMap) Register(id uint32, l3egrId opennsl.L3EgressID) bool {
	_, ok := m.Map.LoadOrStore(NewL3EgressIDKey(id), l3egrId)
	return !ok
}

//
// Unregister removes id.
//
func (m *L3EgressIDMap) Unregister(id uint32) {
	m.Map.Delete(NewL3EgressIDKey(id))
}

//
// Get returns L3EgressID by id.
//
func (m *L3EgressIDMap) Get(id uint32) (opennsl.L3EgressID, bool) {
	if v, ok := m.Map.Load(NewL3EgressIDKey(id)); ok {
		return v.(opennsl.L3EgressID), true
	}
	return 0, false
}

//
// Traverse enumerates all entries.
//
func (m *L3EgressIDMap) Traverse(f func(L3EgressIDKey, opennsl.L3EgressID) bool) {
	m.Map.Range(func(key, value interface{}) bool {
		return f(key.(L3EgressIDKey), value.(opennsl.L3EgressID))
	})
}

//
// TrunkIDKey is key of TrunkIDMap
//
type TrunkIDKey struct {
	LagId uint32
	Vid   uint16
}

//
// NewTrunkIDKey returns new TrunkIDKey
//
func NewTrunkIDKey(lagId uint32, vid uint16) TrunkIDKey {
	return TrunkIDKey{
		LagId: lagId,
		Vid:   vid,
	}
}

//
// String returns string
//
func (k *TrunkIDKey) String() string {
	return fmt.Sprintf("lag:0x%08x vid:%d", k.LagId, k.Vid)
}

//
// TrunkIDMap has lag-id/vid(key) and trunk-id(value)
//
type TrunkIDMap struct {
	sync.Map
}

//
// NewTrunkIDMap returns new TrunkIDMap
//
func NewTrunkIDMap() *TrunkIDMap {
	return &TrunkIDMap{}
}

//
// Register registers id and trunkId.
//
func (m *TrunkIDMap) Register(lagid uint32, vid uint16, trunkId opennsl.Trunk) bool {
	_, ok := m.Map.LoadOrStore(NewTrunkIDKey(lagid, vid), trunkId)
	return !ok
}

//
// Unregister removes id.
//
func (m *TrunkIDMap) Unregister(lagid uint32, vid uint16) {
	m.Map.Delete(NewTrunkIDKey(lagid, vid))
}

//
// Get returns L3IfaceID by id.
//
func (m *TrunkIDMap) Get(lagid uint32, vid uint16) (opennsl.Trunk, bool) {
	if v, ok := m.Map.Load(NewTrunkIDKey(lagid, vid)); ok {
		return v.(opennsl.Trunk), true
	}
	return 0, false
}

//
// Traverse enumerates all entries.
//
func (m *TrunkIDMap) Traverse(f func(TrunkIDKey, opennsl.Trunk) bool) {
	m.Map.Range(func(key, value interface{}) bool {
		return f(key.(TrunkIDKey), value.(opennsl.Trunk))
	})
}
