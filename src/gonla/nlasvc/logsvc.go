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
	dump  uint32
	log   *log.Entry
	level log.Level
}

func NewNLALogService(dump uint32) *NLALogService {
	return &NLALogService{
		dump:  dump,
		log:   NewLogger("NLALogService"),
		level: log.DebugLevel,
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
	n.log.Debugf("NLM : %v", nlmsg)
	if n.dump != 0 {
		for _, line := range nlalib.HexSlice(nlmsg.Data) {
			n.log.Debugf("NLM : %s", line)
		}
	}
}

func (n *NLALogService) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {
	nlamsg.LogLink(n.log, n.level, link)
}

func (n *NLALogService) NetlinkAddr(nlmsg *nlamsg.NetlinkMessage, addr *nlamsg.Addr) {
	nlamsg.LogAddr(n.log, n.level, addr)
}

func (n *NLALogService) NetlinkNeigh(nlmsg *nlamsg.NetlinkMessage, neigh *nlamsg.Neigh) {
	nlamsg.LogNeigh(n.log, n.level, neigh)
}

func (n *NLALogService) NetlinkRoute(nlmsg *nlamsg.NetlinkMessage, route *nlamsg.Route) {
	nlamsg.LogRoute(n.log, n.level, route)
}

func (n *NLALogService) NetlinkNode(nlmsg *nlamsg.NetlinkMessage, node *nlamsg.Node) {
	nlamsg.LogNode(n.log, n.level, node)
}

func (n *NLALogService) NetlinkVpn(nlmsg *nlamsg.NetlinkMessage, vpn *nlamsg.Vpn) {
	nlamsg.LogVpn(n.log, n.level, vpn)
}

func (n *NLALogService) NetlinkBridgeVlanInfo(nlmsg *nlamsg.NetlinkMessage, brvlan *nlamsg.BridgeVlanInfo) {
	nlamsg.LogBridgeVlanInfo(n.log, n.level, brvlan)
}
