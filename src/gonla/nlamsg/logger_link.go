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
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func LogNetlinkLinkAttrs(logger LogLogger, level log.Level, m *netlink.LinkAttrs) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "LinkAttrs: index     : %d", m.Index)
	logger.Logf(level, "LinkAttrs: mtu       : %d", m.MTU)
	logger.Logf(level, "LinkAttrs: name      : '%s'", m.Name)
	logger.Logf(level, "LinkAttrs: mac       : '%s'", m.HardwareAddr)
	logger.Logf(level, "LinkAttrs: flags     : '%s'", m.Flags)
	logger.Logf(level, "LinkAttrs: parent    : %d", m.ParentIndex)
	logger.Logf(level, "LinkAttrs: master    : %d", m.MasterIndex)
	logger.Logf(level, "LinkAttrs: encap     : '%s'", m.EncapType)
	logger.Logf(level, "LinkAttrs: operstate : %s", m.OperState)
}

func LogNetlinkLink(logger LogLogger, level log.Level, m netlink.Link) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Link: type :%s", m.Type())
	LogNetlinkLinkAttrs(logger, level, m.Attrs())

	switch link := m.(type) {
	case *netlink.Device:
		logger.Logf(level, "Link(Device):")

	case *netlink.Dummy:
		logger.Logf(level, "Link(Dummy):")

	case *netlink.Ifb:
		logger.Logf(level, "Link(Ifb):")

	case *netlink.Bridge:
		if v := link.MulticastSnooping; v != nil {
			logger.Logf(level, "Link(Bridge): M.C.Snoop  : %t", *v)
		}
		if v := link.HelloTime; v != nil {
			logger.Logf(level, "Link(Bridge): hello time : %d", *v)
		}
		if v := link.VlanFiltering; v != nil {
			logger.Logf(level, "Link(Bridge): vlan filter: %t", *v)
		}
	case *netlink.Vlan:
		logger.Logf(level, "Link(Vlan): vid  : %d", link.VlanId)
		logger.Logf(level, "Link(Vlan): proto: %d", link.VlanProtocol)

	case *netlink.Macvlan:
		logger.Logf(level, "Link(MacVlan): mode : %s", link.Mode)
		if addrs := link.MACAddrs; addrs != nil {
			for _, addr := range addrs {
				logger.Logf(level, "Link(MacVlan): addr : %s", addr)
			}
		}

	case *netlink.Macvtap:
		logger.Logf(level, "Link(MacVtap): mode : %s", link.Mode)
		if addrs := link.MACAddrs; addrs != nil {
			for _, addr := range addrs {
				logger.Logf(level, "Link(MacVtap): addr : %s", addr)
			}
		}

	case *netlink.Tuntap:
		logger.Logf(level, "Link(Tuntap): mode    : %s", link.Mode)
		logger.Logf(level, "Link(Tuntap): flags   : %s", link.Flags)
		logger.Logf(level, "Link(Tuntap): non-per : %s", link.NonPersist)

	case *netlink.Veth:
		logger.Logf(level, "Link(Veth): peer-name: %s", link.PeerName)
		logger.Logf(level, "Link(Veth): peer-mac : %s", link.PeerHardwareAddr)

	case *netlink.GenericLink:
		logger.Logf(level, "Link(Generic): link-type: %s", link.LinkType)

	case *netlink.Vxlan:
		logger.Logf(level, "Link(Generic): vxlanid: %d", link.VxlanId)

	case *netlink.IPVlan:
		logger.Logf(level, "Link(Generic): mode: %s", link.Mode)
		logger.Logf(level, "Link(Generic): flag: %s", link.Flag)

	case *netlink.Bond:
		logger.Logf(level, "Link(Bond): mode        : %s", link.Mode)
		logger.Logf(level, "Link(Bond): act-slave   : %d", link.ActiveSlave)
		logger.Logf(level, "Link(Bond): miimon      : %d", link.Miimon)
		logger.Logf(level, "Link(Bond): arp-interval: %d", link.ArpInterval)
		if targets := link.ArpIpTargets; targets != nil {
			for _, target := range targets {
				logger.Logf(level, "Link(Bond): arp-targets : %d", target)
			}
		}
		logger.Logf(level, "Link(Bond): primary     : %d", link.Primary)

	case *netlink.Gretap:
		logger.Logf(level, "Link(Gretap): ikey       : %d", link.IKey)
		logger.Logf(level, "Link(Gretap): okey       : %d", link.OKey)
		logger.Logf(level, "Link(Gretap): encap-sport: %d->%d", link.EncapSport)
		logger.Logf(level, "Link(Gretap): encap-dport: %d->%d", link.EncapDport)
		logger.Logf(level, "Link(Gretap): addr       : %s->%s", link.Local, link.Remote)
		logger.Logf(level, "Link(Gretap): link       : %d", link.Link)

	case *netlink.Iptun:
		logger.Logf(level, "Link(Iptun): ttl        : %d", link.Ttl)
		logger.Logf(level, "Link(Iptun): tos        : %d", link.Tos)
		logger.Logf(level, "Link(Iptun): pmtu-disc  : %d", link.PMtuDisc)
		logger.Logf(level, "Link(Iptun): local      : %s", link.Local)
		logger.Logf(level, "Link(Iptun): remote     : %s", link.Remote)
		logger.Logf(level, "Link(Iptun): encap-sport: %d", link.EncapSport)
		logger.Logf(level, "Link(Iptun): encap-dport: %d", link.EncapDport)
		logger.Logf(level, "Link(Iptun): encap-type : %d", link.EncapType)
		logger.Logf(level, "Link(Iptun): encap-flags: %d", link.EncapFlags)
		logger.Logf(level, "Link(Iptun): flow-based : %t", link.FlowBased)

	case *netlink.Sittun:
		logger.Logf(level, "Link(Sittun): link       : %d", link.Link)
		logger.Logf(level, "Link(Sittun): ttl        : %d", link.Ttl)
		logger.Logf(level, "Link(Sittun): tos        : %d", link.Tos)
		logger.Logf(level, "Link(Sittun): pmtu-disc  : %d", link.PMtuDisc)
		logger.Logf(level, "Link(Sittun): local      : %s", link.Local)
		logger.Logf(level, "Link(Sittun): remote     : %s", link.Remote)
		logger.Logf(level, "Link(Sittun): encap-sport: %d", link.EncapSport)
		logger.Logf(level, "Link(Sittun): encap-dport: %d", link.EncapDport)
		logger.Logf(level, "Link(Sittun): encap-type : %d", link.EncapType)
		logger.Logf(level, "Link(Sittun): encap-flags: %d", link.EncapFlags)

	case *netlink.Vti:
		logger.Logf(level, "Link(Vti): link  : %d", link.Link)

	case *netlink.Gretun:
		logger.Logf(level, "Link(Gretun): link  : %d", link.Link)

	case *netlink.Vrf:
		logger.Logf(level, "Link(Vrf): table: %d", link.Table)

	case *netlink.GTP:
		logger.Logf(level, "Link(Vrf): fd0: %d", link.FD0)

	case *netlink.Xfrmi:
		logger.Logf(level, "Link(Vrf): ifid: %d", link.Ifid)

	case *netlink.IPoIB:
		logger.Logf(level, "Link(IPoIB): pkey: %d", link.Pkey)

	default:
		logger.Logf(level, "Link(%s):", link.Type())
	}
}

func LogLink(logger LogLogger, level log.Level, m *Link) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Link: nid :%d", m.NId)
	logger.Logf(level, "Link: lnid:%d", m.LnId)

	LogNetlinkLink(logger, level, m.Link)
}
