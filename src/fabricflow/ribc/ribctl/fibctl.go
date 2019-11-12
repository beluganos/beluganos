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

package ribctl

import fibcapi "fabricflow/fibc/api"

const (
	FIBCTypeTCP     = "tcp"
	FIBCTypeGrpc    = "grpc"
	FIBCTypeDefault = FIBCTypeGrpc
)

type FIBController interface {
	Start() error
	Stop()
	Conn() <-chan bool
	Recv() <-chan *fibcapi.VmMonitorReply
	Hello(*fibcapi.Hello) error
	PortConfig(*fibcapi.PortConfig) error
	FlowMod(*fibcapi.FlowMod) error
	GroupMod(*fibcapi.GroupMod) error
	FIBCType() string
}

func NewFIBController(fibcType, addr, reId string) FIBController {
	switch fibcType {
	case FIBCTypeTCP:
		return NewFIBTcpController(addr)

	default:
		return NewFIBGrpcController(addr, reId)

	}
}
