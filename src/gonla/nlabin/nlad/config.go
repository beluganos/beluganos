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

const DEFAULT_MNG_IF = "eth0"
const AUTO_NID uint8 = 255

type NodeConfig struct {
	NId   uint8  `toml:"nid"`
	MngIF string `toml:"mng-if"`
}

type NLAConfig struct {
	Core string `toml:"core"`
	Api  string `toml:"api"`
}

type LogConfig struct {
	Level  uint8  `toml:"level"`
	Dump   uint32 `toml:"dump"`
	Output string `toml:"output"`
}

type Config struct {
	Node NodeConfig `toml:"node"`
	NLA  NLAConfig  `toml:"nla"`
	Log  LogConfig  `toml:"log"`
}

func (c *Config) IsMaster() bool {
	return len(c.NLA.Api) != 0
}

func ReadConfig(path string, config *Config) error {
	if _, err := toml.DecodeFile(path, config); err != nil {
		return err
	}

	if len(config.Node.MngIF) == 0 {
		config.Node.MngIF = DEFAULT_MNG_IF
	}

	return nil
}
