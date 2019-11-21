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
	"fmt"
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
	linkCh     chan netlink.LinkUpdate
	confDB     *ConfigDB

	log *log.Entry
}

func NewServer() *Server {
	return &Server{
		nsInterval: NSIntervalDefault,
		workers:    map[string]Worker{},
		linkCh:     make(chan netlink.LinkUpdate),
		confDB:     NewConfigDB(),

		log: log.WithFields(log.Fields{"module": "server"}),
	}
}

func (s *Server) SetNSInterval(t time.Duration) {
	s.nsInterval = t
}

func (s *Server) SetConfig(path, typ, name string) error {
	cfg := NewConfig()
	cfg.SetConfigFile(path, typ)
	if err := cfg.Load(); err != nil {
		return err
	}

	c := cfg.Get(name)
	if c == nil {
		return fmt.Errorf("Config not found. name='%s'", name)
	}

	s.confDB.Update(c)
	for key, val := range s.confDB.Features() {
		s.log.Infof("feature: %s = %t", key, val)
	}

	return nil
}

func (s *Server) Start(done chan struct{}) error {

	if err := s.subscribeLinks(); err != nil {
		return err
	}

	if err := netlink.LinkSubscribe(s.linkCh, done); err != nil {
		return err
	}

	go s.Serve(done)
	return nil
}

func (s *Server) Serve(done chan struct{}) {
	s.log.Debugf("Serve: start")

	for {
		select {
		case m := <-s.linkCh:
			if ok := s.confDB.Has(m.Attrs().Name); ok {
				switch m.Header.Type {
				case unix.RTM_NEWLINK:
					s.startWorker(m.Link)

				case unix.RTM_DELLINK:
					s.stopWorker(m.Link)
				}
			}

		case <-done:
			s.log.Debugf("Serve: exit")
			return
		}
	}
}

func (s *Server) subscribeLinks() error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, link := range links {
		if ok := s.confDB.Has(link.Attrs().Name); ok {
			s.newWorker(link.Attrs().Name)
		}
	}

	return nil
}

func (s *Server) newWorker(ifname string) {
	w := NewWorker(ifname, s.nsInterval, s.confDB.Features())
	s.workers[ifname] = w
	w.Start()
	s.log.Debugf("Serve: Worker added. %s", ifname)
}

func (s *Server) startWorker(link netlink.Link) {
	ifname := link.Attrs().Name
	if _, ok := s.workers[ifname]; ok {
		s.log.Debugf("startWorker: Worker already exist. %s", ifname)
		return
	}

	s.newWorker(ifname)
}

func (s *Server) stopWorker(link netlink.Link) {
	ifname := link.Attrs().Name
	w, ok := s.workers[ifname]
	if !ok {
		s.log.Warnf("stopWorker: Worker not found. %s", ifname)
		return
	}

	w.Stop()

	delete(s.workers, ifname)
	s.log.Debugf("stopWorker: Worker deleted. %s", ifname)
}
