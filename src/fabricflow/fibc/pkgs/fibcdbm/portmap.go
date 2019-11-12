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
	"fmt"
	"sync"
)

//
// PortKey is key of PortMap.
//
type PortKey struct {
	ReID   string
	Ifname string
}

//
// NewPortKey returns new PortKey.
//
func NewPortKey(reID, ifname string) *PortKey {
	return &PortKey{
		ReID:   reID,
		Ifname: ifname,
	}
}

//
// String is stringer.
//
func (e *PortKey) String() string {
	return fmt.Sprintf("reid:'%s' ifname:'%s'", e.ReID, e.Ifname)
}

//
// Key returns key of PortMap.
//
func (e *PortKey) toKey() string {
	return fmt.Sprintf("%s@%s", e.ReID, e.Ifname)
}

//
// Equals compares keys.
//
func (e *PortKey) Equals(other *PortKey) bool {
	if e == nil || other == nil {
		return false
	}
	return (e.ReID == other.ReID) && (e.Ifname == other.Ifname)
}

//
// Clone returns new PortKey.
//
func (e *PortKey) Clone() *PortKey {
	if e == nil {
		return nil
	}

	dst := &PortKey{}
	*dst = *e
	return dst
}

//
// PortValue is port value of port entry
//
type PortValue struct {
	DpID   uint64
	ReID   string
	PortID uint32
	Enter  bool
}

//
// NewPortValue returns new PortValue.
//
func NewPortValue(dpID uint64, portID uint32, enter bool) *PortValue {
	return &PortValue{
		DpID:   dpID,
		PortID: portID,
		Enter:  enter,
	}

}

//
// NewPortValueR returns new PortValue.
//
func NewPortValueR(reID string, portID uint32, enter bool) *PortValue {
	return &PortValue{
		ReID:   reID,
		PortID: portID,
		Enter:  enter,
	}
}

//
// IsValid returns valid or not.
//
func (v *PortValue) IsValid() bool {
	return (v != nil)
}

//
// IsInvalid returns valid or not.
//
func (v *PortValue) IsInvalid() bool {
	return (v == nil)
}

//
//
//
func (v *PortValue) isVMPort() bool {
	return len(v.ReID) != 0
}

//
// String is stringer.
//
func (v *PortValue) String() string {
	if v.IsInvalid() {
		return fmt.Sprintf("invalid port")
	}

	if v.isVMPort() {
		return fmt.Sprintf("reid:'%s' port:%d enter:%t", v.ReID, v.PortID, v.Enter)
	}

	return fmt.Sprintf("dpid:%d port:%d enter:%t", v.DpID, v.PortID, v.Enter)
}

//
// Equals compares values.
//
func (v *PortValue) Equals(other *PortValue) bool {
	if v.IsInvalid() || other.IsInvalid() {
		return false
	}

	if v.isVMPort() {
		return (v.ReID == other.ReID) && (v.PortID == other.PortID)
	}

	return (v.DpID == other.DpID) && (v.PortID == other.PortID)
}

//
// Clone returns new PortValue.
//
func (v *PortValue) Clone() *PortValue {
	if v.IsInvalid() {
		return nil
	}

	dst := &PortValue{}
	*dst = *v
	return dst
}

//
// Reset clear values.
//
func (v *PortValue) Reset() {
	if v.IsValid() {
		v.DpID = 0
		v.ReID = ""
		v.PortID = 0
		v.Enter = false
	}
}

//
// IsAssociated returns associated or not.
//
func (v *PortValue) IsAssociated() bool {
	if v.IsInvalid() {
		return false
	}

	return v.Enter
}

//
// Update update PortId.
//
func (v *PortValue) Update(dpID uint64, portID uint32) (upd bool) {
	if v.DpID != dpID {
		upd = true
		v.DpID = dpID
	}

	if v.PortID != portID {
		upd = true
		v.PortID = portID
	}

	return
}

//
// UpdateR update PortId.
//
func (v *PortValue) UpdateR(reID string, portID uint32) (upd bool) {
	if v.ReID != reID {
		upd = true
		v.ReID = reID
	}

	if v.PortID != portID {
		upd = true
		v.PortID = portID
	}

	return
}

//
// UpdatePort update PortId.
//
func (v *PortValue) UpdatePort(portID uint32) (upd bool) {
	if v.PortID != portID {
		upd = true
		v.PortID = portID
	}

	return
}

//
// UpdateEnter set enter status.
//
func (v *PortValue) UpdateEnter(enter bool) (upd bool) {
	if v.Enter != enter {
		upd = true
		v.Enter = enter
	}

	return
}

