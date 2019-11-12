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
	fibcapi "fabricflow/fibc/api"
	fibcnet "fabricflow/fibc/net"
	"fmt"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type RIBController struct {
	nid    uint8
	reId   string
	label  uint32
	ifdb   *IfDB
	nla    *NLAController
	fib    FIBController
	flowdb *FlowConfig
	useNId bool
	log    *log.Entry
}

func NewRIBController(nid uint8, reId string, label uint32, useNId bool, nla *NLAController, fib FIBController, flowdb *FlowConfig) *RIBController {
	return &RIBController{
		nid:    nid,
		reId:   reId,
		label:  label,
		nla:    nla,
		fib:    fib,
		ifdb:   NewIfDB(),
		flowdb: flowdb,
		useNId: useNId,
		log:    log.WithFields(log.Fields{"module": "RIBController"}),
	}
}

func (r *RIBController) Serve(done <-chan struct{}) {
	r.log.Infof("Serve: Start")

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
				r.FIBCConnected()
			}

		case msg := <-r.nla.Recv():
			nlamsg.DispatchUnion(msg, r)

		case msg := <-r.fib.Recv():
			if err := msg.Dispatch(r); err != nil {
				r.log.Errorf("Serve: Dispatch error. %v %s", r, err)
			}

		case <-done:
			r.log.Infof("Serve: Exit")
			return
		}
	}
}

func (r *RIBController) Start(done <-chan struct{}) {
	go r.Serve(done)
}

func (r *RIBController) FIBCConnected() {
	r.log.Debugf("Connected:")

	r.ifdb.Clear()
	r.SendHello()
	nlmsg := nlamsg.NetlinkMessage{}
	nlmsg.Header.Type = unix.RTM_NEWLINK
	r.nla.GetLinks(nlamsg.NODE_ID_ALL, func(link *nlamsg.Link) error {
		r.NetlinkLink(&nlmsg, link)
		return nil
	})
}

func (r *RIBController) FIBCDpStatus(hdr *fibcnet.Header, msg *fibcapi.DpStatus) {
	r.log.Debugf("DpStatus:")
	fibcapi.LogDpStatus(r.log, log.DebugLevel, msg)

	if r.fib.FIBCType() == FIBCTypeGrpc {
		if msg.Status == fibcapi.DpStatus_ENTER {
			r.SendLoopbackFlows(fibcapi.FlowMod_ADD, r.nid, 0)
		}
	}
}

func (r *RIBController) FIBCPortStatus(hdr *fibcnet.Header, msg *fibcapi.PortStatus) {
	r.log.Debugf("PortStatus:")
	fibcapi.LogPortStatus(r.log, log.DebugLevel, msg)

	var ifentry IfDBEntry
	if f := r.ifdb.Update(msg.PortId, func(e *IfDBEntry) IfDBField {
		e.Associated = (msg.Status == fibcapi.PortStatus_UP)
		e.CopyTo(&ifentry)
		return IfDBFieldAny
	}); f.IsNull() {
		r.log.Errorf("PortStatus: ifentry not found. portId:%d", msg.PortId)
		return
	}

	nid := ifentry.NId

	if msg.Status == fibcapi.PortStatus_UP {
		if ifentry.LinkType == fibcapi.LinkType_BRIDGE {
			r.log.Debugf("PortStatus: skip bridge device. %d %s", msg.PortId, ifentry.LinkType)
			return
		}

		if ifentry.LinkType == fibcapi.LinkType_BRIDGE_SLAVE {
			// Setuo L2 access/trunk port.
			r.nla.GetBridgeVlanInfos(nid, ifentry.Index, func(brvlan *nlamsg.BridgeVlanInfo) {
				if err := r.SendBridgeVlanFlows(fibcapi.FlowMod_ADD, brvlan); err != nil {
					r.log.Errorf("PortStatus: Send BRVLAN %s error.", brvlan)
				}
			})

			return
		}

		if nid != r.nid {
			if err := r.SendMPLSFlowVRF(fibcapi.FlowMod_ADD, nid); err != nil {
				r.log.Errorf("PortStatus: add MPLS(VRF) error. %s", err)
			}
		}

		r.nla.GetLinks(nid, func(link *nlamsg.Link) error {
			if NewPortId(link) == msg.PortId {
				if err := r.SendLinkFlows(fibcapi.FlowMod_ADD, link); err != nil {
					r.log.Errorf("PortStatus: add LINK error. %s", err)
				}
			}
			return nil
		})

		r.nla.GetAddrs(nid, func(addr *nlamsg.Addr) error {
			if nid == addr.NId && ifentry.Index == int(addr.Index) {
				if err := r.SendAddrFlows(fibcapi.FlowMod_ADD, addr); err != nil {
					r.log.Errorf("PortStatus: add ACL(Addr) error. %s", err)
					return err
				}
			}
			return nil
		})

		r.nla.GetNeighs(nid, func(neigh *nlamsg.Neigh) error {
			if nid == neigh.NId && ifentry.Index == int(neigh.LinkIndex) {
				if err := r.SendNeighFlows(fibcapi.FlowMod_ADD, neigh); err != nil {
					r.log.Errorf("PortStatus: add NEIGH error. %s", err)
					return err
				}
			}
			return nil
		})
	}
}

