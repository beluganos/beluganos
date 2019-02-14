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
	"io/ioutil"
	"log/syslog"
	"os"

	lib "fabricflow/fibs/fibslib"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	flag "github.com/spf13/pflag"
)

//
// Args is arguments.
//
type Args struct {
	DataPath   string
	DataFormat string
	HandlerCfg string
	Verbose    bool
	Stdout     bool
}

const (
	CONFIG_FILENAME_DEFAULT = "/etc/beluganos/fibssnmp.yaml"
	CONFIG_FILETYPE_DEFAULT = "yaml"
)

//
// NewArgs returns new instance.
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
	flag.StringVarP(&a.DataPath, "data-path", "", lib.FIBS_STATS_FILEPATH, "stats filepath.")
	flag.StringVarP(&a.DataFormat, "data-format", "", CONFIG_FILETYPE_DEFAULT, "stats file format.")
	flag.StringVarP(&a.HandlerCfg, "handlers", "", CONFIG_FILENAME_DEFAULT, "config file.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail messages.")
	flag.BoolVarP(&a.Stdout, "stdout", "", false, "show detail messages on stdout.")
	flag.Parse()
}

func initLog(verbose bool, stdout bool) {

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if !stdout {
		log.SetOutput(ioutil.Discard)

		hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
		if err == nil {
			log.AddHook(hook)
		}
	}
}

func main() {
	args := NewArgs()

	initLog(args.Verbose, args.Stdout)

	log.Infof("fibssnmp start.")

	cfg, err := ReadConfigFile(args.HandlerCfg)
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	ds := NewDataServer()
	if err := ds.Init(args.DataPath, args.DataFormat); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	handlers, err := NewStatsHandlers(cfg.Handlers, ds)
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	for _, hdr := range handlers {
		log.Debugf("Handler: %s", hdr)
	}

	ss := NewSnmpServer()
	ss.AddStatsHandlers(handlers...)
	ss.Serve(os.Stdin, os.Stdout)
}
