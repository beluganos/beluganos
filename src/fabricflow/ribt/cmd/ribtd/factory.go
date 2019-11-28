// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/osrg/gobgp/pkg/packet/bgp"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const (
	TunnelIfNameDelim = "_"
	TunnelCounterMin  = 4096
	TunnelTypeIp4     = "ipip"
	TunnelTypeIp6     = "ip6tnl"
)

//
// TunnelFactory is factory for TunnelTable.
//
type TunnelFactory struct {
	counter uint32
	device  string
	prefix  string
}

func NewTunnelFactory(device string) *TunnelFactory {
	return &TunnelFactory{
		counter: TunnelCounterMin,
		device:  device,
		prefix:  device + TunnelIfNameDelim,
	}
}

func (t *TunnelFactory) NextCounter() uint32 {
	return atomic.AddUint32(&t.counter, 1)
}

func (t *TunnelFactory) NewIfName(counter uint32) string {
	return fmt.Sprintf("%s%x", t.prefix, counter)
}

func (t *TunnelFactory) ParseIfName(ifname string) (string, uint32, bool) {
	if ok := strings.HasPrefix(ifname, t.prefix); !ok {
		log.Debugf("TunnelFactory.ParseIfName No prefix %s", ifname)
		return ifname, 0, false
	}

	subfix := ifname[len(t.prefix):]
	v, err := strconv.ParseUint(subfix, 16, 32)
	if err != nil {
		log.Debugf("TunnelFactory.ParseIfName Parse failed. %s", ifname)
		return ifname, 0, false
	}

	log.Debugf("TunnelFactory.ParseIfName %s %d", t.device, v)
	return t.device, uint32(v), true
}

func (t *TunnelFactory) LinkList(f func(netlink.Link, uint32)) error {
	links, err := netlink.LinkList()
	if err != nil {
		log.Errorf("TunnelFactory.LinkList error. %s", err)
		return err
	}

	for _, link := range links {
		ifname := link.Attrs().Name
		if _, id, ok := t.ParseIfName(ifname); ok {
			linkType := link.Type()
			if linkType == TunnelTypeIp4 || linkType == TunnelTypeIp6 {
				f(link, id)
			}
		}
	}

	return nil
}

func (t *TunnelFactory) RouteList(link netlink.Link, f func(uint16, *netlink.Route)) error {
	routes, err := netlink.RouteList(link, unix.AF_INET)
	if err != nil {
		return nil
	}
	for _, route := range routes {
		if route.Dst.IP.IsGlobalUnicast() {
			f(bgp.AFI_IP, &route)
		}
	}

	routes, err = netlink.RouteList(link, unix.AF_INET6)
	if err != nil {
		return nil
	}
	for _, route := range routes {
		if route.Dst.IP.IsGlobalUnicast() {
			f(bgp.AFI_IP6, &route)
		}
	}

	return nil
}
