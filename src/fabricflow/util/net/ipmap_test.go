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
	"testing"
)

type IPMapTestGenerator struct {
	IPs []net.IP
}

func (g *IPMapTestGenerator) New(key net.IP) (net.IP, error) {
	if len(g.IPs) == 0 {
		return nil, fmt.Errorf("No Pool")
	}
	ip := g.IPs[0]
	g.IPs = g.IPs[1:]
	return ip, nil
}

func (g *IPMapTestGenerator) Free(k net.IP, v net.IP) {

}

func NewIPMapTestGenerator(ips []net.IP) IPMapGenerator {
	return &IPMapTestGenerator{
		IPs: ips,
	}
}

func TestIPMap_New(t *testing.T) {
	g := NewIPMapTestGenerator([]net.IP{})
	m := NewIPMap(g)
	if m == nil {
		t.Errorf("NewIPMap error. %v", m)
	}
}

func TestIPMap_Value(t *testing.T) {
	g := NewIPMapTestGenerator([]net.IP{net.IPv4(1, 1, 1, 1)})
	m := NewIPMap(g)

	ip, err := m.Value(net.IPv4(10, 1, 1, 1))
	if err != nil {
		t.Errorf("IPMap.Value error. %s", err)
	}
	if ip == nil {
		t.Errorf("IPMap.Value unmatch. %v", ip)
	}

	ip, err = m.Value(net.IPv4(10, 1, 1, 2))
	if err == nil {
		t.Errorf("IPMap.Value must be error. %s", err)
	}
	if ip != nil {
		t.Errorf("IPMap.Value unmatch. %v", ip)
	}

	ip, err = m.Value(net.IPv4(10, 1, 1, 1))
	if err != nil {
		t.Errorf("IPMap.Value error. %s", err)
	}
	if ip == nil {
		t.Errorf("IPMap.Value unmatch. %v", ip)
	}
}
