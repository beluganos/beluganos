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

package govsw

import (
	"fmt"
)

type LinkStatus uint16

const (
	LinkStatusNone LinkStatus = iota
	LinkStatusUp
	LinkStatusDown
)

var linkStatus_names = map[LinkStatus]string{
	LinkStatusNone: "none",
	LinkStatusUp:   "up",
	LinkStatusDown: "down",
}

func (v LinkStatus) String() string {
	if s, ok := linkStatus_names[v]; ok {
		return s
	}

	return fmt.Sprintf("LinkStatus(%d)", v)
}

type LinkCmd uint16

const (
	LinkCmdNone LinkCmd = iota
	LinkCmdAdd
	LinkDmdDel
)

var linkCmd_names = map[LinkCmd]string{
	LinkCmdNone: "none",
	LinkCmdAdd:  "add",
	LinkDmdDel:  "del",
}

func (v LinkCmd) String() string {
	if s, ok := linkCmd_names[v]; ok {
		return s
	}

	return fmt.Sprintf("LinkCmd(%d)", v)
}

type DB struct {
	dpID uint64
	link *LinkDB
	name *NameDB
}

func NewDB() *DB {
	return &DB{
		link: NewLinkDB(),
		name: NewNameDB(),
	}
}

func (db *DB) DpID() uint64 {
	return db.dpID
}

func (db *DB) SetDpID(dpID uint64) {
	db.dpID = dpID
}

func (db *DB) Link() *LinkDB {
	return db.link
}

func (db *DB) Ifname() *NameDB {
	return db.name
}
