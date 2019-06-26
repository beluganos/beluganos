// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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
	"testing"

	"github.com/beluganos/go-opennsl/opennsl"
)

func TestL3IfaceIDMap(t *testing.T) {
	m := L3IfaceIDMap{}

	if v := m.Register(10, 100, 1000); v != true {
		t.Errorf("L3IfaceIDMap Regiser error.")
	}

	if v, b := m.Get(10, 100); v != 1000 || !b {
		t.Errorf("L3IfaceIDMap Get error. %d %t", v, b)
	}

	if _, b := m.Get(10, 101); b {
		t.Errorf("L3IfaceIDMap Get must be error. %t", b)
	}

	if _, b := m.Get(11, 100); b {
		t.Errorf("L3IfaceIDMap Get must be error. %t", b)
	}

	m.Unregister(10, 100)

	if _, b := m.Get(10, 100); b {
		t.Errorf("L3IfaceIDMap Get must be error. %t", b)
	}
}

func TestL3IfaceIDMapTraverse(t *testing.T) {
	m := L3IfaceIDMap{}
	m.Register(10, 100, 1000)
	m.Register(10, 101, 1001)
	m.Register(11, 100, 1100)

	ids := map[L3IfaceIDKey]opennsl.L3IfaceID{}
	m.Traverse(func(key L3IfaceIDKey, ifaceId opennsl.L3IfaceID) bool {
		ids[key] = ifaceId
		return true
	})

	if v := len(ids); v != 3 {
		t.Errorf("L3IfaceIDMap Traverse error. count=%d", v)
	}

	if _, v := ids[NewL3IfaceIDKey(10, 100)]; !v {
		t.Errorf("L3IfaceIDMap Traverse error.")
	}
	if _, v := ids[NewL3IfaceIDKey(10, 101)]; !v {
		t.Errorf("L3IfaceIDMap Traverse error.")
	}
	if _, v := ids[NewL3IfaceIDKey(11, 100)]; !v {
		t.Errorf("L3IfaceIDMap Traverse error.")
	}
}
