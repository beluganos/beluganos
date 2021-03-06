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

type L2Addr struct {
}

func NewL2Addr() *L2Addr {
	return &L2Addr{}
}

func (a *L2Addr) Name() string {
	return "l2-addr"
}

func (a *L2Addr) Dump(w io.Writer, client api.GoNSLApiClient) error {
	reply, err := client.GetL2Addrs(context.Background(), api.NewGetL2AddrsRequest())
	if err != nil {
		return err
	}

	for _, l2addr := range reply.Addrs {
		fmt.Fprintf(w, "L2Addr: flags:%08x mac:%s vid:%d, port:%d\n",
			l2addr.GetFlags(), l2addr.GetMac(), l2addr.GetVid(), l2addr.GetPort())
	}

	return nil
}
