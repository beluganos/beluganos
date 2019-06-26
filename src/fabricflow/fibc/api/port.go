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
)

//
// PortStatus
//
func (*PortStatus) Type() uint16 {
	return uint16(FFM_PORT_STATUS)
}

func (p *PortStatus) Bytes() ([]byte, error) {
	return proto.Marshal(p)
}

func NewPortStatus(reId string, portId uint32, ifname string, status PortStatus_Status) *PortStatus {
	return &PortStatus{
		ReId:   reId,
		PortId: portId,
		Ifname: ifname,
		Status: status,
	}
}

func NewPortStatusFromBytes(data []byte) (*PortStatus, error) {
	ps := &PortStatus{}
	if err := proto.Unmarshal(data, ps); err != nil {
		return nil, err
	}

	return ps, nil
}

//
// FFPort
//
const (
	FFPORT_STATE_NONE     = 0
	FFPORT_STATE_LINKDOWN = 0x01 // OFPPS_LINK_DOWN
	FFPORT_STATE_BLOCKED  = 0x02 // OFPPS_BLOCKED
	FFPORT_STATE_LIVE     = 0x04 // OFPPS_LIVE
)

func NewFFPort(portNo uint32) *FFPort {
	return &FFPort{
		PortNo: portNo,
		State:  FFPORT_STATE_NONE,
	}
}

//
// FFPortStatus
//
func (*FFPortStatus) Type() uint16 {
	return uint16(FFM_FF_PORT_STATUS)
}

func (p *FFPortStatus) Bytes() ([]byte, error) {
	return proto.Marshal(p)
}

func NewFFPortStatus(dpId uint64, port *FFPort, reason FFPortStatus_Reason) *FFPortStatus {
	return &FFPortStatus{
		DpId:   dpId,
		Port:   port,
		Reason: reason,
	}
}

func NewFFPortStatusFromBytes(data []byte) (*FFPortStatus, error) {
	ps := &FFPortStatus{}
	if err := proto.Unmarshal(data, ps); err != nil {
		return nil, err
	}

	return ps, nil
}

//
// PortConfig
//
func (*PortConfig) Type() uint16 {
	return uint16(FFM_PORT_CONFIG)
}

func (p *PortConfig) Bytes() ([]byte, error) {
	return proto.Marshal(p)
}

func NewPortConfig(cmd, reID, ifname string, portId uint32, status PortStatus_Status) *PortConfig {
	return &PortConfig{
		Cmd:    PortConfig_Cmd(PortConfig_Cmd_value[cmd]),
		ReId:   reID,
		Ifname: ifname,
		PortId: portId,
		Status: status,
		Link:   "",
	}
}

func NewPortConfigFromBytes(data []byte) (*PortConfig, error) {
	pc := &PortConfig{}
	if err := proto.Unmarshal(data, pc); err != nil {
		return nil, err
	}

	return pc, nil
}

//
// FFPortMod
//
func (*FFPortMod) Type() uint16 {
	return uint16(FFM_FF_PORT_MOD)
}

func (p *FFPortMod) Bytes() ([]byte, error) {
	return proto.Marshal(p)
}

func NewFFPortMod(dpId uint64, portNo uint32, status PortStatus_Status, hwaddr string) *FFPortMod {
	return &FFPortMod{
		DpId:   dpId,
		PortNo: portNo,
		HwAddr: hwaddr,
		Status: status,
	}
}

func NewFFPortModFromBytes(data []byte) (*FFPortMod, error) {
	pm := &FFPortMod{}
	if err := proto.Unmarshal(data, pm); err != nil {
		return nil, err
	}

	return pm, nil
}
