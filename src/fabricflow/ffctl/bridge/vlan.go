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

package bridge

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

func addAccessVlanLink(link netlink.Link, vid uint16) error {
	return netlink.BridgeVlanAdd(link, vid, true, true, false, true)
}

func delAccessVlanLink(link netlink.Link, vid uint16) error {
	return netlink.BridgeVlanDel(link, vid, true, true, false, true)
}

func addTrunkVlanLink(link netlink.Link, vid uint16) error {
	return netlink.BridgeVlanAdd(link, vid, false, false, false, true)
}

func delTrunkVlanLink(link netlink.Link, vid uint16) error {
	return netlink.BridgeVlanDel(link, vid, false, false, false, true)
}

func stringBridgeVlanInfoFlags(flags uint16) string {
	ss := []string{}
	if (flags & nl.BRIDGE_VLAN_INFO_MASTER) != 0 {
		ss = append(ss, "master")
	}
	if (flags & nl.BRIDGE_VLAN_INFO_PVID) != 0 {
		ss = append(ss, "pvid")
	}
	if (flags & nl.BRIDGE_VLAN_INFO_UNTAGGED) != 0 {
		ss = append(ss, "untagged")
	}
	if (flags & nl.BRIDGE_VLAN_INFO_RANGE_BEGIN) != 0 {
		ss = append(ss, "range-begin")
	}
	if (flags & nl.BRIDGE_VLAN_INFO_RANGE_END) != 0 {
		ss = append(ss, "range-end")
	}

	return strings.Join(ss, " ")
}

func showBrVlanInfo(vid uint16, brvlanInfoMap map[int32][]*nl.BridgeVlanInfo) {
	for ifindex, brvlanInfos := range brvlanInfoMap {
		for _, brvlanInfo := range brvlanInfos {
			if vid != 0 && vid != brvlanInfo.Vid {
				continue
			}

			link, err := netlink.LinkByIndex(int(ifindex))
			if err != nil {
				fmt.Printf("failed to get link. %s\n", err)
				continue
			}

			flags := stringBridgeVlanInfoFlags(brvlanInfo.Flags)
			fmt.Printf("%s %d %s\n", link.Attrs().Name, brvlanInfo.Vid, flags)
		}
	}
}

func stringBrVlanPortCommand(link netlink.Link, vid uint16, brname string, flags uint16) []string {

	ifname := link.Attrs().Name
	return []string{
		fmt.Sprintf("# create port %s vid %d", ifname, vid),
		fmt.Sprintf("ip link set %s down", ifname),
		fmt.Sprintf("ip link set %s master %s", ifname, brname),
		fmt.Sprintf("ip link set %s mtu %d", ifname, link.Attrs().MTU),
		fmt.Sprintf("ip link set %s up", ifname),
		fmt.Sprintf("bridge vlan add vid %d dev %s %s", vid, ifname, stringBridgeVlanInfoFlags(flags)),
	}
}

func stringBrVlanBridgeCommand(brname string) []string {
	return []string{
		fmt.Sprintf("# create bridge %s", brname),
		fmt.Sprintf("ip link add %s type bridge vlan_filtering 1", brname),
		fmt.Sprintf("ip link set %s multicast off", brname),
		fmt.Sprintf("ip link set %s up", brname),
	}
}

func showBrVlanCommand(vid uint16, brvlanInfoMap map[int32][]*nl.BridgeVlanInfo) {
	brnames := map[string]struct{}{}
	cmdlist := []string{}
	for ifindex, brvlanInfos := range brvlanInfoMap {
		for _, brvlanInfo := range brvlanInfos {
			if vid != 0 && vid != brvlanInfo.Vid {
				continue
			}

			link, err := netlink.LinkByIndex(int(ifindex))
			if err != nil {
				fmt.Printf("failed to get link. %s\n", err)
				continue
			}

			masterIndex := int(link.Attrs().MasterIndex)
			if masterIndex == 0 || masterIndex == int(ifindex) {
				continue
			}

			bridge, err := netlink.LinkByIndex(masterIndex)
			if err != nil {
				fmt.Printf("failed to get master device. %s %s\n", link.Attrs().Name, err)
				continue
			}

			brname := bridge.Attrs().Name
			if _, ok := brnames[brname]; !ok {
				brnames[brname] = struct{}{}
				cmdlist = append(
					cmdlist,
					stringBrVlanBridgeCommand(brname)...,
				)
			}

			cmdlist = append(
				cmdlist,
				stringBrVlanPortCommand(link, brvlanInfo.Vid, brname, brvlanInfo.Flags)...,
			)
		}
	}

	for _, s := range cmdlist {
		fmt.Println(s)
	}
}

