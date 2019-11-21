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
	"regexp"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Ifnames   []string        `mapstructure:"ifnames"`
	Patterns  []string        `mapstructure:"patterns"`
	BlackList []string        `mapstructure:"blacklist"`
	Features  map[string]bool `mapstructure:"features"`
}

type Configs struct {
	Config map[string]*Config `mapstructure:"ribn"`

	viper *viper.Viper
}

func NewConfig() *Configs {
	return &Configs{
		Config: map[string]*Config{},

		viper: viper.New(),
	}
}

func (c *Configs) SetConfigFile(path, typ string) *Configs {
	c.viper.SetConfigFile(path)
	c.viper.SetConfigType(typ)

	return c
}

func (c *Configs) Load() error {
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return c.viper.UnmarshalExact(c)
}

func (c *Configs) Get(name string) *Config {
	if cfg, ok := c.Config[name]; ok {
		return cfg
	}

	return nil
}

type ConfigDB struct {
	ifnames   map[string]struct{}
	patterns  map[string]*regexp.Regexp
	blacklist map[string]struct{}
	features  map[string]bool

	mutex sync.RWMutex
}

func NewConfigDB() *ConfigDB {
	db := &ConfigDB{}
	db.reset()
	return db
}

func (db *ConfigDB) reset() {
	db.ifnames = map[string]struct{}{}
	db.patterns = map[string]*regexp.Regexp{}
	db.blacklist = map[string]struct{}{}
	db.features = map[string]bool{}
}

func (db *ConfigDB) addPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	db.patterns[pattern] = re

	return nil
}

func (db *ConfigDB) AddPattern(pattern string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.addPattern(pattern)
}

func (db *ConfigDB) addIfname(ifname string) {
	db.ifnames[ifname] = struct{}{}
}

func (db *ConfigDB) AddIfname(ifname string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.addIfname(ifname)
}

func (db *ConfigDB) addBlacklist(ifname string) {
	db.blacklist[ifname] = struct{}{}
}

func (db *ConfigDB) AddBlacklist(ifname string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.addBlacklist(ifname)
}

func (db *ConfigDB) addFeature(name string, b bool) {
	db.features[name] = b
}

func (db *ConfigDB) AddFeature(name string, b bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.addFeature(name, b)
}

func (db *ConfigDB) has(ifname string) bool {
	if _, ok := db.blacklist[ifname]; ok {
		return false
	}

	if _, ok := db.ifnames[ifname]; ok {
		return true
	}

	for _, pattern := range db.patterns {
		if match := pattern.MatchString(ifname); match {
			return true
		}
	}

	return false
}

func (db *ConfigDB) Has(ifname string) bool {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.has(ifname)
}

func (db *ConfigDB) update(cfg *Config) {
	db.reset()

	for _, ifname := range cfg.Ifnames {
		db.addIfname(ifname)
	}

	for _, pattern := range cfg.Patterns {
		db.addPattern(pattern)
	}

	for _, blacklist := range cfg.BlackList {
		db.addBlacklist(blacklist)
	}

	for name, b := range cfg.Features {
		db.addFeature(name, b)
	}
}

func (db *ConfigDB) Update(cfg *Config) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.update(cfg)
}

func (db *ConfigDB) Features() map[string]bool {
	return db.features
}
