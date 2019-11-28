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

package ribssrv

import (
	"fabricflow/ribs/api/ribsapi"
	"gonla/nlalib"

	api "github.com/osrg/gobgp/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//
// CoreAPIClient is client for core api.
//
type CoreAPIClient struct {
	client ribsapi.RIBSCoreApiClient
	connCh chan *nlalib.ConnInfo
	pathCh chan *RibUpdate

	log *log.Entry
}

//
// NewCoreAPIClient returns new CoreAPIClient.
//
func NewCoreAPIClient(pathCh chan *RibUpdate) *CoreAPIClient {
	return &CoreAPIClient{
		client: nil,
		connCh: make(chan *nlalib.ConnInfo),
		pathCh: pathCh,

		log: log.WithFields(log.Fields{"module": "apictl"}),
	}
}

//
// Conn returns connection chan.
//
func (c *CoreAPIClient) Conn() <-chan *nlalib.ConnInfo {
	return c.connCh
}

//
// Recv returns receiver chan.
//
func (c *CoreAPIClient) Recv() <-chan *RibUpdate {
	return c.pathCh
}

//
// Start starts main thread.
//
func (c *CoreAPIClient) Start(addr string) error {
	con, err := nlalib.NewClientConn(addr, c.connCh)
	if err != nil {
		return err
	}

	c.client = ribsapi.NewRIBSCoreApiClient(con)
	return nil
}

//
// ModRib send mod rib request.
//
func (c *CoreAPIClient) ModRib(path *api.Path, rt string) error {
	rib := RibUpdate{
		Path: path,
		Rt:   rt,
	}

	req, err := rib.ToAPI()
	if err != nil {
		return err
	}

	if _, err := c.client.ModRib(context.Background(), req); err != nil {
		return err
	}

	return nil
}

//
// Monitor monitor messages.
//
func (c *CoreAPIClient) Monitor(nid uint8, rt string) error {
	req := &ribsapi.MonitorRibRequest{
		NId: uint32(nid),
		Rt:  rt,
	}

	stream, err := c.client.MonitorRib(context.Background(), req)
	if err != nil {
		c.log.Errorf("Monitor: %s", err)
		return err
	}

	c.log.Infof("Monitor: START")

	go func() {
	FOR_LOOP:
		for {
			msg, err := stream.Recv()
			if err != nil {
				c.log.Infof("Monitor: EXIT. %s", err)
				break FOR_LOOP
			}

			path, err := NewRibUpdateFromAPI(msg)
			if err != nil {
				c.log.Warnf("Monitor: convert error. %s", err)
				continue FOR_LOOP
			}

			c.pathCh <- path
		}
	}()

	return nil
}

//
// SyncRib send sync message.
//
func (c *CoreAPIClient) SyncRib(rt string) error {
	req := &ribsapi.SyncRibRequest{
		Rt: rt,
	}

	_, err := c.client.SyncRib(context.Background(), req)
	return err
}