func showBrVlanAsYaml(vid uint16, brvlanInfoMap map[int32][]*nl.BridgeVlanInfo) {
	cfg := NewBrVlanConfig()
	for ifindex, brvlanInfos := range brvlanInfoMap {
		for _, brvlanInfo := range brvlanInfos {
			if vid != 0 && vid != brvlanInfo.Vid {
				continue
			}

			link, err := netlink.LinkByIndex(int(ifindex))
			if err != nil {
				fmt.Printf("failed to get link. %s\n", err)
				continue
			}

			ifname := link.Attrs().Name

			masterIndex := int(link.Attrs().MasterIndex)
			if masterIndex == 0 || masterIndex == int(ifindex) {
				continue
			}

			bridge, err := netlink.LinkByIndex(masterIndex)
			if err != nil {
				fmt.Printf("failed to get master device. %s %s\n", ifname, err)
				continue
			}

			brname := bridge.Attrs().Name

			brcfg, ok := cfg.Network.Bridges[brname]
			if !ok {
				brcfg = NewBrVlanBridgeConfig()
				cfg.Network.Bridges[brname] = brcfg
			}

			brcfg.Interfaces = append(brcfg.Interfaces, ifname)

			if brvlanInfo.Vid != 1 {
				vlancfg, ok := cfg.Network.Vlans[ifname]
				if !ok {
					vlancfg = NewBrVlanVlanConfig()
					cfg.Network.Vlans[ifname] = vlancfg
				}

				mask := uint16(nl.BRIDGE_VLAN_INFO_UNTAGGED | nl.BRIDGE_VLAN_INFO_UNTAGGED)
				if (brvlanInfo.Flags & mask) == 0 {
					// truk port
					vlancfg.Ids = append(vlancfg.Ids, brvlanInfo.Vid)
				} else if (brvlanInfo.Flags & mask) == mask {
					// access port
					vlancfg.Id = brvlanInfo.Vid
				}
			}
		}
	}

	if b, err := cfg.Yaml(); err == nil {
		fmt.Printf("%s\n", b)
	}
}

type BridgeVlanCmd struct {
	ConfigFile string
	ConfigType string
}

func (c *BridgeVlanCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *BridgeVlanCmd) setApplyFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.ConfigFile, "config-file", "", "/etc/beluganos/bridge_vlan.yaml", "config file path.")
	cmd.PersistentFlags().StringVarP(&c.ConfigType, "config-type", "", "yaml", "config file type.")
	return c.setFlags(cmd)
}

func (c *BridgeVlanCmd) show(vlans []string) error {
	brblanInfos, err := netlink.BridgeVlanList()
	if err != nil {
		return err
	}

	if len(vlans) == 0 {
		showBrVlanInfo(0, brblanInfos)
		return nil
	}

	for _, vlan := range vlans {
		vid, err := strconv.ParseUint(vlan, 0, 16)
		if err != nil {
			return err
		}

		showBrVlanInfo(uint16(vid), brblanInfos)
	}

	return nil
}

func (c *BridgeVlanCmd) dump(vlans []string) error {
	brblanInfos, err := netlink.BridgeVlanList()
	if err != nil {
		return err
	}

	if len(vlans) == 0 {
		showBrVlanCommand(0, brblanInfos)
		return nil
	}

	for _, vlan := range vlans {
		vid, err := strconv.ParseUint(vlan, 0, 16)
		if err != nil {
			return err
		}

		showBrVlanCommand(uint16(vid), brblanInfos)
	}

	return nil
}

func (c *BridgeVlanCmd) dumpYaml(vlans []string) error {
	brblanInfos, err := netlink.BridgeVlanList()
	if err != nil {
		return err
	}

	if len(vlans) == 0 {
		showBrVlanAsYaml(0, brblanInfos)
		return nil
	}

	for _, vlan := range vlans {
		vid, err := strconv.ParseUint(vlan, 0, 16)
		if err != nil {
			return err
		}

		showBrVlanAsYaml(uint16(vid), brblanInfos)
	}

	return nil
}

