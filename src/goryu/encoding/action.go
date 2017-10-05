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
	"strconv"
	"strings"
)

//
// Generic Action
//
func DecodeBaseAction(fields []string) (ofproto.Action, error) {
	a := &ofproto.BaseAction{}
	a.Name = fields[0]
	if len(fields) > 1 {
		a.StrValue = fields[1]
		if n, err := strconv.ParseInt(fields[1], 0, 64); err == nil {
			a.IntValue = n
		}
	}
	return a, nil
}

//
// GotoTable Action
//
func DecodeGotoTableAction(fields []string) (ofproto.Action, error) {
	a := &ofproto.GotoTableAction{}
	n, err := strconv.ParseInt(fields[1], 0, 64)
	if err != nil {
		return nil, err
	}
	a.TableNo = uint8(n)
	return a, nil
}

//
// Group Action
//
func DecodeGroupAction(fields []string) (ofproto.Action, error) {
	a := &ofproto.GroupAction{}
	n, err := strconv.ParseInt(fields[1], 0, 64)
	if err != nil {
		return nil, err
	}
	a.GroupId = uint32(n)
	return a, nil
}

//
// Decode Action
//
func DecodeAction(action string) ofproto.Action {
	var a ofproto.Action = nil
	var err error
	fields := strings.SplitN(action, ":", 2)
	switch fields[0] {
	case "GOTO_TABLE":
		a, err = DecodeGotoTableAction(fields)
	case "GROUP":
		a, err = DecodeGroupAction(fields)
	default:
		a, err = DecodeBaseAction(fields)
	}

	if err != nil {
		return nil
	}

	return a
}

func DecodeActions(actions []interface{}) ([]ofproto.Action, []ofproto.Action) {
	applyActions := []ofproto.Action{}
	writeActions := []ofproto.Action{}

	for _, action := range actions {
		switch action.(type) {
		case map[string]interface{}:
			m := action.(map[string]interface{})
			writeActions, _ = DecodeActions(m["WRITE_ACTIONS"].([]interface{}))
		case string:
			applyActions = append(applyActions, DecodeAction(action.(string)))
		default:
			fmt.Printf("unonown action type. %v", action)
		}
	}
	return applyActions, writeActions
}
