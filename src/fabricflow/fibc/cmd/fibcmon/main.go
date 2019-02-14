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
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"
	"flag"
	"fmt"
	"time"
)

func monitor(con *fibcnet.Connection) {

	defer con.Close()

	hello := fibcapi.NewHello("0.0.0.0")
	if err := con.Write(hello, 0); err != nil {
		fmt.Printf("send: hello error. %s", err)
		return
	}

	for {
		hdr, data, err := con.Read()
		if err != nil {
			fmt.Printf("fibcon.Read error. %s\n", err)
			return
		}

		switch fibcapi.FFM(hdr.Type) {
		case fibcapi.FFM_PORT_STATUS:
			msg, err := fibcapi.NewPortStatusFromBytes(data)
			if err != nil {
				fmt.Printf("recv: NewPortStatusFromBytes error. %s\n", err)
				return
			}
			fmt.Printf("recv: %v\n", msg)

		default:
			fmt.Printf("recv: Unknown Message. type=%d\n", hdr.Type)
		}
	}
}

func main() {

	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:50070", "fibc addr.")
	flag.Parse()

	con := fibcnet.NewConnection(addr)
	for {
		if err := con.Connect(); err == nil {
			monitor(con)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}
