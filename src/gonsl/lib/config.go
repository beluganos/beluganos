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
	configDefaultAddr        = "127.0.0.1"
	configDefaultPort uint16 = 50051
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
// DpConfig is part of gonsl config.
//
type DpConfig struct {
	DpID    uint64      `mapstructure:"dpid"`
	Addr    string      `mapstructure:"addr"`
	Port    uint16      `mapstructure:"port"`
	Unit    int         `mapstructure:"unit"`
	OpenNSL *ONSLConfig `mapstructure:"opennsl"`
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
// Config is gonsl config file.
//
type Config struct {
	Dpaths map[string]*DpConfig `mapstructure:"dpaths"`
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
