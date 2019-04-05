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

package fibslib

import (
	"strconv"
	"strings"
)

const (
	SNMP_LISTEN_ADDR = "127.0.0.1:161"
	SNMP_DAEMON_ADDR = "127.0.0.1:8161"
	SNMP_COMMUNITY   = "beluganos-internal"
	SNMP_VERSION     = 1 // v2c
	SNMP_OID_IFACES  = ".1.2.3.4.5"
)

// SnmpType is snmp value type.
type SnmpType string

const (
	// SnmpTypeInteger is INTEGER type.
	SnmpTypeInteger = "integer"
	// SnmpTypeString is STRING type.
	SnmpTypeString = "string"
)

const (
	SNMP_OID_ifNumber                   = ".1.3.6.1.2.1.2.1.0"
	SNMP_OID_ifIndex                    = ".1.3.6.1.2.1.2.2.1.1"
	SNMP_OID_ifDescr                    = ".1.3.6.1.2.1.2.2.1.2"
	SNMP_OID_ifType                     = ".1.3.6.1.2.1.2.2.1.3"
	SNMP_OID_ifMtu                      = ".1.3.6.1.2.1.2.2.1.4"
	SNMP_OID_ifSpeed                    = ".1.3.6.1.2.1.2.2.1.5"
	SNMP_OID_ifPhysAddress              = ".1.3.6.1.2.1.2.2.1.6"
	SNMP_OID_ifAdminStatus              = ".1.3.6.1.2.1.2.2.1.7"
	SNMP_OID_ifOperStatus               = ".1.3.6.1.2.1.2.2.1.8"
	SNMP_OID_ifLastChange               = ".1.3.6.1.2.1.2.2.1.9"
	SNMP_OID_ifInOctets                 = ".1.3.6.1.2.1.2.2.1.10"
	SNMP_OID_ifInUcastPkts              = ".1.3.6.1.2.1.2.2.1.11"
	SNMP_OID_ifInNUcastPkts             = ".1.3.6.1.2.1.2.2.1.12"
	SNMP_OID_ifInDiscards               = ".1.3.6.1.2.1.2.2.1.13"
	SNMP_OID_ifInErrors                 = ".1.3.6.1.2.1.2.2.1.14"
	SNMP_OID_ifInUnknownProtos          = ".1.3.6.1.2.1.2.2.1.15"
	SNMP_OID_ifOutOctets                = ".1.3.6.1.2.1.2.2.1.16"
	SNMP_OID_ifOutUcastPkts             = ".1.3.6.1.2.1.2.2.1.17"
	SNMP_OID_ifOutDiscards              = ".1.3.6.1.2.1.2.2.1.19"
	SNMP_OID_ifOutErrors                = ".1.3.6.1.2.1.2.2.1.20"
	SNMP_OID_ifOutQLen                  = ".1.3.6.1.2.1.2.2.1.21"
	SNMP_OID_ifSpecific                 = ".1.3.6.1.2.1.2.2.1.22"
	SNMP_OID_ifName                     = ".1.3.6.1.2.1.31.1.1.1.1"
	SNMP_OID_ifInMulticastPkts          = ".1.3.6.1.2.1.31.1.1.1.2"
	SNMP_OID_ifInBroadcastPkts          = ".1.3.6.1.2.1.31.1.1.1.3"
	SNMP_OID_ifOutMulticastPkts         = ".1.3.6.1.2.1.31.1.1.1.4"
	SNMP_OID_ifOutBroadcastPkts         = ".1.3.6.1.2.1.31.1.1.1.5"
	SNMP_OID_ifHCInOctets               = ".1.3.6.1.2.1.31.1.1.1.6"
	SNMP_OID_ifHCInUcastPkts            = ".1.3.6.1.2.1.31.1.1.1.7"
	SNMP_OID_ifHCInMulticastPkts        = ".1.3.6.1.2.1.31.1.1.1.8"
	SNMP_OID_ifHCInBroadcastPkts        = ".1.3.6.1.2.1.31.1.1.1.9"
	SNMP_OID_ifHCOutOctets              = ".1.3.6.1.2.1.31.1.1.1.10"
	SNMP_OID_ifHCOutUcastPkts           = ".1.3.6.1.2.1.31.1.1.1.11"
	SNMP_OID_ifHCOutMulticastPkts       = ".1.3.6.1.2.1.31.1.1.1.12"
	SNMP_OID_ifHCOutBroadcastPkts       = ".1.3.6.1.2.1.31.1.1.1.13"
	SNMP_OID_ifHighSpeed                = ".1.3.6.1.2.1.31.1.1.1.15"
	SNMP_OID_ifPromiscuousMode          = ".1.3.6.1.2.1.31.1.1.1.16"
	SNMP_OID_ifConnectorPresent         = ".1.3.6.1.2.1.31.1.1.1.17"
	SNMP_OID_ifAlias                    = ".1.3.6.1.2.1.31.1.1.1.18"
	SNMP_OID_ifCounterDiscontinuityTime = ".1.3.6.1.2.1.31.1.1.1.19"
)

//
// ParseOID parses string and returns []uint.
//
func ParseOID(oid string) []uint {
	items := strings.Split(oid, ".")
	oids := []uint{}
	for _, item := range items {
		if len(item) > 0 {
			v, _ := strconv.Atoi(item)
			oids = append(oids, uint(v))
		}
	}
	return oids
}

//
// DelPrefixOID remove prefix from oid.
//
func DelPrefixOID(oid []uint, prefix []uint) ([]uint, bool) {
	if len(oid) < len(prefix) {
		return oid, false
	}
	return oid[len(prefix):], true
}

//
// ReplaceOID replace prefix of oid by newOid
//
func ReplaceOID(oid []uint, prefix []uint, newOid []uint) ([]uint, bool) {
	remOid, ok := DelPrefixOID(oid, prefix)
	if !ok {
		return oid, false
	}
	newOid = CloneOID(newOid)
	return append(newOid, remOid...), true
}

//
// CompareOID compares oid.
//
func CompareOID(oid1, oid2 []uint) int {
	if diff := len(oid1) - len(oid2); diff != 0 {
		return diff
	}

	for index, v1 := range oid1 {
		v2 := oid2[index]
		if v1 < v2 {
			return -1
		}
		if v1 > v2 {
			return 1
		}
	}
	return 0
}

//
// CloneOID creates new oid instance.
//
func CloneOID(oid []uint) []uint {
	newOid := make([]uint, len(oid))
	copy(newOid, oid)
	return newOid
}
