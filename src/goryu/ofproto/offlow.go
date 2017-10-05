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

package ofproto

import (
	"fmt"
)

type FlowEntry struct {
	Length       uint16
	TableId      uint8
	DurationSec  uint32
	DurationNsec uint32
	Priority     uint16
	IdleTimeout  uint16
	HardTimeout  uint16
	Flags        uint16
	Cookie       uint64
	PacketCount  uint64
	ByteCount    uint64
	Match        Match
	ApplyActions []Action
	WriteActions []Action
}

func (e *FlowEntry) String() string {
	return fmt.Sprintf("tbl=%d pri=%d cnt=%d/%d m=%s a=%s w=%s",
		e.TableId, e.Priority, e.PacketCount, e.ByteCount, e.Match, e.ApplyActions, e.WriteActions)
}

type FlowMod struct {
	Dpid        uint64  `json:"dpid"`
	TableId     uint8   `json:"table_id"`
	Cookie      uint64  `json:"cookie"`
	CookieMask  uint64  `json:"cookie_mask"`
	Priority    uint16  `json:"priority"`
	IdleTimeout uint16  `json:"idle_timeout"`
	HardTimeout uint16  `json:"hard_timeout"`
	Flags       uint16  `json:"flags"`
	Match       Match   `json:"match"`
	Actions     Actions `json:"actions"`
}

func NewFlowMod(dpid uint64) *FlowMod {
	return &FlowMod{
		Dpid:    dpid,
		Match:   Match{},
		Actions: Actions{},
	}
}

func NewFlowClear(dpid uint64) *FlowMod {
	f := NewFlowMod(dpid)
	f.TableId = 0xff
	return f
}
