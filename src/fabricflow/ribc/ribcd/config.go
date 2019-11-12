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
	"fabricflow/ribc/ribctl"
	"fmt"

	"github.com/BurntSushi/toml"
)

const LABEL_BASE_DEFAULT = 0xfff00

type NodeConfig struct {
	NId       uint8  `toml:"nid"`
	Label     uint32 `toml:"label"`
	ReId      string `toml:"reid"`
	NIdIfname string `toml:"nid_from_ifaddr"`
	DupIfname bool   `toml:"allow_duplicate_ifname"`
}

func (c *NodeConfig) String() string {
	return fmt.Sprintf("nid:%d label:%d reid:'%s' dup-if:%t nid-if:'%s'",
		c.NId, c.Label, c.ReId, c.DupIfname, c.NIdIfname)
}

type NlaConfig struct {
	Api string `toml:"api"`
}

func (c *NlaConfig) String() string {
	return fmt.Sprintf("api:'%s'", c.Api)
}

type RibcConfig struct {
	Fibc     string `toml:"fibc"`
	FibcType string `toml:"fibc_type"`
	Disable  bool   `toml:"disable"`
}

func (c *RibcConfig) String() string {
	return fmt.Sprintf("fibc:'%s' type:'%s' disable:%t", c.Fibc, c.FibcType, c.Disable)
}

func (c *RibcConfig) GetFibcType() string {
	if len(c.FibcType) == 0 {
		return ribctl.FIBCTypeDefault
	}
	return c.FibcType
}

type Config struct {
	Node NodeConfig `toml:"node"`
	NLA  NlaConfig  `toml:"nla"`
	Ribc RibcConfig `toml:"ribc"`
}

func (c *Config) String() string {
	return fmt.Sprintf("node:{%s} nla:{%s} ribc:{%s}", &c.Node, &c.NLA, &c.Ribc)
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	config.Node.DupIfname = true // default value
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return nil, err
	}

	if config.Node.Label < 17 {
		config.Node.Label = LABEL_BASE_DEFAULT
	}

	return config, nil
}
