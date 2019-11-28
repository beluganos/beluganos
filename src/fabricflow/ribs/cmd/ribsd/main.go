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
	"fabricflow/ribs/pkgs/ribscfg"
	"fabricflow/ribs/pkgs/ribssrv"
	fflibnet "fabricflow/util/net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

//
// App is application
//
type App struct {
	configFile string
	verbose    bool
	trace      bool

	log *log.Entry
}

//
// NewApp returns new App.
//
func NewApp() *App {
	return &App{
		log: log.WithFields(log.Fields{"module": "app"}),
	}
}

func (a *App) parseArgs() {
	flag.StringVarP(&a.configFile, "config-file", "c", "/etc/beluganos/risd.conf", "config filename.")
	flag.BoolVarP(&a.verbose, "verbose", "v", false, "show detail messages.")
	flag.BoolVarP(&a.trace, "trace", "", false, "show more detail messages.")

	flag.Parse()
}

func (a App) newMicService(cfg *ribscfg.Config) *ribssrv.MicService {
	a.log.Infof("main(mic):")

	return &ribssrv.MicService{
		RibsService: ribssrv.RibsService{
			NLAAddr:  cfg.NLA.API,
			BgpdAddr: cfg.GetBgpdAddr(),
			CoreAddr: cfg.Ribs.Core,
			APIAddr:  cfg.Ribs.API,
			Family:   cfg.Ribs.Bgp.RouteFamily,
		},
		SyncTime:  time.Duration(cfg.Ribs.SyncTime) * time.Millisecond,
		NexthopNW: cfg.Ribs.Nexthops.Args,
	}
}

func (a App) newRicService(cfg *ribscfg.Config) *ribssrv.RicService {
	a.log.Infof("main(ric):")

	return &ribssrv.RicService{
		RibsService: ribssrv.RibsService{
			NLAAddr:  cfg.NLA.API,
			BgpdAddr: cfg.GetBgpdAddr(),
			CoreAddr: cfg.Ribs.Core,
			Family:   cfg.Ribs.Bgp.RouteFamily,
			RT:       cfg.Ribs.Vrf.Rt,
		},
		RD:      cfg.Ribs.Vrf.Rd,
		Labels:  []uint32{cfg.VrfLabel()},
		DummyIF: cfg.Ribs.Vrf.Iface,
	}
}

func (a App) startRibsServer(cfg *ribscfg.Config, done <-chan struct{}) error {
	if cfg.Node.NId == 0 {
		return a.newMicService(cfg).Start(done)
	}

	return a.newRicService(cfg).Start(done)
}

func (a App) run() error {
	a.parseArgs()

	if a.trace {
		log.SetLevel(log.TraceLevel)
	} else if a.verbose {
		log.SetLevel(log.DebugLevel)
	}

	var cfg ribscfg.Config
	if err := ribscfg.ReadConfig(a.configFile, &cfg); err != nil {
		a.log.Errorf("run: ReadConfig error. %s", err)
		return err
	}

	a.log.Infof("config: %s", a.configFile)
	ribssrv.LogConfig(a.log, log.InfoLevel, &cfg)

	if cfg.Ribs.Disable {
		a.log.Infof("RIBS is disabled.")
		return nil
	}

	if nid, err := fflibnet.GetNIdFromLink(cfg.Node.NIdIfname); err == nil {
		a.log.Infof("nid is created from %s. nid:%d", cfg.Node.NIdIfname, nid)
		cfg.Node.NId = nid
	}

	done := make(chan struct{})
	defer close(done)

	if err := a.startRibsServer(&cfg, done); err != nil {
		a.log.Errorf("run: start ribs server error. %s", err)
		return err
	}

	<-done

	return nil
}

func main() {
	if err := NewApp().run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
