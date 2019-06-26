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

package fflibnet

import (
	"fmt"
	"strings"

	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

//
// nl.BRIDGE_FLAGS_XXX
//
var bridgeFlags_name = map[uint16]string{
	nl.BRIDGE_FLAGS_MASTER: "master",
	nl.BRIDGE_FLAGS_SELF:   "self",
}

var bridgeFlags_value = map[string]uint16{
	"master": nl.BRIDGE_FLAGS_MASTER,
	"self":   nl.BRIDGE_FLAGS_SELF,
}

func StringBridgeFlag(v uint16) string {
	if s, ok := bridgeFlags_name[v]; ok {
		return s
	}
	return fmt.Sprintf("BridgeFlag(%d)", v)
}

func ParseBridgeFlag(s string) (uint16, error) {
	if v, ok := bridgeFlags_value[s]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("Invalid bridge flag. %s", s)
}

//
// unix.NUD_XXX
//
var bridgeState_names = map[int]string{
	unix.NUD_INCOMPLETE: "incomplete",
	unix.NUD_REACHABLE:  "reachable",
	unix.NUD_STALE:      "stale",
	unix.NUD_DELAY:      "delay",
	unix.NUD_PROBE:      "probe",
	unix.NUD_FAILED:     "failed",
	unix.NUD_NOARP:      "noarp",
	unix.NUD_PERMANENT:  "permanent",
	unix.NUD_NONE:       "none",
}

var bridgeState_values = map[string]int{
	"incomplete": unix.NUD_INCOMPLETE,
	"reachable":  unix.NUD_REACHABLE,
	"stale":      unix.NUD_STALE,
	"delay":      unix.NUD_DELAY,
	"probe":      unix.NUD_PROBE,
	"failed":     unix.NUD_FAILED,
	"noarp":      unix.NUD_NOARP,
	"permanent":  unix.NUD_PERMANENT,
	"none":       unix.NUD_NONE,
}

func StringBridgeState(v int) string {
	if s, ok := bridgeState_names[v]; ok {
		return s
	}
	return fmt.Sprintf("BridgeState(%d)", v)
}

func StringBridgeStates(v int, delim string) string {
	if v == unix.NUD_NONE {
		return StringBridgeState(v)
	}

	names := []string{}
	for name, val := range bridgeState_values {
		if (v & val) != 0 {
			names = append(names, name)
		}
	}

	return strings.Join(names, delim)
}

func ParseBridgeState(s string) (int, error) {
	if v, ok := bridgeState_values[s]; ok {
		return v, nil
	}
	return unix.NUD_NONE, fmt.Errorf("Invalid brige state. %s", s)
}
