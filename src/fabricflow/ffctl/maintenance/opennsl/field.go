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

package opennsl

import (
	"context"
	"fmt"
	api "gonsl/api"
	"io"
)

type FieldEntry struct {
}

func NewFieldEntry() *FieldEntry {
	return &FieldEntry{}
}

func (e *FieldEntry) Name() string {
	return "field-entry"
}

func (e *FieldEntry) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetFieldEntries(context.Background(), api.NewGetFieldEntriesRequest())
	if err != nil {
		return err
	}

	for _, entry := range reply.Entries {
		e.dumpEntry(w, entry)
	}

	return nil
}

func (e *FieldEntry) dumpEntry(w io.Writer, entry *api.FieldEntry) {
	switch entry.GetEntryType() {
	case api.FieldEntry_ETH_TYPE:
		e := entry.GetEthType()
		fmt.Fprintf(w, "FieldEntry: in_port:%d eth_type:0x%04x\n", e.GetInPort(), e.GetEthType())

	case api.FieldEntry_DST_IP:
		e := entry.GetDstIp()
		fmt.Fprintf(w, "FieldEntry: in_port:%d eth_type:0x%04x ip_dst:%s\n", e.GetInPort(), e.GetEthType(), e.GetIpDst())

	case api.FieldEntry_IP_PROTO:
		e := entry.GetIpProto()
		fmt.Fprintf(w, "FieldEntry: in_port:%d eth_type:0x%04x ip_proto:%d\n", e.GetInPort(), e.GetEthType(), e.GetIpProto())

	case api.FieldEntry_ETH_DST:
		e := entry.GetEthDst()
		fmt.Fprintf(w, "FieldEntry: in_port:%d eth_Dst:%s/%s", e.GetInPort(), e.GetEthDst(), e.GetEthMask())

	default:
		fmt.Fprintf(w, "Invalid EntryType. %d\n", entry.GetEntryType())
	}
}
