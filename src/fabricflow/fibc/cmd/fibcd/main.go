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

package main

import (
	"fabricflow/fibc/pkgs/fibccfg"
	"fabricflow/fibc/pkgs/fibcsrv"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	argsConfigDir    = "/etc/beluganos/fibc.d"
	argsConfigType   = "yaml"
	argsConfigPrefix = ""
	argsListenNW     = "tcp"
	argsListenIP     = "localhost"
	argsListenPort   = 50061
	argNcConfigPath  = "/tmp/ncmi.yaml"
	argNcConfigType  = "yaml"
)

//
// App is application
//
type App struct {
	ConfigDir    string
	ConfigType   string
	ConfigPrefix string
	ListenNW     string
	ListenIP     string
	ListenPort   uint16
	NcConfigPath string
	NcConfigType string

	Verbose bool
	Trace   bool

	log *log.Entry
}

func (a *App) parseArgs() {
	flag.StringVarP(&a.ConfigDir, "config-dir", "d", argsConfigDir, "config dir path.")
	flag.StringVarP(&a.ConfigType, "config-type", "", argsConfigType, "config file type.")
	flag.StringVarP(&a.ConfigPrefix, "config-prefix", "", argsConfigPrefix, "config prefix.")
	flag.StringVarP(&a.ListenNW, "listen-network", "", argsListenNW, "listen network.")
	flag.StringVarP(&a.ListenIP, "listen-addr", "a", argsListenIP, "listen address.")
	flag.Uint16VarP(&a.ListenPort, "listen-port", "p", argsListenPort, "listen port.")
	flag.StringVarP(&a.NcConfigPath, "netconf-config-file", "", argNcConfigPath, "netconf config path.")
	flag.StringVarP(&a.NcConfigType, "netconf-config-type", "", argNcConfigType, "netconf config type.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show deail messages.")
	flag.BoolVarP(&a.Trace, "trace", "", false, "show more deail messages.")
	flag.Parse()
}

func (a *App) dumpArgs() {
	a.log.Infof("config-dir     : '%s'", a.ConfigDir)
	a.log.Infof("config-type    : '%s'", a.ConfigType)
	a.log.Infof("config-prefix  : '%s'", a.ConfigPrefix)
	a.log.Infof("listen-network : '%s'", a.ListenNW)
	a.log.Infof("listen-address : '%s'", a.ListenIP)
	a.log.Infof("listen-port    : %d", a.ListenPort)
	a.log.Infof("nc-config-file : '%s'", a.NcConfigPath)
	a.log.Infof("nc-config-type : '%s'", a.NcConfigType)
	a.log.Infof("verbose        : %t", a.Verbose)
	a.log.Infof("trace          : %t", a.Trace)
}

func (a *App) dumpConfig(cfg *fibccfg.Config) {
	for _, router := range cfg.Routers {
		a.log.Debugf("router: %v", router)

		for _, port := range router.Ports {
			a.log.Debugf(" - port: %v", port)
		}
	}

	for _, dpath := range cfg.DPaths {
		a.log.Debugf("dpath: %v", dpath)
	}
}

func newApp() *App {
	app := App{
		log: log.WithFields(log.Fields{"module": "main"}),
	}
	app.parseArgs()
	return &app
}

func (a *App) run() error {
	if a.Trace {
		log.SetLevel(log.TraceLevel)
	} else if a.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	a.dumpArgs()

	cfg, err := fibccfg.ReadConfigs(a.ConfigDir, a.ConfigPrefix, a.ConfigType)
	if err != nil {
		a.log.Errorf("ListFile error. %s", err)
		return err
	}

	a.dumpConfig(cfg)

	listenAddr := fmt.Sprintf("%s:%d", a.ListenIP, a.ListenPort)
	lis, err := net.Listen(a.ListenNW, listenAddr)
	if err != nil {
		a.log.Errorf("Listen error. %s %s", listenAddr, err)
		return err
	}

	s := fibcsrv.NewServer()
	s.SetConfig(cfg)
	s.SetNetconfConfig(a.NcConfigPath, a.NcConfigType)
	s.Serve(lis)

	return nil
}

func main() {
	if err := newApp().run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
