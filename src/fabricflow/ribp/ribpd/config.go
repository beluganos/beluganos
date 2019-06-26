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
	"github.com/BurntSushi/toml"
)

type NodeConfig struct {
	NId       uint8  `toml:"nid"`
	ReId      string `toml:"reid"`
	NIdIfname string `toml:"nid_from_ifaddr"`
	DupIfname bool   `toml:"allow_duplicate_ifname"`
}

type RibpConfig struct {
	Api      string   `toml:"api"`
	Interval int      `tomi:"interval"`
	Excludes []string `toml:"exclude_ifaces"`
}

type Config struct {
	Node NodeConfig `toml:"node"`
	Ribp RibpConfig `toml:"ribp"`
}

func ReadConfig(path string, cfg *Config) error {
	cfg.Node.DupIfname = true // default value
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return err
	}

	if cfg.Ribp.Api == "" {
		cfg.Ribp.Api = "127.0.0.1:50091"
	}

	return nil
}
