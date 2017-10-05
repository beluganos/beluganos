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

func TestIPGenerator_NextIP(t *testing.T) {
	_, nw, _ := net.ParseCIDR("10.0.0.3/30")
	ipgen := NewIPGenerator(nw)

	// 10.0.0.1
	ip, err := ipgen.NextIP()
	if err != nil {
		t.Errorf("IPGenerator.NextIP error. %s", err)
	}
	if !ip.Equal(net.IPv4(10, 0, 0, 1)) {
		t.Errorf("IPGenerator.NextIP unmatch. %s", ip)
	}

	// 10.0.0.2
	ip, err = ipgen.NextIP()
	if err != nil {
		t.Errorf("IPGenerator.NextIP error. %s", err)
	}
	if !ip.Equal(net.IPv4(10, 0, 0, 2)) {
		t.Errorf("IPGenerator.NextIP unmatch. %s", ip)
	}

	// 10.0.0.3
	ip, err = ipgen.NextIP()
	if err == nil {
		t.Errorf("IPGenerator.NextIP must be error. %s", err)
	}

	// 10.0.0.4
	ip, err = ipgen.NextIP()
	if err == nil {
		t.Errorf("IPGenerator.NextIP must be error. %s", err)
	}
}
