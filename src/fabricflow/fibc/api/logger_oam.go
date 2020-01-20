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

	log "github.com/sirupsen/logrus"
)

func LogOAMAuditRouteCntRequest(logger LogLogger, level log.Level, m *OAM_AuditRouteCntRequest) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "oam.AuditRouteCntRequest:")
}

func LogOAMAuditRouteCntReply(logger LogLogger, level log.Level, m *OAM_AuditRouteCntReply) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "oam.AuditRouteCntReply: count: %d", m.Count)
}

type logOAMHandler struct {
	level  log.Level
	logger LogLogger
}

func (h *logOAMHandler) FIBCOAMRequest(hdr *fibcnet.Header, oam *OAM_Request) error {
	if isSkipLog(h.level) {
		return nil
	}

	h.logger.Logf(h.level, "OAM.Request: dpid   : %d", oam.DpId)
	h.logger.Logf(h.level, "OAM.Request: oamtype: %s", oam.OamType)
	h.logger.Logf(h.level, "OAM.Request: xid    : %d", hdr.Xid)

	return nil
}

func (h *logOAMHandler) FIBCOAMReply(hdr *fibcnet.Header, oam *OAM_Reply) error {
	if isSkipLog(h.level) {
		return nil
	}

	h.logger.Logf(h.level, "OAM.Reply: dpid   : %d", oam.DpId)
	h.logger.Logf(h.level, "OAM.Reply: oamtype: %s", oam.OamType)
	h.logger.Logf(h.level, "OAM.Reply: xid    : %d", hdr.Xid)

	return nil
}

func (h *logOAMHandler) FIBCOAMAuditRouteCntRequest(hdr *fibcnet.Header, req *OAM_Request, oam *OAM_AuditRouteCntRequest) error {
	LogOAMAuditRouteCntRequest(h.logger, h.level, oam)
	return nil
}

func (h *logOAMHandler) FIBCOAMAuditRouteCntReply(hdr *fibcnet.Header, req *OAM_Reply, oam *OAM_AuditRouteCntReply) error {
	LogOAMAuditRouteCntReply(h.logger, h.level, oam)
	return nil
}

func LogOAMRequest(logger LogLogger, level log.Level, m *OAM_Request, xid uint32) {
	if isSkipLog(level) {
		return
	}

	h := logOAMHandler{
		level:  level,
		logger: logger,
	}

	if err := m.Dispatch(xid, &h); err != nil {
		logger.Logf(level, "OAM.Request: %s", m.Body)
	}
}

func LogOAMReply(logger LogLogger, level log.Level, m *OAM_Reply, xid uint32) {
	if isSkipLog(level) {
		return
	}

	h := logOAMHandler{
		level:  level,
		logger: logger,
	}

	if err := m.Dispatch(xid, &h); err != nil {
		logger.Logf(level, "OAM.Reply: %s", m.Body)
	}
}
