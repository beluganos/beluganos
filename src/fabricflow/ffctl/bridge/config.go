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

package bridge

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type BrVlanBridgeConfig struct {
	VlanFiltering int      `mapstructure:"vlan_filtering"`
	Interfaces    []string `mapstructure:"interfaces"`
}

func NewBrVlanBridgeConfig() *BrVlanBridgeConfig {
	return &BrVlanBridgeConfig{
		VlanFiltering: 1,
		Interfaces:    []string{},
	}
}

func (c *BrVlanBridgeConfig) String() string {
	return fmt.Sprintf("vlan_flt:%d %v", c.VlanFiltering, c.Interfaces)
}

type BrVlanVlanConfig struct {
	Id  uint16   `mapstructure:"id" yaml:",omitempty"`
	Ids []uint16 `mapstructure:"ids" yaml:",omitempty"`
}

func NewBrVlanVlanConfig() *BrVlanVlanConfig {
	return &BrVlanVlanConfig{
		Ids: []uint16{},
	}
}

func (c *BrVlanVlanConfig) String() string {
	if len(c.Ids) > 0 {
		return fmt.Sprintf("%v", c.Ids)
	}

	return fmt.Sprintf("%d", c.Id)
}

type BrVlanNetworkConfig struct {
	Vlans   map[string]*BrVlanVlanConfig   `mapstructure:"vlans" yaml:",omitempty"`
	Bridges map[string]*BrVlanBridgeConfig `mapstructure:"bridges" yaml:",omitempty"`
}

func NewBrVlanNetworkConfig() *BrVlanNetworkConfig {
	return &BrVlanNetworkConfig{
		Vlans:   map[string]*BrVlanVlanConfig{},
		Bridges: map[string]*BrVlanBridgeConfig{},
	}
}

type BrVlanConfig struct {
	Network *BrVlanNetworkConfig `mapstructure:"network"`
	viper   *viper.Viper
}

func NewBrVlanConfig() *BrVlanConfig {
	return &BrVlanConfig{
		Network: NewBrVlanNetworkConfig(),
		viper:   viper.New(),
	}
}

func (c *BrVlanConfig) SetConfigFile(path, format string) *BrVlanConfig {
	c.viper.SetConfigFile(path)
	c.viper.SetConfigType(format)
	return c
}

func (c *BrVlanConfig) Load() error {
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return c.viper.Unmarshal(c)
}

func (c *BrVlanConfig) Yaml() ([]byte, error) {
	return yaml.Marshal(c)
}
