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
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gonla/nlaapi"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"google.golang.org/grpc"
	"io"
	"os"
	"strings"
)

func strNode(node *nlaapi.Node) string {
	return fmt.Sprintf("NODE_:%03d %s", node.NId, node.GetIP())
}

func printNodes(c nlaapi.NLAApiClient) {
	stream, err := c.GetNodes(context.Background(), &nlaapi.GetNodesRequest{})
	if err != nil {
		log.Errorf("GetNodes error. %v", err)
		return
	}

	for {
		node, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Printf("%s\n", strNode(node))
	}
}

func strLink(link *nlaapi.Link) string {
	n := link.ToNative()
	a := n.Attrs()
	return fmt.Sprintf("LINK_:%03d:%04d %-16s %-8s %-8s %18s %s i=%d,p=%d,m=%d",
		link.NId,
		link.LnId,
		a.Name,
		a.Flags,
		a.OperState,
		a.HardwareAddr,
		link.Type,
		a.Index,
		a.ParentIndex,
		a.MasterIndex,
	)
}

func printLinks(c nlaapi.NLAApiClient) {
	stream, err := c.GetLinks(context.Background(), &nlaapi.GetLinksRequest{})
	if err != nil {
		log.Errorf("GetLinks error. %v", err)
		return
	}

	for {
		link, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Printf("%s\n", strLink(link))
	}
}

func strAddr(addr *nlaapi.Addr) string {
	return fmt.Sprintf("ADDR_:%03d:%04d %-32s %-15s i=%d",
		addr.NId,
		addr.AdId,
		addr.GetIPNet(),
		addr.Label,
		addr.Index,
	)
}

func printAddrs(c nlaapi.NLAApiClient) {
	stream, err := c.GetAddrs(context.Background(), &nlaapi.GetAddrsRequest{})
	if err != nil {
		log.Errorf("GetAddrs error. %v", err)
		return
	}

	for {
		addr, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Printf("%s\n", strAddr(addr))
	}
}

func strNeighState(state int32) string {
	s := "|"
	if (state & (1 << 0)) != 0 {
		s = s + "INCOMP|"
	}
	if (state & (1 << 1)) != 0 {
		s = s + "REACH|"
	}
	if (state & (1 << 2)) != 0 {
		s = s + "STALE|"
	}
	if (state & (1 << 3)) != 0 {
		s = s + "DELAY|"
	}
	if (state & (1 << 4)) != 0 {
		s = s + "PROBE|"
	}
	if (state & (1 << 5)) != 0 {
		s = s + "FAILED|"
	}
	if (state & (1 << 6)) != 0 {
		s = s + "NOARP|"
	}
	if (state & (1 << 7)) != 0 {
		s = s + "PARMANENT|"
	}

	return s
}

func strNeigh(neigh *nlaapi.Neigh) string {
	return fmt.Sprintf("NEIGH:%03d:%04d %-32s %-18s %s i=%d",
		neigh.NId,
		neigh.NeId,
		neigh.GetIP(),
		neigh.NetHardwareAddr(),
		strNeighState(neigh.State),
		neigh.LinkIndex,
	)
}

func printNeighs(c nlaapi.NLAApiClient) {
	stream, err := c.GetNeighs(context.Background(), &nlaapi.GetNeighsRequest{})
	if err != nil {
		log.Errorf("GetNeighs error. %v", err)
		return
	}

	for {
		neigh, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Printf("%s\n", strNeigh(neigh))
	}
}

func strRouteTable(t int32) string {
	switch t {
	case 255:
		return "Local"
	case 254:
		return "Main"
	default:
		return fmt.Sprintf("%d", t)
	}
}

func strRouteType(t int32) string {
	switch t {
	case 1:
		return "UC"
	case 2:
		return "Lo"
	case 3:
		return "BC"
	case 4:
		return "AC"
	case 5:
		return "MC"
	default:
		return fmt.Sprintf("%2d", t)
	}
}

