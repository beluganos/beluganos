// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package monitor

import "github.com/spf13/viper"

type IfaceStatsCounter struct {
	Label string `mapstructure:"label"`
	Name  string `mapstructure:"name"`
}

type IfaceStatsConfig struct {
	Label    string               `mapstructure:"label"`
	Counters []*IfaceStatsCounter `mapstructure:"counters"`
}

type IfaceDPathConfig struct {
	DpID   uint64   `mapstructure:"dpid"`
	Ifaces []string `mapstructute:"ifaces"`
}

func NewIfaceDPathConfig() *IfaceDPathConfig {
	return &IfaceDPathConfig{
		DpID:   0,
		Ifaces: []string{},
	}
}

type IfaceConfig struct {
	Stats  map[string][]*IfaceStatsConfig `mapstructure:"ifstats"`
	DPaths map[string]*IfaceDPathConfig   `mapstructute:"dpaths"`

	viper *viper.Viper
}

func NewIfaceConfig() *IfaceConfig {
	return &IfaceConfig{
		Stats:  map[string][]*IfaceStatsConfig{},
		DPaths: map[string]*IfaceDPathConfig{},
		viper:  viper.New(),
	}
}

func (c *IfaceConfig) SetConfig(filePath, fileType string) *IfaceConfig {
	c.viper.SetConfigFile(filePath)
	c.viper.SetConfigType(fileType)

	return c
}

func (c *IfaceConfig) Read() error {
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	if err := c.viper.UnmarshalExact(c); err != nil {
		return err
	}

	return nil
}

func (c *IfaceConfig) StatsConfig(name string) ([]*IfaceStatsConfig, bool) {
	sc, ok := c.Stats[name]
	return sc, ok
}

func (c *IfaceConfig) DPathConfig(name string) (*IfaceDPathConfig, bool) {
	pc, ok := c.DPaths[name]
	return pc, ok
}
