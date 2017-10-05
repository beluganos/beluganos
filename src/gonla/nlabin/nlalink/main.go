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
	"golang.org/x/net/context"
	"gonla/nlaapi"
	"google.golang.org/grpc"
	"os"
	"syscall"
)

type Args struct {
	addr string
	cmd  string
	name string
	nid  uint
}

func (a *Args) Parse() error {
	flag.StringVar(&a.addr, "addr", "127.0.0.1:50062", "nlad address")
	flag.StringVar(&a.cmd, "cmd", "", "command(up/down)")
	flag.StringVar(&a.name, "name", "", "interface name")
	flag.UintVar(&a.nid, "nid", 0, "node-id")
	flag.Parse()

	if len(a.cmd) == 0 {
		return fmt.Errorf("Invalid command. '%s'", a.cmd)
	}

	if len(a.name) == 0 {
		return fmt.Errorf("Invalid interface name. '%s'", a.name)
	}

	return nil
}

func (a *Args) OperState() nlaapi.LinkOperState {
	switch a.cmd {
	case "up":
		return nlaapi.LinkOperState_OperUp
	case "down":
		return nlaapi.LinkOperState_OperDown
	default:
		return nlaapi.LinkOperState_OperUnknown
	}
}

func send(c nlaapi.NLAApiClient, args *Args) {
	nid := uint8(args.nid)
	link := nlaapi.NewDeviceLink(nid, 0)
	attr := link.GetDevice().GetLinkAttrs()
	attr.Name = args.name
	attr.OperState = args.OperState()
	req := nlaapi.NewNetlinkMessageUnion(nid, syscall.RTM_SETLINK, link)
	if _, err := c.ModNetlink(context.Background(), req); err != nil {
		fmt.Printf("ModVpn error. %v\n", err)
	}
}

func main() {
	args := Args{}
	if err := args.Parse(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:50062", opts...)
	if err != nil {
		fmt.Printf("grpc.Dial error. %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	c := nlaapi.NewNLAApiClient(conn)
	send(c, &args)
}
