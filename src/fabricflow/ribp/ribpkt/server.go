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

package ribpkt

import (
	"fabricflow/fibc/net"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"syscall"
	"time"
)

type Server struct {
	interval time.Duration
	reId     string
	vrf      uint8
	links    map[int]netlink.Link // key: ifindex
	linkCh   chan netlink.LinkUpdate
	ctrlCh   chan string
}

func NewServer(reId string, vrf uint8, interval int) *Server {
	return &Server{
		interval: time.Duration(interval),
		reId:     reId,
		vrf:      vrf,
		links:    make(map[int]netlink.Link),
		linkCh:   make(chan netlink.LinkUpdate),
		ctrlCh:   make(chan string),
	}
}

func (s *Server) CtrlCh() chan<- string {
	return s.ctrlCh
}

func (s *Server) LinkCh() chan<- netlink.LinkUpdate {
	return s.linkCh
}

func (s *Server) Start(done <-chan struct{}) error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, link := range links {
		s.AddLink(link)
	}

	go s.Serve(done)
	return nil
}

func (s *Server) Serve(done <-chan struct{}) {

	ticker := func() *time.Ticker {
		interval := s.interval
		if interval == 0 {
			interval = 3600 * 1000
		}
		return time.NewTicker(interval * time.Millisecond)
	}()

	for {
		select {
		case link := <-s.linkCh:
			if link.Header.Type == syscall.RTM_NEWLINK {
				s.AddLink(link)
				log.Infof("Add: %v %v", link.Header, link.Link)
			} else {
				s.DelLink(link)
				log.Infof("Del: %v %v", link.Header, link.Link)
			}
		case <-s.ctrlCh:
			s.SendPackets()
		case <-ticker.C:
			if s.interval > 0 {
				s.SendPackets()
			}
		case <-done:
			return
		}
	}
}

func (s *Server) AddLink(link netlink.Link) {
	if link.Type() != "device" && link.Type() != "veth" {
		return
	}

	if link.Attrs().Name == "lo" {
		return
	}

	s.links[link.Attrs().Index] = link
	log.Infof("ADD %s", link.Attrs().Name)
}

func (s *Server) DelLink(link netlink.Link) {
	delete(s.links, link.Attrs().Index)
	log.Infof("DEL %s", link.Attrs().Name)
}

func (s *Server) SendPackets() {
	for _, link := range s.links {
		if err := s.SendPacket(link.Attrs()); err != nil {
			log.Errorf("Send error. %s", err)
		}
	}
}

func (s *Server) SendPacket(attrs *netlink.LinkAttrs) error {

	ifname := fmt.Sprintf("%d/%s", s.vrf, attrs.Name)

	data, err := fibcnet.NewFFPacket(s.reId, attrs.HardwareAddr, ifname).Bytes()
	if err != nil {
		return err
	}

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	sa := fibcnet.SockAddr(attrs.Index)
	if err := syscall.Sendto(fd, data, 0, sa); err != nil {
		return err
	}

	log.Debugf("SEND %s", ifname)

	return nil
}