func (c *BridgeVlanCmd) addBridge(ifname string) error {
	bridge := newBridge(ifname)
	if err := netlink.LinkAdd(bridge); err != nil {
		return err
	}

	log.Debugf("add bridge %s", ifname)

	if err := netlink.LinkSetMulticastOff(bridge); err != nil {
		return err
	}

	log.Debugf("set bridge %s multicast off", ifname)

	if err := netlink.LinkSetUp(bridge); err != nil {
		return err
	}

	log.Debugf("set bridge %s up", ifname)
	return nil
}

func (c *BridgeVlanCmd) delBridge(ifname string) error {
	bridge := netlink.Bridge{}
	bridge.Attrs().Name = ifname
	if err := netlink.LinkDel(&bridge); err != nil {
		return err
	}

	log.Debugf("del bridge %s", ifname)

	return nil
}

func (c *BridgeVlanCmd) addPortToBridge(ifname, brname string) error {
	bridge := newBridge(brname)
	link := newLink(ifname)

	if err := netlink.LinkSetDown(link); err != nil {
		return err
	}

	log.Debugf("set link %s down", ifname)

	if err := netlink.LinkSetMaster(link, bridge); err != nil {
		return err
	}

	log.Debugf("set link %s master %s", ifname, brname)

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	log.Debugf("set link %s up", ifname)

	return nil
}

func (c *BridgeVlanCmd) delPortFromBridge(ifname string) error {
	link := newLink(ifname)

	if err := netlink.LinkSetDown(link); err != nil {
		return err
	}

	log.Debugf("set link %s down", ifname)

	if err := netlink.LinkSetNoMaster(link); err != nil {
		return err
	}

	log.Debugf("set link %s nomaster", ifname)

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	log.Debugf("set link %s up", ifname)

	return nil
}

func (c *BridgeVlanCmd) addAccessVlan(ifname string, vid uint16) error {
	link := newLink(ifname)

	if err := delAccessVlanLink(link, 1); err != nil {
		return err
	}

	log.Debugf("del access vlan 1 from %s", ifname)

	if err := addAccessVlanLink(link, uint16(vid)); err != nil {
		return err
	}

	log.Debugf("add access vlan %d to %s", vid, ifname)

	return nil
}

func (c *BridgeVlanCmd) delAccessVlan(ifname string, vid uint16) error {
	link := newLink(ifname)

	if err := delAccessVlanLink(link, uint16(vid)); err != nil {
		return err
	}

	log.Debugf("del access vlan %d from %s", vid, ifname)

	if err := addAccessVlanLink(link, 1); err != nil {
		return err
	}

	log.Debugf("add access vlan 1 to %s", ifname)

	return nil
}

func (c *BridgeVlanCmd) addTrunkVlan(ifname string, vid uint16) error {
	link := newLink(ifname)

	if err := delTrunkVlanLink(link, 1); err != nil {
		return err
	}

	log.Debugf("del trunk vlan 1 from %s", ifname)

	if err := addTrunkVlanLink(link, uint16(vid)); err != nil {
		return err
	}

	log.Debugf("add trunk vlan %d to %s", vid, ifname)

	return nil
}

func (c *BridgeVlanCmd) delTrunkVlan(ifname string, vid uint16) error {
	link := newLink(ifname)

	if err := delTrunkVlanLink(link, uint16(vid)); err != nil {
		return err
	}

	log.Debugf("del trunk vlan %d from %s", vid, ifname)

	brblanMap, err := netlink.BridgeVlanList()
	if err != nil {
		return err
	}

	if brvlans, ok := brblanMap[int32(link.Attrs().Index)]; ok && len(brvlans) != 0 {
		// other vlan exists.
		log.Debugf("some trunk vlan exist on %s", ifname)
		return nil
	}

	// default port type is access port.
	if err := addAccessVlanLink(link, 1); err != nil {
		return err
	}

	log.Debugf("add access vlan 1 to %s", ifname)

	return nil
}

