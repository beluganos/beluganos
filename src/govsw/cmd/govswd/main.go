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
	"fmt"
	"govsw/pkgs/govsw"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	ARGS_CONFIG_PATH = "/etc/beluganos/govswd.yaml"
	ARGS_CONFIG_TYPE = "yaml"
	ARGS_API_ADDR    = "localhost"
	ARGS_API_PORT    = govsw.VSWAPI_PORT
	ARGS_FIBC_ADDR   = "localhost"
	ARGS_FIBC_PORT   = 50061
	ARGS_STRIP_WPKT  = 4
	ARGS_STRIP_RPKT  = 0
)

type Args struct {
	DpName string

	ConfigPath string
	ConfigType string

	ApiAddr  string
	ApiPort  uint16
	FibcAddr string
	FibcPort uint16

	StripWritePkt uint16
	StripReadPkt  uint16

	Verbose bool
	Trace   bool
}

func (a *Args) Parse() {
	flag.StringVarP(&a.DpName, "dp", "", "default", "Datapath name.")
	flag.StringVarP(&a.ConfigPath, "config-file", "c", ARGS_CONFIG_PATH, "Config filename.")
	flag.StringVarP(&a.ConfigType, "config-type", "", ARGS_CONFIG_TYPE, "Config file type.")
	flag.StringVarP(&a.ApiAddr, "api-addr", "", ARGS_API_ADDR, "api listen address.")
	flag.Uint16VarP(&a.ApiPort, "api-port", "", ARGS_API_PORT, "api listen port.")
	flag.StringVarP(&a.FibcAddr, "fibc-addr", "", ARGS_FIBC_ADDR, "fibcd address.")
	flag.Uint16VarP(&a.FibcPort, "fibc-port", "", ARGS_FIBC_PORT, "fibcd port.")
	flag.Uint16VarP(&a.StripWritePkt, "strip-write-pkt", "", ARGS_STRIP_WPKT, "strip size for write packet.")
	flag.Uint16VarP(&a.StripReadPkt, "strip-read-pkt", "", ARGS_STRIP_RPKT, "strip size for read packet.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail messages.")
	flag.BoolVarP(&a.Trace, "trace", "", false, "show detail messages.")
	flag.Parse()
}

func NewArgs() *Args {
	args := &Args{}
	args.Parse()
	return args
}

func (a *Args) ApiListenAddr() string {
	return fmt.Sprintf("%s:%d", a.ApiAddr, a.ApiPort)
}

func (a *Args) FIBCListenAddr() string {
	return fmt.Sprintf("%s:%d", a.FibcAddr, a.FibcPort)
}

func main() {
	args := NewArgs()

	if args.Trace {
		log.SetLevel(log.TraceLevel)
	} else if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	done := make(chan struct{})
	defer close(done)

	db := govsw.NewDB()
	db.Link().SetStripWPkt(args.StripWritePkt)
	db.Link().SetStripRPkt(args.StripReadPkt)

	cfgsrv := govsw.ConfigServer{
		DpName:  args.DpName,
		CfgPath: args.ConfigPath,
		CfgType: args.ConfigType,
	}

	app := NewMyApp(db, &cfgsrv)

	if err := cfgsrv.Start(app); err != nil {
		log.Errorf("ConfigSerer start error. %s", err)
		os.Exit(1)
	}

	apisrv := &govsw.VswApiServer{
		DB:       db,
		Listener: app,
		Addr:     args.ApiListenAddr(),
	}
	if err := apisrv.Start(); err != nil {
		log.Errorf("VswApiServer start error. %s", err)
		os.Exit(1)
	}

	fibcsrv := FIBCVsApiClient{
		Listener: app,
		Addr:     args.FIBCListenAddr(),
	}
	if err := fibcsrv.Start(done); err != nil {
		log.Errorf("FIBCDpApiClient Start error. %s", err)
		os.Exit(1)
	}

	srv := &govsw.Server{
		Listener: app,
		DB:       db,
		SyncCh:   app.SyncCh(),
	}
	if err := srv.Start(done); err != nil {
		log.Errorf("server start error. %s", err)
		os.Exit(1)
	}

	<-done
}
