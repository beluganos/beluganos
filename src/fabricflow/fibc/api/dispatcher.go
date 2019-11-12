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
	"fmt"
)

func notifyMessage(header *fibcnet.Header, m fibcnet.Message, handler interface{}) error {
	if h, ok := handler.(MessageHandler); ok {
		return h.FIBCMessage(header, m)
	}

	return nil
}

func notifyhHello(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(HelloHandler); ok {
		m, err := NewHelloFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCHello(header, m)
	}

	return nil
}

func notifyPortStatus(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(PortStatusHandler); ok {
		m, err := NewPortStatusFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCPortStatus(header, m)
	}

	return nil
}

func notifyPortConfig(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(PortConfigHandler); ok {
		m, err := NewPortConfigFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCPortConfig(header, m)
	}

	return nil
}

func notifyFlowMod(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(FlowModHandler); ok {
		m, err := NewFlowModFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCFlowMod(header, m)
	}

	return nil
}

func notifyGroupMod(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(GroupModHandler); ok {
		m, err := NewGroupModFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCGroupMod(header, m)
	}

	return nil
}

func nodifyDpStatus(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(DpStatusHandler); ok {
		m, err := NewDpStatusFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCDpStatus(header, m)
	}

	return nil
}

func notifyFFHello(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(FFHelloHandler); ok {
		m, err := NewFFHelloFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCFFHello(header, m)
	}

	return nil
}

func notifyMultipartRequest(header *fibcnet.Header, data []byte, handler interface{}) error {
	m, err := NewFFMultipart_RequestFromBytes(data)
	if err != nil {
		return err
	}
	if err := notifyMessage(header, m, handler); err != nil {
		return err
	}

	if h, ok := handler.(FFMultipartRequestHandler); ok {
		if err := h.FIBCFFMultipartRequest(header, m); err != nil {
			return err
		}
	}

	switch m.MpType {
	case FFMultipart_PORT:
		if h, ok := handler.(FFMultipartPortRequestHandler); ok {
			h.FIBCFFMultipartPortRequest(header, m, m.GetPort())
		}
	case FFMultipart_PORT_DESC:
		if h, ok := handler.(FFMultipartPortDescRequestHandler); ok {
			h.FIBCFFMultipartPortDescRequest(header, m, m.GetPortDesc())
		}
	default:
		return fmt.Errorf("Invalid Mulripart Request. %d", m.MpType)
	}

	return nil
}

func notifyMultipartReply(header *fibcnet.Header, data []byte, handler interface{}) error {
	m, err := NewFFMultipart_ReplyFromBytes(data)
	if err != nil {
		return err
	}

	if err := notifyMessage(header, m, handler); err != nil {
		return err
	}

	if h, ok := handler.(FFMultipartReplyHandler); ok {
		if err := h.FIBCFFMultipartReply(header, m); err != nil {
			return err
		}
	}

	switch m.MpType {
	case FFMultipart_PORT:
		if h, ok := handler.(FFMultipartPortReplyHandler); ok {
			h.FIBCFFMultipartPortReply(header, m, m.GetPort())
		}
	case FFMultipart_PORT_DESC:
		if h, ok := handler.(FFMultipartPortDescReplyHandler); ok {
			h.FIBCFFMultipartPortDescReply(header, m, m.GetPortDesc())
		}
	default:
		return fmt.Errorf("Invalid Mulripart Reply. %d", m.MpType)
	}

	return nil
}

func notifyFFPacketIn(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(FFPacketInHandler); ok {
		m, err := NewFFPacketInFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCFFPacketIn(header, m)
	}

	return nil
}

func notifyFFPacketOut(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(FFPacketOutHandler); ok {
		m, err := NewFFPacketOutFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCFFPacketOut(header, m)
	}

	return nil
}

func notifyFFPortStatus(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(FFPortStatusHandler); ok {
		m, err := NewFFPortStatusFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCFFPortStatus(header, m)
	}

	return nil
}

func notifyFFPortMod(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(FFPortModHandler); ok {
		m, err := NewFFPortModFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCFFPortMod(header, m)
	}

	return nil
}

func notifyL2AddrStatus(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(L2AddrStatusHandler); ok {
		m, err := NewL2AddrStatusFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCL2AddrStatus(header, m)
	}

	return nil
}

func notifyFFL2AddrStatus(header *fibcnet.Header, data []byte, handler interface{}) error {
	if h, ok := handler.(FFL2AddrStatusHandler); ok {
		m, err := NewFFL2AddrStatusFromBytes(data)
		if err != nil {
			return err
		}

		if err := notifyMessage(header, m, handler); err != nil {
			return err
		}

		h.FIBCFFL2AddrStatus(header, m)
	}

	return nil
}