func (r *RIBController) FIBCL2AddrStatus(hdr *fibcnet.Header, msg *fibcapi.L2AddrStatus) {
	r.log.Debugf("L2AddrStatus:")
	fibcapi.LogL2AddrStatus(r.log, log.DebugLevel, msg)

	for _, addr := range msg.Addrs {
		var ifentry IfDBEntry
		if ok := r.ifdb.Select(&ifentry, addr.PortId); !ok {
			r.log.Errorf("L2AddrStatus: ifentry not found. port:%d", addr.PortId)
			continue
		}

		if ifentry.LinkType != fibcapi.LinkType_BRIDGE_SLAVE {
			r.log.Debugf("L2AddrStatus: not bridge slave. %s", addr.HwAddr)
			continue
		}

		if err := r.SetFdb(addr, ifentry.Index); err != nil {
			r.log.Errorf("L2AddrStatus: set fdb error. %s %s", addr.HwAddr, err)
			continue
		}
	}
}

func (r *RIBController) NetlinkNode(nlmsg *nlamsg.NetlinkMessage, node *nlamsg.Node) {
	r.log.Debugf("NODE: nid:%d", node.NId)

	if (nlmsg.Type() == nlalink.RTM_DELNODE) && (node.NId != r.nid) {
		if err := r.SendMPLSFlowVRF(fibcapi.FlowMod_DELETE, node.NId); err != nil {
			r.log.Errorf("NODE: del MPLS(VRF) error. %s", err)
		}
	}
}

