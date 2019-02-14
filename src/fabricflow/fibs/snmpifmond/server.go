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
	"net"
	"time"

	lib "fabricflow/fibs/fibslib"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Server struct {
	builder     *SnmpMessageBuilder
	serverAddr  *net.UDPAddr
	serverConn  net.Conn
	resendTime  time.Duration
	skipIfnames map[string]struct{}
	ifOid       []uint
	linkCh      chan netlink.LinkUpdate
}

const (
	TRAP_RESEND_INTERVAL = 3600 * time.Second
	TRAP_IFNUM_MAX       = 16
)

func NewServer(serverAddr *net.UDPAddr, ifoid []uint) *Server {
	server := &Server{
		builder:     NewSnmpMessageBuilder(lib.SNMP_VERSION, lib.SNMP_COMMUNITY),
		serverAddr:  serverAddr,
		resendTime:  TRAP_RESEND_INTERVAL,
		skipIfnames: map[string]struct{}{},
		ifOid:       ifoid,
		linkCh:      make(chan netlink.LinkUpdate),
	}
	server.RegisterSkipIfname("lo")
	return server
}

func (s *Server) RegisterSkipIfname(ifnames ...string) {
	for _, ifname := range ifnames {
		s.skipIfnames[ifname] = struct{}{}
	}
}

func (s *Server) SetResendInterval(t time.Duration) {
	s.resendTime = t
}

func (s *Server) Serve() {
	log.Infof("Server.Serve START")

	s.sendSnmpTrapAllLinks()
	log.Debugf("Server.Serve trap all ifaces.")

	ticker := time.NewTicker(s.resendTime)
	for {
		select {
		case link := <-s.linkCh:
			if status := linkUpdateStatus(&link); status != LinkStatusNone {
				log.Infof("%s %s", link.Link.Attrs().Name, status)
				log.Debugf("Link: %v", link.Link)
				log.Debugf("Type: %d Flags:%04x", link.Header.Type, link.IfInfomsg.Flags)

				s.sendSnmpTrap([]netlink.Link{link.Link})
			}

		case <-ticker.C:
			s.sendSnmpTrapAllLinks()
			log.Debugf("Server.Serve trap all ifaces.")
		}
	}
}

func (s *Server) Start() error {
	conn, err := net.DialUDP("udp", nil, s.serverAddr)
	if err != nil {
		log.Errorf("Server.Start Dial error. %s %s", s.serverAddr, err)
		return err
	}
	s.serverConn = conn

	go s.Serve()

	netlink.LinkSubscribe(s.linkCh, nil)

	return nil
}

func (s *Server) sendSnmpTrap(links []netlink.Link) {

	restLinks := links

	for {
		sendNum := len(restLinks)
		if sendNum == 0 {
			break
		}
		if sendNum > TRAP_IFNUM_MAX {
			sendNum = TRAP_IFNUM_MAX
		}

		msg, err := s.builder.NewLinksTrap(s.ifOid, restLinks[:sendNum])
		if err != nil {
			log.Errorf("Server.Serve NewLinksTrap error. %s", err)
			return
		}
		buf, err := s.builder.Encode(msg)
		if err != nil {
			log.Errorf("Server.Serve Encode error. %s", err)
			return
		}

		if _, err := s.serverConn.Write(buf); err != nil {
			log.Errorf("Server.Serve Write error. %s", err)
			return
		}

		restLinks = restLinks[sendNum:]
	}
}

func (s *Server) sendSnmpTrapAllLinks() {
	links, err := netlink.LinkList()
	if err != nil {
		log.Errorf("Server.sendAllLinks LinkList error. %s", err)
		return
	}

	trapLinks := []netlink.Link{}
	for _, link := range links {
		ifname := link.Attrs().Name
		if _, skip := s.skipIfnames[ifname]; !skip {
			trapLinks = append(trapLinks, link)
		}
	}

	s.sendSnmpTrap(trapLinks)
}
