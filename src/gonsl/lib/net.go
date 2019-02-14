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

package gonslib

import (
	"net"

	"golang.org/x/sys/unix"
)

//
// IPToAF returns address-family of ip.
//
func IPToAF(ip net.IP) int {
	if ip.To4() != nil {
		return unix.AF_INET
	}

	if ip.To16() != nil {
		return unix.AF_INET6
	}

	return unix.AF_UNSPEC
}

//
// AFToEtherType returns ether-type of family.
//
func AFToEtherType(family int) uint16 {
	switch family {
	case unix.AF_INET:
		return unix.ETH_P_IP
	case unix.AF_INET6:
		return unix.ETH_P_IPV6
	default:
		return 0
	}
}

//
// IPToEtherType returns ether-type of ip.
//
func IPToEtherType(ip net.IP) uint16 {
	return AFToEtherType(IPToAF(ip))
}

func EtherTypeToLen(etherType uint16) int {
	switch etherType {
	case unix.ETH_P_IP:
		return 32
	case unix.ETH_P_IPV6:
		return 128
	default:
		return 0
	}
}
