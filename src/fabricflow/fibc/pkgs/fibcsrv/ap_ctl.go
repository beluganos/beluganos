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
	fibcapi "fabricflow/fibc/api"
	"fabricflow/fibc/pkgs/fibcdbm"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	apAPIPortStatWaitMax = 3 * time.Second
)

var apCtlPortStatsNames = []string{
	"ifInOctets",
	"ifInUcastPkts",
	"ifInNUcastPkts",
	"ifInDiscards",
	"ifInErrors",
	"ifOutOctets",
	"ifOutUcastPkts",
	"ifOutNUcastPkts",
	"ifOutDiscards",
	"ifOutErrors",
}

//
// APCtl is ap api controller
//
type APCtl struct {
	db        *DBCtl
	psTimeout time.Duration

	stats *fibcdbm.StatsGroup
	log   *log.Entry
}

//
// NewAPCtl returns new APCtl.
//
func NewAPCtl(db *DBCtl) *APCtl {
	return &APCtl{
		db:        db,
		psTimeout: apAPIPortStatWaitMax,

		stats: NewAPStats(db),
		log:   log.WithFields(log.Fields{"module": "apctl"}),
	}
}

//
// Monitor process monitor request.
//
func (c *APCtl) Monitor(stream fibcapi.FIBCApApi_MonitorServer, done <-chan struct{}) error {
	c.stats.Inc(APStatsMonitor)

	if done == nil {
		c.stats.Inc(APStatsMonitorErr)

		c.log.Errorf("Monitor: Invalid channel")
		return fmt.Errorf("Invalid channel/")
	}

	e := NewAPAPIMonitorEntry(stream)
	eid := e.EntryID()

	if ok := c.db.APSet().Add(e); !ok {
		c.stats.Inc(APStatsMonitorErr)

		c.log.Warnf("Monitor: ap already exist. eid:%s", eid)
		return fmt.Errorf("ap already exist. eid:%s", eid)
	}

	c.log.Debugf("Monitor: START. eid:%s", eid)
	e.Start(done)

	defer func() {
		c.db.APSet().Delete(eid)
		e.Stop()

		c.log.Debugf("Monitor: EXIT. eid:%s", eid)
	}()

	<-done

	return nil
}

//
// GetPortEntries process get port entries request.
//
func (c *APCtl) GetPortEntries(stream fibcapi.FIBCApApi_GetPortEntriesServer) error {
	c.stats.Inc(APStatsGetPorts)

	c.db.PortMap().Range(func(e *fibcdbm.PortEntry) {
		msg := NewDBPortEntryFromLocal(e)
		if err := stream.Send(msg); err != nil {
			c.log.Errorf("GetPortEntries: send error. %s", err)
		}
	})

	return nil
}

//
// GetIDEntries process get id entries request.
//
func (c *APCtl) GetIDEntries(stream fibcapi.FIBCApApi_GetIDEntriesServer) error {
	c.stats.Inc(APStatsGetIDMap)

	c.db.IDMap().Range(func(e *fibcdbm.IDEntry) {
		msg := NewDBIDEntryFromLocal(e)
		if err := stream.Send(msg); err != nil {
			c.log.Errorf("GetIDEntries: send error. %s", err)
		}
	})

	return nil
}

func (c *APCtl) getApMonitorEntries(stream fibcapi.FIBCApApi_GetDpEntriesServer) {
	c.db.APSet().Range(func(e fibcdbm.DPEntry) {
		entry := e.(*APAPIMonitorEntry)
		msg := fibcapi.DbDpEntry{
			Type: fibcapi.DbDpEntry_APMON,
			Id:   entry.EntryID(),
		}

		if err := stream.Send(&msg); err != nil {
			c.log.Errorf("GetDpEntries: send ap entry error. %s", err)
		}
	})
}

func (c *APCtl) getVMMonitorEntries(stream fibcapi.FIBCApApi_GetDpEntriesServer) {
	c.db.VMSet().Range(func(e fibcdbm.DPEntry) {
		entry := e.(*VMAPIMonitorEntry)
		msg := fibcapi.DbDpEntry{
			Type: fibcapi.DbDpEntry_VMMON,
			Id:   entry.reID,
		}

		if err := stream.Send(&msg); err != nil {
			c.log.Errorf("GetDpEntries: send vm entry error. %s", err)
		}
	})
}

func (c *APCtl) getDpMonitorEntries(stream fibcapi.FIBCApApi_GetDpEntriesServer) {
	c.db.DPSet().Range(func(e fibcdbm.DPEntry) {
		entry := e.(*DPAPIMonitorEntry)
		msg := fibcapi.DbDpEntry{
			Type: fibcapi.DbDpEntry_DPMON,
			Id:   fmt.Sprintf("%d", entry.dpID),
		}

		if err := stream.Send(&msg); err != nil {
			c.log.Errorf("GetDpEntries: send dp entry error. %s", err)
		}
	})
}

func (c *APCtl) getVsMonitorEntries(stream fibcapi.FIBCApApi_GetDpEntriesServer) {
	c.db.VSSet().Range(func(e fibcdbm.DPEntry) {
		entry := e.(*VSAPIMonitorEntry)
		msg := fibcapi.DbDpEntry{
			Type: fibcapi.DbDpEntry_VSMON,
			Id:   fmt.Sprintf("%d", entry.vsID),
		}

		if err := stream.Send(&msg); err != nil {
			c.log.Errorf("GetDpEntries: send vs entry error. %s", err)
		}
	})
}

