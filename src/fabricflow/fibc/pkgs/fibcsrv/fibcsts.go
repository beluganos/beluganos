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

package fibcsrv

import "fabricflow/fibc/pkgs/fibcdbm"

const (
	// APStatsMonitor is monitor message
	APStatsMonitor = "monitor"
	// APStatsMonitorErr is monitor error.
	APStatsMonitorErr = "monitor/err"
	// APStatsGetPorts is get port entry message
	APStatsGetPorts = "getportentries"
	// APStatsGetIDMap is get idmap entries message.
	APStatsGetIDMap = "getidentries"
	// APStatsGetDPSet is get dpset message
	APStatsGetDPSet = "getdpentries"
	// APStatsGetPortStats is get port stats message.
	APStatsGetPortStats = "getportstats"
	// APStatsGetPortStatsErr is get port stats error.
	APStatsGetPortStatsErr = "getportstats/err"
	// APStatsModPortStats is mod port stats message.
	APStatsModPortStats = "modportstats"
	// APStatsModPortStatsErr is mod port stats error.
	APStatsModPortStatsErr = "modportstats/err"
	// APStatsGetStats is get stats message.
	APStatsGetStats = "getstats"
	// APStatsGetStatsErr is get stats error.
	APStatsGetStatsErr = "getstats/err"
)

var apStatsNames = []string{
	APStatsMonitor,
	APStatsMonitorErr,
	APStatsGetPorts,
	APStatsGetIDMap,
	APStatsGetDPSet,
	APStatsGetPortStats,
	APStatsGetPortStatsErr,
	APStatsModPortStats,
	APStatsModPortStatsErr,
	APStatsGetStats,
	APStatsGetStatsErr,
}

//
// NewAPStats returns new StatsGroup.
//
func NewAPStats(db *DBCtl) *fibcdbm.StatsGroup {
	stats := db.Stats().Register("apctl")
	stats.RegisterList(apStatsNames)

	return stats
}

const (
	// VMStatsHello is hello message
	VMStatsHello = "hello"
	// VMStatsHelloErr is hello error.
	VMStatsHelloErr = "hello/err"
	// VMStatsFlowMod is flow mod message
	VMStatsFlowMod = "flowmod"
	// VMStatsFlowModErr is flow mod error.
	VMStatsFlowModErr = "flowmod/err"
	// VMStatsGroupMod is group mod message
	VMStatsGroupMod = "groupmod"
	// VMStatsGroupModErr is group mod error
	VMStatsGroupModErr = "groupmod/err"
	// VMStatsMonitor is monitor message
	VMStatsMonitor = "monitor"
	// VMStatsMonitorErr is monitor error.
	VMStatsMonitorErr = "monitor/err"
	// VMStatsPortConfig is port config message.
	VMStatsPortConfig = "portcfg"
	// VMStatsPortConfigErr is port config error
	VMStatsPortConfigErr = "portcfg/err"
	// VMStatsEnterPortPhy is enter port event
	VMStatsEnterPortPhy = "enter/port/phy"
	// VMStatsEnterPortPhyUpd is update port event.
	VMStatsEnterPortPhyUpd = "enter/port/phy/upd"
	// VMStatsEnterPortPhyErr is enter port error.
	VMStatsEnterPortPhyErr = "enter/port/phy/err"
	// VMStatsLeavePortPhy is leave port event.
	VMStatsLeavePortPhy = "leave/port/phy"
	// VMStatsLeavePortPhyUpd is leave port and update event
	VMStatsLeavePortPhyUpd = "leave/port/phy/upd"
	// VMStatsLeavePortPhyErr is leave port error
	VMStatsLeavePortPhyErr = "leave/port/phy/err"
	// VMStatsEnterPortVir is enter port event
	VMStatsEnterPortVir = "enter/port/vir"
	// VMStatsEnterPortVirUpd is update port event.
	VMStatsEnterPortVirUpd = "enter/port/vir/upd"
	// VMStatsEnterPortVirErr is enter port error.
	VMStatsEnterPortVirErr = "enter/port/vir/err"
	// VMStatsLeavePortVir is leave port event.
	VMStatsLeavePortVir = "leave/port/vir"
	// VMStatsLeavePortVirUpd s leave port and update event
	VMStatsLeavePortVirUpd = "leave/port/vir/upd"
	// VMStatsLeavePortVirErr is leave port error
	VMStatsLeavePortVirErr = "leave/port/vir/err"
	// VMStatsLeaveVM is leave vm event.
	VMStatsLeaveVM = "leave/vm"
)

