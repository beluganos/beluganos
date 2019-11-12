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
// FFultipart_Request
//
func (m *FFMultipart_Request) Type() uint16 {
	return uint16(FFM_FF_MULTIPART_REQUEST)
}

func (m *FFMultipart_Request) Bytes() ([]byte, error) {
	return proto.Marshal(m)
}

func NewFFMultipart_RequestFromBytes(data []byte) (*FFMultipart_Request, error) {
	mp := &FFMultipart_Request{}
	if err := proto.Unmarshal(data, mp); err != nil {
		return nil, err
	}
	return mp, nil
}

//
// FFMultipart_Reply
//
func (m *FFMultipart_Reply) Type() uint16 {
	return uint16(FFM_FF_MULTIPART_REPLY)
}

func (m *FFMultipart_Reply) Bytes() ([]byte, error) {
	return proto.Marshal(m)
}

func NewFFMultipart_ReplyFromBytes(data []byte) (*FFMultipart_Reply, error) {
	mp := &FFMultipart_Reply{}
	if err := proto.Unmarshal(data, mp); err != nil {
		return nil, err
	}
	return mp, nil
}

//
// FFPortStats
//
func NewFFPortStats(port uint32, stats map[string]uint64) *FFPortStats {
	return &FFPortStats{
		PortNo: port,
		Values: stats,
	}
}

//
// Multipart Request (Port)
//
func newFFMultipart_Request_Port(portNo uint32) *FFMultipart_Request_Port {
	return &FFMultipart_Request_Port{
		Port: &FFMultipart_PortRequest{
			PortNo: portNo,
		},
	}
}

func NewFFMultipart_Request_Port(dpId uint64, portNo uint32) *FFMultipart_Request {
	return &FFMultipart_Request{
		DpId:   dpId,
		MpType: FFMultipart_PORT,
		Body:   newFFMultipart_Request_Port(portNo),
	}
}

func newFFMultipart_Reply_Port(stats []*FFPortStats) *FFMultipart_Reply_Port {
	return &FFMultipart_Reply_Port{
		Port: &FFMultipart_PortReply{
			Stats: stats,
		},
	}
}

func NewFFMultipart_Reply_Port(dpId uint64, stats []*FFPortStats) *FFMultipart_Reply {
	return &FFMultipart_Reply{
		DpId:   dpId,
		MpType: FFMultipart_PORT,
		Body:   newFFMultipart_Reply_Port(stats),
	}
}

//
// Multipart Request (PortDesc)
//
func newMultipart_Request_PortDesc(internal bool) *FFMultipart_Request_PortDesc {
	return &FFMultipart_Request_PortDesc{
		PortDesc: &FFMultipart_PortDescRequest{
			Internal: internal,
		},
	}
}

func NewMultipart_Request_PortDesc(dpId uint64, internal bool) *FFMultipart_Request {
	return &FFMultipart_Request{
		DpId:   dpId,
		MpType: FFMultipart_PORT_DESC,
		Body:   newMultipart_Request_PortDesc(internal),
	}
}

func newMultipart_Reply_PortDesc(ports []*FFPort, internal bool) *FFMultipart_Reply_PortDesc {
	return &FFMultipart_Reply_PortDesc{
		PortDesc: &FFMultipart_PortDescReply{
			Port:     ports,
			Internal: internal,
		},
	}
}

func NewMultipart_Reply_PortDesc(dpId uint64, ports []*FFPort, internal bool) *FFMultipart_Reply {
	return &FFMultipart_Reply{
		DpId:   dpId,
		MpType: FFMultipart_PORT_DESC,
		Body:   newMultipart_Reply_PortDesc(ports, internal),
	}
}

func NewFFMultipartRequest(dpId uint64) *FFMultipart_Request {
	return &FFMultipart_Request{
		DpId: dpId,
	}
}

func (r *FFMultipart_Request) SetPort(port *FFMultipart_PortRequest) *FFMultipart_Request {
	r.MpType = FFMultipart_PORT
	r.Body = &FFMultipart_Request_Port{
		Port: port,
	}

	return r
}

func (r *FFMultipart_Request) SetPortDesc(portDesc *FFMultipart_PortDescRequest) *FFMultipart_Request {
	r.MpType = FFMultipart_PORT_DESC
	r.Body = &FFMultipart_Request_PortDesc{
		PortDesc: portDesc,
	}

	return r
}

func NewFFMultipartPortRequest(portNo uint32, names []string) *FFMultipart_PortRequest {
	return &FFMultipart_PortRequest{
		PortNo: portNo,
		Names:  names,
	}
}

func NewFFMultipartPortDescRequest(internal bool) *FFMultipart_PortDescRequest {
	return &FFMultipart_PortDescRequest{
		Internal: internal,
	}
}

func NewDpMultipartRequest(xid uint32, request *FFMultipart_Request) *DpMultipartRequest {
	return &DpMultipartRequest{
		Xid:     xid,
		Request: request,
	}
}
