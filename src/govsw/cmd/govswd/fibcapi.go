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
	fibcapi "fabricflow/fibc/api"
	ffgrpc "fabricflow/util/grpc"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type FIBCVsApiHandler interface {
	FIBCClient(fibcapi.FIBCVsApiClient)
	FIBCConnect()
}

type FIBCVsApiClient struct {
	Addr     string
	Listener FIBCVsApiHandler

	conn   *grpc.ClientConn
	connCh chan *ffgrpc.ClientConnInfo
	client fibcapi.FIBCVsApiClient

	log *log.Entry
}

func (c *FIBCVsApiClient) init() error {
	conn, connCh, err := ffgrpc.NewClientConn(c.Addr)
	if err != nil {
		return err
	}

	c.conn = conn
	c.connCh = connCh
	c.client = fibcapi.NewFIBCVsApiClient(conn)
	c.log = log.WithFields(log.Fields{"module": "fibccli"})

	return nil
}

func (c *FIBCVsApiClient) Serve(done chan struct{}) {
	c.log.Infof("Serve: started.")

FOR_LOOP:
	for {
		select {
		case conn := <-c.connCh:
			c.log.Infof("Serve: conn %s", conn.RemoteAddr)
			c.Listener.FIBCConnect()

		case <-done:
			c.log.Infof("Serve: exit.")
			break FOR_LOOP
		}
	}
}

func (c *FIBCVsApiClient) Start(done chan struct{}) error {
	if err := c.init(); err != nil {
		c.log.Errorf("Start: init error. %s", err)
		return err
	}

	c.Listener.FIBCClient(c.client)

	go c.Serve(done)

	c.log.Infof("Start: success. %s", c.Addr)

	return nil
}
