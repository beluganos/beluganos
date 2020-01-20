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

	log "github.com/sirupsen/logrus"
)

//
// VMCtl is vm controller
//
type VMCtl struct {
	db *DBCtl

	stats *fibcdbm.StatsGroup
	log   *log.Entry
}

//
// NewVMCtl returns new VMCtl
//
func NewVMCtl(db *DBCtl) *VMCtl {
	return &VMCtl{
		db: db,

		stats: NewVMStats(db),
		log:   log.WithFields(log.Fields{"module": "vmctl"}),
	}
}

//
// Hello process hello message.
//
func (c *VMCtl) Hello(hello *fibcapi.Hello) error {
	c.stats.Inc(VMStatsHello)

	fibcapi.LogHello(c.log, log.DebugLevel, hello)

	if _, err := c.db.ConvertIDVMtoDP(hello.ReId); err != nil {
		c.stats.Inc(VMStatsHelloErr)

		c.log.Warnf("Hello: %s", err)
		return err
	}

	return nil
}

func (c *VMCtl) sendDPPortMod(entry *fibcdbm.PortEntry, status fibcapi.PortStatus_Status) {
	msg := NewDPMonitorReplyPortMod(
		entry.DPPort.DpID,
		entry.DPPort.PortID,
		status,
	)

	if err := c.db.SendDPMonitorReply(entry.DPPort.DpID, msg); err != nil {
		c.log.Errorf("PortConfig: %s", err)
	}
}

func (c *VMCtl) enterPortPhy(entry *fibcdbm.PortEntry, status fibcapi.PortStatus_Status) error {
	c.stats.Inc(VMStatsEnterPortPhy)

	c.log.Debugf("enterPort(phy): %s", entry.Key)

	var upd bool

	ok := c.db.PortMap().Update(entry.Key, func(e *fibcdbm.PortEntry) {
		upd = e.VMPort.UpdatePort(entry.VMPort.PortID)
		e.VMPort.UpdateEnter(true)
		e.MasterKey = entry.MasterKey
		entry = e.Clone()
	})

	if !ok {
		c.stats.Inc(VMStatsEnterPortPhyErr)

		c.log.Errorf("enterPort(phy): port not found. %s", entry.Key)
		return fmt.Errorf("port not found. %s", entry.Key)
	}

	if upd {
		c.stats.Inc(VMStatsEnterPortPhyUpd)

		c.log.Debugf("enterPort(phy): %s updated", entry.Key)
		c.db.PortMap().RegisterVMKey(entry.VMPort.ReID, entry.VMPort.PortID, entry.Key)
		c.db.SendVMPortStatusAll(entry, fibcapi.PortStatus_UP)
	}

	c.sendDPPortMod(entry, status)

	return nil
}

func (c *VMCtl) leavePortPhy(key *fibcdbm.PortKey, status fibcapi.PortStatus_Status) error {
	c.stats.Inc(VMStatsLeavePortPhy)

	c.log.Debugf("leavePort(phy): %s", key)

	var (
		upd   bool
		entry *fibcdbm.PortEntry
	)

	ok := c.db.PortMap().Update(key, func(e *fibcdbm.PortEntry) {
		entry = e.Clone()
		e.MasterKey = nil
		e.VSPort.Reset()
		e.VMPort.UpdateEnter(false)
		upd = e.VMPort.UpdatePort(0)
	})

	if !ok {
		c.stats.Inc(VMStatsLeavePortPhyErr)

		c.log.Errorf("leavePort(phy): port not found. %s", key)
		return fmt.Errorf("port not found. %s", key)
	}

	if upd {
		c.stats.Inc(VMStatsLeavePortPhyUpd)

		c.log.Debugf("leavePort(phy): %s updated", key)
		c.db.PortMap().UnregisterVMKey(entry.VMPort.ReID, entry.VMPort.PortID)
		c.db.SendVMPortStatus(entry.Key, entry.VMPort.PortID, fibcapi.PortStatus_DOWN)
	}

	c.sendDPPortMod(entry, status)

	return nil
}

