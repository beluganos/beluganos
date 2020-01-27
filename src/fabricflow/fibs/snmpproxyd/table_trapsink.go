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
	"fmt"
	"io"
	"net"
)

type TrapSinkTable struct {
	entries []string
}

func NewTrapSinkTable() *TrapSinkTable {
	return &TrapSinkTable{
		entries: []string{},
	}
}

func (t *TrapSinkTable) Add(sink string) {
	t.entries = append(t.entries, sink)
}

func (t *TrapSinkTable) GetAll(f func(*net.UDPAddr)) {
	for _, sink := range t.entries {
		if addr, err := net.ResolveUDPAddr("udp", sink); err == nil {
			f(addr)
		}
	}
}

func (t *TrapSinkTable) WriteTo(w io.Writer) (sum int64, err error) {
	for _, e := range t.entries {
		var n int
		n, err = fmt.Fprintf(w, "TrapSink: %s\n", e)
		sum += int64(n)
		if err != nil {
			return
		}
	}
	return
}
