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
// APAPIServer is APAPI server
//
type APAPIServer struct {
	ctl *APCtl
}

//
// NewAPAPIServer returns new APAPIServer
//
func NewAPAPIServer(ctl *APCtl, server *grpc.Server) *APAPIServer {
	s := &APAPIServer{
		ctl: ctl,
	}
	fibcapi.RegisterFIBCApApiServer(server, s)
	return s
}

//
// Monitor process monitor request.
//
func (s *APAPIServer) Monitor(req *fibcapi.ApMonitorRequest, stream fibcapi.FIBCApApi_MonitorServer) error {
	return s.ctl.Monitor(stream, stream.Context().Done())
}

//
// GetPortEntries process get port entries request.
//
func (s *APAPIServer) GetPortEntries(req *fibcapi.ApGetPortEntriesRequest, stream fibcapi.FIBCApApi_GetPortEntriesServer) error {
	return s.ctl.GetPortEntries(stream)
}

//
// AddPortEntry process add port entry request.
//
func (s *APAPIServer) AddPortEntry(ctxt context.Context, req *fibcapi.DbPortEntry) (*fibcapi.ApAddPortEntryReply, error) {
	if err := s.ctl.AddPortEntry(req); err != nil {
		return nil, err
	}

	return &fibcapi.ApAddPortEntryReply{}, nil
}

//
// DelPortEntry process delete port entry request.
//
func (s *APAPIServer) DelPortEntry(ctxt context.Context, req *fibcapi.DbPortKey) (*fibcapi.ApDelPortEntryReply, error) {
	if err := s.ctl.DelPortEntry(req); err != nil {
		return nil, err
	}

	return &fibcapi.ApDelPortEntryReply{}, nil
}

//
// GetIDEntries process get id entries request.
//
func (s *APAPIServer) GetIDEntries(req *fibcapi.ApGetIdEntriesRequest, stream fibcapi.FIBCApApi_GetIDEntriesServer) error {
	return s.ctl.GetIDEntries(stream)
}

//
// AddIDEntry process add id entry request.
//
func (s *APAPIServer) AddIDEntry(ctxt context.Context, req *fibcapi.DbIdEntry) (*fibcapi.ApAddIdEntryReply, error) {
	if err := s.ctl.AddIDEntry(req); err != nil {
		return nil, err
	}

	return &fibcapi.ApAddIdEntryReply{}, nil
}

//
// DelIDEntry process del id entry request.
//
func (s *APAPIServer) DelIDEntry(ctxt context.Context, req *fibcapi.DbIdEntry) (*fibcapi.ApDelIdEntryReply, error) {
	if err := s.ctl.DelIDEntry(req); err != nil {
		return nil, err
	}

	return &fibcapi.ApDelIdEntryReply{}, nil
}

//
// GetDpEntries process get dp entry request
//
func (s *APAPIServer) GetDpEntries(req *fibcapi.ApGetDpEntriesRequest, stream fibcapi.FIBCApApi_GetDpEntriesServer) error {
	return s.ctl.GetDpEntries(req.Type, stream)
}

//
// GetPortStats process get port stats request.
//
func (s *APAPIServer) GetPortStats(req *fibcapi.ApGetPortStatsRequest, stream fibcapi.FIBCApApi_GetPortStatsServer) error {
	return s.ctl.GetPortStats(req.DpId, req.PortNo, req.Names, stream)
}

//
// GetStats process get stats request.
//
func (s *APAPIServer) GetStats(req *fibcapi.ApGetStatsRequest, stream fibcapi.FIBCApApi_GetStatsServer) error {
	return s.ctl.GetStats(stream)
}
