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

package nlamsg

import (
	"fmt"
	"gonla/nlalib"
	"gonla/nlamsg/nlalink"
)

func CopyVpn(src *nlalink.Vpn) *nlalink.Vpn {
	dst := *src
	return &dst
}

//
// Vpn
//
type Vpn struct {
	*nlalink.Vpn
	VpnId uint32 // auto increment
	NId   uint8
}

func (v *Vpn) Copy() *Vpn {
	return &Vpn{
		Vpn:   CopyVpn(v.Vpn),
		VpnId: v.VpnId,
		NId:   v.NId,
	}
}

func (v *Vpn) String() string {
	return fmt.Sprintf("{Dst: %s GW: %s VpnGW: %s Label: %d} VpnId: %d NId: %d",
		v.GetIPNet(), v.NetGw(), v.NetVpnGw(), v.Label, v.VpnId, v.NId)
}

func NewVpn(vpn *nlalink.Vpn, nid uint8, id uint32) *Vpn {
	return &Vpn{
		VpnId: id,
		Vpn:   vpn,
		NId:   nid,
	}
}

func VpnDeserialize(nlmsg *NetlinkMessage) (*Vpn, error) {
	vpn, err := nlalink.VpnDeserialize(nlmsg.Data)
	if err != nil {
		return nil, err
	}

	return NewVpn(vpn, nlmsg.NId, 0), nil
}

func VpnSerialize(vpn *Vpn, msgType uint16) (*NetlinkMessage, error) {
	data, err := nlalink.VpnSerialize(vpn.Vpn)
	if err != nil {
		return nil, err
	}
	return NewNetlinkMessage(nlalib.NewNetlinkMessage(msgType, data), vpn.NId, SRC_API), nil
}
