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
	ffgrpc "fabricflow/util/grpc"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type FIBGrpcController struct {
	addr   string
	reId   string
	conn   *grpc.ClientConn
	connCh chan bool
	recvCh chan *fibcapi.VmMonitorReply
	client fibcapi.FIBCVmApiClient

	done chan struct{}
	log  *log.Entry
}

func NewFIBGrpcController(addr string, reId string) *FIBGrpcController {
	return &FIBGrpcController{
		addr:   addr,
		reId:   reId,
		connCh: make(chan bool),
		recvCh: make(chan *fibcapi.VmMonitorReply),

		log: log.WithFields(log.Fields{"module": "FIBGrpcController"}),
	}
}

func (c *FIBGrpcController) FIBCType() string {
	return FIBCTypeGrpc
}

func (c *FIBGrpcController) monitor() {
	req := fibcapi.VmMonitorRequest{
		ReId: c.reId,
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
	c.client = fibcapi.NewFIBCVmApiClient(conn)
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

func (c *FIBGrpcController) Recv() <-chan *fibcapi.VmMonitorReply {
	return c.recvCh
}

func (c *FIBGrpcController) Hello(hello *fibcapi.Hello) error {
	if c.client == nil {
		return fmt.Errorf("Hello: bad client status.")
	}

	_, err := c.client.SendHello(context.Background(), hello)
	return err
}

func (c *FIBGrpcController) PortConfig(pc *fibcapi.PortConfig) error {
	if c.client == nil {
		return fmt.Errorf("Hello: bad client status.")
	}

	_, err := c.client.SendPortConfig(context.Background(), pc)
	return err
}

func (c *FIBGrpcController) FlowMod(mod *fibcapi.FlowMod) error {
	if c.client == nil {
		return fmt.Errorf("Hello: bad client status.")
	}

	_, err := c.client.SendFlowMod(context.Background(), mod)
	return err
}

func (c *FIBGrpcController) GroupMod(mod *fibcapi.GroupMod) error {
	if c.client == nil {
		return fmt.Errorf("Hello: bad client status.")
	}

	_, err := c.client.SendGroupMod(context.Background(), mod)
	return err
}
