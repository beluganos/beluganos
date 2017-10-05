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
	"fabricflow/ribs/ribsmsg"
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	"gonla/nlalib"
)

//
// RIBS/RIC Service
//
type RicService struct {
	nid    uint8
	rt     string
	ifname string
	ric    *RicController
	cli    *CoreApiClient
}

func NewRicService(nid uint8, rt string, ifname string, apiAddr string) (*RicService, error) {
	client, err := NewCoreApiClient(apiAddr)
	if err != nil {
		return nil, err
	}

	return &RicService{
		nid:    nid,
		rt:     rt,
		ifname: ifname,
		ric:    NewRicController(),
		cli:    client,
	}, nil
}

func (s *RicService) Start(addr string, port uint16, rd string, label uint32) error {
	entry := NewRicEntryOnRic(s.nid, addr, port, s.rt, rd, label, nil)
	Tables.Rics.Add(entry)
	return s.ric.Start(addr)
}

func (s *RicService) Serve() {
	var micConnect bool = false
	var ricConnect bool = false
	for {
		select {
		case con := <-s.ric.Conn():
			log.Infof("RICS: RIC connected.")
			ricConnect = true
			if err := s.ric.Monitor(con.RemoteAddr); err != nil {
				log.Errorf("RICS: ric.Monitor error. %s", err)
			}
			if micConnect {
				if err := s.cli.SyncRib(s.rt); err != nil {
					log.Errorf("RICS: cli.SyncRib error. %s", err)
				}
			}

		case <-s.cli.Conn():
			log.Infof("RICS: MIC connected.")
			micConnect = true
			if err := s.cli.Monitor(s.nid, s.rt); err != nil {
				log.Errorf("RICS: cli.Monitor error. %s", err)
			}
			if ricConnect {
				if err := s.cli.SyncRib(s.rt); err != nil {
					log.Errorf("RICS: cli.SyncRib error. %s", err)
				}
			}

		case rib := <-s.ric.Recv():
			if err := s.SendToMic(rib); err != nil {
				log.Errorf("RICS: SendToMic error %s", err)
			}

		case rib := <-s.cli.Recv():
			if err := s.SendToRic(rib); err != nil {
				log.Errorf("RICS: SendToRic error %s", err)
			}
		}
	}
}

func (s *RicService) SendToMic(rib *ribsmsg.RibUpdate) error {

	f := func(vpnPath *table.Path, ipPath *table.Path, ric *RicEntry) error {
		log.Debugf("RICS: SendToMic %v", vpnPath)
		paths := []*table.Path{vpnPath}
		if err := s.cli.ModRib(ribsmsg.NewRibUpdate("", paths, ric.Rt)); err != nil {
			log.Errorf("RICS: cli.ModRib error. %s", err)
			return err
		}

		return nil
	}

	return s.ric.Translate(rib, f)
}

func (s *RicService) SendToMicAll() {
	go s.ric.GetRibs(s.rt)
}

func setDummyRoute(path *table.Path, ifname string) {
	dst := nlalib.NewIPNetFromIP(path.GetNexthop())
	err := func() error {
		if path.IsWithdraw {
			// log.Debugf("RICS: DelDummyRoute %s %s", dst, ifname)
			// return nlalib.DelIFAddr(dst, ifname)
			log.Debugf("RICS: Skip DelIFAddr %s %s", dst, ifname)
			return nil
		} else {
			log.Debugf("RICS: AddIFAddr %s %s", dst, ifname)
			return nlalib.AddIFAddr(dst, ifname)
		}
	}()
	if err != nil {
		log.Errorf("RICS: setDummyRoute error. %s %s %s", dst, ifname, err)
	}
}

func (s *RicService) SendToRic(rib *ribsmsg.RibUpdate) error {

	for _, path := range rib.Paths {
		log.Debugf("RICS: SendToRic %v", path)
		setDummyRoute(path, s.ifname)
		if err := Tables.Rics.FindByRt(s.rt).SendBgpPath(path); err != nil {
			log.Errorf("RICS: rics.SendBgpPath error. %s %s ", path, err)
		}
	}

	return nil
}
