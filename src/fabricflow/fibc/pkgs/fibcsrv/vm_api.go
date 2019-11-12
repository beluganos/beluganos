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
// VMAPIServer is VMAPI server
//
type VMAPIServer struct {
	ctl *VMCtl
}

//
// NewVMAPIServer returns new VMAPIServer
//
func NewVMAPIServer(ctl *VMCtl, server *grpc.Server) *VMAPIServer {
	s := &VMAPIServer{
		ctl: ctl,
	}
	fibcapi.RegisterFIBCVmApiServer(server, s)
	return s
}

//
// SendHello process hello message.
//
func (s *VMAPIServer) SendHello(ctxt context.Context, hello *fibcapi.Hello) (*fibcapi.HelloReply, error) {
	if err := s.ctl.Hello(hello); err != nil {
		return nil, err
	}

	return &fibcapi.HelloReply{}, nil
}

//
// SendPortConfig process port config message.
//
func (s *VMAPIServer) SendPortConfig(ctxt context.Context, portConfig *fibcapi.PortConfig) (*fibcapi.PortConfigReply, error) {
	if err := s.ctl.PortConfig(portConfig); err != nil {
		return nil, err
	}

	return &fibcapi.PortConfigReply{}, nil
}

//
// SendFlowMod process flow mod message.
//
func (s *VMAPIServer) SendFlowMod(ctxt context.Context, mod *fibcapi.FlowMod) (*fibcapi.FlowModReply, error) {
	if err := s.ctl.FlowMod(mod); err != nil {
		return nil, err
	}

	return &fibcapi.FlowModReply{}, nil
}

//
// SendGroupMod process group mod message.
//
func (s *VMAPIServer) SendGroupMod(ctxt context.Context, mod *fibcapi.GroupMod) (*fibcapi.GroupModReply, error) {
	if err := s.ctl.GroupMod(mod); err != nil {
		return nil, err
	}

	return &fibcapi.GroupModReply{}, nil
}

//
// Monitor process monitor message.
//
func (s *VMAPIServer) Monitor(req *fibcapi.VmMonitorRequest, stream fibcapi.FIBCVmApi_MonitorServer) error {
	return s.ctl.Monitor(req.ReId, stream, stream.Context().Done())
}
