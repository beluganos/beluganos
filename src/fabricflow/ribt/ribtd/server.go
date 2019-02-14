// -*- coding: utf-8 -*-

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

package main

import (
	gobgputil "fabricflow/util/gobgp"
	"fabricflow/util/gobgp/apiutil"
	"net"
	"time"

	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/packet/bgp"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	TUN_ADD_WAIT = 1 * time.Second
)

type Server struct {
	*gobgputil.BgpMonitor
	*Tables
	local4 net.IP
	local6 net.IP

	TunType6   bgp.TunnelType
	TunForce   bgp.TunnelType
	TunDefault bgp.TunnelType
}

func NewServer(addr string, prefix, family string, local4, local6 net.IP) (*Server, error) {
	monitor, err := gobgputil.NewBgpMonitor(addr, family)
	if err != nil {
		return nil, err
	}

	return &Server{
		BgpMonitor: monitor,
		Tables:     NewTables(prefix),
		local4:     local4,
		local6:     local6,
		TunType6:   bgp.TUNNEL_TYPE_IPV6,
		TunForce:   0,
		TunDefault: 0,
	}, nil
}

func (s *Server) SetTunType6(tunType uint16) {
	s.TunType6 = bgp.TunnelType(tunType)
}

func (s *Server) SetTunForce(tunType uint16) {
	s.TunForce = bgp.TunnelType(tunType)
}

func (s *Server) SetTunTypeDefault(tunType uint16) {
	s.TunDefault = bgp.TunnelType(tunType)
}

func (s *Server) Start(done <-chan struct{}) error {
	s.Tunnels().Reset()
	if err := s.Tunnels().Load(); err != nil {
		log.Errorf("Server: LinkLoad error. %s", err)
		return err
	}

	go s.Serve(done, func(path *api.Path) {
		p, err := apiutil.NewNativePath(path)
		if err != nil {
			log.Errorf("Server.Serve NewNativePath error. %s", err)
			return
		}
		s.ProcessPath(p)
	})

	log.Infof("Serve: Started")
	return nil
}

func (s *Server) addTunnel(route *TunnelRoute) (*TunnelEntry, error) {

	var (
		localIP net.IP
		link    netlink.Link
		err     error
	)

	tunType := func() bgp.TunnelType {
		if s.TunForce != 0 {
			return s.TunForce
		}

		if route.TunnelType == 0 {
			return s.TunDefault
		}

		return route.TunnelType
	}()

	switch tunType {
	case bgp.TUNNEL_TYPE_IP_IN_IP:
		localIP = s.local4

	case s.TunType6:
		localIP = s.local6

	default:
		return nil, nil
	}

	ifname, tunId := s.Tunnels().NewIfName()
	link = NewIptun(ifname, route.Nexthop, localIP)
	if link, err = AddLink(link); err != nil {
		log.Errorf("AddLink error. %s", err)
		return nil, err
	}

	tun := NewTunnelEntry(link, tunId)
	s.Tunnels().Put(tun)

	return tun, nil
}

func (s *Server) delTunnel(tun *TunnelEntry) {
	s.Tunnels().Pop(tun.Remote())
	DelLinkByName(tun.Ifname())
}

func (s *Server) ProcessPath(path *apiutil.Path) {

	tunRoute, err := NewTunnelRouteFromPath(path)
	if err != nil {
		log.Errorf("NewBgpRouteFromPath error. %s", err)
		return
	}

	// expRoute := NewExportRouteFromPath(path)

	log.Debugf("Route(tunnel) %s", tunRoute)
	// log.Debugf("Route(export) %s", expRoute)

	if path.IsWithdraw {
		log.Debugf("route del %s.", tunRoute.Prefix)

		tun, ok := s.Tunnels().FindByPrefix(tunRoute.Prefix.String())
		log.Debugf("tunnel %s is_tunnel=%t", tun, ok)
		if ok {
			//if err := DelBgpPath(s.Client(), expRoute); err != nil {
			//	log.Errorf("DeletePath error. %s", err)
			//}

			if err := DelRoute(tunRoute.Prefix, tun.Ifindex()); err != nil {
				log.Errorf("DelRoute error. %s %s", tunRoute, err)
			}

			if n := s.Tunnels().DelRoute(tunRoute, tun); n == 0 {
				s.delTunnel(tun)
				log.Debugf("tunnel del %s", tun.Ifname())
			}
		}

	} else {
		tun, ok := s.Tunnels().FindByRemote(tunRoute.Nexthop.String())
		if !ok {
			if tun, err = s.addTunnel(tunRoute); err != nil {
				log.Errorf("addTunnel error. %s", err)
				return
			}
			if tun == nil {
				log.Debugf("%s is not tunnel remote.", tunRoute.Nexthop)
				return
			}

			time.Sleep(TUN_ADD_WAIT)

			log.Debugf("tunnel add %s", tun.Ifname())
		}

		log.Debugf("tunnel %s", tun)

		if err := AddRoute(tunRoute.Prefix, tun.Ifindex()); err != nil {
			log.Errorf("AddRoute error. %s %s", tunRoute, err)
		}

		s.Tunnels().AddRoute(tunRoute, tun)
		log.Debugf("route add %s dev %s", tunRoute.Prefix, tun.Ifname())

		//if err := AddBgpPath(s.Client(), expRoute); err != nil {
		//	log.Errorf("AddPath error. %s", err)
		//}
	}
}
