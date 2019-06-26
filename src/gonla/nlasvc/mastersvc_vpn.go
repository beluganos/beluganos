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

package nlasvc

import (
	"gonla/nladbm"
	"gonla/nlalib"
	"gonla/nlamsg"

	"github.com/vishvananda/netlink"
)

func (n *NLAMasterService) NewVpnRoute(ricRoute *nlamsg.Route) *nlamsg.Route {

	vpn := nladbm.Vpns().Select(nladbm.NewVpnKey(ricRoute.NId, ricRoute.GetDst(), ricRoute.GetGw()))
	if vpn == nil {
		n.log.Debugf("NewVpnRoute VPN not found. nid:%d dst:%s gw:%s", ricRoute.NId, ricRoute.GetDst(), ricRoute.GetGw())
		return nil
	}

	gwDst := nlalib.NewIPNetFromIP(vpn.NetVpnGw())
	gwRoute := nladbm.Routes().Select(nladbm.NewRouteKey(n.NId, gwDst))
	if gwRoute == nil {
		n.log.Debugf("NewVpnRoute Route for VPN not found. %d, %s", n.NId, gwDst)
		return nil
	}

	// Create MplsInfo := [Mpls(GW)..., Mpls(VPN)]
	labels := []int{}
	enIds := []uint32{}
	if encap := gwRoute.GetMPLSEncap(); encap != nil {
		labels = encap.Labels
		enIds = gwRoute.EnIds
	}
	labels = append(labels, int(vpn.Label))
	enIds = append(enIds, nladbm.Encaps().EncapId(gwDst, vpn.Label))

	// Create new VPN Route instance.
	// *** Do not change original Route ***
	vpnRoute := gwRoute.Copy()
	vpnRoute.Dst = vpn.GetIPNet()
	vpnRoute.NId = vpn.NId
	vpnRoute.VpnGw = vpn.NetVpnGw()
	vpnRoute.LinkIndex = 0
	vpnRoute.Encap = &netlink.MPLSEncap{Labels: labels}
	vpnRoute.EnIds = enIds
	vpnRoute.MultiPath = []*netlink.NexthopInfo{}

	n.log.Debugf("NewVpnRoute %v", vpnRoute)
	return vpnRoute
}

func (n *NLAMasterService) NewVpnRoutes(gwRoute *nlamsg.Route, f func(*nlamsg.Route) error) {

	labels := []int{}
	enIds := []uint32{}
	if gwRoute.GetEncap() != nil {
		if mpls, ok := gwRoute.GetEncap().(*netlink.MPLSEncap); ok {
			labels = mpls.Labels
			enIds = []uint32{nladbm.Encaps().EncapId(gwRoute.GetDst(), 0)}
		}
	}

	err := nladbm.Vpns().WalkByVpnGw(gwRoute.Dst.IP, func(vpn *nlamsg.Vpn) error {
		// Create new VPN Route instance.
		// *** Do not change original Route ***
		vpnLabels := append(labels, int(vpn.Label))
		vpnEnIds := append(enIds, nladbm.Encaps().EncapId(gwRoute.GetDst(), vpn.Label))

		vpnRoute := gwRoute.Copy()
		vpnRoute.Dst = vpn.GetIPNet()
		vpnRoute.NId = vpn.NId
		vpnRoute.VpnGw = vpn.NetVpnGw()
		vpnRoute.LinkIndex = 0
		vpnRoute.Encap = &netlink.MPLSEncap{Labels: vpnLabels}
		vpnRoute.EnIds = vpnEnIds
		vpnRoute.MultiPath = []*netlink.NexthopInfo{}

		n.log.Debugf("NewVpnRoutes %v", vpnRoute)
		return f(vpnRoute)
	})

	if err != nil {
		n.log.Errorf("NewVpnRoutes Walk error. %s", err)
	}
}
