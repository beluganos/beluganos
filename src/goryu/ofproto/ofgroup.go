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

type Bucket struct {
	Weight    uint32
	WatchPort uint32
	WatchGoup uint32
	Actions   []Action
}

func (b *Bucket) String() string {
	// return fmt.Sprintf("we=%d,wp=%d,wg=%d %s", b.Weight, b.WatchPort, b.WatchGoup, b.Actions)
	return fmt.Sprintf("%s", b.Actions)
}

type GroupEntry struct {
	Type    string
	GroupId uint32
	Buckets []*Bucket
}

func (e *GroupEntry) String() string {
	gName := ConvGroupId(e.GroupId)
	return fmt.Sprintf("gid=%08x(%-24s) %s bkts=%s", e.GroupId, gName, e.Type, e.Buckets)
}

type BucketMod struct {
	Weight    uint32  `json:"weight"`
	WatchPort uint32  `json:"watch_port"`
	WatchGoup uint32  `json:"watch_group"`
	Actions   Actions `json:"actions"`
}

func NewBucket(actions Actions) *BucketMod {
	return &BucketMod{
		Actions: actions,
	}
}

type GroupMod struct {
	Dpid    uint64       `json:"dpid"`
	Type    string       `json:"type"`
	GroupId uint32       `json:"group_id"`
	Buckets []*BucketMod `json:"buckets"`
}

func (g *GroupMod) AddBucket(b *BucketMod) {
	g.Buckets = append(g.Buckets, b)
}

func NewGroupMod(dpid uint64, t string, gid uint32) *GroupMod {
	return &GroupMod{
		Dpid:    dpid,
		Type:    t,
		GroupId: gid,
		Buckets: []*BucketMod{},
	}
}

func NewGroupClear(dpid uint64, t string) *GroupMod {
	return NewGroupMod(dpid, t, 0xfffffffc)
}
