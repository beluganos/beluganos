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
	"fmt"

	log "github.com/sirupsen/logrus"
)

//
// DPAPIMonitorEntry is entry for FIBCDpApi.Monitor.
//
type DPAPIMonitorEntry struct {
	stream fibcapi.FIBCDpApi_MonitorServer
	dpID   uint64
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	remote string
>>>>>>> develop
=======
	remote string
>>>>>>> develop
=======
	remote string
>>>>>>> develop
	monCh  chan *fibcapi.DpMonitorReply
	modCh  chan *fibcapi.DpMonitorReply
	db     *DBCtl

	active bool
	log    *log.Entry
}

//
// NewDPAPIMonitorEntry retuens new DpApiMonitorEntry.
//
func NewDPAPIMonitorEntry(stream fibcapi.FIBCDpApi_MonitorServer, dpID uint64, db *DBCtl) *DPAPIMonitorEntry {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	return &DPAPIMonitorEntry{
		stream: stream,
		dpID:   dpID,
=======
=======
>>>>>>> develop
=======
>>>>>>> develop
	remote, _ := GrpcRemoteHostPort(stream)
	return &DPAPIMonitorEntry{
		stream: stream,
		dpID:   dpID,
		remote: remote,
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
		monCh:  make(chan *fibcapi.DpMonitorReply),
		modCh:  make(chan *fibcapi.DpMonitorReply),
		db:     db,

		log: log.WithFields(log.Fields{"module": "dpmon", "dpid": dpID}),
	}
}

//
// EntryID returns entry-id.
//
func (m *DPAPIMonitorEntry) EntryID() string {
	return NewDPAPIMonitorEntryID(m.dpID)
}

//
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> develop
=======
>>>>>>> develop
// Remote returns remote addr.
//
func (m *DPAPIMonitorEntry) Remote() string {
	return m.remote
}

//
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
// NewDPAPIMonitorEntryID returns entry-id
//
func NewDPAPIMonitorEntryID(dpID uint64) string {
	return fmt.Sprintf("%d", dpID)
}

//
// Send enqueue message to send queue.
//
func (m *DPAPIMonitorEntry) Send(msg *fibcapi.DpMonitorReply) {
	if m.active {
		m.monCh <- msg
	}
}

//
// SendMod enqueue message to send queue.
//
func (m *DPAPIMonitorEntry) SendMod(msg *fibcapi.DpMonitorReply) {
	if m.active {
		m.modCh <- msg
	}
}

//
// Start stats serve thread.
//
func (m *DPAPIMonitorEntry) Start(done <-chan struct{}) {
	m.active = true
	go m.Serve(done)

	m.log.Infof("started.")
}

//
// Stop deactivate and close chan.
//
func (m *DPAPIMonitorEntry) Stop() {
	if m.active {
		m.active = false
		close(m.monCh)
	}
	m.log.Infof("stopped.")
}

func (m *DPAPIMonitorEntry) convertMod(msg *fibcapi.DpMonitorReply) error {

	switch mod := msg.Body.(type) {
	case *fibcapi.DpMonitorReply_FlowMod:
		if err := m.db.ConvertFlowMod(mod.FlowMod); err == ENoEffect {
			m.log.Debugf("server: ignore flow mod. %s", err)
			return err

		} else if err != nil {
			m.log.Errorf("server: convert flow mod error. %s", err)
			return err
		}

		fibcapi.LogFlowMod(m.log, log.TraceLevel, mod.FlowMod)
		return nil

	case *fibcapi.DpMonitorReply_GroupMod:
		if err := m.db.ConvertGroupMod(mod.GroupMod); err == ENoEffect {
			m.log.Debugf("server: ignore group mod. %s", err)
			return err

		} else if err != nil {
			m.log.Errorf("server: convert grop mod error. %s", err)
			return err
		}

		fibcapi.LogGroupMod(m.log, log.TraceLevel, mod.GroupMod)
		return nil

	default:
		// pass
		return nil
	}
}

//
// Serve process messages.
//
func (m *DPAPIMonitorEntry) Serve(done <-chan struct{}) {
	m.log.Debugf("Serve: START.")

FOR_LOOP:
	for {
		select {
		case msg := <-m.monCh:
			if msg != nil {
				if err := m.stream.Send(msg); err != nil {
					m.log.Errorf("Serve: send error. %s", err)
					break FOR_LOOP
				}
			}

		case msg := <-m.modCh:
			if msg != nil {
				if err := m.convertMod(msg); err == nil {
					if err := m.stream.Send(msg); err != nil {
						m.log.Errorf("Serve: send mod error. %s", err)
						break FOR_LOOP
					}
				}
			}

		case <-done:
			m.log.Debugf("Serve: EXIT.")
			m.Stop()
			break FOR_LOOP
		}
	}
}