func (c *VMCtl) enterPortVirt(entry *fibcdbm.PortEntry) error {
	c.stats.Inc(VMStatsEnterPortVir)

	c.log.Debugf("leavePort(vir): %s", entry.Key)

	var upd bool

	newEntry := fibcdbm.PortEntry{
		Key:    entry.Key,
		VMPort: &fibcdbm.PortValue{},
	}
	c.db.PortMap().SelectOrRegister(&newEntry, func(e *fibcdbm.PortEntry) bool {
		upd = e.VMPort.UpdateR(entry.Key.ReID, entry.VMPort.PortID)
		e.VMPort.UpdateEnter(true)
		e.DPPort = entry.DPPort
		e.ParentKey = entry.ParentKey
		e.MasterKey = entry.MasterKey
		entry = e.Clone()
		return true
	})

	if upd {
		c.stats.Inc(VMStatsEnterPortVirUpd)

		c.log.Debugf("enterPort(vir): %s updated", entry.Key)
		c.db.PortMap().RegisterVMKey(entry.VMPort.ReID, entry.VMPort.PortID, entry.Key)
		c.db.SendVMPortStatus(entry.Key, entry.VMPort.PortID, fibcapi.PortStatus_UP)
	}

	return nil
}

func (c *VMCtl) leavePortVirt(key *fibcdbm.PortKey) error {
	c.stats.Inc(VMStatsLeavePortVir)

	c.log.Debugf("leavePort(vir): %s", key)

	e := c.db.PortMap().Unregister(key)

	if e == nil {
		c.stats.Inc(VMStatsLeavePortVirErr)

		c.log.Errorf("leavePort(vir): port not found. %s", key)
		return fmt.Errorf("port not found. %s", key)
	}

	c.log.Debugf("leavePort(vir): %s updated", key)

	// c.db.Ports().UnregisterVmKey(e.VmPort.ReId, e.VmPort.PortId) // done by Unregister()
	c.db.SendVMPortStatus(e.Key, e.VMPort.PortID, fibcapi.PortStatus_DOWN)

	return nil
}

//
// PortConfig process port config message.
//
func (c *VMCtl) PortConfig(pc *fibcapi.PortConfig) error {
	c.stats.Inc(VMStatsPortConfig)

	fibcapi.LogPortConfig(c.log, log.DebugLevel, pc)

	entry, err := c.db.NewPortEntryFromPortConfig(pc)
	if err != nil {
		c.stats.Inc(VMStatsPortConfigErr)

		c.log.Errorf("PortConfg: Bad message. %s", err)
		return err
	}

	if len(pc.Link) == 0 && pc.DpPort == 0 {
		// physical device
		switch pc.Cmd {
		case fibcapi.PortConfig_ADD, fibcapi.PortConfig_MODIFY:
			return c.enterPortPhy(entry, pc.Status)

		case fibcapi.PortConfig_DELETE:
			return c.leavePortPhy(entry.Key, pc.Status)

		default:
			c.stats.Inc(VMStatsPortConfigErr)

			c.log.Errorf("PortConfig: Invalid cmd. %s", pc.Cmd)
			return fmt.Errorf("Invalid cmd. %s", pc.Cmd)
		}

	} else {
		// virtual device(vlan, iptun, bridge...)
		switch pc.Cmd {
		case fibcapi.PortConfig_ADD, fibcapi.PortConfig_MODIFY:
			return c.enterPortVirt(entry)

		case fibcapi.PortConfig_DELETE:
			return c.leavePortVirt(entry.Key)

		default:
			c.stats.Inc(VMStatsPortConfigErr)

			c.log.Errorf("PortConfig: Invalid cmd. %s", pc.Cmd)
			return fmt.Errorf("Invalid cmd. %s", pc.Cmd)
		}
	}
}

