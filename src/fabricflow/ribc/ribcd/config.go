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

const LABEL_BASE_DEFAULT = 0xfff00

type NodeConfig struct {
	NId   uint8  `toml:"nid"`
	Label uint32 `toml:"label"`
	ReId  string `toml:"reid"`
}

type NlaConfig struct {
	Api string `toml:"api"`
}

type RibcConfig struct {
	Fibc    string `toml:"fibc"`
	Disable bool   `toml:"disable"`
}

type Config struct {
	Node NodeConfig `toml:"node"`
	NLA  NlaConfig  `toml:"nla"`
	Ribc RibcConfig `toml:"ribc"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return nil, err
	}

	if config.Node.Label < 17 {
		config.Node.Label = LABEL_BASE_DEFAULT
	}

	return config, nil
}
