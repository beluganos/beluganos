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

package fflibnet

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func GetIPv4Addrs(ifname string) ([]netlink.Addr, error) {
	link, err := netlink.LinkByName(ifname)
	if err != nil {
		return nil, err
	}

	return netlink.AddrList(link, netlink.FAMILY_V4)
}

func GetNIdFromLink(ifname string) (uint8, error) {
	if len(ifname) == 0 {
		return 0, fmt.Errorf("Interface name not available.")
	}

	addrs, err := GetIPv4Addrs(ifname)
	if err != nil {
		return 0, err
	}

	if len(addrs) == 0 {
		return 0, fmt.Errorf("IP address not exist. %s", ifname)
	}

	ip := addrs[0].IP
	return uint8(ip[len(ip)-1]), nil
}
