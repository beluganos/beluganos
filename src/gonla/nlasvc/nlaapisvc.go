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
	"gonla/nlactl"
	"gonla/nladbm"
	"gonla/nlamsg"
)

//
// NLA API Service
//
type NLAApiService struct {
	NId    uint8
	Addr   string
	NlMsgs chan *nlamsg.NetlinkMessageUnion
}

func NewNLAApiService(addr string) *NLAApiService {
	return &NLAApiService{
		NId:    0,
		Addr:   addr,
		NlMsgs: make(chan *nlamsg.NetlinkMessageUnion),
	}
}

func (n *NLAApiService) Start(nid uint8, chans *nlactl.NLAChannels) error {
	n.NId = nid
	s := NewNLAApiServer(n.Addr, n.NId)
	if err := s.Start(chans.Api); err != nil {
		return err
	}

	log.Infof("NLAApiService: START")
	return nil
}

func (n *NLAApiService) Stop() {
	log.Infof("NLAApiService: STOP")
}

func (n *NLAApiService) NetlinkMessage(nlmsg *nlamsg.NetlinkMessage) {
	//log.Debugf("NLAApiService: NlMsg")
}

func (n *NLAApiService) sendToClients(nlmsg *nlamsg.NetlinkMessage, m interface{}) {
	msg := nlamsg.NewNetlinkMessageUnion(&nlmsg.Header, m, nlmsg.NId, nlmsg.Src)
	nladbm.Clients().Send(msg)
}

func (n *NLAApiService) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {
	// log.Debugf("NLAApiService: LINK")
	n.sendToClients(nlmsg, link)
}

func (n *NLAApiService) NetlinkAddr(nlmsg *nlamsg.NetlinkMessage, addr *nlamsg.Addr) {
	// log.Debugf("NLAApiService: ADDR")
	n.sendToClients(nlmsg, addr)
}

func (n *NLAApiService) NetlinkNeigh(nlmsg *nlamsg.NetlinkMessage, neigh *nlamsg.Neigh) {
	// log.Debugf("NLAApiService: NEIG")
	n.sendToClients(nlmsg, neigh)
}

func (n *NLAApiService) NetlinkRoute(nlmsg *nlamsg.NetlinkMessage, route *nlamsg.Route) {
	// log.Debugf("NLAApiService: ROUT")
	n.sendToClients(nlmsg, route)
}

func (n *NLAApiService) NetlinkNode(nlmsg *nlamsg.NetlinkMessage, node *nlamsg.Node) {
	// log.Debugf("NLAApiService: NODE")
	n.sendToClients(nlmsg, node)
}

func (n *NLAApiService) NetlinkVpn(nlmsg *nlamsg.NetlinkMessage, vpn *nlamsg.Vpn) {
	// log.Debugf("NLAApiService: VPN")
	n.sendToClients(nlmsg, vpn)
}
