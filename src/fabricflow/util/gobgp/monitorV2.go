// -*- coding: utf-8 -*-
// +build !gobgpv1

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

package gobgputil

import (
	"context"
	"fmt"
	"gonla/nlalib"
	"io"

	api "github.com/osrg/gobgp/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	FamilyIPv4UC = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_UNICAST,
	}
	FamilyIPv6UC = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_UNICAST,
	}
	FamilyIPv4VPN = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_MPLS_VPN,
	}
	FamilyIPv6VPN = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_MPLS_VPN,
	}
	FamilyIPv4Encap = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_ENCAPSULATION,
	}
	FamilyIPv6Encap = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_ENCAPSULATION,
	}
)

var family_values = map[string]*api.Family{
	"ipv4-unicast":       FamilyIPv4UC,
	"ipv6-unicast":       FamilyIPv6UC,
	"l3vpn-ipv4-unicast": FamilyIPv4VPN,
	"l3vpn-ipv6-unicast": FamilyIPv6VPN,
	"ipv4-encap":         FamilyIPv4Encap,
	"ipv6-encap":         FamilyIPv6Encap,
}

func ParseRouteFamily(s string) (*api.Family, error) {
	if v, ok := family_values[s]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("invalid route family. %s", s)
}

type BgpMonitorListener interface {
	BgpConnected(*nlalib.ConnInfo)
	BgpPathUpdate(*api.Path)
}

type BgpMonitor struct {
	client api.GobgpApiClient
	conn   *grpc.ClientConn
	connCh chan *nlalib.ConnInfo
	family *api.Family

	log *log.Entry
}

func NewBgpMonitor(addr string, routeFamily string) (*BgpMonitor, error) {
	family, err := ParseRouteFamily(routeFamily)
	if err != nil {
		return nil, err
	}

	connCh := make(chan *nlalib.ConnInfo)
	conn, err := nlalib.NewClientConn(addr, connCh)
	if err != nil {
		close(connCh)
		return nil, err
	}

	return &BgpMonitor{
		client: api.NewGobgpApiClient(conn),
		conn:   conn,
		connCh: connCh,
		family: family,
		log:    log.WithFields(log.Fields{"module": fmt.Sprintf("bgpmon(%s)", routeFamily)}),
	}, nil
}

func (s *BgpMonitor) Client() api.GobgpApiClient {
	return s.client
}

func (s *BgpMonitor) Serve(done <-chan struct{}, listener BgpMonitorListener) {

	s.log.Infof("Serve: Started")

	pathCh := make(chan *api.Path)
	defer close(pathCh)

FOR_LOOP:
	for {
		select {
		case conn := <-s.connCh:
			s.log.Infof("Serve: connected.")
			listener.BgpConnected(conn)
			go s.MonitorTable(pathCh)

		case path := <-pathCh:
			// s.log.Tracef("Serve: %v", path)
			listener.BgpPathUpdate(path)

		case <-done:
			s.log.Infof("Serve: Exit")
			break FOR_LOOP
		}
	}
}

func (s *BgpMonitor) MonitorTable(ch chan<- *api.Path) {
	s.log.Infof("MonitorTable Start.")

	req := &api.MonitorTableRequest{
		TableType: api.TableType_GLOBAL,
		Family:    s.family,
		Current:   true,
	}

	stream, err := s.client.MonitorTable(context.Background(), req)
	if err != nil {
		s.log.Errorf("MonitorTable error. %s", err)
		return
	}

FOR_LOOP:
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			s.log.Infof("MonitorRib Recv end. %s", err)
			break FOR_LOOP
		}
		if err != nil {
			s.log.Errorf("MonitorRib error. %s", err)
			break FOR_LOOP
		}

		ch <- res.Path
	}

	s.log.Infof("MonitorRib Exit.")
}
