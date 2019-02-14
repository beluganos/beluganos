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
	"strings"
	"testing"

	lib "fabricflow/fibs/fibslib"
)

func TestSnmpReply_WriteTo_string(t *testing.T) {
	oid := ".1.2.3"
	data := ".1.2.3\nstring\ntest\n"
	sb := strings.Builder{}

	total, err := NewSnmpReply(oid, lib.SnmpTypeString, "test").WriteTo(&sb)

	if err != nil {
		t.Errorf("SnmpReply.WriteTo error. %s", err)
	}
	if total != 19 {
		t.Errorf("SnmpReply.WriteTo unmatch. total=%d", total)
	}

	s := sb.String()
	if s != data {
		t.Errorf("SnmpReply.WriteTo unmatch. s=%s", s)
	}
}

func TestSnmpReply_WriteTo_integer(t *testing.T) {
	oid := ".1.2.3"
	data := ".1.2.3\ninteger\n1024\n"
	sb := strings.Builder{}

	total, err := NewSnmpReply(oid, lib.SnmpTypeInteger, 1024).WriteTo(&sb)

	if err != nil {
		t.Errorf("SnmpReply.WriteTo error. %s", err)
	}
	if total != 20 {
		t.Errorf("SnmpReply.WriteTo unmatch. total=%d", total)
	}

	s := sb.String()
	if s != data {
		t.Errorf("SnmpReply.WriteTo unmatch. s=%s", s)
	}
}

func TestParseOID_0(t *testing.T) {
	oidStr := "."
	oid := lib.ParseOID(oidStr)

	if v := len(oid); v != 0 {
		t.Errorf("ParseOID unmatch. len=%v", v)
	}
}

func TestParseOID_1(t *testing.T) {
	oidStr := ".1"
	oid := lib.ParseOID(oidStr)

	if v := len(oid); v != 1 {
		t.Errorf("ParseOID unmatch. len=%v", v)
	}

	if oid[0] != 1 {
		t.Errorf("ParseOID unmatch. oid=%v", oid)
	}
}

func TestParseOID_n(t *testing.T) {
	oidStr := ".1.10.0"
	oid := lib.ParseOID(oidStr)

	if v := len(oid); v != 3 {
		t.Errorf("ParseOID unmatch. len=%v", v)
	}

	if oid[0] != 1 {
		t.Errorf("ParseOID unmatch. oid=%v", oid)
	}
	if oid[1] != 10 {
		t.Errorf("ParseOID unmatch. oid=%v", oid)
	}
	if oid[2] != 0 {
		t.Errorf("ParseOID unmatch. oid=%v", oid)
	}
}

func TestGetIndexOfOID(t *testing.T) {
	oid := []uint{10}
	index, ok := GetIndexOfOID(oid)

	if !ok {
		t.Errorf("GetIndexOfOID error. %v", oid)
	}
	if index != 10 {
		t.Errorf("GetIndexOfOID unmatch. %v", oid)
	}
}

func TestGetIndexOfOID_err(t *testing.T) {
	oid := []uint{1, 10}
	_, ok := GetIndexOfOID(oid)

	if ok {
		t.Errorf("GetIndexOfOID error. %v", oid)
	}
}

func TestGetNextIndexOfOID_0(t *testing.T) {
	oid := []uint{}
	index, ok := GetNextIndexOfOID(oid, 10)

	if !ok {
		t.Errorf("GetIndexOfOID error. %v", oid)
	}
	if index != 10 {
		t.Errorf("GetIndexOfOID unmatch. %v %d", oid, index)
	}
}

func TestGetNextIndexOfOID(t *testing.T) {
	oid := []uint{10}
	index, ok := GetNextIndexOfOID(oid, 0)

	if !ok {
		t.Errorf("GetIndexOfOID error. %v", oid)
	}
	if index != 11 {
		t.Errorf("GetIndexOfOID unmatch. %v %d", oid, index)
	}
}

func TestGetNextIndexOfOID_err(t *testing.T) {
	oid := []uint{1, 10}
	_, ok := GetNextIndexOfOID(oid, 9)

	if ok {
		t.Errorf("GetIndexOfOID error. %v", oid)
	}
}
