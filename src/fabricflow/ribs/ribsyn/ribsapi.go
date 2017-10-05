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
	"google.golang.org/grpc"
	"net"
)

type ApiServer struct {
	addr string
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{
		addr: addr,
	}
}

func (a *ApiServer) Start() error {
	if len(a.addr) == 0 {
		log.Infof("APIS Start: exit.")
		return nil
	}

	listen, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	ribsapi.RegisterRIBSApiServer(s, a)
	go s.Serve(listen)

	log.Infof("APIS Start: %s", a.addr)
	return nil
}

func (a *ApiServer) GetRics(req *ribsapi.GetRicsRequest, stream ribsapi.RIBSApi_GetRicsServer) error {
	return Tables.Rics.Walk(func(key string, entry *RicEntry) error {
		return stream.Send(ribsapi.NewRicEntryFromNative(key, &entry.RicEntry))
	})
}

func (a *ApiServer) GetNexthops(req *ribsapi.GetNexthopsRequest, stream ribsapi.RIBSApi_GetNexthopsServer) error {
	return Tables.Nexthops.Walk(func(nh *ribsmsg.Nexthop) error {
		return stream.Send(ribsapi.NewNexthopFromNative(nh))
	})
}

func (a *ApiServer) GetNexthopMap(req *ribsapi.GetNexthopMapRequest, stream ribsapi.RIBSApi_GetNexthopMapServer) error {
	Tables.NexthopMap.Walk(func(key string, val net.IP) bool {
		if err := stream.Send(ribsapi.NewNexthopMapFromNative(key, val.String())); err != nil {
			return false
		}
		return true
	})
	return nil
}
