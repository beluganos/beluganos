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

package ribctl

import (
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gonla/nlamsg"
	"net"
)

type IfaceMap struct {
	ifmap map[string]string
}

func NewIfaceMap() *IfaceMap {
	return &IfaceMap{
		ifmap: make(map[string]string),
	}
}

func NewIfaceEntry(nid uint8, ifindex int) string {
	return fmt.Sprintf("%d_%d", nid, ifindex)
}

func (i *IfaceMap) Set(ifname string, nid uint8, ifindex int) {
	i.ifmap[ifname] = NewIfaceEntry(nid, ifindex)
}

func (i *IfaceMap) Delete(ifname string) {
	delete(i.ifmap, ifname)
}

func (i *IfaceMap) Contains(nid uint8, ifindex int) bool {
	if ifindex > 0 {
		entry := NewIfaceEntry(nid, ifindex)
		for _, v := range i.ifmap {
			if v == entry {
				return true
			}
		}
	}
	return ifindex == 0
}

type RIBController struct {
	nid   uint8
	reId  string
	label uint32
	ifmap *IfaceMap
	nla   *NLAController
	fib   *FIBController
}

func NewRIBController(nid uint8, reId string, label uint32, nla *NLAController, fib *FIBController) *RIBController {
	return &RIBController{
		nid:   nid,
		reId:  reId,
		label: label,
		nla:   nla,
		fib:   fib,
		ifmap: NewIfaceMap(),
	}
}

func (r *RIBController) Serve(done <-chan struct{}) {
	for {
		select {
		case conn := <-r.nla.Conn():
			if conn != nil {
				r.fib.Start()
			} else {
				r.fib.Stop()
			}
		case connected := <-r.fib.Conn():
			if connected {
				r.SendHello()
				r.SendPortConfigs()
			}
		case msg := <-r.nla.Recv():
			r.OnNetlinkMessage(msg)
		case msg := <-r.fib.Recv():
			r.OnFIBCMessage(msg)
		case <-done:
			return
		}
	}
}

func (r *RIBController) Start(done <-chan struct{}) {
	go r.Serve(done)
}

func (r *RIBController) OnFIBCMessage(msg fibcnet.Message) {
	FIBCDispatch(msg, r)
}

func (r *RIBController) OnDpStatus(msg *fibcapi.DpStatus) {
	log.Debugf("RIBController: OnDpStatus %v", msg)
	if msg.Status == fibcapi.DpStatus_ENTER {
		r.SendLoopbackFlows(fibcapi.FlowMod_ADD, r.nid)
	}
}

func (r *RIBController) OnPortStatus(msg *fibcapi.PortStatus) {
	log.Debugf("RIBController: OnPortStatus %v", msg)

	nid, ifname := ParseLinkName(msg.Ifname)

	if msg.Status == fibcapi.PortStatus_UP {
		if nid != r.nid {
			if err := r.SendMPLSFlowVRF(fibcapi.FlowMod_ADD, nid); err != nil {
				log.Errorf("RIBController: SendMPLSFlowVRF error. %s", err)
			}
		}

		r.nla.GetLinks(func(link *nlamsg.Link) error {
			if NewPortId(link) == msg.PortId {
				if err := r.SendLinkFlows(fibcapi.FlowMod_ADD, link); err != nil {
					log.Errorf("RIBController: SendLinkFlows error. %s", err)
				}

				r.ifmap.Set(NewLinkLinkName(link), link.NId, link.Attrs().Index)
			}
			return nil
		})

		r.nla.GetAddrs(func(addr *nlamsg.Addr) error {
			if NewAddrLinkName(addr) == msg.Ifname {
				if err := r.SendACLFlowByAddr(fibcapi.FlowMod_ADD, addr); err != nil {
					log.Errorf("RIBController: ACL FLow(Addr) error. %s", err)
					return err
				}
			}
			return nil
		})

		r.nla.GetNeighs(func(neigh *nlamsg.Neigh) error {
			link, err := r.nla.GetLink(neigh.NId, neigh.LinkIndex)
			if err != nil {
				log.Warnf("RIBController: Link not found. Neigh %s %s", neigh, err)
				return nil
			}

			if NewPortId(link) == msg.PortId {
				if err := r.SendNeighFlows(fibcapi.FlowMod_ADD, neigh); err != nil {
					log.Errorf("RIBController: SendNeighs error. %s", err)
					return err
				}
			}
			return nil
		})

		if err := r.nla.ModLinkStatus(nid, ifname, "OperUp"); err != nil {
			log.Errorf("RIBController: ModLinkStatus error. %d/%s UP", nid, ifname)
		}
	} else {
		if err := r.nla.ModLinkStatus(nid, ifname, "OperDown"); err != nil {
			log.Errorf("RIBController: ModLinkStatus error. %d/%s Down", nid, ifname)
		}
	}
}

func (r *RIBController) OnNetlinkMessage(nlmsg *nlamsg.NetlinkMessageUnion) {
	nlamsg.DispatchUnion(nlmsg, r)
}

