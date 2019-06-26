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

package fibcapi

import (
	"net"

	"github.com/golang/protobuf/proto"
	"golang.org/x/sys/unix"
)

//
// FlowMod
//
func (f *FlowMod) Type() uint16 {
	return uint16(FFM_FLOW_MOD)
}

func (f *FlowMod) Bytes() ([]byte, error) {
	return proto.Marshal(f)
}

func NewFlowModFromBytes(data []byte) (*FlowMod, error) {
	flow_mod := &FlowMod{}
	if err := proto.Unmarshal(data, flow_mod); err != nil {
		return nil, err
	}

	return flow_mod, nil
}

//
// VLAN Flow Table
//
func NewVLANFlowMatch(inPort, vid, vidMask uint32) *VLANFlow_Match {
	return &VLANFlow_Match{
		InPort:  inPort,
		Vid:     vid,
		VidMask: vidMask,
	}
}

func NewVLANFlowAction(key string, value uint32) *VLANFlow_Action {
	if n, ok := VLANFlow_Action_Name_value[key]; ok {
		return &VLANFlow_Action{
			Name:  VLANFlow_Action_Name(n),
			Value: value,
		}
	}

	return nil
}

func NewVLANFlow(match *VLANFlow_Match, actions []*VLANFlow_Action, gotoTable uint32) *VLANFlow {
	return &VLANFlow{
		Match:     match,
		Actions:   actions,
		GotoTable: gotoTable,
	}
}

func (f *VLANFlow) ToMod(cmd FlowMod_Cmd, reId string) *FlowMod {
	return &FlowMod{
		Cmd:   cmd,
		Table: FlowMod_VLAN,
		ReId:  reId,
		Entry: &FlowMod_Vlan{Vlan: f},
	}
}

func (f VLANFlow) GetAction(name VLANFlow_Action_Name) *VLANFlow_Action {
	for _, action := range f.Actions {
		if action.GetName() == name {
			return action
		}
	}

	return nil
}

func (f VLANFlow) BridgeVlanInfoFlags() (BridgeVlanInfo_Flags, bool) {
	if a := f.GetAction(VLANFlow_Action_SET_VLAN_L2_TYPE); a != nil {
		return BridgeVlanInfo_Flags(a.GetValue()), true
	}

	return BridgeVlanInfo_NOP, false
}

//
// Termination MAC Flow Table
//
func NewTermMACMatch(inPort uint32, ethType uint32, ethDst string, vid uint16) *TerminationMacFlow_Match {
	return &TerminationMacFlow_Match{
		InPort:  inPort,
		EthType: ethType,
		EthDst:  ethDst,
		VlanVid: uint32(vid),
	}
}

func NewTermMACAction(key string, value uint32) *TerminationMacFlow_Action {
	if n, ok := TerminationMacFlow_Action_Name_value[key]; ok {
		return &TerminationMacFlow_Action{
			Name:  TerminationMacFlow_Action_Name(n),
			Value: value,
		}
	}
	return nil
}

func NewTermMACFlow(match *TerminationMacFlow_Match, actions []*TerminationMacFlow_Action, gotoTable uint32) *TerminationMacFlow {
	return &TerminationMacFlow{
		Match:     match,
		Actions:   actions,
		GotoTable: gotoTable,
	}
}

func (f *TerminationMacFlow) ToMod(cmd FlowMod_Cmd, reId string) *FlowMod {
	return &FlowMod{
		Cmd:   cmd,
		Table: FlowMod_TERM_MAC,
		ReId:  reId,
		Entry: &FlowMod_TermMac{TermMac: f},
	}
}

//
// Mpls Flow Table
//
func NewMPLSMatch(label uint32, bos bool) *MPLSFlow_Match {
	return &MPLSFlow_Match{
		Label: label,
		Bos:   bos,
	}
}

func NewMPLSAction(key string, value uint32) *MPLSFlow_Action {
	if n, ok := MPLSFlow_Action_Name_value[key]; ok {
		return &MPLSFlow_Action{
			Name:  MPLSFlow_Action_Name(n),
			Value: value,
		}
	}
	return nil
}

func NewMPLSFlow(match *MPLSFlow_Match, actions []*MPLSFlow_Action, gotoTable uint32, gtype GroupMod_GType, gid uint32) *MPLSFlow {
	return &MPLSFlow{
		Match:     match,
		Actions:   actions,
		GType:     gtype,
		GId:       gid,
		GotoTable: gotoTable,
	}
}

