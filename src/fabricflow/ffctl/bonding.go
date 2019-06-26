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
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
)

//
// BondMode Value
//
type BondModeValue struct {
	netlink.BondMode
}

func (v *BondModeValue) Set(s string) error {
	mode := netlink.StringToBondMode(s)
	if mode == netlink.BOND_MODE_UNKNOWN {
		return fmt.Errorf("Invalid BondMode '%s'", s)
	}

	v.BondMode = mode
	return nil
}

func (v *BondModeValue) Type() string {
	return "BondModeValue"
}

//
// LacpRate Value
//
type BondLacpRateValue struct {
	netlink.BondLacpRate
}

func (v *BondLacpRateValue) Set(s string) error {
	rate := netlink.StringToBondLacpRate(s)
	if rate == netlink.BOND_LACP_RATE_UNKNOWN {
		return fmt.Errorf("Invalid LacpRate. '%s'", s)
	}
	v.BondLacpRate = rate
	return nil
}

func (v *BondLacpRateValue) Type() string {
	return "BondLacpRateValue"
}

//
// BondAdSelect Value
//
type BondAdSelectValue struct {
	netlink.BondAdSelect
}

var bondAdSelect_names = map[netlink.BondAdSelect]string{
	netlink.BOND_AD_SELECT_STABLE:    "stable",
	netlink.BOND_AD_SELECT_BANDWIDTH: "bandwith",
	netlink.BOND_AD_SELECT_COUNT:     "count",
}

var bondAdSelect_values = map[string]netlink.BondAdSelect{
	"stable":   netlink.BOND_AD_SELECT_STABLE,
	"bandwith": netlink.BOND_AD_SELECT_BANDWIDTH,
	"count":    netlink.BOND_AD_SELECT_COUNT,
}

func ParseBondAdSelect(s string) (netlink.BondAdSelect, error) {
	if v, ok := bondAdSelect_values[s]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("Invalid BondAdSelect. '%s'", s)
}

func (v *BondAdSelectValue) String() string {
	if s, ok := bondAdSelect_names[v.BondAdSelect]; ok {
		return s
	}
	return fmt.Sprintf("BondAdSelect(%d)", v)
}

func (v *BondAdSelectValue) Set(s string) error {
	adsel, err := ParseBondAdSelect(s)
	if err != nil {
		return err
	}

	v.BondAdSelect = adsel
	return nil
}

func (v *BondAdSelectValue) Type() string {
	return "BondAdSelectValue"
}

type BondCmd struct {
	Mode     BondModeValue
	MiiMon   int
	MinLinks int
	LacpRate BondLacpRateValue
	AdSelect BondAdSelectValue

	ShowDetail bool
}

func (c *BondCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *BondCmd) setShowFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolVarP(&c.ShowDetail, "detail", "", false, "show details.")
	return cmd
}

func (c *BondCmd) setBondFlags(cmd *cobra.Command) *cobra.Command {
	modes := func() string {
		names := []string{}
		for name, _ := range netlink.StringToBondModeMap {
			names = append(names, name)
		}
		return strings.Join(names, " | ")
	}()

	lacpRates := func() string {
		names := []string{}
		for name, _ := range netlink.StringToBondLacpRateMap {
			names = append(names, name)
		}
		return strings.Join(names, " | ")
	}()

	adSelects := func() string {
		names := []string{}
		for name, _ := range bondAdSelect_values {
			names = append(names, name)
		}
		return strings.Join(names, " | ")
	}()

	cmd.Flags().IntVarP(&c.MiiMon, "miimon", "", -1, "miimon.")
	cmd.Flags().IntVarP(&c.MinLinks, "min-links", "", -1, "minimum links.")
	cmd.Flags().VarP(&c.Mode, "mode", "", fmt.Sprintf("bonde mode. (%s)", modes))
	cmd.Flags().VarP(&c.LacpRate, "lacp-rate", "", fmt.Sprintf("lacp rate. 802.3ad mode only. (%s)", lacpRates))
	cmd.Flags().VarP(&c.AdSelect, "ad-select", "", fmt.Sprintf("ad select. (%s)", adSelects))

	return cmd
}

