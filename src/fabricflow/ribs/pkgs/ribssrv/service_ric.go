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
	"fabricflow/util/gobgp/apiutil"
	"gonla/nlalib"

	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/pkg/packet/bgp"
	log "github.com/sirupsen/logrus"
)

//
// RicService is ric service.
//
type RicService struct {
	RibsService

	NId     uint8
	RD      string
	Labels  []uint32
	DummyIF string

	rd    bgp.RouteDistinguisherInterface
	ecRT  *bgp.PathAttributeExtendedCommunities
	label *bgp.MPLSLabelStack

	client *CoreAPIClient
}

func (s *RicService) initRic() error {
	rd, err := bgp.ParseRouteDistinguisher(s.RD)
	if err != nil {
		s.log.Errorf("initBGPInfo: bad RD. '%s'", s.RD)
		return err
	}
	s.rd = rd

	rt, err := bgp.ParseRouteTarget(s.RT)
	if err != nil {
		s.log.Errorf("initBGPInfo: bad RT. '%s'", s.RT)
		return err
	}
	s.ecRT = bgp.NewPathAttributeExtendedCommunities(
		[]bgp.ExtendedCommunityInterface{rt},
	)

	if s.Labels == nil || len(s.Labels) == 0 {
		s.log.Errorf("initBGPInfo: bad labels. '%v'", s.Labels)
		return err
	}
	s.label = bgp.NewMPLSLabelStack(s.Labels...)

	return nil
}

//
// Start starts main thread,
//
func (s *RicService) Start(done <-chan struct{}) error {
	s.RibsService.init("svc.ric")

	s.log.Debugf("Start:")

	if err := s.initRic(); err != nil {
		return err
	}

	if err := s.startBGPMonitor(done); err != nil {
		return err
	}

	client := NewCoreAPIClient(s.apiPathCh)
	if err := client.Start(s.CoreAddr); err != nil {
		return err
	}
	s.client = client

	go s.serve(done)

	s.log.Debugf("Start: success.")
	return nil
}

func (s *RicService) serve(done <-chan struct{}) {

	s.log.Infof("serve: START")

FOR_LOOP:
	for {
		select {
		case conn, ok := <-s.bgpConnCh:
			if ok {
				s.log.Debugf("serve: bgp connected %s", conn.RemoteAddr)

				s.log.Debugf("Serve: sync rib %s", s.RT)
				if err := s.client.SyncRib(s.RT); err != nil {
					s.log.Errorf("serve: SyncRib error. %s", err)
				}
			}

		case conn, ok := <-s.client.Conn():
			if ok {
				s.log.Debugf("serve: api connected %s", conn.RemoteAddr)

				if err := s.client.Monitor(s.NId, s.RT); err != nil {
					s.log.Errorf("serve: monitor error. %s", err)
				}
			}

		case upd, ok := <-s.bgpPathCh:
			if ok {
				s.log.Tracef("serve: from bgp. RT:%s", upd.Rt)
				LogBgpPath(s.log, log.TraceLevel, upd.Path)
				s.sendToMic(upd)
			}

		case upd, ok := <-s.apiPathCh:
			if ok {
				s.log.Tracef("serve: from api RT:%s", upd.Rt)
				LogBgpPath(s.log, log.TraceLevel, upd.Path)
				s.sendToGoBGP(upd)
			}

		case <-done:
			s.log.Infof("serve: EXIT")
			break FOR_LOOP
		}
	}
}

func (s *RicService) sendToMic(upd *RibUpdate) {
	s.log.Debugf("sendToMic: ")

	path := apiutil.NewNativePath(upd.Path)

	nlriIP, ok := path.GetNlri().(*bgp.IPAddrPrefix)
	if !ok {
		s.log.Errorf("sendToMic: bad nlri type. %s", path.GetNlri())
		return
	}

	nlriVPN := bgp.NewLabeledVPNIPAddrPrefix(
		nlriIP.Length,
		nlriIP.Prefix.String(),
		*s.label,
		s.rd,
	)

	nh := path.GetNexthop()
	s.log.Debugf("sendToMic: Path(IPv4) : %s via %s", nlriIP, nh)

	pattrsVPN := []bgp.PathAttributeInterface{
		bgp.NewPathAttributeNextHop(nh.String()),
		s.ecRT,
	}

	for _, pattr := range path.GetPathAttrs() {
		switch pattr.GetType() {
		case bgp.BGP_ATTR_TYPE_NEXT_HOP:
			// pass
		default:
			pattrsVPN = append(pattrsVPN, pattr)
		}
	}

	path.Nlri = nlriVPN
	path.Attrs = pattrsVPN
	path.Family = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_MPLS_VPN,
	}

	pathVPN := path.NewAPIPath()

	s.log.Debugf("sendToMic: Path(VPNv4): %s via %s rt %s", nlriVPN, nh, s.RT)
	LogBgpPath(s.log, log.TraceLevel, pathVPN)

	if err := s.client.ModRib(pathVPN, s.RT); err != nil {
		s.log.Errorf("sendToMic: send error. %s", err)
		return
	}
}

func (s *RicService) sendToGoBGP(rib *RibUpdate) {
	s.log.Debugf("sendToGoBGP: ")

	path := rib.Path

	s.setDummyRoute(path, s.DummyIF)

	if err := modGoBGPPath(s.bgp.Client(), path); err != nil {
		s.log.Errorf("sendToGoBGP: modGoBGPPath error. %s", err)
	}
}

func (s *RicService) setDummyRoute(p *api.Path, ifname string) {
	path := apiutil.NewNativePath(p)
	nhIP := path.GetNexthop()
	ifaddr := nlalib.NewIPNetFromIP(nhIP)

	if path.IsWithdraw {
		s.log.Debugf("setDummyRoute: del ifaddr. %s %s", ifaddr, ifname)

		if err := nlalib.DelIFAddr(ifaddr, ifname); err != nil {
			s.log.Warnf("setDummyRoute: del ifaddr error. %s", err)
		}

	} else {
		s.log.Debugf("setDummyRoute: add ifaddr. %s %s", ifaddr, ifname)

		if err := nlalib.AddIFAddr(ifaddr, ifname); err != nil {
			s.log.Errorf("setDummyRoute: add ifaddr error. %s", err)
		}
	}
}
