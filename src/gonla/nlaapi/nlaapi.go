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

package nlaapi

import (
	"net"
)

func BytesToIPs(bs [][]byte) []net.IP {
	ips := make([]net.IP, len(bs))
	for i, b := range bs {
		ips[i] = net.IP(b)
	}
	return ips
}

func IPsToBytes(ips []net.IP) [][]byte {
	bs := make([][]byte, len(ips))
	for i, ip := range ips {
		bs[i] = ip
	}
	return bs
}

func IPNetToBytes(ip *net.IPNet) ([]byte, []byte) {
	if ip == nil {
		return []byte{}, []byte{}
	}
	return ip.IP, ip.Mask
}

func BytesToIPNet(ip []byte, mask []byte) *net.IPNet {
	if len(ip) == 0 {
		return nil
	}

	return &net.IPNet{
		IP:   ip,
		Mask: mask,
	}
}
