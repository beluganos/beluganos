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
	"fmt"
	"gonla/nladbm"
	"gonla/nlamsg"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
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

const (
	RECEIVE_BUFFER_SIZE = 0x1000
)

type NLARecvBuffer struct {
	Buffer []byte
	Len    int
}

func NewNLARecvBuffer() *NLARecvBuffer {
	return &NLARecvBuffer{
		Buffer: make([]byte, RECEIVE_BUFFER_SIZE),
	}
}

func (b *NLARecvBuffer) Recvfrom(fd int) (err error) {
	if b.Len, _, err = unix.Recvfrom(fd, b.Buffer, 0); err != nil {
		return
	}

	if b.Len < unix.NLMSG_HDRLEN {
		err = fmt.Errorf("too short msg. %d", b.Len)
	}

	return
}

func (b *NLARecvBuffer) Bytes() []byte {
	return b.Buffer[:b.Len]
}

type NLAServer struct {
	Nid             uint8
	nlmsgs          chan<- *nlamsg.NetlinkMessage
	done            chan struct{}
	recvChanSize    int
	recvSockBufSize int
	log             *log.Entry
}

func NewNLAServer(nid uint8, nlmsgs chan<- *nlamsg.NetlinkMessage, done chan struct{}) *NLAServer {
	fields := log.Fields{
		"module": "NLAServer",
	}

	return &NLAServer{
		Nid:             nid,
		nlmsgs:          nlmsgs,
		done:            done,
		recvChanSize:    16,
		recvSockBufSize: 1024 * 1024,
		log:             log.WithFields(fields),
	}
}

func (s *NLAServer) SetRecvChanSize(chanSize int) {
	s.recvChanSize = chanSize
}

func (s *NLAServer) SetRecvSockBufferSize(sockBufSize int) {
	s.recvSockBufSize = sockBufSize
}

func (s *NLAServer) parseNlMsgs(recvCh <-chan *NLARecvBuffer) {

	statCount := nladbm.Stats().New("NLAServer.nlmsg.count")

	for rb := range recvCh {
		nlmsgs, err := syscall.ParseNetlinkMessage(rb.Bytes())
		if err != nil {
			s.log.Errorf("parseNlMsg: parse error. %s", err)
			return
		}

		for _, nlmsg := range nlmsgs {
			s.nlmsgs <- nlamsg.NewNetlinkMessage(&nlmsg, s.Nid, nlamsg.SRC_KNL)
			statCount.Inc()
		}
	}
}

func (s *NLAServer) Serve(sock *nl.NetlinkSocket) {

	statRecv := nladbm.Stats().New("NLAServer.nlmsg.recv")
	statRetry := nladbm.Stats().New("NLAServer.nlmsg.e_retry")
	statNoBufs := nladbm.Stats().New("NLAServer.nlmsg.e_nobufs")

	recvCh := make(chan *NLARecvBuffer, s.recvChanSize)

	go s.parseNlMsgs(recvCh)

	fd := sock.GetFd()
	for {
		rb := NewNLARecvBuffer()
		if err := rb.Recvfrom(fd); err != nil {
			switch err {
			case unix.EINTR, unix.EAGAIN:
				statRetry.Inc()
				s.log.Debugf("Serve: Recvfrom retry. %s", err)
				continue

			case unix.ENOBUFS:
				statNoBufs.Inc()
				s.log.Warnf("Serve: Recvfrom error. %s", err)
				continue

			default:
				s.log.Errorf("Serve: EXIT. Recvfrom error. %s", err)
				return
			}
		}

		statRecv.Inc()

		recvCh <- rb
	}
}

func (s *NLAServer) Start() error {
	sock, err := nl.Subscribe(syscall.NETLINK_ROUTE, RTNLGRPLIST...)
	if err != nil {
		s.log.Errorf("Start: subscribe error. %s", err)
		if s.done != nil {
			close(s.done)
		}
		return err
	}

	if err := sock.SetReceiveBuffer(s.recvSockBufSize); err != nil {
		s.log.Errorf("Start: sock.SetReceiveBuffer error. %s", err)
		sock.Close()
		if s.done != nil {
			close(s.done)
		}
		return err
	}

	n, _ := sock.GetReceiveBuffer()
	s.log.Infof("Start: socketopt(RCVBUF) = %d", n)

	if s.done != nil {
		go func() {
			<-s.done
			sock.Close()
		}()
	}

	go s.Serve(sock)
	log.Info("Start:")
	return nil
}
