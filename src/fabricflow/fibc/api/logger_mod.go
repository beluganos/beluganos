// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

	log "github.com/sirupsen/logrus"
)

func LogVALNFlow(logger LogLogger, level log.Level, flow *VLANFlow) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FloeMod(VLAN): match  inport: %d (0x%x)", flow.Match.InPort, flow.Match.InPort)
	logger.Logf(level, "FlowMod(VLAN): match  vid   : 0x%x/0x%x", flow.Match.Vid, flow.Match.VidMask)

	if actions := flow.Actions; actions != nil {
		for _, a := range actions {
			logger.Logf(level, "FlowMod(VLAN): action %s:%d", a.Name, a.Value)
		}
	}

	logger.Logf(level, "FlowMod(VLAN): goto   %d", flow.GotoTable)
}

func LogTerminationMacFlow(logger LogLogger, level log.Level, flow *TerminationMacFlow) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FlowMod(TermMac): match inport : %d (0x%x)", flow.Match.InPort, flow.Match.InPort)
	logger.Logf(level, "FlowMod(TermMac): match ethtype: 0x%04x", flow.Match.EthType)
	logger.Logf(level, "FlowMod(TermMac): match dstmac : '%s'", flow.Match.EthDst)
	logger.Logf(level, "FlowMod(TermMac): match vid    : %d", flow.Match.VlanVid)

	if actions := flow.Actions; actions != nil {
		for _, a := range actions {
			logger.Logf(level, "FlowMod(TermMac): action %s:%d", a.Name, a.Value)
		}
	}

	logger.Logf(level, "FlowMod(TermMac): goto   %d", flow.GotoTable)
}

func LogMPLSFlow(logger LogLogger, level log.Level, flow *MPLSFlow) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FlowMod(MPLS1): match  bos  : %t", flow.Match.Bos)
	logger.Logf(level, "FlowMod(MPLS1): match  label: %d", flow.Match.Label)

	if actions := flow.Actions; actions != nil {
		for _, a := range actions {
			logger.Logf(level, "FlowMod(MPLS1): action %s:%d", a.Name, a.Value)
		}
	}

	logger.Logf(level, "FlowMod(MPLS1): group  %s 0x%x", flow.GType, flow.GId)
	logger.Logf(level, "FlowMod(MPLS1): goto   %d", flow.GotoTable)
}

func LogUnicastRoutingFlow(logger LogLogger, level log.Level, flow *UnicastRoutingFlow) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FlowMod(U.C.): match  dip : '%s'", flow.Match.IpDst)
	logger.Logf(level, "FlowMod(U.C.): match  vrf : %d", flow.Match.Vrf)
	logger.Logf(level, "FlowMod(U.C.): match  orig: %s", flow.Match.Origin)

	if a := flow.Action; a != nil {
		logger.Logf(level, "FlowMod(U.C.): action %s:%d", a.Name, a.Value)
	}

	logger.Logf(level, "FlowMod(U.C.): group  %s 0x%x", flow.GType, flow.GId)
}

func LogBridgingFlow(logger LogLogger, level log.Level, flow *BridgingFlow) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FlowMod(Bridge): match dstmac: '%s'", flow.Match.EthDst)
	logger.Logf(level, "FlowMod(Bridge): match vid   : %d", flow.Match.VlanVid)
	logger.Logf(level, "FlowMod(Bridge): match tun   : %d", flow.Match.TunnelId)

	if a := flow.Action; a != nil {
		logger.Logf(level, "FlowMod(Bridge): action %s:%d", a.Name, a.Value)
	}
}

func LogPolicyACLFlow(logger LogLogger, level log.Level, flow *PolicyACLFlow) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FlowMod(ACL): match  inport : %d (0x%x)", flow.Match.InPort, flow.Match.InPort)
	logger.Logf(level, "FlowMod(ACL): match  vrf    : %d", flow.Match.Vrf)
	logger.Logf(level, "FlowMod(ACL): match  ethtype: 0x04%x", flow.Match.EthType)
	logger.Logf(level, "FlowMod(ACL): match  dstmac : '%s'", flow.Match.EthDst)
	logger.Logf(level, "FlowMod(ACL): match  proto  : %d", flow.Match.IpProto)
	logger.Logf(level, "FlowMod(ACL): match  dstip  : '%s'", flow.Match.IpDst)
	logger.Logf(level, "FlowMod(ACL): match  port   : %d->%d", flow.Match.TpSrc, flow.Match.TpDst)

	if a := flow.Action; a != nil {
		logger.Logf(level, "FlowMod(ACL): action %s:%d", a.Name, a.Value)
	}
}

func LogL2InterfaceGroup(logger LogLogger, level log.Level, g *L2InterfaceGroup) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "GroupMod(L2-IF): port  : %d (0x%x)", g.PortId, g.PortId)
	logger.Logf(level, "GroupMod(L2-IF): vrf   : %d", g.Vrf)
	logger.Logf(level, "GroupMod(L2-IF): vid   : %d", g.VlanVid)
	logger.Logf(level, "GroupMod(L2-IF): trans : %t", g.VlanTranslation)
	logger.Logf(level, "GroupMod(L2-IF): mac   : '%s'", g.HwAddr)
	logger.Logf(level, "GroupMod(L2-IF): mtu   : %d", g.Mtu)
	logger.Logf(level, "GroupMod(L2-IF): master: %d", g.Master)
}

