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

//
//
//
func NewVmMonitorReply() *VmMonitorReply {
	return &VmMonitorReply{}
}

func (r *VmMonitorReply) SetDpStatus(dpStatus *DpStatus) *VmMonitorReply {
	r.Body = &VmMonitorReply_DpStatus{
		DpStatus: dpStatus,
	}
	return r
}

func (r *VmMonitorReply) SetPortStatus(portStatus *PortStatus) *VmMonitorReply {
	r.Body = &VmMonitorReply_PortStatus{
		PortStatus: portStatus,
	}
	return r
}

func (r *VmMonitorReply) SetL2AddrStatus(l2addrStatus *L2AddrStatus) *VmMonitorReply {
	r.Body = &VmMonitorReply_L2AddrStatus{
		L2AddrStatus: l2addrStatus,
	}
	return r
}

func (r *VmMonitorReply) SetOAM(oam *OAMRequest) *VmMonitorReply {
	r.Body = &VmMonitorReply_Oam{
		Oam: oam,
	}
	return r
}

//
// NewDpMonitorReply returns new DpMonitorReply.
//
func NewDpMonitorReply() *DpMonitorReply {
	return &DpMonitorReply{}
}

//
// SetFFPacketOut set pktout to body.
//
func (r *DpMonitorReply) SetFFPacketOut(pktout *FFPacketOut) *DpMonitorReply {
	r.Body = &DpMonitorReply_PacketOut{
		PacketOut: pktout,
	}
	return r
}

//
// SetFFPortMod set mod to body.
//
func (r *DpMonitorReply) SetFFPortMod(mod *FFPortMod) *DpMonitorReply {
	r.Body = &DpMonitorReply_PortMod{
		PortMod: mod,
	}
	return r
}

//
// SetFlowMod set mod to body.
//
func (r *DpMonitorReply) SetFlowMod(mod *FlowMod) *DpMonitorReply {
	r.Body = &DpMonitorReply_FlowMod{
		FlowMod: mod,
	}
	return r
}

//
// SetGroupMod set mod to body.
//
func (r *DpMonitorReply) SetGroupMod(mod *GroupMod) *DpMonitorReply {
	r.Body = &DpMonitorReply_GroupMod{
		GroupMod: mod,
	}
	return r
}

//
// SetMultipart set multipart to body.
//
func (r *DpMonitorReply) SetMultipart(multipart *DpMultipartRequest) *DpMonitorReply {
	r.Body = &DpMonitorReply_Multipart{
		Multipart: multipart,
	}
	return r
}

func (r *DpMonitorReply) SetOAM(oam *OAMRequest) *DpMonitorReply {
	r.Body = &DpMonitorReply_Oam{
		Oam: oam,
	}
	return r
}

func NewDpMonitorRequest(dpId uint64, dpType FFHello_DpType) *DpMonitorRequest {
	return &DpMonitorRequest{
		DpId:   dpId,
		DpType: dpType,
	}
}

//
//
//
func NewVsMonitorReply() *VsMonitorReply {
	return &VsMonitorReply{}
}

func (r *VsMonitorReply) SetFFPacketOut(pktout *FFPacketOut) *VsMonitorReply {
	r.Body = &VsMonitorReply_PacketOut{
		PacketOut: pktout,
	}
	return r
}

func (r *VsMonitorReply) SetFFPortMod(portMod *FFPortMod) *VsMonitorReply {
	r.Body = &VsMonitorReply_PortMod{
		PortMod: portMod,
	}
	return r
}

func (r *VsMonitorReply) SetOAM(oam *OAMRequest) *VsMonitorReply {
	r.Body = &VsMonitorReply_Oam{
		Oam: oam,
	}
	return r
}

func NewOAMRequest(dpID uint64) *OAM_Request {
	return &OAM_Request{
		DpId: dpID,
	}
}

func (r *OAM_Request) SetAuditRouteCnt(audit *OAM_AuditRouteCntRequest) *OAM_Request {
	r.OamType = OAM_AUDIT_ROUTE_CNT
	r.Body = &OAM_Request_AuditRouteCnt{
		AuditRouteCnt: audit,
	}
	return r
}

func NewOAMReply(dpID uint64) *OAM_Reply {
	return &OAM_Reply{
		DpId: dpID,
	}
}

func NewOAMReplyVM(reID string) *OAM_Reply {
	return &OAM_Reply{
		ReId: reID,
	}
}

func (r *OAM_Reply) SetAuditRouteCnt(audit *OAM_AuditRouteCntReply) *OAM_Reply {
	r.OamType = OAM_AUDIT_ROUTE_CNT
	r.Body = &OAM_Reply_AuditRouteCnt{
		AuditRouteCnt: audit,
	}
	return r
}

func NewOAMAuditRouteCntRequest() *OAM_AuditRouteCntRequest {
	return &OAM_AuditRouteCntRequest{}
}

func NewOAMAuditRouteCntReply(count uint64) *OAM_AuditRouteCntReply {
	return &OAM_AuditRouteCntReply{
		Count: count,
	}
}
