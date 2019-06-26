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

package gonslib

import (
	"fmt"

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

func SetPortVlanTranslation(unit int, port opennsl.Port, enable int) error {
	log.Infof("Set vlan translation. port=%d, enable=%d", port, enable)

	ctrls := []opennsl.VlanControlPort{
		opennsl.VlanTranslateIngressEnable,
		opennsl.VlanTranslateIngressMissDrop,
		opennsl.VlanTranslateEgressEnable,
		opennsl.VlanTranslateEgressMissDrop,
	}

	// Enable VLAN translations for both ingress and egress.
	for _, ctrl := range ctrls {
		if err := ctrl.Set(unit, port, enable); err != nil {
			log.Errorf("VlanTranslate%s  error. port=%d enable=%d %s", ctrl, port, enable, err)
			return err
		}

		log.Debugf("VlanTranslate%s ok. port=%d enable=%d", ctrl, port, enable)
	}

	if enable == opennsl.TRUE {
		// Set up port's double tagging mode.
		if err := port.DtagModeSet(unit, opennsl.PORT_DTAG_MODE_INTERNAL); err != nil {
			log.Errorf("PortDtagModeSet port=%d mode=%s error. %s", port, opennsl.PORT_DTAG_MODE_INTERNAL, err)
			return err
		}

		log.Debugf("DtagModeSet ok port=%d, mode=%s", port, opennsl.PORT_DTAG_MODE_INTERNAL)
	}

	return nil
}

func SetNativeVlan(unit int, vid opennsl.Vlan, pbmp *opennsl.PBmp, strictlyUntagged bool) error {
	log.Debugf("Set native vlan. vid=%d pbmp=%s strict=%t", vid, pbmp, strictlyUntagged)

	defaultVid, err := opennsl.VlanDefaultGet(unit)
	if err != nil {
		log.Errorf("SetNativeVlan VlanDefaultGet error. unit=%d %s", unit, err)
		return err
	}

	if defaultVid == vid {
		log.Debugf("SetNativeVlan skip. vid=%d", vid)
		return nil
	}

	return pbmp.Each(func(port opennsl.Port) error {
		if err := port.UntaggedVlanSet(unit, vid); err != nil {
			log.Errorf("SetNativeVlan PortUntaggedVlanSet error. port=%d vid=%d %s", port, vid, err)
			// Don't exit.  Keep setting the other ports.
		}

		if strictlyUntagged {
			// If strictly untagged option is set, we want to enable VLAN
			// translation & set up miss-drop flags on this port (although
			// we won't need any actual old VID to new VID mapping).  This
			// is to prevent external frames with VID that happened to match
			// the internal VID to be accepted.  In other words, if a port
			// is strictly untagged, only untagged frame will be allowed.
			log.Debugf("SetNativeVlan with strictry untagged enabld. port=%d", port)
			return SetPortVlanTranslation(unit, port, opennsl.TRUE)
		}

		return nil
	})
}

func ClearNativeVlan(unit int, vid opennsl.Vlan, pbmp *opennsl.PBmp, strictlyUntagged bool) {
	// Get the switch's default vid.
	defaultVid, err := opennsl.VlanDefaultGet(unit)
	if err != nil {
		log.Errorf("ClearNativeVlan DefaultGet error. unit=%d %s", unit, err)
		return
	}

	if defaultVid == vid {
		log.Debugf("ClearNativeVlan skip. vid=%d", vid)
		return
	}

	pbmp.Each(func(port opennsl.Port) error {
		if err := port.UntaggedVlanSet(unit, defaultVid); err != nil {
			log.Errorf("ClearNativeVlan PortUntaggedVlanSet error. port=%d vid=%d %s", port, defaultVid, err)
		}

		if strictlyUntagged {
			// Also clear translation settings if this port
			// was strictly untagged.
			log.Debugf("ClearNativeVlan with strictry untagged disabld. port=%d", port)
			SetPortVlanTranslation(unit, port, opennsl.FALSE)
		}

		return nil
	})
}

func CreateVlan(unit int, vid opennsl.Vlan) error {
	log.Debugf("Create vlan. vid=%d", vid)

	if _, err := vid.Create(unit); err != nil {
		log.Errorf("Create vlan error. vid=%d %s", vid, err)
		return err
	}

	return nil
}

func DestroyVlan(unit int, vid opennsl.Vlan) {
	log.Debugf("Destroy vlan. vid=%d", vid)

	if err := vid.Destroy(unit); err != nil {
		log.Errorf("Destroy vlan error. vid=%d %s", vid, err)
	}
}

func DeleteL2Addrs(unit int, vid opennsl.Vlan) {
	if err := opennsl.L2AddrDeleteByVID(unit, vid, opennsl.L2_DELETE_NO_CALLBACKS); err != nil {
		log.Errorf("Delete l2addrs by vlan %d", vid)
	}
}

func DestroyVlanIfEmpty(unit int, vid opennsl.Vlan) {
	pbmp, _, err := vid.PortGet(unit)
	if err != nil {
		log.Errorf("Delete vlan PortGet error. %d %s", vid, err)
		return
	}

	if pbmp.IsNull() {
		DeleteL2Addrs(unit, vid)
		DestroyVlan(unit, vid)
	}
}

func AddPortsToVlan(unit int, allBmp *opennsl.PBmp, untagBmp *opennsl.PBmp, vid opennsl.Vlan, strictlyUntagged bool) error {

	if untagBmp.IsNotNull() {
		// Update default VLAN ID of the ports if untagged.
		if err := SetNativeVlan(unit, vid, untagBmp, strictlyUntagged); err != nil {
			log.Errorf("setNativeVlan error. vid=%d untags=%s %s", vid, untagBmp, err)
			return err
		}
	}

	if allBmp.IsNotNull() {
		// Finally, add ports to VLAN.
		if err := vid.PortAdd(unit, allBmp, untagBmp); err != nil {
			log.Errorf("VlanPortAdd error. vid=%d ports=%s untags=%s", vid, allBmp, untagBmp)
			return err
		}
	}

	return nil
}

func DelPortsFromVlan(unit int, allBmp *opennsl.PBmp, untagBmp *opennsl.PBmp, vid opennsl.Vlan, strictlyUntagged bool) {

	if allBmp.IsNotNull() {
		// Remove ports from VLAN.
		if _, err := vid.PortRemove(unit, allBmp); err != nil {
			log.Errorf("VlanPortRemove error. vid=%d ports=%s", vid, allBmp)
		}
	}

	if untagBmp.IsNotNull() {
		// Update default VLAN ID of the ports if untagged.
		ClearNativeVlan(unit, vid, untagBmp, strictlyUntagged)
	}

	DestroyVlanIfEmpty(unit, vid)
}

func adjustVlan(unit int, vid opennsl.Vlan) opennsl.Vlan {
	if vid == opennsl.VLAN_ID_NONE {
		vid = opennsl.VlanDefaultMustGet(unit)
	}
	return vid
}

type L3Vlan struct {
	Vlan opennsl.Vlan
	Vid  opennsl.Vlan
	Pbmp *opennsl.PBmp
}

func NewL3Vlan(unit int, vid opennsl.Vlan) *L3Vlan {
	return &L3Vlan{
		Vlan: adjustVlan(unit, 0),
		Vid:  adjustVlan(unit, vid),
		Pbmp: opennsl.NewPBmp(),
	}
}

func (b *L3Vlan) String() string {
	return fmt.Sprintf("vlan:%d vid:%d ports:%s", b.Vlan, b.Vid, b.Pbmp)
}

func (b *L3Vlan) Create(unit int) error {
	if err := CreateVlan(unit, b.Vlan); err != nil {
		return err
	}

	ubmp := func() *opennsl.PBmp {
		if b.Vid == opennsl.VlanDefaultMustGet(unit) {
			return b.Pbmp
		}
		return opennsl.NewPBmp()
	}()

	return b.Vlan.PortAdd(unit, b.Pbmp, ubmp)
}

func (b *L3Vlan) Delete(unit int) error {
	_, err := b.Vlan.PortRemove(unit, b.Pbmp)
	DestroyVlan(unit, b.Vlan)

	return err
}

type BrVlan struct {
	Vid              opennsl.Vlan
	Pbmp             *opennsl.PBmp
	UntagBmp         *opennsl.PBmp
	StrictlyUntagged bool
}

func NewBrVlan(unit int, vid opennsl.Vlan) *BrVlan {
	return &BrVlan{
		Vid:      adjustVlan(unit, vid),
		Pbmp:     opennsl.NewPBmp(),
		UntagBmp: opennsl.NewPBmp(),
	}
}

func (b *BrVlan) String() string {
	return fmt.Sprintf("vid:%d strict:%t %s/%s", b.Vid, b.StrictlyUntagged, b.Pbmp, b.UntagBmp)
}

func (b *BrVlan) Create(unit int) error {
	if err := CreateVlan(unit, b.Vid); err != nil {
		return err
	}

	return AddPortsToVlan(unit, b.Pbmp, b.UntagBmp, b.Vid, b.StrictlyUntagged)
}

func (b *BrVlan) Delete(unit int) {
	DelPortsFromVlan(unit, b.Pbmp, b.UntagBmp, b.Vid, b.StrictlyUntagged)
}
