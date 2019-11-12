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

package main

import "github.com/spf13/viper"

type FibsStatsConfig struct {
	Names []string `mapstructure:"names"`
}

func NewFibsStatsConfig() *FibsStatsConfig {
	return &FibsStatsConfig{
		Names: []string{},
	}
}

type FibsConfig struct {
	FibsStats *FibsStatsConfig `mapstructure:"stats"`
}

func NewFibsConfig() *FibsConfig {
	return &FibsConfig{
		FibsStats: NewFibsStatsConfig(),
	}
}

type Config struct {
	FibsConfig map[string]*FibsConfig `mapstructure:"fibs"`
}

func NewConfig() *Config {
	return &Config{
		FibsConfig: map[string]*FibsConfig{},
	}
}

func (c *Config) GetFibsConfig(name string) *FibsConfig {
	if cfg, ok := c.FibsConfig[name]; ok {
		return cfg
	}

	return nil
}

func (c *Config) ReadFile(configFile, configType string) error {
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType(configType)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.UnmarshalExact(c); err != nil {
		return err
	}

	return nil
}