func (r *RIBController) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {
	r.log.Debugf("LINK: NId:%d LnId:%d", link.NId, link.LnId)

	msgType := nlmsg.Type()

	// RTM_NEWLINK
	if msgType == unix.RTM_NEWLINK {
		r.ifdb.Set(NewIfDBEntryFromLink(link))
		r.log.Debugf("LINK: %d/'%s' registered to ifmap.", link.NId, link.Attrs().Name)

		if err := r.SendPortConfig("ADD", link); err != nil {
			r.log.Errorf("LINK: Add PortConfig error. %v %s", link, err)
		}

		return
	}

	// RTM_DELLINK
	if msgType == unix.RTM_DELLINK {
		if err := r.SendLinkFlows(fibcapi.FlowMod_DELETE, link); err != nil {
			r.log.Errorf("LINK: Link Flows error. %s", err)
		}

		r.ifdb.Delete(NewPortId(link))
		r.log.Debugf("LINK: %d/'%s' unregistered from ifmap.", link.NId, link.Attrs().Name)

		if err := r.SendPortConfig("DELETE", link); err != nil {
			r.log.Errorf("LINK: Del PortConfig error. %v %s", link, err)
		}

		return
	}

	// RTM_SETLINK
	var (
		ifeOld IfDBEntry
		ifeNew IfDBEntry
	)

	fields := r.ifdb.Update(NewPortId(link), func(e *IfDBEntry) IfDBField {
		e.CopyTo(&ifeOld)
		f := e.Update(link)
		e.CopyTo(&ifeNew)
		return f
	})

	r.log.Debugf("LINK: old %d/%d m=%d %s %s",
		ifeOld.NId, ifeOld.Index, ifeOld.MasterIndex, ifeOld.PortStatus, ifeOld.LinkType)
	r.log.Debugf("LINK: new %d/%d m=%d %s %s",
		ifeNew.NId, ifeNew.Index, ifeNew.MasterIndex, ifeNew.PortStatus, ifeNew.LinkType)
	r.log.Debugf("LINK: ifdb updated %s", fields)

	// device status changed.
	if fields.Has(IfDBFieldStatus) {
		// send port condig to set up/down dp port.
		if err := r.SendPortConfig("MODIFY", link); err != nil {
			r.log.Errorf("LINK: Modify PortConfig error. %v %s", link, err)
		}
	}

	if fields.Has(IfDBFieldLinkType) {
		cmd := GetFlowCmdByLinkStatus(link)

		r.log.Debugf("LINK: LinkType changed %s %v", cmd, link)

		if ifeOld.LinkType == fibcapi.LinkType_DEVICE {
			switch ifeNew.LinkType {
			case fibcapi.LinkType_BRIDGE_SLAVE:
				// DEVICE -> BRIDGE_SLAVE: pass
				if err := r.SendLinkFlows(fibcapi.FlowMod_DELETE, link); err != nil {
					r.log.Errorf("LINK: Delete LinkFlows error. %v %s", link, err)
				}

			case fibcapi.LinkType_BOND_SLAVE:
				// DEVICE -> BOND_SLAVE
				if err := r.SendBondSlaveFlows(cmd, link, &ifeNew); err != nil {
					r.log.Errorf("LINK: %s Link error. %v", cmd, link)
				}
			}
		}

		if ifeNew.LinkType == fibcapi.LinkType_DEVICE {
			switch ifeOld.LinkType {
			case fibcapi.LinkType_BRIDGE_SLAVE:
				// BRIDGE_SLAVE -> DEVICE: pass
				if err := r.SendLinkFlows(fibcapi.FlowMod_ADD, link); err != nil {
					r.log.Errorf("LINK: Add LinkFlows error. %v %s", link, err)
				}
			case fibcapi.LinkType_BOND_SLAVE:
				// BOND_SLAVE -> DEVICE
				if err := r.SendBondSlaveFlows(cmd, link, &ifeOld); err != nil {
					r.log.Errorf("LINK: %s Link error. %v", cmd, link)
				}
			}
		}
	}

	r.log.Debugf("LINK: OK (associated) %v", link)
}

func (r *RIBController) NetlinkAddr(nlmsg *nlamsg.NetlinkMessage, addr *nlamsg.Addr) {
	r.log.Debugf("ADDR: NId:%d AdId:%d", addr.NId, addr.AdId)

	cmd := GetFlowCmd(nlmsg.Type())
	if err := r.SendAddrFlows(cmd, addr); err != nil {
		r.log.Errorf("ADDR: %s error. %v %s", cmd, addr, err)
	}

	r.log.Debugf("ADDR: OK %s %v", cmd, addr)
}

func (r *RIBController) NetlinkNeigh(nlmsg *nlamsg.NetlinkMessage, neigh *nlamsg.Neigh) {
	r.log.Debugf("NEIGH: NId;%d NeId:%d", neigh.NId, neigh.NeId)

	if ok := r.ifdb.Associated(neigh.NId, neigh.LinkIndex); !ok {
		r.log.Warnf("NEIGH: Ifindex not found. Neigh %s", neigh)
		return
	}

	cmd := GetFlowCmd(nlmsg.Type())
	if err := r.SendNeighFlows(cmd, neigh); err != nil {
		r.log.Errorf("NEIGH: %s error. %v %s", cmd, neigh, err)
	}

	r.log.Debugf("NEIGH: OK %s %v", cmd, neigh)
}

