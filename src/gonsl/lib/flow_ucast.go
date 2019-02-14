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

package gonslib

import (
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"
	"net"

	"github.com/beluganos/go-opennsl/opennsl"
	"golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
)

//
// FIBCUnicastRoutingFlowMod process FlowMod(Unicast Routing)
//
func (s *Server) FIBCUnicastRoutingFlowMod(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.UnicastRoutingFlow) {
	log.Debugf("Server: UC-RoutingFlow: %v %v %v", hdr, mod, flow)

	ip, ipnet, err := net.ParseCIDR(flow.Match.IpDst)
	if err != nil {
		log.Errorf("Server: UC-RoutingFlow: Invalid IP address. %s %s", flow.Match.IpDst, err)
		return
	}

	vrf := opennsl.Vrf(flow.Match.Vrf)
	neid := flow.GId

	switch flow.GType {
	case fibcapi.GroupMod_L3_UNICAST:
		switch flow.Match.Origin {
		case fibcapi.UnicastRoutingFlow_NEIGH:
			log.Debugf("Server: UC-RoutingFlow(Neigh): %s neid:%08x", ip, neid)

			l3host := opennsl.NewL3Host()

			if IPToAF(ip) == unix.AF_INET {
				l3host.SetIPAddr(ip)
			} else {
				l3host.SetIP6Addr(ip)
				l3host.SetFlags(opennsl.L3_IP6)
			}

			if vrf != 0 {
				l3host.SetVRF(vrf)
			}

			switch mod.Cmd {
			case fibcapi.FlowMod_ADD:
				l3egrID, ok := s.idmaps.L3Egress.Get(neid)
				if !ok {
					log.Errorf("Server: UC-RoutingFlow(Neigh): L3Egress(neid:%08x) not found.", neid)
					return
				}

				log.Debugf("Server: UC-RoutingFlow(Neigh): %s l3eg:%d", ip, l3egrID)

				l3host.SetEgressID(l3egrID)

				if err := l3host.Add(s.Unit()); err != nil {
					log.Errorf("Server: UC-RoutingFlow: L3Host add error. %s", err)
				}

			case fibcapi.FlowMod_MODIFY, fibcapi.FlowMod_MODIFY_STRICT:
				log.Errorf("Server: UC-RoutingFlow: L3Host modify unsupported.")

			case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
				if err := l3host.Delete(s.Unit()); err != nil {
					log.Errorf("Server: UC-RoutingFlow: L3Host delete error. %s", err)
				}

			default:
				log.Errorf("Server: UC-RoutingFlow(Neigh): Invalid Command. %d", mod.Cmd)
			}

		case fibcapi.UnicastRoutingFlow_ROUTE:
			log.Debugf("Server: UC-RoutingFlow(Route): %s neid:%08x", ipnet, neid)

			l3route := opennsl.NewL3Route()

			if IPToAF(ip) == unix.AF_INET {
				l3route.SetIP4Net(ipnet)
			} else {
				l3route.SetIP6Net(ipnet)
				l3route.SetFlags(opennsl.L3_IP6)
			}

			if vrf != 0 {
				l3route.SetVRF(vrf)
			}

			switch mod.Cmd {
			case fibcapi.FlowMod_ADD:
				l3egrID, ok := s.idmaps.L3Egress.Get(neid)
				if !ok {
					log.Errorf("Server: UC-RoutingFlow(Route): L3Egress(neid:%08x) not found.", neid)
					return
				}

				log.Debugf("Server: UC-RoutingFlow(Route): %s l3eg:%d", ipnet, l3egrID)

				l3route.SetEgressID(l3egrID)

				if err := l3route.Add(s.Unit()); err != nil {
					log.Errorf("Server: UC-RoutingFlow: L3Route add error. %s", err)
				}

			case fibcapi.FlowMod_MODIFY, fibcapi.FlowMod_MODIFY_STRICT:
				log.Errorf("Server: UC-RoutingFlow: L3Route modify unsupported.")

			case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
				if err := l3route.Delete(s.Unit()); err != nil {
					log.Errorf("Server: UC-RoutingFlow: L3Route delete error. %s", err)
				}

			default:
				log.Errorf("Server: UC-RoutingFlow(Route): Invalid Command. %d", mod.Cmd)
			}

		default:
			log.Errorf("Server: UC-RoutingFlow: Invalid Origin %d", flow.Match.Origin)
		}

	case fibcapi.GroupMod_L3_ECMP:
		log.Warnf("Server: UC-RoutingFlow(L3_ECMP): %s", ipnet)

	case fibcapi.GroupMod_MPLS_L3_VPN:
		log.Warnf("Server: UC-RoutingFlow(MPLS_L3_VPN): %s", ipnet)

	default:
		log.Errorf("Server: UC-RoutingFlow: Invalid Group type. %d", flow.GType)
	}
}
