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
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

const (
	apAPIMonitorChanSize    = 1024
	apAPIMonitorLogChanSize = 64
)

var apAPIEntryId = uint64(0)

func getAPAPIEntryId() string {
	eid := atomic.AddUint64(&apAPIEntryId, 1)
	return fmt.Sprintf("%d", eid)
}

//
// APAPIMonitorEntry is entry for ApiMonitor
//
type APAPIMonitorEntry struct {
	stream  fibcapi.FIBCApApi_MonitorServer
	entryID string
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	remote  string
>>>>>>> develop
=======
	remote  string
>>>>>>> develop
=======
	remote  string
>>>>>>> develop
=======
	remote  string
>>>>>>> develop
	monCh   chan *fibcapi.ApMonitorReply
	logCh   chan *fibcapi.ApMonitorReplyLog

	active bool
	log    *log.Entry
}

//
// NewAPAPIMonitorEntry returns new ApiMonitorEntry.
//
func NewAPAPIMonitorEntry(stream fibcapi.FIBCApApi_MonitorServer) *APAPIMonitorEntry {
	eid := getAPAPIEntryId()
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	return &APAPIMonitorEntry{
		stream:  stream,
		entryID: eid,
=======
=======
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
	remote, _ := GrpcRemoteHostPort(stream)
	return &APAPIMonitorEntry{
		stream:  stream,
		entryID: eid,
		remote:  remote,
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
		monCh:   make(chan *fibcapi.ApMonitorReply, apAPIMonitorChanSize),
		logCh:   make(chan *fibcapi.ApMonitorReplyLog, apAPIMonitorLogChanSize),
		log:     log.WithFields(log.Fields{"module": "apmon", "eid": eid}),
	}
}

//
// EntryID returns entry-id.
//
func (m *APAPIMonitorEntry) EntryID() string {
	return m.entryID
}

//
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
// Remote returns remote addr.
//
func (m *APAPIMonitorEntry) Remote() string {
	return m.remote
}

//
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
// Send enqueue msg to send queue.
//
func (m *APAPIMonitorEntry) Send(msg *fibcapi.ApMonitorReply) {
	if m.active {
		m.monCh <- msg
	}
}

//
// SendLog enqueue log message to log queue.
//
func (m *APAPIMonitorEntry) SendLog(msg *fibcapi.ApMonitorReplyLog) {
	if m.active {
		m.logCh <- msg
	}
}

//
// Start starts serve thread.
//
func (m *APAPIMonitorEntry) Start(done <-chan struct{}) {
	m.active = true
	go m.Serve(done)

	m.log.Infof("started.")
}

//
// Stop deactivate and close chan
//
func (m *APAPIMonitorEntry) Stop() {
	if m.active {
		m.active = false
		close(m.monCh)
		close(m.logCh)
	}
	m.log.Infof("stopped.")
}

//
// Serve process messages.
//
func (m *APAPIMonitorEntry) Serve(done <-chan struct{}) {
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

		case data := <-m.logCh:
			if data != nil {
				body := fibcapi.ApMonitorReply_Log{
					Log: data,
				}
				msg := fibcapi.ApMonitorReply{
					Body: &body,
				}
				if err := m.stream.Send(&msg); err != nil {
					m.log.Warnf("Serve: send log error. %s", err)
				}
			}

		case <-done:
			m.log.Debugf("Serve: EXIT.")
			m.Stop()
			break FOR_LOOP
		}
	}
}