//
// GetDpEntries process get dp entries request.
//
func (c *APCtl) GetDpEntries(t fibcapi.DbDpEntry_Type, stream fibcapi.FIBCApApi_GetDpEntriesServer) error {
	c.stats.Inc(APStatsGetDPSet)

	if t == fibcapi.DbDpEntry_NOP || t == fibcapi.DbDpEntry_APMON {
		c.getApMonitorEntries(stream)
	}

	if t == fibcapi.DbDpEntry_NOP || t == fibcapi.DbDpEntry_VMMON {
		c.getVMMonitorEntries(stream)
	}

	if t == fibcapi.DbDpEntry_NOP || t == fibcapi.DbDpEntry_DPMON {
		c.getDpMonitorEntries(stream)
	}

	if t == fibcapi.DbDpEntry_NOP || t == fibcapi.DbDpEntry_VSMON {
		c.getVsMonitorEntries(stream)
	}

	return nil
}

//
// AddPortEntry process add port entry request.
//
func (c *APCtl) AddPortEntry(entry *fibcapi.DbPortEntry) error {
	e := NewDBPortEntryFromAPI(entry)
	e.VMPort = &fibcdbm.PortValue{}
	e.VSPort = &fibcdbm.PortValue{}

	c.db.PortMap().Register(e)
	return nil
}

//
// DelPortEntry process del port entry request.
//
func (c *APCtl) DelPortEntry(key *fibcapi.DbPortKey) error {
	k := NewDBPortKeyFromAPI(key)
	c.db.PortMap().Unregister(k)
	return nil
}

//
// AddIDEntry process add id entry request.
//
func (c *APCtl) AddIDEntry(entry *fibcapi.DbIdEntry) error {
	e := NewDBIDEntryFromAPI(entry)
	return c.db.IDMap().Register(e)
}

//
// DelIDEntry process del id entry request.
//
func (c *APCtl) DelIDEntry(entry *fibcapi.DbIdEntry) error {
	if len(entry.ReId) != 0 {
		c.db.IDMap().UnregisterByReID(entry.ReId)
		return nil
	}

	if entry.DpId != 0 {
		c.db.IDMap().UnregisterByDpID(entry.DpId)
		return nil
	}

	return fmt.Errorf("Invalid entry. %s", entry)
}

//
// GetPortStats process get port stats.
//
func (c *APCtl) GetPortStats(dpID uint64, portID uint32, names []string, stream fibcapi.FIBCApApi_GetPortStatsServer) error {
	c.stats.Inc(APStatsGetPortStats)

	c.log.Debugf("PortStats: dpid:%d port:%d names:%v", dpID, portID, names)

	if len(names) == 0 {
		names = apCtlPortStatsNames

		c.log.Debugf("PortStats: use default stats list.")
	}

	w := NewDBMpWaiter()
	xid := c.db.Waiters().Register(w)
	defer c.db.Waiters().Unregister(xid)

	msg := NewDPMonitorReplyMpPort(dpID, portID, names, xid)
	if err := c.db.SendDPMonitorReply(dpID, msg); err != nil {
		c.stats.Inc(APStatsGetPortStatsErr)

		c.log.Errorf("GetPortStats: send monitor reply error. %s", err)
		return err
	}

	c.log.Debugf("PortStats: wait... xid:%d", xid)

	if err := w.Wait(c.psTimeout); err != nil {
		c.stats.Inc(APStatsGetPortStatsErr)

		c.log.Errorf("GetPortStats: xid:%d %s", xid, err)
		return err
	}

	c.log.Debugf("PortStats: wait...done. xid:%d", xid)

	portReply := w.Reply.GetPort()
	if portReply == nil {
		c.stats.Inc(APStatsGetPortStatsErr)

		c.log.Errorf("GetPortStats: invaid reply from dp. %d", xid)
		return fmt.Errorf("invalid reply from dp. %d", xid)
	}

	c.extendPortStats(dpID, portReply.Stats)

	for _, stats := range portReply.Stats {
		if err := stream.Send(stats); err != nil {
			c.stats.Inc(APStatsGetPortStatsErr)

			c.log.Errorf("GetPortStats: send error. %s", err)
			return err
		}
	}

	return nil
}

func (c *APCtl) extendPortStats(dpID uint64, portStatsList []*fibcapi.FFPortStats) {
	boolToIfOperStatus := func(b bool) uint64 {
		// 1.3.6.1.2.1.2.2.1.8
		// 1-up, 2-down
		if b {
			return 1
		}
		return 2
	}

	for _, portStats := range portStatsList {
		portID := portStats.PortNo
		portStats.Values["port_no"] = uint64(portID)
		if portStats.SValues == nil {
			portStats.SValues = map[string]string{}
		}

		portExist := c.db.PortMap().SelectByDP(dpID, portID, func(e *fibcdbm.PortEntry) {
			portStats.SValues["ifName"] = e.Key.Ifname
			portStats.Values["ifOperStatus"] = boolToIfOperStatus(e.DPPort.Enter)
		})

		if !portExist {
			portStats.SValues["ifName"] = ""
			portStats.Values["ifOperStatus"] = boolToIfOperStatus(false)
		}
	}
}

//
// GetStats process get stats request.
//
func (c *APCtl) GetStats(stream fibcapi.FIBCApApi_GetStatsServer) error {
	c.stats.Inc(APStatsGetStats)

	c.db.Stats().Range(func(group, name string, value uint64) {
		msg := fibcapi.StatsEntry{
			Group: group,
			Name:  name,
			Value: value,
		}

		if err := stream.Send(&msg); err != nil {
			c.stats.Inc(APStatsGetStatsErr)

			c.log.Errorf("GetStats: send ap entry error. %s", err)
		}
	})

	return nil
}
