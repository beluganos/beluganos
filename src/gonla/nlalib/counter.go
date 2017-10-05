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

package nlalib

import (
	"sync/atomic"
)

const COUNTERS_MAX = 256

//
// Counter(uint16)
//
type Counter16 struct {
	Count uint32
}

func (c *Counter16) Next() uint16 {
	n := atomic.AddUint32(&c.Count, 1)
	return uint16(n & 0xffff)
}

type Counters16 struct {
	Counters [COUNTERS_MAX]Counter16
}

func NewCounters16() *Counters16 {
	return &Counters16{}
}

func (c *Counters16) Next(index uint8) uint16 {
	return c.Counters[index].Next()
}

//
// Counter(uint32)
//
type Counter32 struct {
	Count uint32
}

func (c *Counter32) Next() uint32 {
	return atomic.AddUint32(&c.Count, 1)
}

type Counters32 struct {
	Counters [COUNTERS_MAX]Counter32
}

func NewCounters32() *Counters32 {
	return &Counters32{}
}

func (c *Counters32) Next(index uint8) uint32 {
	return c.Counters[index].Next()
}