func Dispatch(h *fibcnet.Header, data []byte, handler interface{}) error {

	switch FFM(h.Type) {
	case FFM_HELLO:
		return notifyhHello(h, data, handler)
	case FFM_PORT_STATUS:
		return notifyPortStatus(h, data, handler)
	case FFM_PORT_CONFIG:
		return notifyPortConfig(h, data, handler)
	case FFM_FLOW_MOD:
		return notifyFlowMod(h, data, handler)
	case FFM_GROUP_MOD:
		return notifyGroupMod(h, data, handler)
	case FFM_DP_STATUS:
		return nodifyDpStatus(h, data, handler)
	case FFM_FF_HELLO:
		return notifyFFHello(h, data, handler)
	case FFM_FF_MULTIPART_REQUEST:
		return notifyMultipartRequest(h, data, handler)
	case FFM_FF_MULTIPART_REPLY:
		return notifyMultipartReply(h, data, handler)
	case FFM_FF_PACKET_IN:
		return notifyFFPacketIn(h, data, handler)
	case FFM_FF_PACKET_OUT:
		return notifyFFPacketOut(h, data, handler)
	case FFM_FF_PORT_STATUS:
		return notifyFFPortStatus(h, data, handler)
	case FFM_FF_PORT_MOD:
		return notifyFFPortMod(h, data, handler)
	case FFM_L2ADDR_STATUS:
		return notifyL2AddrStatus(h, data, handler)
	case FFM_FF_L2ADDR_STATUS:
		return notifyFFL2AddrStatus(h, data, handler)
	default:
		return fmt.Errorf("Invalid Message Type. %d", h.Type)
	}
}

func DispatchMsg(header *fibcnet.Header, m fibcnet.Message, handler interface{}) error {
	if err := notifyMessage(header, m, handler); err != nil {
		return err
	}

	switch FFM(header.Type) {
	case FFM_HELLO:
		if h, ok := handler.(HelloHandler); ok {
			h.FIBCHello(header, m.(*Hello))
		}
	case FFM_PORT_STATUS:
		if h, ok := handler.(PortStatusHandler); ok {
			h.FIBCPortStatus(header, m.(*PortStatus))
		}
	case FFM_PORT_CONFIG:
		if h, ok := handler.(PortConfigHandler); ok {
			h.FIBCPortConfig(header, m.(*PortConfig))
		}
	case FFM_FLOW_MOD:
		if h, ok := handler.(FlowModHandler); ok {
			h.FIBCFlowMod(header, m.(*FlowMod))
		}
	case FFM_GROUP_MOD:
		if h, ok := handler.(GroupModHandler); ok {
			h.FIBCGroupMod(header, m.(*GroupMod))
		}
	case FFM_DP_STATUS:
		if h, ok := handler.(DpStatusHandler); ok {
			h.FIBCDpStatus(header, m.(*DpStatus))
		}
	case FFM_FF_HELLO:
		if h, ok := handler.(FFHelloHandler); ok {
			h.FIBCFFHello(header, m.(*FFHello))
		}
	case FFM_FF_MULTIPART_REQUEST:
		if h, ok := handler.(FFMultipartRequestHandler); ok {
			return h.FIBCFFMultipartRequest(header, m.(*FFMultipart_Request))
		}
	case FFM_FF_MULTIPART_REPLY:
		if h, ok := handler.(FFMultipartReplyHandler); ok {
			return h.FIBCFFMultipartReply(header, m.(*FFMultipart_Reply))
		}
	case FFM_FF_PACKET_IN:
		if h, ok := handler.(FFPacketInHandler); ok {
			h.FIBCFFPacketIn(header, m.(*FFPacketIn))
		}
	case FFM_FF_PACKET_OUT:
		if h, ok := handler.(FFPacketOutHandler); ok {
			h.FIBCFFPacketOut(header, m.(*FFPacketOut))
		}
	case FFM_FF_PORT_STATUS:
		if h, ok := handler.(FFPortStatusHandler); ok {
			h.FIBCFFPortStatus(header, m.(*FFPortStatus))
		}
	case FFM_FF_PORT_MOD:
		if h, ok := handler.(FFPortModHandler); ok {
			h.FIBCFFPortMod(header, m.(*FFPortMod))
		}
	case FFM_L2ADDR_STATUS:
		if h, ok := handler.(L2AddrStatusHandler); ok {
			h.FIBCL2AddrStatus(header, m.(*L2AddrStatus))
		}
	case FFM_FF_L2ADDR_STATUS:
		if h, ok := handler.(FFL2AddrStatusHandler); ok {
			h.FIBCFFL2AddrStatus(header, m.(*FFL2AddrStatus))
		}
	default:
		return fmt.Errorf("Invalid Message Type. %d", header.Type)
	}

	return nil
}
