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
	"gonla/nlamsg/nlalink"
	"google.golang.org/grpc"
	"net"
	"os"
)

func sendVpn(args *Args, c nlaapi.NLAApiClient) {
	vpn := &nlaapi.Vpn{
		NId:   args.NId,
		Ip:    args.Dst.IP,
		Mask:  args.Dst.Mask,
		Gw:    args.Gw,
		VpnGw: args.VpnGw,
		Label: args.Label,
	}

	req := &nlaapi.ModVpnRequest{
		Type: args.Cmd,
		Vpn:  vpn,
	}

	if _, err := c.ModVpn(context.Background(), req); err != nil {
		fmt.Printf("ModVpn error. %v\n", err)
	}
}

type Args struct {
	Cmd   uint32
	Dst   *net.IPNet
	Gw    net.IP
	VpnGw net.IP
	Label uint32
	NId   uint32
	Addr  string
}

func (a *Args) Parse() error {
	var cmd string
	var dst string
	var gw string
	var vgw string
	var label uint
	var nid uint

	flag.StringVar(&a.Addr, "addr", "127.0.0.1:50062", "NLA address.")
	flag.StringVar(&cmd, "cmd", "", "command add/del")
	flag.StringVar(&dst, "dst", "100.100.1.0/24", "dest addr.")
	flag.StringVar(&gw, "gw", "1.1.0.1", "gateway addr.")
	flag.StringVar(&vgw, "vgw", "", "vpn gateway addr.")
	flag.UintVar(&label, "label", 20001, "VRF label,")
	flag.UintVar(&nid, "nid", 10, "RIC node id.")
	flag.Parse()

	// Parse Cmd
	switch cmd {
	case "add":
		a.Cmd = nlalink.RTM_NEWVPN
	case "del":
		a.Cmd = nlalink.RTM_DELVPN
	default:
		return fmt.Errorf("Invalid command. %s", cmd)
	}

	// Parse Dst
	_, d, err := net.ParseCIDR(dst)
	if err != nil {
		return err
	}

	// FIX VPN Gw
	if len(vgw) == 0 {
		vgw = gw
	}

	a.Dst = d
	a.Label = uint32(label)
	a.NId = uint32(nid)
	a.Gw = net.ParseIP(gw)
	a.VpnGw = net.ParseIP(vgw)

	if a.Gw == nil || a.VpnGw == nil {
		return fmt.Errorf("Invalid Gateway. GW:%s VPN-GW:%s", a.Gw, a.VpnGw)
	}

	return nil
}

func main() {
	args := Args{}
	if err := args.Parse(); err != nil {
		fmt.Printf("Invalid Argument. %v\n", err)
		os.Exit(1)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(args.Addr, opts...)
	if err != nil {
		fmt.Printf("grpc.Dial error. %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	c := nlaapi.NewNLAApiClient(conn)

	sendVpn(&args, c)
}
