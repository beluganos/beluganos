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
	"encoding/hex"
	"fmt"
	"os"
	"time"

	dhcplib "fabricflow/ffctl/dhcp"

	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	DHCPIfName      = "eth0"
	DHCPRetryMax    = 3
	DHCPRetrySecond = 1
)

var (
	DHCPIfNames     = []string{"eth0", "ens3", "ma1"}
	DHCPIfBlacklist = []string{}
	DHCPIfPatterns  = []string{"eth[0-9]+", "ens[0-9]+", "ma[0-9]+"}
)

type DhcpCmd struct {
	OptCode uint8
	OptName string

	IfNames     []string
	IfBlacklist []string
	IfPatterns  []string

	Timeout   time.Duration
	RetryMax  uint8
	RetryTime time.Duration
}

func (c *DhcpCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringSliceVarP(&c.IfNames, "interface-list", "I", DHCPIfNames, "interface list.")
	cmd.Flags().StringSliceVarP(&c.IfBlacklist, "interface-blacklist", "B", DHCPIfBlacklist, "interface blacklist")
	cmd.Flags().StringSliceVarP(&c.IfPatterns, "interface-patterns", "P", DHCPIfPatterns, "interface patterns.")
	cmd.Flags().Uint8VarP(&c.RetryMax, "retry-max", "", DHCPRetryMax, "retry max count.")
	cmd.Flags().DurationVarP(&c.RetryTime, "retry-interval", "", DHCPRetrySecond*time.Second, "retry interval.")
	cmd.Flags().DurationVarP(&c.Timeout, "dhcp-timeout", "", client4.DefaultReadTimeout, "dhcp timeout.")

	cmd.Flags().Uint8VarP(&c.OptCode, "option-code", "c", 0, "option code.")
	cmd.Flags().StringVarP(&c.OptName, "option-name", "n", "", "option name.")

	return cmd
}

func (c *DhcpCmd) doGetOption(client *client4.Client, option *dhcplib.IPv4Option, ifnames []string) error {
	msg, err := dhcplib.IPv4Exchanges(client, option.OptionCode(), ifnames)
	if err != nil {
		log.Debugf("%s", err)
		return err
	}

	b := msg.Options.Get(option.OptionCode())
	if b == nil {
		return fmt.Errorf("option not found.")
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("%s", hex.Dump(b))
	}

	s, err := option.Decode(b)
	if err != nil {
		return err
	}

	fmt.Printf("%s", s)
	return nil
}

func (c *DhcpCmd) getOption() error {
	var (
		option *dhcplib.IPv4Option
		err    error
	)

	if c.OptCode != 0 {
		option, err = dhcplib.GetIPv4OptionByCode(c.OptCode)
	} else {
		option, err = dhcplib.GetIPv4Option(c.OptName)
	}
	if err != nil {
		log.Debugf("%s", err)
		return err
	}

	ifnames := dhcplib.IfNameList(c.IfPatterns, c.IfNames, c.IfBlacklist)

	client := client4.NewClient()

	for index := uint8(0); index < c.RetryMax; index++ {
		err := c.doGetOption(client, option, ifnames)
		if err == nil {
			return nil
		}

		log.Debugf("%s", err)
		time.Sleep(c.RetryTime)
	}

	log.Errorf("Retry overflow.")
	return fmt.Errorf("Retry overflow.")
}

type DhcpToolsCmd struct {
}

func (c *DhcpToolsCmd) showOptions() error {
	dhcplib.ListIPv4Options(func(name string, option *dhcplib.IPv4Option) {
		fmt.Printf("%s (%d)\n", name, option.Code)
	})
	return nil
}

func (c *DhcpToolsCmd) showDhcpdConf(addr string) error {
	t, err := dhcplib.NewISCDhcpdConfig(addr)
	if err != nil {
		return err
	}

	return t.Execute(os.Stdout)
}

func dhcpToolsCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tools",
		Short: "dhcp tools command.",
	}

	dhcp := DhcpToolsCmd{}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "options-list",
			Short: "show option list.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return dhcp.showOptions()
			},
		},
	)

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "dhcpd-config <server addr/prefix-len>",
			Short: "show sample isc-dhcpd config",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return dhcp.showDhcpdConf(args[0])
			},
		},
	)

	return rootCmd
}

func dhcpIPv4OptionCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "option",
		Short: "dhcp(IPv4) option command.",
	}

	dhcp := DhcpCmd{}

	rootCmd.AddCommand(dhcp.setFlags(
		&cobra.Command{
			Use:     "show",
			Short:   "show dhcp option.",
			Aliases: []string{"dump", "get"},
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return dhcp.getOption()
			},
		},
	))

	return rootCmd
}

func dhcpIPv4Cmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dhcp",
		Short: "dhcp(IPv4) command.",
	}

	rootCmd.AddCommand(
		dhcpIPv4OptionCmd(),
		dhcpToolsCmd(),
	)

	return rootCmd
}
