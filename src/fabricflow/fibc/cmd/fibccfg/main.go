// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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
	"fabricflow/fibc/cmd/lib"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func showRouters(routers []*fibccmd.Router) {
	for _, router := range routers {
		fmt.Printf("%s\n", router.String())
		for _, port := range router.Ports {
			fmt.Printf("  %s\n", port.String())
		}
	}
}

func showDpaths(dpaths []*fibccmd.Datapath) {
	for _, dpath := range dpaths {
		fmt.Printf("%s\n", dpath.String())
	}
}

func main() {
	var addr string
	var cmd string
	var path string
	flag.StringVar(&addr, "a", "127.0.0.1:8080", "ryu addr")
	flag.StringVar(&cmd, "c", "check", "command(add/del/check)")
	flag.StringVar(&path, "f", "fibc.yml", "filename.")
	flag.Parse()

	var cfg fibccmd.Config
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("ReadConfig error. %s", err)
		return
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		fmt.Printf("ReadConfig error. %s", err)
		return
	}

	c := fibccmd.NewConfigClient(addr)
	switch cmd {
	case "add":
		if err := c.Add(&cfg); err != nil {
			fmt.Printf("Add error. %s\n", err)
		}

	case "del":
		if err := c.Del(&cfg); err != nil {
			fmt.Printf("Add error. %s\n", err)
		}

	default:
		showRouters(cfg.Routers)
		showDpaths(cfg.Datapaths)
	}
}
