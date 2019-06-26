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
	"net"

	"github.com/beluganos/go-opennsl/opennsl"
	"golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
)

const (
	fieldCosDefault = 1
)

//
// FieldEntry is interface of field entry.
//
type FieldEntry interface {
	key() string
	setTo(int, uint32, opennsl.FieldEntry)
	getFrom(int, opennsl.FieldEntry) error
}

const (
	fieldPriHigher = iota + 100
	fieldPriEthDst
	fieldPriEthType
	fieldPriDstIPv4
	fieldPriDstIPv6
	fieldPriIPProto
	fieldPriLower
)

//
// FieldGroups has opennsl field_groups.
//
type FieldGroups struct {
	EthDst  *FieldGroup
	EthType *FieldGroup
	DstIPv4 *FieldGroup
	DstIPv6 *FieldGroup
	IPProto *FieldGroup
}

//
// NewFieldGroups returns new instance.
//
func NewFieldGroups(unit int) *FieldGroups {
	return &FieldGroups{
		EthDst:  NewFieldGroupEthDst(unit),
		EthType: NewFieldGroupEthType(unit),
		DstIPv4: NewFieldGroupDstIPv4(unit),
		DstIPv6: NewFieldGroupDstIPv6(unit),
		IPProto: NewFieldGroupIPProto(unit),
	}
}

//
// FieldGroup has cos, field_group and entries.
//
type FieldGroup struct {
	unit    int
	cos     uint32
	group   opennsl.FieldGroup
	entries map[string]opennsl.FieldEntry
}

//
// NewFieldGroup returns new instance.
//
func NewFieldGroup(unit int, cos uint32, pri int, qs ...opennsl.FieldQualify) *FieldGroup {
	qset := opennsl.NewFieldQSet()
	qset.Add(qs...)
	group, err := opennsl.FieldGroupCreate(unit, qset, pri)
	if err != nil {
		log.Errorf("FieldGroupCreate error. %s", err)
		return nil
	}

	return &FieldGroup{
		unit:    unit,
		cos:     cos,
		group:   group,
		entries: map[string]opennsl.FieldEntry{},
	}
}

//
// InstallEntry installs field group entry.
//
func (f *FieldGroup) installEntry(key string, entry opennsl.FieldEntry) error {
	if _, ok := f.entries[key]; ok {
		return fmt.Errorf("FieldEntry already exist. key='%s'", key)
	}

	if err := entry.Install(f.unit); err != nil {
		return err
	}

	f.entries[key] = entry
	return nil
}

//
// UninstallEntry uninstalls field group entry.
//
func (f *FieldGroup) uninstallEntry(key string) {
	entry, ok := f.entries[key]
	if !ok {
		log.Warnf("FieldEntry not found. key='%s'", key)
		return
	}

	delete(f.entries, key)

	if err := entry.Destroy(f.unit); err != nil {
		log.Warnf("FieldEntry remove error. %s", err)
	}
}

//
// AddEntry installs field entry.
//
func (f *FieldGroup) AddEntry(e FieldEntry) error {
	entry, err := f.group.EntryCreate(f.unit)
	if err != nil {
		log.Errorf("EntryCreate error. %s", err)
		return err
	}

	e.setTo(f.unit, f.cos, entry)
	return f.installEntry(e.key(), entry)
}

//
// DeleteEntry uninstall field entry.
//
func (f *FieldGroup) DeleteEntry(e FieldEntry) {
	f.uninstallEntry(e.key())
}

//
// GetEntry get field entry form H.W.
//
func (f *FieldGroup) GetEntry(e FieldEntry, entry opennsl.FieldEntry) error {
	return e.getFrom(f.unit, entry)
}

//
// GetEntries get all field entry from H.W.
//
func (f *FieldGroup) GetEntries() ([]opennsl.FieldEntry, error) {
	return f.group.EntryMultiGet(f.unit, -1)
}

//
// NewFieldGroupEthDst created new FieldGroup for FieldEntryEthDst.
//
func NewFieldGroupEthDst(unit int) *FieldGroup {
	return NewFieldGroup(
		unit, fieldCosDefault, fieldPriEthDst,
		opennsl.FieldQualifyDstMac,
		opennsl.FieldQualifyInPort,
	)
}

//
// FieldEntryEthDst is field entry (EthDst).
//
type FieldEntryEthDst struct {
	Dest   net.HardwareAddr
	Mask   net.HardwareAddr
	InPort opennsl.Port
}

//
// NewFieldEntryEthDst returns new FieldEntryEthDst.
func NewFieldEntryEthDst(dest, mask net.HardwareAddr, inPort opennsl.Port) *FieldEntryEthDst {
	return &FieldEntryEthDst{
		Dest:   dest,
		Mask:   mask,
		InPort: inPort,
	}
}

