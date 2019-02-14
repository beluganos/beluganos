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
	"fabricflow/fibc/api"
	"fabricflow/fibc/lib"
	"fabricflow/fibc/net"

	log "github.com/sirupsen/logrus"
)

//
// Server is main service of gonsld.
//
type Server struct {
	client     *fibcnet.Client
	dpID       uint64
	ONSLConfig string
	Unit       int
	Fields     *FieldGroups
	idmaps     *IDMaps
}

//
// NewServer cerates new instance of Server.
//
func NewServer(cfg *DpConfig) *Server {
	return &Server{
		client: fibcnet.NewClient(cfg.GetHost()),
		dpID:   cfg.DpID,
		Unit:   cfg.Unit,
		Fields: NewFieldGroups(cfg.Unit),
		idmaps: NewIDMaps(),
	}
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

	for {
		select {
		case connected := <-s.client.Conn():
			if connected {
				log.Infof("Server: connected.")

				hello := fibcapi.NewFFHello(s.dpID)
				if err := s.client.Write(hello, 0); err != nil {
					log.Errorf("Server: Write error. %s", err)
				}

			} else {
				log.Infof("Server: connection closed.")
			}

		case pkt := <-rxCh:
			log.Debugf("pkt  : %p len:%d tot:%d", pkt, pkt.PktLen(), pkt.TotLen())
			log.Debugf("unit : %d", pkt.Unit())
			log.Debugf("flags: %d", pkt.Flags())
			log.Debugf("cos  : %d", pkt.Cos())
			log.Debugf("vid  : %d", pkt.VID())
			log.Debugf("port : src:%d dst:%d", pkt.SrcPort(), pkt.DstPort())
			log.Debugf("rx   : port    : %d", pkt.RxPort())
			log.Debugf("rx   : untagged: %d", pkt.RxUntagged())
			log.Debugf("rx   : matched : %d", pkt.RxMatched())
			log.Debugf("rx   : reasons : %d", pkt.RxReasons())
			log.Debugf("blk  : #%d", pkt.BlkCount())

			rxBuf.Reset()
			if _, err := pkt.WriteTo(&rxBuf); err != nil {
				log.Errorf("Server: pkt.WriteTo error. %s", err)
			} else {
				pktIn := fibcapi.NewFFPacketIn(s.dpID, uint32(pkt.SrcPort()), rxBuf.Bytes())
				if err := s.client.Write(pktIn, 0); err != nil {
					log.Errorf("Server: client.Write error. %s", err)
				}
			}

		case linkInfo := <-linkCh:
			log.Debugf("Server: LinkInfo: %v", linkInfo)

			port := fibcapi.NewFFPort(linkInfo.PortNo())
			port.State = linkInfo.PortState()
			portStatus := fibcapi.NewFFPortStatus(s.dpID, port, fibcapi.FFPortStatus_MODIFY)
			if err := s.client.Write(portStatus, 0); err != nil {
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
	go s.RecvMain()
	go s.Serve(done)

	log.Infof("Server: started.")
	return nil
}
