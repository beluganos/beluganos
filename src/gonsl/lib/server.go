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

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

//
// Server is main service of gonsld.
//
type Server struct {
	client    FIBController
	dpCfg     *DpConfig
	logCfg    *LogConfig
	fields    *FieldGroups
	idmaps    *IDMaps
	vlanPorts *VlanPortTable
	l2addrCh  chan []*L2addrmonEntry

	log *log.Entry
}

//
// NewServer cerates new instance of Server.
//
func NewServer(dpCfg *DpConfig, logCfg *LogConfig) *Server {
	return &Server{
		client:    NewFIBController(dpCfg.GetFIBCType(), dpCfg.GetHost(), dpCfg.DpID),
		dpCfg:     dpCfg,
		logCfg:    logCfg,
		fields:    NewFieldGroups(dpCfg.Unit),
		idmaps:    NewIDMaps(),
		vlanPorts: NewVlanPortTableFromConfig(&dpCfg.BlockBcast),
		l2addrCh:  make(chan []*L2addrmonEntry),

		log: log.WithFields(log.Fields{"module": "server"}),
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

func (s *Server) recvMessages(done <-chan struct{}) {
	s.log.Debugf("Receiver: started")

FOR_LOOP:
	for {
		select {
		case msg := <-s.client.Recv():
			if err := msg.Dispatch(s); err != nil {
				s.log.Warnf("Serve: Dispatch error. %s", err)
			}

		case <-done:
			s.log.Debugf("Receiver: exit.")
			break FOR_LOOP
		}
	}
}

func (s *Server) sendMessages(done <-chan struct{}) {
	s.log.Debugf("Sender: started")

	var rxBuf bytes.Buffer

	defer s.client.Stop()
	rxCh := s.RxStart(done)
	linkCh := s.LinkmonStart(done)
	s.L2AddrMonStart(done)

FOR_LOOP:
	for {
		select {
		case connected := <-s.client.Conn():
			if connected {
				s.log.Infof("Sender: connected.")

				hello := fibcapi.NewFFHello(s.DpID())
				if err := s.client.Hello(hello); err != nil {
					s.log.Errorf("Sender: Write error. %s", err)
				}

			} else {
				s.log.Infof("Sender: connection closed.")
			}

		case pkt := <-rxCh:
			s.dumpRxPkt(pkt)

			rxBuf.Reset()
			if _, err := pkt.WriteTo(&rxBuf); err != nil {
				s.log.Errorf("Sender: pkt.WriteTo error. %s", err)
			} else {
				pktIn := fibcapi.NewFFPacketIn(s.DpID(), uint32(pkt.SrcPort()), rxBuf.Bytes())
				if err := s.client.PacketIn(pktIn); err != nil {
					s.log.Errorf("Sender: client.Write error. %s", err)
				}
			}

		case linkInfo := <-linkCh:
			s.log.Debugf("Sender: LinkInfo: %v", linkInfo)

			port := fibcapi.NewFFPort(linkInfo.PortNo())
			port.State = linkInfo.PortState()
			portStatus := fibcapi.NewFFPortStatus(s.DpID(), port, fibcapi.FFPortStatus_MODIFY)
			if err := s.client.PortStatus(portStatus); err != nil {
				s.log.Errorf("Sender: client.Write error. %s", err)
			}

		case entries := <-s.l2addrCh:
			s.log.Debugf("Sender: L2AddrEntry %d", len(entries))

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
			if err := s.client.L2AddrStatus(status); err != nil {
				s.log.Errorf("Sender: client.Write error. %s", err)
			}

		case <-done:
			s.log.Infof("Sender: Exit.")
			break FOR_LOOP
		}
	}
}

//
// Start starts submodules.
//
func (s *Server) Start(done <-chan struct{}) error {
	s.log.Infof("Start:")

	if block := s.dpCfg.BlockBcast.Block(); !block {
		if err := PortDefaultVlanConfig(s.Unit()); err != nil {
			s.log.Errorf("Start: PortDefaultVlanConfig error. %s", err)
			return err
		}

		s.log.Infof("Start: PortDefaultVlanConfig ok.")
	}

	s.log.Infof("FIBCConroller: %s", s.client)

	go s.client.Start()
	go s.sendMessages(done)
	go s.recvMessages(done)

	return nil
}
