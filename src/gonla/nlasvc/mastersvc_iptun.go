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
	"gonla/nladbm"
	"gonla/nlamsg"
	"net"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type NLAMasterIptun struct {
	mutex          sync.Mutex
	master         *NLAMasterService
	updateInterval time.Duration
	log            *log.Entry
}

func NewNLAMasterIptun(master *NLAMasterService, updateInterval time.Duration) *NLAMasterIptun {
	return &NLAMasterIptun{
		master:         master,
		updateInterval: updateInterval,
		log:            NewLogger("NLAMasterIptun"),
	}
}

func (n *NLAMasterIptun) NewNeigh(tun *nlamsg.Iptun, route *nlamsg.Route) {
	nid := tun.NId

	phyln := nladbm.Links().Select(nladbm.NewLinkKey(nid, route.GetLinkIndex()))
	if phyln == nil {
		n.log.Warnf("NewNeigh phy-link %d/%d not found.", nid, route.GetLinkIndex())
		return
	}

	neigh := nladbm.Neighs().Select(nladbm.NewNeighKey(nid, route.GetGw()))
	if neigh == nil {
		n.log.Warnf("NewNeigh neigh %d/%s not found.", nid, route.GetGw())
		return
	}

	tun.LocalMAC = phyln.Attrs().HardwareAddr
	n.log.Infof("NewNeigh remote %d/%s", nid, tun.Remote())

	neigh = neigh.Copy()
	neigh.NeId = 0
	neigh.IP = tun.Remote()
	neigh.LinkIndex = tun.Attrs().Index
	neigh.PhyLink = phyln.Attrs().Index
	neigh.Tunnel = nlamsg.NewNeighIptun(
		tun.Type(),
		tun.Local(),
	)

	nlmsg := &nlamsg.NetlinkMessage{}
	nlmsg.NId = nid
	nlmsg.Src = nlamsg.SRC_KNL
	nlmsg.Header.Type = syscall.RTM_NEWNEIGH

	n.log.Debugf("NewNeigh %s", neigh)

	n.master.NetlinkNeigh(nlmsg, neigh)
}

func (n *NLAMasterIptun) DelNeigh(tun *nlamsg.Iptun) {
	nid := tun.NId

	n.log.Infof("DelNeigh neigh %d/%s", nid, tun.Remote())

	neigh := nlamsg.NewNeigh(nil, nid, 0)
	neigh.NeId = 0
	neigh.IP = tun.Remote()
	neigh.LinkIndex = tun.Attrs().Index
	neigh.PhyLink = 0
	neigh.Tunnel = nlamsg.NewNeighIptun(
		tun.Type(),
		tun.Local(),
	)

	nlmsg := &nlamsg.NetlinkMessage{}
	nlmsg.NId = nid
	nlmsg.Src = nlamsg.SRC_KNL
	nlmsg.Header.Type = syscall.RTM_DELNEIGH

	tun.LocalMAC = nil

	n.log.Debugf("DelNeigh %s", neigh)

	n.master.NetlinkNeigh(nlmsg, neigh)
}

func (n *NLAMasterIptun) RemoteUp(nid uint8, remote net.IP) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	route := nladbm.Routes().SelectByTunRemote(nid, remote)

	if route == nil {
		n.log.Warnf("RemoteUp route %d/%s not found.", nid, remote)
		return
	}

	tunKey := nladbm.NewIptunKey(nid, remote)
	if tun := nladbm.Links().SelectTun(tunKey); tun != nil {
		n.NewNeigh(tun, route)
	}
}

func (n *NLAMasterIptun) RemoteRouteUp(route *nlamsg.Route) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	nladbm.Links().WalkTunByRemote(route.NId, route.Dst, func(iptun *nlamsg.Iptun) error {
		n.NewNeigh(iptun, route)
		return nil
	})
}

func (n *NLAMasterIptun) RemoteDown(nid uint8, remote net.IP) {
	n.remoteDown(nid, remote, false)
}

func (n *NLAMasterIptun) remoteDown(nid uint8, remote net.IP, checkRoute bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if checkRoute {
		route := nladbm.Routes().SelectByTunRemote(nid, remote)
		if route != nil {
			n.log.Debugf("remoteDown neigh %d/%s not down.", nid, remote)
			return
		}
	}

	tunKey := nladbm.NewIptunKey(nid, remote)
	if tun := nladbm.Links().SelectTun(tunKey); tun != nil {
		n.DelNeigh(tun)
	}
}

func (n *NLAMasterIptun) RemoteRouteDown(route *nlamsg.Route) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	nladbm.Links().WalkTunByRemote(route.NId, route.Dst, func(iptun *nlamsg.Iptun) error {
		n.DelNeigh(iptun)
		return nil
	})
}

func (n *NLAMasterIptun) NewIptunRoute(route *nlamsg.Route) *nlamsg.Route {
	nid := route.NId
	ifindex := route.GetLinkIndex()

	link := nladbm.Links().Select(nladbm.NewLinkKey(nid, ifindex))
	if link == nil {
		n.log.Errorf("NewIptunRoute link %d/%d not found.", nid, ifindex)
		return nil
	}

	iptun := link.Iptun()
	if iptun == nil {
		n.log.Debugf("NewIptunRoute link %d/%d is not iptun.", nid, ifindex)
		return nil
	}

	neigh := nladbm.Neighs().Select(nladbm.NewNeighKey(nid, iptun.Remote))
	if neigh == nil {
		n.log.Debugf("NewIptunRoute neigh %d/%s nod found.", nid, iptun.Remote)
		return nil
	}

	iptunRoute := route.Copy()
	iptunRoute.Gw = neigh.IP

	return iptunRoute
}

func (n *NLAMasterIptun) NewRoutes(neigh *nlamsg.Neigh, f func(*nlamsg.Route)) {
	nid := neigh.NId
	gw := neigh.IP

	if tunRemote := neigh.IsTunnelRemote(); !tunRemote {
		n.log.Debugf("NewRoutes neigh %d/%s not tunnel-remote.", nid, gw)
		return
	}

	nladbm.Routes().WalkByLinkFree(nid, neigh.LinkIndex, func(route *nlamsg.Route) error {
		iptunRoute := route.Copy()
		iptunRoute.Gw = gw
		f(iptunRoute)
		return nil
	})
}

func (n *NLAMasterIptun) Serve() {
	if n.updateInterval == 0 {
		n.log.Infof("Disabled.")
		return
	}

	ticker := time.NewTicker(n.updateInterval)

	n.log.Infof("START")

	for {
		select {
		case <-ticker.C:
			n.log.Infof("Serve auto update start")

			nladbm.Links().WalkTun(func(tun *nlamsg.Iptun) error {
				if local := tun.LocalMAC; local == nil {
					n.RemoteUp(tun.NId, tun.Remote())
				} else {
					n.remoteDown(tun.NId, tun.Remote(), true)
				}
				return nil
			})

			n.log.Infof("Serve auto update end")
		}
	}
}
