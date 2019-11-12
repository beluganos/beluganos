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
	"time"

	log "github.com/sirupsen/logrus"
)

//
// Server is main server.
//
type Server struct {
	fibc    FIBController
	Path    string
	UpdTime time.Duration

	StatsNames []string
}

//
// NewServer returns new instance.
//
func NewServer(fibc FIBController, path string, updTime time.Duration) *Server {
	return &Server{
		fibc:    fibc,
		Path:    path,
		UpdTime: updTime,

		StatsNames: []string{},
	}
}

func (s *Server) SetStatsNames(names []string) {
	s.StatsNames = names
}

//
// Serve serve main process.
//
func (s *Server) Serve(done <-chan struct{}) {
	log.Infof("Server: started.")

	ticker := time.NewTicker(s.UpdTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Update()

		case <-done:
			log.Infof("Server: Exit")
			return
		}
	}
}

//
// Update gets port stats by http and update port stat file.
//
func (s *Server) Update() {
	dpIds, err := s.fibc.Dps()
	if err != nil {
		log.Errorf("Update: Dps error. %s", err)
		return
	}

	log.Debugf("Update: %v", dpIds)

	if len(dpIds) == 0 {
		log.Errorf("Update: DpId not found.")
		return
	}

	dpid := dpIds[0]

	log.Debugf("Update: dpid = %d", dpid)

	ps, err := s.fibc.PortStats(dpid, s.StatsNames)
	if err != nil {
		log.Errorf("Update: PortStats error. %s", err)
		return
	}

	tempPath, err := writeToTemp(ps)
	if err != nil {
		log.Errorf("Update: write to temp error. %s", err)
		return
	}

	log.Debugf("Update: temp:%s -> %s", tempPath, s.Path)

	if err := moveFile(tempPath, s.Path); err != nil {
		log.Errorf("Update: move file error. %s", err)
		return
	}
}
