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
	EthDst  *FieldGroupEthDst
	EthType *FieldGroupEthType
	DstIPv4 *FieldGroupDstIP
	DstIPv6 *FieldGroupDstIP
	IPProto *FieldGroupIPProto
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
func (f *FieldGroup) InstallEntry(key string, entry opennsl.FieldEntry) error {
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
func (f *FieldGroup) UninstallEntry(key string) {
	entry, ok := f.entries[key]
	if !ok {
		log.Warnf("FieldEntry not found. key='%s'", key)
		return
	}

	delete(f.entries, key)

	if err := entry.Remove(f.unit); err != nil {
		log.Warnf("FieldEntry remove error. %s", err)
	}
}

//
// GetEntries returns field group entry.
//
func (f *FieldGroup) GetEntries() ([]opennsl.FieldEntry, error) {
	return f.group.EntryMultiGet(f.unit, -1)
}

//
// FieldGroupEthDst is opennsl field group(ether-dest)
//
type FieldGroupEthDst struct {
	*FieldGroup
}

func NewFieldGroupEthDst(unit int) *FieldGroupEthDst {
	return &FieldGroupEthDst{
		FieldGroup: NewFieldGroup(
			unit, fieldCosDefault, fieldPriEthDst,
			opennsl.FieldQualifyDstMac,
		),
	}
}

func (f *FieldGroupEthDst) Key(dest, mask net.HardwareAddr) string {
	return fmt.Sprintf("%s/%s", dest, mask)
}

func (f *FieldGroupEthDst) AddEntry(dest, mask net.HardwareAddr) error {
	entry, err := f.group.EntryCreate(f.unit)
	if err != nil {
		log.Errorf("EntryCreate error. %s", err)
		return err
	}

	entry.Qualify().DstMAC(f.unit, dest, mask)
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCosQCpuNew(f.cos))
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCopyToCpu())

	return f.InstallEntry(f.Key(dest, mask), entry)
}

//
// DeleteEntry uninstall and remove field group entry.
//
func (f *FieldGroupEthDst) DeleteEntry(dest, mask net.HardwareAddr) {
	f.UninstallEntry(f.Key(dest, mask))
}

//
// GetEntry returns field group entry.
//
func (f *FieldGroupEthDst) GetEntry(entry opennsl.FieldEntry) (net.HardwareAddr, net.HardwareAddr, error) {
	return entry.Qualify().DstMACGet(f.unit)
}

//
// FieldGroupEthType is opennsl field group(ether-type)
//
type FieldGroupEthType struct {
	*FieldGroup
}

//
// NewFieldGroupEthType returns new instance.
//
func NewFieldGroupEthType(unit int) *FieldGroupEthType {
	return &FieldGroupEthType{
		FieldGroup: NewFieldGroup(
			unit, fieldCosDefault, fieldPriEthType,
			opennsl.FieldQualifyEtherType,
		),
	}
}

//
// Key creates new key.
//
func (f *FieldGroupEthType) Key(etherType uint16) string {
	return fmt.Sprintf("%d", etherType)
}

//
// AddEntry creates and installs field group entry.
//
func (f *FieldGroupEthType) AddEntry(etherType uint16) error {
	entry, err := f.group.EntryCreate(f.unit)
	if err != nil {
		log.Errorf("EntryCreate error. %s", err)
		return err
	}

	entry.Qualify().EtherType(f.unit, opennsl.Ethertype(etherType), 0xffff)
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCosQCpuNew(f.cos))
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCopyToCpu())

	return f.InstallEntry(f.Key(etherType), entry)
}

//
// DeleteEntry uninstall and remove field group entry.
//
func (f *FieldGroupEthType) DeleteEntry(etherType uint16) {
	f.UninstallEntry(f.Key(etherType))
}

//
// GetEntry returns field group entry.
//
func (f *FieldGroupEthType) GetEntry(entry opennsl.FieldEntry) (opennsl.Ethertype, error) {
	ethType, _, err := entry.Qualify().EtherTypeGet(f.unit)
	return ethType, err
}

//
// FieldGroupDstIP is field group(dst-ip)
//
type FieldGroupDstIP struct {
	*FieldGroup
	etherType opennsl.Ethertype
}

//
// NewFieldGroupDstIP returns new instalce.
//
func NewFieldGroupDstIPv4(unit int) *FieldGroupDstIP {
	return &FieldGroupDstIP{
		FieldGroup: NewFieldGroup(
			unit, fieldCosDefault, fieldPriDstIPv4,
			opennsl.FieldQualifyEtherType,
			opennsl.FieldQualifyDstIp,
		),
		etherType: opennsl.Ethertype(unix.ETH_P_IP),
	}
}

