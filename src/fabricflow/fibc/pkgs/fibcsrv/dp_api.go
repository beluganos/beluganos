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
// DPAPIServer is DPAPI server.
//
type DPAPIServer struct {
	ctl *DPCtl
}

//
// NewDPAPIServer returns new DPAPIServer
//
func NewDPAPIServer(ctl *DPCtl, server *grpc.Server) *DPAPIServer {
	s := &DPAPIServer{
		ctl: ctl,
	}
	fibcapi.RegisterFIBCDpApiServer(server, s)
	return s
}

//
// SendHello process Hello message.
//
func (s *DPAPIServer) SendHello(ctxt context.Context, hello *fibcapi.FFHello) (*fibcapi.FFHelloReply, error) {
	if err := s.ctl.Hello(hello.DpId, hello.DpType); err != nil {
		return nil, err
	}

	return &fibcapi.FFHelloReply{}, nil
}

//
// SendPacketIn process packet in message.
//
func (s *DPAPIServer) SendPacketIn(ctxt context.Context, pktin *fibcapi.FFPacketIn) (*fibcapi.FFPacketInReply, error) {
	if err := s.ctl.PacketIn(pktin.DpId, pktin.PortNo, pktin.Data); err != nil {
		return nil, err
	}

	return &fibcapi.FFPacketInReply{}, nil
}

//
// SendPortStatus process port status message.
//
func (s *DPAPIServer) SendPortStatus(ctxt context.Context, ptst *fibcapi.FFPortStatus) (*fibcapi.FFPortStatusReply, error) {
	if err := s.ctl.PortStatus(ptst.DpId, ptst.Port.PortNo, ptst.Port.State); err != nil {
		return nil, err
	}

	return &fibcapi.FFPortStatusReply{}, nil
}

//
// SendL2AddrStatus process l2 addr message.
//
func (s *DPAPIServer) SendL2AddrStatus(ctxt context.Context, l2addr *fibcapi.FFL2AddrStatus) (*fibcapi.L2AddrStatusReply, error) {
	if err := s.ctl.L2AddrStatus(l2addr.DpId, l2addr.Addrs); err != nil {
		return nil, err
	}

	return &fibcapi.L2AddrStatusReply{}, nil
}

//
// SendMultipartReply process multipart reply message
//
func (s *DPAPIServer) SendMultipartReply(ctxt context.Context, mp *fibcapi.DpMultipartReply) (*fibcapi.DpMultipartReplyAck, error) {

	if err := s.ctl.MultipartReply(mp.Xid, mp.Reply); err != nil {
		return nil, err
	}

	return &fibcapi.DpMultipartReplyAck{}, nil
}

//
// Monitor process monitor message.
//
func (s *DPAPIServer) Monitor(req *fibcapi.DpMonitorRequest, stream fibcapi.FIBCDpApi_MonitorServer) error {
	return s.ctl.Monitor(req.DpId, stream, stream.Context().Done())
}
