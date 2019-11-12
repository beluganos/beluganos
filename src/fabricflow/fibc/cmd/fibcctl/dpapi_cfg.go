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

import (
	"fmt"

	"github.com/spf13/viper"
)

//
// DPAPIMultipartPortStatsConfig is DPAPI Multipart Port Status command data.
//
type DPAPIMultipartPortStatsConfig struct {
	Values map[string]uint64
}

//
// NewDPAPIMultipartPortStatsConfig returns nwq DPAPIMultipartPortStatsConfig.
//
func NewDPAPIMultipartPortStatsConfig(v map[string]uint64) *DPAPIMultipartPortStatsConfig {
	return &DPAPIMultipartPortStatsConfig{
		Values: v,
	}
}

//
// PortID returns portID
//
func (c *DPAPIMultipartPortStatsConfig) PortID() uint32 {
	if v, ok := c.Values["port_no"]; ok {
		return uint32(v)
	}
	return 0
}

//
// GetValues returns map of values.
//
func (c *DPAPIMultipartPortStatsConfig) GetValues() map[string]uint64 {
	values := map[string]uint64{}
	for name, value := range c.Values {
		if name != "port_no" {
			values[name] = value
		}
	}

	return values
}

//
// String is stringer.
//
func (c *DPAPIMultipartPortStatsConfig) String() string {
	return fmt.Sprintf("%v", c.Values)
}

//
// DPAPIMultipartConfig is DPAPI Multipart config.
//
type DPAPIMultipartConfig struct {
	PortStats []map[string]uint64 `mapstructure:"portstats"`
}

//
// GetPortStats retuen port status.
//
func (c *DPAPIMultipartConfig) GetPortStats() []*DPAPIMultipartPortStatsConfig {
	cc := []*DPAPIMultipartPortStatsConfig{}

	if c != nil && c.PortStats != nil {
		for _, ps := range c.PortStats {
			cc = append(cc, NewDPAPIMultipartPortStatsConfig(ps))
		}
	}

	return cc
}

//
// DPAPIConfig is DPAPI config
//
type DPAPIConfig struct {
	Multipart *DPAPIMultipartConfig `mapstructure:"multipart"`
}

//
// ReadFile read config file.
//
func (c *DPAPIConfig) ReadFile(path, t string) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(t)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
