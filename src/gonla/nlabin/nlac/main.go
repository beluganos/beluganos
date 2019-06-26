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
	"gonla/nlaapi"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

type Command struct {
	rootCmd *cobra.Command
	addr    string
	client  nlaapi.NLAApiClient
}

func NewCommand() *Command {
	c := &Command{}
	c.Init()
	return c
}

func (c *Command) Init() {
	rootCmd := &cobra.Command{
		Use:   "nlac",
		Short: "nlac is cli for nlad.",
		Run: func(cmd *cobra.Command, args []string) {
			printAll(c.client)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&c.addr, "addr", "", "127.0.0.1:50062", "NLA API address.")

	monCmd := &cobra.Command{
		Use:   "mon",
		Short: "Monitor messages.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			monNetlink(c.client)
		},
	}
	rootCmd.AddCommand(monCmd)

	targets := []string{
		"link", "addr", "neigh", "route", "mpls", "node", "vpn", "encap", "stat", "iptun", "all",
	}
	showCmd := &cobra.Command{
		Use:       "show [name...]",
		Short:     fmt.Sprintf("show %s", targets),
		Args:      cobra.OnlyValidArgs,
		ValidArgs: targets,
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				switch arg {
				case "link":
					printLinks(c.client)
				case "addr":
					printAddrs(c.client)
				case "neigh":
					printNeighs(c.client)
				case "route":
					printRoutes(c.client)
				case "mpls":
					printMplss(c.client)
				case "node":
					printNodes(c.client)
				case "vpn":
					printVpns(c.client)
				case "encao":
					printEncapInfos(c.client)
				case "stat":
					printStats(c.client)
				case "iptun":
					printIptuns(c.client)
				case "brvlan":
					printBrVlanInfo(c.client)
				case "all":
					printAll(c.client)
				default:
					fmt.Printf("Unknown data. %s\n", arg)
				}
			}
		},
	}
	rootCmd.AddCommand(showCmd)

	c.rootCmd = rootCmd
}

func (c *Command) Execute() error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(c.addr, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	c.client = nlaapi.NewNLAApiClient(conn)

	return c.rootCmd.Execute()
}

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

func strLinkArgs(link netlink.Link) string {
	switch l := link.(type) {
	case *netlink.Vlan:
		return fmt.Sprintf("vid=%d", l.VlanId)

	case *netlink.Veth:
		return fmt.Sprintf("peer='%s'", l.PeerName)

	case *netlink.Bridge:
		s := ""
		if mcsnoop := l.MulticastSnooping; mcsnoop != nil {
			s = fmt.Sprintf("%ssnoop=%t ", s, *mcsnoop)
		}
		if hello := l.HelloTime; hello != nil {
			s = fmt.Sprintf("%shello=%d ", s, *hello)
		}
		if vfilter := l.VlanFiltering; vfilter != nil {
			s = fmt.Sprintf("%svlanfilt=%t ", s, *vfilter)
		}
		return s

	case *netlink.Bond:
		adinfo := func() string {
			if l.AdInfo == nil {
				return ""
			}
			return fmt.Sprintf("ag:%d pt#%d '%s'", l.AdInfo.AggregatorId, l.AdInfo.NumPorts, l.AdInfo.PartnerMac)
		}()

		return fmt.Sprintf("%s act#%d act:'%s' %s", l.Mode, l.ActiveSlave, l.AdActorSystem, adinfo)

	default:
		return ""
	}
}

func strLinkSlaveInfo(attrs *netlink.LinkAttrs) string {
	switch si := attrs.SlaveInfo.(type) {
	case *netlink.BondSlaveInfo:
		return fmt.Sprintf("slave:bond %s ag:%d %s", si.State, si.AggregatorId, si.PermanentHwAddr)
	default:
		return ""
	}
}

