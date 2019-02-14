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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const (
	NSIntervalDefault = 15 * time.Minute
)

type Server struct {
	nsInterval time.Duration
	workers    map[string]Worker
}

func (s *Server) SetNSInterval(t time.Duration) {
	s.nsInterval = t
}

func NewServer() *Server {
	return &Server{
		nsInterval: NSIntervalDefault,
		workers:    map[string]Worker{},
	}
}

func (s *Server) Start(done chan struct{}) error {

	if err := s.subscribeLinks(); err != nil {
		return err
	}

	ch := make(chan *LinkUpdate)

	if err := MonitorLink(ch, done); err != nil {
		return err
	}

	go s.Serve(ch, done)
	return nil
}

func (s *Server) Serve(ch <-chan *LinkUpdate, done chan struct{}) {
	log.Debugf("Serve: start")

	for {
		select {
		case m := <-ch:
			if ok := checkLinkType(m.Link); ok {
				switch m.Type {
				case unix.RTM_NEWLINK:
					s.startWorker(m)

				case unix.RTM_DELLINK:
					s.stopWorker(m)
				}
			}

		case <-done:
			log.Debugf("Serve: exit")
			return
		}
	}
}

func (s *Server) newWorker(ifname string) {
	w := NewWorker(ifname, s.nsInterval)
	if err := w.Start(); err != nil {
		log.Errorf("Serve: Worker.Start error. %s %s", ifname, err)
		return
	}

	s.workers[ifname] = w
	log.Debugf("Serve: Worker added. %s", ifname)
}

func (s *Server) subscribeLinks() error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, link := range links {
		if ok := checkLinkType(link); ok {
			s.newWorker(link.Attrs().Name)
		}
	}

	return nil
}

func (s *Server) startWorker(m *LinkUpdate) {
	ifname := m.Link.Attrs().Name
	if _, ok := s.workers[ifname]; ok {
		log.Debugf("Serve: Worker already exist. %s", ifname)
		return
	}

	time.AfterFunc(3*time.Second, func() {
		s.newWorker(ifname)
	})
}

func (s *Server) stopWorker(m *LinkUpdate) {
	ifname := m.Link.Attrs().Name
	w, ok := s.workers[ifname]
	if !ok {
		log.Warnf("Serve: Worker not found. %s", ifname)
		return
	}

	w.Stop()

	delete(s.workers, ifname)
	log.Debugf("Serve: Worker deleted. %s", ifname)
}
