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

package nlalib

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"net"
)

type ConnInfo struct {
	LocalAddr  net.IP
	RemoteAddr net.IP
}

func NewConnInfo(local net.IP, remote net.IP) *ConnInfo {
	return &ConnInfo{
		LocalAddr:  local,
		RemoteAddr: remote,
	}
}

type ConnStatsHandler struct {
	ch chan<- *ConnInfo
}

func NewConnStatsHandler(ch chan<- *ConnInfo) stats.Handler {
	return &ConnStatsHandler{
		ch: ch,
	}
}

func (h *ConnStatsHandler) TagRPC(ctxt context.Context, info *stats.RPCTagInfo) context.Context {
	return ctxt
}

func (h *ConnStatsHandler) HandleRPC(ctxt context.Context, st stats.RPCStats) {
	// nothing to do.
}

func (h *ConnStatsHandler) TagConn(ctxt context.Context, info *stats.ConnTagInfo) context.Context {
	local, _, _ := net.SplitHostPort(info.LocalAddr.String())
	remote, _, _ := net.SplitHostPort(info.RemoteAddr.String())
	h.ch <- NewConnInfo(net.ParseIP(local), net.ParseIP(remote))

	return ctxt
}

func (h *ConnStatsHandler) HandleConn(ctxt context.Context, st stats.ConnStats) {
	// nothing to do.
}

func NewClientConn(addr string, ch chan<- *ConnInfo) (*grpc.ClientConn, error) {

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if ch != nil {
		opts = append(opts, grpc.WithStatsHandler(NewConnStatsHandler(ch)))
	}
	return grpc.Dial(addr, opts...)
}
