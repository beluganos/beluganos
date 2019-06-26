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

package nlamsg

import (
	"bytes"
	"gonla/nlalib"
	"testing"
)

func TestBridgeVlanInfoSerialize(t *testing.T) {
	br := BridgeVlanInfo{}
	br.Flags = 0x1111
	br.Vid = 0x2222
	br.Index = 0x33333333
	br.Name = "eth1234"
	br.MasterIndex = 0x44444444
	br.Mtu = 0x5555
	br.BrId = 0x66666666
	br.NId = 0x77

	d := []byte{
		0x11, 0x11, 0x22, 0x22, // Flags, Vid
		0x33, 0x33, 0x33, 0x33, // Index
		0x65, 0x74, 0x68, 0x31, // "eth1"
		0x32, 0x33, 0x34, 0x00, // "234"
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x44, 0x44, 0x44, 0x44, //MasterIndex
		0x00, 0x00, 0x55, 0x55, // Mtu (uint32)
	}

	nlmsg, err := BridgeVlanInfoSerialize(&br, 0x8888)

	if err != nil {
		t.Errorf("BridgeVlanInfoSerialize error. %s", err)
	}
	if v := bytes.Compare(nlmsg.Data, d); v != 0 {
		t.Errorf("BridgeVlanInfoSerialize unmatch. %v %v", d, nlmsg.Data)
	}
	if v := nlmsg.NId; v != 0x77 {
		t.Errorf("BridgeVlanInfoSerialize unmatch. nid=%d", v)
	}
	if v := nlmsg.Src; v != SRC_NOP {
		t.Errorf("BridgeVlanInfoSerialize unmatch. src=%d", v)
	}
	if v := nlmsg.Header.Type; v != 0x8888 {
		t.Errorf("BridgeVlanInfoSerialize unmatch. type=%d", v)
	}
}

func TestBridgeVlanInfoDeserialize(t *testing.T) {
	d := []byte{
		0x11, 0x11, 0x22, 0x22, // Flags, Vid
		0x33, 0x44, 0x33, 0x44, // Index
		0x65, 0x74, 0x68, 0x31, // "eth1"
		0x32, 0x33, 0x34, 0x00, // "234"
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x44, 0x44, 0x44, 0x44, //MasterIndex
		0x00, 0x00, 0x55, 0x55, // Mtu (uint32)
	}

	nlmsg := NewNetlinkMessage(nlalib.NewNetlinkMessage(0, d), 0, SRC_NOP)
	br, err := BridgeVlanInfoDeserialize(nlmsg)

	if err != nil {
		t.Errorf("BridgeVlanInfoDeserialize error. %s", err)
	}
	if v := br.Flags; v != 0x1111 {
		t.Errorf("BridgeVlanInfoDeserialize unmatch. flags=%d", v)
	}
	if v := br.Vid; v != 0x2222 {
		t.Errorf("BridgeVlanInfoDeserialize unmatch. vid=%d", v)
	}
	if v := br.Index; v != 0x33443344 {
		t.Errorf("BridgeVlanInfoDeserialize unmatch. index=%d", v)
	}
	if v := br.Name; v != "eth1234" {
		t.Errorf("BridgeVlanInfoDeserialize unmatch. name=%v", v)
	}
}
