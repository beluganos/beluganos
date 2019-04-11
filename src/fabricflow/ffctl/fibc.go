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
