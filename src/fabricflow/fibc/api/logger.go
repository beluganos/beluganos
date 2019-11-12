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
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

type LogLogger interface {
	Logf(log.Level, string, ...interface{})
}

const (
	LogPacketDumpSize = 256
)

func isSkipLog(level log.Level) bool {
	return log.GetLevel() < level
}

func LogHello(logger LogLogger, level log.Level, hello *Hello) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "Hello: reid: '%s'", hello.ReId)
}

func LogDpStatus(logger LogLogger, level log.Level, m *DpStatus) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "DpStatus: reid  : '%s'", m.ReId)
	logger.Logf(level, "DpStatus: status: %s", m.Status)
}

func LogPortStatus(logger LogLogger, level log.Level, m *PortStatus) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "PortStatus: reid  : '%s'", m.ReId)
	logger.Logf(level, "PortStatus: port  : %d", m.PortId)
	logger.Logf(level, "PortStatus: ifname: '%s'", m.Ifname)
	logger.Logf(level, "PortStatus: status: %s", m.Status)
}

func LogPortConfig(logger LogLogger, level log.Level, m *PortConfig) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "PortConfig: cmd   : %s", m.Cmd)
	logger.Logf(level, "PortConfig: reid  : '%s'", m.ReId)
	logger.Logf(level, "PortConfig: port  : %d", m.PortId)
	logger.Logf(level, "PortConfig: ifname: '%s'", m.Ifname)
	logger.Logf(level, "PortConfig: link  : '%s'", m.Link)
	logger.Logf(level, "PortConfig: master: '%s'", m.Master)
	logger.Logf(level, "PortConfig: dpport: %d", m.DpPort)
	logger.Logf(level, "PortConfig: status: %s", m.Status)
}

func LogFFHello(logger LogLogger, level log.Level, m *FFHello) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFHello: dpid: %d", m.DpId)
	logger.Logf(level, "FFHello: type: %s", m.DpType)
}

func LogFFPort(logger LogLogger, level log.Level, m *FFPort) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFPort: port  : %d", m.PortNo)
	logger.Logf(level, "FFPort: hwaddr: '%s'", m.HwAddr)
	logger.Logf(level, "FFPort: name  : '%s'", m.Name)
	logger.Logf(level, "FFPort: config: '%s'", m.Config)
	logger.Logf(level, "FFPort: state : %d", m.State)
	logger.Logf(level, "FFPort: curr  : %d", m.Curr)
	logger.Logf(level, "FFPort: adv   : %d", m.Advertised)
	logger.Logf(level, "FFPort: cur_spd: %d", m.CurrSpeed)
	logger.Logf(level, "FFPort: mac_spd: %d", m.MaxSpeed)
}

func LogFFPortStats(logger LogLogger, level log.Level, m *FFPortStats) {
	if isSkipLog(level) {
		return
	}

	if values := m.Values; values != nil {
		for k, v := range values {
			logger.Logf(level, "FFPortStats: '%s' = %d", k, v)
		}
	}
	if values := m.SValues; values != nil {
		for k, v := range values {
			logger.Logf(level, "FFPortStats: '%s' = '%s'", k, v)
		}
	}
}

func LogFFPacketIn(logger LogLogger, level log.Level, m *FFPacketIn, dumpData bool) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFPacketIn: dpid: %d", m.DpId)
	logger.Logf(level, "FFPacketIn: port: %d", m.PortNo)

	if dumpData {
		dumpSize := len(m.Data)
		if dumpSize > LogPacketDumpSize {
			dumpSize = LogPacketDumpSize
		}

		logger.Logf(level, "FFPacketIn:\n%s", hex.Dump(m.Data[:dumpSize]))
	}
}

func LogFFPacketOut(logger LogLogger, level log.Level, m *FFPacketOut, dumpData bool) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFPacketOut: dpid: %d", m.DpId)
	logger.Logf(level, "FFPacketOut: port: %d", m.PortNo)

	if dumpData {
		dumpSize := len(m.Data)
		if dumpSize > LogPacketDumpSize {
			dumpSize = LogPacketDumpSize
		}

		logger.Logf(level, "FFPacketOut:\n%s", hex.Dump(m.Data[:dumpSize]))
	}
}

func LogFFPacket(logger LogLogger, level log.Level, m *FFPacket) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFPacket: dpid  : %d", m.DpId)
	logger.Logf(level, "FFPacket: port  : %d", m.PortNo)
	logger.Logf(level, "FFPacket: reid  : '%s'", m.ReId)
	logger.Logf(level, "FFPacket: ifname: '%s'", m.Ifname)
}

func LogFFPortStatus(logger LogLogger, level log.Level, m *FFPortStatus) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFPortStatus: dpid  : %d", m.DpId)
	logger.Logf(level, "FFPortStatus: reason: %d", m.Reason)
	LogFFPort(logger, level, m.Port)
}

func LogFFPortMod(logger LogLogger, level log.Level, m *FFPortMod) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFPortMod: dpid  : %d", m.DpId)
	logger.Logf(level, "FFPortMod: port  : %d", m.PortNo)
	logger.Logf(level, "FFPortMod: mac   : '%s'", m.HwAddr)
	logger.Logf(level, "FFPortMod: status: %s", m.Status)
}

func LogFFL2AddrStatus(logger LogLogger, level log.Level, m *FFL2AddrStatus) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "FFL2AddrStatus: dpid: %d", m.DpId)
	if addrs := m.Addrs; addrs != nil {
		for _, addr := range addrs {
			LogL2Addr(logger, level, addr)
		}
	}
}

func LogL2AddrStatus(logger LogLogger, level log.Level, m *L2AddrStatus) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "L2AddrStatus: reid: '%s'", m.ReId)
	if addrs := m.Addrs; addrs != nil {
		for _, addr := range addrs {
			LogL2Addr(logger, level, addr)
		}
	}
}

func LogL2Addr(logger LogLogger, level log.Level, m *L2Addr) {
	if isSkipLog(level) {
		return
	}

	logger.Logf(level, "L2Addr: mac   : '%s'", m.HwAddr)
	logger.Logf(level, "L2Addr: port  : %d", m.PortId)
	logger.Logf(level, "L2Addr: vid   : %d", m.VlanVid)
	logger.Logf(level, "L2Addr: ifname: '%s'", m.Ifname)
	logger.Logf(level, "L2Addr: reason: %s", m.Reason)
}