func (r *RIBController) NetlinkRoute(nlmsg *nlamsg.NetlinkMessage, route *nlamsg.Route) {
	r.log.Debugf("ROUTE: NId:%d RtId:%d", route.NId, route.RtId)

	if ok := r.ifdb.Associated(route.NId, route.GetLinkIndex()); !ok {
		r.log.Warnf("ROUTE: Ifindex not found. Route %s", route)
		return
	}

	cmd := GetFlowCmd(nlmsg.Type())
	if err := r.SendRouteFlows(cmd, route); err != nil {
		r.log.Errorf("ROUTE: %s error. %v %s", cmd, route, err)
	}

	r.log.Debugf("ROUTE: OK %s %v", cmd, route)
}

func (r *RIBController) NetlinkBridgeVlanInfo(nlmsg *nlamsg.NetlinkMessage, brvlan *nlamsg.BridgeVlanInfo) {
	r.log.Debugf("BRVLAN: NId:%d BrId:%d", brvlan.NId, brvlan.BrId)

	cmd := GetFlowCmd(nlmsg.Type())
	if err := r.SendBridgeVlanFlows(cmd, brvlan); err != nil {
		r.log.Errorf("BRVLAN: %s error.", brvlan)
	}

	r.log.Debugf("BRVLAN: OK %s %v", cmd, brvlan)
}

func (r *RIBController) SendBondSlaveFlows(cmd fibcapi.FlowMod_Cmd, link *nlamsg.Link, ife *IfDBEntry) error {

	r.log.Debugf("BondSlaveFlows: %s %v", cmd, link)

	var ifeMaster IfDBEntry
	if ok := r.ifdb.SelectBy(&ifeMaster, ife.NId, ife.MasterIndex); !ok {
		log.Errorf("BondSlaveFlows: master not found. %d/%d", ife.NId, ife.MasterIndex)
		return fmt.Errorf("master not found. %d/%d", ife.NId, ife.MasterIndex)
	}

	log.Debugf("BondSlaveFlows: master nid:%d index:%d port_id:0x%08x",
		ifeMaster.NId, ife.MasterIndex, ifeMaster.PortId())

	grpCmd := FlowCmdToGroupCmd(cmd)

	if err := r.SendL2InterfaceGroup(grpCmd, link, &ifeMaster); err != nil {
		r.log.Errorf("BondSlaveFlows: L2 Interface Group error. %s", err)
		return err
	}

	return nil
}

func (r *RIBController) SendLinkFlows(cmd fibcapi.FlowMod_Cmd, link *nlamsg.Link) error {
	r.log.Debugf("LinkFlows: %s %v", cmd, link)

	grpCmd := FlowCmdToGroupCmd(cmd)

	var ifeMaster IfDBEntry
	if masterIndex := link.Attrs().MasterIndex; masterIndex != 0 {
		if ok := r.ifdb.SelectBy(&ifeMaster, link.NId, masterIndex); !ok {
			log.Warnf("LinkFlows: master not found. %d/%d %v", link.NId, masterIndex, link)
		}
	}

	r.log.Debugf("LinkFlows: master nid:%d index:%d port_id:0x%08x",
		ifeMaster.NId, ifeMaster.Index, ifeMaster.PortId())

	if grpCmd != fibcapi.GroupMod_DELETE {
		if err := r.SendL2InterfaceGroup(grpCmd, link, &ifeMaster); err != nil {
			r.log.Errorf("LinkFlows: L2 Interface Group error. %s", err)
			return err
		}
	}

	if err := r.SendVLANFlow(cmd, link); err != nil {
		r.log.Errorf("LinkFlows: VLAN Flow error. %s", err)
		return err
	}

	if err := r.SendTermMACFlow(cmd, link); err != nil {
		r.log.Errorf("LinkFlows: TermMAC Flow error. %s", err)
		return err
	}

	if err := r.SendACLFlowByLink(cmd, link); err != nil {
		r.log.Errorf("LinkFlows: PolicyACL(link) Flow error. %s", err)
		return err
	}

	if err := r.SendLoopbackFlows(cmd, link.NId, NewPortId(link)); err != nil {
		r.log.Errorf("LinkFlows: PolicyACL(Lo) Flow error. %s", err)
		return err
	}

	if grpCmd == fibcapi.GroupMod_DELETE {
		if err := r.SendL2InterfaceGroup(grpCmd, link, &ifeMaster); err != nil {
			r.log.Errorf("LinkFlows: L2 Interface Group error. %s", err)
			return err
		}
	}

	return nil
}

