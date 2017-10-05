// coding: utf-8 -*-

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
	"flag"
	log "github.com/sirupsen/logrus"
	"gonla/nlactl"
	"gonla/nladbm"
	"gonla/nlasvc"
	"os"
)

func readConfig(config *Config) error {
	var path string
	flag.StringVar(&path, "config", "/etc/nlad/nlad.toml", "file name")
	flag.Parse()

	if err := ReadConfig(path, config); err != nil {
		return err
	}

	log.Infof("config: %v", config)
	return nil
}

func initlog(c *Config) {
	log.SetLevel(log.Level(c.Log.Level))
	if len(c.Log.Output) > 0 {
		f, err := os.OpenFile(c.Log.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Warnf("os.OpenFile error. %s", err)
			return
		}

		log.SetOutput(f)
	}
}

func services(c *Config) []nlactl.NLAService {
	if c.IsMaster() {
		nlaapi := nlasvc.NewNLAApiService(c.NLA.Api)
		return []nlactl.NLAService{
			nlasvc.NewNLALogService(c.Log.Dump),
			nlasvc.NewNLAMasterService(nlaapi),
			nlasvc.NewNLACoreApiService(c.NLA.Core),
			nlasvc.NewNLANetlinkService(),
		}

	} else {
		return []nlactl.NLAService{
			nlasvc.NewNLALogService(c.Log.Dump),
			nlasvc.NewNLASlaveService(c.NLA.Core),
			nlasvc.NewNLANetlinkService(),
		}
	}
}

func main() {

	var cfg Config
	if err := readConfig(&cfg); err != nil {
		log.Errorf("readConfig error. %s", err)
		return
	}

	initlog(&cfg)

	done := make(chan struct{})
	svcs := services(&cfg)
	m := nlactl.NewNLAManager(cfg.Node.NId, done, svcs...)
	s := nlactl.NewNLAServer(cfg.Node.NId, m.Chans.NlMsg, done)
	nladbm.Create()
	m.Start()
	s.Start()

	<-done
}
