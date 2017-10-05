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

package ribsyn

import (
	"fabricflow/ribs/ribsapi"
	"fabricflow/ribs/ribsmsg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gonla/nlalib"
	"google.golang.org/grpc"
	"net"
)

type CoreApiServer struct {
	recvCh chan *ribsmsg.RibUpdate
	connCh chan *RicEntry
	syncCh chan string
}

func NewCoreApiServer() *CoreApiServer {
	return &CoreApiServer{
		recvCh: make(chan *ribsmsg.RibUpdate),
		connCh: make(chan *RicEntry),
		syncCh: make(chan string),
	}
}

func (c *CoreApiServer) Start(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	ribsapi.RegisterRIBSCoreApiServer(s, c)
	go s.Serve(listen)

	log.Infof("CORE: Start: %s", addr)
	return nil
}

func (c *CoreApiServer) MonitorRib(req *ribsapi.MonitorRibRequest, stream ribsapi.RIBSCoreApi_MonitorRibServer) error {
	entry := NewRicEntryOnMic(uint8(req.NId), req.Rt)

	defer func() {
		entry.RicEntry.Leave = true
		c.connCh <- entry
	}()

	c.connCh <- entry

	done := stream.Context().Done()
	for {
		select {
		case <-done:
			log.Infof("CORE: Stream closed. %s", req.Rt)
			return nil
		case rib := <-entry.Ch:
			if err := stream.Send(rib); err != nil {
				log.Errorf("CORE: Stream.Send error. %s %s", req.Rt, err)
				return err
			}
		}
	}
}

func (c *CoreApiServer) ModRib(ctxt context.Context, rib *ribsapi.RibUpdate) (*ribsapi.ModRibReply, error) {
	r, err := rib.ToNative()
	if err != nil {
		return &ribsapi.ModRibReply{}, err
	}

	log.Debugf("CORE: recv %v", r)

	c.recvCh <- r
	return &ribsapi.ModRibReply{}, nil
}

func (c *CoreApiServer) SyncRib(ctxt context.Context, req *ribsapi.SyncRibRequest) (*ribsapi.SyncRibReply, error) {
	c.syncCh <- req.Rt
	return &ribsapi.SyncRibReply{}, nil
}

func (c *CoreApiServer) Recv() <-chan *ribsmsg.RibUpdate {
	return c.recvCh
}

func (c *CoreApiServer) Conn() <-chan *RicEntry {
	return c.connCh
}

func (c *CoreApiServer) Sync() <-chan string {
	return c.syncCh
}

type CoreApiClient struct {
	Api   ribsapi.RIBSCoreApiClient
	con   *grpc.ClientConn
	conCh chan *nlalib.ConnInfo
	ribCh chan *ribsmsg.RibUpdate
}

func NewCoreApiClient(addr string) (*CoreApiClient, error) {
	conCh := make(chan *nlalib.ConnInfo)
	con, err := nlalib.NewClientConn(addr, conCh)
	if err != nil {
		return nil, err
	}

	return &CoreApiClient{
		Api:   ribsapi.NewRIBSCoreApiClient(con),
		con:   con,
		conCh: conCh,
		ribCh: make(chan *ribsmsg.RibUpdate),
	}, nil
}

func (c *CoreApiClient) Conn() <-chan *nlalib.ConnInfo {
	return c.conCh
}

func (c *CoreApiClient) Recv() <-chan *ribsmsg.RibUpdate {
	return c.ribCh
}

func (c *CoreApiClient) ModRib(rib *ribsmsg.RibUpdate) error {
	r, err := ribsapi.NewRibUpdateFromNative(rib)
	if err != nil {
		return err
	}

	_, err = c.Api.ModRib(context.Background(), r)
	return err
}

func (c *CoreApiClient) Monitor(nid uint8, rt string) error {
	req := &ribsapi.MonitorRibRequest{
		NId: uint32(nid),
		Rt:  rt,
	}
	stream, err := c.Api.MonitorRib(context.Background(), req)
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}

			rib, err := msg.ToNative()
			if err != nil {
				log.Errorf("CORE: ToNative error. %s", err)
				continue
			}

			c.ribCh <- rib
		}
		log.Errorf("CORE: Monitor EXIT. %s", rt)
	}()

	return nil
}

func (a *CoreApiClient) SyncRib(rt string) error {
	req := &ribsapi.SyncRibRequest{
		Rt: rt,
	}

	_, err := a.Api.SyncRib(context.Background(), req)
	return err
}
