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

func NewPortStatusFromBytes(data []byte) (*PortStatus, error) {
	ps := &PortStatus{}
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

func NewPortConfig(cmd, reID, ifname string, value uint32) *PortConfig {
	return &PortConfig{
		Cmd:    PortConfig_Cmd(PortConfig_Cmd_value[cmd]),
		ReId:   reID,
		Ifname: ifname,
		Value:  value,
	}
}
