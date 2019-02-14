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
	api "gonsl/api"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

//
// Args is argument of onsldump.
//
type Args struct {
	ServerAddr string
	Verbose    bool
}

//
// Init parse arguments.
//
func (a *Args) Init() {
	flag.StringVarP(&a.ServerAddr, "addr", "a", "localhost:50061", "Server address.")
	flag.BoolVarP(&a.Verbose, "verboce", "v", false, "show detail messages.")
	flag.Parse()
}

func main() {
	args := Args{}
	args.Init()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	conn, err := grpc.Dial(args.ServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("Connect error. %s", err)
		os.Exit(1)
	}

	defer conn.Close()

	client := api.NewGoNSLApiClient(conn)
	dumpFieldEntries(client)
	dumpVlans(client)
	dumpL2Addrs(client)
	dumpL3Ifaces(client)
	dumpL3Egresses(client)
	dumpL3Hosts(client)
	dumpL3Routes(client)
	dumpIDMapEntries(client)
}
