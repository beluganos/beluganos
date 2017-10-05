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
	"fabricflow/ribs/ribsyn"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	var configPath string
	var verbose bool
	flag.StringVar(&configPath, "config", "/etc/fabricflow/ribsd.conf", "config file path.")
	flag.BoolVar(&verbose, "verbose", false, "show detail log.")
	flag.Parse()

	var cfg Config
	if err := ReadConfig(configPath, &cfg); err != nil {
		log.Errorf("MAIN: ReadConfig error. %s", err)
		os.Exit(1)
	}

	log.Infof("MAIN: config=:%s", cfg)

	if cfg.Ribs.Disable {
		log.Infof("MAIN: Disabled.")
		os.Exit(0)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	ribsyn.CreateTables(cfg.Ribs.Nexthops.Args)

	startApi(&cfg)

	if cfg.Node.NId == 0 {
		micMain(&cfg)
	} else {
		ricMain(&cfg)
	}
}

func startApi(cfg *Config) {
	apiSrv := ribsyn.NewApiServer(cfg.Ribs.Api)
	if err := apiSrv.Start(); err != nil {
		log.Errorf("MAIN: ApiServer.Start error. %s", err)
		os.Exit(1)
	}
}

func micMain(cfg *Config) {
	log.Infof("MAIN: MIC mode")

	s := ribsyn.NewMicService(cfg.Ribs.SyncTime)
	if err := s.Start(cfg.NLA.Api, cfg.GetBgpAddr(), cfg.Ribs.Core); err != nil {
		log.Errorf("MAIN: MicService.Start error. %s", err)
		os.Exit(1)
	}

	s.Serve()
}

func ricMain(cfg *Config) {
	log.Infof("MAIN: RIC mode")

	s, err := ribsyn.NewRicService(cfg.Node.NId, cfg.Ribs.Vrf.Rt, cfg.Ribs.Vrf.Iface, cfg.Ribs.Core)
	if err != nil {
		log.Errorf("MAIN: RicService create error. %s", err)
		os.Exit(1)
	}

	if err := s.Start(cfg.Ribs.Bgp.Addr, cfg.Ribs.Bgp.Port, cfg.Ribs.Vrf.Rd, cfg.VrfLabel()); err != nil {
		log.Errorf("MAIN: RicService.Start error. %s", err)
		os.Exit(1)
	}

	s.Serve()
}
