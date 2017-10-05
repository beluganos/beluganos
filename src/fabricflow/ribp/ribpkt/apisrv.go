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

package ribpkt

import (
	"fabricflow/ribp/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type ApiServer struct {
	addr string
	ch   chan<- string
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{
		addr: addr,
		ch:   nil,
	}
}

func (a *ApiServer) Start(ch chan<- string) error {
	listen, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}

	a.ch = ch

	s := grpc.NewServer()
	ribpapi.RegisterRIBPApiServer(s, a)
	go s.Serve(listen)

	log.Infof("RIBPApiServer: START")
	return nil
}

func (a *ApiServer) SendFFPacket(ctxt context.Context, req *ribpapi.FFPacketRequest) (*ribpapi.SendFFPacketReply, error) {
	log.Infof("API  %s", req.Ifname)
	a.ch <- req.Ifname
	return &ribpapi.SendFFPacketReply{}, nil
}
