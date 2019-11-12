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

package govsw

import (
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
)

type Listener interface {
	PacketIn(*Packet)
}

type Server struct {
	Listener Listener
	DB       *DB
	SyncCh   <-chan string

	linkMon *LinkMonitor
	pktInCh chan *Packet

	log *log.Entry
}

func (s *Server) init() {
	s.linkMon = &LinkMonitor{}
	s.pktInCh = make(chan *Packet)
	s.log = log.WithFields(log.Fields{"module": "server"})
}

func (s *Server) syncLinks(targetIfname string) {
	links, err := netlink.LinkList()
	if err != nil {
		s.log.Errorf("list link error. %s", err)
		return
	}

FOR_LOOP:
	for _, link := range links {
		ifname := link.Attrs().Name
		if len(targetIfname) > 0 && targetIfname != "*" {
			if ifname != targetIfname {
				continue FOR_LOOP
			}
		}

		upd := netlink.LinkUpdate{}
		upd.Link = link
		upd.Header.Type = func() uint16 {
			if has := s.DB.Ifname().Has(ifname); has {
				return unix.RTM_NEWLINK
			}
			return unix.RTM_DELLINK
		}()

		s.linkMon.Put(upd)
	}
}

func (s *Server) Serve(done <-chan struct{}) {
FOR_LOOP:
	for {
		select {
		case upd := <-s.linkMon.Recv():
			state := upd.Link.Attrs().OperState
			ifname := upd.Link.Attrs().Name

			switch upd.Header.Type {
			case unix.RTM_NEWLINK:
				if has := s.DB.Ifname().Has(ifname); !has {
					s.log.Debugf("Serve: NEWLINK %s not registered.", ifname)
					continue FOR_LOOP
				}

				s.log.Debugf("Serve: NEWLINK %s", ifname)

				s.DB.Link().GetOrAdd(upd.Link, func(link *Link) {
					if state == netlink.OperUp {
						s.log.Debugf("Serve: NEWLINK/Up %s", ifname)

						if err := link.Start(s.pktInCh); err != nil {
							s.log.Errorf("Serve: Link Start error. %s", err)
						}

					} else if state == netlink.OperDown {
						s.log.Debugf("Serve: NEWLINK/Down %s", ifname)

						link.Stop()

					} else {
						s.log.Debugf("Serve: NEWLINK/%d %s", state, ifname)
					}
				})

			case unix.RTM_DELLINK:
				if link, ok := s.DB.Link().Delete(upd.Link.Attrs().Index); ok {
					s.log.Debugf("Serve: DELLINK %s", ifname)

					link.Destroy()
				}

			default:
				s.log.Debugf("Serve: Bad type. %d", upd.Header.Type)
			}

		case ifname := <-s.SyncCh:
			go s.syncLinks(ifname)

		case pkt := <-s.pktInCh:
			s.Listener.PacketIn(pkt)

		case <-done:
			s.log.Infof("Server: exit(done).")
			break FOR_LOOP
		}
	}
}

func (s *Server) Start(done <-chan struct{}) error {
	s.init()

	if err := s.linkMon.Start(done); err != nil {
		s.log.Errorf("Start: link reload error. %s", err)
		return err
	}

	go s.Serve(done)

	s.log.Infof("Start: success.")

	return nil
}
