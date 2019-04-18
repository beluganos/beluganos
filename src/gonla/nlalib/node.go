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

package nlalib

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

func NewNodeIdFromIP(ip net.IP) uint8 {
	if l := len(ip); l > 0 {
		return uint8(ip[l-1])
	}
	return 0
}

func NewNodeIdFromIF(ifname string) (uint8, error) {
	link, err := netlink.LinkByName(ifname)
	if err != nil {
		return 0, err
	}

	addrs, err := netlink.AddrList(link, nl.FAMILY_V4)
	if err != nil {
		return 0, err
	}

	for _, addr := range addrs {
		if ipnet := addr.IPNet; ipnet != nil {
			if ones, bits := ipnet.Mask.Size(); ones != bits && ones != 0 {
				return NewNodeIdFromIP(ipnet.IP), nil
			}
		}
	}

	return 0, fmt.Errorf("global ipv4-addr not found.")
}