func strRouteDest(route *nlaapi.Route) string {
	netDst := route.NetDst()
	if netDst != nil {
		return netDst.String()
	}

	newDst := route.NewDst
	if newDst != nil {
		return fmt.Sprintf("%d to %s", route.MplsDst, newDst.Dest)
	}

	return fmt.Sprintf("%d", route.MplsDst)
}

func strNexthopInfo(n *nlaapi.NexthopInfo) string {
	ss := []string{}
	ss = append(ss, fmt.Sprintf("i=%d", n.LinkIndex))
	ss = append(ss, fmt.Sprintf("via %s", n.NetGw()))
	if d := n.NewDst.GetDest(); d != nil {
		ss = append(ss, fmt.Sprintf("dst %s", d))
	}
	if e := n.Encap.GetEncap(); e != nil {
		ss = append(ss, fmt.Sprintf("enc %s", e))
	}
	return strings.Join(ss, " ")
}

func strNexthopInfos(ns []*nlaapi.NexthopInfo) string {
	ss := make([]string, len(ns))
	for i, n := range ns {
		ss[i] = strNexthopInfo(n)
	}
	return strings.Join(ss, ",")
}

func strRoute(route *nlaapi.Route) string {
	return fmt.Sprintf("ROUTE:%03d:%04d %-32s src %-32s via %-32s %v %s %s i=%d nh=[%s] en=%v",
		route.NId,
		route.RtId,
		strRouteDest(route),
		route.NetSrc(),
		route.NetGw(),
		route.Encap.GetEncap(),
		strRouteType(route.Type),
		strRouteTable(route.Table),
		route.LinkIndex,
		strNexthopInfos(route.MultiPath),
		route.EnIds,
	)
}

func printRoutes(c nlaapi.NLAApiClient) {
	stream, err := c.GetRoutes(context.Background(), &nlaapi.GetRoutesRequest{})
	if err != nil {
		log.Errorf("GetRoutes error. %v", err)
		return
	}

	for {
		route, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Printf("%s\n", strRoute(route))
	}
}

func strMpls(route *nlaapi.Route) string {
	return fmt.Sprintf("MPLS :%03d:%04d %-32s src %-32s via %-32s %v %s %s i=%d nh=[%s] en=%v",
		route.NId,
		route.RtId,
		strRouteDest(route),
		route.NetSrc(),
		route.NetGw(),
		route.Encap.GetEncap(),
		strRouteType(route.Type),
		strRouteTable(route.Table),
		route.LinkIndex,
		strNexthopInfos(route.MultiPath),
		route.EnIds,
	)
}

func printMplss(c nlaapi.NLAApiClient) {
	stream, err := c.GetMplss(context.Background(), &nlaapi.GetMplssRequest{})
	if err != nil {
		log.Errorf("GetRoutes error. %v", err)
		return
	}

	for {
		mpls, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Printf("%s\n", strMpls(mpls))
	}
}

func strVpn(vpn *nlaapi.Vpn) string {
	return fmt.Sprintf("VPN  :%03d:%04d %-32s %-15s(%-15s) %5d",
		vpn.NId,
		vpn.VpnId,
		vpn.GetIPNet(),
		vpn.NetGw(),
		vpn.NetVpnGw(),
		vpn.Label,
	)
}

func printVpns(c nlaapi.NLAApiClient) {
	stream, err := c.GetVpns(context.Background(), &nlaapi.GetVpnsRequest{})
	if err != nil {
		log.Errorf("GetVpns error. %v", err)
		return
	}
	for {
		vpn, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Println(strVpn(vpn))
	}
}

func strEncapInfo(e *nlaapi.EncapInfo) string {
	return fmt.Sprintf("Encap:0000:%04d %-32s %d",
		e.EnId,
		e.NetDst(),
		e.Vrf)
}

func printEncapInfos(c nlaapi.NLAApiClient) {
	stream, err := c.GetEncapInfos(context.Background(), &nlaapi.GetEncapInfosRequest{})
	if err != nil {
		log.Errorf("GetEncapInfos error. %v", err)
		return
	}
	for {
		e, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Println(strEncapInfo(e))
	}
}

