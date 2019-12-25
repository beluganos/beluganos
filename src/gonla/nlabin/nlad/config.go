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
	"net"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	DEFAULT_MNG_IF              = "eth0"
	AUTO_NID             uint8  = 255
	BRVLAN_UPDATE_SECOND uint32 = 1800
	BRVLAN_CHAN_SIZE     int    = 4096 * 4
)

type NodeConfig struct {
	NId   uint8  `toml:"nid"`
	MngIF string `toml:"mng-if"`
}

func (c *NodeConfig) Adjust() {
	if len(c.MngIF) == 0 {
		c.MngIF = DEFAULT_MNG_IF
	}
}

func (c *NodeConfig) String() string {
	return fmt.Sprintf("NId:%d, MngIF:'%s'", c.NId, c.MngIF)
}

type IptunRemote struct {
	*net.IPNet
}

func (c *IptunRemote) UnmarshalText(text []byte) error {
	_, ipnet, err := net.ParseCIDR(string(text))
	if err != nil {
		return err
	}

	c.IPNet = ipnet
	return nil
}

type IptunConfig struct {
	NId     uint8         `toml:"nid"`
	Remotes []IptunRemote `toml:"remotes"`
}

func (c *IptunConfig) String() string {
	return fmt.Sprintf("nid:%d %v", c.NId, c.Remotes)
}

type BridgeVlanConfig struct {
	UpdateSec uint32 `toml:"update_sec"`
	ChanSize  int    `toml:"chan_size"`
}

func (c *BridgeVlanConfig) UpdateTime() time.Duration {
	return time.Duration(c.UpdateSec) * time.Second
}

func (c *BridgeVlanConfig) String() string {
	return fmt.Sprintf("update_sec:%d, chan_size:%d", c.UpdateSec, c.ChanSize)
}

type NLAConfig struct {
	Core            string `toml:"core"`
	Api             string `toml:"api"`
	RecvChanSize    int    `toml:"recv_chan_size"`
	RecvSockBufSize int    `toml:"recv_sock_buf"`

	Iptun      []IptunConfig    `toml:"iptun"`
	BridgeVlan BridgeVlanConfig `toml:"bridge_vlan"`
}

func (c *NLAConfig) Adjust() {
	if c.RecvChanSize <= 0 {
		c.RecvChanSize = 65536
	}

	if c.RecvSockBufSize <= 0 {
		c.RecvSockBufSize = 1024 * 1024
	}

	if c.BridgeVlan.UpdateSec == 0 {
		c.BridgeVlan.UpdateSec = BRVLAN_UPDATE_SECOND
	}

	if c.BridgeVlan.ChanSize == 0 {
		c.BridgeVlan.ChanSize = BRVLAN_CHAN_SIZE
	}
}

func (c *NLAConfig) String() string {
	return fmt.Sprintf("Core:'%s', Api:'%s', RecvChan:%d, RecvSock:%d, iptun:{%v}, brvlan:{%s}", c.Core, c.Api, c.RecvChanSize, c.RecvSockBufSize, &c.Iptun, &c.BridgeVlan)
}

type LogConfig struct {
	Level  uint8  `toml:"level"`
	Dump   uint32 `toml:"dump"`
	Output string `toml:"output"`
}

func (c *LogConfig) String() string {
	return fmt.Sprintf("Level:%d, Dump:%d, Output:'%s'", c.Level, c.Dump, c.Output)
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

	config.Node.Adjust()
	config.NLA.Adjust()

	return nil
}
