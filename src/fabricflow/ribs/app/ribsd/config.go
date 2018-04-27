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
	"fmt"
	"github.com/BurntSushi/toml"
)

const GOBGPD_GRPC_PORT = 50051
const VIRT_NEXTHOPS = "127.0.0.1/32"

type NodeConfig struct {
	NId       uint8  `toml:"nid"`
	Label     uint32 `toml:"label"`
	NIdIfname string `toml:"nid_from_ifaddr"`
}

type NLAConfig struct {
	Api string `toml:"api"`
}

type BgpConfig struct {
	Addr string `toml:"addr"`
	Port uint16 `toml:"port"`
}

func (c *BgpConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

type VrfConfig struct {
	Iface string `toml:"iface"`
	Rt    string `toml:"rt"`
	Rd    string `toml:"rd"`
}

type NHConfig struct {
	Mode string `toml:"mode"`
	Args string `toml:"args"`
}

type RibsConfig struct {
	Disable  bool      `toml:"disable"`
	Core     string    `toml:"core"`
	Api      string    `toml:"api"`
	SyncTime int64     `toml:"resync"`
	Nexthops NHConfig  `toml:"nexthops"`
	Bgp      BgpConfig `toml:"bgpd"`
	Vrf      VrfConfig `toml:"vrf"`
}

type Config struct {
	Node NodeConfig `toml:"node"`
	NLA  NLAConfig  `toml:"nla"`
	Ribs RibsConfig `toml:"ribs"`
}

func (c *Config) GetBgpAddr() string {
	return c.Ribs.Bgp.GetAddr()
}

func (c *Config) VrfLabel() uint32 {
	return uint32(c.Node.NId) + c.Node.Label
}

func ReadConfig(path string, cfg *Config) error {
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return err
	}

	if cfg.Ribs.Bgp.Port == 0 {
		cfg.Ribs.Bgp.Port = GOBGPD_GRPC_PORT
	}

	if cfg.Ribs.Nexthops.Args == "" {
		cfg.Ribs.Nexthops.Args = VIRT_NEXTHOPS
	}
	return nil
}
