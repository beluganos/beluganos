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

package ethtools

import (
	"fmt"
	"net"
	"regexp"
	"sort"

	"github.com/safchain/ethtool"
	"github.com/spf13/cobra"
)

type EthtoolCmd struct {
	excludes []string
}

func NewEthtoolCmd() *EthtoolCmd {
	return &EthtoolCmd{
		excludes: []string{},
	}
}

func (c *EthtoolCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *EthtoolCmd) setChangeFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringSliceVarP(&c.excludes, "exclude", "e", []string{}, "exclude ifname patterns (comma separated or multiple)")
	return c.setFlags(cmd)
}

func (c *EthtoolCmd) ethtool(f func(*ethtool.Ethtool) error) error {
	et, err := ethtool.NewEthtool()
	if err != nil {
		return err
	}
	defer et.Close()

	return f(et)
}

func (c *EthtoolCmd) listInterfaces(pattern string, excludes []string) ([]net.Interface, error) {

	patternRegex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	excludeRegexes := []*regexp.Regexp{}
	for _, exclude := range excludes {
		re, err := regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
		excludeRegexes = append(excludeRegexes, re)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	targetIfaces := []net.Interface{}
FOR_LOOP:
	for _, iface := range ifaces {
		for _, re := range excludeRegexes {
			if match := re.MatchString(iface.Name); match {
				continue FOR_LOOP
			}

		}

		if match := patternRegex.MatchString(iface.Name); match {
			targetIfaces = append(targetIfaces, iface)
			continue FOR_LOOP
		}
	}

	return targetIfaces, nil
}

func (c *EthtoolCmd) listFeatures(ifname string) error {
	return c.ethtool(func(et *ethtool.Ethtool) error {
		features, err := et.Features(ifname)
		if err != nil {
			return err
		}

		names := []string{}
		for name := range features {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			fmt.Printf("%s = %t\n", name, features[name])
		}

		return nil
	})
}

func (c *EthtoolCmd) changeFeature(pattern string, vals map[string]bool) error {
	ifaces, err := c.listInterfaces(pattern, c.excludes)
	if err != nil {
		return err
	}

	return c.ethtool(func(et *ethtool.Ethtool) error {
		for _, iface := range ifaces {
			fmt.Printf("chage featurs. %s\n", iface.Name)
			if err := et.Change(iface.Name, vals); err != nil {
				fmt.Printf("change feature error. %s %s\n", iface.Name, err)
			}
		}
		return nil
	})
}

func ethtoolFeatureCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "feature",
		Aliases: []string{"f"},
		Short:   "ethtool features command.",
	}

	et := NewEthtoolCmd()

	rootCmd.AddCommand(et.setFlags(
		&cobra.Command{
			Use:     "list <interface name>",
			Aliases: []string{"ls"},
			Args:    cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, arg := range args {
					if err := et.listFeatures(arg); err != nil {
						return err
					}
				}
				return nil
			},
		},
	))

	rootCmd.AddCommand(et.setChangeFlags(
		&cobra.Command{
			Use:  "on <interface name> <features...>",
			Args: cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				vals := map[string]bool{}
				for _, ifname := range args[1:] {
					vals[ifname] = true
				}
				return et.changeFeature(args[0], vals)
			},
		},
	))

	rootCmd.AddCommand(et.setChangeFlags(
		&cobra.Command{
			Use:  "off <interface name> <features...>",
			Args: cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				vals := map[string]bool{}
				for _, ifname := range args[1:] {
					vals[ifname] = false
				}
				return et.changeFeature(args[0], vals)
			},
		},
	))

	return rootCmd
}

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ethtool",
		Aliases: []string{"et"},
		Short:   "ethtool command.",
	}

	rootCmd.AddCommand(
		ethtoolFeatureCmd(),
	)
	return rootCmd
}
