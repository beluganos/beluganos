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
	"fabricflow/util/netplan"
	"fabricflow/util/sysctl"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/vishvananda/netlink"
)

const (
	SYSCTL_CONF  = "/etc/sysctl.d/30-beluganos.conf"
	NETPLAN_CONF = "/etc/netplan/20-beluganos.yaml"
	STDOUT_MODE  = "-"
	TEMP_MODE    = "temp"
)

func parseVid(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}

	if v > 0xffff {
		return 0, fmt.Errorf("Invalid VlanID. '%s'", s)
	}

	return uint16(v), nil
}

func parseVlanProto(s string) (netlink.VlanProtocol, error) {
	if len(s) == 0 {
		return netlink.VLAN_PROTOCOL_UNKNOWN, nil
	}

	proto := netlink.StringToVlanProtocol(s)
	if proto == netlink.VLAN_PROTOCOL_UNKNOWN {
		return proto, fmt.Errorf("Invalid vlan proto. %s", s)
	}

	return proto, nil
}

func vlanLinkList() ([]*netlink.Vlan, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	vlans := []*netlink.Vlan{}
	for _, link := range links {
		vlan, ok := link.(*netlink.Vlan)
		if ok {
			vlans = append(vlans, vlan)
		}
	}

	return vlans, nil
}

func isStdout(p string) bool {
	return (len(p) == 0) || (p == STDOUT_MODE)
}

func openOutputFile(p string) (*os.File, error) {
	if isStdout(p) {
		return os.Stdout, nil
	}

	if p == TEMP_MODE {
		return ioutil.TempFile("", p)
	}

	return os.Create(p)
}

func sysctlMplsInputPath(ifname string, vid uint16) string {
	if vid == 0 {
		return fmt.Sprintf("net.mpls.conf.%s.input", ifname)
	}

	return fmt.Sprintf("net.mpls.conf.%s/%d.input", ifname, vid)
}

func sysctlRpFilterPath(ifname string, vid uint16) string {
	if vid == 0 {
		return fmt.Sprintf("net.ipv4.conf.%s.rp_filter", ifname)
	}

	return fmt.Sprintf("net.ipv4.conf.%s/%d.rp_filter", ifname, vid)
}

func sysctlReadConfig(path string) (*sysctl.SysctlConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := sysctl.ReadConfig(f)
	return cfg, nil
}

func sysctlWriteConfig(output string, cfg *sysctl.SysctlConfig) error {
	temp, err := openOutputFile(output)
	if err != nil {
		return err
	}

	defer func() {
		if !isStdout(output) {
			temp.Close()
		}
	}()

	if _, err := cfg.WriteTo(temp); err != nil {
		if !isStdout(output) {
			os.Remove(temp.Name())
		}
		return err
	}

	return nil
}

func netplanReadConfig(path string) (map[interface{}]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg, err := netplan.ReadConfig(f)
	return cfg, err
}

func netplanWriteConfig(output string, m map[interface{}]interface{}) error {
	temp, err := openOutputFile(output)
	if err != nil {
		return err
	}

	defer func() {
		if !isStdout(output) {
			temp.Close()
		}
	}()

	if err := netplan.WriteConfig(temp, m); err != nil {
		if !isStdout(output) {
			os.Remove(temp.Name())
		}
		return err
	}

	return nil
}
