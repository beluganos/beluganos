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
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const (
	LINK_MON_CHAN_SIZE = 512
)

type LinkMonitor struct {
	ch  chan netlink.LinkUpdate
	log *log.Entry
}

func (m *LinkMonitor) init() {
	m.ch = make(chan netlink.LinkUpdate, LINK_MON_CHAN_SIZE)
	m.log = log.WithFields(log.Fields{"module": "linkmon"})
}

func (m *LinkMonitor) linkMonitorSendAll() {
	links, err := netlink.LinkList()
	if err != nil {
		m.log.Errorf("LinkList error. %s", err)
		return
	}

	for _, link := range links {
		linkUpdate := netlink.LinkUpdate{}
		linkUpdate.Link = link
		linkUpdate.Header.Type = unix.RTM_NEWLINK
		m.ch <- linkUpdate
	}
}

func (m *LinkMonitor) Put(upd netlink.LinkUpdate) {
	m.ch <- upd
}

func (m *LinkMonitor) Recv() <-chan netlink.LinkUpdate {
	return m.ch
}

func (m *LinkMonitor) Start(done <-chan struct{}) error {
	m.init()

	m.linkMonitorSendAll()

	if err := netlink.LinkSubscribe(m.ch, done); err != nil {
		m.log.Errorf("Start: LinkSubscribe error. %s", err)
		close(m.ch)
		return err
	}

	m.log.Infof("Start: success.")

	return nil
}