func NewFieldGroupDstIPv6(unit int) *FieldGroupDstIP {
	return &FieldGroupDstIP{
		FieldGroup: NewFieldGroup(
			unit, fieldCosDefault, fieldPriDstIPv6,
			opennsl.FieldQualifyEtherType,
			opennsl.FieldQualifyDstIp6,
		),
		etherType: opennsl.Ethertype(unix.ETH_P_IPV6),
	}
}

//
// Key creates new key.
//
func (f *FieldGroupDstIP) Key(ipDst *net.IPNet) string {
	return ipDst.String()
}

//
// AddEntry creates and install field group entry.
//
func (f *FieldGroupDstIP) AddEntry(ipDst *net.IPNet) error {
	entry, err := f.group.EntryCreate(f.unit)
	if err != nil {
		log.Errorf("EntryCreate error. %s", err)
		return err
	}

	entry.Qualify().EtherType(f.unit, f.etherType, 0xffff)
	if f.etherType == unix.ETH_P_IP {
		entry.Qualify().DstIp(f.unit, ipDst.IP, ipDst.Mask)
	} else {
		entry.Qualify().DstIp6(f.unit, ipDst.IP, ipDst.Mask)
	}
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCosQCpuNew(f.cos))
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCopyToCpu())

	return f.InstallEntry(f.Key(ipDst), entry)
}

//
// DeleteEntry uninstall and remove field group entry.
//
func (f *FieldGroupDstIP) DeleteEntry(ipDst *net.IPNet) {
	f.UninstallEntry(f.Key(ipDst))
}

//
// GetEntry returns field group entry.
//
func (f *FieldGroupDstIP) GetEntry(entry opennsl.FieldEntry) (opennsl.Ethertype, *net.IPNet, error) {
	ethType, _, err := entry.Qualify().EtherTypeGet(f.unit)
	if err != nil {
		return 0, nil, err
	}

	ip, mask, err := func() (net.IP, net.IPMask, error) {
		if ethType == unix.ETH_P_IPV6 {
			return entry.Qualify().DstIp6Get(f.unit)
		}
		return entry.Qualify().DstIpGet(f.unit)
	}()
	if err != nil {
		return ethType, nil, err
	}

	dstip := &net.IPNet{
		IP:   ip,
		Mask: mask,
	}

	return ethType, dstip, nil
}

//
// FieldGroupIPProto is field group(ip-proto)
//
type FieldGroupIPProto struct {
	*FieldGroup
}

//
// NewFieldGroupIPProto returns new instance.
//
func NewFieldGroupIPProto(unit int) *FieldGroupIPProto {
	return &FieldGroupIPProto{
		FieldGroup: NewFieldGroup(
			unit, fieldCosDefault, fieldPriIPProto,
			opennsl.FieldQualifyEtherType,
			opennsl.FieldQualifyIpProtocol,
		),
	}
}

//
// Key creates new key.
//
func (f *FieldGroupIPProto) Key(etherType uint16, proto uint8) string {
	return fmt.Sprintf("%d_%d", etherType, proto)
}

//
// AddEntry create and install field group entry.
//
func (f *FieldGroupIPProto) AddEntry(etherType uint16, proto uint8) error {
	entry, err := f.group.EntryCreate(f.unit)
	if err != nil {
		log.Errorf("EntryCreate error. %s", err)
		return err
	}

	entry.Qualify().EtherType(f.unit, opennsl.Ethertype(etherType), 0xffff)
	entry.Qualify().IpProtocol(f.unit, proto, 0xff)
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCosQCpuNew(f.cos))
	entry.Action().AddP(f.unit, opennsl.NewFieldActionCopyToCpu())

	return f.InstallEntry(f.Key(etherType, proto), entry)
}

//
// DeleteEntry uninstalls and remove field group entry.
//
func (f *FieldGroupIPProto) DeleteEntry(etherType uint16, proto uint8) {
	f.UninstallEntry(f.Key(etherType, proto))
}

//
// GetEntry returns field group entry.
//
func (f *FieldGroupIPProto) GetEntry(entry opennsl.FieldEntry) (opennsl.Ethertype, uint8, error) {
	ethType, _, err := entry.Qualify().EtherTypeGet(f.unit)
	if err != nil {
		return 0, 0, err
	}

	proto, _, err := entry.Qualify().IpProtocolGet(f.unit)
	if err != nil {
		return ethType, 0, err
	}

	return ethType, proto, nil
}