//
// PortEntry is entry of PortMap.
//
type PortEntry struct {
	Key       *PortKey
	ParentKey *PortKey
	MasterKey *PortKey

	VMPort *PortValue
	DPPort *PortValue
	VSPort *PortValue

	version uint64
}

//
// NewPortEntry returns new PortEntry.
//
func NewPortEntry(key *PortKey) *PortEntry {
	return &PortEntry{
		Key: key,
	}
}

//
// Clone returns new PortEntry.
//
func (e *PortEntry) Clone() *PortEntry {
	if e == nil {
		return nil
	}

	return &PortEntry{
		Key:       e.Key.Clone(),
		ParentKey: e.ParentKey.Clone(),
		MasterKey: e.MasterKey.Clone(),

		VMPort: e.VMPort.Clone(),
		DPPort: e.DPPort.Clone(),
		VSPort: e.VSPort.Clone(),

		version: e.version,
	}
}

//
// IsAssociated judges if port is associated.
//
func (e *PortEntry) IsAssociated() bool {
	b := e.VMPort.IsAssociated()

	if e.DPPort.IsValid() {
		b = b && e.DPPort.IsAssociated()
	}

	if e.VSPort.IsValid() {
		b = b && e.VSPort.IsAssociated()
	}

	return b
}

//
// PortMap is table of ports(vm,vs,dp).
//
type PortMap struct {
	mutex   sync.RWMutex
	entries map[string]*PortEntry
	vmKey   *PortMapVMKey
	dpKey   *PortMapDPKey
	vsKey   *PortMapDPKey
	version uint64
}

//
// NewPortMap returns new PortMap.
//
func NewPortMap() *PortMap {
	return &PortMap{
		entries: map[string]*PortEntry{},
		vmKey:   NewPortMapVMKey(),
		dpKey:   NewPortMapDPKey(),
		vsKey:   NewPortMapDPKey(),
	}
}

func (m *PortMap) find(key *PortKey) *PortEntry {
	if e, ok := m.entries[key.toKey()]; ok {
		return e
	}

	return nil
}

func (m *PortMap) add(e *PortEntry) {
	m.entries[e.Key.toKey()] = e
}

func (m *PortMap) delete(key *PortKey) {
	delete(m.entries, key.toKey())
}

//
// VerUp updates version.
//
func (m *PortMap) VerUp() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.version++
}

//
// Register adds entry if not exist.
//
func (m *PortMap) Register(e *PortEntry) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if old := m.find(e.Key); old != nil {
		return false
	}

	e.version = m.version
	m.add(e)
	m.dpKey.AddPort(e.DPPort, e.Key)

	return true
}

//
// Unregister deletes entry.
//
func (m *PortMap) Unregister(key *PortKey) *PortEntry {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if e := m.find(key); e != nil {
		m.delete(key)
		m.vmKey.DeletePort(e.VMPort)
		m.dpKey.DeletePort(e.DPPort)
		m.vsKey.DeletePort(e.VSPort)
		return e
	}

	return nil
}

//
// RegisterVMKey registers vm key.
//
func (m *PortMap) RegisterVMKey(reID string, portID uint32, key *PortKey) {
	m.vmKey.Add(reID, portID, key)
}

//
// UnregisterVMKey deletes vm key.
//
func (m *PortMap) UnregisterVMKey(reID string, portID uint32) {
	m.vmKey.Delete(reID, portID)
}

//
// RegisterVSKey registers vs key.
//
func (m *PortMap) RegisterVSKey(vsID uint64, portID uint32, key *PortKey) {
	m.vsKey.Add(vsID, portID, key)
}

//
// UnregisterVSKey deletes vs key.
//
func (m *PortMap) UnregisterVSKey(vsID uint64, portID uint32) {
	m.vsKey.Delete(vsID, portID)
}

//
// SelectOrRegister select entry if eist, or add entry if not exist.
//
func (m *PortMap) SelectOrRegister(e *PortEntry, f func(*PortEntry) bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old := m.find(e.Key)
	if old == nil {
		m.add(e)
		m.dpKey.AddPort(e.DPPort, e.Key)
		old = e
	}

	if ok := f(old); ok {
		old.version = m.version
	}
}

//
// Select returns port entry if exist.
//
func (m *PortMap) Select(key *PortKey, f func(*PortEntry)) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if e := m.find(key); e != nil {
		f(e)
		return true
	}

	return false
}

//
// Update returns port entry if exist.
//
func (m *PortMap) Update(key *PortKey, f func(*PortEntry)) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if e := m.find(key); e != nil {
		f(e)
		return true
	}

	return false
}

