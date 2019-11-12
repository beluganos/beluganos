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

package fflibnet

import (
	"net"
)

func IncIP(ip net.IP) {
	endpos := 0
	if ip.To4() != nil {
		endpos = len(ip) - 4
	}
	for pos := len(ip) - 1; pos >= endpos; pos-- {
		ip[pos]++
		if ip[pos] > 0 {
			break
		}
	}
}

func IncIPNet(ipnet *net.IPNet) {
	ones, _ := ipnet.Mask.Size()

	if ones == 0 {
		return
	}

	index := uint32((ones - 1) / 8)
	byteShift := (index+1)*8 - uint32(ones)
	byteDiff := byte(1 << byteShift)

	ip := ipnet.IP

	for {
		ip[index] += byteDiff
		if ip[index] > 0 || index == 0 {
			break
		}

		index--
		byteDiff = 1
	}
}

func ToBroadcast(nw *net.IPNet) net.IP {
	bc := nw.IP
	for pos, mask := range nw.Mask {
		bc[pos] |= ^mask
	}
	return bc
}
