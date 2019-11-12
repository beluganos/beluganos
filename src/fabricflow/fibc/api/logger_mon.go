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
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type logDpMonitorReplyHandler struct {
	level   log.Level
	logger  LogLogger
	dumpPkt bool
}

func (h *logDpMonitorReplyHandler) FIBCFFPacketOut(hdr *fibcnet.Header, msg *FFPacketOut) {
	LogFFPacketOut(h.logger, h.level, msg, h.dumpPkt)
}

func (h *logDpMonitorReplyHandler) FIBCFFPortMod(hdr *fibcnet.Header, msg *FFPortMod) {
	LogFFPortMod(h.logger, h.level, msg)
}

func (h *logDpMonitorReplyHandler) FIBCFlowMod(hdr *fibcnet.Header, msg *FlowMod) {
	LogFlowMod(h.logger, h.level, msg)
}

func (h *logDpMonitorReplyHandler) FIBCGroupMod(hdr *fibcnet.Header, msg *GroupMod) {
	LogGroupMod(h.logger, h.level, msg)
}

func (h *logDpMonitorReplyHandler) FIBCFFMultipartPortRequest(hdr *fibcnet.Header, mp *FFMultipart_Request, req *FFMultipart_PortRequest) {
	LogFFMultipartPortRequest(h.logger, h.level, req)
}

func (h *logDpMonitorReplyHandler) FIBCFFMultipartPortReply(hdr *fibcnet.Header, mp *FFMultipart_Reply, reply *FFMultipart_PortReply) {
	LogFFMultipartPortReply(h.logger, h.level, reply)
}

func (h *logDpMonitorReplyHandler) FIBCFFMultipartPortDescRequest(hdr *fibcnet.Header, mp *FFMultipart_Request, req *FFMultipart_PortDescRequest) {
	LogFFMultipartPortDescRequest(h.logger, h.level, req)
}

func (h *logDpMonitorReplyHandler) FIBCFFMultipartPortDescReply(hdr *fibcnet.Header, mp *FFMultipart_Reply, reply *FFMultipart_PortDescReply) {
	LogFFMultipartPortDescReply(h.logger, h.level, reply)
}

func LogDpMonitorReply(logger LogLogger, level log.Level, msg *DpMonitorReply, dumpPkt bool) {
	if isSkipLog(level) {
		return
	}

	h := logDpMonitorReplyHandler{
		level:   level,
		dumpPkt: dumpPkt,
		logger:  logger,
	}

	if err := msg.Dispatch(&h); err != nil {
		logger.Logf(level, "DpMonitorReply: %s", msg.Body)
	}
}

const (
	apAAPILogFormat = "03:04:05.000000000"
	// apAPILogFormat = "2006-01-02 03:04:05.000000000"
)

type logApMonitorReplyHandler struct {
	level  log.Level
	logger LogLogger
}

func (h *logApMonitorReplyHandler) FIBCApMonitorReplyLog(hdr *fibcnet.Header, msg *ApMonitorReplyLog) {
	fmt.Printf("%s %s\n",
		time.Unix(0, msg.Time).Format(apAAPILogFormat),
		strings.TrimRight(msg.Line, "\r\n"),
	)
}

func LogApMonitorReply(logger LogLogger, level log.Level, msg *ApMonitorReply) {
	if isSkipLog(level) {
		return
	}

	h := logApMonitorReplyHandler{
		level:  level,
		logger: logger,
	}

	if err := msg.Dispatch(&h); err != nil {
		logger.Logf(level, "ApMonitorReply: %s", msg.Body)
	}
}

type logVmMonitorReplyHandler struct {
	level  log.Level
	logger LogLogger
}

func (h *logVmMonitorReplyHandler) FIBCPortStatus(hdr *fibcnet.Header, msg *PortStatus) {
	LogPortStatus(h.logger, h.level, msg)
}

func (h *logVmMonitorReplyHandler) FIBCDpStatus(hdr *fibcnet.Header, msg *DpStatus) {
	LogDpStatus(h.logger, h.level, msg)
}

func (h *logVmMonitorReplyHandler) FIBCL2AddrStatus(hdr *fibcnet.Header, msg *L2AddrStatus) {
	LogL2AddrStatus(h.logger, h.level, msg)
}

func LogVmMonitorReply(logger LogLogger, level log.Level, msg *VmMonitorReply) {
	if isSkipLog(level) {
		return
	}

	h := logVmMonitorReplyHandler{
		level:  level,
		logger: logger,
	}

	if err := msg.Dispatch(&h); err != nil {
		logger.Logf(level, "VmMonitorReply: %s", msg.Body)
	}
}

type logVsMonitorReplyHandler struct {
	level   log.Level
	logger  LogLogger
	dumpPkt bool
}

func (h *logVsMonitorReplyHandler) FIBCFFPacketOut(hdr *fibcnet.Header, msg *FFPacketOut) {
	LogFFPacketOut(h.logger, h.level, msg, h.dumpPkt)
}

func (h *logVsMonitorReplyHandler) FIBCFFPortMod(hdr *fibcnet.Header, msg *FFPortMod) {
	LogFFPortMod(h.logger, h.level, msg)
}

func LogVsMonitorReply(logger LogLogger, level log.Level, msg *VsMonitorReply, dumpPkt bool) {
	if isSkipLog(level) {
		return
	}

	h := logVsMonitorReplyHandler{
		level:   level,
		logger:  logger,
		dumpPkt: dumpPkt,
	}

	if err := msg.Dispatch(&h); err != nil {
		logger.Logf(level, "VsMonitorReply: %s", msg.Body)
	}
}
