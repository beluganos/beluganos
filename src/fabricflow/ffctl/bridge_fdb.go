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
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type BridgeFdb struct {
	link netlink.Link
	fdb  *BridgeFdbCLI
}

func (f *BridgeFdb) String() string {
	return fmt.Sprintf("%s %d %s", f.fdb.Mac, f.fdb.Vlan, f.fdb.Dev)
}

func newFDBEntry(hwaddr net.HardwareAddr, vid uint16, link netlink.Link) *netlink.Neigh {
	return &netlink.Neigh{
		LinkIndex:    link.Attrs().Index,
		IP:           nil,
		Family:       unix.AF_BRIDGE,
		HardwareAddr: hwaddr,
		Vlan:         int(vid),
		Flags:        unix.NTF_MASTER,
		State:        unix.NUD_PERMANENT,
	}
}

func parseFDBEntry(mac, vlan, ifname string) (*netlink.Neigh, error) {
	vid, err := strconv.ParseUint(vlan, 0, 16)
	if err != nil {
		return nil, err
	}

	hwaddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	link, err := netlink.LinkByName(ifname)
	if err != nil {
		return nil, err
	}

	return newFDBEntry(hwaddr, uint16(vid), link), nil
}

func fdbList(all bool) ([]*BridgeFdb, error) {

	fdbCLIs, err := ExecBridgeFdbShow()
	if err != nil {
		return nil, err
	}

	ifmap := map[string]netlink.Link{}
	ifmac := map[string]struct{}{}
	fdbs := []*BridgeFdb{}

	for _, fdbCLI := range fdbCLIs {
		link, ok := ifmap[fdbCLI.Dev]
		if !ok {
			// Get Link from kernel.
			var err error
			if link, err = netlink.LinkByName(fdbCLI.Dev); err != nil {
				continue
			}

			// add link to ifmap, and add if-mac to ifmac.
			ifmap[fdbCLI.Dev] = link
			ifmac[link.Attrs().HardwareAddr.String()] = struct{}{}
		}

		if !all {
			if _, ok := ifmac[fdbCLI.Mac]; ok {
				continue
			}

			if mcaddr := isMulticastHardwareAddr(fdbCLI.GetHardwareAddr()); mcaddr {
				continue
			}
		}

		fdbs = append(fdbs, &BridgeFdb{link: link, fdb: fdbCLI})
	}

	return fdbs, nil
}

type BridgeFdbCmd struct {
	useLocal bool
}

func (c *BridgeFdbCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.useLocal, "include-local", "", false, "include local entries.")
	return cmd
}

func (c *BridgeFdbCmd) show() error {
	fdbs, err := fdbList(c.useLocal)
	if err != nil {
		return err
	}

	for _, fdb := range fdbs {
		fmt.Printf("%s dev %s vlan %d %s\n",
			fdb.fdb.Mac,
			fdb.fdb.Dev,
			fdb.fdb.Vlan,
			fdb.fdb.State,
		)
	}

	return nil
}

func (c *BridgeFdbCmd) add(hwaddr, vid, ifname string) error {
	fdb, err := parseFDBEntry(hwaddr, vid, ifname)
	if err != nil {
		return err
	}

	if err := netlink.NeighAppend(fdb); err != nil {
		return err
	}

	return nil
}

func (c *BridgeFdbCmd) del(hwaddr, vid, ifname string) error {
	fdb, err := parseFDBEntry(hwaddr, vid, ifname)
	if err != nil {
		return err
	}

	if err := netlink.NeighDel(fdb); err != nil {
		return err
	}

	return nil
}

func (c *BridgeFdbCmd) clear() error {
	fdbs, err := fdbList(c.useLocal)
	if err != nil {
		return err
	}

	for _, fdb := range fdbs {
		neigh := newFDBEntry(fdb.fdb.GetHardwareAddr(), fdb.fdb.Vlan, fdb.link)
		if err := netlink.NeighDel(neigh); err != nil {
			fmt.Printf("%s %s\n", err, fdb)
		}
	}

	return nil
}

func bridgeFdbCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fdb",
		Short: "bridge fdb command.",
	}

	fdb := &BridgeFdbCmd{}
	rootCmd.AddCommand(fdb.setFlags(
		&cobra.Command{
			Use:   "show",
			Short: "show fdb",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return fdb.show()
			},
		},
	))

	rootCmd.AddCommand(fdb.setFlags(
		&cobra.Command{
			Use:   "add <hwaddr> <vid> <ifname>",
			Short: "add fdb entry.",
			Args:  cobra.ExactArgs(3),
			RunE: func(cmd *cobra.Command, args []string) error {
				return fdb.add(args[0], args[1], args[2])
			},
		},
	))

	rootCmd.AddCommand(fdb.setFlags(
		&cobra.Command{
			Use:   "del <hwaddr> <vid> <ifname>",
			Short: "delete fdb entry.",
			Args:  cobra.ExactArgs(3),
			RunE: func(cmd *cobra.Command, args []string) error {
				return fdb.del(args[0], args[1], args[2])
			},
		},
	))

	rootCmd.AddCommand(fdb.setFlags(
		&cobra.Command{
			Use:   "clear",
			Short: "delete all fdb entry.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return fdb.clear()
			},
		},
	))

	return rootCmd
}
