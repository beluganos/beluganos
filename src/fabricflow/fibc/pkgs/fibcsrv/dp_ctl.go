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

import (
	"encoding/hex"
	fibcapi "fabricflow/fibc/api"
	"fabricflow/fibc/pkgs/fibcdbm"
	"fmt"

	log "github.com/sirupsen/logrus"
)

//
// DPCtl is dp controller.
//
type DPCtl struct {
	db *DBCtl

	stats *fibcdbm.StatsGroup
	log   *log.Entry
}

//
// NewDPCtl returns new DPCtl
//
func NewDPCtl(db *DBCtl) *DPCtl {
	return &DPCtl{
		db: db,

		stats: NewDPStats(db),
		log:   log.WithFields(log.Fields{"module": "dpctl"}),
	}
}

//
// Hello process hello message.
//
func (c *DPCtl) Hello(dpID uint64, dpType fibcapi.FFHello_DpType) error {
	c.stats.Inc(DPStatsHello)

	c.log.Debugf("SendHello: dpid:%d type:%s", dpID, dpType)
	return nil
}

//
// Monitor process monitor message.
//
func (c *DPCtl) Monitor(dpID uint64, stream fibcapi.FIBCDpApi_MonitorServer, done <-chan struct{}) error {
	c.stats.Inc(DPStatsMonitor)

	if done == nil {
		c.stats.Inc(DPStatsMonitorErr)

		c.log.Errorf("Monitor: Invalid channel")
		return fmt.Errorf("Invalid channel/")
	}

	// check dpID is registered or not.
	if _, err := c.db.ConvertIDDPtoVM(dpID); err != nil {
		c.stats.Inc(DPStatsMonitorErr)

		c.log.Warnf("Monitor: %s", err)
		return err
	}

	e := NewDPAPIMonitorEntry(stream, dpID, c.db)
	eid := e.EntryID()

	if ok := c.db.DPSet().Add(e); !ok {
		c.stats.Inc(DPStatsMonitorErr)

		c.log.Warnf("Monitor: dp already exist. dpid:%d", dpID)
		return fmt.Errorf("dp already exist. dpid:%d", dpID)
	}

	c.log.Debugf("Monitor: START. eid:%s", eid)
	e.Start(done)

	defer func() {
		c.db.DPSet().Delete(eid)
		e.Stop()
		c.leaveDP(dpID)

		c.log.Debugf("Monitor: EXIT. eid:%s", eid)
	}()

	if err := c.enterDP(dpID); err != nil {
		c.stats.Inc(DPStatsMonitorErr)

		c.log.Errorf("Monitor: enter dp error. %s", err)
		return err
	}

	<-done

	return nil
}

//
// MultipartReply process multipart reply message
//
func (c *DPCtl) MultipartReply(xid uint32, reply *fibcapi.FFMultipart_Reply) error {
	c.stats.Inc(DPStatsMP)

	c.log.Debugf("MP.Reply: dpid:%d %s xid:%d", reply.DpId, reply.MpType, xid)

	if ptdesc, ok := reply.GetBody().(*fibcapi.FFMultipart_Reply_PortDesc); ok {
		c.stats.Inc(DPStatsMPPortDesc)

		if ptdesc.PortDesc.Internal {
			// 'reply' is reply for request sent by FIBC.
			for _, port := range ptdesc.PortDesc.Port {
				c.enterPort(reply.DpId, port)
			}

			c.db.UpdateNetconfConfig(reply.DpId, ptdesc.PortDesc.Port)
		}
	}

	c.db.Waiters().Select(xid, func(w fibcdbm.Waiter) {
		fibcapi.LogFFMultipartReply(c.log, log.DebugLevel, reply, xid)
		w.Set(nil, reply)
	})

	return nil
}

func (c *DPCtl) enterDP(dpID uint64) error {
	c.stats.Inc(DPStatsEnterDP)

	c.log.Debugf("enterDP: dpid:%d", dpID)

	pdesc := NewDPMonitorReplyMpPortDescInternal(dpID)
	return c.db.SendDPMonitorReply(dpID, pdesc)
}

func (c *DPCtl) leaveDP(dpID uint64) {
	c.stats.Inc(DPStatsLeaveDP)

	c.log.Debugf("leaveDP: dpid:%d", dpID)

	leaves := []*fibcdbm.PortEntry{}
	c.db.PortMap().ListByDP(dpID, func(entry *fibcdbm.PortEntry) {
		entry.VSPort.Reset()
		entry.DPPort.UpdateEnter(false)
		leaves = append(leaves, entry.Clone())
	})

	for _, e := range leaves {
		c.stats.Inc(DPStatsLeavePort)

		c.log.Debugf("leaveDP: DP leave dpid:%d port:%d", dpID, e.DPPort.PortID)
		c.db.SendVMPortStatus(e.Key, e.VMPort.PortID, fibcapi.PortStatus_DOWN)
	}
}

