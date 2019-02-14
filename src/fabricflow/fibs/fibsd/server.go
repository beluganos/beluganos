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
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

//
// Server is main server.
//
type Server struct {
	URL     string
	Path    string
	UpdTime time.Duration
}

//
// NewServer returns new instance.
//
func NewServer(addr string, path string, updTime time.Duration) *Server {
	return &Server{
		URL:     fmt.Sprintf("http://%s", addr),
		Path:    path,
		UpdTime: updTime,
	}
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

func (s *Server) WritePortStats(ps PortStats) (string, error) {
	file, err := ioutil.TempFile("", "fibs_stats")
	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}

	defer file.Close()

	if level := log.GetLevel(); level == log.DebugLevel {
		ps.WriteTo(log.StandardLogger().Writer())
	}

	os.Chmod(file.Name(), 0644)

	w := bufio.NewWriter(file)
	ps.WriteTo(w)
	w.Flush()

	return file.Name(), nil
}

//
// Update gets port stats by http and update port stat file.
//
func (s *Server) Update() {
	dps := DpsMsg{}
	if err := dps.HTTPGet(s.URL); err != nil {
		log.Errorf("%s", err)
		return
	}

	log.Debugf("%v", dps)

	if len(dps.DpIds) == 0 {
		log.Errorf("DpId not found.")
		return
	}

	dpid := dps.DpIds[0]

	log.Debugf("dpid = %d", dpid)

	psmsg := PortStatsMsg{}
	if err := psmsg.HTTPGet(s.URL, dpid); err != nil {
		log.Errorf("%s", err)
		return
	}

	ps, _ := psmsg.PortStats(dpid)
	tempPath, err := s.WritePortStats(ps)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	log.Debugf("temp:%s -> %s", tempPath, s.Path)

	if err := moveFile(tempPath, s.Path); err != nil {
		log.Errorf("%s", err)
		return
	}
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func moveFile(srcPath, dstPath string) error {
	if err := os.Rename(srcPath, dstPath); err == nil {
		return nil
	}

	if err := copyFile(srcPath, dstPath); err != nil {
		return err
	}

	return os.Remove(srcPath)
}
