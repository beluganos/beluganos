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

package ffgrpc

import (
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type ClientConnInfo struct {
	LocalAddr  net.Addr
	RemoteAddr net.Addr
}

func NewClientConnInfo(local, remote net.Addr) *ClientConnInfo {
	return &ClientConnInfo{
		LocalAddr:  local,
		RemoteAddr: remote,
	}
}

func (c *ClientConnInfo) String() string {
	return fmt.Sprintf("local:%s remote:%s", c.LocalAddr, c.RemoteAddr)
}

func NewClientConnChan() chan *ClientConnInfo {
	return make(chan *ClientConnInfo)
}

type ClientConnHandler struct {
	ch  chan *ClientConnInfo
	log *log.Entry
}

func NewClientConnHandler(ch chan *ClientConnInfo) *ClientConnHandler {
	return &ClientConnHandler{
		ch:  ch,
		log: log.WithFields(log.Fields{"module": "grpc.conn"}),
	}
}

func (h *ClientConnHandler) TagRPC(ctxt context.Context, info *stats.RPCTagInfo) context.Context {
	// h.log.Debugf("TagRPC %v", info)
	return ctxt
}

func (h *ClientConnHandler) HandleRPC(ctxt context.Context, st stats.RPCStats) {
	// h.log.Debugf("HandleRPC %v", st)
}

func (h *ClientConnHandler) TagConn(ctxt context.Context, info *stats.ConnTagInfo) context.Context {
	// h.log.Debugf("TagConn %v", info)
	h.ch <- NewClientConnInfo(info.LocalAddr, info.RemoteAddr)
	return ctxt
}

func (h *ClientConnHandler) HandleConn(ctxt context.Context, st stats.ConnStats) {
	// h.log.Debugf("HandleConn %v", st)
}

func NewClientConn(addr string) (*grpc.ClientConn, chan *ClientConnInfo, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	ch := NewClientConnChan()
	opts = append(opts, grpc.WithStatsHandler(NewClientConnHandler(ch)))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		close(ch)
		return nil, nil, err
	}

	return conn, ch, nil
}
