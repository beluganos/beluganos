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

package main

import (
	"fmt"
	"goryu/ofproto"
	"goryu/ryulib"
)

const DPID = 14

func flowtest1(client *ryulib.RyuClient) {
	m := ofproto.Match{}
	m.Set("in_port", 1).Set("eth_type", 0x0800)
	w := ofproto.Actions{}
	w.SetField("eth_dst", "11:22:33:44:55:66").Output(1)
	a := ofproto.Actions{}
	a.WriteActions(w).GotoTable(10)

	f := ofproto.NewFlowMod(DPID)
	f.Match = m
	f.Actions = a

	if err := client.ModFlow("add", f); err != nil {
		fmt.Printf("Mod Flow error. %s", err)
	}
}

func grouptest1(client *ryulib.RyuClient) {
	a := ofproto.Actions{}
	a.SetField("eth_dst", "11:22:33:44:55:66").Output(1)
	g := ofproto.NewGroupMod(DPID, "ALL", 0x00000001)
	g.AddBucket(ofproto.NewBucket(a))

	if err := client.ModGroup("add", g); err != nil {
		fmt.Printf("Mof Group error. %s\n", err)
	}
}

func flowclear(client *ryulib.RyuClient) {
	f := ofproto.NewFlowClear(DPID)
	if err := client.ModFlow("delete", f); err != nil {
		fmt.Printf("Mod Flow error. %s", err)
	}
}

func groupclear(client *ryulib.RyuClient) {
	for _, gtype := range []string{"FF", "ALL", "SELECT", "INDIRECT"} {
		g := ofproto.NewGroupClear(DPID, gtype)
		if err := client.ModGroup("delete", g); err != nil {
			fmt.Printf("Mod Group error. %s\n", err)
		}
	}
}

func main() {

	c := ryulib.NewClient("http://127.0.0.1:8080/stats")
	flowtest1(c)
	flowclear(c)
	grouptest1(c)
	groupclear(c)
}
