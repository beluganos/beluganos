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

package fibclib

import (
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"
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
	FIBCHello(*fibcnet.Header, *fibcapi.Hello)
}

//
// PortStatus
//
type PortStatusHandler interface {
	FIBCPortStatus(*fibcnet.Header, *fibcapi.PortStatus)
}

//
// PortConfig
//
type PortConfigHandler interface {
	FIBCPortConfig(*fibcnet.Header, *fibcapi.PortConfig)
}

//
// FlowMod
//
type FlowModHandler interface {
	FIBCFlowMod(*fibcnet.Header, *fibcapi.FlowMod)
}

//
// GroupMod
//
type GroupModHandler interface {
	FIBCGroupMod(*fibcnet.Header, *fibcapi.GroupMod)
}

//
// DpStatus
//
type DpStatusHandler interface {
	FIBCDpStatus(*fibcnet.Header, *fibcapi.DpStatus)
}

//
// FFHello
//
type FFHelloHandler interface {
	FIBCFFHello(*fibcnet.Header, *fibcapi.FFHello)
}

//
// Multipart
//
type FFMultipartRequestHandler interface {
	FIBCFFMultipartRequest(*fibcnet.Header, *fibcapi.FFMultipart_Request) error
}

type FFMultipartReplyHandler interface {
	FIBCFFMultipartReply(*fibcnet.Header, *fibcapi.FFMultipart_Reply) error
}

//
// Multipart.Port
//
type FFMultipartPortRequestHandler interface {
	FIBCFFMultipartPortRequest(*fibcnet.Header, *fibcapi.FFMultipart_Request, *fibcapi.FFMultipart_PortRequest)
}

type FFMultipartPortReplyHandler interface {
	FIBCFFMultipartPortReply(*fibcnet.Header, *fibcapi.FFMultipart_Reply, *fibcapi.FFMultipart_PortReply)
}

//
// Multipart.PortDesc
//
type FFMultipartPortDescRequestHandler interface {
	FIBCFFMultipartPortDescRequest(*fibcnet.Header, *fibcapi.FFMultipart_Request, *fibcapi.FFMultipart_PortDescRequest)
}

type FFMultipartPortDescReplyHandler interface {
	FIBCFFMultipartPortDescReply(*fibcnet.Header, *fibcapi.FFMultipart_Reply, *fibcapi.FFMultipart_PortDescReply)
}

//
// PacketIn/Out
//
type FFPacketInHandler interface {
	FIBCFFPacketIn(*fibcnet.Header, *fibcapi.FFPacketIn)
}

type FFPacketOutHandler interface {
	FIBCFFPacketOut(*fibcnet.Header, *fibcapi.FFPacketOut)
}

//
// PortStatus
//
type FFPortStatusHandler interface {
	FIBCFFPortStatus(*fibcnet.Header, *fibcapi.FFPortStatus)
}

//
// PortMod
//
type FFPortModHandler interface {
	FIBCFFPortMod(*fibcnet.Header, *fibcapi.FFPortMod)
}
