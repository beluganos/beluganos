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
	"fmt"
)

func (r *ApMonitorReply) Dispatch(i interface{}) error {
	return DispatchApMonitorReply(r, i)
}

func DispatchApMonitorReply(r *ApMonitorReply, i interface{}) error {
	if r == nil {
		return fmt.Errorf("Invalid reply. %v", r)
	}

	switch body := r.Body.(type) {
	case *ApMonitorReply_Log:
		if h, ok := i.(ApMonitorReplyLogHandler); ok {
			hdr := fibcnet.Header{
				Type: uint16(FFM_AP_MON_REPLY),
			}
			h.FIBCApMonitorReplyLog(&hdr, body.Log)
			return nil
		}

	default:
		return fmt.Errorf("Invalid type. %v", body)
	}

	return fmt.Errorf("handler not implements. %v", r)
}

func (r *VmMonitorReply) Dispatch(i interface{}) error {
	return DispatchVmMonitorReply(r, i)
}

func DispatchVmMonitorReply(r *VmMonitorReply, i interface{}) error {
	if r == nil {
		return fmt.Errorf("Invalid reply. %v", r)
	}

	switch body := r.Body.(type) {
	case *VmMonitorReply_PortStatus:
		if h, ok := i.(PortStatusHandler); ok {
			hdr := fibcnet.Header{Type: uint16(FFM_PORT_STATUS)}
			h.FIBCPortStatus(&hdr, body.PortStatus)
			return nil
		}

	case *VmMonitorReply_DpStatus:
		if h, ok := i.(DpStatusHandler); ok {
			hdr := fibcnet.Header{Type: uint16(FFM_DP_STATUS)}
			h.FIBCDpStatus(&hdr, body.DpStatus)
			return nil
		}

	case *VmMonitorReply_L2AddrStatus:
		if h, ok := i.(L2AddrStatusHandler); ok {
			hdr := fibcnet.Header{Type: uint16(FFM_L2ADDR_STATUS)}
			h.FIBCL2AddrStatus(&hdr, body.L2AddrStatus)
			return nil
		}

	default:
		return fmt.Errorf("Invalid type. %v", body)
	}

	return fmt.Errorf("handler not implements. %v", r)
}

func (r *VsMonitorReply) Dispatch(i interface{}) error {
	return DispatchVsMonitorReply(r, i)
}

func DispatchVsMonitorReply(r *VsMonitorReply, i interface{}) error {
	if r == nil {
		return fmt.Errorf("Invalid reply. %v", r)
	}

	switch body := r.Body.(type) {
	case *VsMonitorReply_PacketOut:
		if h, ok := i.(FFPacketOutHandler); ok {
			hdr := fibcnet.Header{Type: uint16(FFM_FF_PACKET_OUT)}
			h.FIBCFFPacketOut(&hdr, body.PacketOut)
			return nil
		}

	case *VsMonitorReply_PortMod:
		if h, ok := i.(FFPortModHandler); ok {
			hdr := fibcnet.Header{Type: uint16(FFM_FF_PORT_MOD)}
			h.FIBCFFPortMod(&hdr, body.PortMod)
			return nil
		}

	default:
		return fmt.Errorf("Invalid type. %v", body)
	}

	return fmt.Errorf("handler not implements. %v", r)
}

func (r *DpMonitorReply) Dispatch(i interface{}) error {
	hdr := fibcnet.Header{Type: uint16(FFM_DP_MON_REPLY)}
	return DispatchDpMonitorReply(&hdr, r, i)
}

func DispatchDpMonitorReply(hdr *fibcnet.Header, r *DpMonitorReply, i interface{}) error {
	if r == nil {
		return fmt.Errorf("Invalid reply. %v", r)
	}

	switch msg := r.Body.(type) {
	case *DpMonitorReply_PacketOut:
		if h, ok := i.(FFPacketOutHandler); ok {
			h.FIBCFFPacketOut(hdr, msg.PacketOut)
			return nil
		}

	case *DpMonitorReply_PortMod:
		if h, ok := i.(FFPortModHandler); ok {
			h.FIBCFFPortMod(hdr, msg.PortMod)
			return nil
		}

	case *DpMonitorReply_FlowMod:
		if h, ok := i.(FlowModHandler); ok {
			h.FIBCFlowMod(hdr, msg.FlowMod)
			return nil
		}

	case *DpMonitorReply_GroupMod:
		if h, ok := i.(GroupModHandler); ok {
			h.FIBCGroupMod(hdr, msg.GroupMod)
			return nil
		}

	case *DpMonitorReply_Multipart:
		hdr.Xid = msg.Multipart.Xid
		return DispatchFFMultipartRequest(
			hdr,
			msg.Multipart.Request,
			i,
		)

	default:
		return fmt.Errorf("Invalid type. %v", msg)
	}

	return fmt.Errorf("handler not implements. %v", r)
}
