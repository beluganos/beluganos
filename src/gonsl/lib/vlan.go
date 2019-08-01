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
	fibcapi "fabricflow/fibc/api"
	"fmt"
	"sync"

	"github.com/beluganos/go-opennsl/opennsl"
)

type VlanPortKey struct {
	Port opennsl.Port
	Vid  opennsl.Vlan
}

func NewVlanPortKey(port opennsl.Port, vid opennsl.Vlan) *VlanPortKey {
	return &VlanPortKey{
		Port: port,
		Vid:  vid,
	}
}

func (k *VlanPortKey) String() string {
	return fmt.Sprintf("%d_%d", k.Port, k.Vid)
}

type VlanPortEntry struct {
	Vid opennsl.Vlan
}

func NewVlanPortEntry(vid opennsl.Vlan) *VlanPortEntry {
	return &VlanPortEntry{
		Vid: vid,
	}
}

func NewVlanPortEntryDefault() *VlanPortEntry {
	return NewVlanPortEntry(opennsl.VLAN_ID_DEFAULT)
}

type VlanPortTable struct {
	entries    map[string]*VlanPortEntry
	mutex      sync.RWMutex
	minPort    opennsl.Port
	maxPort    opennsl.Port
	baseVid    opennsl.Vlan
	defaultVid opennsl.Vlan
}

func NewVlanPortTable() *VlanPortTable {
	return &VlanPortTable{
		entries:    map[string]*VlanPortEntry{},
		minPort:    0,
		maxPort:    0,
		baseVid:    1,
		defaultVid: opennsl.VLAN_ID_DEFAULT,
	}
}

func (t *VlanPortTable) SetMinPort(minPort opennsl.Port) {
	t.minPort = minPort
}

func (t *VlanPortTable) SetMaxPort(maxPort opennsl.Port) {
	t.maxPort = maxPort
}

func (t *VlanPortTable) SetBaseVID(vid opennsl.Vlan) {
	t.baseVid = vid
}

func (t *VlanPortTable) SetDefaultVID(vid opennsl.Vlan) {
	t.defaultVid = vid
}

func (t *VlanPortTable) find(key string) *VlanPortEntry {
	if e, ok := t.entries[key]; ok {
		return e
	}
	return nil
}

func (t *VlanPortTable) Insert(key *VlanPortKey, entry *VlanPortEntry) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	k := key.String()

	if e := t.find(k); e != nil {
		return false
	}

	t.entries[k] = entry

	return true
}

func (t *VlanPortTable) ConvVID(port opennsl.Port, vid opennsl.Vlan) opennsl.Vlan {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	key := NewVlanPortKey(port, vid)
	if e := t.find(key.String()); e != nil {
		return e.Vid
	}

	if key.Vid != t.defaultVid {
		return key.Vid
	}

	if _, linkType := fibcapi.ParseDPPortId(uint32(port)); linkType.IsVirtual() {
		return t.defaultVid
	}

	if port := key.Port; (port >= t.minPort) && (port <= t.maxPort) {
		return opennsl.Vlan(port-t.minPort+1) + t.baseVid
	}

	return t.defaultVid
}

func (t *VlanPortTable) Has(vid opennsl.Vlan) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if vid <= t.baseVid {
		return false
	}

	if maxPort := opennsl.Port(vid-t.baseVid) + t.minPort - 1; maxPort > t.maxPort {
		return false
	}

	return true
}

func NewVlanPortTableFromConfig(config *BlockBcastConfig) *VlanPortTable {
	table := NewVlanPortTable()
	table.SetBaseVID(opennsl.Vlan(config.Range.GetBaseVID()))
	table.SetMinPort(opennsl.Port(config.Range.Min))
	table.SetMaxPort(opennsl.Port(config.Range.Max))

	for _, port := range config.Ports {
		e := NewVlanPortEntry(opennsl.Vlan(port.PVid))
		k := NewVlanPortKey(opennsl.Port(port.Port), opennsl.Vlan(port.Vid))
		table.Insert(k, e)
	}

	return table
}
