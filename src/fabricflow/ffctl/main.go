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
	"fabricflow/ffctl/bonding"
	"fabricflow/ffctl/bridge"
	"fabricflow/ffctl/container"
	dhcplib "fabricflow/ffctl/dhcp"
	"fabricflow/ffctl/ethtools"
	"fabricflow/ffctl/maintenance"
	"fabricflow/ffctl/mkpb"
	"fabricflow/ffctl/monitor"
	"fabricflow/ffctl/oam"
	"fabricflow/ffctl/ovs"
	"fabricflow/ffctl/service"
	"fabricflow/ffctl/vlan"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	cfgPathDefault = "/etc/beluganos/fibc.conf"
)

func rootCmd(name string) *cobra.Command {
	var verbose bool
	var showCompletion bool

	rootCmd := &cobra.Command{
		Use:   name,
		Short: "Beluganos control command.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				log.SetLevel(log.DebugLevel)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if showCompletion {
				cmd.GenBashCompletion(os.Stdout)
			} else {
				cmd.Usage()
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show detail messages.")
	rootCmd.PersistentFlags().BoolVar(&showCompletion, "show-completion", false, "Show bash-comnpletion")

	rootCmd.AddCommand(
		bonding.NewCmd(),
		bridge.NewCmd(),
		container.NewCmd(),
		dhcplib.NewIPv4Cmd(),
		ethtools.NewCmd(),
		mkpb.NewCmd(),
		monitor.NewCmd(),
		maintenance.NewCmd(),
		oam.NewCmd(),
		ovs.NewCmd(),
		service.NewCmd(),
		vlan.NewCmd(),
	)

	return rootCmd
}

func main() {
	if err := rootCmd("ffctl").Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
