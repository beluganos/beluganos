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

package gonslib

import (
	fibcapi "fabricflow/fibc/api"
	"fmt"
)

const (
	FIBCTypeTCP     = "tcp"
	FIBCTypeGRPC    = "grpc"
	FIBCTypeDefault = FIBCTypeGRPC
)

type FIBController interface {
	Start() error
	Stop()
	Conn() <-chan bool
	Recv() <-chan *fibcapi.DpMonitorReply
	Hello(*fibcapi.FFHello) error
	PacketIn(*fibcapi.FFPacketIn) error
	PortStatus(*fibcapi.FFPortStatus) error
	L2AddrStatus(*fibcapi.FFL2AddrStatus) error
	MultipartReply(*fibcapi.FFMultipart_Reply, uint32) error
	OAMReply(*fibcapi.OAM_Reply, uint32) error

	fmt.Stringer
}

func NewFIBController(fibcType, addr string, dpId uint64) FIBController {
	switch fibcType {
	case FIBCTypeTCP:
		return NewFIBTcpController(addr, dpId)

	default:
		return NewFIBGrpcController(addr, dpId)
	}
}