func (r *RIBController) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {
	log.Debugf("RIBController: LINK nid:%d LnId:%d", link.NId, link.LnId)

	cmd := GetPortConfigCmd(nlmsg.Type())
	if err := r.SendPortConfig(cmd, link); err != nil {
		log.Errorf("RIBController: SendPortConfig(%s) error. %v %s", cmd, link, err)
	}

	log.Debugf("RIBController: LINK %s %v", cmd, link)
}

func (r *RIBController) NetlinkAddr(nlmsg *nlamsg.NetlinkMessage, addr *nlamsg.Addr) {
	log.Debugf("RIBController: ADDR NId:%d AdId:%d", addr.NId, addr.AdId)

	if ok := r.ifmap.Contains(addr.NId, int(addr.Index)); !ok {
		log.Warnf("RIBController: Ifindex not found. Addr %s", addr)
		return
	}

	cmd := GetFlowCmd(nlmsg.Type())
	if err := r.SendACLFlowByAddr(cmd, addr); err != nil {
		log.Errorf("RIBController: ADDR %s error. %v %s", cmd, addr, err)
	}

	log.Debugf("RIBController: ADDR %s %v", cmd, addr)
}

func (r *RIBController) NetlinkNeigh(nlmsg *nlamsg.NetlinkMessage, neigh *nlamsg.Neigh) {
	log.Debugf("RIBController: NEIGH NId;%d NeId:%d", neigh.NId, neigh.NeId)

	if ok := r.ifmap.Contains(neigh.NId, neigh.LinkIndex); !ok {
		log.Warnf("RIBController: Ifindex not found. Neigh %s", neigh)
		return
	}

	cmd := GetFlowCmd(nlmsg.Type())
	if err := r.SendNeighFlows(cmd, neigh); err != nil {
		log.Errorf("RIBController: NEIGH %s error. %v %s", cmd, neigh, err)
	}

	log.Debugf("RIBController: NEIGH %s %v", cmd, neigh)
}

func (r *RIBController) NetlinkRoute(nlmsg *nlamsg.NetlinkMessage, route *nlamsg.Route) {
	log.Debugf("RIBController: ROUTE NId:%d RtId:%d", route.NId, route.RtId)

	if ok := r.ifmap.Contains(route.NId, route.GetLinkIndex()); !ok {
		log.Warnf("RIBController: Ifindex not found. Route %s", route)
		return
	}

	cmd := GetFlowCmd(nlmsg.Type())
	if err := r.SendRouteFlows(cmd, route); err != nil {
		log.Errorf("RIBController: ROUTE %s error. %v %s", cmd, route, err)
	}

	log.Debugf("RIBController: ROUTE %s %v", cmd, route)
}

func (r *RIBController) SendLinkFlows(cmd fibcapi.FlowMod_Cmd, link *nlamsg.Link) error {
	grpCmd := FlowCmdToGroupCmd(cmd)

	if err := r.SendL2InterfaceGroup(grpCmd, link); err != nil {
		log.Errorf("RIBController: L2 Interface Group error. %s", err)
		return err
	}

	if err := r.SendVLANFlow(cmd, link); err != nil {
		log.Errorf("RIBController: VLAN Flow error. %s", err)
		return err
	}

	if err := r.SendTermMACFlow(cmd, link); err != nil {
		log.Errorf("RIBController: TermMAC Flow error. %s", err)
		return err
	}

	return nil
}

func (r *RIBController) SendNeighFlows(cmd fibcapi.FlowMod_Cmd, neigh *nlamsg.Neigh) error {
	grpCmd := FlowCmdToGroupCmd(cmd)

	if cmd == fibcapi.FlowMod_DELETE {
		if err := r.SendUnicastRoutingFlowNeigh(cmd, neigh); err != nil {
			log.Errorf("RIBController: Unicast Routing(Neigh) error. %s", err)
			return err
		}
	}

	if err := r.SendL3UnicastGroup(grpCmd, neigh); err != nil {
		log.Errorf("RIBController: L3 Unicast Group error. %s", err)
		return err
	}

	if err := r.SendMPLSInterfaceGroup(grpCmd, neigh); err != nil {
		log.Errorf("RIBController: MPLS Interface error. %s", err)
		return err
	}

	if cmd != fibcapi.FlowMod_DELETE {
		if err := r.SendUnicastRoutingFlowNeigh(cmd, neigh); err != nil {
			log.Errorf("RIBController: Unicast Routing(Neigh) error. %s", err)
			return err
		}
	}

	return nil
}

