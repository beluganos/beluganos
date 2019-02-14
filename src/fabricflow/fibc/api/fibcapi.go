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

package fibcapi

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	ETHTYPE_IPV4    = unix.ETH_P_IP      // 0x0800
	ETHTYPE_IPV6    = unix.ETH_P_IPV6    // 0x86dd
	ETHTYPE_MPLS    = unix.ETH_P_MPLS_UC // 0x8847
	ETHTYPE_LACP    = unix.ETH_P_SLOW    // 0x8809
	ETHTYPE_ARP     = unix.ETH_P_ARP     // 0x0806
	ETHTYPE_VLAN_Q  = unix.ETH_P_8021Q   // 0x8100
	ETHTYPE_VLAN_AD = unix.ETH_P_8021AD  // 0x88a8
)

const (
	HWADDR_NONE             = "00:00:00:00:00:00"
	HWADDR_DUMMY            = "02:00:00:00:00:00"
	HWADDR_EXACT_MASK       = "ff:ff:ff:ff:ff:ff"
	HWADDR_MULTICAST4       = "01:00:5e:00:00:00"
	HWADDR_MULTICAST4_MASK  = "ff:ff:ff:80:00:00"
	HWADDR_MULTICAST4_MATCH = HWADDR_MULTICAST4 + "/" + HWADDR_MULTICAST4_MASK
	HWADDR_MULTICAST6       = "33:33:00:00:00:00"
	HWADDR_MULTICAST6_MASK  = "ff:ff:00:00:00:00"
	HWADDR_MULTICAST6_MATCH = HWADDR_MULTICAST6 + "/" + HWADDR_MULTICAST6_MASK
	HWADDR_ISIS_LEVEL1      = "01:80:C2:00:00:14"
	HWADDR_ISIS_LEVEL2      = "01:80:C2:00:00:15"
)

var (
	HardwareAddrNone           = net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	HardwareAddrDummy          = net.HardwareAddr{0x02, 0x00, 0x00, 0x00, 0x00, 0x00}
	HardwareAddrExactMask      = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	HardwareAddrMulticast4     = net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x00}
	HardwareAddrMulticast4Mask = net.HardwareAddr{0xff, 0xff, 0xff, 0x80, 0x00, 0x00}
	HardwareAddrMulticast6     = net.HardwareAddr{0x33, 0x33, 0x00, 0x00, 0x00, 0x00}
	HardwareAddrMulticast6Mask = net.HardwareAddr{0xff, 0xff, 0x00, 0x00, 0x00, 0x00}
	HardwareAddrISISLevel1     = net.HardwareAddr{0x01, 0x80, 0xC2, 0x00, 0x00, 0x14}
	HardwareAddrISISLevel2     = net.HardwareAddr{0x01, 0x80, 0xC2, 0x00, 0x00, 0x15}
)

const (
	OFPVID_UNTAGGED = 0x0001
	OFPVID_PRESENT  = 0x1000
	OFPVID_NONE     = 0x0000
	OFPVID_ABSENT   = 0x0000
)

func AdjustVlanVID16(vlanId uint16) uint16 {
	if vlanId != OFPVID_NONE {
		return vlanId
	}
	return OFPVID_UNTAGGED
}

func AdjustVlanVID(vlanId uint16) uint32 {
	return uint32(AdjustVlanVID16(vlanId))
}

const (
	IPPROTO_ICMP4 = unix.IPPROTO_ICMP   // 1
	IPPROTO_TCP   = unix.IPPROTO_TCP    // 6
	IPPROTO_UDP   = unix.IPPROTO_UDP    // 17
	IPPROTO_ICMP6 = unix.IPPROTO_ICMPV6 // 58
	IPPROTO_OSPF  = 89
)

const (
	TCPPORT_BGP = 179
	TCPPORT_LDP = 646
)

const (
	MCADDR_ALLHOSTS   = "224.0.0.1"
	MCADDR_ALLROUTERS = "224.0.0.2"
	MCADDR_OSPF_HELLO = "224.0.0.5"
	MCADDR_OSPF_ALLDR = "224.0.0.6"
)

const (
	PRIORITY_DEFAULT = 0
	PRIORITY_NORMAL  = 32800
	PRIORITY_HIGHEST = 65530
)

const (
	MPLSTYPE_NONE      = 0x00
	MPLSTYPE_VPS       = 0x01
	MPLSTYPE_UNICAST   = 0x08
	MPLSTYPE_MULTICAST = 0x10
	MPLSTYPE_PHP       = 0x20
)

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

