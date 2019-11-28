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
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type Args struct {
	Addr      string
	Prefix    string
	Family    string
	LocalLink string
	LocalNW4  net.IPNet
	LocalNW6  net.IPNet
	LocalWait uint
	DumpTbl   time.Duration
	Verbose   bool

	TunType6   uint16
	TunForce   uint16
	TunDefault uint16

	APIAddr string
}

const (
	ARGS_DUMP_TABLE_DEFAULT = 3 * time.Second
	ARGS_LOCAL_LINK_DEFAULT = "lo"
	ARGS_LOCAL_NW4_DEFAULT  = "127.0.0.1/32"
	ARGS_LOCAL_NW6_DEFAULT  = "::1/128"
	ARGS_LOCAL_ADDR_WAIT    = 30
	ARGS_TUNTYPE6_DEFAULT   = 14 // IPv6 Tunnel
	ARGS_APIADDR_DEFAULT    = "localhost:50099"
)

func (a *Args) Parse() error {
	_, local4, _ := net.ParseCIDR(ARGS_LOCAL_NW4_DEFAULT)
	_, local6, _ := net.ParseCIDR(ARGS_LOCAL_NW6_DEFAULT)

	flag.StringVarP(&a.Addr, "gobgpd-api", "a", "localhost:50051", "GoBGP api address.")
	flag.StringVarP(&a.Prefix, "tunnel-prefix", "", "tun", "prefix of tunnel device name.")
	flag.StringVarP(&a.Family, "route-family", "", "ipv4-unicast", "route family.")
	flag.StringVarP(&a.LocalLink, "tunnel-local-if", "", ARGS_LOCAL_LINK_DEFAULT, "tunnel local ifname.")
	flag.IPNetVarP(&a.LocalNW4, "tunnel-local-nw4", "", *local4, "tunnel local network(IPv4).")
	flag.IPNetVarP(&a.LocalNW6, "tunnel-local-nw6", "", *local6, "tunnel local network(IPv6).")
	flag.UintVarP(&a.LocalWait, "tunnel-local-wait", "", ARGS_LOCAL_ADDR_WAIT, "tunnel local addr wait sec.")
	flag.DurationVarP(&a.DumpTbl, "dump-table", "", ARGS_DUMP_TABLE_DEFAULT, "Dump table interval.")
	flag.Uint16VarP(&a.TunType6, "tunnel-type-ipv6", "", ARGS_TUNTYPE6_DEFAULT, "tunnel type value(IPv6).")
	flag.Uint16VarP(&a.TunForce, "tunnel-type-force", "", 0, "tunnel type value(all route).")
	flag.Uint16VarP(&a.TunDefault, "tunnel-type-default", "", 0, "tunnel type value(no encap).")

	flag.StringVarP(&a.APIAddr, "api-addr", "", ARGS_APIADDR_DEFAULT, "ribt api listen address.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail messages.")
	flag.Parse()

	return nil
}

func NewArgs() (args *Args, err error) {
	args = &Args{}
	err = args.Parse()
	return
}

func (a *Args) Dump() {
	log.Infof("gobgpd-api       : %s", a.Addr)
	log.Infof("route-family     : %s", a.Family)
	log.Infof("tunnel-prefix    : %s", a.Prefix)
	log.Infof("tunnel-local-if  : %s", a.LocalLink)
	log.Infof("tunnel-local-nw4 : %s", &a.LocalNW4)
	log.Infof("tunnel-local-nw6 : %s", &a.LocalNW6)
	log.Infof("tunnel-local-try : %d", &a.LocalWait)
	log.Infof("tunnel-type-ipv6 : %d", a.TunType6)
	log.Infof("tunnel-type-force: %d", a.TunForce)
	log.Infof("tunnel-type-deflt: %d", a.TunDefault)
	log.Infof("api listen addr  : %s", a.APIAddr)
	log.Infof("dump-table       : %s", a.DumpTbl)
}

func (a *Args) TunLocalAddrs() (net.IP, net.IP, error) {
	local4, err := LinkAddr(a.LocalLink, &a.LocalNW4)
	if err != nil {
		log.Debugf("Tunnel local address(ipv4) not found. %s", err)
		return nil, nil, err
	}

	local6, err := LinkAddr(a.LocalLink, &a.LocalNW6)
	if err != nil {
		log.Debugf("Tunnel local address(ipv6) not found. %s", err)
		return nil, nil, err
	}

	log.Infof("Tunnel local address(IPv4): %s", local4.IP)
	log.Infof("Tunnel local address(IPv6): %s", local6.IP)

	return local4.IP, local6.IP, nil
}

func (a *Args) WaitForTunLocalAddrs() (net.IP, net.IP, error) {
	var retry uint
	for {
		if local4, local6, err := a.TunLocalAddrs(); err == nil {
			return local4, local6, nil
		}

		retry++
		if retry > a.LocalWait {
			return nil, nil, fmt.Errorf("Local address not detected.")
		}

		log.Debugf("Waiting for local address ...")
		time.Sleep(time.Second)
	}
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

	args.Dump()

	local4, local6, err := args.WaitForTunLocalAddrs()
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	server, err := NewServer(args.Addr, args.Prefix, args.Family, local4, local6)
	if err != nil {
		log.Errorf("NewServer error. %s", err)
		os.Exit(1)
	}

	server.SetTunType6(args.TunType6)
	server.SetTunForce(args.TunForce)
	server.SetTunTypeDefault(args.TunDefault)
	server.SetAPIAddr(args.APIAddr)

	done := make(chan struct{})

	if err := server.Start(done); err != nil {
		log.Errorf("LinkLoad error. %s", err)
		os.Exit(1)
	}

	if args.DumpTbl > 0 {
		go dumpTable(server.Tables, args.DumpTbl, done)
	}

	<-done
}
