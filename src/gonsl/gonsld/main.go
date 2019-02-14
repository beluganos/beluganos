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
	"gonsl/lib"
	"os"
	"os/signal"

	"github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
)

func watchSignal(done chan struct{}) {

	ch := make(chan os.Signal, 1)
	defer close(ch)

	signal.Notify(ch, os.Interrupt)
	<-ch

	log.Infof("Interrupt signal.")

	close(done)
}

func startDaemon(args *gonslib.Args) *daemon.Context {
	ctx := &daemon.Context{
		PidFileName: args.PidFile,
		PidFilePerm: 0644,
		LogFileName: args.LogFile,
		LogFilePerm: 0640,
	}

	child, err := ctx.Reborn()
	if err != nil {
		log.Errorf("Daemonize error. %s", err)
		os.Exit(1)
	}

	if child != nil {
		// Parent process
		os.Exit(0)
	}

	// Child process.
	return ctx
}

func main() {
	args := gonslib.NewArgs()

	if args.Daemon {
		ctx := startDaemon(args)
		defer ctx.Release()
	}

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	cfg := gonslib.NewConfig()
	if _, err := cfg.ReadFile(args); err != nil {
		log.Errorf("Config read error. %s", err)
		os.Exit(1)
	}

	log.Debugf("Config: %v", cfg)

	dpcfg, err := cfg.GetDpConfig(args.DpName)
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	log.Debugf("DpConfig: %v", dpcfg)

	if args.UseSim {
		gonslib.SimInit(dpcfg.Unit)
	}

	if err := gonslib.DriverInit(dpcfg.Unit, dpcfg.OpenNSL); err != nil {
		log.Errorf("DriverInit error. %s", err)
		os.Exit(1)
	}
	defer gonslib.DriverExit()

	log.Infof("OpenNSL Driver initialized. unit=%d config=%v", dpcfg.Unit, dpcfg.OpenNSL)

	done := make(chan struct{})
	s := gonslib.NewServer(dpcfg)
	if err := s.Start(done); err != nil {
		log.Errorf("Server start error. %s", err)
		os.Exit(1)
	}

	api := gonslib.NewAPIServer(s)
	if err := api.Start(args.APIAddr); err != nil {
		log.Errorf("Server start error. %s", err)
		os.Exit(1)
	}

	go watchSignal(done)
	log.Debugf("Initialize ok: %v", dpcfg)

	<-done
}
