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

func printIDMapEntry(entry *api.IDMapEntry) {
	fmt.Printf("IDMap: %-8s[%s] = %d\n", entry.Name, entry.Key, entry.Value)
}

func dumpIDMapEntries(client api.GoNSLApiClient) {
	reply, err := client.GetIDMapEntries(context.Background(), api.NewGetIDMapEntriesRequest())
	if err != nil {
		log.Errorf("GetL2Addrs error. %s", err)
		return
	}

	for _, entry := range reply.Entries {
		printIDMapEntry(entry)
	}
}
