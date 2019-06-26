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
	"fabricflow/ribc/ribctl"
	fflibnet "fabricflow/util/net"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type Args struct {
	ConfigPath     string
	FlowConfigPath string
	FlowConfigName string
	FlowConfigType string // yaml, toml, ...
	Verbose        bool
}

func (a *Args) Parse() {
	flag.StringVar(&a.ConfigPath, "config", "/etc/fabricflow/ribcd.conf", "config file path.")
	flag.StringVar(&a.FlowConfigPath, "flow-config-path", "", "Flow config path.")
	flag.StringVar(&a.FlowConfigName, "flow-config-name", ribctl.FLOWDB_BUILTIN_CONFIG, "Flow config name.")
	flag.StringVar(&a.FlowConfigType, "flow-config-type", "yaml", "Flow config type.")
	flag.BoolVar(&a.Verbose, "verbose", false, "show detail log.")
	flag.Parse()
}

func NewArgs() *Args {
	a := &Args{}
	a.Parse()
	return a
}

func getFlowConfig(args *Args) (*ribctl.FlowConfig, error) {
	flowdb := ribctl.NewFlowDB()
	if len(args.FlowConfigPath) != 0 {
		flowdb.SetConfigFile(args.FlowConfigPath, args.FlowConfigType)
		if err := flowdb.Load(); err != nil {
			return nil, err
		}
	}

	flowcfg := flowdb.Config(args.FlowConfigName)
	if flowcfg == nil {
		return nil, fmt.Errorf("FlowConfig not found. %s", args.FlowConfigName)
	}

	return flowcfg, nil
}

func dumpArgs(a *Args) {
	log.Infof("Args: ConfigPath     : `%s`", a.ConfigPath)
	log.Infof("Args: FlowConfigPath : `%s`", a.FlowConfigPath)
	log.Infof("Args: FlowConfigType : `%s`", a.FlowConfigType)
	log.Infof("Args: FlowConfigName : `%s`", a.FlowConfigName)
	log.Infof("Args: Verbose        : %t", a.Verbose)
}

func dumpConfig(c *Config) {
	log.Infof("CONFIG: Node.NId        : %d", c.Node.NId)
	log.Infof("CONFIG: Node.Label      : %d", c.Node.Label)
	log.Infof("CONFIG: Node.ReId       : '%s'", c.Node.ReId)
	log.Infof("CONFIG: Node.NId-ifname : '%s'", c.Node.NIdIfname)
	log.Infof("CONFIG: Node.Dup-ifname : %t", c.Node.DupIfname)
	log.Infof("CONFIG: NLA.Api         : '%s'", c.NLA.Api)
	log.Infof("CONFIG: RIBC.FIBC       : '%s'", c.Ribc.Fibc)
	log.Infof("CONFIG: RIBC.Disable    : %t", c.Ribc.Disable)
}

func main() {
	args := NewArgs()
	dumpArgs(args)

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	config, err := LoadConfig(args.ConfigPath)
	if err != nil {
		log.Errorf("RIBC: LoadConfig error. %s", err)
		os.Exit(1)
	}

	dumpConfig(config)

	if config.Ribc.Disable {
		log.Infof("RIBC: Disable")
		os.Exit(0)
	}

	flowcfg, err := getFlowConfig(args)
	if err != nil {
		log.Errorf("RIBC: get flow config error. %s", err)
		os.Exit(1)
	}

	nid, err := fflibnet.GetNIdFromLink(config.Node.NIdIfname)
	if err != nil {
		log.Infof("node.nid(conig:%d) is used as nid. reason:'%s'", nid, err)
		nid = config.Node.NId
	}

	nla := ribctl.NewNLAController(config.NLA.Api)
	fib := ribctl.NewFIBController(config.Ribc.Fibc)
	rib := ribctl.NewRIBController(nid, config.Node.ReId, config.Node.Label, config.Node.DupIfname, nla, fib, flowcfg)

	if err := nla.Start(); err != nil {
		log.Errorf("NewNLAMonitor Start error. %s", err)
		os.Exit(1)
	}

	done := make(chan struct{})

	rib.Start(done)

	<-done
}
