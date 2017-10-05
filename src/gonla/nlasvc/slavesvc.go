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

package nlasvc

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gonla/nlaapi"
	"gonla/nlactl"
	"gonla/nlalib"
	"gonla/nlamsg"
)

type NLASlaveService struct {
	client nlaapi.NLACoreApiClient
	NId    uint8
	Addr   string
	Chans  *nlactl.NLAChannels
}

func NewNLASlaveService(addr string) *NLASlaveService {
	p := &NLASlaveService{
		client: nil,
		NId:    0,
		Addr:   addr,
		Chans:  nil,
	}
	return p
}

func (n *NLASlaveService) Start(nid uint8, chans *nlactl.NLAChannels) error {
	n.NId = nid
	connected, err := n.Connect(n.Addr)
	if err != nil {
		return err
	}

	n.Chans = chans

	go n.Serve(connected)
	log.Infof("SlaveService: START")
	return nil
}

func (n *NLASlaveService) Stop() {
	// nothing to do
	log.Infof("SlaveService: STOP")
}

func (n *NLASlaveService) Connect(addr string) (<-chan *nlalib.ConnInfo, error) {
	connected := make(chan *nlalib.ConnInfo)
	conn, err := nlalib.NewClientConn(addr, connected)
	if err != nil {
		close(connected)
		return nil, err
	}

	n.client = nlaapi.NewNLACoreApiClient(conn)
	log.Infof("SlaveService: Connected %s", addr)

	return connected, nil
}

func (n *NLASlaveService) Serve(connected <-chan *nlalib.ConnInfo) {
	for {
		select {
		case ci := <-connected:
			go n.MonNetlinkMessage(ci)
		}
	}
}

func (n *NLASlaveService) MonNetlinkMessage(ci *nlalib.ConnInfo) {

	req := &nlaapi.Node{
		NId: uint32(n.NId),
		Ip:  ci.LocalAddr,
	}
	stream, err := n.client.MonNetlinkMessage(context.Background(), req)
	if err != nil {
		log.Errorf("SlaveService: MonNetlinkMessage error. %s", err)
		return
	}

	log.Infof("SlaveService: MonNetlinkMessage START")

	SubscribeNetlinkResources(n.Chans.NlMsg, n.NId)

	for {
		nlmsg, err := stream.Recv()
		if err != nil {
			log.Errorf("SlaveService: MonNetlinkMessage EXIT. %s", err)
			return
		}

		n.Chans.Api <- nlmsg.ToNative()
	}
}

func (n *NLASlaveService) NetlinkMessage(nlmsg *nlamsg.NetlinkMessage) {
	if nlmsg.Src != nlamsg.SRC_KNL {
		return
	}

	msg := nlaapi.NewNetlinkMessageFromNative(nlmsg)
	if _, err := n.client.SendNetlinkMessage(context.Background(), msg); err != nil {
		log.Errorf("SlaveService: SendNetlinkMessage error. %s", err)
		return
	}

	log.Debugf("SlaveService: Send to master. %v", &msg.Header)
}
