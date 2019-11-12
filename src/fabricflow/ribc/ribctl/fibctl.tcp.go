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
	fibcapi "fabricflow/fibc/api"
	fibcnet "fabricflow/fibc/net"

	log "github.com/sirupsen/logrus"
)

type FIBTcpController struct {
	*fibcnet.Client
	recvCh chan *fibcapi.VmMonitorReply
	log    *log.Entry
}

func NewFIBTcpController(addr string) *FIBTcpController {
	return &FIBTcpController{
		Client: fibcnet.NewClient(addr),
		recvCh: make(chan *fibcapi.VmMonitorReply),
		log:    log.WithFields(log.Fields{"module": "FIBTcpController"}),
	}
}

func (f *FIBTcpController) FIBCType() string {
	return FIBCTypeTCP
}

func (f *FIBTcpController) monitor() {
	f.Client.Start(func(client *fibcnet.Client) {
		f.log.Infof("Monitor: START.")

		for {
			hdr, data, err := client.Read()
			if err != nil {
				f.log.Errorf("Monitor: EXIT. Read error. %s", err)
				return
			}

			reply, err := fibcapi.ParseVmMonitorReply(hdr, data)
			if err != nil {
				f.log.Warnf("Monitor: Parse Monitor reply error. %s", err)
			} else {
				f.recvCh <- reply
			}
		}
	})

}

func (f *FIBTcpController) Start() error {
	go f.monitor()
	return nil
}

func (f *FIBTcpController) Recv() <-chan *fibcapi.VmMonitorReply {
	return f.recvCh
}

func (f *FIBTcpController) Hello(hello *fibcapi.Hello) error {
	return f.Client.Write(hello, 0)
}

func (f *FIBTcpController) PortConfig(pc *fibcapi.PortConfig) error {
	return f.Client.Write(pc, 0)
}

func (f *FIBTcpController) FlowMod(mod *fibcapi.FlowMod) error {
	return f.Client.Write(mod, 0)
}

func (f *FIBTcpController) GroupMod(mod *fibcapi.GroupMod) error {
	return f.Client.Write(mod, 0)
}
