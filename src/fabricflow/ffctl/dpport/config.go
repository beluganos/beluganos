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

package dpport

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

type VlanConfig struct {
	Eth uint
	Vid uint
}

func NewVlanConfig(eth, vid uint) *VlanConfig {
	return &VlanConfig{
		Eth: eth,
		Vid: vid,
	}
}

func (c *VlanConfig) String() string {
	if c.Vid != 0 {
		return fmt.Sprintf("%d.%d", c.Eth, c.Vid)
	}
	return fmt.Sprintf("%d", c.Eth)
}

func ParseVlanConfig(s string) (*VlanConfig, error) {
	ss := strings.Split(s, ".")
	switch len(ss) {
	case 0:
		return nil, fmt.Errorf("Bad vlan. %s", s)

	case 1:
		pport, err := strconv.ParseUint(ss[0], 10, 32)
		if err != nil {
			return nil, err
		}

		return &VlanConfig{Eth: uint(pport), Vid: 0}, nil

	case 2:
		pport, err := strconv.ParseUint(ss[0], 10, 32)
		if err != nil {
			return nil, err
		}

		vlan, err := strconv.ParseUint(ss[1], 10, 32)
		if err != nil {
			return nil, err
		}

		return &VlanConfig{Eth: uint(pport), Vid: uint(vlan)}, nil

	default:
		return nil, fmt.Errorf("Bad vlan. %s", s)
	}
}

type PortConfig struct {
	Eth  []uint   `yaml:"eth"`
	Vlan []string `yaml:"vlan"`
}

func NewPortConfig() *PortConfig {
	return &PortConfig{
		Eth:  []uint{},
		Vlan: []string{},
	}
}

func (c *PortConfig) Filter(src map[uint]uint) map[uint]uint {
	if len(c.Eth) == 0 {
		log.Debugf("DoPorts is empty. Not filtered.")
		return src
	}

	m := map[uint]uint{}
	for _, pport := range c.Eth {
		if lport, ok := src[pport]; ok {
			m[pport] = lport
		}
	}

	return m
}

func (c *PortConfig) Devices() []*VlanConfig {
	ifaces := c.Vlans()
	for _, eth := range c.Eth {
		ifaces = append(ifaces, NewVlanConfig(eth, 0))
	}

	sort.Slice(ifaces, func(i, j int) bool {
		return strings.Compare(ifaces[i].String(), ifaces[j].String()) < 0
	})
	return ifaces
}

func (c *PortConfig) Vlans() []*VlanConfig {
	vlans := []*VlanConfig{}
	for _, vlan := range c.Vlan {
		vc, err := ParseVlanConfig(vlan)
		if err != nil {
			log.Warnf("Invalid VLAN. %s", vlan)
		} else {
			vlans = append(vlans, vc)
		}
	}

	return vlans
}

type PortsConfig struct {
	Ports map[string]*PortConfig `yaml:"ports"`
}

func NewPortsConfig() *PortsConfig {
	return &PortsConfig{
		Ports: map[string]*PortConfig{},
	}
}

func (c *PortsConfig) PortConfig(name string) (*PortConfig, bool) {
	if c == nil {
		return NewPortConfig(), false
	}

	pc, ok := c.Ports[name]
	if !ok {
		return NewPortConfig(), false
	}

	return pc, ok
}

func ReadConfig(path string) *PortsConfig {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	c := NewPortsConfig()
	if err := yaml.Unmarshal(buf, c); err != nil {
		return nil
	}

	return c
}