func (e *FieldEntryEthDst) key() string {
	return fmt.Sprintf("%s/%s_%d", e.Dest, e.Mask, e.InPort)
}

//
// String returns string.
//
func (e *FieldEntryEthDst) String() string {
	return fmt.Sprintf("%s/%s in_port:%d", e.Dest, e.Mask, e.InPort)
}

func (e *FieldEntryEthDst) setTo(unit int, cos uint32, entry opennsl.FieldEntry) {
	entry.Qualify().DstMAC(unit, e.Dest, e.Mask)
	entry.Qualify().InPort(unit, e.InPort, 0)
	entry.Action().AddP(unit, opennsl.NewFieldActionCosQCpuNew(cos))
	entry.Action().AddP(unit, opennsl.NewFieldActionCopyToCpu())
}

func (e *FieldEntryEthDst) getFrom(unit int, entry opennsl.FieldEntry) error {
	dest, mask, err := entry.Qualify().DstMACGet(unit)
	if err != nil {
		return err
	}

	inPort, _, err := entry.Qualify().InPortGet(unit)
	if err != nil {
		return err
	}

	e.Dest = dest
	e.Mask = mask
	e.InPort = inPort

	return nil
}

//
// NewFieldGroupEthType create new FieldGroup for FieldEntryEthType.
//
func NewFieldGroupEthType(unit int) *FieldGroup {
	return NewFieldGroup(
		unit, fieldCosDefault, fieldPriEthType,
		opennsl.FieldQualifyEtherType,
		opennsl.FieldQualifyInPort,
	)
}

//
// FieldEntryEthType is field entry (EthType).
//
type FieldEntryEthType struct {
	EthType uint16
	InPort  opennsl.Port
}

//
// NewFieldEntryEthType returns new FieldEntryEthType.
//
func NewFieldEntryEthType(ethType uint16, inPort opennsl.Port) *FieldEntryEthType {
	return &FieldEntryEthType{
		EthType: ethType,
		InPort:  inPort,
	}
}

func (e *FieldEntryEthType) key() string {
	return fmt.Sprintf("%d_%d", e.EthType, e.InPort)
}

func (e *FieldEntryEthType) String() string {
	return fmt.Sprintf("%04x in_port:%d", e.EthType, e.InPort)
}

func (e *FieldEntryEthType) setTo(unit int, cos uint32, entry opennsl.FieldEntry) {
	entry.Qualify().InPort(unit, e.InPort, 0)
	entry.Qualify().EtherType(unit, opennsl.Ethertype(e.EthType), 0xffff)
	entry.Action().AddP(unit, opennsl.NewFieldActionCosQCpuNew(cos))
	entry.Action().AddP(unit, opennsl.NewFieldActionCopyToCpu())
}

func (e *FieldEntryEthType) getFrom(unit int, entry opennsl.FieldEntry) error {
	ethType, _, err := entry.Qualify().EtherTypeGet(unit)
	if err != nil {
		return err
	}

	inPort, _, err := entry.Qualify().InPortGet(unit)
	if err != nil {
		return err
	}

	e.EthType = uint16(ethType)
	e.InPort = inPort

	return nil
}

//
// NewFieldGroupDstIPv4 returns new FieldGroup for FieldEntryDstIP(v4)
//
func NewFieldGroupDstIPv4(unit int) *FieldGroup {
	return NewFieldGroup(
		unit, fieldCosDefault, fieldPriDstIPv4,
		opennsl.FieldQualifyEtherType,
		opennsl.FieldQualifyDstIp,
		opennsl.FieldQualifyInPort,
	)
}

//
// NewFieldGroupDstIPv6 returns new FieldGroup for FieldEntryDstIP(v6)
//
func NewFieldGroupDstIPv6(unit int) *FieldGroup {
	return NewFieldGroup(
		unit, fieldCosDefault, fieldPriDstIPv6,
		opennsl.FieldQualifyEtherType,
		opennsl.FieldQualifyDstIp6,
		opennsl.FieldQualifyInPort,
	)
}

//
// FieldEntryDstIP is field entry (DstIPv4 or DstIPv6).
//
type FieldEntryDstIP struct {
	EthType opennsl.Ethertype
	Dest    *net.IPNet
	InPort  opennsl.Port
}

//
// NewFieldEntryDstIP returns new FieldEntryDstIP(IPv4 or IPv6)
//
func NewFieldEntryDstIP(dest *net.IPNet, inPort opennsl.Port, ethType uint16) *FieldEntryDstIP {
	return &FieldEntryDstIP{
		EthType: opennsl.Ethertype(ethType),
		Dest:    dest,
		InPort:  inPort,
	}
}

