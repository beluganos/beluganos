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
	"sort"

	"github.com/spf13/viper"
)

type OptionConfig struct {
	L2SWBridge           string   `mapstructure:"l2sw-bridge"`
	L2SWAgingSec         uint32   `mapstructure:"l2sw-aging-sec"`
	L2SWSweepSec         uint32   `mapstructure:"l2sw-sweep-sec"`
	L2SWNotifyLimit      uint32   `mapstructure:"l2sw-notify-limit"`
	L3PortStart          uint32   `mapstructure:"l3-port-start"`
	L3PortEnd            uint32   `mapstructure:"l3-port-end"`
	L3VlanBase           uint16   `mapstructure:"l3-vlan-base"`
	NetlinkSocketBufSize uint64   `mapstructure:"netlink-socket-buf-size"`
	NLARecvChannelSize   uint64   `mapstructure:"nla-recv-channel-size"`
<<<<<<< HEAD
=======
	NLABrVlanUpdateSec   uint32   `mapstructure:"nla-brvlan-update-sec"`
	NLABrVlanChanSize    uint32   `mapstructure:"nla-brvlan-channel-size"`
>>>>>>> develop
	NLACorePort          uint16   `mapstructure:"nla-core-port"`
	NLAAPIPort           uint16   `mapstructure:"nla-api-port"`
	FibcAPIAddr          string   `mapstructure:"fibc-api-addr"`
	FibcAPIPort          uint16   `mapstructure:"fibc-api-port"`
	RibsCorePort         uint16   `mapstructure:"ribs-core-port"`
	RibsAPIPort          uint16   `mapstructure:"ribs-api-port"`
	RibpAPIPort          uint16   `mapstructure:"ribp-api-port"`
	RibxLogLevel         uint8    `mapstructure:"ribx-log-level"`
	RibxLogDump          uint8    `mapstructure:"ribx-log-dump"`
	LXDBridge            string   `mapstructure:"lxd-bridge"`
	LXDMtu               uint16   `mapstructure:"lxd-mtu"`
	LXDMngInterface      string   `mapstructure:"lxd-mng-interface"`
	LXDConfigMode        bool     `mapstructure:"lxd-config-mode"`
	SnmpproxydAddr       string   `mapstructure:"snmpproxyd-addr"`
	SnmpproxydTrap2Sink  []string `mapstructure:"snmpproxyd-trap2sink"`
	SnmpproxydIfResend   string   `mapstructure:"snmpproxyd-if-resend"`
	SnmpdLinkmonInterval uint32   `mapstructure:"snmpd-linkmon-interval"`
	SnmpdListenPort      uint16   `mapstructure:"snmpd-listen-port"`
	SysctlMPLSLabelMax   uint32   `mapstructure:"sysctl-mpls-label-max"`
	GoBGPAs              uint32   `mapstructure:"gobgp-as"`
	GoBGPZAPIVersion     uint16   `mapstructure:"gobgp-zapi-version"`
	GoBGPZAPIEnable      bool     `mapstructure:"gobgp-zapi-enable"`
	VPMNexthopNetwork    string   `mapstructure:"vpn-nexthop-network"`
	VPNPseudoBridge      string   `mapstructure:"vpn-pseudo-bridge"`
	GonsldListenAddr     string   `mapstructure:"gonsld-listen-addr"`
	GonsldListenPort     uint16   `mapstructure:"gonsld-listen-port"`
}

var optionConfigDefaults = map[string]interface{}{
	"l2sw-bridge":             "l2swbr0",
	"l2sw-aging-sec":          uint32(3600),
	"l2sw-sweep-sec":          uint32(3),
	"l2sw-notify-limit":       uint32(250),
	"l3-port-start":           uint32(1),
	"l3-port-end":             uint32(190),
	"l3-vlan-base":            uint16(3900),
	"netlink-socket-buf-size": uint64(8388608),
	"nla-recv-channel-size":   uint64(65536),
<<<<<<< HEAD
=======
	"nla-brvlan-update-sec":   uint32(60 * 30),
	"nla-brvlan-channel-size": uint32(4096 * 4),
>>>>>>> develop
	"nla-core-port":           uint16(50061),
	"nla-api-port":            uint16(50062),
	"fibc-api-addr":           "192.169.1.1",
	"fibc-api-port":           uint16(50070),
	"ribs-core-port":          uint16(50071),
	"ribs-api-port":           uint16(50072),
	"ribp-api-port":           uint16(50091),
	"ribx-log-level":          uint8(0),
	"ribx-log-dump":           uint8(0),
	"lxd-bridge":              "lxdbr0",
	"lxd-mtu":                 uint16(9000),
	"lxd-mng-interface":       "eth0",
	"lxd-config-mode":         false,
	"snmpproxyd-addr":         "192.169.1.1",
	"snmpproxyd-trap2sink":    []string{"192.160.122.1:161"},
	"snmpproxyd-if-resend":    "15m",
	"snmpd-linkmon-interval":  10,
	"snmpd-listen-port":       uint16(1161),
	"sysctl-mpls-label-max":   10240,
	"gobgp-as":                uint32(65001),
	"gobgp-zapi-version":      uint16(6),
	"gobgp-zapi-enable":       false,
	"vpn-nexthop-network":     "1.1.0.0/24",
	"vpn-pseudo-bridge":       "ffbr0",
	"gonsld-listen-addr":      "",
	"gonsld-listen-port":      uint16(50061),
}

func setLXDConfigMode(b bool) {
	optionConfigDefaults["lxd-config-mode"] = b
}

func (c *OptionConfig) setDefault(v *viper.Viper) {
	for key, val := range optionConfigDefaults {
		name := fmt.Sprintf("option.%s", key)
		v.SetDefault(name, val)
	}
}

type OptionConfigEntry struct {
	Key string
	Val interface{}
}

func (c *OptionConfig) List() []*OptionConfigEntry {
	keys := []string{}
	for key := range optionConfigDefaults {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	entries := make([]*OptionConfigEntry, len(keys))
	for index, key := range keys {
		entries[index] = &OptionConfigEntry{
			Key: key,
			Val: optionConfigDefaults[key],
		}
	}

	return entries
}
