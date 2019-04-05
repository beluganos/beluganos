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

type Args struct {
	IfCommunity string
	IfOID       string
	SnmpdAddr   *net.UDPAddr
	ResendTime  time.Duration
	SkipIfnames []string
	Verbose     bool
}

func (a *Args) Parse() error {
	var (
		snmpdAddr string
		err       error
	)
	skipIfnames := []string{"lo", "eth0"}

	flag.StringVarP(&a.IfCommunity, "if-notify-community", "", lib.SNMP_COMMUNITY, "iface notify community.")
	flag.StringVarP(&a.IfOID, "if-noify-oid", "", lib.SNMP_OID_IFACES, "iface noify OID.")
	flag.StringVarP(&snmpdAddr, "snmpd-addr", "", lib.SNMP_LISTEN_ADDR, "snmpd address:port.")
	flag.DurationVarP(&a.ResendTime, "trap-resend", "", TRAP_RESEND_INTERVAL, "Trap resend interval time.")
	flag.StringSliceVarP(&a.SkipIfnames, "skip", "", skipIfnames, "Skip ifnames.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail message.")
	flag.Parse()

	if a.SnmpdAddr, err = net.ResolveUDPAddr("udp", snmpdAddr); err != nil {
		return err
	}

	return nil
}

func printArgs(a *Args) {
	log.Infof("IfCom     : '%s'", a.IfCommunity)
	log.Infof("IfOID     : '%s'", a.IfOID)
	log.Infof("Snmpd     : '%s'", a.SnmpdAddr)
	log.Infof("TrapResend: %s", a.ResendTime)
	log.Infof("SkipIfaces: %v", a.SkipIfnames)
}

func NewArgs() (*Args, error) {
	args := &Args{}
	if err := args.Parse(); err != nil {
		return nil, err
	}
	return args, nil
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

	ifoid := lib.ParseOID(args.IfOID)

	s := NewServer(args.SnmpdAddr, args.IfCommunity, ifoid)
	s.RegisterSkipIfname(args.SkipIfnames...)
	s.SetResendInterval(args.ResendTime)
	if err := s.Start(); err != nil {
		log.Errorf("Server Start error. %s", err)
		os.Exit(1)
	}

	done := make(chan struct{})
	<-done
}
