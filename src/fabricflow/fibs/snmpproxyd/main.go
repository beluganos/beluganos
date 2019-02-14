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
	"net"
	"os"
	"time"

	lib "fabricflow/fibs/fibslib"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	CONFIG_FILE_DEFAULT = "/etc/beluganos/snmpconvd.conf"
	DUMP_FILE_DEFAULT   = "/tmp/snmpproxy_tables"
	DUMP_TIME_DEFAULT   = 15 * time.Second
	DUMP_TIME_MIN       = 3 * time.Second
)

type Args struct {
	Config        string
	Table         string
	IfOID         string
	ListenAddr    *net.UDPAddr
	SnmpdAddr     *net.UDPAddr
	DumpTableTime time.Duration
	DumpTableFile string
	Verbose       bool
}

func (a *Args) Parse() error {
	var (
		err        error
		listenAddr string
		snmpdAddr  string
	)

	flag.StringVarP(&a.Config, "config-file", "c", CONFIG_FILE_DEFAULT, "config file.")
	flag.StringVarP(&a.Table, "table", "t", "default", "table name.")
	flag.StringVarP(&a.IfOID, "if-notify-oid", "", lib.SNMP_OID_IFACES, "iface notify OID.")
	flag.StringVarP(&listenAddr, "listen-addr", "", lib.SNMP_LISTEN_ADDR, "Listen address:port.")
	flag.StringVarP(&snmpdAddr, "snmpd-addr", "", lib.SNMP_DAEMON_ADDR, "snmpd address:port.")
	flag.DurationVarP(&a.DumpTableTime, "dump-table-time", "", DUMP_TIME_DEFAULT, "dump-table interval")
	flag.StringVarP(&a.DumpTableFile, "dump-table-file", "", DUMP_FILE_DEFAULT, "dump-table filename.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail message.")
	flag.Parse()

	if a.ListenAddr, err = net.ResolveUDPAddr("udp", listenAddr); err != nil {
		return err
	}
	if a.SnmpdAddr, err = net.ResolveUDPAddr("udp", snmpdAddr); err != nil {
		return err
	}

	return nil
}

func NewArgs() (*Args, error) {
	args := &Args{}
	if err := args.Parse(); err != nil {
		return nil, err
	}
	return args, nil
}

func dumpTables(t *Tables, path string, interval time.Duration) {
	if interval == 0 {
		log.Infof("dump table disabled.")
		return
	}
	if interval < DUMP_TIME_MIN {
		interval = DUMP_TIME_MIN
	}

	f, err := os.Create(path)
	if err != nil {
		log.Errorf("Open dump-file error. %s", err)
		return
	}
	defer f.Close()

	tiker := time.NewTicker(interval)
	for {
		f.Truncate(0)
		f.Seek(0, 0)
		t.WriteTo(f)
		f.Sync()
		<-tiker.C
	}
}

func printArgs(a *Args) {
	log.Debugf("config: '%s'", a.Config)
	log.Debugf("Table : '%s'", a.Table)
	log.Debugf("IfOID : '%s'", a.IfOID)
	log.Debugf("Listen: '%s'", a.ListenAddr)
	log.Debugf("Snmpd : '%s'", a.SnmpdAddr)
	log.Debugf("Dump  : %s '%s", a.DumpTableTime, a.DumpTableFile)
}

func main() {
	args, err := NewArgs()
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	printArgs(args)

	configs, err := ReadConfig(args.Config)
	if err != nil {
		log.Errorf("ReadConfig error. %s", err)
		os.Exit(1)
	}

	log.Debugf("%v", configs)

	config, ok := configs[args.Table]
	if !ok {
		log.Errorf("Config[%s] not found.", args.Table)
		os.Exit(1)
	}

	log.Debugf("%v", config)

	s, err := NewProxyServer(args.ListenAddr, args.SnmpdAddr, args.IfOID)
	if err != nil {
		log.Errorf("NewUDPServer error. %s", err)
		os.Exit(1)
	}

	for _, c := range config.OidMap {
		e := NewOidMapEntry(c.Name, c.Oid, c.Local)
		log.Debugf("OidMap %s", e)
		s.OidMapTable().Add(e)
	}

	s.IfMapTable().Add(NewIfMapOidMap(config.IfMap.GetOidMap()))
	s.IfMapTable().Add(NewIfMapShift(config.IfMap.GetShift()))

	for ifname, portId := range config.Trap2Map {
		e := NewTrapMapEntry(ifname, -1, portId)
		log.Debugf("TrapMap %s", e)
		s.TrapMapTable().add(e)
	}

	for _, v := range config.Trap2Sink {
		log.Debugf("TrapSink %s", v)
		s.TrapSinkTable().Add(v.Addr)
	}

	go dumpTables(s.Tables, args.DumpTableFile, args.DumpTableTime)
	s.Start()
}
