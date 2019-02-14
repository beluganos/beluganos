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
	"sort"
)

//
// PortStats is stats in port stats file.
//
type PortStats map[string]interface{}

//
// PortNo returns port_no.
//
func (p PortStats) PortNo() (int, bool) {
	portNo, ok := p["port_no"]
	if !ok {
		return 0, false
	}

	switch portNo.(type) {
	case int:
		return portNo.(int), true
	case string:
		return 0, false
	default:
		return 0, false

	}
}

//
// PortStatsList is array of PortStats.
//
type PortStatsList []PortStats

//
// GetNext returns next port stats.
//
func (p PortStatsList) GetNext(portNo int) (PortStats, bool) {
	for _, ps := range p {
		if n, ok := ps.PortNo(); n >= portNo && ok {
			return ps, true
		}
	}
	return nil, false
}

//
// Get returns port stats.
//
func (p PortStatsList) Get(portNo int) (PortStats, bool) {
	for _, ps := range p {
		if n, ok := ps.PortNo(); n == portNo && ok {
			return ps, true
		}
	}
	return nil, false
}

//
// Validate removes invalid port stats.
//
func (p PortStatsList) Validate() PortStatsList {
	psList := PortStatsList{}
	for _, ps := range p {
		if _, ok := ps.PortNo(); ok {
			psList = append(psList, ps)
		}
	}
	return psList
}

//
// Sort sorts by port_no.
//
func (p PortStatsList) Sort() {
	sort.Slice(p, func(i, j int) bool {
		ii, _ := p[i].PortNo()
		jj, _ := p[j].PortNo()
		return ii < jj
	})
}