func strLink(link *nlaapi.Link) string {
	n := link.ToNative()
	a := n.Attrs()
	linkArgs := strLinkArgs(n.Link)
	slaveInfo := strLinkSlaveInfo(a)
	return fmt.Sprintf("LINK_:%03d:%04d %-16s %-8s %-8s %18s %s i=%d,p=%d,m=%d %s %s",
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
		linkArgs,
		slaveInfo,
	)
}

func printLinks(c nlaapi.NLAApiClient) {
	stream, err := c.GetLinks(context.Background(), nlaapi.NewGetLinksRequest(nlaapi.NODE_ID_ALL))
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
	stream, err := c.GetAddrs(context.Background(), nlaapi.NewGetAddrsRequest(nlaapi.NODE_ID_ALL))
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
	return fmt.Sprintf("NEIGH:%03d:%04d %-32s %-18s %s i=%d/%d,t=%s,v=%d,%d",
		neigh.NId,
		neigh.NeId,
		neigh.GetIP(),
		neigh.NetHardwareAddr(),
		strNeighState(neigh.State),
		neigh.LinkIndex,
		neigh.PhyLink,
		neigh.Tunnel,
		neigh.VlanId,
		neigh.Vni,
	)
}

func printNeighs(c nlaapi.NLAApiClient) {
	stream, err := c.GetNeighs(context.Background(), nlaapi.NewGetNeighsRequest(nlaapi.NODE_ID_ALL))
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
	stream, err := c.GetRoutes(context.Background(), nlaapi.NewGetRoutesRequest(nlaapi.NODE_ID_ALL))
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
	stream, err := c.GetMplss(context.Background(), nlaapi.NewGetMplssRequest(nlaapi.NODE_ID_ALL))
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

func strBrVlanInfo(br *nlaapi.BridgeVlanInfo) string {
	return fmt.Sprintf("BRVLN:%03d:%04d %-16s %d %s %s i=%d,m=%d",
		br.NId,
		br.BrId,
		br.Name,
		br.Vid,
		br.Flags,
		br.PortType(),
		br.Index,
		br.MasterIndex,
	)
}

func printBrVlanInfo(c nlaapi.NLAApiClient) {
	stream, err := c.GetBridgeVlanInfos(context.Background(), &nlaapi.GetBridgeVlanInfosRequest{})
	if err != nil {
		log.Errorf("GetBridgeVlanInfos error. %v", err)
		return
	}
	for {
		br, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv error. %v", err)
			break
		}
		fmt.Println(strBrVlanInfo(br))
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

func strIptun(e *nlaapi.Iptun) string {
	n := e.ToNative()
	a := n.Attrs()
	return fmt.Sprintf("IPTUN:%03d:%04d %-16s remote %-32s local %-32s mac %s mode %s ln=%d",
		n.NId,
		n.TnlId,
		a.Name,
		e.GetRemoteIP(),
		e.GetLocalIP(),
		e.GetLocalMACAddr(),
		n.Type(),
		n.LnId,
	)
}

func printIptuns(c nlaapi.NLAApiClient) {
	stream, err := c.GetIptuns(context.Background(), &nlaapi.GetIptunsRequest{})
	if err != nil {
		log.Errorf("GetTunnels error. %v", err)
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
		fmt.Println(strIptun(e))
	}
}

func strStat(e *nlaapi.Stat) string {
	return fmt.Sprintf("STAT_: %-24s = %d", e.Key, e.Val)
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
		fmt.Println(strStat(stat))
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
	printIptuns(c)
	printBrVlanInfo(c)
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
	case nlalink.RTMGRP_BRIDGE:
		return strBrVlanInfo(m.Msg.GetBrVlanInfo())
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
		now := time.Now().Format("15:04:05")
		t := nlamsg.NlMsgTypeStr(nlmsg.Type())
		fmt.Printf("%s %-16s %s\n", now, t, strNlMsgUion(nlmsg))
	}
}

func main() {

	c := NewCommand()

	if err := c.Execute(); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
