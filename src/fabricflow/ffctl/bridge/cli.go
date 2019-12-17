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

package bridge

import (
	"encoding/json"
	fflibnet "fabricflow/util/net"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type BridgeFdbCLI struct {
	Mac    string   `json:"mac"`
	Dev    string   `json:"dev"`              // device name
	Master string   `json:"master,omitempty"` // device name
	Vlan   uint16   `json:"vlan,omitempty"`   // vlan vid
	Flags  []string `json:"flags,omitempty"`  // nl.BRIDGE_FLAGS_XXX
	State  string   `json:"state,omitempty"`  // unix.NUD_XXX
}

func ParseBridgeFdbCLI(b []byte) ([]*BridgeFdbCLI, error) {
	data := []*BridgeFdbCLI{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func ExecBridgeFdbShow() ([]*BridgeFdbCLI, error) {
	out, err := exec.Command("bridge", "-j", "fdb", "show").Output()
	if err != nil {
		return nil, err
	}

	fdbCLIs, err := ParseBridgeFdbCLI(out)
	if err != nil {
		return nil, err
	}

	return fdbCLIs, nil
}

func (b *BridgeFdbCLI) String() string {
	ss := []string{b.Mac, "dev", b.Dev}

	if b.Vlan != 0 {
		ss = append(ss, "vlan", fmt.Sprintf("%d", b.Vlan))
	}

	if len(b.Master) != 0 {
		ss = append(ss, "master", b.Master)
	}

	if len(b.Flags) != 0 {
		ss = append(ss, b.Flags...)
	}

	if len(b.State) != 0 {
		ss = append(ss, b.State)
	}

	return strings.Join(ss, " ")
}

func (b *BridgeFdbCLI) GetHardwareAddr() net.HardwareAddr {
	if hwaddr, err := net.ParseMAC(b.Mac); err == nil {
		return hwaddr
	}
	return nil
}

func (b *BridgeFdbCLI) GetFlags() []uint16 {
	flags := []uint16{}
	for _, s := range b.Flags {
		if flag, err := fflibnet.ParseBridgeFlag(s); err == nil {
			flags = append(flags, flag)
		}
	}

	return flags
}

func (b *BridgeFdbCLI) GetState() int {
	state, _ := fflibnet.ParseBridgeState(b.State)
	return state
}
