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

	log "github.com/sirupsen/logrus"
)

//
// VMAPIMonitorEntry is entry foe FIBCVmApi.Monitor.
//
type VMAPIMonitorEntry struct {
	stream fibcapi.FIBCVmApi_MonitorServer
	reID   string
	monCh  chan *fibcapi.VmMonitorReply

	active bool
	log    *log.Entry
}

//
// NewVMAPIMonitorEntry returns new VmApiMonitorEntry.
//
func NewVMAPIMonitorEntry(stream fibcapi.FIBCVmApi_MonitorServer, reID string) *VMAPIMonitorEntry {
	return &VMAPIMonitorEntry{
		stream: stream,
		reID:   reID,
		monCh:  make(chan *fibcapi.VmMonitorReply),
		log:    log.WithFields(log.Fields{"module": "vmmon", "reid": reID}),
	}
}

//
// EntryID returns entry-id.
//
func (m *VMAPIMonitorEntry) EntryID() string {
	return NewVMAPIMonitorEntryID(m.reID)
}

//
// NewVMAPIMonitorEntryID returns entry-id.
//
func NewVMAPIMonitorEntryID(reID string) string {
	return reID
}

//
// Send enqueue message to send queue.
//
func (m *VMAPIMonitorEntry) Send(msg *fibcapi.VmMonitorReply) {
	if m.active {
		m.monCh <- msg
	}
}

//
// Start starts serve thread.
//
func (m *VMAPIMonitorEntry) Start(done <-chan struct{}) {
	m.active = true
	go m.Serve(done)

	m.log.Infof("started.")
}

//
// Stop deactivate and close chan
//
func (m *VMAPIMonitorEntry) Stop() {
	if m.active {
		m.active = false
		close(m.monCh)
	}

	m.log.Infof("stopped.")
}

//
// Serve process messages.
//
func (m *VMAPIMonitorEntry) Serve(done <-chan struct{}) {
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

		case <-done:
			m.log.Debugf("Serve: EXIT.")
			m.Stop()
			break FOR_LOOP
		}
	}
}
