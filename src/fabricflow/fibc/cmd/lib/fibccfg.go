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

package fibccmd

import (
	"fmt"
)

const (
	CFGCMD_ADD      = "add"
	CFGCMD_DEL      = "delete"
	CFGTBL_PORT     = "port"
	CFGTBL_DPID     = "dp"
	CFGTBL_REID     = "re"
	CFGMODE_GENERIC = "generic"
	CFGMODE_OFDPA2  = "ofdpa2"
	CFGMODE_OVS     = "ovs"
)

type Port struct {
	Name   string   `yaml:"name" json:"name"`
	Port   uint32   `yaml:"port" json:"port"`
	Link   string   `yaml:"link" json:"link"`
	Slaves []string `yaml:"slaves" json:"slaves"`
}

func NewPort(port uint32, name string) *Port {
	return &Port{
		Name:   name,
		Port:   port,
		Link:   "",
		Slaves: []string{},
	}
}

func NewVLANPort(port uint32, name string, vid uint16) *Port {
	return &Port{
		Name:   fmt.Sprintf("%s.%d", name, vid),
		Port:   port,
		Link:   name,
		Slaves: []string{},
	}
}

func (p *Port) String() string {
	return fmt.Sprintf("name:'%s', port:%d, link='%s', slaves=%v", p.Name, p.Port, p.Link, p.Slaves)
}

type Router struct {
	Desc     string  `yaml:"desc" json:"desc"`
	ReId     string  `yaml:"re_id" json:"re_id"`
	Datapath string  `yaml:"datapath" json:"datapath"`
	Ports    []*Port `yaml:"ports" json:"ports"`
}

func NewRouter(reid string, dpname string) *Router {
	return &Router{
		ReId:     reid,
		Datapath: dpname,
		Ports:    []*Port{},
	}
}

func (r *Router) AddPort(ports ...*Port) {
	r.Ports = append(r.Ports, ports...)
}

func (r *Router) String() string {
	return fmt.Sprintf("ReId:'%s', Dpath:'%s'", r.ReId, r.Datapath)
}

type Datapath struct {
	Name string `yaml:"name" json:"name"`
	Dpid uint64 `yaml:"dp_id" json:"dp_id"`
	Mode string `yaml:"mode" json:"mode"`
}

func NewDatapath(name string, dpid uint64, mode string) *Datapath {
	return &Datapath{
		Name: name,
		Dpid: dpid,
		Mode: mode,
	}
}

func (d *Datapath) String() string {
	return fmt.Sprintf("name:'%s', dpid:%d, mode:'%s'", d.Name, d.Dpid, d.Mode)
}

type Config struct {
	Routers   []*Router   `yaml:"routers" json:"routers"`
	Datapaths []*Datapath `yaml:"datapaths" json:"datapaths"`
}

type ConfigClient struct {
	url string
}

func NewConfigClient(addr string) *ConfigClient {
	return &ConfigClient{
		url: fmt.Sprintf("http://%s/fib/portmap", addr),
	}
}

func (m *ConfigClient) Add(cfg *Config) error {
	if err := m.ModDpaths(CFGCMD_ADD, cfg.Datapaths); err != nil {
		return err
	}

	if err := m.ModRouters(CFGCMD_ADD, cfg.Routers); err != nil {
		return err
	}

	if err := m.ModPorts(CFGCMD_ADD, cfg.Routers); err != nil {
		return err
	}

	return nil
}

func (m *ConfigClient) Del(cfg *Config) error {
	if err := m.ModPorts(CFGCMD_DEL, cfg.Routers); err != nil {
		return err
	}

	if err := m.ModRouters(CFGCMD_DEL, cfg.Routers); err != nil {
		return err
	}

	if err := m.ModDpaths(CFGCMD_DEL, cfg.Datapaths); err != nil {
		return err
	}

	return nil
}

func (m *ConfigClient) ModDpaths(cmd string, dpaths []*Datapath) error {
	url := fmt.Sprintf("%s/%s/%s", m.url, CFGTBL_DPID, cmd)
	for _, dpath := range dpaths {
		if err := httpJson(url, dpath); err != nil {
			return err
		}
	}
	return nil
}

func (m *ConfigClient) ModRouters(cmd string, routers []*Router) error {
	url := fmt.Sprintf("%s/%s/%s", m.url, CFGTBL_REID, cmd)
	for _, router := range routers {
		if err := httpJson(url, router); err != nil {
			return err
		}
	}
	return nil
}

func (m *ConfigClient) ModPorts(cmd string, routers []*Router) error {
	url := fmt.Sprintf("%s/%s/%s", m.url, CFGTBL_PORT, cmd)
	for _, router := range routers {
		if err := httpJson(url, router); err != nil {
			return err
		}
	}
	return nil
}
