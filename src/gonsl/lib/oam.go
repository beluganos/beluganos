// -*- coding: utf-8 -*-

// Copyright (C) 2020 Nippon Telegraph and Telephone Corporation.
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
	fibcnet "fabricflow/fibc/net"

	"github.com/beluganos/go-opennsl/opennsl"
	log "github.com/sirupsen/logrus"
)

const OAMAuditRouteCntTraverseNum = 1024

func (s *Server) fibcOAMAuditRouteCntReply(count uint64, xid uint32) {
	reply := fibcapi.OAM_Reply{
		DpId:    s.DpID(),
		OamType: fibcapi.OAM_AUDIT_ROUTE_CNT,
		Body: &fibcapi.OAM_Reply_AuditRouteCnt{
			AuditRouteCnt: &fibcapi.OAM_AuditRouteCntReply{
				Count: count,
			},
		},
	}

	s.log.Debugf("OAM(AuditRouteCnt): Send OAMReply. xid=%d count=%d", xid, count)

	if err := s.client.OAMReply(&reply, xid); err != nil {
		s.log.Errorf("OAM(AuditRouteCnt): Send OAMReply error. xid=%d %s", xid, err)
	}
}

func (s *Server) fibcOAMAuditRoutCnt(flag uint32) (uint64, error) {
	var (
		startIndex uint32
		count      uint64
	)

	for {
		endIndex := startIndex + OAMAuditRouteCntTraverseNum
		routeCount := uint64(0)

		s.log.Debugf("OAM(AuditRouteCnt): L3RouteTraverse flag=0x%x sta=%d ent=%d",
			flag, startIndex, endIndex)

		err := opennsl.L3RouteTraverse(s.Unit(), flag, startIndex, endIndex, func(rtUnit int, rtIndex int, route *opennsl.L3Route) opennsl.OpenNSLError {
			routeCount++
			return opennsl.E_NONE
		})

		if err != nil {
			s.log.Errorf("OAM(AuditRouteCnt): L3RouteTraverse flag=0x%x sta=%d ent=%d error. %s",
				flag, startIndex, endIndex, err)
			return 0, err
		}

		s.log.Debugf("OAM(AuditRouteCnt): L3RouteTraverse flag=0x%x sta=%d ent=%d route=%d",
			flag, startIndex, endIndex, routeCount)

		if routeCount == 0 {
			s.log.Debugf("OAM(AuditRouteCnt): L3RouteTraverse flag=0x%x sta=%d ent=%d end.",
				flag, startIndex, endIndex)
			break
		}

		count += routeCount
		startIndex += (OAMAuditRouteCntTraverseNum + 1)
	}

	s.log.Debugf("OAM(AuditRouteCnt): flag=0x%x count=%d", flag, count)

	return count, nil
}

func (s *Server) FIBCOAMAuditRouteCntRequest(hdr *fibcnet.Header, oam *fibcapi.OAM_Request, audit *fibcapi.OAM_AuditRouteCntRequest) error {
	s.log.Debugf("OAM(AuditRouteCnt): %v", hdr)
	fibcapi.LogOAMRequest(s.log, log.DebugLevel, oam, hdr.Xid)

	countV4, err := s.fibcOAMAuditRoutCnt(0)
	if err != nil {
		return err
	}

	countV6, err := s.fibcOAMAuditRoutCnt(uint32(opennsl.L3_IP6))
	if err != nil {
		return err
	}

	go s.fibcOAMAuditRouteCntReply(countV4+countV6, hdr.Xid)

	s.log.Debugf("OAM(AuditRouteCnt): end.")

	return nil
}
