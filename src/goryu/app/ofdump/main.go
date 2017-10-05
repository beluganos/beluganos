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
	"flag"
	"fmt"
	"goryu/ryulib"
)

func main() {

	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "ryu addr")
	flag.Parse()

	baseUrl := fmt.Sprintf("http://%s/stats", addr)
	c := ryulib.NewClient(baseUrl)

	dpids, err := c.GetSwitches()
	if err != nil {
		fmt.Printf("GET ERROR %s\n", err)
		return
	}

	fmt.Printf("%v\n", dpids)

	for _, dpid := range dpids {
		fmt.Printf("--- dp %d/0x%x\n", dpid, dpid)
		desc, err := c.GetDesc(dpid)
		fmt.Printf("%s %s\n", desc, err)

		flow, err := c.GetFlow(dpid)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}

		for _, entry := range flow {
			fmt.Printf("%s\n", entry)
		}

		groups, err := c.GetGroup(dpid)
		if err != nil {
			fmt.Printf("GetGroup error. %s\n", err)
			continue
		}
		for _, group := range groups {
			fmt.Println(group)
		}
	}

}
