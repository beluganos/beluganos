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

const (
	ETHTYPE_IPV4    = 0x0800
	ETHTYPE_IPV6    = 0x86dd
	ETHTYPE_MPLS    = 0x8847
	ETHTYPE_LACP    = 0x8809
	ETHTYPE_ARP     = 0x0806
	ETHTYPE_VLAN_Q  = 0x8100
	ETHTYPE_VLAN_AD = 0x88a8
)

const (
	HWADDR_MULTICAST4       = "01:00:5e:00:00:00"
	HWADDR_MULTICAST4_MASK  = "ff:ff:ff:80:00:00"
	HWADDR_MULTICAST4_MATCH = HWADDR_MULTICAST4 + "/" + HWADDR_MULTICAST4_MASK
	HWADDR_MULTICAST6       = "33:33:00:00:00:00"
	HWADDR_MULTICAST6_MASK  = "ff:ff:00:00:00:00"
	HWADDR_MULTICAST6_MATCH = HWADDR_MULTICAST6 + "/" + HWADDR_MULTICAST6_MASK
)

const (
	OFPVID_UNTAGGED = 0x0001
	OFPVID_PRESENT  = 0x1000
	OFPVID_NONE     = 0x0000
	OFPVID_ABSENT   = 0x0000
)

const (
	IPPROTO_ICMP4 = 1
	IPPROTO_TCP   = 6
	IPPROTO_UDP   = 11
	IPPROTO_ICMP6 = 58
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
