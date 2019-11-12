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
)

// VLANFlow
type FIBCVLANFlowModHandler interface {
	FIBCVLANFlowMod(*fibcnet.Header, *FlowMod, *VLANFlow)
}

// TerminationMacFlow
type FIBCTerminationMacFlowModHandler interface {
	FIBCTerminationMacFlowMod(*fibcnet.Header, *FlowMod, *TerminationMacFlow)
}

// MPLSFlow
type FIBCMPLSFlowModHandler interface {
	FIBCMPLSFlowMod(*fibcnet.Header, *FlowMod, *MPLSFlow)
}

// UnicastRoutingFlow
type FIBCUnicastRoutingFlowModHandler interface {
	FIBCUnicastRoutingFlowMod(*fibcnet.Header, *FlowMod, *UnicastRoutingFlow)
}

// BridgingFlow
type FIBCBridgingFlowModHandler interface {
	FIBCBridgingFlowMod(*fibcnet.Header, *FlowMod, *BridgingFlow)
}

// PolicyACLFlow
type FIBCPolicyACLFlowModHandler interface {
	FIBCPolicyACLFlowMod(*fibcnet.Header, *FlowMod, *PolicyACLFlow)
}

// L2InterfaceGroup
type FIBCL2InterfaceGroupModHandler interface {
	FIBCL2InterfaceGroupMod(*fibcnet.Header, *GroupMod, *L2InterfaceGroup)
}

// L3UnicastGroup
type FIBCL3UnicastGroupModHandler interface {
	FIBCL3UnicastGroupMod(*fibcnet.Header, *GroupMod, *L3UnicastGroup)
}

// MPLSInterfaceGroup
type FIBCMPLSInterfaceGroupModHandler interface {
	FIBCMPLSInterfaceGroupMod(*fibcnet.Header, *GroupMod, *MPLSInterfaceGroup)
}

// MPLSLabelGroup
type FIBCMPLSLabelL2VpnGroupModHandler interface {
	FIBCMPLSLabelL2VpnGroupMod(*fibcnet.Header, *GroupMod, *MPLSLabelGroup)
}

type FIBCMPLSLabelL3VpnGroupModHandler interface {
	FIBCMPLSLabelL3VpnGroupMod(*fibcnet.Header, *GroupMod, *MPLSLabelGroup)
}

type FIBCMPLSLabelTun1GroupModHandler interface {
	FIBCMPLSLabelTun1GroupMod(*fibcnet.Header, *GroupMod, *MPLSLabelGroup)
}

type FIBCMPLSLabelTun2GroupModHandler interface {
	FIBCMPLSLabelTun2GroupMod(*fibcnet.Header, *GroupMod, *MPLSLabelGroup)
}

type FIBCMPLSLabelSwapGroupModHandler interface {
	FIBCMPLSLabelSwapGroupMod(*fibcnet.Header, *GroupMod, *MPLSLabelGroup)
}