func (c *DPCtl) enterPort(dpID uint64, port *fibcapi.FFPort) {
	c.stats.Inc(DPStatsEnterPort)

	c.log.Debugf("enterPort: dpid:%d port:%d", dpID, port.PortNo)

	var e *fibcdbm.PortEntry

	ok := c.db.PortMap().UpdateByDP(dpID, port.PortNo, func(entry *fibcdbm.PortEntry) {
		entry.DPPort.UpdateEnter(true)
		e = entry.Clone()
	})

	if !ok {
		c.stats.Inc(DPStatsEnterPortErr)

		c.log.Warnf("enterPort: entry not found. dpid:%d port:%d", dpID, port.PortNo)
		return
	}

	status := func() fibcapi.PortStatus_Status {
		if (port.State & fibcapi.FFPORT_STATE_LINKDOWN) != 0 {
			return fibcapi.PortStatus_DOWN
		}
		return fibcapi.PortStatus_UP
	}()

	c.log.Debugf("enterPort: DP enter dpid:%d port:%d", dpID, port.PortNo)
	c.db.SendVMPortStatusAll(e, status)
}

//
// PacketIn process packet in message.
//
func (c *DPCtl) PacketIn(dpID uint64, portID uint32, data []byte) error {
	c.stats.Inc(DPStatsPacketIn)

	if log.IsLevelEnabled(log.TraceLevel) {
		c.log.Tracef("PacketIn: dpid:%d port:%d", dpID, portID)
		c.log.Tracef("PacketIn:\n%s", hex.Dump(data))
	}

	vsID, vsPort, err := c.db.ConvertPortDPtoVS(dpID, portID)
	if err != nil {
		c.stats.Inc(DPStatsPacketInErr)

		c.log.Errorf("PacketIn: %s", err)
		return err
	}

	msg := NewVSMonitorReplyPacketOut(vsID, vsPort, data)
	if err := c.db.SendVSMonitorReply(vsID, msg); err != nil {
		c.stats.Inc(DPStatsPacketInErr)

		c.log.Errorf("PacketIn: %s", err)
		return err
	}

	return nil
}

//
// PortStatus process port status message.
//
func (c *DPCtl) PortStatus(dpID uint64, portID uint32, dpState uint32) error {
	c.stats.Inc(DPStatsPortStatus)

	c.log.Debugf("PortStatus: dpid:%d port:%d state:0x%x", dpID, portID, dpState)

	status := func() fibcapi.PortStatus_Status {
		if dpState&fibcapi.FFPORT_STATE_LINKDOWN != 0 {
			return fibcapi.PortStatus_DOWN
		}
		return fibcapi.PortStatus_UP
	}()

	vsID, vsPort, err := c.db.ConvertPortDPtoVS(dpID, portID)
	if err != nil {
		c.stats.Inc(DPStatsPortStatusErr)

		c.log.Errorf("PortStatus: %s", err)
		return err
	}

	msg := NewVSMonitorReplyPortMod(vsID, vsPort, status)
	if err := c.db.SendVSMonitorReply(vsID, msg); err != nil {
		c.stats.Inc(DPStatsPortStatusErr)

		c.log.Errorf("PortStatus: %s", err)
		return err
	}

	return nil
}

//
// L2AddrStatus process l2 addr status message
//
func (c *DPCtl) L2AddrStatus(dpID uint64, addrs []*fibcapi.L2Addr) error {
	c.stats.Inc(DPStatsL2AddrStatus)

	c.log.Debugf("L2AddrStatus: dpid:%d #addrs:%d", dpID, len(addrs))
	if log.IsLevelEnabled(log.TraceLevel) {
		for _, addr := range addrs {
			c.log.Tracef("L2AddrStatus: %s", addr)
		}
	}

	reID, vmAddrs, err := c.db.ConvertL2Addrs(dpID, addrs)
	if err != nil {
		c.stats.Inc(DPStatsL2AddrStatusErr)

		c.log.Errorf("L2AddrStatus: %s", err)
		return err
	}

	msg := NewVMMonitorReplyL2AddrStatus(reID, vmAddrs)
	if err := c.db.SendVMMonitorReply(reID, msg); err != nil {
		c.stats.Inc(DPStatsL2AddrStatusErr)

		c.log.Errorf("L2AddrStatus: %s", err)
		return err
	}

	return nil
}

func (c *DPCtl) OAMReply(xid uint32, reply *fibcapi.OAM_Reply) error {
	c.db.Waiters().Select(xid, func(w fibcdbm.Waiter) {
		fibcapi.LogOAMReply(c.log, log.DebugLevel, reply, xid)
		w.Set("dp", reply)
	})
	return nil
}
