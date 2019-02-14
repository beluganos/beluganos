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

package gonslib

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	configDefaultAddr           = "127.0.0.1"
	configDefaultPort    uint16 = 50051
	configDefaultBaseVid uint16 = 1
)

//
// ONSLConfig is opennsl config file.
//
type ONSLConfig struct {
	Config     string `mapstructure:"cfg_fname"`
	ConfigPost string `mapstructure:"cfg_post_fname"`
	Flags      uint   `mapstructure:"flags"`
	RmConfig   string `mapstructure:"rmcfg_fname"`
	WbConfig   string `mapstructure:"wb_fname"`
}

//
// BlockBcastRangeConfig is port entry of BlockBcastConfig.
//
type BlockBcastRangeConfig struct {
	Min     int    `mapstructure:"min"`
	Max     int    `mapstructure:"max"`
	BaseVID uint16 `mapstructure:"base_vid"`
}

func (c *BlockBcastRangeConfig) String() string {
	return fmt.Sprintf("min:%d max:%d base_vid:%d", c.Min, c.Max, c.BaseVID)
}

func (c *BlockBcastRangeConfig) Block() bool {
	return (c.Min > 0) && (c.Max >= c.Min)
}

func (c *BlockBcastRangeConfig) GetBaseVID() uint16 {
	if c.BaseVID == 0 {
		return configDefaultBaseVid
	}
	return c.BaseVID
}

//
// BlockBcastPortConfig is port entry of BlockBcastConfig.
//
type BlockBcastPortConfig struct {
	Port int    `mapstructure:"port"`
	Vid  uint16 `mapstructure:"vid"`
	PVid uint16 `mapstructure:"pvid"`
}

func (c *BlockBcastPortConfig) String() string {
	return fmt.Sprintf("port:%d vid:%d to %d", c.Port, c.Vid, c.PVid)
}

//
// VlanPortConfig is vlan config that port belongs to.
//
type BlockBcastConfig struct {
	Range BlockBcastRangeConfig   `mapstructure:"range"`
	Ports []*BlockBcastPortConfig `mapstructure:"ports"`
}

func (c *BlockBcastConfig) String() string {
	return fmt.Sprintf("%s", &c.Range)
}

func (c *BlockBcastConfig) Block() bool {
	return c.Range.Block() || (len(c.Ports) > 0)
}

//
// DpConfig is part of gonsl config.
//
type DpConfig struct {
	DpID       uint64           `mapstructure:"dpid"`
	Addr       string           `mapstructure:"addr"`
	Port       uint16           `mapstructure:"port"`
	Unit       int              `mapstructure:"unit"`
	BlockBcast BlockBcastConfig `mapstructure:"block_bcast"`
	OpenNSL    *ONSLConfig      `mapstructure:"opennsl"`
}

//
// String retuns config information.
//
func (c *DpConfig) String() string {
	return fmt.Sprintf("DpID=%d, addr='%s', port=%d", c.DpID, c.GetAddr(), c.GetPort())
}

//
// GetAddr returns address of server.
//
func (c *DpConfig) GetAddr() string {
	if len(c.Addr) == 0 {
		return configDefaultAddr
	}
	return c.Addr
}

//
// GetPort returns port number of server.
//
func (c *DpConfig) GetPort() uint16 {
	if c.Port == 0 {
		return configDefaultPort
	}
	return c.Port
}

//
// GetHost returns ipaddr:port
//
func (c *DpConfig) GetHost() string {
	return fmt.Sprintf("%s:%d", c.GetAddr(), c.GetPort())
}

//
// LogConfig is logging config.
//
type LogConfig struct {
	RxDetail bool `mapstructure:"rx_detail"`
	RxDump   bool `mapstructure:"rx_dump"`
}

//
// String retuns config information.
//
func (c *LogConfig) String() string {
	return fmt.Sprintf("rx_dump=%t", c.RxDump)
}

//
// Config is gonsl config file.
//
type Config struct {
	Dpaths  map[string]*DpConfig `mapstructure:"dpaths"`
	Logging LogConfig            `mapstructure:"logging"`
}

//
// NewConfig returns new instance.
//
func NewConfig() *Config {
	return &Config{
		Dpaths: map[string]*DpConfig{},
	}
}

//
// ReadFile reads config file.
//
func (c *Config) ReadFile(args *Args) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(args.ConfigFile)
	v.SetConfigType(args.ConfigType)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.UnmarshalExact(c); err != nil {
		return nil, err
	}

	return v, nil
}

//
// GetDpConfig returns dp cofig.
//
func (c *Config) GetDpConfig(name string) (*DpConfig, error) {
	if dpcfg, ok := c.Dpaths[name]; ok {
		return dpcfg, nil
	}
	return nil, fmt.Errorf("%s not found in dp config", name)
}
