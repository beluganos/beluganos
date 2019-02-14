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
	"io"
)

type Tables struct {
	workerTable   *WorkerTable
	oidMapTable   *OidMapTable
	ifMapTable    *IfMapTable
	trapMapTable  *TrapMapTable
	trapSinkTable *TrapSinkTable
}

func NewTables() *Tables {
	return &Tables{
		workerTable:   NewWorkerTable(),
		oidMapTable:   NewOidMapTable(),
		ifMapTable:    NewIfMapTable(),
		trapMapTable:  NewTrapMapTable(),
		trapSinkTable: NewTrapSinkTable(),
	}
}

func (t *Tables) WorkerTable() *WorkerTable {
	return t.workerTable
}

func (t *Tables) OidMapTable() *OidMapTable {
	return t.oidMapTable
}

func (t *Tables) IfMapTable() *IfMapTable {
	return t.ifMapTable
}

func (t *Tables) TrapMapTable() *TrapMapTable {
	return t.trapMapTable
}

func (t *Tables) TrapSinkTable() *TrapSinkTable {
	return t.trapSinkTable
}

func (t *Tables) WriteTo(w io.Writer) (sum int64, err error) {
	var n int64

	n, err = t.workerTable.WriteTo(w)
	sum += n
	if err != nil {
		return
	}

	n, err = t.oidMapTable.WriteTo(w)
	sum += n
	if err != nil {
		return
	}

	n, err = t.ifMapTable.WriteTo(w)
	sum += n
	if err != nil {
		return
	}

	n, err = t.trapMapTable.WriteTo(w)
	sum += n
	if err != nil {
		return
	}

	n, err = t.trapSinkTable.WriteTo(w)
	sum += n
	if err != nil {
		return
	}

	return
}
