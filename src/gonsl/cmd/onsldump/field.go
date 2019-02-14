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
	api "gonsl/api"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func printFieldEntry(entry *api.FieldEntry) {
	switch entry.GetEntryType() {
	case api.FieldEntry_ETH_TYPE:
		e := entry.GetEthType()
		fmt.Printf("FieldEntry: eth_type:0x%04x\n", e.GetEthType())

	case api.FieldEntry_DST_IP:
		e := entry.GetDstIp()
		fmt.Printf("FieldEntry: eth_type=0x%04x ip_dst:%s\n", e.GetEthType(), e.GetIpDst())

	case api.FieldEntry_IP_PROTO:
		e := entry.GetIpProto()
		fmt.Printf("FieldEntry: eth_type:0x%04x ip_proto:%d\n", e.GetEthType(), e.GetIpProto())

	default:
		log.Errorf("Invalid EntryType. %d", entry.GetEntryType())
	}
}

func dumpFieldEntries(client api.GoNSLApiClient) {
	fieldEntries, err := client.GetFieldEntries(context.Background(), api.NewGetFieldEntriesRequest())
	if err != nil {
		log.Errorf("GetFieldEntries error. %s", err)
		return
	}

	for _, entry := range fieldEntries.Entries {
		printFieldEntry(entry)
	}
}