func (c *BondCmd) addBond(ifname string) error {
	link := netlink.NewLinkBond(
		netlink.LinkAttrs{
			Name: ifname,
		},
	)
	link.Mode = c.Mode.BondMode
	link.Miimon = c.MiiMon
	link.MinLinks = c.MinLinks
	if c.Mode.BondMode == netlink.BOND_MODE_802_3AD {
		link.LacpRate = c.LacpRate.BondLacpRate
	}
	link.AdSelect = c.AdSelect.BondAdSelect

	if err := netlink.LinkAdd(link); err != nil {
		return err
	}

	log.Debugf("add bond %s", ifname)

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	log.Debugf("set bond %s up", ifname)
	return nil
}

func (c *BondCmd) delBond(ifname string) error {
	link := netlink.NewLinkBond(
		netlink.LinkAttrs{
			Name: ifname,
		},
	)
	if err := netlink.LinkDel(link); err != nil {
		return err
	}

	log.Debugf("del bond %s", ifname)
	return nil
}

func (c *BondCmd) addSlave(bondname, ifname string) error {
	bond := netlink.NewLinkBond(
		netlink.LinkAttrs{
			Name: bondname,
		},
	)

	link := newLink(ifname)

	if err := netlink.LinkSetDown(link); err != nil {
		return err
	}

	log.Debugf("set link %s down", ifname)

	if err := netlink.LinkSetBondSlave(link, bond); err != nil {
		return err
	}

	log.Debugf("set link %s master %s", ifname, bondname)

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	log.Debugf("set link %s up", ifname)

	return nil
}

