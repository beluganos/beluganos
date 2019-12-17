// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package bridge

import (
	"net"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
)

func maskHardwareAddr(mac, mask net.HardwareAddr) net.HardwareAddr {
	b := []byte{}
	for index, s := range mac {
		b = append(b, s&mask[index])
	}

	return b
}

func cmpHrdwareAddr(src, dst net.HardwareAddr) bool {
	for index, b := range src {
		if dst[index] != b {
			return false
		}
	}

	return true
}

func matchHardwareAddr(mac, src, mask net.HardwareAddr) bool {
	s := maskHardwareAddr(mac, mask)
	return cmpHrdwareAddr(s, src)
}

func isMulticastHardwareAddr(mac net.HardwareAddr) bool {
	mcIPv4MAC, _ := net.ParseMAC("01:00:5E:00:00:00")
	mcIPv4Msk, _ := net.ParseMAC("FF:FF:FF:80:00:00")
	if match := matchHardwareAddr(mac, mcIPv4MAC, mcIPv4Msk); match {
		return true
	}

	mcIPv6MAC, _ := net.ParseMAC("33:33:00:00:00:00")
	mcIPv6Msk, _ := net.ParseMAC("Ff:FF:00:00:00:00")
	if match := matchHardwareAddr(mac, mcIPv6MAC, mcIPv6Msk); match {
		return true
	}

	fullMsk, _ := net.ParseMAC("ff:ff:ff:ff:ff:ff")
	macList := []string{
		"01:00:0C:CC:CC:CC",
		"01:00:0C:CC:CC:CD",
		"01:80:C2:00:00:00",
		"01:80:C2:00:00:01",
		"01:80:C2:00:00:02",
		"01:80:C2:00:00:03",
		"01:80:C2:00:00:08",
		"01:80:C2:00:00:0E",
		"01:80:C2:00:00:21",
	}
	for _, macstr := range macList {
		m, _ := net.ParseMAC(macstr)
		if match := matchHardwareAddr(mac, m, fullMsk); match {
			return true
		}
	}

	return false
}

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bridge",
		Short: "bridge command.",
	}

	rootCmd.AddCommand(
		bridgeVlanCmd(),
		bridgeFdbCmd(),
	)

	return rootCmd
}

func newBridge(ifname string) *netlink.Bridge {
	flanFiltering := true
	bridge := netlink.Bridge{}
	bridge.VlanFiltering = &flanFiltering
	bridge.Attrs().Name = ifname

	return &bridge
}

func newLink(name string) netlink.Link {
	dev := netlink.Device{}
	dev.Attrs().Name = name
	return &dev
}
