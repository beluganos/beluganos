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

func LogFFMultipartPortRequest(logger LogLogger, level log.Level, m *FFMultipart_PortRequest) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "MP.PortRequest: port: %d", m.PortNo)
	if names := m.Names; names != nil {
		for _, name := range names {
			logger.Logf(level, "MP.PortRequest: name: '%s'", name)
		}
	}
}

func LogFFMultipartPortDescRequest(logger LogLogger, level log.Level, m *FFMultipart_PortDescRequest) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "MP.PortDescRequest: internal: %t", m.Internal)
}

func LogFFMultipartPortReply(logger LogLogger, level log.Level, m *FFMultipart_PortReply) {
	if isSkipLog(level) {
		return
	}

	if stats := m.Stats; stats != nil {
		for index, stat := range stats {
			logger.Logf(level, "MP.PortReply[%d]:", index)
			LogFFPortStats(logger, level, stat)
		}
	}
}

func LogFFMultipartPortDescReply(logger LogLogger, level log.Level, m *FFMultipart_PortDescReply) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "MP.PortDescReply: internal: %t", m.Internal)
	if ports := m.Port; ports != nil {
		for index, port := range ports {
			logger.Logf(level, "MP.PortDescReply[%d]: ", index)
			LogFFPort(logger, level, port)
		}
	}
}

type logMPHandler struct {
	level  log.Level
	logger LogLogger
}

func (h *logMPHandler) FIBCFFMultipartRequest(hdr *fibcnet.Header, mp *FFMultipart_Request) error {
	if isSkipLog(h.level) {
		return nil
	}

	h.logger.Logf(h.level, "MP.Request: dpid  : %d", mp.DpId)
	h.logger.Logf(h.level, "MP.Request: mptype: %s", mp.MpType)
	h.logger.Logf(h.level, "MP.Request: xid   : %d", hdr.Xid)

	return nil
}

func (h *logMPHandler) FIBCFFMultipartPortRequest(hdr *fibcnet.Header, mp *FFMultipart_Request, req *FFMultipart_PortRequest) {
	LogFFMultipartPortRequest(h.logger, h.level, req)
}

func (h *logMPHandler) FIBCFFMultipartPortReply(hdr *fibcnet.Header, mp *FFMultipart_Reply, reply *FFMultipart_PortReply) {
	LogFFMultipartPortReply(h.logger, h.level, reply)
}

func (h *logMPHandler) FIBCFFMultipartReply(hdr *fibcnet.Header, mp *FFMultipart_Reply) error {
	if isSkipLog(h.level) {
		return nil
	}

	h.logger.Logf(h.level, "MP.Reply: dpid  : %d", mp.DpId)
	h.logger.Logf(h.level, "MP.Reply: mptype: %s", mp.MpType)
	h.logger.Logf(h.level, "MP.Reply: xid   : %d", hdr.Xid)

	return nil
}

func (h *logMPHandler) FIBCFFMultipartPortDescRequest(hdr *fibcnet.Header, mp *FFMultipart_Request, req *FFMultipart_PortDescRequest) {
	LogFFMultipartPortDescRequest(h.logger, h.level, req)
}

func (h *logMPHandler) FIBCFFMultipartPortDescReply(hdr *fibcnet.Header, mp *FFMultipart_Reply, reply *FFMultipart_PortDescReply) {
	LogFFMultipartPortDescReply(h.logger, h.level, reply)
}

func LogFFMultipartRequest(logger LogLogger, level log.Level, m *FFMultipart_Request, xid uint32) {
	if isSkipLog(level) {
		return
	}

	h := logMPHandler{
		level:  level,
		logger: logger,
	}

	if err := m.Dispatch(xid, &h); err != nil {
		logger.Logf(level, "MP.Request: %s", m.Body)
	}
}

func LogFFMultipartReply(logger LogLogger, level log.Level, m *FFMultipart_Reply, xid uint32) {
	if isSkipLog(level) {
		return
	}

	h := logMPHandler{
		level:  level,
		logger: logger,
	}

	if err := m.Dispatch(xid, &h); err != nil {
		logger.Logf(level, "MP.Reply: %s", m.Body)
	}
}
