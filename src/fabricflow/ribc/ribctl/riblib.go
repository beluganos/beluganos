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

package ribctl

import (
	"fabricflow/fibc/api"
	"fmt"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"net"
	"strconv"
	"strings"
	"syscall"
)

func NewVRFLabel(base uint32, nid uint8) uint32 {
	return base + uint32(nid)
}

func NewLinkName(nid uint8, ifname string) string {
	return fmt.Sprintf("%d/%s", nid, ifname)
}

func ParseLinkName(name string) (uint8, string) {
	i := strings.Index(name, "/")
	if i < 0 {
		return 0, name
	}

	nid := func() uint8 {
		if n, err := strconv.Atoi(name[:i]); err == nil {
			return uint8(n)
		}

		return 0
	}()

	return nid, name[i+1:]
}

func NewLinkLinkName(link *nlamsg.Link) string {
	return NewLinkName(link.NId, link.Attrs().Name)
}

func NewAddrLinkName(addr *nlamsg.Addr) string {
	return NewLinkName(addr.NId, addr.Label)
}

func NewPortId(link *nlamsg.Link) uint32 {
	return (uint32(link.NId) << 16) + uint32(link.LnId)
}

func ParsePortId(linkId uint32) (uint8, uint16) {
	return uint8(linkId >> 16), uint16(linkId & 0xffff)
}

func NewNeighId(neigh *nlamsg.Neigh) uint32 {
	if neigh != nil {
		return (uint32(neigh.NId) << 16) + uint32(neigh.NeId)
	} else {
		return 0
	}
}

func GetPortConfigCmd(t uint16) string {
	switch t {
	case syscall.RTM_NEWLINK:
		return "ADD"
	case syscall.RTM_DELLINK:
		return "DELETE"
	case syscall.RTM_SETLINK:
		return "MODIFY"
	default:
		return "NOP"
	}
}

func GetGroupCmd(t uint16) fibcapi.GroupMod_Cmd {
	switch t {
	case syscall.RTM_NEWLINK, syscall.RTM_NEWADDR, syscall.RTM_NEWNEIGH, syscall.RTM_NEWROUTE:
		return fibcapi.GroupMod_ADD

	case syscall.RTM_SETLINK, nlalink.RTM_SETADDR, nlalink.RTM_SETNEIGH, nlalink.RTM_SETROUTE:
		return fibcapi.GroupMod_MODIFY

	case syscall.RTM_DELLINK, syscall.RTM_DELADDR, syscall.RTM_DELNEIGH, syscall.RTM_DELROUTE:
		return fibcapi.GroupMod_DELETE

	default:
		return fibcapi.GroupMod_NOP
	}
}

func GetFlowCmd(t uint16) fibcapi.FlowMod_Cmd {
	switch t {
	case syscall.RTM_NEWLINK, syscall.RTM_NEWADDR, syscall.RTM_NEWNEIGH, syscall.RTM_NEWROUTE:
		return fibcapi.FlowMod_ADD

	case syscall.RTM_SETLINK, nlalink.RTM_SETADDR, nlalink.RTM_SETNEIGH, nlalink.RTM_SETROUTE:
		return fibcapi.FlowMod_MODIFY

	case syscall.RTM_DELLINK, syscall.RTM_DELADDR, syscall.RTM_DELNEIGH, syscall.RTM_DELROUTE:
		return fibcapi.FlowMod_DELETE

	default:
		return fibcapi.FlowMod_NOP
	}
}

func FlowCmdToGroupCmd(cmd fibcapi.FlowMod_Cmd) fibcapi.GroupMod_Cmd {
	switch cmd {
	case fibcapi.FlowMod_ADD:
		return fibcapi.GroupMod_ADD
	case fibcapi.FlowMod_MODIFY:
		return fibcapi.GroupMod_MODIFY
	case fibcapi.FlowMod_DELETE:
		return fibcapi.GroupMod_DELETE
	default:
		return fibcapi.GroupMod_NOP
	}
}

func NewIPNetFromIP(ip net.IP) *net.IPNet {
	bitlen := 128
	if ip.To4() != nil {
		bitlen = 32
	}
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(bitlen, bitlen),
	}
}
