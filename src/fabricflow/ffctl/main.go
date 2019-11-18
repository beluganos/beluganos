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
	"os"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
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
		containerCmd(),
		statusCmd(),
		fibcCmd(),
		ovsCmd(),
		playbookCmd(),
		vlanCmd(),
		bridgeCmd(),
		bondCmd(),
		MonitorCmd(),
	)

	return rootCmd
}

func main() {
	if err := rootCmd("ffctl").Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
