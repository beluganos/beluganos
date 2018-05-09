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

package fibcapi

import (
	"testing"
)

//
// PortStatus
//
func TestPortStatus_Type(t *testing.T) {
	p := PortStatus{}

	if v := p.Type(); v != uint16(FFM_PORT_STATUS) {
		t.Errorf("PortStatus Type unmatch. %d", v)
	}
}

func TestPortStatus_Bytes(t *testing.T) {
	p := PortStatus{
		Status: PortStatus_UP,
		ReId:   "1.1.1.1",
		PortId: 0x12345678,
		Ifname: "ethX",
	}

	b, err := p.Bytes()

	if err != nil {
		t.Errorf("PortStatus  Bytes error. %v", err)
	}

	if v := len(b); v == 0 {
		t.Errorf("PortStatus  Bytes unmatch. %v", v)
	}
}

func TestPortStatus_NewFromBytes(t *testing.T) {
	p := &PortStatus{
		Status: PortStatus_UP,
		ReId:   "1.1.1.1",
		PortId: 0x12345678,
		Ifname: "ethX",
	}
	b, _ := p.Bytes()

	d, err := NewPortStatusFromBytes(b)

	if err != nil {
		t.Errorf("NewPortStatusFromBytes error. %s", err)
	}
	if d.Status != PortStatus_UP {
		t.Errorf("NewPortStatusFromBytes unmatch. Status=%d", d.Status)
	}
	if d.ReId != "1.1.1.1" {
		t.Errorf("NewPortStatusFromBytes unmatch. ReId=%s", d.ReId)
	}
	if d.PortId != 0x12345678 {
		t.Errorf("NewPortStatusFromBytes unmatch. PortId=%d", d.PortId)
	}
	if d.Ifname != "ethX" {
		t.Errorf("NewPortStatusFromBytes unmatch. Ifname=%s", d.Ifname)
	}
}

func TestPortStatus_NewFromBytes_err(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5}

	_, err := NewPortStatusFromBytes(b)

	if err == nil {
		t.Errorf("NewPortStatusFromBytes must be error. %s", err)
	}
}

//
// PortConfig
//
func TestPortConfig_Type(t *testing.T) {
	p := PortConfig{}

	if v := p.Type(); v != uint16(FFM_PORT_CONFIG) {
		t.Errorf("PortStatus Type unmatch. %d", v)
	}
}

func TestPortConfig_Bytes(t *testing.T) {
	p := PortConfig{
		Cmd:    PortConfig_ADD,
		ReId:   "1.1.1.1",
		Ifname: "ethX",
		PortId: 100,
		Status: PortStatus_UP,
	}

	b, err := p.Bytes()

	if err != nil {
		t.Errorf("PortConfig Bytes error. %s", err)
	}

	if v := len(b); v == 0 {
		t.Errorf("PortConfig Bytes unmatch. %d", v)
	}
}

func TestPortConfig_New(t *testing.T) {
	m := map[string]PortConfig_Cmd{
		"ADD":    PortConfig_ADD,
		"MODIFY": PortConfig_MODIFY,
		"DELETE": PortConfig_DELETE,
	}

	for key, val := range m {
		p := NewPortConfig(key, "1.1.1.1", "ethX", 100, PortStatus_UP)
		if v := p.Cmd; v != val {
			t.Errorf("NewPortConfig Cmd unmatch. %d", v)
		}
		if v := p.ReId; v != "1.1.1.1" {
			t.Errorf("NewPortConfig ReId unmatch. %s", v)
		}
		if v := p.Ifname; v != "ethX" {
			t.Errorf("NewPortConfig Ifname unmatch. %s", v)
		}
		if v := p.PortId; v != 100 {
			t.Errorf("NewPortConfig PortId unmatch. %d", v)
		}
		if v := p.Status; v != PortStatus_UP {
			t.Errorf("NewPortConfig Status unmatch. %d", v)
		}
	}
}
