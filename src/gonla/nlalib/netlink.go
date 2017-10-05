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

package nlalib

import (
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"net"
	"syscall"
)

const (
	INVALID_HARDWAREADDR = "00:00:00:00:00:00"
	EMPTY_HARDWAREADDR   = ""
)

func GetNetlinkLinks() ([][]byte, error) {
	req := nl.NewNetlinkRequest(syscall.RTM_GETLINK, syscall.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(syscall.AF_UNSPEC))
	return req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWLINK)
}

func GetNetlinkAddrs() ([][]byte, error) {
	req := nl.NewNetlinkRequest(syscall.RTM_GETADDR, syscall.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(syscall.AF_UNSPEC))
	return req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWADDR)
}

func GetNetlinkNeighs() ([][]byte, error) {
	req := nl.NewNetlinkRequest(syscall.RTM_GETNEIGH, syscall.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(syscall.AF_UNSPEC))
	return req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWNEIGH)
}

func GetNetlinkRoutes() ([][]byte, error) {
	req := nl.NewNetlinkRequest(syscall.RTM_GETROUTE, syscall.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(syscall.AF_UNSPEC))
	return req.Execute(syscall.NETLINK_ROUTE, syscall.RTM_NEWROUTE)
}

func NewNlMsghdr(t uint16, length uint32) syscall.NlMsghdr {
	return syscall.NlMsghdr{
		Type: t,
		Len:  length,
	}
}

func NewNetlinkMessage(t uint16, data []byte) *syscall.NetlinkMessage {
	return &syscall.NetlinkMessage{
		Header: NewNlMsghdr(t, uint32(len(data))),
		Data:   data,
	}
}

func NewLinkAndAddr(addr *net.IPNet, ifname string) (netlink.Link, *netlink.Addr, error) {
	ln, err := netlink.LinkByName(ifname)
	if err != nil {
		return nil, nil, err
	}

	a := &netlink.Addr{
		IPNet: addr,
		Label: ifname,
	}
	return ln, a, nil
}

func AddIFAddr(addr *net.IPNet, ifname string) error {
	ln, a, err := NewLinkAndAddr(addr, ifname)
	if err != nil {
		return err
	}

	return netlink.AddrAdd(ln, a)
}

func DelIFAddr(addr *net.IPNet, ifname string) error {
	ln, a, err := NewLinkAndAddr(addr, ifname)
	if err != nil {
		return err
	}

	return netlink.AddrDel(ln, a)
}

func NewDummyRoute(dst *net.IPNet, ifname string) (*netlink.Route, error) {
	link, err := netlink.LinkByName(ifname)
	if err != nil {
		return nil, err
	}

	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Scope:     netlink.SCOPE_LINK,
		Dst:       dst,
		Protocol:  0x80,
	}
	return route, nil
}

func SetDummyRoute(dst *net.IPNet, ifname string) error {
	route, err := NewDummyRoute(dst, ifname)
	if err != nil {
		return err
	}
	return netlink.RouteAdd(route)
}

func DelDummyRoute(dst *net.IPNet, ifname string) error {
	route, err := NewDummyRoute(dst, ifname)
	if err != nil {
		return err
	}
	return netlink.RouteDel(route)
}

func IsInvalidHardwareAddr(hwaddr net.HardwareAddr) bool {
	switch hwaddr.String() {
	case INVALID_HARDWAREADDR, EMPTY_HARDWAREADDR:
		return true
	default:
		return false
	}
}
