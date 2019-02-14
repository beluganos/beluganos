// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

//
// StatsHandlerEntry is handler entry.
//
type StatsHandlerEntry struct {
	Oid     string
	Handler StatsHandler
}

//
// SnmpServer is main server.
//
type SnmpServer struct {
	handlers []*StatsHandlerEntry
}

//
// NewSnmpServer returns new instance.
//
func NewSnmpServer() *SnmpServer {
	return &SnmpServer{
		handlers: []*StatsHandlerEntry{},
	}
}

//
// AddStatsHandler adds handler entry of port stats.
//
func (s *SnmpServer) AddStatsHandler(handler StatsHandler) {
	h := &StatsHandlerEntry{
		Oid:     handler.Oid(),
		Handler: handler,
	}
	s.handlers = append(s.handlers, h)
}

//
// AddStatsHandlers adds handler entries of port stats.
//
func (s *SnmpServer) AddStatsHandlers(handlers ...StatsHandler) {
	for _, handler := range handlers {
		s.AddStatsHandler(handler)
	}
}

//
// getStatsHandler returns entry.
//
func (s *SnmpServer) getStatsHandler(oid string) StatsHandler {
	oid = oid + "."
	for _, entry := range s.handlers {
		prefix := entry.Oid + "."
		if ok := strings.HasPrefix(oid, prefix); ok {
			return entry.Handler
		}
	}

	return nil
}

//
// OnGet process get request.
//
func (s *SnmpServer) OnGet(oid string) *SnmpReply {
	log.Debugf("GET: '%s'", oid)

	if handler := s.getStatsHandler(oid); handler != nil {
		return handler.Get(oid)
	}

	return nil
}

//
// OnGetNext process getnext request.
//
func (s *SnmpServer) OnGetNext(oid string) *SnmpReply {
	log.Debugf("GET NEXT: '%s'", oid)

	if handler := s.getStatsHandler(oid); handler != nil {
		return handler.GetNext(oid)
	}

	return nil
}

//
// Serve receives and dispatch request.
//
func (s *SnmpServer) Serve(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		cmd := scanner.Text()
		log.Debugf("Stdin: '%s'", cmd)

		switch cmd {
		case "":
			log.Debugf("EOL")
			return

		case "PING":
			fmt.Fprintln(w, "PONG")

		case "get":
			if !scanner.Scan() {
				break
			}

			oid := scanner.Text()
			s.OnGet(oid).WriteTo(w)

		case "getnext":
			if !scanner.Scan() {
				break
			}

			oid := scanner.Text()
			s.OnGetNext(oid).WriteTo(w)

		default:
			log.Debugf("Unknown cmd. '%s'", cmd)
		}
	}
}