func LogL3UnicastGroup(logger LogLogger, level log.Level, g *L3UnicastGroup) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "GroupMod(L3-UC): port  : %d (0x%x)", g.PortId, g.PortId)
	logger.Logf(level, "GroupMod(L3-UC): phy   : %d", g.PhyPortId)
	logger.Logf(level, "GroupMod(L3-UC): vid   : %d", g.VlanVid)
	logger.Logf(level, "GroupMod(L3-UC): neigh : %d", g.NeId)
	logger.Logf(level, "GroupMod(L3-UC): tun   : %s", g.TunType)
	logger.Logf(level, "GroupMod(L3-UC): remote:'%s'", g.EthDst)
	logger.Logf(level, "GroupMod(L3-UC): local :'%s'", g.EthSrc)
}

func LogMPLSInterfaceGroup(logger LogLogger, level log.Level, g *MPLSInterfaceGroup) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "GroupMod(MPLS-IF): port  : %d (0x%x)", g.PortId, g.PortId)
	logger.Logf(level, "GroupMod(MPLS-IF): vid   : %d", g.VlanVid)
	logger.Logf(level, "GroupMod(MPLS-IF): neigh : %d", g.NeId)
	logger.Logf(level, "GroupMod(MPLS-IF): dstmac: '%s'", g.EthDst)
	logger.Logf(level, "GroupMod(MPLS-IF): srcmac: '%s'", g.EthSrc)
}

func LogMPLSLabelGroup(logger LogLogger, level log.Level, g *MPLSLabelGroup) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "GroupMod(MPLS-LABEL): label  : %d->%d", g.DstId, g.NewLabel)
	logger.Logf(level, "GroupMod(MPLS-LABEL): neigh  : %d", g.NeId)
	logger.Logf(level, "GroupMod(MPLS-LABEL): encapid: %d", g.NewDstId)
	logger.Logf(level, "GroupMod(MPLS-LABEL): group  : %s", g.GType)
}

func LogFlowMod(logger LogLogger, level log.Level, mod *FlowMod) {
	if isSkipLog(level) {
		return
	}

	h := logModHandler{
		level:  level,
		logger: logger,
	}

	logger.Logf(h.level, "FlowMod: %s %s %s", mod.Cmd, mod.Table, mod.ReId)

	if err := mod.Dispatch(&h); err != nil {
		logger.Logf(level, "FlowMod: %s", mod.Entry)
	}
}

func LogGroupMod(logger LogLogger, level log.Level, mod *GroupMod) {
	if isSkipLog(level) {
		return
	}

	h := logModHandler{
		level:  level,
		logger: logger,
	}

	logger.Logf(h.level, "GroupMod: %s %s %s", mod.Cmd, mod.GType, mod.ReId)

	if err := mod.Dispatch(&h); err != nil {
		logger.Logf(level, "GroupMod: %s", mod.Entry)
	}
}

type logModHandler struct {
	level  log.Level
	logger LogLogger
}

func (h *logModHandler) FIBCVLANFlowMod(hdr *fibcnet.Header, mod *FlowMod, flow *VLANFlow) {
	LogVALNFlow(h.logger, h.level, flow)
}

func (h *logModHandler) FIBCTerminationMacFlowMod(hdr *fibcnet.Header, mod *FlowMod, flow *TerminationMacFlow) {
	LogTerminationMacFlow(h.logger, h.level, flow)
}

func (h *logModHandler) FIBCMPLSFlowMod(hdr *fibcnet.Header, mod *FlowMod, flow *MPLSFlow) {
	LogMPLSFlow(h.logger, h.level, flow)
}

func (h *logModHandler) FIBCUnicastRoutingFlowMod(hdr *fibcnet.Header, mod *FlowMod, flow *UnicastRoutingFlow) {
	LogUnicastRoutingFlow(h.logger, h.level, flow)
}

func (h *logModHandler) FIBCBridgingFlowMod(hdr *fibcnet.Header, mod *FlowMod, flow *BridgingFlow) {
	LogBridgingFlow(h.logger, h.level, flow)
}

func (h *logModHandler) FIBCPolicyACLFlowMod(hdr *fibcnet.Header, mod *FlowMod, flow *PolicyACLFlow) {
	LogPolicyACLFlow(h.logger, h.level, flow)
}

func (h *logModHandler) FIBCL2InterfaceGroupMod(hdr *fibcnet.Header, mod *GroupMod, grp *L2InterfaceGroup) {
	LogL2InterfaceGroup(h.logger, h.level, grp)
}

func (h *logModHandler) FIBCL3UnicastGroupMod(hdr *fibcnet.Header, mod *GroupMod, grp *L3UnicastGroup) {
	LogL3UnicastGroup(h.logger, h.level, grp)
}

func (h *logModHandler) FIBCMPLSInterfaceGroupMod(hdr *fibcnet.Header, mod *GroupMod, grp *MPLSInterfaceGroup) {
	LogMPLSInterfaceGroup(h.logger, h.level, grp)
}

func (h *logModHandler) FIBCMPLSLabelL2VpnGroupMod(hdr *fibcnet.Header, mod *GroupMod, grp *MPLSLabelGroup) {
	LogMPLSLabelGroup(h.logger, h.level, grp)
}
