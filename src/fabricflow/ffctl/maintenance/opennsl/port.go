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

type Port struct {
}

func NewPort() *Port {
	return &Port{}
}

func (p *Port) Name() string {
	return "port"
}

func (p *Port) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetPortInfos(context.Background(), api.NewGetPortInfosRequest())
	if err != nil {
		return err
	}

	for _, pinfo := range reply.PortInfos {
		fmt.Fprintf(w, "PortInfo; %d status:%d untag:%d\n",
			pinfo.Port,
			pinfo.LinkStatus,
			pinfo.UntaggedVlan,
		)
	}

	return nil
}
