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

package fibcsrv

import (
	fibcapi "fabricflow/fibc/api"
	"fabricflow/fibc/pkgs/fibcdbm"
)

//
// NewVMMonitorReplyPortStatus returns new VmMonitorReply
//
// VmMonitorReply {
//   oneof body {
//    PortStatus port_status
// }
//
// PortStatus {
//   Status status
//   string reID
//   uint32 portID
//   string ifname
// }
//
func NewVMMonitorReplyPortStatus(reID string, portID uint32, ifname string, status fibcapi.PortStatus_Status) *fibcapi.VmMonitorReply {
	return fibcapi.NewVmMonitorReply().SetPortStatus(
		fibcapi.NewPortStatus(
			reID,
			portID,
			ifname,
			status,
		),
	)
}

//
// NewVMMonitorReplyL2AddrStatus returns new VmMonitorReply
//
// VmMonitorReply {
//   oneof body {
//     L2AddrStatus l2_addr_status
//   }
// }
//
// L2AddrStatus {
//   string          reID
//   repeated L2Addr addrs
// }
//
func NewVMMonitorReplyL2AddrStatus(reID string, addrs []*fibcapi.L2Addr) *fibcapi.VmMonitorReply {
	return fibcapi.NewVmMonitorReply().SetL2AddrStatus(
		fibcapi.NewL2AddrStatus(reID, addrs),
	)
}

//
// NewDPMonitorReplyMpPort returns new DpMonitorReply
//
// DpMonitorReply {
//   oneof body {
//     DpMultipartRequest multipart
//   }
// }
//
// DpMultipartRequest {
//   uint32 xid
//   FFMultipart.Request request
// }
//
// FFMultipart {
//   Request {
//     uint64 dpID
//     MpType mp_type
//     oneof body {
//       PortRequest port
//     }
//   }
// }
//
// PortRequest {
//   uint32 portId
//   repeated string names
// }
//
func NewDPMonitorReplyMpPort(dpID uint64, portID uint32, names []string, xid uint32) *fibcapi.DpMonitorReply {
	return fibcapi.NewDpMonitorReply().SetMultipart(
		fibcapi.NewDpMultipartRequest(
			xid,
			fibcapi.NewFFMultipartRequest(dpID).SetPort(
				fibcapi.NewFFMultipartPortRequest(portID, names),
			),
		),
	)
}

//
// NewDPMonitorReplyMpPortDesc returns new DpMonitorReply
//
// DpMonitorReply {
//   oneof body {
//     DpMultipartRequest multipart
//   }
// }
//
// DpMultipartRequest {
//   uint32 xid
//   FFMultipart.Request request
// }
//
// FFMultipart {
//   Request {
//     uint64 dpID
//     MpType mp_type
//     oneof body {
//       PortDescRequest port_desc
//     }
//   }
// }
//
// PortDescRequest {
//   bool internal
// }
//
func NewDPMonitorReplyMpPortDesc(dpID uint64, internal bool, xid uint32) *fibcapi.DpMonitorReply {
	return fibcapi.NewDpMonitorReply().SetMultipart(
		fibcapi.NewDpMultipartRequest(
			xid,
			fibcapi.NewFFMultipartRequest(dpID).SetPortDesc(
				fibcapi.NewFFMultipartPortDescRequest(internal),
			),
		),
	)
}

//
// NewDPMonitorReplyMpPortDescInternal returns new DpMonitorReply
//
func NewDPMonitorReplyMpPortDescInternal(dpID uint64) *fibcapi.DpMonitorReply {
	return NewDPMonitorReplyMpPortDesc(dpID, true, 0)
}

//
// NewDPMonitorReplyFlowMod returns new DpMonitorReply
//
// DpMonitorReply {
//   oneof body {
//     FlowMod flow_mod
//   }
// }
//
func NewDPMonitorReplyFlowMod(mod *fibcapi.FlowMod) *fibcapi.DpMonitorReply {
	return fibcapi.NewDpMonitorReply().SetFlowMod(mod)
}

//
// NewDPMonitorReplyGroupMod returns new DpMonitorReply
//
// DpMonitorReply {
//   oneof body {
//     GroupMod group_mod
//   }
// }
//
func NewDPMonitorReplyGroupMod(mod *fibcapi.GroupMod) *fibcapi.DpMonitorReply {
	return fibcapi.NewDpMonitorReply().SetGroupMod(mod)
}

//
// NewDPMonitorReplyPacketOut retuens new DpMonitorReply
//
// DpMonitorReply {
//   oneof body {
//     FFPacketOut packet_out
//   }
// }
//
// FFPacketOut {
//   uint64 dpID
//   uint32 portID
//   bytes  data
// }
//
func NewDPMonitorReplyPacketOut(dpID uint64, portID uint32, data []byte) *fibcapi.DpMonitorReply {
	return fibcapi.NewDpMonitorReply().SetFFPacketOut(
		fibcapi.NewFFPacketOut(dpID, portID, data),
	)
}

//
// NewDPMonitorReplyPortMod returns new DpMonitorReply
//
// DpMonitorReply {
//   oneof body {
//     FFPortMod port_mod
//   }
// }
//
// FFPortMod {
//   uint64 dp_id
//   uint32 port_no
//   string hw_addr (not used at govsw mode)
//   PortStatus.Status status
// }
//
func NewDPMonitorReplyPortMod(dpID uint64, portID uint32, status fibcapi.PortStatus_Status) *fibcapi.DpMonitorReply {
	return fibcapi.NewDpMonitorReply().SetFFPortMod(
		&fibcapi.FFPortMod{
			DpId:   dpID,
			PortNo: portID,
			Status: status,
		},
	)
}

