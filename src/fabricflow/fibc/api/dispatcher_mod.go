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

package fibcapi

import (
	fibcnet "fabricflow/fibc/net"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (f *FlowMod) Dispatch(h interface{}) error {
	hdr := fibcnet.Header{Type: uint16(FFM_FLOW_MOD)}
	return DispatchFlowMod(&hdr, f, h)
}

func DispatchFlowMod(hdr *fibcnet.Header, mod *FlowMod, handler interface{}) error {

	switch mod.Table {
	case FlowMod_VLAN:
		if h, ok := handler.(FIBCVLANFlowModHandler); ok {
			h.FIBCVLANFlowMod(hdr, mod, mod.GetVlan())
			return nil
		}

	case FlowMod_TERM_MAC:
		if h, ok := handler.(FIBCTerminationMacFlowModHandler); ok {
			h.FIBCTerminationMacFlowMod(hdr, mod, mod.GetTermMac())
			return nil
		}

	case FlowMod_MPLS1:
		if h, ok := handler.(FIBCMPLSFlowModHandler); ok {
			h.FIBCMPLSFlowMod(hdr, mod, mod.GetMpls1())
			return nil
		}

	case FlowMod_UNICAST_ROUTING:
		if h, ok := handler.(FIBCUnicastRoutingFlowModHandler); ok {
			h.FIBCUnicastRoutingFlowMod(hdr, mod, mod.GetUnicast())
			return nil
		}

	case FlowMod_BRIDGING:
		if h, ok := handler.(FIBCBridgingFlowModHandler); ok {
			h.FIBCBridgingFlowMod(hdr, mod, mod.GetBridging())
			return nil
		}

	case FlowMod_POLICY_ACL:
		if h, ok := handler.(FIBCPolicyACLFlowModHandler); ok {
			h.FIBCPolicyACLFlowMod(hdr, mod, mod.GetAcl())
			return nil
		}

	default:
		log.Warnf("DispatchFlowMod: not dispatched. %s", mod.Table)
		return fmt.Errorf("invalid type. %s", mod.Table)
	}

	return fmt.Errorf("handler not implemented. %s", mod.Table)
}

func (g *GroupMod) Dispatch(h interface{}) error {
	hdr := fibcnet.Header{Type: uint16(FFM_GROUP_MOD)}
	return DispatchGroupMod(&hdr, g, h)
}

func DispatchGroupMod(hdr *fibcnet.Header, mod *GroupMod, handler interface{}) error {

	switch mod.GType {
	case GroupMod_L2_INTERFACE:
		if h, ok := handler.(FIBCL2InterfaceGroupModHandler); ok {
			h.FIBCL2InterfaceGroupMod(hdr, mod, mod.GetL2Iface())
			return nil
		}

	case GroupMod_L3_UNICAST:
		if h, ok := handler.(FIBCL3UnicastGroupModHandler); ok {
			h.FIBCL3UnicastGroupMod(hdr, mod, mod.GetL3Unicast())
			return nil
		}

	case GroupMod_MPLS_INTERFACE:
		if h, ok := handler.(FIBCMPLSInterfaceGroupModHandler); ok {
			h.FIBCMPLSInterfaceGroupMod(hdr, mod, mod.GetMplsIface())
			return nil
		}

	case GroupMod_MPLS_L2_VPN:
		if h, ok := handler.(FIBCMPLSLabelL2VpnGroupModHandler); ok {
			h.FIBCMPLSLabelL2VpnGroupMod(hdr, mod, mod.GetMplsLabel())
			return nil
		}

	case GroupMod_MPLS_L3_VPN:
		if h, ok := handler.(FIBCMPLSLabelL3VpnGroupModHandler); ok {
			h.FIBCMPLSLabelL3VpnGroupMod(hdr, mod, mod.GetMplsLabel())
			return nil
		}

	case GroupMod_MPLS_TUNNEL1:
		if h, ok := handler.(FIBCMPLSLabelTun1GroupModHandler); ok {
			h.FIBCMPLSLabelTun1GroupMod(hdr, mod, mod.GetMplsLabel())
			return nil
		}

	case GroupMod_MPLS_TUNNEL2:
		if h, ok := handler.(FIBCMPLSLabelTun2GroupModHandler); ok {
			h.FIBCMPLSLabelTun2GroupMod(hdr, mod, mod.GetMplsLabel())
			return nil
		}

	case GroupMod_MPLS_SWAP:
		if h, ok := handler.(FIBCMPLSLabelSwapGroupModHandler); ok {
			h.FIBCMPLSLabelSwapGroupMod(hdr, mod, mod.GetMplsLabel())
			return nil
		}

	default:
		log.Warnf("DispatchGroupMod: not dispatched. %s", mod.GType)
		return fmt.Errorf("invalid group type. %s", mod.GType)
	}

	return fmt.Errorf("handler not implemented. %s", mod.GType)
}
