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

package govsw

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type IfaceConfig struct {
	Names     []string `mapstructure:"names"`
	Patterns  []string `mapstructure:"patterns"`
	BlackList []string `mapstructure:"blacklist"`
}

type DpConfig struct {
	DpId   uint64      `mapstructure:"dpid"`
	Ifaces IfaceConfig `mapstructure:"ifaces"`
}

type Config struct {
	DpConfigs map[string]*DpConfig `mapstructure:"dpaths"`
}

func NewConfig() *Config {
	return &Config{
		DpConfigs: map[string]*DpConfig{},
	}
}

func (c *Config) DpConfig(dpname string) (*DpConfig, bool) {
	if cfg, ok := c.DpConfigs[dpname]; ok {
		return cfg, true
	}
	return nil, false
}

type ConfigServerListener interface {
	ConfigChanged(*DpConfig)
}

type ConfigServer struct {
	DpName  string
	CfgPath string
	CfgType string
	viper   *viper.Viper
}

func (s *ConfigServer) Reload() error {
	return s.viper.ReadInConfig()
}

func (s *ConfigServer) Write() error {
	return s.viper.WriteConfig()
}

func (s *ConfigServer) SetIfNames(ifnames []string) {
	path := fmt.Sprintf("dpaths.%s.ifaces.names", s.DpName)
	s.viper.Set(path, ifnames)
}

func (s *ConfigServer) SetIfPatterns(patterns []string) {
	path := fmt.Sprintf("dpaths.%s.ifaces.patterns", s.DpName)
	s.viper.Set(path, patterns)
}

func (s *ConfigServer) Get() (*DpConfig, error) {
	c := NewConfig()
	if err := s.viper.UnmarshalExact(c); err != nil {
		return nil, err
	}

	if cfg, ok := c.DpConfig(s.DpName); ok {
		return cfg, nil
	}

	return nil, fmt.Errorf("dp config not found. %s", s.DpName)
}

func (s *ConfigServer) Start(listener ConfigServerListener) error {
	s.viper = viper.New()
	s.viper.SetConfigFile(s.CfgPath)
	s.viper.SetConfigType(s.CfgType)

	if err := s.Reload(); err != nil {
		return err
	}

	cfg, err := s.Get()
	if err != nil {
		return err
	}

	listener.ConfigChanged(cfg)

	s.viper.WatchConfig()
	s.viper.OnConfigChange(func(e fsnotify.Event) {
		cfg, err := s.Get()
		if err == nil {
			listener.ConfigChanged(cfg)
		}
	})

	return nil
}
