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

package fibcapi

import (
	fibcnet "fabricflow/fibc/net"
	"fmt"
)

func (r *FFMultipart_Request) Dispatch(xid uint32, i interface{}) error {
	hdr := fibcnet.Header{
		Type: uint16(FFM_FF_MULTIPART_REQUEST),
		Xid:  xid,
	}

	return DispatchFFMultipartRequest(&hdr, r, i)
}

func DispatchFFMultipartRequest(hdr *fibcnet.Header, msg *FFMultipart_Request, i interface{}) error {
	if h, ok := i.(FFMultipartRequestHandler); ok {
		h.FIBCFFMultipartRequest(hdr, msg)
	}

	switch mp := msg.Body.(type) {
	case *FFMultipart_Request_Port:
		if h, ok := i.(FFMultipartPortRequestHandler); ok {
			h.FIBCFFMultipartPortRequest(hdr, msg, mp.Port)
			return nil
		}

	case *FFMultipart_Request_PortDesc:
		if h, ok := i.(FFMultipartPortDescRequestHandler); ok {
			h.FIBCFFMultipartPortDescRequest(hdr, msg, mp.PortDesc)
			return nil
		}

	default:
		return fmt.Errorf("Invalid mp type. %v", mp)
	}

	return fmt.Errorf("handler not implements. %d %s", msg.DpId, msg.MpType)
}

func (r *FFMultipart_Reply) Dispatch(xid uint32, i interface{}) error {
	hdr := fibcnet.Header{
		Type: uint16(FFM_FF_MULTIPART_REPLY),
		Xid:  xid,
	}

	return DispatchFFMultipartReply(&hdr, r, i)
}

func DispatchFFMultipartReply(hdr *fibcnet.Header, msg *FFMultipart_Reply, i interface{}) error {
	if h, ok := i.(FFMultipartReplyHandler); ok {
		h.FIBCFFMultipartReply(hdr, msg)
	}

	switch mp := msg.Body.(type) {
	case *FFMultipart_Reply_Port:
		if h, ok := i.(FFMultipartPortReplyHandler); ok {
			h.FIBCFFMultipartPortReply(hdr, msg, mp.Port)
			return nil
		}

	case *FFMultipart_Reply_PortDesc:
		if h, ok := i.(FFMultipartPortDescReplyHandler); ok {
			h.FIBCFFMultipartPortDescReply(hdr, msg, mp.PortDesc)
			return nil
		}

	default:
		return fmt.Errorf("Invalid mp type. %v", mp)
	}

	return fmt.Errorf("handler not implements. %d %s", msg.DpId, msg.MpType)
}
