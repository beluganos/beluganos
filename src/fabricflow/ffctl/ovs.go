// -*- coding: utf-8 -*-

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
