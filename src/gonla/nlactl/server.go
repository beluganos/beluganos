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

package nlactl

import (
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink/nl"
	"gonla/nlamsg"
	"syscall"
)

var RTNLGRPLIST = []uint{
	syscall.RTNLGRP_LINK,
	syscall.RTNLGRP_NEIGH,
	syscall.RTNLGRP_IPV4_IFADDR,
	syscall.RTNLGRP_IPV4_ROUTE,
	syscall.RTNLGRP_IPV6_IFADDR,
	syscall.RTNLGRP_IPV6_ROUTE,
	nl.RTNLGRP_MPLS_ROUTE,
}

type NLAServer struct {
	Nid    uint8
	nlmsgs chan<- *nlamsg.NetlinkMessage
	done   chan struct{}
}

func NewNLAServer(nid uint8, nlmsgs chan<- *nlamsg.NetlinkMessage, done chan struct{}) *NLAServer {
	return &NLAServer{
		Nid:    nid,
		nlmsgs: nlmsgs,
		done:   done,
	}
}

func (s *NLAServer) Serve(sock *nl.NetlinkSocket) {
	for {
		nlmsgs, err := sock.Receive()
		if err != nil {
			log.Errorf("NLAServer EXIT. %s", err)
			return
		}

		for _, nlmsg := range nlmsgs {
			s.nlmsgs <- nlamsg.NewNetlinkMessage(&nlmsg, s.Nid, nlamsg.SRC_KNL)
		}
	}
}

func (s *NLAServer) Start() error {
	sock, err := nl.Subscribe(syscall.NETLINK_ROUTE, RTNLGRPLIST...)
	if err != nil {
		log.Errorf("NLAServer: nl.SubscribeAt error. %s", err)
		close(s.done)
		return err
	}

	if s.done != nil {
		go func() {
			<-s.done
			sock.Close()
		}()
	}

	go s.Serve(sock)
	log.Info("NLAServer: START")
	return nil
}
