// coding: utf-8 -*-

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
	"fabricflow/ribc/ribctl"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	var configPath string
	var verbose bool
	flag.StringVar(&configPath, "config", "/etc/fabricflow/ribcd.conf", "config file path.")
	flag.BoolVar(&verbose, "verbose", false, "show detail log.")
	flag.Parse()

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Errorf("LoadConfig error. %s", err)
		os.Exit(1)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("RIBC: config=%s", config)

	if config.Ribc.Disable {
		log.Infof("RIBC: Disable")
		os.Exit(0)
	}

	nla := ribctl.NewNLAController(config.NLA.Api)
	fib := ribctl.NewFIBController(config.Ribc.Fibc)
	rib := ribctl.NewRIBController(config.Node.NId, config.Node.ReId, config.Node.Label, nla, fib)

	if err := nla.Start(); err != nil {
		log.Errorf("NewNLAMonitor Start error. %s", err)
		os.Exit(1)
	}

	done := make(chan struct{})

	rib.Start(done)

	<-done
}
