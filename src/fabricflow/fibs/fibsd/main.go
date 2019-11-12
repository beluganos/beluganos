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
	lib "fabricflow/fibs/fibslib"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

//
// Args is arguments.
//
type Args struct {
	Addr     string
	FibcType string
	Path     string
	UpdTime  time.Duration
	Verbose  bool

	ConfigFile string
	ConfigType string
	ConfigName string
}

//
// NewArgs returns new argument.
//
func NewArgs() *Args {
	args := Args{}
	args.Init()
	return &args
}

//
// Init parse and get arguments.
//
func (a *Args) Init() {
	flag.StringVarP(&a.Addr, "fibc-addr", "a", "localhost:8080", "fibc address.")
	flag.StringVarP(&a.FibcType, "fibc-type", "", FIBCTypeDefault, "fib connection type.")
	flag.StringVarP(&a.Path, "stats-path", "p", lib.FIBS_STATS_FILEPATH, "fibc stats filename.")
	flag.DurationVarP(&a.UpdTime, "update-time", "u", 15*time.Minute, "update stats interval time.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail messages.")

	flag.StringVarP(&a.ConfigFile, "config-file", "", "", "config file path.")
	flag.StringVarP(&a.ConfigType, "config-type", "", "yaml", "config file type.")
	flag.StringVarP(&a.ConfigName, "config-name", "", "default", "config name.")

	flag.Parse()
}

func main() {
	args := NewArgs()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Args: addr   = '%s'", args.Addr)
	log.Infof("Args: type   = '%s'", args.FibcType)
	log.Infof("Args: path   = '%s'", args.Path)
	log.Infof("Args: update = %s", args.UpdTime)

	fibctl := NewFIBController(args.FibcType, args.Addr)

	done := make(chan struct{})
	s := NewServer(fibctl, args.Path, args.UpdTime)
	if len(args.ConfigFile) > 0 {
		cfg := NewConfig()
		if err := cfg.ReadFile(args.ConfigFile, args.ConfigType); err != nil {
			log.Warnf("ReadConfig error. %s", err)
		} else {
			if fibscfg := cfg.GetFibsConfig(args.ConfigName); fibscfg != nil {
				s.SetStatsNames(fibscfg.FibsStats.Names)
			} else {
				log.Warnf("FibsConfig not found. %s", args.ConfigName)
			}
		}
	}

	s.Serve(done)

	<-done
}
