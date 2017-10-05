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
	"github.com/golang/protobuf/proto"
	"net"
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

//
// Termination MAC Flow Table
//
func NewTermMACMatch(ethType uint32, ethDst string) *TerminationMacFlow_Match {
	return &TerminationMacFlow_Match{
		InPort:  0,
		EthType: ethType,
		EthDst:  ethDst,
		VlanVid: 0,
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
func NewUnicastRoutingMatch(ipDst *net.IPNet, vrf uint8) *UnicastRoutingFlow_Match {
	return &UnicastRoutingFlow_Match{
		IpDst: ipDst.String(),
		Vrf:   uint32(vrf),
	}
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
func NewPolicyACLFlowByAddr(ipDst net.IP, vrf uint8) *PolicyACLFlow {
	return &PolicyACLFlow{
		Match: &PolicyACLFlow_Match{
			IpDst: ipDst.String(),
			Vrf:   uint32(vrf),
		},
		Action: &PolicyACLFlow_Action{
			Name: PolicyACLFlow_Action_OUTPUT,
		},
	}
}