func (r *RIBController) SendAddrFlows(cmd fibcapi.FlowMod_Cmd, addr *nlamsg.Addr) error {
	r.log.Debugf("SendAddrFlowss: %s %s", cmd, addr)

	var ife IfDBEntry
	if ok := r.ifdb.SelectBy(&ife, addr.NId, int(addr.Index)); !ok {
		r.log.Errorf("SendAddrFlows: iface not found. %v", addr)
		return fmt.Errorf("iface not found. nid:%d index:%d", addr.NId, addr.Index)
	}

	switch ife.LinkType {
	case fibcapi.LinkType_BOND:
		// TODO: fix inPort=0 -> portId for each slave.
		// now we can use L2SW or LAG (not both).
		if err := r.SendACLFlowByAddr(cmd, addr, 0); err != nil {
			r.log.Errorf("PortStatus: add ACL(Addr) error. %s", err)
			return err
		}

	case fibcapi.LinkType_BOND_SLAVE:
		// pass

	default:
		if r.fib.FIBCType() == FIBCTypeTCP {
			if err := r.SendACLFlowByAddr(cmd, addr, 0); err != nil {
				// for openflow mode.
				r.log.Errorf("PortStatus: add ACL(Addr) error. port:0 %s", err)
			}
		} else {
			if err := r.SendACLFlowByAddr(cmd, addr, ife.PortId()); err != nil {
				r.log.Errorf("PortStatus: add ACL(Addr) error. port:%d %s", ife.PortId(), err)
				return err
			}
		}
	}

	return nil
}

func (r *RIBController) SendFdbFlows(cmd fibcapi.FlowMod_Cmd, neigh *nlamsg.Neigh) error {
	r.log.Debugf("FdbFlows: %s %s", cmd, neigh)

	var ife IfDBEntry
	if ok := r.ifdb.SelectBy(&ife, neigh.NId, neigh.LinkIndex); !ok {
		r.log.Errorf("FdbFlows: iface not found. %s", neigh)
		return fmt.Errorf("iface not found. %s", neigh)
	}

	if ife.LinkType != fibcapi.LinkType_BRIDGE_SLAVE {
		// link is not bridge slave.
		r.log.Debugf("FdbFlows: not bridge slave. %s", neigh)
		return nil
	}

	if err := r.SendBridgingFlow(cmd, neigh, ife.PortId()); err != nil {
		r.log.Errorf("FdbFlows: Bridge flow error. %s", err)
		return err
	}

	return nil
}

func (r *RIBController) SendNeighFlows(cmd fibcapi.FlowMod_Cmd, neigh *nlamsg.Neigh) error {
	if neigh.IsFdbEntry() {
		return r.SendFdbFlows(cmd, neigh)
	}

	r.log.Debugf("NeighFlows: %s %s", cmd, neigh)

	if neigh.NeId == 0 {
		r.log.Debugf("NeighFlows: ignore %s %s", cmd, neigh)
		return nil
	}

	grpCmd := FlowCmdToGroupCmd(cmd)

	if grpCmd == fibcapi.GroupMod_DELETE {
		if err := r.SendUnicastRoutingFlowNeigh(cmd, neigh); err != nil {
			r.log.Errorf("NeighFlows:: Unicast Routing(Neigh) error. %s", err)
			return err
		}
	}

	if err := r.SendL3UnicastGroup(grpCmd, neigh); err != nil {
		r.log.Errorf("NEighFlows: L3 Unicast Group error. %s", err)
		return err
	}

	if err := r.SendMPLSInterfaceGroup(grpCmd, neigh); err != nil {
		r.log.Errorf("NeighFlows: MPLS Interface error. %s", err)
		return err
	}

	if grpCmd != fibcapi.GroupMod_DELETE {
		if err := r.SendUnicastRoutingFlowNeigh(cmd, neigh); err != nil {
			r.log.Errorf("NeighFlows: Unicast Routing(Neigh) error. %s", err)
			return err
		}
	}

	return nil
}

