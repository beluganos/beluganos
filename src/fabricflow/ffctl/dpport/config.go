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

func ParsePortLine(s string) (uint, []uint, error) {
	ss := strings.Split(s, ".")
	if len(ss) == 0 {
		return 0, nil, fmt.Errorf("bad iface, '%s'", s)
	}

	index, err := strconv.ParseUint(ss[0], 10, 32)
	if err != nil {
		return 0, nil, err
	}

	vids := []uint{}
	for _, vidstr := range ss[1:] {
		vid, err := strconv.ParseUint(vidstr, 10, 32)
		if err != nil {
			return 0, nil, err
		}

		vids = append(vids, uint(vid))
	}

	return uint(index), vids, nil
}

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
	pport, vids, err := ParsePortLine(s)
	if err != nil {
		return nil, err
	}

	switch len(vids) {
	case 0:
		return &VlanConfig{Eth: pport, Vid: 0}, nil

	case 1:
		return &VlanConfig{Eth: pport, Vid: vids[0]}, nil

	default:
		return nil, fmt.Errorf("Bad vlan. %s", s)
	}
}

type L2SWAccessConfig struct {
	Eth uint
	Vid uint
}

func NewL2SWAccessConfig(eth, vid uint) *L2SWAccessConfig {
	return &L2SWAccessConfig{
		Eth: eth,
		Vid: vid,
	}
}

func (c *L2SWAccessConfig) String() string {
	if c.Vid != 0 {
		return fmt.Sprintf("%d.%d", c.Eth, c.Vid)
	}
	return fmt.Sprintf("%d", c.Eth)
}

func ParseL2SWAccessConfig(s string) (*L2SWAccessConfig, error) {
	port, vids, err := ParsePortLine(s)
	if err != nil {
		return nil, err
	}

	if len(vids) != 1 {
		return nil, fmt.Errorf("bad access port. %s", s)
	}

	return NewL2SWAccessConfig(port, vids[0]), nil
}

type L2SWTrunkConfig struct {
	Eth  uint
	Vids []uint
}

func NewL2SWTrunkConfig(eth uint, vids []uint) *L2SWTrunkConfig {
	return &L2SWTrunkConfig{
		Eth:  eth,
		Vids: vids,
	}
}

func (c *L2SWTrunkConfig) String() string {
	vv := []string{fmt.Sprintf("%d", c.Eth)}
	for _, vid := range c.Vids {
		vv = append(vv, fmt.Sprintf("%d", vid))
	}
	return strings.Join(vv, ".")
}

func ParseL2SWTrunkConfig(s string) (*L2SWTrunkConfig, error) {
	port, vids, err := ParsePortLine(s)
	if err != nil {
		return nil, err
	}

	if len(vids) == 0 {
		return nil, fmt.Errorf("bad trunk port. %s", s)
	}

	return NewL2SWTrunkConfig(port, vids), nil
}

type L2SWConfig struct {
	Access []string `yaml:"access"`
	Trunk  []string `yaml:"trunk"`
}

type PortConfig struct {
	Eth  []uint      `yaml:"eth"`
	Vlan []string    `yaml:"vlan"`
	L2sw *L2SWConfig `yaml:"l2sw"`
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

func (c *PortConfig) L2SWAccessPorts() []*L2SWAccessConfig {
	if c.L2sw == nil || c.L2sw.Access == nil {
		return nil
	}

	cfgs := []*L2SWAccessConfig{}
	for _, s := range c.L2sw.Access {
		cfg, err := ParseL2SWAccessConfig(s)
		if err != nil {
			log.Errorf("Invalid L2SW(access). %s", s)
		} else {
			cfgs = append(cfgs, cfg)
		}
	}

	return cfgs
}

func (c *PortConfig) L2SWTrunkPorts() []*L2SWTrunkConfig {
	if c.L2sw == nil || c.L2sw.Trunk == nil {
		return nil
	}

	cfgs := []*L2SWTrunkConfig{}
	for _, s := range c.L2sw.Trunk {
		cfg, err := ParseL2SWTrunkConfig(s)
		if err != nil {
			log.Errorf("Invalid L2SW(trunk). %s", s)
		} else {
			cfgs = append(cfgs, cfg)
		}
	}

	return cfgs
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
