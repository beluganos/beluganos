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
	"fabricflow/util/container/interfacemap"
	"fabricflow/util/netplan"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
)

type VlanCmd struct {
	sysctlPath  string
	sysctlOut   string
	netplanPath string
	netplanOut  string
	dryRun      bool

	vlanProto string
	mtu       uint16
	persist   bool
}

func (c *VlanCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.vlanProto, "vlan-proto", "", "802.1q", "vlan protocol. (802.1q | 802.1ad)")
	cmd.Flags().Uint16VarP(&c.mtu, "mtu", "", 0, "MTU")
	cmd.Flags().BoolVarP(&c.persist, "apply-to-config", "", false, "apply to config file.")
	return c.setConfigFlags(cmd)
}

func (c *VlanCmd) setDumpFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *VlanCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.sysctlPath, "sysctl", "s", SYSCTL_CONF, "sysctl.conf path.")
	cmd.Flags().StringVarP(&c.sysctlOut, "sysctl-out", "S", SYSCTL_CONF, "sysctl.conf output path.")
	cmd.Flags().StringVarP(&c.netplanPath, "netplan", "n", NETPLAN_CONF, "netpkan.yaml path.")
	cmd.Flags().StringVarP(&c.netplanOut, "netplan-out", "N", NETPLAN_CONF, "netpkan.yaml output path.")
	cmd.Flags().BoolVarP(&c.dryRun, "dry-run", "", false, "dry run mode.")

	return cmd
}

func (c *VlanCmd) SysctlOut() string {
	if c.dryRun {
		return STDOUT_MODE
	}
	return c.sysctlOut
}

func (c *VlanCmd) NetplanOut() string {
	if c.dryRun {
		return STDOUT_MODE
	}
	return c.netplanOut
}

func (c *VlanCmd) addToConfig(ifname string, vlanID string) error {
	vid, err := parseVid(vlanID)
	if err != nil {
		return err
	}

	sysctlCfg, err := sysctlReadConfig(c.sysctlPath)
	if err != nil {
		return err
	}

	log.Debugf("read from %s success.", c.sysctlPath)

	netplanCfg, err := netplanReadConfig(c.netplanPath)
	if err != nil {
		return err
	}

	log.Debugf("read from %s success.", c.netplanPath)

	sysctlCfg.Set(sysctlMplsInputPath(ifname, vid), "1")
	sysctlCfg.Set(sysctlRpFilterPath(ifname, vid), "0")

	m, ok := interfacemap.SelectOrInsert(netplanCfg, netplan.NewVlanPath(ifname, vid)...)
	if !ok {
		return fmt.Errorf("Netplan config already exist and not map. %s %d", ifname, vid)
	}
	m["link"] = ifname
	m["id"] = vid

	if err := sysctlWriteConfig(c.SysctlOut(), sysctlCfg); err != nil {
		return err
	}

	log.Debugf("write to %s success.", c.SysctlOut())

	if err := netplanWriteConfig(c.NetplanOut(), netplanCfg); err != nil {
		return err
	}

	log.Debugf("write to %s success.", c.NetplanOut())

	return nil
}

func (c *VlanCmd) delFromConfig(ifname string, vlanID string) error {
	vid, err := parseVid(vlanID)
	if err != nil {
		return err
	}

	sysctlCfg, err := sysctlReadConfig(c.sysctlPath)
	if err != nil {
		return err
	}

	log.Debugf("read from %s success.", c.sysctlPath)

	netplanCfg, err := netplanReadConfig(c.netplanPath)
	if err != nil {
		return err
	}

	log.Debugf("read from %s success.", c.netplanPath)

	sysctlCfg.Del(sysctlMplsInputPath(ifname, vid))
	sysctlCfg.Del(sysctlRpFilterPath(ifname, vid))

	if ok := interfacemap.Remove(netplanCfg, netplan.NewVlanPath(ifname, vid)...); !ok {
		return fmt.Errorf("Invalid ifname or vlanID. %s %d", ifname, vid)
	}

	if err := sysctlWriteConfig(c.SysctlOut(), sysctlCfg); err != nil {
		return err
	}

	log.Debugf("write to %s success.", c.SysctlOut())

	if err := netplanWriteConfig(c.NetplanOut(), netplanCfg); err != nil {
		return err
	}

	log.Debugf("write to %s success.", c.NetplanOut())

	return nil
}