func (f *MPLSFlow) ToMod(cmd FlowMod_Cmd, reId string) *FlowMod {
	return &FlowMod{
		Cmd:   cmd,
		Table: FlowMod_MPLS1,
		ReId:  reId,
		Entry: &FlowMod_Mpls1{Mpls1: f},
	}
}

//
// Unicast Routing Flow Table
//
func NewUnicastRoutingMatch(ipDst *net.IPNet, vrf uint8, origin UnicastRoutingFlow_Origin) *UnicastRoutingFlow_Match {
	return &UnicastRoutingFlow_Match{
		IpDst:  ipDst.String(),
		Vrf:    uint32(vrf),
		Origin: origin,
	}
}

func NewUnicastRoutingMatchNeigh(ip net.IP, vrf uint8) *UnicastRoutingFlow_Match {
	ipDst := NewIPNetFromIP(ip)
	return NewUnicastRoutingMatch(ipDst, vrf, UnicastRoutingFlow_NEIGH)
}

func NewUnicastRoutingMatchRoute(ipDst *net.IPNet, vrf uint8) *UnicastRoutingFlow_Match {
	return NewUnicastRoutingMatch(ipDst, vrf, UnicastRoutingFlow_ROUTE)
}

func NewUnicastRoutingAction(key string, value uint32) *UnicastRoutingFlow_Action {
	if n, ok := UnicastRoutingFlow_Action_Name_value[key]; ok {
		return &UnicastRoutingFlow_Action{
			Name:  UnicastRoutingFlow_Action_Name(n),
			Value: value,
		}
	}

	return nil
}

func NewUnicastRoutingFlow(match *UnicastRoutingFlow_Match, action *UnicastRoutingFlow_Action, gtype GroupMod_GType, gid uint32) *UnicastRoutingFlow {
	return &UnicastRoutingFlow{
		Match:  match,
		Action: action,
		GType:  gtype,
		GId:    gid,
	}
}

func (f *UnicastRoutingFlow) ToMod(cmd FlowMod_Cmd, reId string) *FlowMod {
	return &FlowMod{
		Cmd:   cmd,
		Table: FlowMod_UNICAST_ROUTING,
		ReId:  reId,
		Entry: &FlowMod_Unicast{Unicast: f},
	}
}

//
// Bridging Flow Table
//

func (f *BridgingFlow) ToMod(cmd FlowMod_Cmd, reId string) *FlowMod {
	return &FlowMod{
		Cmd:   cmd,
		Table: FlowMod_BRIDGING,
		ReId:  reId,
		Entry: &FlowMod_Bridging{Bridging: f},
	}
}

func NewBridgingFlowMatch(ethDst string, vid uint16, tunId uint32) *BridgingFlow_Match {
	return &BridgingFlow_Match{
		EthDst:   ethDst,
		VlanVid:  uint32(vid),
		TunnelId: tunId,
	}
}

func NewBridgingFlowAction(key string, value uint32) *BridgingFlow_Action {
	if n, ok := BridgingFlow_Action_Name_value[key]; ok {
		return &BridgingFlow_Action{
			Name:  BridgingFlow_Action_Name(n),
			Value: value,
		}
	}

	return nil
}

func NewBridgingFlow(match *BridgingFlow_Match, action *BridgingFlow_Action) *BridgingFlow {
	return &BridgingFlow{
		Match:  match,
		Action: action,
	}
}

//
// Policy ACL Flow Table
//
func (f *PolicyACLFlow) ToMod(cmd FlowMod_Cmd, reId string) *FlowMod {
	return &FlowMod{
		Cmd:   cmd,
		Table: FlowMod_POLICY_ACL,
		ReId:  reId,
		Entry: &FlowMod_Acl{Acl: f},
	}
}

//
// Policy ACL Flow (match ip_dst and send controller)
//
func NewPolicyACLFlowByAddr(family int32, ipDst net.IP, vrf uint8, inPort uint32) *PolicyACLFlow {
	return &PolicyACLFlow{
		Match: &PolicyACLFlow_Match{
			InPort: inPort,
			IpDst:  ipDst.String(),
			Vrf:    uint32(vrf),
			EthType: func() uint32 {
				switch family {
				case unix.AF_INET:
					return unix.ETH_P_IP
				case unix.AF_INET6:
					return unix.ETH_P_IPV6
				default:
					return 0
				}
			}(),
		},
		Action: &PolicyACLFlow_Action{
			Name: PolicyACLFlow_Action_OUTPUT,
		},
	}
}
