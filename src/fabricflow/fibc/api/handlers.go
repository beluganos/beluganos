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
	fibcnet "fabricflow/fibc/net"
)

//
// Message
//

type MessageHandler interface {
	FIBCMessage(*fibcnet.Header, fibcnet.Message) error
}

//
// Hello
//
type HelloHandler interface {
	FIBCHello(*fibcnet.Header, *Hello)
}

//
// PortStatus
//
type PortStatusHandler interface {
	FIBCPortStatus(*fibcnet.Header, *PortStatus)
}

//
// PortConfig
//
type PortConfigHandler interface {
	FIBCPortConfig(*fibcnet.Header, *PortConfig)
}

//
// FlowMod
//
type FlowModHandler interface {
	FIBCFlowMod(*fibcnet.Header, *FlowMod)
}

//
// GroupMod
//
type GroupModHandler interface {
	FIBCGroupMod(*fibcnet.Header, *GroupMod)
}

//
// DpStatus
//
type DpStatusHandler interface {
	FIBCDpStatus(*fibcnet.Header, *DpStatus)
}

//
// FFHello
//
type FFHelloHandler interface {
	FIBCFFHello(*fibcnet.Header, *FFHello)
}

//
// PacketIn/Out
//
type FFPacketInHandler interface {
	FIBCFFPacketIn(*fibcnet.Header, *FFPacketIn)
}

type FFPacketOutHandler interface {
	FIBCFFPacketOut(*fibcnet.Header, *FFPacketOut)
}

//
// PortStatus
//
type FFPortStatusHandler interface {
	FIBCFFPortStatus(*fibcnet.Header, *FFPortStatus)
}

//
// PortMod
//
type FFPortModHandler interface {
	FIBCFFPortMod(*fibcnet.Header, *FFPortMod)
}

//
// L2AddrStatus
//
type L2AddrStatusHandler interface {
	FIBCL2AddrStatus(*fibcnet.Header, *L2AddrStatus)
}

//
// FFL2AddrStatus
//
type FFL2AddrStatusHandler interface {
	FIBCFFL2AddrStatus(*fibcnet.Header, *FFL2AddrStatus)
}

//
// ApMonitorReplyLog
//
type ApMonitorReplyLogHandler interface {
	FIBCApMonitorReplyLog(*fibcnet.Header, *ApMonitorReplyLog)
}
