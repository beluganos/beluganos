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
	"strings"
)

//
// Action interface
//
type Action interface {
}

//
// Action (Base)
//
type BaseAction struct {
	Name     string
	StrValue string
	IntValue int64
}

func (a *BaseAction) String() string {
	if a.IntValue == 0 {
		return fmt.Sprintf("%s:\"%s\"", a.Name, strings.Trim(a.StrValue, " "))
	} else {
		return fmt.Sprintf("%s:0x%x", a.Name, a.IntValue)
	}
}

//
// Action (GOTO_TABLE)
//
type GotoTableAction struct {
	TableNo uint8
}

func (a *GotoTableAction) String() string {
	return fmt.Sprintf("GOTO_TABLE:%02d(%s)", a.TableNo, StrTable(a.TableNo))
}

//
// Action (Group)
//
type GroupAction struct {
	GroupId uint32
}

func (a *GroupAction) String() string {
	return fmt.Sprintf("GROUP:%08x(%s)", a.GroupId, ConvGroupId(a.GroupId))
}

type Actions []map[string]interface{}

func NewAction0(t string) map[string]interface{} {
	return map[string]interface{}{
		"type": t,
	}
}

func NewAction(t string, key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": t,
		key:    value,
	}
}

func (a *Actions) Append(m map[string]interface{}) *Actions {
	*a = append(*a, m)
	return a
}

func (a *Actions) Output(port uint32) *Actions {
	return a.Append(NewAction("OUTPUT", "port", port))
}

func (a *Actions) CopyTTLOut() *Actions {
	return a.Append(NewAction0("COPY_TTL_OUT"))
}

func (a *Actions) CopyTTLIn() *Actions {
	return a.Append(NewAction0("COPY_TTL_IN"))
}

func (a *Actions) SetMPLSTTL(mpls_ttl uint8) *Actions {
	return a.Append(NewAction("SET_MPLS_TTL", "mpls_ttl", mpls_ttl))
}

func (a *Actions) DecMPLSTTL() *Actions {
	return a.Append(NewAction0("DEC_MPLS_TTL"))
}

func (a *Actions) PushVLAN(ethertype uint16) *Actions {
	return a.Append(NewAction("PUSH_VLAN", "ethertype", ethertype))
}

func (a *Actions) PopVLAN() *Actions {
	return a.Append(NewAction0("POP_VLAN"))
}

func (a *Actions) PushMPLS(ethertype uint16) *Actions {
	return a.Append(NewAction("PUSH_MPLS", "ethertype", ethertype))
}

func (a *Actions) PopMPLS(ethertype uint16) *Actions {
	return a.Append(NewAction("POP_MPLS", "ethertype", ethertype))
}

func (a *Actions) SetQueue(queue_id uint32) *Actions {
	return a.Append(NewAction("SET_QUEUE", "queue_id", queue_id))
}

func (a *Actions) Group(group_id uint32) *Actions {
	return a.Append(NewAction("GROUP", "group_id", group_id))
}

func (a *Actions) SetNWTTL(nw_ttl uint8) *Actions {
	return a.Append(NewAction("SET_NW_TTL", "nw_ttl", nw_ttl))
}

func (a *Actions) DecNWTTL() *Actions {
	return a.Append(NewAction0("DEC_NW_TTL"))
}

func (a *Actions) SetField(field string, value interface{}) *Actions {
	m := map[string]interface{}{
		"type":  "SET_FIELD",
		"field": field,
		"value": value,
	}
	return a.Append(m)
}

func (a *Actions) PushPBB(ethertype uint16) *Actions {
	return a.Append(NewAction("PUSH_PBB", "ethertype", ethertype))
}

func (a *Actions) PopPBB() *Actions {
	return a.Append(NewAction0("POP_PBB"))
}

func (a *Actions) WriteActions(actions Actions) *Actions {
	return a.Append(NewAction("WRITE_ACTIONS", "actions", actions))
}

func (a *Actions) ClearActions() *Actions {
	return a.Append(NewAction0("CLEAR_ACTIONS"))
}

func (a *Actions) GotoTable(table_id uint8) *Actions {
	return a.Append(NewAction("GOTO_TABLE", "table_id", table_id))
}

func (a *Actions) WriteMetadata(metadata uint64, metadata_mask uint64) *Actions {
	m := map[string]interface{}{
		"type":          "WRITE_METADATA",
		"metadata":      metadata,
		"metadata_mask": metadata_mask,
	}
	return a.Append(m)
}

func (a *Actions) Meter(meter_id uint32) *Actions {
	return a.Append(NewAction("METER", "meter_id", meter_id))
}

func (a *Actions) SetVRF(vrf uint8, useMetadata bool) *Actions {
	if useMetadata {
		return a.WriteMetadata(uint64(vrf)<<VRF_METADATA_SHIFT, VRF_METADATA_MASK)
	} else {
		return a.SetField("vrf", vrf)
	}
}

func (a *Actions) SetMPLSType(mpls_type uint8, useMetadata bool) *Actions {
	if useMetadata {
		return a.WriteMetadata(uint64(mpls_type)<<MPLSTYPE_METADATA_SHIFT, MPLSTYPE_METADATA_MASK)
	} else {
		return a.SetField("mpls_type", mpls_type)
	}
}