func (r *RIBController) SendRouteFlows(cmd fibcapi.FlowMod_Cmd, route *nlamsg.Route) error {

	if route.GetDst() != nil {
		if route.GetMPLSEncap() == nil {
			// IP Routing
			if err := r.SendUnicastRoutingFlow(cmd, route); err != nil {
				log.Errorf("RIBController: Unicast Routing(IP) error. %s", err)
				return err
			}
		} else {
			// PUSH (single label)
			//  Unicast Routing flow
			//   -> MPLS L3 VPN (0x92LLLLLL) L:LDP Label
			//    -> MPLS Interface (0x90VVNNNN) V:VRF/N:NeId
			//
			// PUSH (double label)
			//  Unicast Routing flow
			//   -> MPLS L3 VPN (0x92LLLLLL) L:VPN Label
			//     -> MPLS Tun Label1 (0x93LLLLLL) L:LDP Label
			//      -> MPLS Interface(0x90VVNNNN) V:VRF/N:NeId
			grpCmd := FlowCmdToGroupCmd(cmd)

			if cmd == fibcapi.FlowMod_DELETE {
				if err := r.SendUnicastRoutingFlowMPLS(cmd, route); err != nil {
					log.Errorf("RIBController: Unicast Routing(MPLS) error. %s", err)
					return err
				}
			}

			if err := r.SendMPLSLabelGroupMPLS(grpCmd, route); err != nil {
				log.Errorf("RIBController: SendMPLSLabelGroupMPLS error. %s", err)
				return err
			}

			if err := r.SendMPLSLabelGroupVPN(grpCmd, route); err != nil {
				log.Errorf("RIBController: SendMPLSLabelGroupVPN error. %s", err)
				return err
			}

			if cmd != fibcapi.FlowMod_DELETE {
				if err := r.SendUnicastRoutingFlowMPLS(cmd, route); err != nil {
					log.Errorf("RIBController: Unicast Routing(MPLS) error. %s", err)
					return err
				}
			}
		}
	}

	if route.MPLSDst != nil {
		if route.GetMPLSNewDst() == nil {
			// POP (single label -> no label)
			//  MPLS1(BOS=1)  DEC_MPLS_TTL,SET(VRF)
			//   -> MPLS L3 Type, POP_MPLS, MPLS_TYPE=L3 Unicast(built-in)
			//    -> MPLS Label Trust(built-in)
			//     -> MPLS Type(built-in/L3 Unicast)
			//      -> Unicast Routing
			if err := r.SendMPLSFlowPop1(cmd, route); err != nil {
				log.Errorf("SendMPLSFlowPop1 error. %s", err)
				return err
			}

			// POP (double labels -> single label)
			//  MPLS1(BOS=0)  DEC_MPLS_TTL,POP_MPLS,GROUP(MPLS Interface)
			//   -> MPLS Label Trust(built-in/skip)
			//    -> MPLS Type(built-in/miss)
			//     -> ACL
			//      -> MPLS Interface(0x90VVNNNN) V:VRF/N:NeId
			//
			if err := r.SendMPLSFlowPop2(cmd, route); err != nil {
				log.Errorf("SendMPLSFlowPop2 error. %s", err)
				return err
			}

		} else {
			// SWAP
			//  MPLS1(BOS=0/1)  DEC_MPLS_TTL, GROUP(MPLS Swap Label)
			//   -> MPLS Label Trust(built-in/skip)
			//    -> MPLS Type(built-in/miss)
			//     -> ACL
			//      -> MPLS Swap Label(0x95LLLLLL) L:LLabel
			//       -> MPLS Interface
			grpCmd := FlowCmdToGroupCmd(cmd)

			if cmd != fibcapi.FlowMod_DELETE {
				if err := r.SendMPLSLabelGroupSwap(grpCmd, route); err != nil {
					log.Errorf("SendMPLSLabelGroupSwap error. %s", err)
					return err
				}
			}

			if err := r.SendMPLSFlowSwap(cmd, route, true); err != nil {
				log.Errorf("SendMPLSFlowSwap(BOS=1) error. %s", err)
				return err
			}

			if err := r.SendMPLSFlowSwap(cmd, route, false); err != nil {
				log.Errorf("SendMPLSFlowSwap(BOS=0) error. %s", err)
				return err
			}

			if cmd == fibcapi.FlowMod_DELETE {
				if err := r.SendMPLSLabelGroupSwap(grpCmd, route); err != nil {
					log.Errorf("SendMPLSLabelGroupSwap error. %s", err)
					return err
				}
			}
		}
	}

	return nil
}

func (r *RIBController) SendLoopbackFlows(cmd fibcapi.FlowMod_Cmd, nid uint8) {
	links := make(map[int32]struct{}, 0)
	r.nla.GetLinks(func(link *nlamsg.Link) error {
		if nid == link.NId {
			if (link.Attrs().Flags & net.FlagLoopback) != 0 {
				links[int32(link.Attrs().Index)] = struct{}{}
			}
		}
		return nil
	})

	r.nla.GetAddrs(func(addr *nlamsg.Addr) error {
		if addr.NId == nid {
			if _, ok := links[addr.Index]; ok {
				if err := r.SendACLFlowByAddr(cmd, addr); err != nil {
					log.Errorf("RIBController: ACL FLow(Addr) error. %s", err)
					return err
				}
			}
		}
		return nil
	})
}
