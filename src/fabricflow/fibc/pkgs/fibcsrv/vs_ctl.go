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
// VSCtl is vs controller
//
type VSCtl struct {
	db *DBCtl

	stats *fibcdbm.StatsGroup
	log   *log.Entry
}

//
// NewVSCtl returns new VSCtl
//
func NewVSCtl(db *DBCtl) *VSCtl {
	return &VSCtl{
		db: db,

		stats: NewVSStats(db),
		log:   log.WithFields(log.Fields{"module": "vsctl"}),
	}
}

//
// Hello process hello message.
//
func (c *VSCtl) Hello(dpID uint64, dpType fibcapi.FFHello_DpType) error {
	c.stats.Inc(VSStatsHello)

	c.log.Debugf("Hello: vsid:%d type:%s", dpID, dpType)

	if dpType != fibcapi.FFHello_FFVS {
		c.stats.Inc(VSStatsHelloErr)

		c.log.Warnf("Hello: Invalid type. %s", dpType)
		return fmt.Errorf("Invalid type. %s", dpType)
	}

	c.log.Infof("Hello: vsid:%d", dpID)
	return nil
}

//
// PacketIn process packet in message.
//
func (c *VSCtl) PacketIn(vsID uint64, portID uint32, data []byte) error {
	c.stats.Inc(VSStatsPacketIn)

	if log.IsLevelEnabled(log.TraceLevel) {
		c.log.Tracef("PacketIn: vsid:%d port:%d", vsID, portID)
		c.log.Tracef("PacketIn:\n%s", hex.Dump(data))
	}

	dpID, dpPort, err := c.db.ConvertPortVStoDP(vsID, portID)
	if err != nil {
		c.stats.Inc(VSStatsPacketInErr)

		c.log.Errorf("PacketIn: %s", err)
		return err
	}

	msg := NewDPMonitorReplyPacketOut(dpID, dpPort, data)
	if err := c.db.SendDPMonitorReply(dpID, msg); err != nil {
		c.stats.Inc(VSStatsPacketInErr)

		c.log.Errorf("PacketIn: %s", err)
		return err
	}

	return nil
}

//
// Monitor process monitor message.
//
func (c *VSCtl) Monitor(vsID uint64, stream fibcapi.FIBCVsApi_MonitorServer, done <-chan struct{}) error {
	c.stats.Inc(VSStatsMonitor)

	if done == nil {
		c.stats.Inc(VSStatsMonitorErr)

		c.log.Errorf("Monitor: Invalid channel")
		return fmt.Errorf("Invalid channel/")
	}

	e := NewVSAPIMonitorEntry(stream, vsID)
	eid := e.EntryID()

	if ok := c.db.VSSet().Add(e); !ok {
		c.stats.Inc(VSStatsMonitorErr)

		c.log.Warnf("Monitor: vs already exist. vsid:%d", vsID)
		return fmt.Errorf("vs already exist. vsid:%d", vsID)
	}

	c.log.Debugf("Monitor: START. eid:%s", eid)
	e.Start(done)

	defer func() {
		c.db.VSSet().Delete(eid)
		e.Stop()
		c.leaveVS(vsID)

		c.log.Debugf("Monitor: EXIT. eid:%s", eid)
	}()

	<-done

	return nil
}

func (c *VSCtl) leaveVS(vsID uint64) {
	c.stats.Inc(VSStatsLeaveVS)

	c.log.Debugf("leaveVS: vsid:%d", vsID)

	c.db.PortMap().ListByVS(vsID, func(entry *fibcdbm.PortEntry) {
		c.log.Debugf("leaveVS: vsid:%d port:%d", vsID, entry.VSPort.PortID)

		entry.VSPort.Reset()
	})
}

//
// FFPacket process ffpacket message.
//
func (c *VSCtl) FFPacket(vsID uint64, portID uint32, reID, ifname string) error {
	c.stats.Inc(VSStatsFFPacket)

	var (
		e   *fibcdbm.PortEntry
		upd bool
	)

	key := fibcdbm.NewPortKey(reID, ifname)
	ok := c.db.PortMap().Update(key, func(entry *fibcdbm.PortEntry) {
		if valid := entry.VSPort.IsValid(); !valid {
			return
		}

		entry.VSPort.Enter = true
		upd = entry.VSPort.Update(vsID, portID)
		e = entry.Clone()
	})

	if !ok {
		c.stats.Inc(VSStatsFFPacketErr)

		return fmt.Errorf("port entry not found. reid:%s ifname:%s", reID, ifname)
	}

	if e == nil {
		c.stats.Inc(VSStatsFFPacketErr)

		return fmt.Errorf("VS port not valid. vsid:%d port:%d", vsID, portID)
	}

	if upd {
		c.stats.Inc(VSStatsFFPacketUpd)

		c.log.Debugf("FFPacket: %s updated", e.Key)
		c.db.PortMap().RegisterVSKey(vsID, portID, key)
		c.db.SendVMPortStatusAll(e, fibcapi.PortStatus_UP)
	}

	return nil
}
