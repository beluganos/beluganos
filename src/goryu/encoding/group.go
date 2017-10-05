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

func DecodeGroupEntries(entries []map[string]interface{}) ([]*ofproto.GroupEntry, error) {
	es := make([]*ofproto.GroupEntry, len(entries))
	for i, entry := range entries {
		es[i] = DecodeGroupEntry(entry)
	}

	return es, nil
}

func DecodeGroupEntry(entry map[string]interface{}) *ofproto.GroupEntry {
	e := &ofproto.GroupEntry{}
	for name, field := range entry {
		switch name {
		case "buckets":
			e.Buckets = DecodeBuckets(field.([]interface{}))
		case "group_id":
			e.GroupId = uint32(field.(float64))
		case "type":
			e.Type = field.(string)
		default:
			fmt.Printf("Unknown Field. %s %v\n", name, field)
		}
	}
	return e
}

func DecodeBuckets(entries []interface{}) []*ofproto.Bucket {
	bs := make([]*ofproto.Bucket, len(entries))
	for i, e := range entries {
		bs[i] = DecodeBucket(e.(map[string]interface{}))
	}
	return bs
}

func DecodeBucket(entry map[string]interface{}) *ofproto.Bucket {
	b := &ofproto.Bucket{}
	for name, field := range entry {
		switch name {
		case "weight":
			b.Weight = uint32(field.(float64))
		case "watch_port":
			b.WatchPort = uint32(field.(float64))
		case "watch_group":
			b.WatchGoup = uint32(field.(float64))
		case "actions":
			b.Actions, _ = DecodeActions(field.([]interface{}))
		default:
			fmt.Printf("Unknown Field. %s %v\n", name, field)
		}
	}
	return b
}
