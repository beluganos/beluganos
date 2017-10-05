// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Port struct {
	Name   string   `yaml:"name" json:"name"`
	Port   uint32   `yaml:"port" json:"port"`
	Link   string   `yaml:"link" json:"link"`
	Slaves []string `yaml:"slaves" json:"slaves"`
}

func (p *Port) String() string {
	return fmt.Sprintf("Port{Name: '%s' Port: %d Link: '%s' Slaves: %s}", p.Name, p.Port, p.Link, p.Slaves)
}

type Router struct {
	Desc     string `yaml:"desc" json:"desc"`
	ReId     string `yaml:"re_id" json:"re_id"`
	Datapath string `yaml:"datapath" json:"datapath"`
	Ports    []Port `yaml:"ports" json:"ports"`
}

func (r *Router) String() string {
	return fmt.Sprintf("Router{Desc: '%s' ReId: '%s' Dp: '%s'}", r.Desc, r.ReId, r.Datapath)
}

func (r *Router) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type Datapath struct {
	Name string `yaml:"name" json:"name"`
	Dpid uint64 `yaml:"dp_id" json:"dp_id"`
	Mode string `yaml:"mode" json:"mode"`
}

func (d *Datapath) String() string {
	return fmt.Sprintf("Datapath{Name: '%s' Dpid: %d, Mode: '%s'}", d.Name, d.Dpid, d.Mode)
}

func (d *Datapath) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

type Config struct {
	Routers   []Router   `yaml:"routers"`
	Datapaths []Datapath `yaml:"datapaths"`
}

func ReadConfig(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}
