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
	"testing"
)

func TestIncIP(t *testing.T) {
	ip := net.ParseIP("10.0.0.0")
	IncIP(ip)
	if s := ip.String(); s != "10.0.0.1" {
		t.Errorf("IncIP unmatch. %v", ip)
	}

	ip = net.ParseIP("10.0.0.255")
	IncIP(ip)
	if s := ip.String(); s != "10.0.1.0" {
		t.Errorf("IncIP unmatch. %v", ip)
	}

	ip = net.ParseIP("255.255.255.255")
	IncIP(ip)
	if s := ip.String(); s != "0.0.0.0" {
		t.Errorf("IncIP unmatch. %v", ip)
	}
}

func TestToBroadcast(t *testing.T) {
	_, nw, _ := net.ParseCIDR("10.0.1.1/30")
	bc := ToBroadcast(nw)

	if !bc.Equal(net.IPv4(10, 0, 1, 3)) {
		t.Errorf("ToBroadcast unmatch. %v", bc)
	}

	_, nw, _ = net.ParseCIDR("10.0.1.1/24")
	bc = ToBroadcast(nw)

	if !bc.Equal(net.IPv4(10, 0, 1, 255)) {
		t.Errorf("ToBroadcast unmatch. %v", bc)
	}
}
