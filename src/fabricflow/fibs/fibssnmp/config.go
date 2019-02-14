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

	"gopkg.in/yaml.v2"
)

//
// HandlerConfig is handler config
//
type HandlerConfig struct {
	Oid  string `yaml:"oid"`
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (c *HandlerConfig) String() string {
	return fmt.Sprintf("'%s', '%s', '%s'", c.Oid, c.Name, c.Type)
}

//
// Config is gonsld config.
//
type Config struct {
	Handlers []*HandlerConfig `yaml:"handlers"`
}

//
// ReadConfigFile reads config from file.
//
func ReadConfigFile(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ReadConfig(buf)
}

//
// ReadConfig reads config from binary.
//
func ReadConfig(buf []byte) (*Config, error) {
	cfg := &Config{}
	if err := yaml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
