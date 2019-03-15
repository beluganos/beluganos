// -*- coding: utf-8 -*-

package main

import (
	"github.com/spf13/cobra"
)

const (
	fibcDaemonName = "fibcd"
)

func fibcCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fibc <start, stop, status>",
		Short: "fibc control",
		Run: func(cmd *cobra.Command, args []string) {
			execAndOutput("systemctl", "status", fibcDaemonName)
		},
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start fibcd",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return execAndOutput("sudo", "systemctl", "start", fibcDaemonName)
		},
	}

	rootCmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop fibcd",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return execAndOutput("sudo", "systemctl", "stop", fibcDaemonName)
		},
	}

	rootCmd.AddCommand(stopCmd)

	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart fibcd",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			execAndOutput("sudo", "systemctl", "stop", fibcDaemonName)
			return execAndOutput("sudo", "systemctl", "start", fibcDaemonName)
		},
	}

	rootCmd.AddCommand(restartCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show status",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			execAndOutput("systemctl", "status", fibcDaemonName)
		},
	}

	rootCmd.AddCommand(statusCmd)

	return rootCmd
}
