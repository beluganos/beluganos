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
	"sort"

	"github.com/safchain/ethtool"
	"github.com/spf13/cobra"
)

type EthtoolCmd struct {
}

func (c *EthtoolCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *EthtoolCmd) ethtool(f func(*ethtool.Ethtool) error) error {
	et, err := ethtool.NewEthtool()
	if err != nil {
		return err
	}
	defer et.Close()

	return f(et)
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

func (c *EthtoolCmd) changeFeature(ifname string, vals map[string]bool) error {
	return c.ethtool(func(et *ethtool.Ethtool) error {
		return et.Change(ifname, vals)
	})
}

func ethtoolFeatureCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "feature",
		Aliases: []string{"f"},
		Short:   "ethtool features command.",
	}

	et := EthtoolCmd{}

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

	rootCmd.AddCommand(et.setFlags(
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

	rootCmd.AddCommand(et.setFlags(
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

func ethtoolCmd() *cobra.Command {
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
