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

package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gonla/nlaapi"
	"gonla/nlalib"
	"gonla/nlamsg"
)

type NLAController struct {
	recvCh chan *nlamsg.NetlinkMessageUnion
	client nlaapi.NLAApiClient
}

func NewNLAController() *NLAController {
	return &NLAController{
		recvCh: make(chan *nlamsg.NetlinkMessageUnion),
		client: nil,
	}
}

func (n *NLAController) Recv() <-chan *nlamsg.NetlinkMessageUnion {
	return n.recvCh
}

func (n *NLAController) Start(addr string) error {
	ch := make(chan *nlalib.ConnInfo)
	conn, err := nlalib.NewClientConn(addr, ch)
	if err != nil {
		close(ch)
		return err
	}
	n.client = nlaapi.NewNLAApiClient(conn)

	go func() {
		for {
			ci := <-ch
			go n.Monitor(ci)
		}
	}()

	return nil
}

func (n *NLAController) Monitor(ci *nlalib.ConnInfo) {
	stream, err := n.client.MonNetlink(context.Background(), &nlaapi.MonNetlinkRequest{})
	if err != nil {
		log.Errorf("NLAController: Monitor error. %s", err)
		return
	}

	log.Infof("NLAController: Monitor START")

	for {
		nlmsg, err := stream.Recv()
		if err != nil {
			log.Infof("NLAController: Monitor EXIT. %s", err)
			break
		}

		log.Debugf("NLAController: Monitor %v", nlmsg)
		n.recvCh <- nlmsg.ToNative()
	}
}
