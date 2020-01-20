// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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

import (
	fibcapi "fabricflow/fibc/api"
	"gonla/nlaapi"
	"gonla/nlamsg"
)

func (r *RIBController) SendOAMAuditRouteCnt(xid uint32, count uint64) error {
	audit := fibcapi.NewOAMAuditRouteCntReply(count)
	reply := fibcapi.NewOAMReplyVM(r.reId).SetAuditRouteCnt(audit)
	if err := r.fib.OAMReply(reply, xid); err != nil {
		return err
	}
	return nil
}

func (r *RIBController) serveAuditRouteCnt(xid uint32) {
	var count uint64

	r.nla.GetRoutes(nlaapi.NODE_ID_ALL, func(route *nlamsg.Route) error {
		if route.GetDst() == nil || route.GetGw() == nil {
			return nil
		}
		if route.GetMPLSEncap() != nil {
			return nil
		}
		if ok := r.ifdb.Associated(route.NId, route.GetLinkIndex()); !ok {
			return nil
		}

		// r.log.Debugf("OAM(AuditRouteCnt): nid:%d %s oif:%d",
		//	route.NId, route.GetDst(), route.GetLinkIndex())

		count++

		return nil
	})

	if err := r.SendOAMAuditRouteCnt(xid, count); err != nil {
		r.log.Errorf("OAM(AuditRouteCnt): send error. xid:%d %s", xid, err)
	}
}

func (r *RIBController) StartAuditRouteCnt(xid uint32) error {
	go r.serveAuditRouteCnt(xid)
	return nil
}
