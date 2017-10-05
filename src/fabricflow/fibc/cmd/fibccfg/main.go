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
	"bytes"
	"flag"
	"fmt"
	"net/http"
)

func showRouters(routers []Router) {
	for _, router := range routers {
		fmt.Printf("%s\n", router.String())
		for _, port := range router.Ports {
			fmt.Printf("  %s\n", port.String())
		}
	}
}

func showDpaths(dpaths []Datapath) {
	for _, dpath := range dpaths {
		fmt.Printf("%s\n", dpath.String())
	}
}

func modRouters(url string, routers []Router) {
	for _, router := range routers {
		b, err := router.ToJSON()
		if err != nil {
			fmt.Printf("ToJSON error. %s", err)
			return
		}

		if _, err = http.Post(url, "application/json", bytes.NewReader(b)); err != nil {
			fmt.Printf("http POST error. %s", err)
			return
		}
	}
}

func modDpaths(url string, dpaths []Datapath) {
	for _, dpath := range dpaths {
		b, err := dpath.ToJSON()
		if err != nil {
			fmt.Printf("ToJSON error. %s", err)
			return
		}

		if _, err := http.Post(url, "application/json", bytes.NewReader(b)); err != nil {
			fmt.Printf("http POST error. %s", err)
			return
		}
	}
}

func main() {
	var addr string
	var cmd string
	var path string
	var table string
	flag.StringVar(&addr, "a", "127.0.0.1:8080", "ryu addr")
	flag.StringVar(&cmd, "c", "check", "command(add/delete/check)")
	flag.StringVar(&table, "t", "port", "table(port/dp/re)")
	flag.StringVar(&path, "f", "fibc.yml", "filename.")
	flag.Parse()

	var cfg Config
	if err := ReadConfig(path, &cfg); err != nil {
		fmt.Printf("ReadConfig error. %s", err)
		return
	}

	switch cmd {
	case "add", "delete":
		url := fmt.Sprintf("http://%s/fib/portmap/%s/%s", addr, table, cmd)
		if table == "dp" {
			modDpaths(url, cfg.Datapaths)
		} else {
			modRouters(url, cfg.Routers)
		}
	default:
		showRouters(cfg.Routers)
		showDpaths(cfg.Datapaths)
	}
}
