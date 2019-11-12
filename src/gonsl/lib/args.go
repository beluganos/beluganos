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

package gonslib

import (
	flag "github.com/spf13/pflag"
)

//
// Args is gonsl arguments.
//
type Args struct {
	ConfigFile string
	ConfigType string
	UseSim     bool
	DpName     string
	APIAddr    string
	Daemon     bool
	PidFile    string
	LogFile    string
	Verbose    bool
	Trace      bool
}

//
// NewArgs returns new instance.
//
func NewArgs() *Args {
	args := &Args{}
	args.Parse()
	return args
}

//
// Parse parses argument and get.
//
func (a *Args) Parse() {
	flag.StringVar(&a.DpName, "dp", "default", "Datapath name.")
	flag.StringVarP(&a.ConfigFile, "config-file", "c", "/etc/gonsl/gonsld.yaml", "Config filename.")
	flag.StringVar(&a.ConfigType, "config-type", "yaml", "Config file type.")
	flag.StringVarP(&a.APIAddr, "api-addr", "a", ":50061", "Listen addr for grpc.")
	flag.BoolVar(&a.UseSim, "use-sim", false, "use simulator.")
	flag.BoolVar(&a.Daemon, "daemon", false, "run as daemon.")
	flag.StringVar(&a.PidFile, "pid", "/var/run/gonsld.pid", "PID file.")
	flag.StringVar(&a.LogFile, "log-file", "/var/log/gonsld.log", "log file.")
	flag.BoolVarP(&a.Verbose, "verboce", "v", false, "show detail messages.")
	flag.BoolVarP(&a.Trace, "trace", "", false, "show more detail messages.")
	flag.Parse()
}
