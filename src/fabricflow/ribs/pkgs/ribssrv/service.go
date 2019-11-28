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

package ribssrv

import (
	"fabricflow/ribs/api/ribsapi"
	gobgputil "fabricflow/util/gobgp"
	"gonla/nlalib"
	"net"

	api "github.com/osrg/gobgp/api"
	gobgpapi "github.com/osrg/gobgp/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//
// RibsService is base service.
//
type RibsService struct {
	NLAAddr  string
	BgpdAddr string
	CoreAddr string
	APIAddr  string
	Family   string
	RT       string

	nla *NLAController
	bgp *gobgputil.BgpMonitor

	bgpConnCh chan *nlalib.ConnInfo
	bgpPathCh chan *RibUpdate
	apiPathCh chan *RibUpdate

	family *gobgpapi.Family

	log *log.Entry
}

func (s *RibsService) startNLACtrl() error {
	s.nla = NewNLAController()
	if err := s.nla.Start(s.NLAAddr); err != nil {
		s.log.Errorf("Start(NLA): %s", err)
		return err
	}

	s.log.Debugf("Start(NLA): success.")
	return nil
}

//
// BgpConnected process bgp connected message.
//
func (s *RibsService) BgpConnected(conn *nlalib.ConnInfo) {
	s.bgpConnCh <- conn
}

//
// BgpPathUpdate process bgp path update message.
//
func (s *RibsService) BgpPathUpdate(path *api.Path) {
	s.bgpPathCh <- NewRibUpdate(path, s.RT)
}

func (s *RibsService) startBGPMonitor(done <-chan struct{}) error {
	family, err := gobgputil.ParseRouteFamily(s.Family)
	if err != nil {
		return err
	}
	s.family = family

	bgp, err := gobgputil.NewBgpMonitor(s.BgpdAddr, s.Family)
	if err != nil {
		s.log.Errorf("Start(BGP): %s", err)
		return err
	}
	s.bgp = bgp

	go s.bgp.Serve(done, s)

	s.log.Debugf("Start(BGP): success.")
	return nil
}

func (s *RibsService) startCoreAPIServer(apisrv ribsapi.RIBSCoreApiServer) error {
	listen, err := net.Listen("tcp", s.CoreAddr)
	if err != nil {
		s.log.Errorf("Start(CoreAPI): %s", err)
		return err
	}

	server := grpc.NewServer()
	ribsapi.RegisterRIBSCoreApiServer(server, apisrv)
	go server.Serve(listen)

	s.log.Debugf("Start(CoreAPI): success.")
	return nil
}

func (s *RibsService) startAPIServer(apisrv ribsapi.RIBSApiServer) error {
	listen, err := net.Listen("tcp", s.APIAddr)
	if err != nil {
		s.log.Errorf("Start(API): %s", err)
		return err
	}

	server := grpc.NewServer()
	ribsapi.RegisterRIBSApiServer(server, apisrv)
	go server.Serve(listen)

	s.log.Debugf("Start(API): success.")
	return nil
}

func (s *RibsService) init(module string) {
	s.log = log.WithFields(log.Fields{"module": module})
	s.bgpConnCh = make(chan *nlalib.ConnInfo)
	s.bgpPathCh = make(chan *RibUpdate)
	s.apiPathCh = make(chan *RibUpdate)
}
