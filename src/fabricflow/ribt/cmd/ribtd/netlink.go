// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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
	"net"

	"github.com/vishvananda/netlink"
)

func NewIptun(name string, remote, local net.IP) *netlink.Iptun {
	tun := &netlink.Iptun{}
	tun.Attrs().Name = name
	tun.Remote = remote
	tun.Local = local

	return tun
}

func AddLink(link netlink.Link) (netlink.Link, error) {
	if err := netlink.LinkAdd(link); err != nil {
		return nil, err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		netlink.LinkDel(link)
		return nil, err
	}

	newLink, err := netlink.LinkByName(link.Attrs().Name)
	if err != nil {
		netlink.LinkDel(link)
		return nil, err
	}

	return newLink, nil
}

func DelLinkByName(name string) error {
	link := netlink.Dummy{}
	link.Attrs().Name = name
	netlink.LinkSetDown(&link)
	return netlink.LinkDel(&link)
}

func AddRoute(dst *net.IPNet, ifindex int) error {
	route := netlink.Route{
		LinkIndex: ifindex,
		Dst:       dst,
	}

	return netlink.RouteAdd(&route)
}

func DelRoute(dst *net.IPNet, ifindex int) error {
	route := netlink.Route{
		LinkIndex: ifindex,
		Dst:       dst,
	}

	return netlink.RouteDel(&route)
}

func SelectAddr(network *net.IPNet, addrs []netlink.Addr) (*netlink.Addr, error) {
	for _, addr := range addrs {
		if ok := network.Contains(addr.IP); ok {
			return &addr, nil
		}
	}

	return nil, fmt.Errorf("Addr in %s not found.", network)
}

func LinkAddrs(ifname string, family int) ([]netlink.Addr, error) {
	link, err := netlink.LinkByName(ifname)
	if err != nil {
		return nil, err
	}

	return netlink.AddrList(link, family)
}

func LinkAddr(ifname string, network *net.IPNet) (*netlink.Addr, error) {
	family := func() int {
		if network.IP.To4() != nil {
			return netlink.FAMILY_V4
		}
		return netlink.FAMILY_V6
	}()

	addrs, err := LinkAddrs(ifname, family)
	if err != nil {
		return nil, err
	}

	return SelectAddr(network, addrs)
}
