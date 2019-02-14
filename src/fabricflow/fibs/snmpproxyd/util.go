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
	"github.com/PromonLogicalis/snmp"
	log "github.com/sirupsen/logrus"
)

func dumpSnmpPdu(pdu *snmp.Pdu) {
	log.Debugf("Identifier: %d", pdu.Identifier)
	log.Debugf("ErrStatus : %d", pdu.ErrorStatus)
	log.Debugf("ErrIndex  : %d", pdu.ErrorIndex)
	log.Debugf("Variables : #%d", len(pdu.Variables))
	for _, v := range pdu.Variables {
		log.Debugf("Variable  : %s %v", v.Name, v.Value)
	}
}

func dumpSnmpBulkPdu(pdu *snmp.BulkPdu) {
	log.Debugf("Identifier    : %d", pdu.Identifier)
	log.Debugf("NonRepeaters  : %d", pdu.NonRepeaters)
	log.Debugf("MaxRepetitions: %d", pdu.MaxRepetitions)
	log.Debugf("Variables     : #%d", len(pdu.Variables))
	for _, v := range pdu.Variables {
		log.Debugf("Variable      : %s %v", v.Name, v.Value)
	}
}

func dumpSnmpTrapV1Pdu(pdu *snmp.V1TrapPdu) {
	log.Debugf("Enterprise  : %s", pdu.Enterprise)
	log.Debugf("AgentAddr   : %s", pdu.AgentAddr)
	log.Debugf("GenericTrap : %d", pdu.GenericTrap)
	log.Debugf("SpecificTrap: %d", pdu.SpecificTrap)
	log.Debugf("Timestamp   : %v", pdu.Timestamp)
	log.Debugf("Variables   : #%d", len(pdu.Variables))
	for _, v := range pdu.Variables {
		log.Debugf("Variable    : %s %v", v.Name, v.Value)
	}
}
