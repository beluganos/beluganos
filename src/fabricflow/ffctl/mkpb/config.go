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

package mkpb

import (
	"fmt"

	"github.com/spf13/viper"
)

type GlobalConfig struct {
	ReID   string `mapstructure:"re-id"`
	DpID   uint64 `mapstructure:"dp-id"`
	DpType string `mapstructure:"dp-type"`
	DpMode string `mapstructure:"dp-mode"`
	DpAddr string `mapstructure:"dp-addr"`
	Vpn    bool   `mapstructure:"vpn"`
}

func NewGlobalConfig() *GlobalConfig {
	return &GlobalConfig{}
}

func (c *GlobalConfig) String() string {
	return fmt.Sprintf("reid:%s dpid:%d type:%s vpn:%t dp:%s", c.ReID, c.DpID, c.DpType, c.Vpn, c.DpAddr)
}

type L2SWConfig struct {
	Access map[uint32]uint16   `mapstructure:"access"`
	Trunk  map[uint32][]uint16 `mapstructure:"trunk"`
}

func NewL2SWConfig() *L2SWConfig {
	return &L2SWConfig{
		Access: map[uint32]uint16{},
		Trunk:  map[uint32][]uint16{},
	}
}

func (c *L2SWConfig) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("access:%v, trunk:%v", c.Access, c.Trunk)
}

func (c *L2SWConfig) RangeAccess(f func(uint32, uint16)) {
	if c == nil {
		return
	}

	for k, v := range c.Access {
		f(k, v)
	}
}

func (c *L2SWConfig) RangeTrunk(f func(uint32, []uint16)) {
	if c == nil {
		return
	}

	for k, v := range c.Trunk {
		f(k, v)
	}
}

type RouterConfig struct {
	Name    string              `mapstructure:"name"`
	NodeID  uint8               `mapstructure:"nid"`
	Eth     []uint32            `mapstructure:"eth"`
	Vlan    map[uint32][]uint16 `mapstructure:"vlan"`
	L2SW    *L2SWConfig         `mapstructute:"l2sw"`
	Daemons []string            `mapstructure:"daemons"`
	RtRd    []string            `mapstructure:"rt-rd"`
}

func NewRouterConfig() *RouterConfig {
	return &RouterConfig{
		Eth:     []uint32{},
		Vlan:    map[uint32][]uint16{},
		Daemons: []string{},
		RtRd:    []string{},
	}
}

func (c *RouterConfig) String() string {
	return fmt.Sprintf("%s: {nid:%d, eth:%v vlan:%v, l2sw:{%s}, rt/rd:%v}", c.Name, c.NodeID, c.Eth, c.Vlan, c.L2SW, c.RtRd)
}

func (c *RouterConfig) RT() string {
	if c.RtRd != nil && len(c.RtRd) > 0 {
		return c.RtRd[0]
	}
	return ""
}

func (c *RouterConfig) RD() string {
	if c.RtRd != nil && len(c.RtRd) > 1 {
		return c.RtRd[1]
	}
	return ""
}

type Config struct {
	Global *GlobalConfig   `mapstructure:"global"`
	Router []*RouterConfig `mapstructure:"router"`
	Option OptionConfig    `mapstructure:"option"`

	v *viper.Viper
}

func NewConfig() *Config {
	return &Config{
		Global: NewGlobalConfig(),
		Router: []*RouterConfig{},
		v:      viper.New(),
	}
}

func (c *Config) setDefault() {
	c.Option.setDefault(c.v)
}

func (c *Config) GetRouter(name string) *RouterConfig {
	for _, r := range c.Router {
		if r.Name == name {
			return r
		}
	}
	return nil
}

func (c *Config) GetMICRouter() *RouterConfig {
	if len(c.Router) == 0 {
		return nil
	}
	return c.Router[0]
}

func (c *Config) GetRicRouters() []*RouterConfig {
	if len(c.Router) == 0 {
		return nil
	}
	return c.Router[1:]
}

func (c *Config) SetConfig(configFile, configType string) {
	c.v.SetConfigFile(configFile)
	c.v.SetConfigType(configType)
}

func (c *Config) Load() error {
	c.setDefault()
	if err := c.v.ReadInConfig(); err != nil {
		return err
	}

	return c.v.Unmarshal(c)
}
