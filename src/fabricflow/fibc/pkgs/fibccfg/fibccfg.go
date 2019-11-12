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

package fibccfg

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/spf13/viper"
)

//
// PortConfig is port configuration.
//
type PortConfig struct {
	Name   string `mapstructure:"name"`
	PortID uint16 `mapstructure:"port"`
}

//
// String is stringer.
//
func (c *PortConfig) String() string {
	return fmt.Sprintf("{name:'%s', port:%d}", c.Name, c.PortID)
}

//
// RouterConfig is router configuration.
//
type RouterConfig struct {
	Desc  string        `mapstructure:"desc"`
	ReID  string        `mapstructure:"re_id"`
	DPath string        `mapstructure:"datapath"`
	Ports []*PortConfig `mapstructure:"ports"`
}

//
// String is stringer.
//
func (c *RouterConfig) String() string {
	return fmt.Sprintf("{desc:'%s', re_id:'%s', dpath:'%s'}", c.Desc, c.ReID, c.DPath)
}

//
// DPathConfig is datapath configuration.
//
type DPathConfig struct {
	Name string `mapstructure:"name"`
	DpID uint64 `mapstructure:"dp_id"`
	Mode string `mapstructure:"mode"`
}

//
// String is stringer.
//
func (c *DPathConfig) String() string {
	return fmt.Sprintf("{name:'%s', dp_id:%d, mode:'%s'}", c.Name, c.DpID, c.Mode)
}

//
// Config is config values.
//
type Config struct {
	Routers []*RouterConfig `mapstructure:"routers"`
	DPaths  []*DPathConfig  `mapstructure:"datapaths"`
}

//
// NewConfig returns new Config
//
func NewConfig() *Config {
	return &Config{
		Routers: []*RouterConfig{},
		DPaths:  []*DPathConfig{},
	}
}

//
// DPathConfig returns DPathConfig value.
//
func (c *Config) DPathConfig(name string) (*DPathConfig, bool) {
	for _, cfg := range c.DPaths {
		if name == cfg.Name {
			return cfg, true
		}
	}

	return nil, false
}

//
// Merge mreges src config.
//
func (c *Config) Merge(src *Config) {
	c.Routers = append(c.Routers, src.Routers...)
	c.DPaths = append(c.DPaths, src.DPaths...)
}

var suffixBlakckList = []string{
	"~",
	".bak",
	".orig",
	".org",
	".old",
}

func hasBlackListSuffix(s string) bool {
	for _, suffix := range suffixBlakckList {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}

//
// ListFiles returns filenames.
//
func ListFiles(dirname string, prefix string) ([]string, error) {
	finfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, finfo := range finfos {
		if finfo.IsDir() {
			continue
		}

		if !finfo.Mode().IsRegular() {
			continue
		}

		if strings.HasPrefix(finfo.Name(), ".") {
			continue
		}

		if hasBlackListSuffix(finfo.Name()) {
			continue
		}

		if len(prefix) > 0 {
			if ok := strings.HasPrefix(finfo.Name(), prefix); !ok {
				continue
			}
		}

		paths = append(paths, path.Join(dirname, finfo.Name()))
	}

	return paths, nil
}

//
// ReadConfig reads config file.
//
func ReadConfig(configPath, configType string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType(configType)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	c := NewConfig()
	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}

//
// ReadConfigs reads config files.
//
func ReadConfigs(dirname, prefix, configType string) (*Config, error) {
	paths, err := ListFiles(dirname, prefix)
	if err != nil {
		return nil, err
	}

	cfg := NewConfig()
	for _, path := range paths {
		c, err := ReadConfig(path, configType)
		if err != nil {
			continue
		}

		cfg.Merge(c)
	}

	return cfg, nil
}
