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

package mkpb

const (
	SnmpTypeNone      = ""
	SnmpTypeInteger32 = "integer"
	SnmpTypeInteger64 = "integer64"
	SnmpTypeCounter32 = "counter"
	SnmpTypeCounter64 = "counter64"
	SnmpTypeGauge32   = "gauge"
	SnmpTypeString    = "string"
	SnmpTypeOID       = "objectid"
	SnmpTypeTimeTicks = "timeticks"
	SnmpTypeIPAddr    = "ipaddress"
)

const (
	SnmpModeNone  = ""
	SnmpModeProxy = "proxy"
)

type SnmpdOidEntry struct {
	Name     string
	SnmpType string
	Oid      []uint32
	Local    []uint32
	Proxy    string
}

func NewSnmpdOidEntry(name string, snmpType string, oid []uint32, local []uint32) *SnmpdOidEntry {
	return &SnmpdOidEntry{
		Name:     name,
		SnmpType: snmpType,
		Oid:      oid,
		Local:    local,
	}
}

func NewSnmpdProxyOidEntry(name string, oid []uint32, local []uint32) *SnmpdOidEntry {
	return &SnmpdOidEntry{
		Name:  name,
		Oid:   oid,
		Local: local,
	}
}

var snmpOidMap = []*SnmpdOidEntry{
	NewSnmpdOidEntry(
		"ifIndex",
		SnmpTypeInteger32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 1},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 1},
	),
	//NewSnmpdOidEntry(
	//	"ifDescr",
	//	SnmpTypeString,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 2},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 2},
	//),
	//NewSnmpdOidEntry(
	//	"ifType",
	//	SnmpTypeString,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 3},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 3},
	//),
	//NewSnmpdOidEntry(
	//	"ifMtu",
	//	SnmpTypeInteger32,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 4},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 4},
	//),
	//NewSnmpdOidEntry(
	//	"ifSpeed",
	//	SnmpTypeGauge32,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 5},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 5},
	//),
	//NewSnmpdOidEntry(
	//	SnmpTypeIPAddr,
	//	SnmpTypeString,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 6},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 6},
	//),
	//NewSnmpdOidEntry(
	//	"ifAdminStatus",
	//	SnmpTypeInteger32,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 7},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 7},
	//),
	NewSnmpdOidEntry(
		"ifOperStatus",
		SnmpTypeInteger32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 8},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 8},
	),
	//NewSnmpdOidEntry(
	//	"ifLastChange",
	//	SnmpTypeTimeTicks,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 9},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 9},
	//),
	NewSnmpdOidEntry(
		"ifInOctets",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 10},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 10},
	),
	NewSnmpdOidEntry(
		"ifInUcastPkts",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 11},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 11},
	),
	NewSnmpdOidEntry(
		"ifInNUcastPkts",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 12},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 12},
	),
	NewSnmpdOidEntry(
		"ifInDiscards",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 13},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 13},
	),
	NewSnmpdOidEntry(
		"ifInErrors",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 14},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 14},
	),
	NewSnmpdOidEntry(
		"ifInUnknownProtos",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 15},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 15},
	),
	NewSnmpdOidEntry(
		"ifOutOctets",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 16},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 16},
	),
	NewSnmpdOidEntry(
		"ifOutUcastPkts",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 17},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 17},
	),
	NewSnmpdOidEntry(
		"ifOutNUcastPkts",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 18},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 18},
	),
	NewSnmpdOidEntry(
		"ifOutDiscards",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 19},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 19},
	),
	NewSnmpdOidEntry(
		"ifOutErrors",
		SnmpTypeCounter32,
		[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 20},
		[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 20},
	),
	//NewSnmpdOidEntry(
	//	"ifOutQLen",
	//	SnmpTypeGauge32,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 21},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 21},
	//),
	//NewSnmpdOidEntry(
	//	"ifSpecific",
	//	SnmpTypeOID,
	//	[]uint32{1, 3, 6, 1, 2, 1, 2, 2, 1, 22},
	//	[]uint32{1, 3, 6, 1, 4, 99999, 2, 2, 1, 22},
	//),
	NewSnmpdOidEntry(
		"ifName",
		SnmpTypeString,
		[]uint32{1, 3, 6, 1, 2, 1, 31, 1, 1, 1, 1},
		[]uint32{1, 3, 6, 1, 4, 99999, 31, 1, 1, 1, 1},
	),
}

var snmpONLOidMap = []*SnmpdOidEntry{
	NewSnmpdProxyOidEntry(
		"proxy-to-wbsw",
		[]uint32{1, 3, 6, 1, 1234, 0, 1, 3},
		[]uint32{1, 3},
	),
	NewSnmpdProxyOidEntry(
		"ONL-mibs",
		[]uint32{1, 3, 6, 1, 4, 1, 42623, 1},
		[]uint32{1, 3, 6, 1, 4, 1, 42623, 1},
	),
}

func NewSnmpONLOidMap(addr string) []*SnmpdOidEntry {
	entries := []*SnmpdOidEntry{}
	for _, oidmap := range snmpONLOidMap {
		entry := *oidmap
		entry.Proxy = addr
		entries = append(entries, &entry)
	}
	return entries
}
