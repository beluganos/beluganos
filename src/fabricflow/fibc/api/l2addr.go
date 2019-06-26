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
)

//
// L2Addr
//
func NewL2AddrVM(hwaddr net.HardwareAddr, vid uint16, portId uint32, reason L2Addr_Reason, ifname string) *L2Addr {
	return &L2Addr{
		HwAddr:  hwaddr.String(),
		VlanVid: uint32(vid),
		PortId:  portId,
		Reason:  reason,
		Ifname:  ifname,
	}
}

func NewL2AddrDP(hwaddr net.HardwareAddr, vid uint16, portId uint32, reason L2Addr_Reason) *L2Addr {
	return NewL2AddrVM(hwaddr, vid, portId, reason, "")
}

func NewL2AddrFromBytes(data []byte) (*L2Addr, error) {
	a := &L2Addr{}
	if err := proto.Unmarshal(data, a); err != nil {
		return nil, err
	}

	return a, nil
}

//
// L2AddrStatus
//
func (*L2AddrStatus) Type() uint16 {
	return uint16(FFM_L2ADDR_STATUS)
}

func (a *L2AddrStatus) Bytes() ([]byte, error) {
	return proto.Marshal(a)
}

func NewL2AddrStatus(reId string, addrs []*L2Addr) *L2AddrStatus {
	return &L2AddrStatus{
		ReId:  reId,
		Addrs: addrs,
	}
}

func NewL2AddrStatusFromBytes(data []byte) (*L2AddrStatus, error) {
	a := &L2AddrStatus{}
	if err := proto.Unmarshal(data, a); err != nil {
		return nil, err
	}

	return a, nil
}

//
// FFL2AddrStatus
//
func (*FFL2AddrStatus) Type() uint16 {
	return uint16(FFM_FF_L2ADDR_STATUS)
}

func (a *FFL2AddrStatus) Bytes() ([]byte, error) {
	return proto.Marshal(a)
}

func NewFFL2AddrStatus(dpId uint64, addrs []*L2Addr) *FFL2AddrStatus {
	return &FFL2AddrStatus{
		DpId:  dpId,
		Addrs: addrs,
	}
}

func NewFFL2AddrStatusFromBytes(data []byte) (*FFL2AddrStatus, error) {
	a := &FFL2AddrStatus{}
	if err := proto.Unmarshal(data, a); err != nil {
		return nil, err
	}

	return a, nil
}
