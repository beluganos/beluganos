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

func (r *OAM_Request) Dispatch(xid uint32, i interface{}) error {
	hdr := fibcnet.Header{
		Type: uint16(FFM_OAM_REQUEST),
		Xid:  xid,
	}

	return DispatchOAMRequest(&hdr, r, i)
}

func DispatchOAMRequest(hdr *fibcnet.Header, msg *OAM_Request, i interface{}) error {
	if h, ok := i.(OAMRequestHandler); ok {
		h.FIBCOAMRequest(hdr, msg)
	}

	switch oam := msg.Body.(type) {
	case *OAM_Request_AuditRouteCnt:
		if h, ok := i.(OAMAuditRouteCntRequestHandler); ok {
			return h.FIBCOAMAuditRouteCntRequest(hdr, msg, oam.AuditRouteCnt)
		}

	default:
		return fmt.Errorf("Invalid oam type. %v", oam)
	}

	return fmt.Errorf("handler not implements. %d %s", msg.DpId, msg.OamType)
}

func (r *OAM_Reply) Dispatch(xid uint32, i interface{}) error {
	hdr := fibcnet.Header{
		Type: uint16(FFM_OAM_REPLY),
		Xid:  xid,
	}

	return DispatchOAMReply(&hdr, r, i)
}

func DispatchOAMReply(hdr *fibcnet.Header, msg *OAM_Reply, i interface{}) error {
	if h, ok := i.(OAMReplyHandler); ok {
		h.FIBCOAMReply(hdr, msg)
	}

	switch oam := msg.Body.(type) {
	case *OAM_Reply_AuditRouteCnt:
		if h, ok := i.(OAMAuditRouteCntReplyHandler); ok {
			return h.FIBCOAMAuditRouteCntReply(hdr, msg, oam.AuditRouteCnt)
		}

	default:
		return fmt.Errorf("Invalid oam type. %v", oam)
	}

	return fmt.Errorf("handler not implements. %d %s", msg.DpId, msg.OamType)
}
