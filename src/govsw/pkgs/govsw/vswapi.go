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

package govsw

import (
	"context"
	"fmt"
	"govsw/api/vswapi"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const VSWAPI_PORT = 50099

type VswApiListener interface {
	VswAPISync(string)
	VswAPISaveConfig() error
}

type VswApiServer struct {
	DB       *DB
	Listener VswApiListener
	Addr     string

	log *log.Entry
}

func (s *VswApiServer) Serve(lis net.Listener) {
	s.log.Infof("Serve: started. %s", lis.Addr())

	server := grpc.NewServer()
	vswapi.RegisterVswApiServer(server, s)
	server.Serve(lis)

	s.log.Infof("Serve: exit.")
}

func (s *VswApiServer) Start() error {
	s.log = log.WithFields(log.Fields{"module": "vswapi"})

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		s.log.Errorf("Listen error. %s %s", s.Addr, err)
		return err
	}

	go s.Serve(lis)

	s.log.Infof("Start: success. %s", s.Addr)
	return nil
}

func (s *VswApiServer) ModIfname(ctxt context.Context, req *vswapi.ModIfnameRequest) (*vswapi.ModIfnameReply, error) {
	switch req.Cmd {
	case vswapi.ModIfnameRequest_ADD:
		s.DB.Ifname().RegisterIfname(req.Ifname)

	case vswapi.ModIfnameRequest_REG:
		s.DB.Ifname().RegisterPattern(req.Ifname)

	case vswapi.ModIfnameRequest_DELETE:
		s.DB.Ifname().Delete(req.Ifname)

	case vswapi.ModIfnameRequest_SYNC:
		s.Listener.VswAPISync(req.Ifname)

	default:
		return nil, fmt.Errorf("Invalid command. %s", req.Cmd)
	}

	return &vswapi.ModIfnameReply{}, nil
}

func (s *VswApiServer) GetIfnames(req *vswapi.GetIfnamesRequest, stream vswapi.VswApi_GetIfnamesServer) error {
	s.DB.Ifname().Range(func(kind, ifname string) {
		reply := vswapi.GetIfnamesReply{
			Ifname: ifname,
			Kind:   kind,
		}

		if err := stream.Send(&reply); err != nil {
			s.log.Errorf("send ifname error. %s", err)
		}
	})

	return nil
}

func (s *VswApiServer) ModLink(ctxt context.Context, req *vswapi.ModLinkRequest) (*vswapi.ModLinkReply, error) {
	operState, err := func() (bool, error) {
		switch req.Cmd {
		case vswapi.ModLinkRequest_UP:
			return true, nil

		case vswapi.ModLinkRequest_DOWN:
			return false, nil

		default:
			return false, fmt.Errorf("Invalid command. %s", req.Cmd)
		}
	}()
	if err != nil {
		s.log.Errorf("ModLink: %s", err)
		return nil, err
	}

	if err := s.DB.Link().GetByName(req.Ifname, func(link *Link) error {
		return link.SetOperStatus(operState)
	}); err != nil {
		s.log.Errorf("ModLink: %s", err)
		return nil, err
	}

	return &vswapi.ModLinkReply{}, nil
}

func (s *VswApiServer) GetLinks(req *vswapi.GetLinksRequest, stream vswapi.VswApi_GetLinksServer) error {
	s.DB.Link().Range(func(ifindex int, link *Link) {
		reply := vswapi.GetLinksReply{
			Index: int32(ifindex),
			Name:  link.Name(),
		}
		if err := stream.Send(&reply); err != nil {
			s.log.Errorf("send link error. %s", err)
		}
	})

	return nil
}

func (s *VswApiServer) GetStats(req *vswapi.GetStatsRequest, stream vswapi.VswApi_GetStatsServer) error {
	replies := []*vswapi.GetStatsReply{}
	s.DB.Link().Range(func(ifindex int, link *Link) {
		reply := &vswapi.GetStatsReply{
			Group:  fmt.Sprintf("link/%s", link.Name()),
			Values: map[string]uint64{},
		}

		link.Stats().Range(func(name string, value uint64) {
			reply.Values[name] = value
		})

		replies = append(replies, reply)
	})

	for _, reply := range replies {
		if err := stream.Send(reply); err != nil {
			return err
		}
	}

	return nil
}

func (s *VswApiServer) SaveConfig(ctxt context.Context, req *vswapi.SaveConfigRequest) (*vswapi.SaveConfigReply, error) {
	if err := s.Listener.VswAPISaveConfig(); err != nil {
		s.log.Errorf("SaveConfig: %s", err)
		return nil, err
	}

	return &vswapi.SaveConfigReply{}, nil
}
