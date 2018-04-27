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
	"net"
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

func IpToString(ip net.IP) string {
	if ip == nil || len(ip) == 0 {
		return ""
	}
	return ip.String()
}
