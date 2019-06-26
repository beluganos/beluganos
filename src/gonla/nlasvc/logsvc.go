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
	"gonla/nlactl"
	"gonla/nlalib"
	"gonla/nlamsg"

	log "github.com/sirupsen/logrus"
)

type NLALogService struct {
	dump uint32
	log  *log.Entry
}

func NewNLALogService(dump uint32) *NLALogService {
	return &NLALogService{
		dump: dump,
		log:  NewLogger("NLALogService"),
	}
}

func (n *NLALogService) Start(uint8, *nlactl.NLAChannels) error {
	// nothing to do.
	return nil
}

func (n *NLALogService) Stop() {
	// nothing to do.
}

func (n *NLALogService) NetlinkMessage(nlmsg *nlamsg.NetlinkMessage) {
	n.log.Debugf("NLM : %s %v", nlmsg.Src, nlmsg)
	if n.dump != 0 {
		for _, line := range nlalib.HexSlice(nlmsg.Data) {
			n.log.Debugf("NLM : %s", line)
		}
	}
}

func (n *NLALogService) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {
	n.log.Debugf("LINK: %s %v %v", nlmsg.Src, &nlmsg.Header, link)
}

func (n *NLALogService) NetlinkAddr(nlmsg *nlamsg.NetlinkMessage, addr *nlamsg.Addr) {
	n.log.Debugf("ADDR: %s %v %v", nlmsg.Src, &nlmsg.Header, addr)
}

func (n *NLALogService) NetlinkNeigh(nlmsg *nlamsg.NetlinkMessage, neigh *nlamsg.Neigh) {
	n.log.Debugf("NEIG: %s %v %v", nlmsg.Src, &nlmsg.Header, neigh)
}

func (n *NLALogService) NetlinkRoute(nlmsg *nlamsg.NetlinkMessage, route *nlamsg.Route) {
	n.log.Debugf("ROUT: %s %v %v", nlmsg.Src, &nlmsg.Header, route)
}

func (n *NLALogService) NetlinkNode(nlmsg *nlamsg.NetlinkMessage, node *nlamsg.Node) {
	n.log.Debugf("NODE: %s %v %v", nlmsg.Src, &nlmsg.Header, node.IP())
}

func (n *NLALogService) NetlinkVpn(nlmsg *nlamsg.NetlinkMessage, vpn *nlamsg.Vpn) {
	n.log.Debugf("VPN : %s %v %v", nlmsg.Src, &nlmsg.Header, vpn)
}

func (n *NLALogService) NetlinkBridgeVlanInfo(nlmsg *nlamsg.NetlinkMessage, brvlan *nlamsg.BridgeVlanInfo) {
	n.log.Debugf("BRVI: %s %v %v", nlmsg.Src, &nlmsg.Header, brvlan)
}
