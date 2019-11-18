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

package monitor

import (
	"sync"
	"time"
)

type IfaceStatsHistory struct {
	history     []*IfaceStatsDataSets
	historySize uint16
	interval    time.Duration
	mutex       sync.Mutex
}

func NewIfaceStatsHistory() *IfaceStatsHistory {
	return &IfaceStatsHistory{
		history:     []*IfaceStatsDataSets{},
		historySize: 5,
		interval:    time.Second,
	}
}

func (d *IfaceStatsHistory) count() uint16 {
	return uint16(len(d.history))
}

func (d *IfaceStatsHistory) SetInterval(interval time.Duration) {
	if interval < time.Second {
		interval = time.Second
	}
	d.interval = interval
}

func (d *IfaceStatsHistory) SetHistorySize(hs uint16) {
	if hs < 2 {
		hs = 2
	}
	d.historySize = hs
}

func (d *IfaceStatsHistory) last() *IfaceStatsDataSets {
	if n := d.count(); n < 1 {
		return nil
	} else {
		return d.history[n-1]
	}
}

func (d *IfaceStatsHistory) pop() *IfaceStatsDataSets {
	if n := d.count(); n == 0 {
		return nil
	}

	ds := d.history[0]
	d.history = d.history[1:]

	return ds
}

func (d *IfaceStatsHistory) add(ds *IfaceStatsDataSets) {
	d.history = append(d.history, ds)
}

func (d *IfaceStatsHistory) Add(ds *IfaceStatsDataSets) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	ds.SetDiff(d.last())
	ds.SetAverage(d.interval.Seconds(), d.history...)

	d.add(ds)
	if n := d.count(); n >= d.historySize {
		d.pop()
	}
}