func (c *BondCmd) delSlave(bondname, ifname string) error {
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

func (c *BondCmd) listBonds() ([]*netlink.Bond, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	bonds := []*netlink.Bond{}
	for _, link := range links {
		if bond, ok := link.(*netlink.Bond); ok {
			bonds = append(bonds, bond)
		}
	}

	return bonds, nil
}

func (c *BondCmd) listSlaves(bondname string) ([]netlink.Link, error) {
	bond, err := netlink.LinkByName(bondname)
	if err != nil {
		return nil, err
	}

	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	bondIndex := bond.Attrs().Index

	slaves := []netlink.Link{}
	for _, link := range links {
		if masterIndex := link.Attrs().MasterIndex; masterIndex == bondIndex {
			slaves = append(slaves, link)
		}
	}

	return slaves, nil
}

func stringBondArpValidate(v netlink.BondArpValidate) string {
	switch v {
	case netlink.BOND_ARP_VALIDATE_NONE:
		return "None"
	case netlink.BOND_ARP_VALIDATE_ACTIVE:
		return "Active"
	case netlink.BOND_ARP_VALIDATE_BACKUP:
		return "Backup"
	case netlink.BOND_ARP_VALIDATE_ALL:
		return "All"
	default:
		return fmt.Sprintf("BondArpValidate(%d)", v)
	}
}

func stringBondArpAllTargets(v netlink.BondArpAllTargets) string {
	switch v {
	case netlink.BOND_ARP_ALL_TARGETS_ANY:
		return "Any"
	case netlink.BOND_ARP_ALL_TARGETS_ALL:
		return "All"
	default:
		return fmt.Sprintf("BondArpAllTargets(%d)", v)
	}
}

func stringPrimaryReselect(v netlink.BondPrimaryReselect) string {
	switch v {
	case netlink.BOND_PRIMARY_RESELECT_ALWAYS:
		return "always"
	case netlink.BOND_PRIMARY_RESELECT_BETTER:
		return "better"
	case netlink.BOND_PRIMARY_RESELECT_FAILURE:
		return "filure"
	default:
		return fmt.Sprintf("BondPrimaryReselect(%d)", v)
	}
}

func stringBondFailOverMac(v netlink.BondFailOverMac) string {
	switch v {
	case netlink.BOND_FAIL_OVER_MAC_NONE:
		return "None"
	case netlink.BOND_FAIL_OVER_MAC_ACTIVE:
		return "Active"
	case netlink.BOND_FAIL_OVER_MAC_FOLLOW:
		return "Follow"
	default:
		return fmt.Sprintf("BondFailOverMac(%d)", v)
	}
}

func stringBondAdSelect(v netlink.BondAdSelect) string {
	switch v {
	case netlink.BOND_AD_SELECT_STABLE:
		return "stable"
	case netlink.BOND_AD_SELECT_BANDWIDTH:
		return "bandwith"
	case netlink.BOND_AD_SELECT_COUNT:
		return "count"
	default:
		return fmt.Sprintf("BondAdSelect(%d)", v)
	}
}

func stringLinkAttrs(attrs *netlink.LinkAttrs) string {
	return fmt.Sprintf("%d %s %s %s m:%d p:%d",
		attrs.Index,
		attrs.HardwareAddr,
		attrs.Flags,
		attrs.OperState,
		attrs.MasterIndex,
		attrs.ParentIndex,
	)
}

func (c *BondCmd) showBonds() error {
	bonds, err := c.listBonds()
	if err != nil {
		return err
	}

	for _, bond := range bonds {
		fmt.Printf("%s : %s %s\n", bond.Attrs().Name, bond.Type(), stringLinkAttrs(bond.Attrs()))

		if c.ShowDetail {
			fmt.Printf("mode              : %s\n", bond.Mode)
			fmt.Printf("active slav  e    : %d\n", bond.ActiveSlave)
			fmt.Printf("miimon            : %d\n", bond.Miimon)
			fmt.Printf("up-delay          : %d\n", bond.UpDelay)
			fmt.Printf("down-delay        : %d\n", bond.DownDelay)
			fmt.Printf("use-carrier       : %d\n", bond.UseCarrier)
			fmt.Printf("arp-interval      : %d\n", bond.ArpInterval)
			fmt.Printf("arp-ip-target     : %v\n", bond.ArpIpTargets)
			fmt.Printf("arp-validate      : %s\n", stringBondArpValidate(bond.ArpValidate))
			fmt.Printf("arp-all-targets   : %s\n", stringBondArpAllTargets(bond.ArpAllTargets))
			fmt.Printf("primary           : %d\n", bond.Primary)
			fmt.Printf("primary-reselect  : %s\n", stringPrimaryReselect(bond.PrimaryReselect))
			fmt.Printf("failover-mac      : %s\n", stringBondFailOverMac(bond.FailOverMac))
			fmt.Printf("xmit-hash-policy  : %s\n", bond.XmitHashPolicy)
			fmt.Printf("resend-igmp       : %d\n", bond.ResendIgmp)
			fmt.Printf("num-peer-norif    : %d\n", bond.NumPeerNotif)
			fmt.Printf("allslaves-active  : %d\n", bond.AllSlavesActive)
			fmt.Printf("min-links         : %d\n", bond.MinLinks)
			fmt.Printf("lp-interval       : %d\n", bond.LpInterval)
			fmt.Printf("lacp-rate         : %s\n", bond.LacpRate)
			fmt.Printf("ad-select         : %s\n", stringBondAdSelect(bond.AdSelect))
			if adinfo := bond.AdInfo; adinfo == nil {
				fmt.Printf("ad-info           : nil\n")
			} else {
				fmt.Printf("ad-info.aggregator-id  : %d\n", adinfo.AggregatorId)
				fmt.Printf("ad-info.num-ports      : %d\n", adinfo.NumPorts)
				fmt.Printf("ad-info.actor-key      : %d\n", adinfo.ActorKey)
				fmt.Printf("ad-info.partner-key    : %d\n", adinfo.PartnerKey)
				fmt.Printf("ad-info.partner-mac    : %s\n", adinfo.PartnerMac)
			}
			fmt.Printf("ad-actor-sys-prio : %d\n", bond.AdActorSysPrio)
			fmt.Printf("ad-user-port-key  : %d\n", bond.AdUserPortKey)
			fmt.Printf("ad-actor-system   : %s\n", bond.AdActorSystem)
			fmt.Printf("tib-dynamic-lb    ; %d\n", bond.TlbDynamicLb)
		}

		c.showSlaves(bond.Attrs().Name)
	}

	return nil
}

func (c *BondCmd) showSlaves(bond string) error {
	slaves, err := c.listSlaves(bond)
	if err != nil {
		return err
	}

	for _, slave := range slaves {
		fmt.Printf("%s : %s %s\n", slave.Attrs().Name, slave.Type(), stringLinkAttrs(slave.Attrs()))
		if c.ShowDetail {
			if slaveInfo := slave.Attrs().SlaveInfo; slaveInfo != nil {
				switch si := slaveInfo.(type) {
				case *netlink.BondSlaveInfo:
					fmt.Printf("bond-slave.state             : %s\n", si.State)
					fmt.Printf("bond-slave.mii-status        : %d\n", si.MiiStatus)
					fmt.Printf("bond-slave.link-failure-count: %d\n", si.LinkFailureCount)
					fmt.Printf("bond-slave.parmanent-mac     : %s\n", si.PermanentHwAddr)
					fmt.Printf("bond-slave.queue-id          : %d\n", si.QueueId)
					fmt.Printf("bond-slave.aggregator-id     : %d\n", si.AggregatorId)
					fmt.Printf("bond-slave.actor-oper-oprt-state     : %d\n", si.ActorOperPortState)
					fmt.Printf("bond-slave.ad-partner-oper-port-state: %d\n", si.AdPartnerOperPortState)
				default:
					fmt.Printf("slave-info        : %s %v\n", si.SlaveType(), si)
				}
			} else {
				fmt.Printf("slave-info        : nil\n")
			}
		}
	}

	return nil
}

func bondCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bond",
		Short: "Setting bonding device.",
		Long: `Setting bonding device. 
This command is used on container.`,
	}

	bond := &BondCmd{}

	rootCmd.AddCommand(bond.setBondFlags(
		&cobra.Command{
			Use:   "add-bond <bond>",
			Short: "add bonding device.",
			Long:  "Add bonding device for beluganos l2sw function.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmf *cobra.Command, args []string) error {
				return bond.addBond(args[0])
			},
		},
	))

	rootCmd.AddCommand(bond.setBondFlags(
		&cobra.Command{
			Use:   "del-bond <bond>",
			Short: "Delete bonding device.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmf *cobra.Command, args []string) error {
				return bond.delBond(args[0])
			},
		},
	))

	rootCmd.AddCommand(bond.setFlags(
		&cobra.Command{
			Use:   "add-slaves <bond> <slave ...>",
			Short: "Add slaves to bond.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, arg := range args[1:] {
					if err := bond.addSlave(args[0], arg); err != nil {
						return err
					}
				}
				return nil
			},
		},
	))

	rootCmd.AddCommand(bond.setFlags(
		&cobra.Command{
			Use:   "del-slaves <bond> <slave ...>",
			Short: "Delete slaves from bond.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, arg := range args[1:] {
					if err := bond.delSlave(args[0], arg); err != nil {
						return err
					}
				}
				return nil
			},
		},
	))

	rootCmd.AddCommand(bond.setShowFlags(
		&cobra.Command{
			Use:   "show [bond]",
			Short: "Show bond and slaves,",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return bond.showBonds()
				}

				return bond.showSlaves(args[0])
			},
		},
	))

	return rootCmd
}
