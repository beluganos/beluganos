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

package ryuenc

import (
	"fmt"
	"goryu/ofproto"
)

func DecodeFlowEntries(rs []map[string]interface{}) []*ofproto.FlowEntry {
	es := make([]*ofproto.FlowEntry, len(rs))
	for i, r := range rs {
		es[i] = DecodeFlowEntry(r)
	}
	return es
}

func DecodeFlowEntry(r map[string]interface{}) *ofproto.FlowEntry {
	e := &ofproto.FlowEntry{}
	for name, field := range r {
		switch name {
		case "length":
			e.Length = uint16(field.(float64))
		case "table_id":
			e.TableId = uint8(field.(float64))
		case "duration_sec":
			e.DurationSec = uint32(field.(float64))
		case "duration_nsec":
			e.DurationNsec = uint32(field.(float64))
		case "priority":
			e.Priority = uint16(field.(float64))
		case "idle_timeout":
			e.IdleTimeout = uint16(field.(float64))
		case "hard_timeout":
			e.HardTimeout = uint16(field.(float64))
		case "flags":
			e.Flags = uint16(field.(float64))
		case "cookie":
			e.Cookie = uint64(field.(float64))
		case "packet_count":
			e.PacketCount = uint64(field.(float64))
		case "byte_count":
			e.ByteCount = uint64(field.(float64))
		case "match":
			e.Match = field.(map[string]interface{})
		case "actions":
			e.ApplyActions, e.WriteActions = DecodeActions(field.([]interface{}))
		default:
			fmt.Printf("Unknown Field. %s %v\n", name, field)
		}
	}
	return e
}