//
// NewFieldEntryDstIPv4 returns new FieldEntryDstIP(IPv4)
//
func NewFieldEntryDstIPv4(dest *net.IPNet, inPort opennsl.Port) *FieldEntryDstIP {
	return NewFieldEntryDstIP(dest, inPort, unix.ETH_P_IP)
}

//
// NewFieldEntryDstIPv6 returns new FieldEntryDstIP(IPv6)
//
func NewFieldEntryDstIPv6(dest *net.IPNet, inPort opennsl.Port) *FieldEntryDstIP {
	return NewFieldEntryDstIP(dest, inPort, unix.ETH_P_IPV6)
}

func (e *FieldEntryDstIP) key() string {
	return fmt.Sprintf("%s_%d", e.Dest, e.InPort)
}

func (e *FieldEntryDstIP) String() string {
	return fmt.Sprintf("%s in_port=%d", e.Dest, e.InPort)
}

func (e *FieldEntryDstIP) setTo(unit int, cos uint32, entry opennsl.FieldEntry) {
	entry.Qualify().InPort(unit, e.InPort, 0)
	entry.Qualify().EtherType(unit, e.EthType, 0xffff)
	if e.EthType == unix.ETH_P_IP {
		entry.Qualify().DstIp(unit, e.Dest.IP, e.Dest.Mask)
	} else {
		entry.Qualify().DstIp6(unit, e.Dest.IP, e.Dest.Mask)
	}
	entry.Action().AddP(unit, opennsl.NewFieldActionCosQCpuNew(cos))
	entry.Action().AddP(unit, opennsl.NewFieldActionCopyToCpu())
}

func (e *FieldEntryDstIP) getFrom(unit int, entry opennsl.FieldEntry) error {
	ethType, _, err := entry.Qualify().EtherTypeGet(unit)
	if err != nil {
		return err
	}

	ip, mask, err := func() (net.IP, net.IPMask, error) {
		if ethType == unix.ETH_P_IPV6 {
			return entry.Qualify().DstIp6Get(unit)
		}
		return entry.Qualify().DstIpGet(unit)
	}()
	if err != nil {
		return err
	}

	inPort, _, err := entry.Qualify().InPortGet(unit)
	if err != nil {
		return err
	}

	e.EthType = ethType
	e.Dest = &net.IPNet{IP: ip, Mask: mask}
	e.InPort = inPort

	return nil
}

//
// NewFieldGroupIPProto created new FieldGroup for FieldEntryIPProto.
//
func NewFieldGroupIPProto(unit int) *FieldGroup {
	return NewFieldGroup(
		unit, fieldCosDefault, fieldPriIPProto,
		opennsl.FieldQualifyEtherType,
		opennsl.FieldQualifyIpProtocol,
		opennsl.FieldQualifyInPort,
	)
}

//
// FieldEntryIPProto is field entry (IP proto).
//
type FieldEntryIPProto struct {
	EthType uint16
	IPProto uint8
	InPort  opennsl.Port
}

//
// NewFieldEntryIPProto returns new FieldEntryIPProto.
//
func NewFieldEntryIPProto(ipProto uint8, ethType uint16, inPort opennsl.Port) *FieldEntryIPProto {
	return &FieldEntryIPProto{
		EthType: ethType,
		IPProto: ipProto,
		InPort:  inPort,
	}
}

func (e *FieldEntryIPProto) key() string {
	return fmt.Sprintf("%d_%d_%d", e.EthType, e.IPProto, e.InPort)
}

func (e *FieldEntryIPProto) String() string {
	return fmt.Sprintf("%d eth_type:%04x in_port:%d", e.EthType, e.IPProto, e.InPort)
}

func (e *FieldEntryIPProto) setTo(unit int, cos uint32, entry opennsl.FieldEntry) {
	entry.Qualify().InPort(unit, e.InPort, 0)
	entry.Qualify().EtherType(unit, opennsl.Ethertype(e.EthType), 0xffff)
	entry.Qualify().IpProtocol(unit, e.IPProto, 0xff)
	entry.Action().AddP(unit, opennsl.NewFieldActionCosQCpuNew(cos))
	entry.Action().AddP(unit, opennsl.NewFieldActionCopyToCpu())
}

func (e *FieldEntryIPProto) getFrom(unit int, entry opennsl.FieldEntry) error {
	ethType, _, err := entry.Qualify().EtherTypeGet(unit)
	if err != nil {
		return err
	}

	proto, _, err := entry.Qualify().IpProtocolGet(unit)
	if err != nil {
		return err
	}

	inPort, _, err := entry.Qualify().InPortGet(unit)
	if err != nil {
		return err
	}

	e.EthType = uint16(ethType)
	e.IPProto = proto
	e.InPort = inPort

	return nil
}
