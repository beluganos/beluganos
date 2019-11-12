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
	"context"
	fibcapi "fabricflow/fibc/api"
	ffgrpc "fabricflow/util/grpc"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type FIBGrpcController struct {
	addr   string
	dpId   uint64
	conn   *grpc.ClientConn
	connCh chan bool
	recvCh chan *fibcapi.DpMonitorReply
	client fibcapi.FIBCDpApiClient

	done chan struct{}
	log  *log.Entry
}

func NewFIBGrpcController(addr string, dpId uint64) *FIBGrpcController {
	return &FIBGrpcController{
		addr:   addr,
		dpId:   dpId,
		connCh: make(chan bool),
		recvCh: make(chan *fibcapi.DpMonitorReply),

		log: log.WithFields(log.Fields{"module": "FIBGrpcCtrl", "fibcd": addr, "dpId": dpId}),
	}
}

func (c *FIBGrpcController) String() string {
	return fmt.Sprintf("FIBGrpcController(%s, %d)", c.addr, c.dpId)
}

func (c *FIBGrpcController) monitor() {
	req := fibcapi.DpMonitorRequest{
		DpId: c.dpId,
	}

	stream, err := c.client.Monitor(context.Background(), &req)
	if err != nil {
		c.log.Errorf("monitor: Monitor error. %s", err)
		return
	}

	c.log.Debugf("Monitor: started.")

FOR_LOOP:
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			c.log.Info("monitor: exit.")
			break FOR_LOOP
		}
		if err != nil {
			c.log.Errorf("monitor: exit. recv errr. %s", err)
			break FOR_LOOP
		}
		if m == nil {
			c.log.Warnf("monitor: invalid message. %v", m)
			continue
		}

		c.recvCh <- m
	}
}

func (c *FIBGrpcController) serve(connCh chan *ffgrpc.ClientConnInfo, done chan struct{}) {
	c.log.Debugf("Serve: START.")

	defer close(connCh)

FOR_LOOP:
	for {
		select {
		case conn := <-connCh:
			if conn != nil {
				c.log.Debugf("Serve: connected.")

				c.connCh <- true

				c.monitor()

			} else {
				c.log.Debugf("Serve: disconected.")
				c.connCh <- false
			}

		case <-done:
			c.log.Debugf("Serve: EXIT.")
			break FOR_LOOP
		}
	}
}

func (c *FIBGrpcController) Start() error {
	conn, connCh, err := ffgrpc.NewClientConn(c.addr)
	if err != nil {
		c.log.Errorf("Start: client create error. %s", err)
		return err
	}

	c.conn = conn
	c.client = fibcapi.NewFIBCDpApiClient(conn)
	c.done = make(chan struct{})

	go c.serve(connCh, c.done)

	c.log.Infof("Start: success.")

	return nil
}

func (c *FIBGrpcController) Stop() {
	c.log.Debugf("Stop:")
	if c.done != nil {
		close(c.done)
		c.done = nil

		c.conn.Close()

		c.log.Debugf("Stop: channels closed.")
	}

}

func (c *FIBGrpcController) Conn() <-chan bool {
	return c.connCh
}

func (c *FIBGrpcController) Recv() <-chan *fibcapi.DpMonitorReply {
	return c.recvCh
}

func (c *FIBGrpcController) Hello(hello *fibcapi.FFHello) error {
	if c.client == nil {
		return fmt.Errorf("Hello: bad client status.")
	}

	_, err := c.client.SendHello(context.Background(), hello)
	return err
}

func (c *FIBGrpcController) PacketIn(pktin *fibcapi.FFPacketIn) error {
	if c.client == nil {
		return fmt.Errorf("PacketIn: bad client status.")
	}

	_, err := c.client.SendPacketIn(context.Background(), pktin)
	return err
}

func (c *FIBGrpcController) PortStatus(ps *fibcapi.FFPortStatus) error {
	if c.client == nil {
		return fmt.Errorf("PortStatus: bad client status.")
	}

	_, err := c.client.SendPortStatus(context.Background(), ps)
	return err
}

func (c *FIBGrpcController) L2AddrStatus(l2addr *fibcapi.FFL2AddrStatus) error {
	if c.client == nil {
		return fmt.Errorf("L2AddrStatus: bad client status.")
	}

	_, err := c.client.SendL2AddrStatus(context.Background(), l2addr)
	return err
}

func (c *FIBGrpcController) MultipartReply(reply *fibcapi.FFMultipart_Reply, xid uint32) error {
	mp := fibcapi.DpMultipartReply{
		Xid:   xid,
		Reply: reply,
	}

	_, err := c.client.SendMultipartReply(context.Background(), &mp)
	return err
}
