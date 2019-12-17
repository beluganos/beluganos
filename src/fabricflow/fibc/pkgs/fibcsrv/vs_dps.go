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
// VSAPIMonitorEntry is entry for FIBCVsApi.Monitor.
//
type VSAPIMonitorEntry struct {
	stream fibcapi.FIBCVsApi_MonitorServer
	vsID   uint64
<<<<<<< HEAD
=======
	remote string
>>>>>>> develop
	monCh  chan *fibcapi.VsMonitorReply

	active bool
	log    *log.Entry
}

//
// NewVSAPIMonitorEntry returns new VsApiMonitorEntry.
//
func NewVSAPIMonitorEntry(stream fibcapi.FIBCVsApi_MonitorServer, vsID uint64) *VSAPIMonitorEntry {
<<<<<<< HEAD
	return &VSAPIMonitorEntry{
		stream: stream,
		vsID:   vsID,
=======
	remote, _ := GrpcRemoteHostPort(stream)
	return &VSAPIMonitorEntry{
		stream: stream,
		vsID:   vsID,
		remote: remote,
>>>>>>> develop
		monCh:  make(chan *fibcapi.VsMonitorReply),
		log:    log.WithFields(log.Fields{"module": "vsmon", "vsid": vsID}),
	}
}

//
// EntryID returns entry-id.
//
func (m *VSAPIMonitorEntry) EntryID() string {
	return NewVSAPIMonitorEntryID(m.vsID)
}

//
<<<<<<< HEAD
=======
// Remote returns remote addr.
//
func (m *VSAPIMonitorEntry) Remote() string {
	return m.remote
}

//
>>>>>>> develop
// NewVSAPIMonitorEntryID returns entry-id.
//
func NewVSAPIMonitorEntryID(vsID uint64) string {
	return fmt.Sprintf("%d", vsID)
}

//
// Send enqueue message to send queue.
//
func (m *VSAPIMonitorEntry) Send(msg *fibcapi.VsMonitorReply) {
	if m.active {
		m.monCh <- msg
	}
}

//
// Start starts serve thread.
//
func (m *VSAPIMonitorEntry) Start(done <-chan struct{}) {
	m.active = true
	go m.Serve(done)

	m.log.Infof("started.")
}

//
// Stop deactivste and close chan
//
func (m *VSAPIMonitorEntry) Stop() {
	if m.active {
		m.active = false
		close(m.monCh)
	}
	m.log.Infof("stopped.")
}

//
// Serve process message.
//
func (m *VSAPIMonitorEntry) Serve(done <-chan struct{}) {
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