func ParseMaskedMAC(mac string) (net.HardwareAddr, net.HardwareAddr, error) {
	items := strings.SplitN(mac, "/", 2)
	addr, mask, err := func() (string, string, error) {
		switch len(items) {
		case 1:
			return items[0], HWADDR_EXACT_MASK, nil
		case 2:
			return items[0], items[1], nil
		default:
			return "", "", fmt.Errorf("Invalid MAC %s", mac)
		}
	}()
	if err != nil {
		return nil, nil, err
	}

	hwAddr, err := net.ParseMAC(addr)
	if err != nil {
		return nil, nil, err
	}

	hwMask, err := net.ParseMAC(mask)
	if err != nil {
		return nil, nil, err
	}

	return hwAddr, hwMask, nil
}

func NewMaskedMAC(mac, mask string) string {
	if len(mask) == 0 {
		mask = HWADDR_EXACT_MASK
	}
	return fmt.Sprintf("%s/%s", mac, mask)
}

func CompMaskedMAC(mac, base, mask net.HardwareAddr) bool {
	for index := 0; index < 6; index++ {
		if v := (mac[index] & mask[index]); v != base[index] {
			return false
		}
	}
	return true
}

//
// HardwareAddr types
//
type HardwareAddrType uint32

const (
	HWADDR_TYPE_NONE           = HardwareAddrType(0)
	HWADDR_TYPE_IPV4           = HardwareAddrType(1 << 0)
	HWADDR_TYPE_IPV6           = HardwareAddrType(1 << 1)
	HWADDR_TYPE_UNICAST        = HardwareAddrType(1 << 2)
	HWADDR_TYPE_MULTICAST      = HardwareAddrType(1 << 3)
	HWADDR_TYPE_OTHERS         = HardwareAddrType(1 << 4)
	HWADDR_TYPE_UNICAST_IPV4   = HWADDR_TYPE_UNICAST | HWADDR_TYPE_IPV4
	HWADDR_TYPE_UNICAST_IPV6   = HWADDR_TYPE_UNICAST | HWADDR_TYPE_IPV6
	HWADDR_TYPE_MULTICAST_IPV4 = HWADDR_TYPE_MULTICAST | HWADDR_TYPE_IPV4
	HWADDR_TYPE_MULTICAST_IPV6 = HWADDR_TYPE_MULTICAST | HWADDR_TYPE_IPV6
)

var hardwareAddrType_names = map[HardwareAddrType]string{
	HWADDR_TYPE_NONE:      "none",
	HWADDR_TYPE_IPV4:      "ipv4",
	HWADDR_TYPE_IPV6:      "ipv6",
	HWADDR_TYPE_UNICAST:   "unicast",
	HWADDR_TYPE_MULTICAST: "multicast",
	HWADDR_TYPE_OTHERS:    "others",
}

func (t HardwareAddrType) String() string {
	if t == HWADDR_TYPE_NONE {
		return "none"
	}

	names := []string{}
	for v, name := range hardwareAddrType_names {
		if t.Has(v) {
			names = append(names, name)
		}
	}
	return strings.Join(names, "|")
}

func (t HardwareAddrType) Has(v HardwareAddrType) bool {
	return (t & v) != 0
}

func ParseHardwareAddrType(mac net.HardwareAddr) (HardwareAddrType, error) {
	if len(mac) == 0 {
		return HWADDR_TYPE_NONE, fmt.Errorf("Invalid MAC. '%s'", mac)
	}

	if len(mac) != 6 {
		return HWADDR_TYPE_OTHERS, nil
	}

	if ok := CompMaskedMAC(mac, HardwareAddrMulticast4, HardwareAddrMulticast4Mask); ok {
		return HWADDR_TYPE_MULTICAST_IPV4, nil
	}

	if ok := CompMaskedMAC(mac, HardwareAddrMulticast6, HardwareAddrMulticast6Mask); ok {
		return HWADDR_TYPE_MULTICAST_IPV6, nil
	}

	return HWADDR_TYPE_UNICAST, nil
}

//
// TunnelType
//
var tunneType_native_names = map[TunnelType_Type]string{
	TunnelType_NOP:  "",
	TunnelType_IPIP: "ipip",
	TunnelType_IPV6: "ip6tnl",
	TunnelType_GRE4: "gre",
	TunnelType_GRE6: "ip6gre",
}

var tunnelType_native_values = map[string]TunnelType_Type{
	"ipip":   TunnelType_IPIP,
	"ip6tnl": TunnelType_IPV6,
	"gre":    TunnelType_GRE4,
	"ip6gre": TunnelType_GRE6,
}

func ParseTunnelTypeFromNative(s string) (TunnelType_Type, error) {
	if v, ok := tunnelType_native_values[s]; ok {
		return v, nil
	}
	return TunnelType_NOP, fmt.Errorf("Invalid TunnelType. %s", s)
}
