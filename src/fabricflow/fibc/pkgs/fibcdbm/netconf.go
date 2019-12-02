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

package fibcdbm

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type NetconfPortConfig struct {
	HwAddr string `mapstructure:"hw_addr"`
	Name   string `mapstructure:"name"`
	PortNo uint32 `mapstructure:"port_no"`
}

func (c *NetconfPortConfig) String() string {
	return fmt.Sprintf("'%s' %d '%s'", c.Name, c.PortNo, c.HwAddr)
}

func NewNetconfPortConfig(name, hwaddr string, portNo uint32) *NetconfPortConfig {
	return &NetconfPortConfig{
		HwAddr: hwaddr,
		Name:   name,
		PortNo: portNo,
	}
}

type NetconfDpConfig struct {
	Ports []*NetconfPortConfig `mapstructure:"ports"`
}

func NewNetconfDpConfig() *NetconfDpConfig {
	return &NetconfDpConfig{
		Ports: []*NetconfPortConfig{},
	}
}

func (c *NetconfDpConfig) AddPort(port *NetconfPortConfig) {
	c.Ports = append(c.Ports, port)
}

type NetconfConfig struct {
	Dps map[uint64]*NetconfDpConfig `mapstructure:"dps"`

	v *viper.Viper
}

func NewNetconfConfig() *NetconfConfig {
	return &NetconfConfig{
		Dps: map[uint64]*NetconfDpConfig{},
		v:   viper.New(),
	}
}

func (c *NetconfConfig) createIfNotExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	return nil
}

func (c *NetconfConfig) SetConfig(configPath, configType string) error {
	if err := c.createIfNotExist(configPath); err != nil {
		return err
	}

	c.v.SetConfigFile(configPath)
	c.v.SetConfigType(configType)

	return nil
}

func (c *NetconfConfig) Load() error {
	if err := c.v.ReadInConfig(); err != nil {
		return err
	}

	return c.v.Unmarshal(c)
}

func (c *NetconfConfig) Save() error {
	c.v.Set("dps", c.Dps)
	return c.v.WriteConfig()
}
