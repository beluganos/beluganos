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

package main

import (
	"fmt"
	api "gonsl/api"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func printPortInfo(pinfo *api.PortInfo) {
	fmt.Printf("PortInfo; %d status:%d untag:%d\n",
		pinfo.Port,
		pinfo.LinkStatus,
		pinfo.UntaggedVlan,
	)
}

func dumpPortInfos(client api.GoNSLApiClient) {
	pinfos, err := client.GetPortInfos(context.Background(), api.NewGetPortInfosRequest())
	if err != nil {
		log.Errorf("GetPortInfos error. %s", err)
		return
	}

	for _, pinfo := range pinfos.PortInfos {
		printPortInfo(pinfo)
	}
}
