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
	"github.com/golang/protobuf/proto"
)

//
// FFPacketIn
//
func (p *FFPacketIn) Type() uint16 {
	return uint16(FFM_FF_PACKET_IN)
}

func (p *FFPacketIn) Bytes() ([]byte, error) {
	return proto.Marshal(p)
}

func (p *FFPacketIn) DataLen() uint32 {
	return uint32(len(p.Data))
}

func NewFFPacketIn(dpId uint64, portNo uint32, data []byte) *FFPacketIn {
	return &FFPacketIn{
		DpId:   dpId,
		PortNo: portNo,
		Data:   data,
	}
}
func NewFFPacketInFromBytes(data []byte) (*FFPacketIn, error) {
	p := &FFPacketIn{}
	if err := proto.Unmarshal(data, p); err != nil {
		return nil, err
	}

	return p, nil
}

//
// FFPacketOut
//
func (p *FFPacketOut) Type() uint16 {
	return uint16(FFM_FF_PACKET_IN)
}

func (p *FFPacketOut) Bytes() ([]byte, error) {
	return proto.Marshal(p)
}

func (p *FFPacketOut) DataLen() uint32 {
	return uint32(len(p.Data))
}

func NewFFPacketOut(dpId uint64, portNo uint32, data []byte) *FFPacketOut {
	return &FFPacketOut{
		DpId:   dpId,
		PortNo: portNo,
		Data:   data,
	}
}

func NewFFPacketOutFromBytes(data []byte) (*FFPacketOut, error) {
	p := &FFPacketOut{}
	if err := proto.Unmarshal(data, p); err != nil {
		return nil, err
	}

	return p, nil
}
