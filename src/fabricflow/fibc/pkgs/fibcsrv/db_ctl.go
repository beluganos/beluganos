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
// DBMpWaiter is waiter for multipart message,
//
type DBMpWaiter struct {
	*fibcdbm.SimpleWaiter
	Reply *fibcapi.FFMultipart_Reply
}

//
// Set sets reply and close wait channel.
//
func (w *DBMpWaiter) Set(v interface{}) {
	if reply, ok := v.(*fibcapi.FFMultipart_Reply); ok {
		w.Reply = reply
		w.SimpleWaiter.Close()
		return
	}

	w.SimpleWaiter.SetError(fmt.Errorf("Invalid reply. %v", v))
}

//
// NewDBMpWaiter returns new DBMpWaiter
//
func NewDBMpWaiter() *DBMpWaiter {
	return &DBMpWaiter{
		SimpleWaiter: fibcdbm.NewSimpleWaiter(),
	}
}

//
// DBCtl is db controller.
//
type DBCtl struct {
	apset *fibcdbm.DPSet
	vmset *fibcdbm.DPSet
	vsset *fibcdbm.DPSet
	dpset *fibcdbm.DPSet
	idmap *fibcdbm.IDMap
	ptmap *fibcdbm.PortMap
	stats *fibcdbm.StatsTable

	waits *fibcdbm.WaiterTable

	nccfg *fibcdbm.NetconfConfig

	log *log.Entry
}

//
// NewDBCtl returns new DBCtl
//
func NewDBCtl() *DBCtl {
	return &DBCtl{
		apset: fibcdbm.NewDPSet(),
		vmset: fibcdbm.NewDPSet(),
		vsset: fibcdbm.NewDPSet(),
		dpset: fibcdbm.NewDPSet(),
		idmap: fibcdbm.NewIDMap(),
		ptmap: fibcdbm.NewPortMap(),
		stats: fibcdbm.NewStatsTable(),
		waits: fibcdbm.NewWaiterTable(),
		nccfg: fibcdbm.NewNetconfConfig(),

		log: log.WithFields(log.Fields{"module": "dbctl"}),
	}
}

//
// APSet returns APSet table
//
func (c *DBCtl) APSet() *fibcdbm.DPSet {
	return c.apset
}

//
// VMSet returns VMSet table.
//
func (c *DBCtl) VMSet() *fibcdbm.DPSet {
	return c.vmset
}

//
// VSSet returns VSSet table/
//
func (c *DBCtl) VSSet() *fibcdbm.DPSet {
	return c.vsset
}

//
// DPSet returns DPSet table.
//
func (c *DBCtl) DPSet() *fibcdbm.DPSet {
	return c.dpset
}

//
// IDMap returns idmap table.
//
func (c *DBCtl) IDMap() *fibcdbm.IDMap {
	return c.idmap
}

//
// PortMap returns portmap table.
//
func (c *DBCtl) PortMap() *fibcdbm.PortMap {
	return c.ptmap
}

//
// Stats returns stats table.
//
func (c *DBCtl) Stats() *fibcdbm.StatsTable {
	return c.stats
}

//
// Waiters returns waiter table.
//
func (c *DBCtl) Waiters() *fibcdbm.WaiterTable {
	return c.waits
}

//
// NetconfConfig returns netconf config table.
//
func (c *DBCtl) NetconfConfig() *fibcdbm.NetconfConfig {
	return c.nccfg
}

//
// SendVMPortStatus sends PortStats message.
//
func (c *DBCtl) SendVMPortStatus(key *fibcdbm.PortKey, portID uint32, status fibcapi.PortStatus_Status) error {
	if status != fibcapi.PortStatus_DOWN {
		if ok, _ := c.PortMap().IsAssociated(key); !ok {
			c.log.Debugf("PortStatus: not associated. %s", key)
			return nil
		}
	}

	c.log.Infof("PortStatus: %s port:%d %s ", key, portID, status)

	msg := NewVMMonitorReplyPortStatus(key.ReID, portID, key.Ifname, status)
	return c.SendVMMonitorReply(key.ReID, msg)
}

//
// SendVMPortStatusAll sends PortStats message.
//
func (c *DBCtl) SendVMPortStatusAll(entry *fibcdbm.PortEntry, status fibcapi.PortStatus_Status) {
	entries := []*fibcdbm.PortEntry{entry}

	c.PortMap().ListByParent(entry.Key, func(e *fibcdbm.PortEntry) bool {
		entries = append(entries, e.Clone())
		return true
	})

	for _, e := range entries {
		if err := c.SendVMPortStatus(e.Key, e.VMPort.PortID, status); err != nil {
			c.log.Warnf("send port status error. %s", err)
		}
	}
}

//
// SendAPMonitorReply send ap monitor reply.
//
func (c *DBCtl) SendAPMonitorReply(msg *fibcapi.ApMonitorReply) {
	// NOTE: DO NOT USE logging.
	c.APSet().Range(func(e fibcdbm.DPEntry) {
		mon := e.(*APAPIMonitorEntry)
		mon.Send(msg)

	})
}