func (m *PortMap) findByVM(reID string, portID uint32) *PortEntry {
	var (
		key *PortKey
		e   *PortEntry
	)

	if key = m.vmKey.findby(reID, portID); key == nil {
		return nil
	}

FOR_LOOP:
	for {
		if e = m.find(key); e == nil {
			return nil
		}

		if key = e.ParentKey; key == nil {
			break FOR_LOOP
		}
	}

	return e
}

//
// SelectByVM returns port entry if exist.
//
func (m *PortMap) SelectByVM(reID string, portID uint32, f func(*PortEntry)) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if e := m.findByVM(reID, portID); e != nil {
		f(e)
		return true
	}

	return false
}

//
// UpdateByVM returns port entry if exist.
//
func (m *PortMap) UpdateByVM(reID string, portID uint32, f func(*PortEntry)) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if e := m.findByVM(reID, portID); e != nil {
		f(e)
		return true
	}

	return false
}

func (m *PortMap) findByDP(dpID uint64, portID uint32) *PortEntry {
	key := m.dpKey.findby(dpID, portID)
	if key == nil {
		return nil
	}

	return m.find(key)
}

//
// SelectByDP returns port entry if exist.
//
func (m *PortMap) SelectByDP(dpID uint64, portID uint32, f func(*PortEntry)) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if e := m.findByDP(dpID, portID); e != nil {
		f(e)
		return true
	}

	return false
}

//
// UpdateByDP returns port entry if exist.
//
func (m *PortMap) UpdateByDP(dpID uint64, portID uint32, f func(*PortEntry)) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if e := m.findByDP(dpID, portID); e != nil {
		f(e)
		return true
	}

	return false
}

func (m *PortMap) findByVS(vsID uint64, portID uint32) *PortEntry {
	key := m.vsKey.findby(vsID, portID)
	if key == nil {
		return nil
	}

	return m.find(key)
}

//
// SelectByVS returns port entry if exist.
//
func (m *PortMap) SelectByVS(vsID uint64, portID uint32, f func(*PortEntry)) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if e := m.findByVS(vsID, portID); e != nil {
		f(e)
		return true
	}

	return false
}

//
// UpdateByVS returns port entry if exist.
//
func (m *PortMap) UpdateByVS(vsID uint64, portID uint32, f func(*PortEntry)) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if e := m.findByVS(vsID, portID); e != nil {
		f(e)
		return true
	}

	return false
}

//
// ListByVM returns port entry list.
//
func (m *PortMap) ListByVM(reID string, f func(*PortEntry)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, e := range m.entries {
		if vmPort := e.VMPort; vmPort.ReID == reID {
			f(e)
		}
	}
}

//
// ListByDP returns port entry list.
//
func (m *PortMap) ListByDP(dpID uint64, f func(*PortEntry)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, e := range m.entries {
		if dpPort := e.DPPort; (dpPort != nil) && (dpPort.DpID == dpID) {
			f(e)
		}
	}
}

//
// ListByVS returns port entry list.
//
func (m *PortMap) ListByVS(vsID uint64, f func(*PortEntry)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, e := range m.entries {
		if vsPort := e.VSPort; (vsPort != nil) && (vsPort.DpID == vsID) {
			f(e)
		}
	}
}

func (m *PortMap) listByParent(parentKey *PortKey, f func(*PortEntry) bool) {
	for _, e := range m.entries {
		if ok := e.ParentKey.Equals(parentKey); ok {
			if recursive := f(e); recursive {
				m.listByParent(e.Key, f)
			}
		}
	}
}

//
// ListByParent returns port entry list.
//
func (m *PortMap) ListByParent(parentKey *PortKey, f func(*PortEntry) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	m.listByParent(parentKey, f)
}

func (m *PortMap) isAssociated(key *PortKey) (bool, error) {
	// follow parent link until physical device.
	if e := m.find(key); e != nil {
		if pkey := e.ParentKey; pkey != nil {
			// e is not physical device.
			return m.isAssociated(pkey)
		}

		// e is physical device.
		return e.IsAssociated(), nil
	}

	// can not folow parent link.
	return false, fmt.Errorf("entry not found. %s", key)
}

//
// IsAssociated judges if port is associated.
//
func (m *PortMap) IsAssociated(key *PortKey) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.isAssociated(key)
}

//
// GC enumerate entries which version is not equal to PortMap.
//
func (m *PortMap) GC(f func(*PortEntry) bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delEntries := []*PortEntry{}
	for _, e := range m.entries {
		if e.version != m.version {
			if ok := f(e); ok {
				e.version = m.version
			} else {
				delEntries = append(delEntries, e)
			}
		}
	}

	for _, e := range delEntries {
		m.delete(e.Key)
		m.dpKey.deletePort(e.DPPort)
	}
}

