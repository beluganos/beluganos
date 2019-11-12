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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//
// Command is rot command.
//
type Command struct {
	Verbose    bool
	Completion bool
}

//
// NewCommand returns new command.
//
func NewCommand() *Command {
	return &Command{}
}

func (c *Command) execute(name string) error {
	rootCmd := &cobra.Command{
		Use:   name,
		Short: "FIBC command.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if c.Verbose {
				log.SetLevel(log.DebugLevel)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if c.Completion {
				cmd.GenBashCompletion(os.Stdout)
			} else {
				cmd.Usage()
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(
		&c.Verbose, "verbose", "v", false, "Show detail messages.")
	rootCmd.PersistentFlags().BoolVar(
		&c.Completion, "show-completion", false, "Show bash-comnpletion")

	rootCmd.AddCommand(
		dpAPICmd(),
		apAPICmd(),
		vmAPICmd(),
		vsAPICmd(),
		dbCmd(),
		monCmd(),
		testCmd(),
	)

	return rootCmd.Execute()
}

func main() {
	if err := NewCommand().execute("fibcctl"); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