//
// SendAPMonitorReplyLog send ap monitor reply.
//
func (c *DBCtl) SendAPMonitorReplyLog(msg *fibcapi.ApMonitorReplyLog) {
	// NOTE: DO NOT USE logging.
	c.APSet().Range(func(e fibcdbm.DPEntry) {
		mon := e.(*APAPIMonitorEntry)
		mon.SendLog(msg)

	})
}

//
// SendVMMonitorReply send vm monitor reply.
//
func (c *DBCtl) SendVMMonitorReply(reID string, msg *fibcapi.VmMonitorReply) error {
	eid := NewVMAPIMonitorEntryID(reID)
	ok := c.VMSet().Select(eid, func(e fibcdbm.DPEntry) {
		mon := e.(*VMAPIMonitorEntry)
		mon.Send(msg)
	})

	if !ok {
		return fmt.Errorf("vm entry not found. reid:%s", eid)
	}

	return nil
}

//
// SendDPMonitorReply send dp monitor reply.
//
func (c *DBCtl) SendDPMonitorReply(dpID uint64, msg *fibcapi.DpMonitorReply) error {
	eid := NewDPAPIMonitorEntryID(dpID)
	ok := c.DPSet().Select(eid, func(e fibcdbm.DPEntry) {
		mon := e.(*DPAPIMonitorEntry)
		mon.Send(msg)
	})

	if !ok {
		return fmt.Errorf("dp entry not found. dpid=%d", dpID)
	}

	return nil
}

//
// SendDPMonitorMod sends dp monitor reply.
//
func (c *DBCtl) SendDPMonitorMod(dpID uint64, msg *fibcapi.DpMonitorReply) error {
	eid := NewDPAPIMonitorEntryID(dpID)
	ok := c.DPSet().Select(eid, func(e fibcdbm.DPEntry) {
		mon := e.(*DPAPIMonitorEntry)
		mon.SendMod(msg)
	})

	if !ok {
		return fmt.Errorf("dp entry not found. dpid=%d", dpID)
	}

	return nil
}

//
// SendVSMonitorReply sends vs monitor reply.
//
func (c *DBCtl) SendVSMonitorReply(vsID uint64, msg *fibcapi.VsMonitorReply) error {
	eid := NewVSAPIMonitorEntryID(vsID)
	ok := c.VSSet().Select(eid, func(e fibcdbm.DPEntry) {
		mon := e.(*VSAPIMonitorEntry)
		mon.Send(msg)
	})

	if !ok {
		return fmt.Errorf("vs entry not found. vsid:%d", vsID)
	}

	return nil
}

//
// NewPortEntryFromPortConfig returns new Port entry.
//
func (c *DBCtl) NewPortEntryFromPortConfig(pc *fibcapi.PortConfig) (*fibcdbm.PortEntry, error) {
	dpID, ok := c.IDMap().SelectByReID(pc.ReId)
	if !ok {
		return nil, fmt.Errorf("re_id not found. re_id:%s", pc.ReId)
	}

	entry := &fibcdbm.PortEntry{}
	entry.Key = fibcdbm.NewPortKey(pc.ReId, pc.Ifname)
	entry.VMPort = fibcdbm.NewPortValueR(pc.ReId, pc.PortId, true)

	if pc.DpPort != 0 {
		entry.DPPort = fibcdbm.NewPortValue(dpID, pc.DpPort, true)
	}

	if len(pc.Link) != 0 {
		entry.ParentKey = fibcdbm.NewPortKey(pc.ReId, pc.Link)
	}

	if len(pc.Master) != 0 {
		entry.MasterKey = fibcdbm.NewPortKey(pc.ReId, pc.Master)
	}

	return entry, nil
}

func (c *DBCtl) UpdateNetconfConfig(dpID uint64, ports []*fibcapi.FFPort) {
	c.log.Debugf("UpdateNetconfConfig: dpID=%d", dpID)

	if err := c.nccfg.Load(); err != nil {
		c.log.Errorf("UpdateNetconfConfig: load error. %s", err)
		return
	}

	dpCfg := fibcdbm.NewNetconfDpConfig()
	for _, port := range ports {
		ifname := port.Name
		if len(ifname) == 0 {
			c.PortMap().SelectByDP(dpID, port.PortNo, func(e *fibcdbm.PortEntry) {
				ifname = e.Key.Ifname
			})
		}

		portCfg := fibcdbm.NewNetconfPortConfig(ifname, port.HwAddr, port.PortNo)
		dpCfg.AddPort(portCfg)

		c.log.Debugf("UpdateNetconfConfig: add %s", portCfg)
	}

	c.nccfg.Dps[dpID] = dpCfg
	if err := c.nccfg.Save(); err != nil {
		c.log.Errorf("UpdateNetconfConfig: save error. %s", err)
		return
	}

	c.log.Debugf("UpdateNetconfConfig: success. dpID=%d", dpID)
}
