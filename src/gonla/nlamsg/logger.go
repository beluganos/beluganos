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

package nlamsg

import (
	"encoding/hex"
	"gonla/nlamsg/nlalink"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

const (
	LogDataDumpSize = 256
)

type LogLogger interface {
	Logf(log.Level, string, ...interface{})
}

func isSkipLog(level log.Level) bool {
	return log.GetLevel() < level
}

func LogNlMsghdr(logger LogLogger, level log.Level, m *syscall.NlMsghdr) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "NlHdr: len  : %d", m.Len)
	logger.Logf(level, "NlHdr: type : %d", m.Type)
	logger.Logf(level, "NlHdr: flags: 0x%x", m.Flags)
}

func LogNetlinkMessage(logger LogLogger, level log.Level, m *NetlinkMessage, dumpData bool) {
	if isSkipLog(level) {
		return
	}

	LogNlMsghdr(logger, level, &m.Header)

	logger.Logf(level, "NlMsg: nid  : %d", m.NId)
	logger.Logf(level, "NlMsg: src  : %s", m.Src)

	if dumpData {
		dumpSize := len(m.Data)
		if dumpSize > LogDataDumpSize {
			dumpSize = LogDataDumpSize
		}

		logger.Logf(level, "NlMsg:\n%s", hex.Dump(m.Data[:dumpSize]))
	}
}

func LogNetlinkAddr(logger LogLogger, level log.Level, m *netlink.Addr) {
	if isSkipLog(level) {
		return
	}

	if ipnet := m.IPNet; ipnet != nil {
		logger.Logf(level, "Addr: ipnet    : %s", m.IPNet)
	}
	logger.Logf(level, "Addr: flags    : 0x%0x", m.Flags)
	logger.Logf(level, "Addr: scope    : %s", m.IPNet)
	if peer := m.Peer; peer != nil {
		logger.Logf(level, "Addr: peer     : %s", peer)
	}
	logger.Logf(level, "Addr: bcast    : %s", m.Broadcast)
	logger.Logf(level, "Addr: pref-lft : %d", m.PreferedLft)
	logger.Logf(level, "Addr: valid-lft: %d", m.ValidLft)
}

func LogAddr(logger LogLogger, level log.Level, m *Addr) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Addr: nid      : %d", m.NId)
	logger.Logf(level, "Addr: adid     : %d", m.AdId)
	logger.Logf(level, "Addr: index    : %d", m.Index)
	logger.Logf(level, "Addr: family   : %d", m.Family)

	LogNetlinkAddr(logger, level, m.Addr)
}

func LogNetlinkNeigh(logger LogLogger, level log.Level, m *netlink.Neigh) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Neigh: link     : %d", m.LinkIndex)
	logger.Logf(level, "Neigh: family   : %d", m.Family)
	logger.Logf(level, "Neigh: state    : %d", m.State)
	logger.Logf(level, "Neigh: type     : %d", m.Type)
	logger.Logf(level, "Neigh: flags    : 0x%x", m.Flags)
	logger.Logf(level, "Neigh: ip       : %s", m.IP)
	logger.Logf(level, "Neigh: mac      : %s", m.HardwareAddr)
	logger.Logf(level, "Neigh: ll ip    : %s", m.LLIPAddr)
	logger.Logf(level, "Neigh: vid      : %d", m.Vlan)
	logger.Logf(level, "Neigh: vni      : %d", m.VNI)
	logger.Logf(level, "Neigh: master   : %d", m.MasterIndex)
}

func LogNeigh(logger LogLogger, level log.Level, m *Neigh) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Neigh: nid      : %d", m.NId)
	logger.Logf(level, "Neigh: neid     : %d", m.NeId)
	logger.Logf(level, "Neigh: phy-link : %d", m.PhyLink)
	if tun := m.Tunnel; tun != nil {
		logger.Logf(level, "Neigh: tunnel   : %s", tun)
	}

	LogNetlinkNeigh(logger, level, m.Neigh)
}

func LogNetlinkRoute(logger LogLogger, level log.Level, m *netlink.Route) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Route: link     : %d", m.LinkIndex)
	logger.Logf(level, "Route: i-link   : %d", m.ILinkIndex)
	logger.Logf(level, "Route: scope    : %d", m.Scope)
	logger.Logf(level, "Route: dst      : %s", m.Dst)
	logger.Logf(level, "Route: src      : %s", m.Src)
	if mpaths := m.MultiPath; mpaths != nil {
		for _, mpath := range mpaths {
			logger.Logf(level, "Route: multipath: %s", mpath)
		}
	}
	logger.Logf(level, "Route: protocol : %d", m.Protocol)
	logger.Logf(level, "Route: priority : %d", m.Priority)
	logger.Logf(level, "Route: table    : %d", m.Table)
	logger.Logf(level, "Route: type     : %d", m.Type)
	logger.Logf(level, "Route: tos      : %d", m.Tos)
	logger.Logf(level, "Route: flags    : 0x%x", m.Flags)
	if mplsDst := m.MPLSDst; mplsDst != nil {
		logger.Logf(level, "Route: mpls-dst : %d", *mplsDst)
	}
	if encap := m.Encap; encap != nil {
		logger.Logf(level, "Route: encap    : %s", m.Encap)
	}
	logger.Logf(level, "Route: mtu      : %d", m.MTU)
	logger.Logf(level, "Route: adv-mss  : %d", m.AdvMSS)
	logger.Logf(level, "Route: hop-limit: %d", m.Hoplimit)
}

