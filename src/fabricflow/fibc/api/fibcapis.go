// -*- coding: utf-8 -*-

package fibcapi

import (
	fibcnet "fabricflow/fibc/net"
	"fmt"
)

func ParseDbDpEntryType(s string) (DbDpEntry_Type, error) {
	if v, ok := DbDpEntry_Type_value[s]; ok {
		return DbDpEntry_Type(v), nil
	}

	return DbDpEntry_NOP, fmt.Errorf("Invalid DbDpEntry_Type. %s", s)
}

func ParseVmMonitorReply(hdr *fibcnet.Header, data []byte) (*VmMonitorReply, error) {
	reply := VmMonitorReply{}
	if err := Dispatch(hdr, data, &reply); err != nil {
		return nil, err
	}

	return &reply, nil
}

func (r *VmMonitorReply) FIBCDpStatus(hdr *fibcnet.Header, msg *DpStatus) {
	r.Body = &VmMonitorReply_DpStatus{DpStatus: msg}
}

func (r *VmMonitorReply) FIBCPortStatus(hdr *fibcnet.Header, msg *PortStatus) {
	r.Body = &VmMonitorReply_PortStatus{PortStatus: msg}
}

func (r *VmMonitorReply) FIBCL2AddrStatus(hdr *fibcnet.Header, msg *L2AddrStatus) {
	r.Body = &VmMonitorReply_L2AddrStatus{L2AddrStatus: msg}
}

func ParseDpMonitprReply(hdr *fibcnet.Header, data []byte) (*DpMonitorReply, error) {
	reply := DpMonitorReply{}
	if err := Dispatch(hdr, data, &reply); err != nil {
		return nil, err
	}

	return &reply, nil
}

func (r *DpMonitorReply) FIBCFFPacketOut(hdr *fibcnet.Header, msg *FFPacketOut) {
	r.Body = &DpMonitorReply_PacketOut{
		PacketOut: msg,
	}
}

func (r *DpMonitorReply) FIBCFFPortMod(hdr *fibcnet.Header, msg *FFPortMod) {
	r.Body = &DpMonitorReply_PortMod{
		PortMod: msg,
	}
}

func (r *DpMonitorReply) FIBCFlowMod(hdr *fibcnet.Header, msg *FlowMod) {
	r.Body = &DpMonitorReply_FlowMod{
		FlowMod: msg,
	}
}

func (r *DpMonitorReply) FIBCGroupMod(hdr *fibcnet.Header, msg *GroupMod) {
	r.Body = &DpMonitorReply_GroupMod{
		GroupMod: msg,
	}
}

func (r *DpMonitorReply) FIBCFFMultipartRequest(hdr *fibcnet.Header, msg *FFMultipart_Request) error {
	r.Body = &DpMonitorReply_Multipart{
		Multipart: &DpMultipartRequest{
			Xid:     hdr.Xid,
			Request: msg,
		},
	}
	return nil
}
