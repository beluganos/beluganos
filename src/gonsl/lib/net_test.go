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
	"testing"

	"golang.org/x/sys/unix"
)

func TestIPToAF(t *testing.T) {
	var ip net.IP

	ip = net.ParseIP("1.1.1.1")
	if af := IPToAF(ip); af != unix.AF_INET {
		t.Errorf("IPToAF unmatch. %s %d", ip, af)
	}

	ip = net.ParseIP("127.0.0.1")
	if af := IPToAF(ip); af != unix.AF_INET {
		t.Errorf("IPToAF unmatch. %s %d", ip, af)
	}

	ip = net.ParseIP("2001:db8::1")
	if af := IPToAF(ip); af != unix.AF_INET6 {
		t.Errorf("IPToAF unmatch. %s %d", ip, af)
	}

	ip = net.ParseIP("::1")
	if af := IPToAF(ip); af != unix.AF_INET6 {
		t.Errorf("IPToAF unmatch. %s %d", ip, af)
	}

	ip = net.ParseIP("")
	if af := IPToAF(ip); af != 0 {
		t.Errorf("IPToAF unmatch. %s %d", ip, af)
	}
}

func TestAFToEtherType(t *testing.T) {
	if e := AFToEtherType(unix.AF_INET); e != unix.ETH_P_IP {
		t.Errorf("AFToEtherType unmatch. %d", e)
	}

	if e := AFToEtherType(unix.AF_INET6); e != unix.ETH_P_IPV6 {
		t.Errorf("AFToEtherType unmatch. %d", e)
	}

	if e := AFToEtherType(123); e != 0 {
		t.Errorf("AFToEtherType unmatch. %d", e)
	}
}

func TestIPToEtherType(t *testing.T) {
	var ip net.IP

	ip = net.ParseIP("1.1.1.1")
	if e := IPToEtherType(ip); e != unix.ETH_P_IP {
		t.Errorf("IPToEtherType unmatch. %s %d", ip, e)
	}

	ip = net.ParseIP("127.0.0.1")
	if e := IPToEtherType(ip); e != unix.ETH_P_IP {
		t.Errorf("IPToEtherType unmatch. %s %d", ip, e)
	}

	ip = net.ParseIP("2001:db8::1")
	if e := IPToEtherType(ip); e != unix.ETH_P_IPV6 {
		t.Errorf("IPToEtherType unmatch. %s %d", ip, e)
	}

	ip = net.ParseIP("::1")
	if e := IPToEtherType(ip); e != unix.ETH_P_IPV6 {
		t.Errorf("IPToEtherType unmatch. %s %d", ip, e)
	}

	ip = net.ParseIP("")
	if e := IPToEtherType(ip); e != 0 {
		t.Errorf("IPToEtherType unmatch. %s %d", ip, e)
	}

}
