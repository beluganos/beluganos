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
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

const (
	INVALID_HARDWAREADDR = "00:00:00:00:00:00"
	EMPTY_HARDWAREADDR   = ""
)

func GetNetlinkLinks() ([][]byte, error) {
	req := nl.NewNetlinkRequest(unix.RTM_GETLINK, unix.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(unix.AF_UNSPEC))
	return req.Execute(unix.NETLINK_ROUTE, unix.RTM_NEWLINK)
}

func GetNetlinkAddrs() ([][]byte, error) {
	req := nl.NewNetlinkRequest(unix.RTM_GETADDR, unix.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(unix.AF_UNSPEC))
	return req.Execute(unix.NETLINK_ROUTE, unix.RTM_NEWADDR)
}

func GetNetlinkNeighs() ([][]byte, error) {
	req := nl.NewNetlinkRequest(unix.RTM_GETNEIGH, unix.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(unix.AF_UNSPEC))
	return req.Execute(unix.NETLINK_ROUTE, unix.RTM_NEWNEIGH)
}

func GetNetlinkRoutes() ([][]byte, error) {
	req := nl.NewNetlinkRequest(unix.RTM_GETROUTE, unix.NLM_F_DUMP)
	req.AddData(nl.NewIfInfomsg(unix.AF_UNSPEC))
	return req.Execute(unix.NETLINK_ROUTE, unix.RTM_NEWROUTE)
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

func NewFDBEntry(hwaddr net.HardwareAddr, index int, vid uint16, flags int, state int) *netlink.Neigh {
	return &netlink.Neigh{
		LinkIndex:    index,
		IP:           nil,
		Family:       unix.AF_BRIDGE,
		HardwareAddr: hwaddr,
		Vlan:         int(vid),
		Flags:        flags | unix.NTF_MASTER,
		State:        state | unix.NUD_PERMANENT,
	}
}

var bridgeVlanInfoFlags_names = map[uint16]string{
	nl.BRIDGE_VLAN_INFO_MASTER:      "MASTER",
	nl.BRIDGE_VLAN_INFO_PVID:        "PVID",
	nl.BRIDGE_VLAN_INFO_UNTAGGED:    "UNTAGGED",
	nl.BRIDGE_VLAN_INFO_RANGE_BEGIN: "RANGE_BEGIN",
	nl.BRIDGE_VLAN_INFO_RANGE_END:   "RANGE_END",
}

var bridgeVlanInfoFlags_values = map[string]uint16{
	"MASTER":      nl.BRIDGE_VLAN_INFO_MASTER,
	"PVID":        nl.BRIDGE_VLAN_INFO_PVID,
	"UNTAGGED":    nl.BRIDGE_VLAN_INFO_UNTAGGED,
	"RANGE_BEGIN": nl.BRIDGE_VLAN_INFO_RANGE_BEGIN,
	"RANGE_END":   nl.BRIDGE_VLAN_INFO_RANGE_END,
}

func StringBridgeVlanInfoFlags(flags uint16) string {
	names := []string{}
	for val, name := range bridgeVlanInfoFlags_names {
		if (flags & val) != 0 {
			names = append(names, name)
		}
	}
	return strings.Join(names, ",")
}

func ParseBridgeVlanInfoFlags(s string) (uint16, error) {
	names := strings.Split(s, ",")
	var flags uint16
	for _, name := range names {
		v, ok := bridgeVlanInfoFlags_values[name]
		if !ok {
			return 0, fmt.Errorf("Invalid Flags. %s", s)
		}

		flags |= v
	}

	return flags, nil
}
