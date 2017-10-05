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

package nlaapi

import (
	"github.com/vishvananda/netlink"
	"testing"
)

func TestLinkType(t *testing.T) {
	datas := map[netlink.Link]string{
		&netlink.Device{}: LINK_TYPE_DEVICE,
		&netlink.Bridge{}: LINK_TYPE_BRIDGE,
		&netlink.Vlan{}:   LINK_TYPE_VLAN,
		&netlink.Vxlan{}:  LINK_TYPE_VXLAN,
		&netlink.Vti{}:    LINK_TYPE_VTI,
		&netlink.Veth{}:   LINK_TYPE_VETH,
		&netlink.Bond{}:   LINK_TYPE_BOND,
	}

	for link, s := range datas {
		if v := link.Type(); v != s {
			t.Errorf("Link.Type unmatch. %s %s", v, s)
		}
	}
}

//
// LinkAddrs
//
func TestLinkAttrs_GetHardwareAddr(t *testing.T) {
	a := LinkAttrs{}

	if v := a.NetHardwareAddr().String(); v != "" {
		t.Errorf("LinkAttrs.GetHardwareAddr unmatch. %s", v)
	}

	a.HardwareAddr = []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}

	if v := a.NetHardwareAddr().String(); v != "11:22:33:44:55:66" {
		t.Errorf("LinkAttrs.GetHardwareAddr unmatch. %s", v)
	}
}

func TestLinkAttrs_ToLocal(t *testing.T) {
	a := LinkAttrs{
		Index:        1,
		Mtu:          2,
		TxQLen:       3,
		Name:         "eth1",
		HardwareAddr: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		Flags:        4,
		RawFlags:     5,
		ParentIndex:  6,
		MasterIndex:  7,
		Alias:        "eth1:10",
		Promisc:      8,
		EncapType:    "encap1",
		OperState:    LinkOperState_OperUp,
	}

	n := netlink.LinkAttrs{}
	LinkAttrsToNative(&a, &n)

	if n.Index != 1 ||
		n.MTU != 2 ||
		n.TxQLen != 3 ||
		n.Name != "eth1" ||
		n.HardwareAddr.String() != "11:22:33:44:55:66" ||
		n.Flags != 4 ||
		n.RawFlags != 5 ||
		n.ParentIndex != 6 ||
		n.MasterIndex != 7 ||
		n.Alias != "eth1:10" ||
		n.Promisc != 8 ||
		n.OperState != 6 {
		t.Errorf("LinkAttrs.ToNative unmatch. %v", n)
	}
}

//
// LinkOperState
//

func TestLinkOperState(t *testing.T) {
	datas := map[LinkOperState]netlink.LinkOperState{
		LinkOperState_OperUnknown:        netlink.OperUnknown,
		LinkOperState_OperNotPresent:     netlink.OperNotPresent,
		LinkOperState_OperDown:           netlink.OperDown,
		LinkOperState_OperLowerLayerDown: netlink.OperLowerLayerDown,
		LinkOperState_OperTesting:        netlink.OperTesting,
		LinkOperState_OperDormant:        netlink.OperDormant,
		LinkOperState_OperUp:             netlink.OperUp,
	}

	for v1, v2 := range datas {
		if int32(v1) != int32(v2) {
			t.Errorf("LinkOperState unmatch %d %d", v1, v2)
		}
	}
}
