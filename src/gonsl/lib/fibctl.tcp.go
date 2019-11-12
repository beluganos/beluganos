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

package gonslib

import (
	fibcapi "fabricflow/fibc/api"
	fibcnet "fabricflow/fibc/net"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type FIBTcpController struct {
	addr string
	dpId uint64
	*fibcnet.Client
	recvCh chan *fibcapi.DpMonitorReply

	log *log.Entry
}

func NewFIBTcpController(addr string, dpId uint64) *FIBTcpController {
	return &FIBTcpController{
		addr:   addr,
		dpId:   dpId,
		Client: fibcnet.NewClient(addr),
		recvCh: make(chan *fibcapi.DpMonitorReply),

		log: log.WithFields(log.Fields{"module": "FIBTcpCtrl", "fibcd": addr, "dpId": dpId}),
	}
}

func (c *FIBTcpController) String() string {
	return fmt.Sprintf("FIBTcpController(%s, %d)", c.addr, c.dpId)
}

func (c *FIBTcpController) monitor() {
	c.Client.Start(func(client *fibcnet.Client) {
		c.log.Debugf("Monitor: START.")

		for {
			hdr, data, err := client.Read()
			if err != nil {
				c.log.Errorf("Monitor:  EXIT. Read error. %s", err)
				return
			}

			reply, err := fibcapi.ParseDpMonitprReply(hdr, data)
			if err != nil {
				c.log.Warnf("Monitor: Parse Monitor reply error. %s", err)
			} else {
				c.recvCh <- reply
			}
		}
	})
}

func (c *FIBTcpController) Start() error {
	go c.monitor()
	return nil
}

func (c *FIBTcpController) Recv() <-chan *fibcapi.DpMonitorReply {
	return c.recvCh
}

func (c *FIBTcpController) Hello(hello *fibcapi.FFHello) error {
	return c.Client.Write(hello, 0)
}

func (c *FIBTcpController) PacketIn(pktin *fibcapi.FFPacketIn) error {
	return c.Client.Write(pktin, 0)
}

func (c *FIBTcpController) PortStatus(ps *fibcapi.FFPortStatus) error {
	return c.Client.Write(ps, 0)
}

func (c *FIBTcpController) L2AddrStatus(status *fibcapi.FFL2AddrStatus) error {
	return c.Client.Write(status, 0)
}
func (c *FIBTcpController) MultipartReply(reply *fibcapi.FFMultipart_Reply, xid uint32) error {
	return c.Client.Write(reply, xid)
}
