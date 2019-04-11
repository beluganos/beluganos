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
	"github.com/spf13/cobra"
)

const (
	ovsBridgeDefault = "dp0"
	ofcAddrDefault   = "tcp:127.0.0.1:6633"
)

func doInitOVS(bridge string, ofc string) error {
	if err := execAndOutput("sudo", "ovs-vsctl", "add-br", bridge); err != nil {
		return err
	}

	if err := execAndOutput("sudo", "ovs-vsctl", "set-controller", bridge, ofc); err != nil {
		return err
	}

	return nil
}

func doCleanOVS(bridge string) error {
	return execAndOutput("sudo", "ovs-vsctl", "del-br", bridge)
}

func doStatusOVS(bridge string) error {
	return execAndOutput("sudo", "ovs-vsctl", "show")
}

func ovsCmd() *cobra.Command {
	var bridge string
	var ofcAddr string

	rootCmd := &cobra.Command{
		Use:   "ovs",
		Short: "OVS commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doStatusOVS(bridge)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&bridge, "bridge", "", ovsBridgeDefault, "ovs-bridge name.")
	rootCmd.PersistentFlags().StringVarP(&ofcAddr, "controller", "", ofcAddrDefault, "ovs-controller address.")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize OVS",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInitOVS(bridge, ofcAddr)
		},
	}

	rootCmd.AddCommand(initCmd)

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean OVS",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doCleanOVS(bridge)
		},
	}

	rootCmd.AddCommand(cleanCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doStatusOVS(bridge)
		},
	}

	rootCmd.AddCommand(statusCmd)

	return rootCmd
}
