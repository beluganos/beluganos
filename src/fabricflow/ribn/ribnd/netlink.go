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
	"net"

	"golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type LinkUpdate struct {
	Link netlink.Link
	Type uint16
}

func (m *LinkUpdate) String() string {
	return fmt.Sprintf("%d %s", m.Type, m.Link.Attrs().Name)
}

func MonitorLink(ch chan<- *LinkUpdate, done chan struct{}) error {
	subscr := make(chan netlink.LinkUpdate)
	if err := netlink.LinkSubscribe(subscr, done); err != nil {
		return nil
	}

	go func() {
		defer close(subscr)

		log.Debugf("MonitorLink start.")

		for {
			select {
			case m := <-subscr:
				log.Debugf("IfInfoMsg: %v", m.IfInfomsg)
				log.Debugf("NlMsgHdr : %v", m.Header)

				switch m.Header.Type {
				case unix.RTM_NEWLINK, unix.RTM_DELLINK:
					ch <- &LinkUpdate{
						Link: m.Link,
						Type: m.Header.Type,
					}
				}

			case <-done:
				log.Debugf("MonitorLink exit.")
				return
			}
		}
	}()

	return nil
}

func checkLinkType(link netlink.Link) bool {
	if flags := link.Attrs().Flags; (flags & net.FlagLoopback) != 0 {
		return false
	}
	if linkType := link.Type(); linkType != "device" && linkType != "veth" {
		return false
	}

	return true
}