func (c *BridgeVlanCmd) apply() error {
	cfg := NewBrVlanConfig()
	if err := cfg.SetConfigFile(c.ConfigFile, c.ConfigType).Load(); err != nil {
		return err
	}

	if cfg.Network.Bridges != nil {
		for brname, bridge := range cfg.Network.Bridges {
			if err := c.addBridge(brname); err != nil {
				log.Warnf("add bridge error. %s %s", brname, err)
			}

			if bridge.Interfaces == nil {
				continue
			}

			for _, ifname := range bridge.Interfaces {
				if err := c.addPortToBridge(ifname, brname); err != nil {
					log.Warnf("add bridge to port error. %s %s %s", ifname, brname, err)
				}
			}
		}
	}

	if cfg.Network.Vlans != nil {
		for ifname, vlan := range cfg.Network.Vlans {
			if vlan.Ids != nil && len(vlan.Ids) > 0 {
				for _, vid := range vlan.Ids {
					if err := c.addTrunkVlan(ifname, vid); err != nil {
						log.Errorf("add trunk vlan error. %s %d %s", ifname, vid, err)
					}
				}

			} else if vlan.Id != 0 {
				if err := c.addAccessVlan(ifname, vlan.Id); err != nil {
					log.Errorf("add access vlan error. %s %d %s", ifname, vlan.Id, err)
				}
			}
		}
	}

	return nil
}

func bridgeVlanCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vlan",
		Short: "bridge vlan command.",
	}

	brvlan := &BridgeVlanCmd{}
	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "show [vid...]",
			Short: "show ports by vlan.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return brvlan.show(args)
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "dump [vid...]",
			Short: "dump ports by vlan and show commands.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return brvlan.dump(args)
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "dump-yaml [vid...]",
			Short: "dump ports by vlan and show yaml.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return brvlan.dumpYaml(args)
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "add-br <ifname>",
			Short: "add bridge device.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return brvlan.addBridge(args[0])
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "del-br <ifname>",
			Short: "delete bridge device.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return brvlan.delBridge(args[0])
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "add-ports <bridge> <ifname...>",
			Short: "add ports to bridge.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				bridge := args[0]
				for _, ifname := range args[1:] {
					if err := brvlan.addPortToBridge(ifname, bridge); err != nil {
						return err
					}
				}
				return nil
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "del-ports <bridge> <ifname...>",
			Short: "delete ports from bridge.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				// bridge := args[0]
				for _, ifname := range args[1:] {
					if err := brvlan.delPortFromBridge(ifname); err != nil {
						return err
					}
				}
				return nil
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "add-access <vid> <ifname...>",
			Short: "add access vlan to ports.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				vid, err := strconv.ParseUint(args[0], 0, 16)
				if err != nil {
					return err
				}

				for _, ifname := range args[1:] {
					if err := brvlan.addAccessVlan(ifname, uint16(vid)); err != nil {
						return err
					}
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "del-access <vid> <ifname...>",
			Short: "delete access vlan from ports.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				vid, err := strconv.ParseUint(args[0], 0, 16)
				if err != nil {
					return err
				}

				for _, ifname := range args[1:] {
					if err := brvlan.delAccessVlan(ifname, uint16(vid)); err != nil {
						return err
					}
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "add-trunk <ifname> <vid...>",
			Short: "add trunk vlans to port.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				ifname := args[0]

				for _, arg := range args[1:] {
					vid, err := strconv.ParseUint(arg, 0, 16)
					if err != nil {
						return err
					}

					if err := brvlan.addTrunkVlan(ifname, uint16(vid)); err != nil {
						return err
					}
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(brvlan.setFlags(
		&cobra.Command{
			Use:   "del-trunk <ifname> <vid...>",
			Short: "delete trunk vlans from port.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				ifname := args[0]

				for _, arg := range args[1:] {
					vid, err := strconv.ParseUint(arg, 0, 16)
					if err != nil {
						return err
					}

					if err := brvlan.delTrunkVlan(ifname, uint16(vid)); err != nil {
						return err
					}
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(brvlan.setApplyFlags(
		&cobra.Command{
			Use:   "net-apply",
			Short: "apply bridge vlan settings",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return brvlan.apply()
			},
		},
	))

	return rootCmd
}
