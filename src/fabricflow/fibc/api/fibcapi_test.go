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

package fibcapi

import (
	"net"
	"testing"
)

func TestCompMaskedMAC4(t *testing.T) {
	base := HardwareAddrMulticast4
	mask := HardwareAddrMulticast4Mask

	var mac net.HardwareAddr

	// match
	mac, _ = net.ParseMAC("01:00:5e:00:00:00")
	if ok := CompMaskedMAC(mac, base, mask); !ok {
		t.Errorf("CompMaskedMAC must match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("01:00:5e:00:00:01")
	if ok := CompMaskedMAC(mac, base, mask); !ok {
		t.Errorf("CompMaskedMAC must match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("01:00:5e:7f:ff:ff")
	if ok := CompMaskedMAC(mac, base, mask); !ok {
		t.Errorf("CompMaskedMAC must match. %s %s/%s", mac, base, mask)
	}

	// not match
	mac, _ = net.ParseMAC("01:00:5e:80:00:00")
	if ok := CompMaskedMAC(mac, base, mask); ok {
		t.Errorf("CompMaskedMAC must not match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("00:00:5e:00:00:00")
	if ok := CompMaskedMAC(mac, base, mask); ok {
		t.Errorf("CompMaskedMAC must not match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("01:01:5e:00:00:00")
	if ok := CompMaskedMAC(mac, base, mask); ok {
		t.Errorf("CompMaskedMAC must not match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("01:00:5f:10:00:00")
	if ok := CompMaskedMAC(mac, base, mask); ok {
		t.Errorf("CompMaskedMAC must not match. %s %s/%s", mac, base, mask)
	}
}

func TestCompMaskedMAC6(t *testing.T) {
	base := HardwareAddrMulticast6
	mask := HardwareAddrMulticast6Mask

	var mac net.HardwareAddr

	// match
	mac, _ = net.ParseMAC("33:33:00:00:00:00")
	if ok := CompMaskedMAC(mac, base, mask); !ok {
		t.Errorf("CompMaskedMAC must match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("33:33:00:00:00:01")
	if ok := CompMaskedMAC(mac, base, mask); !ok {
		t.Errorf("CompMaskedMAC must match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("33:33:ff:ff:ff:ff")
	if ok := CompMaskedMAC(mac, base, mask); !ok {
		t.Errorf("CompMaskedMAC must match. %s %s/%s", mac, base, mask)
	}

	// not match
	mac, _ = net.ParseMAC("33:34:00:00:00:00")
	if ok := CompMaskedMAC(mac, base, mask); ok {
		t.Errorf("CompMaskedMAC must not match. %s %s/%s", mac, base, mask)
	}
	mac, _ = net.ParseMAC("23:33:00:00:00:00")
	if ok := CompMaskedMAC(mac, base, mask); ok {
		t.Errorf("CompMaskedMAC must not match. %s %s/%s", mac, base, mask)
	}
}