func LogRoute(logger LogLogger, level log.Level, m *Route) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Route: nid      : %d", m.NId)
	logger.Logf(level, "Route: rtid     : %d", m.RtId)
	logger.Logf(level, "Route: vpn-gw   : %s", m.VpnGw)
	if enids := m.EnIds; enids != nil {
		for _, enid := range enids {
			logger.Logf(level, "Route: encap-id : %d", enid)
		}
	}

	LogNetlinkRoute(logger, level, m.Route)
}

func LogNetlinkNode(logger LogLogger, level log.Level, m *nlalink.Node) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Node: ip : %s", m.IP())
}

func LogNode(logger LogLogger, level log.Level, m *Node) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Node: nid: %d", m.NId)
}

func LogNetlinkVpn(logger LogLogger, level log.Level, m *nlalink.Vpn) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Vpn: ipnet  : %s", m.GetIPNet())
	logger.Logf(level, "Vpn: gw     : %s", m.NetGw())
	logger.Logf(level, "Vpn: vpn-gw : %s", m.NetVpnGw())
	logger.Logf(level, "Vpn: label  : %d", m.Label)
}

func LogVpn(logger LogLogger, level log.Level, m *Vpn) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Vpn: nid   : %d", m.NId)
	logger.Logf(level, "Vpn: vpnid : %d", m.VpnId)

	LogNetlinkVpn(logger, level, m.Vpn)
}

func LogNetlinkBridgeVlanInfo(logger LogLogger, level log.Level, m *nl.BridgeVlanInfo) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "BrVlanInfo: flags  : 0x%x", m.Flags)
	logger.Logf(level, "BrVlanInfo: vid    : %d", m.Vid)
}

func LogBridgeVlanInfo(logger LogLogger, level log.Level, m *BridgeVlanInfo) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "BrVlanInfo: nid    : %d", m.NId)
	logger.Logf(level, "BrVlanInfo: brid   : %d", m.BrId)
	logger.Logf(level, "BrVlanInfo: name   : '%s'", m.Name)
	logger.Logf(level, "BrVlanInfo: index  : %d", m.Index)
	logger.Logf(level, "BrVlanInfo: master : %d", m.MasterIndex)
	logger.Logf(level, "BrVlanInfo: mtu    : %d", m.Mtu)

	LogNetlinkBridgeVlanInfo(logger, level, &m.BridgeVlanInfo)
}

type logNetlinkMessageUnion struct {
	level  log.Level
	logger LogLogger
}

func (m *logNetlinkMessageUnion) NetlinkLink(nlmsg *NetlinkMessage, link *Link) {
	LogLink(m.logger, m.level, link)
}

func (m *logNetlinkMessageUnion) NetlinkAddr(nlmsg *NetlinkMessage, addr *Addr) {
	LogAddr(m.logger, m.level, addr)
}

func (m *logNetlinkMessageUnion) NetlinkNeigh(nlmsg *NetlinkMessage, neigh *Neigh) {
	LogNeigh(m.logger, m.level, neigh)
}

func (m *logNetlinkMessageUnion) NetlinkRoute(nlmsg *NetlinkMessage, route *Route) {
	LogRoute(m.logger, m.level, route)
}

func (m *logNetlinkMessageUnion) NetlinkNode(nlmsg *NetlinkMessage, node *Node) {
	LogNode(m.logger, m.level, node)
}

func (m *logNetlinkMessageUnion) NetlinkVpn(nlmsg *NetlinkMessage, vpn *Vpn) {
	LogVpn(m.logger, m.level, vpn)
}

func (m *logNetlinkMessageUnion) NetlinkBridgeVlanInfo(nlmsg *NetlinkMessage, brvlan *BridgeVlanInfo) {
	LogBridgeVlanInfo(m.logger, m.level, brvlan)
}

func LogNetlinkMesssage(logger LogLogger, level log.Level, m *NetlinkMessage) {
	if isSkipLog(level) {
		return
	}

	h := logNetlinkMessageUnion{
		level:  level,
		logger: logger,
	}

	LogNlMsghdr(logger, level, &m.Header)

	logger.Logf(level, "NlMsg: nid: %d", m.NId)
	logger.Logf(level, "NlMsg: src: %d", m.Src)

	if err := Dispatch(m, &h); err != nil {
		logger.Logf(level, "NlMsg: %s", m.NetlinkMessage)
	}
}

func LogNetlinkMessageUnion(logger LogLogger, level log.Level, m *NetlinkMessageUnion) {
	if isSkipLog(level) {
		return
	}

	h := logNetlinkMessageUnion{
		level:  level,
		logger: logger,
	}

	LogNlMsghdr(logger, level, &m.Header)

	logger.Logf(level, "NlMsgUni: nid: %d", m.NId)
	logger.Logf(level, "NlMsgUni: src: %d", m.Src)

	if err := DispatchUnion(m, &h); err != nil {
		logger.Logf(level, "NlMsgUni: %s", m.Msg)
	}
}
