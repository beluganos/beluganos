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

func printVlanEntry(vlan *api.VlanEntry) {
	fmt.Printf("Vlan: vid:%d ports:%v, untag:%v\n",
		vlan.Vid,
		vlan.Ports,
		vlan.UntagPorts,
	)
}

func dumpVlans(client api.GoNSLApiClient) {
	vlans, err := client.GetVlans(context.Background(), api.NewGetVlansRequest())
	if err != nil {
		log.Errorf("GetVlans error. %s", err)
		return
	}

	for _, vlan := range vlans.Vlans {
		printVlanEntry(vlan)
	}
}
