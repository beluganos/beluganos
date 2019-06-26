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

package gonslib

import (
	"bytes"
	fibcapi "fabricflow/fibc/api"
	fibclib "fabricflow/fibc/lib"
	fibcnet "fabricflow/fibc/net"

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

//
// Server is main service of gonsld.
//
type Server struct {
	client    *fibcnet.Client
	dpCfg     *DpConfig
	logCfg    *LogConfig
	fields    *FieldGroups
	idmaps    *IDMaps
	vlanPorts *VlanPortTable
	l2addrCh  chan []*L2addrmonEntry
}

//
// NewServer cerates new instance of Server.
//
func NewServer(dpCfg *DpConfig, logCfg *LogConfig) *Server {
	return &Server{
		client:    fibcnet.NewClient(dpCfg.GetHost()),
		dpCfg:     dpCfg,
		logCfg:    logCfg,
		fields:    NewFieldGroups(dpCfg.Unit),
		idmaps:    NewIDMaps(),
		vlanPorts: NewVlanPortTableFromConfig(&dpCfg.BlockBcast),
		l2addrCh:  make(chan []*L2addrmonEntry),
	}
}

func (s *Server) Unit() int {
	return s.dpCfg.Unit
}

func (s *Server) DpID() uint64 {
	return s.dpCfg.DpID
}

func (s *Server) Fields() *FieldGroups {
	return s.fields
}

func (s *Server) LogConfig() *LogConfig {
	return s.logCfg
}

func (s *Server) VlanPorts() *VlanPortTable {
	return s.vlanPorts
}

//
// RecvMain reveives and dispatch messages from fibcd.
//
func (s *Server) RecvMain() {
	log.Debugf("Server: RecvMain started.")

	s.client.Start(func(client *fibcnet.Client) {
		for {
			hdr, data, err := s.client.Read()
			if err != nil {
				log.Errorf("Server: Client EXIT. Read error. %s", err)
				return
			}

			log.Debugf("Server: Recv %v", hdr)

			if err := fibclib.Dispatch(hdr, data, s); err != nil {
				log.Errorf("Dispatch error. %s", err)
			}
		}
	})
}

//
// Serve is main loop of Service.
//
func (s *Server) Serve(done <-chan struct{}) {
	log.Debugf("Server: Serve started.")

	var rxBuf bytes.Buffer

	defer s.client.Stop()
	rxCh := s.RxStart(done)
	linkCh := s.LinkmonStart(done)
	s.L2AddrMonStart(done)

	for {
		select {
		case connected := <-s.client.Conn():
			if connected {
				log.Infof("Server: connected.")

				hello := fibcapi.NewFFHello(s.DpID())
				if err := s.client.Write(hello, 0); err != nil {
					log.Errorf("Server: Write error. %s", err)
				}

			} else {
				log.Infof("Server: connection closed.")
			}

		case pkt := <-rxCh:
			s.dumpRxPkt(pkt)

			rxBuf.Reset()
			if _, err := pkt.WriteTo(&rxBuf); err != nil {
				log.Errorf("Server: pkt.WriteTo error. %s", err)
			} else {
				pktIn := fibcapi.NewFFPacketIn(s.DpID(), uint32(pkt.SrcPort()), rxBuf.Bytes())
				if err := s.client.Write(pktIn, 0); err != nil {
					log.Errorf("Server: client.Write error. %s", err)
				}
			}

		case linkInfo := <-linkCh:
			log.Debugf("Server: LinkInfo: %v", linkInfo)

			port := fibcapi.NewFFPort(linkInfo.PortNo())
			port.State = linkInfo.PortState()
			portStatus := fibcapi.NewFFPortStatus(s.DpID(), port, fibcapi.FFPortStatus_MODIFY)
			if err := s.client.Write(portStatus, 0); err != nil {
				log.Errorf("Server: client.Write error. %s", err)
			}

		case entries := <-s.l2addrCh:
			log.Debugf("Server: L2AddrEntry %d", len(entries))

			reason := func(oper opennsl.L2CallbackOper) fibcapi.L2Addr_Reason {
				switch oper {
				case opennsl.L2_CALLBACK_ADD:
					return fibcapi.L2Addr_ADD

				case opennsl.L2_CALLBACK_DELETE:
					return fibcapi.L2Addr_DELETE

				default:
					return fibcapi.L2Addr_NOP
				}
			}

			addrs := make([]*fibcapi.L2Addr, len(entries))
			for index, entry := range entries {
				addrs[index] = fibcapi.NewL2AddrDP(
					entry.L2Addr.MAC(),
					uint16(entry.L2Addr.VID()),
					uint32(entry.L2Addr.Port()),
					reason(entry.Oper),
				)
			}

			status := fibcapi.NewFFL2AddrStatus(s.DpID(), addrs)
			if err := s.client.Write(status, 0); err != nil {
				log.Errorf("Server: client.Write error. %s", err)
			}

		case <-done:
			log.Infof("Server: Exit.")
			return
		}
	}
}

//
// Start starts submodules.
//
func (s *Server) Start(done <-chan struct{}) error {
	if block := s.dpCfg.BlockBcast.Block(); !block {
		if err := PortDefaultVlanConfig(s.Unit()); err != nil {
			log.Errorf("Server: PortDefaultVlanConfig error. %s", err)
			return err
		}

		log.Infof("Server: PortDefaultVlanConfig ok.")
	}

	go s.RecvMain()
	go s.Serve(done)

	log.Infof("Server: started.")
	return nil
}
