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
	"context"
	fibcapi "fabricflow/fibc/api"

	"google.golang.org/grpc"
)

//
// VSAPIServer is VSAPI server
//
type VSAPIServer struct {
	ctl *VSCtl
}

//
// NewVSAPIServer returns new VSAPIServer
//
func NewVSAPIServer(ctl *VSCtl, server *grpc.Server) *VSAPIServer {
	s := &VSAPIServer{
		ctl: ctl,
	}
	fibcapi.RegisterFIBCVsApiServer(server, s)
	return s
}

//
// SendHello process hello message.
//
func (s *VSAPIServer) SendHello(ctxt context.Context, hello *fibcapi.FFHello) (*fibcapi.FFHelloReply, error) {
	if err := s.ctl.Hello(hello.DpId, hello.DpType); err != nil {
		return nil, err
	}

	return &fibcapi.FFHelloReply{}, nil
}

//
// SendPacketIn process packet in message.
//
func (s *VSAPIServer) SendPacketIn(ctxt context.Context, pktin *fibcapi.FFPacketIn) (*fibcapi.FFPacketInReply, error) {
	if err := s.ctl.PacketIn(pktin.DpId, pktin.PortNo, pktin.Data); err != nil {
		return nil, err
	}

	return &fibcapi.FFPacketInReply{}, nil
}

//
// SendFFPacket process ffpacket message
//
func (s *VSAPIServer) SendFFPacket(ctxt context.Context, pkt *fibcapi.FFPacket) (*fibcapi.FFPacketReply, error) {
	if err := s.ctl.FFPacket(pkt.DpId, pkt.PortNo, pkt.ReId, pkt.Ifname); err != nil {
		return nil, err
	}

	return &fibcapi.FFPacketReply{}, nil
}

//
// Monitor process monitor message.
//
func (s *VSAPIServer) Monitor(req *fibcapi.VsMonitorRequest, stream fibcapi.FIBCVsApi_MonitorServer) error {
	return s.ctl.Monitor(req.VsId, stream, stream.Context().Done())
}

//
// SendOAMReply process oam request message.
//
func (s *VSAPIServer) SendOAMReply(ctxt context.Context, oam *fibcapi.OAMReply) (*fibcapi.OAMReplyAck, error) {
	if err := s.ctl.OAMReply(oam.Xid, oam.Reply); err != nil {
		return nil, err
	}
	return &fibcapi.OAMReplyAck{}, nil
}
