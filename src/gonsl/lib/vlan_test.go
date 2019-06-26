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
	"testing"

	"github.com/beluganos/go-opennsl/opennsl"
)

func TestVlanPortTable_ConvVlan_no_change(t *testing.T) {
	var (
		vid opennsl.Vlan
	)

	defaultVid := opennsl.Vlan(2)
	tbl := NewVlanPortTable()
	tbl.SetDefaultVID(defaultVid)

	vid = tbl.ConvVID(opennsl.Port(50), opennsl.Vlan(2))
	if vid != defaultVid {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(50), opennsl.Vlan(1))
	if vid != opennsl.Vlan(1) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(50), opennsl.Vlan(3))
	if vid != opennsl.Vlan(3) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}
}

func TestVlanPortTable_ConvVlan_minmax(t *testing.T) {
	var (
		vid opennsl.Vlan
	)

	tbl := NewVlanPortTable()
	tbl.SetMinPort(opennsl.Port(50))
	tbl.SetMaxPort(opennsl.Port(55))

	vid = tbl.ConvVID(opennsl.Port(49), opennsl.Vlan(1))
	if vid != opennsl.Vlan(1) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(50), opennsl.Vlan(1))
	if vid != opennsl.Vlan(2) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(55), opennsl.Vlan(1))
	if vid != opennsl.Vlan(7) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(56), opennsl.Vlan(1))
	if vid != opennsl.Vlan(1) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(49), opennsl.Vlan(2))
	if vid != opennsl.Vlan(2) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(50), opennsl.Vlan(2))
	if vid != opennsl.Vlan(2) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(55), opennsl.Vlan(2))
	if vid != opennsl.Vlan(2) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}

	vid = tbl.ConvVID(opennsl.Port(56), opennsl.Vlan(2))
	if vid != opennsl.Vlan(2) {
		t.Errorf("ConvVID unmatch. vid=%d", vid)
	}
}

func TestVlanPortTableHasDefaultVIDBase(t *testing.T) {
	tbl := NewVlanPortTable()
	tbl.SetMinPort(opennsl.Port(50))
	tbl.SetMaxPort(opennsl.Port(55))
	tbl.SetBaseVID(opennsl.Vlan(1))

	// port: 49 -> x
	// port: 50 -> 2
	// port: 51 -> 3
	// ...
	// port: 55 -> 7
	// port: 56 -> x

	if b := tbl.Has(1); b {
		t.Errorf("Has 1 unmatch. %t", b)
	}

	if b := tbl.Has(2); !b {
		t.Errorf("Has 2 unmatch. %t", b)
	}

	if b := tbl.Has(7); !b {
		t.Errorf("Has 7 unmatch. %t", b)
	}

	if b := tbl.Has(8); b {
		t.Errorf("Has 8 unmatch. %t", b)
	}
}

func TestVlanPortTableHas(t *testing.T) {
	tbl := NewVlanPortTable()
	tbl.SetMinPort(opennsl.Port(50))
	tbl.SetMaxPort(opennsl.Port(55))
	tbl.SetBaseVID(opennsl.Vlan(4000))

	// port: 49 -> x
	// port: 50 -> 4001
	// port: 51 -> 4002
	// ...
	// port: 55 -> 4006
	// port: 56 -> x

	if b := tbl.Has(4000); b {
		t.Errorf("Has 4000 unmatch. %t", b)
	}

	if b := tbl.Has(4001); !b {
		t.Errorf("Has 4001 unmatch. %t", b)
	}

	if b := tbl.Has(4006); !b {
		t.Errorf("Has 4006 unmatch. %t", b)
	}

	if b := tbl.Has(4007); b {
		t.Errorf("Has 4007 unmatch. %t", b)
	}
}
