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

package main

import (
	"testing"
)

func TestPortStatsList_GetNext(t *testing.T) {
	ps1 := PortStats{"port_no": 1}
	ps2 := PortStats{"port_no": 2}
	ps4 := PortStats{"port_no": 4}
	pslist := PortStatsList{ps1, ps2, ps4}

	var ps PortStats
	var b bool

	ps, b = pslist.GetNext(0)
	if !b {
		t.Errorf("PortStatsList.GetNext error. %v", ps)
	}
	if n, _ := ps.PortNo(); n != 1 {
		t.Errorf("PortStatsList.GetNext unmatch. %v %d", ps, n)
	}

	ps, b = pslist.GetNext(1)
	if !b {
		t.Errorf("PortStatsList.GetNext error. %v", ps)
	}
	if n, _ := ps.PortNo(); n != 1 {
		t.Errorf("PortStatsList.GetNext unmatch. %v %d", ps, n)
	}

	ps, b = pslist.GetNext(2)
	if !b {
		t.Errorf("PortStatsList.GetNext error. %v", ps)
	}
	if n, _ := ps.PortNo(); n != 2 {
		t.Errorf("PortStatsList.GetNext unmatch. %v %d", ps, n)
	}

	ps, b = pslist.GetNext(3)
	if !b {
		t.Errorf("PortStatsList.GetNext error. %v", ps)
	}
	if n, _ := ps.PortNo(); n != 4 {
		t.Errorf("PortStatsList.GetNext unmatch. %v %d", ps, n)
	}

	ps, b = pslist.GetNext(4)
	if !b {
		t.Errorf("PortStatsList.GetNext error. %v", ps)
	}
	if n, _ := ps.PortNo(); n != 4 {
		t.Errorf("PortStatsList.GetNext unmatch. %v %d", ps, n)
	}

	_, b = pslist.GetNext(5)
	if b {
		t.Errorf("PortStatsList.GetNext must be error. %v", ps)
	}
}

func TestPortStatsList_Get(t *testing.T) {
	ps1 := PortStats{"port_no": 1}
	ps2 := PortStats{"port_no": 2}
	ps4 := PortStats{"port_no": 4}
	pslist := PortStatsList{ps1, ps2, ps4}

	var ps PortStats
	var b bool

	_, b = pslist.Get(0)
	if b {
		t.Errorf("PortStatsList.GetNext must be error. %v", ps)
	}

	ps, b = pslist.Get(1)
	if !b {
		t.Errorf("PortStatsList.GetNext error. %v", ps)
	}
	if n, _ := ps.PortNo(); n != 1 {
		t.Errorf("PortStatsList.GetNext unmatch. %v %d", ps, n)
	}

	ps, b = pslist.Get(2)
	if !b {
		t.Errorf("PortStatsList.GetNext error. %v", ps)
	}
	if n, _ := ps.PortNo(); n != 2 {
		t.Errorf("PortStatsList.GetNext unmatch. %v %d", ps, n)
	}

	_, b = pslist.Get(3)
	if b {
		t.Errorf("PortStatsList.GetNext mst be error. %v", ps)
	}
}

func TestPortStatsList_Validate(t *testing.T) {

	ps1 := PortStats{"port_no": 1}
	ps2 := PortStats{"port_n_": 2}
	ps3 := PortStats{"port_no": 3}
	pslist := PortStatsList{ps1, ps2, ps3}

	pslist = pslist.Validate()

	if len(pslist) != 2 {
		t.Errorf("PortStatsList.Validate unmatch. %v", pslist)
	}

	if ps := pslist[0]; ps["port_no"] != 1 {
		t.Errorf("PortStatsList.Validate unmatch. %v", ps)
	}
	if ps := pslist[1]; ps["port_no"] != 3 {
		t.Errorf("PortStatsList.Validate unmatch. %v", ps)
	}
}

func TestPortStatsList_Sort(t *testing.T) {

	ps1 := PortStats{"port_no": 1}
	ps2 := PortStats{"port_no": 2}
	ps3 := PortStats{"port_no": 3}
	pslist := PortStatsList{ps3, ps1, ps2}

	pslist.Sort()

	if len(pslist) != 3 {
		t.Errorf("PortStatsList.Validate unmatch. %v", pslist)
	}

	if ps := pslist[0]; ps["port_no"] != 1 {
		t.Errorf("PortStatsList.Validate unmatch. %v", ps)
	}
	if ps := pslist[1]; ps["port_no"] != 2 {
		t.Errorf("PortStatsList.Validate unmatch. %v", ps)
	}
	if ps := pslist[2]; ps["port_no"] != 3 {
		t.Errorf("PortStatsList.Validate unmatch. %v", ps)
	}
}