func printStats(c nlaapi.NLAApiClient) {
	stream, err := c.GetStats(context.Background(), &nlaapi.GetStatsRequest{})
	if err != nil {
		log.Errorf("GetStats error. %v", err)
		return
	}
	for {
		stat, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Println(stat)
	}
}

func printAll(c nlaapi.NLAApiClient) {
	printLinks(c)
	printAddrs(c)
	printNeighs(c)
	printRoutes(c)
	printMplss(c)
	printVpns(c)
	printNodes(c)
	printEncapInfos(c)
	printStats(c)
}

func strNlMsgUion(m *nlaapi.NetlinkMessageUnion) string {
	switch m.Group() {
	case nlalink.RTMGRP_LINK:
		return strLink(m.Msg.GetLink())
	case nlalink.RTMGRP_ADDR:
		return strAddr(m.Msg.GetAddr())
	case nlalink.RTMGRP_NEIGH:
		return strNeigh(m.Msg.GetNeigh())
	case nlalink.RTMGRP_ROUTE:
		return strRoute(m.Msg.GetRoute())
	case nlalink.RTMGRP_NODE:
		return strNode(m.Msg.GetNode())
	case nlalink.RTMGRP_VPN:
		return strVpn(m.Msg.GetVpn())
	default:
		return fmt.Sprintf("%v", m)
	}
}

func monNetlink(c nlaapi.NLAApiClient) {
	stream, err := c.MonNetlink(context.Background(), &nlaapi.MonNetlinkRequest{})
	if err != nil {
		log.Errorf("MonNetlink error. %v", err)
		return
	}

	for {
		nlmsg, err := stream.Recv()
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		t := nlamsg.NlMsgTypeStr(nlmsg.Type())
		fmt.Printf("%-12s %s\n", t, strNlMsgUion(nlmsg))
	}
}

type Args struct {
	Cmd  string
	Name string
	Addr string
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" nlac [command] [name]")
	fmt.Println(" - command:")
	fmt.Println("   - dump : dump [name] data. (default)")
	fmt.Println("   - mon  : monitor netlink messages.")
	fmt.Println(" - name:")
	fmt.Println("   - link : Dump all links.")
	fmt.Println("   - addr : Dump all addresses.")
	fmt.Println("   - neigh: Dump all neighbors.")
	fmt.Println("   - route: Dump all routes (IPv4/6).")
	fmt.Println("   - mpls : Dump add routes (MPL).")
	fmt.Println("   - vpn  : Dump all VPNs.")
	fmt.Println("   - node : Dump all nodes.")
	fmt.Println("   - encap: Dump all encap infos.")
	fmt.Println("   - all  : Dump all target datas.(default)")

	os.Exit(1)
}

func getargs() (args *Args) {
	addr := flag.String("addr", "127.0.0.1:50062", "NLA API address.")
	flag.Parse()

	args = &Args{
		Cmd:  "dump",
		Name: "all",
		Addr: *addr,
	}

	a := flag.Args()
	if len(a) >= 1 {
		args.Cmd = a[0]
	}

	if len(a) >= 2 {
		args.Name = a[1]
	}

	return
}

func main() {

	args := getargs()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(args.Addr, opts...)
	if err != nil {
		log.Errorf("grpc.Dial error. %v", err)
		return
	}
	defer conn.Close()

	c := nlaapi.NewNLAApiClient(conn)

	switch args.Cmd {
	case "dump":
		switch args.Name {
		case "link":
			printLinks(c)
		case "addr":
			printAddrs(c)
		case "neigh":
			printNeighs(c)
		case "route":
			printRoutes(c)
		case "mpls":
			printMplss(c)
		case "node":
			printNodes(c)
		case "vpn":
			printVpns(c)
		case "encao":
			printEncapInfos(c)
		case "stat":
			printStats(c)
		case "all":
			printAll(c)
		default:
			printUsage()
		}
	case "mon":
		monNetlink(c)

	default:
		printUsage()
	}
}
