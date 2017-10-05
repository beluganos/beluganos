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
	"fmt"
	"net"
)

type IPGenerator struct {
	network *net.IPNet
	baseIP  net.IP
	currIP  net.IP
	bcIP    net.IP
}

func NewIPGenerator(nw *net.IPNet) *IPGenerator {
	baseIP := nw.IP.Mask(nw.Mask)
	bcIP := ToBroadcast(nw)
	return &IPGenerator{
		network: nw,
		baseIP:  baseIP,
		currIP:  baseIP,
		bcIP:    bcIP,
	}
}

func (g *IPGenerator) IsBroadcast(ip net.IP) bool {
	return ip.Equal(g.bcIP)
}

func (g *IPGenerator) Reset() {
	g.currIP = g.baseIP
}

func (g *IPGenerator) NextIP() (net.IP, error) {
	IncIP(g.currIP)

	if ok := g.network.Contains(g.currIP); !ok {
		return nil, fmt.Errorf("IP Net Overflow")
	}

	if isBc := g.IsBroadcast(g.currIP); isBc {
		return nil, fmt.Errorf("IP Net Overflow")
	}

	nextIP := make([]byte, len(g.currIP))
	copy(nextIP, g.currIP)
	return nextIP, nil
}