func (r *RIBController) SendRouteFlows(cmd fibcapi.FlowMod_Cmd, route *nlamsg.Route) error {
	r.log.Debugf("RouteFlows: %s %v", cmd, route)

	if route.GetDst() != nil {
		if route.GetMPLSEncap() == nil {
			// IP Routing
			if err := r.SendUnicastRoutingFlow(cmd, route); err != nil {
				r.log.Errorf("RouteFlows: Unicast Routing(IP) error. %s", err)
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

			if grpCmd == fibcapi.GroupMod_DELETE {
				if err := r.SendUnicastRoutingFlowMPLS(cmd, route); err != nil {
					r.log.Errorf("RouteFlows: Unicast Routing(MPLS) error. %s", err)
					return err
				}
			}

			if err := r.SendMPLSLabelGroupMPLS(grpCmd, route); err != nil {
				r.log.Errorf("RouteFlows: MPLS Label Group(MPLS) error. %s", err)
				return err
			}

			if err := r.SendMPLSLabelGroupVPN(grpCmd, route); err != nil {
				r.log.Errorf("RouteFlows: MPLS Label Group(VPN) error. %s", err)
				return err
			}

			if grpCmd != fibcapi.GroupMod_DELETE {
				if err := r.SendUnicastRoutingFlowMPLS(cmd, route); err != nil {
					r.log.Errorf("RouteFlows: Unicast Routing(MPLS) error. %s", err)
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
				r.log.Errorf("RouteFlows: MPLS Flow(Pop1) error. %s", err)
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
				r.log.Errorf("RouteFlows: MPLS Flow(Pop2) error. %s", err)
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

			if grpCmd != fibcapi.GroupMod_DELETE {
				if err := r.SendMPLSLabelGroupSwap(grpCmd, route); err != nil {
					r.log.Errorf("RouteFlow: MPLS Label Group(Swap) error. %s", err)
					return err
				}
			}

			if err := r.SendMPLSFlowSwap(cmd, route, true); err != nil {
				r.log.Errorf("RouteFlow: MPLS Flow (Swap/BOS=1) error. %s", err)
				return err
			}

			if err := r.SendMPLSFlowSwap(cmd, route, false); err != nil {
				r.log.Errorf("RouteFlows: MPLS Flow (Swap/BOS=0) error. %s", err)
				return err
			}

			if grpCmd == fibcapi.GroupMod_DELETE {
				if err := r.SendMPLSLabelGroupSwap(grpCmd, route); err != nil {
					r.log.Errorf("RouteFlows: MPLS Label Group(Swap) error. %s", err)
					return err
				}
			}
		}
	}

	return nil
}

func (r *RIBController) SendLoopbackFlows(cmd fibcapi.FlowMod_Cmd, nid uint8, inPort uint32) error {
	r.log.Debugf("LoFlows: %s nid:%d in_port:%d", cmd, nid, inPort)

	links := make(map[int32]struct{}, 0)
	r.nla.GetLinks(nid, func(link *nlamsg.Link) error {
		if (link.Attrs().Flags & net.FlagLoopback) != 0 {
			links[int32(link.Attrs().Index)] = struct{}{}
		}
		return nil
	})

	return r.nla.GetAddrs(nid, func(addr *nlamsg.Addr) error {
		if _, ok := links[addr.Index]; ok {
			if err := r.SendACLFlowByAddr(cmd, addr, inPort); err != nil {
				r.log.Errorf("LoFlows: ACL FLow(Addr) error. %s", err)
				return err
			}
		}
		return nil
	})
}

func (r *RIBController) SendBridgeVlanFlows(cmd fibcapi.FlowMod_Cmd, brvlan *nlamsg.BridgeVlanInfo) error {
	r.log.Debugf("BrVlanFlows: %s %v", cmd, brvlan)

	var ife IfDBEntry
	if ok := r.ifdb.SelectBy(&ife, brvlan.NId, brvlan.Index); !ok {
		r.log.Errorf("BrVlanFlows: iface not found. %s", brvlan)
		return fmt.Errorf("iface not found. %s", brvlan)
	}

	if err := r.SendVLANBridgeVlanFlow(cmd, brvlan, ife.PortId()); err != nil {
		r.log.Errorf("BrVlanFlows: VLAN Flow(BrVlan) error. %s", err)
		return err
	}

	return nil
}