func (c *VMCtl) leaveVM(reID string) {
	c.stats.Inc(VMStatsLeaveVM)

	c.log.Debugf("leaveVM: reid:'%s'", reID)

	delKeys := []*fibcdbm.PortKey{}
	c.db.PortMap().ListByVM(reID, func(e *fibcdbm.PortEntry) {
		if e.VSPort == nil {
			delKeys = append(delKeys, e.Key.Clone())
		} else {
			e.VMPort.UpdateEnter(false)
			e.VMPort.UpdatePort(0)
			e.MasterKey = nil

			c.log.Debugf("leaveVM: reset. %s", e.Key)
		}
	})

	for _, key := range delKeys {
		c.db.PortMap().Unregister(key)

		// c.db.Ports().UnregisterVmKey(e.VmPort.ReId, e.VmPort.PortId) // done by Unregister()
		c.log.Debugf("leaveVM: vm port deleted. %s", key)
	}
}

//
// Monitor process monitor message.
//
func (c *VMCtl) Monitor(reID string, stream fibcapi.FIBCVmApi_MonitorServer, done <-chan struct{}) error {
	c.stats.Inc(VMStatsMonitor)

	if done == nil {
		c.stats.Inc(VMStatsMonitorErr)

		c.log.Errorf("Monitor: Invalid channel")
		return fmt.Errorf("Invalid channel/")
	}

	if _, err := c.db.ConvertIDVMtoDP(reID); err != nil {
		c.stats.Inc(VMStatsMonitorErr)

		c.log.Warnf("Monitor: %s", err)
		return err
	}

	e := NewVMAPIMonitorEntry(stream, reID)
	eid := e.EntryID()

	if ok := c.db.VMSet().Add(e); !ok {
		c.stats.Inc(VMStatsMonitorErr)

		c.log.Warnf("Monitor: vm already exist. reid:'%s'", reID)
		return fmt.Errorf("vm already exist. reid:'%s'", reID)
	}

	c.log.Debugf("Monitor: START. eid:%s", eid)
	e.Start(done)

	defer func() {
		c.db.VMSet().Delete(eid)
		e.Stop()
		c.leaveVM(reID)

		c.log.Debugf("Monitor: EXIT. eid:%s", eid)
	}()

	<-done

	return nil
}

//
// FlowMod process flow mod message.
//
func (c *VMCtl) FlowMod(mod *fibcapi.FlowMod) error {
	c.stats.Inc(VMStatsFlowMod)

	fibcapi.LogFlowMod(c.log, log.DebugLevel, mod)

	dpID, err := c.db.ConvertIDVMtoDP(mod.ReId)
	if err != nil {
		c.stats.Inc(VMStatsFlowModErr)

		c.log.Errorf("FlowMod: convert error. %s", err)
		return err
	}

	msg := NewDPMonitorReplyFlowMod(mod)
	if err := c.db.SendDPMonitorMod(dpID, msg); err != nil {
		c.stats.Inc(VMStatsFlowModErr)

		c.log.Errorf("FlowMod: send error. %s", err)
		return err
	}

	return nil
}

//
// GroupMod process group mod message.
//
func (c *VMCtl) GroupMod(mod *fibcapi.GroupMod) error {
	c.stats.Inc(VMStatsGroupMod)

	fibcapi.LogGroupMod(c.log, log.DebugLevel, mod)

	dpID, err := c.db.ConvertIDVMtoDP(mod.ReId)
	if err != nil {
		c.stats.Inc(VMStatsGroupModErr)

		c.log.Errorf("GroupMod: convert error. %s", err)
		return err
	}

	msg := NewDPMonitorReplyGroupMod(mod)
	if err := c.db.SendDPMonitorMod(dpID, msg); err != nil {
		c.stats.Inc(VMStatsGroupModErr)

		c.log.Errorf("GroupMod: send error. %s", err)
		return err
	}

	return nil
}

func (c *VMCtl) OAMReply(xid uint32, reply *fibcapi.OAM_Reply) error {
	c.db.Waiters().Select(xid, func(w fibcdbm.Waiter) {
		fibcapi.LogOAMReply(c.log, log.DebugLevel, reply, xid)
		w.Set(reply.ReId, reply)
	})
	return nil
}
