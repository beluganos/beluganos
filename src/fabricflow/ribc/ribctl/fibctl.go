// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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

package ribctl

import (
	fibcnet "fabricflow/fibc/net"

	log "github.com/sirupsen/logrus"
)

type FIBCData struct {
	Hdr  *fibcnet.Header
	Data []byte
}

func NewFIBCData(h *fibcnet.Header, data []byte) *FIBCData {
	return &FIBCData{
		Hdr:  h,
		Data: data,
	}
}

type FIBController struct {
	*fibcnet.Client
	recvCh chan *FIBCData
	log    *log.Entry
}

func NewFIBController(addr string) *FIBController {
	return &FIBController{
		Client: fibcnet.NewClient(addr),
		recvCh: make(chan *FIBCData),
		log:    log.WithFields(log.Fields{"module": "FIBCController"}),
	}
}

func (f *FIBController) Send(msg fibcnet.Message, xid uint32) error {
	return f.Client.Write(msg, xid)
}

func (f *FIBController) Recv() <-chan *FIBCData {
	return f.recvCh
}

func (f *FIBController) Start() {
	go f.Client.Start(func(client *fibcnet.Client) {
		f.log.Infof("Monitor: START.")

		for {
			hdr, data, err := client.Read()
			if err != nil {
				f.log.Errorf("Monitor: EXIT. Read error. %s", err)
				return
			}

			f.recvCh <- NewFIBCData(hdr, data)
		}
	})
}