func (c *VlanCmd) addVlan(ifname string, vlanID string) error {
	vid, err := parseVid(vlanID)
	if err != nil {
		return err
	}

	proto, err := parseVlanProto(c.vlanProto)
	if err != nil {
		return err
	}

	log.Debugf("vlan-protocol %s", proto)

	parentLink, err := netlink.LinkByName(ifname)
	if err != nil {
		return err
	}

	link := netlink.Vlan{
		VlanId:       int(vid),
		VlanProtocol: proto,
	}
	link.Attrs().Name = fmt.Sprintf("%s.%d", ifname, vid)
	link.Attrs().ParentIndex = parentLink.Attrs().Index
	link.Attrs().MTU = int(c.mtu)

	log.Debugf("vlan-link %s", link.Attrs().Name)

	return netlink.LinkAdd(&link)
}

func (c *VlanCmd) delVlan(ifname string, vlanID string) error {
	vid, err := parseVid(vlanID)
	if err != nil {
		return err
	}

	proto, err := parseVlanProto(c.vlanProto)
	if err != nil {
		return err
	}

	log.Debugf("vlan-protocol %s", proto)

	link := netlink.Vlan{
		VlanId:       int(vid),
		VlanProtocol: proto,
	}
	link.Attrs().Name = fmt.Sprintf("%s.%d", ifname, vid)

	log.Debugf("vlan-link %s", link.Attrs().Name)

	return netlink.LinkDel(&link)
}

func (c *VlanCmd) dumpVlan() error {
	vlans, err := vlanLinkList()
	if err != nil {
		return err
	}

	for _, vlan := range vlans {
		parentLink, err := netlink.LinkByIndex(vlan.Attrs().ParentIndex)
		if err != nil {
			fmt.Printf("%s: vid %d %s link unknown\n",
				vlan.Attrs().Name, vlan.VlanId, vlan.VlanProtocol)
		} else {
			fmt.Printf("%s: vid %d %s link %s\n",
				vlan.Attrs().Name, vlan.VlanId, vlan.VlanProtocol, parentLink.Attrs().Name)
		}
	}

	return nil
}

func (c *VlanCmd) dumpCommand() error {
	vlans, err := vlanLinkList()
	if err != nil {
		return err
	}

	for _, vlan := range vlans {
		parentLink, err := netlink.LinkByIndex(vlan.Attrs().ParentIndex)
		if err != nil {
			return err
		}

		fmt.Printf("ip link add link %s name %s type vlan id %d protocol %s\n",
			parentLink.Attrs().Name, vlan.Attrs().Name, vlan.VlanId, vlan.VlanProtocol)
	}

	return nil
}

func vlanCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vlan",
		Short: "VLAN command.",
	}

	vlan := &VlanCmd{}

	rootCmd.AddCommand(vlan.setFlags(
		&cobra.Command{
			Use:   "add <ifname> <vid>",
			Short: "create vlan device as <ifname>.<vid>.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := vlan.addVlan(args[0], args[1]); err != nil {
					return err
				}
				if vlan.persist {
					return vlan.addToConfig(args[0], args[1])
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(vlan.setFlags(
		&cobra.Command{
			Use:   "del <ifname> <vid>",
			Short: "delete vlan device.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := vlan.delVlan(args[0], args[1]); err != nil {
					return err
				}
				if vlan.persist {
					return vlan.delFromConfig(args[0], args[1])
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(vlan.setDumpFlags(
		&cobra.Command{
			Use:   "show",
			Short: "show vlan devices",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vlan.dumpVlan()
			},
		},
	))

	rootCmd.AddCommand(vlan.setDumpFlags(
		&cobra.Command{
			Use:   "dump",
			Short: "show ip commands",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vlan.dumpCommand()
			},
		},
	))

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "vlan config command.",
	}

	configCmd.AddCommand(vlan.setConfigFlags(
		&cobra.Command{
			Use:   "add <ifname> <vid>",
			Short: "Add Vlan device.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vlan.addToConfig(args[0], args[1])
			},
		},
	))

	configCmd.AddCommand(vlan.setConfigFlags(
		&cobra.Command{
			Use:   "del <ifname> <vid>",
			Short: "Delete Vlan device.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vlan.delFromConfig(args[0], args[1])
			},
		},
	))

	rootCmd.AddCommand(configCmd)

	return rootCmd
}