var vmStatsNames = []string{
	VMStatsHello,
	VMStatsHelloErr,
	VMStatsFlowMod,
	VMStatsFlowModErr,
	VMStatsGroupMod,
	VMStatsGroupModErr,
	VMStatsMonitor,
	VMStatsMonitorErr,
	VMStatsPortConfig,
	VMStatsPortConfigErr,
	VMStatsEnterPortPhy,
	VMStatsEnterPortPhyUpd,
	VMStatsEnterPortPhyErr,
	VMStatsLeavePortPhy,
	VMStatsLeavePortPhyUpd,
	VMStatsLeavePortPhyErr,
	VMStatsEnterPortVir,
	VMStatsEnterPortVirUpd,
	VMStatsEnterPortVirErr,
	VMStatsLeavePortVir,
	VMStatsLeavePortVirUpd,
	VMStatsLeavePortVirErr,
	VMStatsLeaveVM,
}

//
// NewVMStats returns new StatsGroup.q
//
func NewVMStats(db *DBCtl) *fibcdbm.StatsGroup {
	stats := db.Stats().Register("vmctl")
	stats.RegisterList(vmStatsNames)

	return stats
}

const (
	// DPStatsHello is hello message.
	DPStatsHello = "hello"
	// DPStatsMonitor is monitor message.
	DPStatsMonitor = "monitor"
	// DPStatsMonitorErr is monitor error.
	DPStatsMonitorErr = "monitor/err"
	// DPStatsMP is multipart message.
	DPStatsMP = "mp"
	// DPStatsMPPortDesc is mp portdesc message
	DPStatsMPPortDesc = "mp/portdesc"
	// DPStatsEnterDP enter dp event.
	DPStatsEnterDP = "enter/dp"
	// DPStatsLeaveDP is leave dp event.
	DPStatsLeaveDP = "leave/dp"
	// DPStatsEnterPort is enter port event.
	DPStatsEnterPort = "enter/port"
	// DPStatsEnterPortErr is enter port error.
	DPStatsEnterPortErr = "enter/port/err"
	// DPStatsLeavePort is leave port event.
	DPStatsLeavePort = "leave/port"
	// DPStatsPacketIn is packet in message.
	DPStatsPacketIn = "pktin"
	// DPStatsPacketInErr is packet in error.
	DPStatsPacketInErr = "pktin/err"
	// DPStatsPortStatus is port status message
	DPStatsPortStatus = "portstatus"
	// DPStatsPortStatusErr is port status error.
	DPStatsPortStatusErr = "portstatus/err"
	// DPStatsL2AddrStatus is l2 addr status message
	DPStatsL2AddrStatus = "l2addrstatus"
	// DPStatsL2AddrStatusErr is l2 addr status error.
	DPStatsL2AddrStatusErr = "l2addrstatus/err"
)

var dpStatsNames = []string{
	DPStatsHello,
	DPStatsMonitor,
	DPStatsMonitorErr,
	DPStatsMP,
	DPStatsMPPortDesc,
	DPStatsEnterDP,
	DPStatsLeaveDP,
	DPStatsEnterPort,
	DPStatsEnterPortErr,
	DPStatsLeavePort,
	DPStatsPacketIn,
	DPStatsPacketInErr,
	DPStatsPortStatus,
	DPStatsPortStatusErr,
	DPStatsL2AddrStatus,
	DPStatsL2AddrStatusErr,
}

//
// NewDPStats returns new StatsGroup.
//
func NewDPStats(db *DBCtl) *fibcdbm.StatsGroup {
	stats := db.Stats().Register("dpctl")
	stats.RegisterList(dpStatsNames)

	return stats
}

const (
	// VSStatsHello is hello message
	VSStatsHello = "hello"
	// VSStatsHelloErr is hello error,
	VSStatsHelloErr = "hello/err"
	// VSStatsPacketIn is packet in message
	VSStatsPacketIn = "pktin"
	// VSStatsPacketInErr is packet in error.
	VSStatsPacketInErr = "pktin/err"
	// VSStatsMonitor is monitor message
	VSStatsMonitor = "monitor"
	// VSStatsMonitorErr is monitor error.
	VSStatsMonitorErr = "monitor/err"
	// VSStatsFFPacket is ffpacket message
	VSStatsFFPacket = "ffpkt"
	// VSStatsFFPacketUpd is ffpacket update
	VSStatsFFPacketUpd = "ffpkt/upd"
	// VSStatsFFPacketErr is ffpacket error
	VSStatsFFPacketErr = "ffpkt/err"
	// VSStatsLeaveVS is leave vs event.
	VSStatsLeaveVS = "leave/vs"
)

var vsStatsNames = []string{
	VSStatsHello,
	VSStatsHelloErr,
	VSStatsPacketIn,
	VSStatsPacketInErr,
	VSStatsMonitor,
	VSStatsMonitorErr,
	VSStatsFFPacket,
	VSStatsFFPacketUpd,
	VSStatsFFPacketErr,
	VSStatsLeaveVS,
}

//
// NewVSStats returns new StatsGroup.
//
func NewVSStats(db *DBCtl) *fibcdbm.StatsGroup {
	stats := db.Stats().Register("vsctl")
	stats.RegisterList(vsStatsNames)

	return stats
}