//
// NewVSMonitorReplyPacketOut returns new VsMonitorReply
//
// VsMonitorReply {
//   oneof body {
//     FFPacketOut packet_out
//   }
// }
//
// FFPacketOut {
//   uint64 dpID
//   uint32 portID
//   bytes  data
// }
//
func NewVSMonitorReplyPacketOut(vsID uint64, portID uint32, data []byte) *fibcapi.VsMonitorReply {
	return fibcapi.NewVsMonitorReply().SetFFPacketOut(
		fibcapi.NewFFPacketOut(vsID, portID, data),
	)

}

//
// NewVSMonitorReplyPortMod returns new VsMonitorReply
//
// VsMonitorReply {
//   oneof body {
//     FFPortMod port_mod
//   }
// }
//
// FFPortMod {
//   uint64 dp_id
//   uint32 port_no
//   string hw_addr (not used at govsw mode)
//   PortStatus.Status status
// }
//
func NewVSMonitorReplyPortMod(vsID uint64, portID uint32, status fibcapi.PortStatus_Status) *fibcapi.VsMonitorReply {
	return fibcapi.NewVsMonitorReply().SetFFPortMod(
		&fibcapi.FFPortMod{
			DpId:   vsID,
			PortNo: portID,
			Status: status,
		},
	)
}

//
// NewDBPortKeyFromLocal converts fibcdbm.PortKey to fibcapi.DbPortKey.
//
func NewDBPortKeyFromLocal(key *fibcdbm.PortKey) *fibcapi.DbPortKey {
	if key == nil {
		return &fibcapi.DbPortKey{}
	}
	return &fibcapi.DbPortKey{
		ReId:   key.ReID,
		Ifname: key.Ifname,
	}
}

//
// NewDBPortKeyFromAPI converts fibcapi.DbPortKey to fibcdbm.PortKey.
//
func NewDBPortKeyFromAPI(key *fibcapi.DbPortKey) *fibcdbm.PortKey {
	if len(key.ReId) == 0 && len(key.Ifname) == 0 {
		return nil
	}
	return &fibcdbm.PortKey{
		ReID:   key.ReId,
		Ifname: key.Ifname,
	}
}

//
// NewDBPortValueFromLocal converts fibcdbm.PortValue to fibcapi.DbPortValue.
//
func NewDBPortValueFromLocal(val *fibcdbm.PortValue) *fibcapi.DbPortValue {
	if val == nil {
		return &fibcapi.DbPortValue{}
	}

	return &fibcapi.DbPortValue{
		DpId:   val.DpID,
		ReId:   val.ReID,
		PortId: val.PortID,
		Enter:  val.Enter,
	}
}

//
// NewDBPortValueFromAPI converts fibcapi.DbPortValue to fibcdbm.PortValue.
//
func NewDBPortValueFromAPI(val *fibcapi.DbPortValue) *fibcdbm.PortValue {
	if val.DpId == 0 && len(val.ReId) == 0 {
		return nil
	}

	return &fibcdbm.PortValue{
		DpID:   val.DpId,
		ReID:   val.ReId,
		PortID: val.PortId,
		Enter:  val.Enter,
	}
}

//
// NewDBPortEntryFromLocal converts fibcdbm.PortEntry to fibcapi.DbPortEntry.
//
func NewDBPortEntryFromLocal(e *fibcdbm.PortEntry) *fibcapi.DbPortEntry {
	return &fibcapi.DbPortEntry{
		Key:       NewDBPortKeyFromLocal(e.Key),
		ParentKey: NewDBPortKeyFromLocal(e.ParentKey),
		MasterKey: NewDBPortKeyFromLocal(e.MasterKey),

		VmPort: NewDBPortValueFromLocal(e.VMPort),
		DpPort: NewDBPortValueFromLocal(e.DPPort),
		VsPort: NewDBPortValueFromLocal(e.VSPort),
	}
}

//
// NewDBPortEntryFromAPI converts fibcapi.DbPortEntry to fibcdbm.PortEntry.
//
func NewDBPortEntryFromAPI(e *fibcapi.DbPortEntry) *fibcdbm.PortEntry {
	return &fibcdbm.PortEntry{
		Key:       NewDBPortKeyFromAPI(e.Key),
		ParentKey: NewDBPortKeyFromAPI(e.ParentKey),
		MasterKey: NewDBPortKeyFromAPI(e.MasterKey),

		VMPort: NewDBPortValueFromAPI(e.VmPort),
		DPPort: NewDBPortValueFromAPI(e.DpPort),
		VSPort: NewDBPortValueFromAPI(e.VsPort),
	}
}

//
// NewDBIDEntryFromLocal convets fibcdbm.IDEntry to fibcapi.DbIdEntry
//
func NewDBIDEntryFromLocal(e *fibcdbm.IDEntry) *fibcapi.DbIdEntry {
	return &fibcapi.DbIdEntry{
		ReId: e.ReID,
		DpId: e.DpID,
	}
}

//
// NewDBIDEntryFromAPI converts fibcapi.DbIdEntry to fibcdbm.IDEntry.
//
func NewDBIDEntryFromAPI(e *fibcapi.DbIdEntry) *fibcdbm.IDEntry {
	return &fibcdbm.IDEntry{
		ReID: e.ReId,
		DpID: e.DpId,
	}
}
