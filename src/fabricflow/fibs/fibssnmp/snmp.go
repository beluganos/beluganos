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

	lib "fabricflow/fibs/fibslib"
)

//
// SnmpReply is reply of get/getnext request.
//
type SnmpReply struct {
	Oid   string
	Type  lib.SnmpType
	Value interface{}
}

//
// NewSnmpReply returns new instance.
//
func NewSnmpReply(oid string, t lib.SnmpType, v interface{}) *SnmpReply {
	return &SnmpReply{
		Oid:   oid,
		Type:  t,
		Value: v,
	}
}

//
// WriteTo writes reply to writer.
//
func (r *SnmpReply) WriteTo(w io.Writer) (total int64, err error) {
	var writeSize int
	if r != nil {
		if writeSize, err = fmt.Fprintln(w, r.Oid); err != nil {
			goto EXIT
		}
		total += int64(writeSize)

		if writeSize, err = fmt.Fprintln(w, r.Type); err != nil {
			goto EXIT
		}
		total += int64(writeSize)

		if writeSize, err = fmt.Fprintln(w, r.Value); err != nil {
			goto EXIT
		}
		total += int64(writeSize)

	} else {
		writeSize, err = fmt.Fprintln(w, "NONE")
		total = int64(writeSize)
	}

EXIT:
	return
}

//
// GetIndexOfOID returns index of OID.
//
func GetIndexOfOID(oid []uint) (uint, bool) {
	if len(oid) == 1 {
		return oid[0], true
	}
	return 0, false
}

//
// GetNextIndexOfOID returns next index of OID.
//
func GetNextIndexOfOID(oid []uint, defaultIndex uint) (uint, bool) {
	switch len(oid) {
	case 0:
		return defaultIndex, true
	case 1:
		return oid[0] + 1, true
	default:
		return defaultIndex, false
	}
}