//
// Range enumrate all entries.
//
func (m *PortMap) Range(f func(e *PortEntry)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, e := range m.entries {
		f(e)
	}
}

//
// PortMapDPKey is dp key table
//
type PortMapDPKey struct {
	mutex   sync.Mutex
	entries map[string]*PortKey
}

//
// NewPortMapDPKey returns nnew PortMapDPKey
//
func NewPortMapDPKey() *PortMapDPKey {
	return &PortMapDPKey{
		entries: map[string]*PortKey{},
	}
}

//
// NewPortMapDPKeyKey returns nnew PortMapDPKey
//
func NewPortMapDPKeyKey(dpID uint64, portID uint32) string {
	return fmt.Sprintf("%d_%d", dpID, portID)
}

func (m *PortMapDPKey) find(key string) *PortKey {
	if k, ok := m.entries[key]; ok {
		return k
	}

	return nil
}

func (m *PortMapDPKey) findby(dpID uint64, portID uint32) *PortKey {
	key := NewPortMapDPKeyKey(dpID, portID)
	return m.find(key)
}

func (m *PortMapDPKey) add(dpID uint64, portID uint32, pkey *PortKey) {
	key := NewPortMapDPKeyKey(dpID, portID)
	m.entries[key] = pkey
}

//
// Add adds port.
//
func (m *PortMapDPKey) Add(dpID uint64, portID uint32, pkey *PortKey) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.add(dpID, portID, pkey)
}

func (m *PortMapDPKey) addPort(pval *PortValue, pkey *PortKey) {
	if pval.IsValid() {
		m.add(pval.DpID, pval.PortID, pkey)
	}
}

//
// AddPort adds port.
//
func (m *PortMapDPKey) AddPort(pval *PortValue, pkey *PortKey) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.addPort(pval, pkey)
}

func (m *PortMapDPKey) delete(dpID uint64, portID uint32) {
	key := NewPortMapDPKeyKey(dpID, portID)
	if v := m.find(key); v != nil {
		delete(m.entries, key)
	}
}

//
// Delete deletes port.
//
func (m *PortMapDPKey) Delete(dpID uint64, portID uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.delete(dpID, portID)
}

func (m *PortMapDPKey) deletePort(pval *PortValue) {
	if pval.IsValid() {
		m.delete(pval.DpID, pval.PortID)
	}
}

//
// DeletePort deletes port.
//
func (m *PortMapDPKey) DeletePort(pval *PortValue) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.deletePort(pval)
}

//
// PortMapVMKey is mv key table
//
type PortMapVMKey struct {
	mutex   sync.Mutex
	entries map[string]*PortKey
}

//
// NewPortMapVMKey returns new PortMapVMKey
//
func NewPortMapVMKey() *PortMapVMKey {
	return &PortMapVMKey{
		entries: map[string]*PortKey{},
	}
}

//
// NewPortMapVMKeyKey returns key of PortMapVMKey
//
func NewPortMapVMKeyKey(reID string, portID uint32) string {
	return fmt.Sprintf("%s_%d", reID, portID)
}

func (p *PortMapVMKey) find(key string) *PortKey {
	if e, ok := p.entries[key]; ok {
		return e
	}

	return nil
}

func (p *PortMapVMKey) findby(reID string, portID uint32) *PortKey {
	key := NewPortMapVMKeyKey(reID, portID)
	return p.find(key)
}

func (p *PortMapVMKey) add(reID string, portID uint32, pkey *PortKey) {
	key := NewPortMapVMKeyKey(reID, portID)
	p.entries[key] = pkey
}

//
// Add adds port
//
func (p *PortMapVMKey) Add(reID string, portID uint32, pkey *PortKey) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.add(reID, portID, pkey)
}

func (p *PortMapVMKey) addPort(port *PortValue, pkey *PortKey) {
	p.add(port.ReID, port.PortID, pkey)
}

//
// AddPort adds port.
//
func (p *PortMapVMKey) AddPort(port *PortValue, pkey *PortKey) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.addPort(port, pkey)
}

func (p *PortMapVMKey) delete(reID string, portID uint32) {
	key := NewPortMapVMKeyKey(reID, portID)
	if e := p.find(key); e != nil {
		delete(p.entries, key)
	}
}

//
// Delete deletes port.
//
func (p *PortMapVMKey) Delete(reID string, portID uint32) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.delete(reID, portID)
}

func (p *PortMapVMKey) deletePort(port *PortValue) {
	p.delete(port.ReID, port.PortID)
}

//
// DeletePort deletes port.
//
func (p *PortMapVMKey) DeletePort(port *PortValue) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.deletePort(port)
}
