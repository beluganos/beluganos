// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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
	"fabricflow/ribp/ribpkt"
	"fabricflow/util/net"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func main() {

	var path string
	var verbose bool
	flag.StringVar(&path, "config", "/etc/fabricflow/ribpd.conf", "config flle")
	flag.BoolVar(&verbose, "verbose", false, "show detail log")
	flag.Parse()

	var cfg Config
	if err := ReadConfig(path, &cfg); err != nil {
		log.Errorf("MAIN: ReadConfig error. %s", err)
		os.Exit(1)
	}

	log.Infof("MAIN: config=%v", cfg)

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	done := make(chan struct{})

	nid, err := fflibnet.GetNIdFromLink(cfg.Node.NIdIfname)
	if err != nil {
		log.Infof("node.nid(conig:%d) is used as nid. reason:'%s'", nid, err)
		nid = cfg.Node.NId
	}

	s := ribpkt.NewServer(cfg.Node.ReId, nid, cfg.Ribp.Interval, cfg.Node.DupIfname)
	if err := s.Start(done); err != nil {
		log.Errorf("RIBP Server Start error. %s", err)
		return
	}

	a := ribpkt.NewApiServer(cfg.Ribp.Api)
	if err := a.Start(s.CtrlCh()); err != nil {
		log.Errorf("API Server Start error. %s", err)
		return
	}

	if err := netlink.LinkSubscribe(s.LinkCh(), done); err != nil {
		log.Errorf("Netlink Subscribe error. %s", err)
		return
	}

	<-done
}
