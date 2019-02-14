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

package main

import (
	"fmt"
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v2"
)

//
// ConfigOidMap is config(/oidmap/<name>)
//
type ConfigOidMap struct {
	Name  string `yaml:"name"`
	Oid   string `yaml:"oid"`
	Local string `yaml:"local"`
}

func (c *ConfigOidMap) String() string {
	return fmt.Sprintf("%s oid:'%s', Local:'%s'", c.Name, c.Oid, c.Local)
}

type ConfigIfMapEntry struct {
	Min uint `yaml:"min"`
	Max uint `yaml:"max"`
}

func (c *ConfigIfMapEntry) String() string {
	return fmt.Sprintf("min:%d, max:%d", c.Min, c.Max)
}

//
// ConfigIfMap is config(/ifmap/<name>)
//
type ConfigIfMap struct {
	OidMap *ConfigIfMapEntry `yaml:"oidmap"`
	Shift  *ConfigIfMapEntry `yaml:"shift"`
}

func (c *ConfigIfMap) String() string {
	return fmt.Sprintf("oidmap:%v, shift:%v", c.OidMap, c.Shift)
}

func (c *ConfigIfMap) GetOidMap() (uint, uint) {
	return c.OidMap.Min, c.OidMap.Max
}

func (c *ConfigIfMap) GetShift() (uint, uint) {
	return c.Shift.Min, c.Shift.Max
}

//
// ConfigTrap2Map is config(/trap2map/<name>)
//
type ConfigTrap2Map map[string]uint

//
// ConfigTrap2sink is config(/trap2sink/<name>)
//
type ConfigTrap2sink struct {
	Addr *net.UDPAddr
}

func (c *ConfigTrap2sink) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data struct {
		Addr string `yaml:"addr"`
	}
	if err := unmarshal(&data); err != nil {
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", data.Addr)
	if err != nil {
		return err
	}
	c.Addr = addr
	return nil
}

//
// Config is config(/ifindex)
//
type Config struct {
	OidMap    []*ConfigOidMap    `yaml:"oidmap"`
	IfMap     *ConfigIfMap       `yaml:"ifmap"`
	Trap2Map  ConfigTrap2Map     `yaml:"trap2map"`
	Trap2Sink []*ConfigTrap2sink `yaml:"trap2sink"`
}

//
// ReadConfig load config file.
//
func ReadConfig(path string) (map[string]*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var configs struct {
		Configs map[string]*Config `yaml:"snmpproxy"`
	}
	if err = yaml.Unmarshal(buf, &configs); err != nil {
		return nil, err
	}

	return configs.Configs, err
}
