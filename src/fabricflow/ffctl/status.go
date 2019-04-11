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

	"github.com/spf13/cobra"
)

func doStatusCmd(args []string) {
	if len(args) == 0 {
		execAndOutput("systemctl", "status", "--no-pager", "fibcd")
		execAndOutput("lxc", "list")
		execAndOutput("sudo", "ovs-vsctl", "show")

		return
	}

	for _, arg := range args {
		switch arg {
		case "fibc":
			execAndOutput("systemctl", "status", "--no-pager", "fibcd")
		case "lxc":
			execAndOutput("lxc", "list")
		case "ovs":
			execAndOutput("sudo", "ovs-vsctl", "show")
		}
	}
}

func statusCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "status [fibc, lxc, ovs]",
		Short: "Show status",
		Args: func(cmd *cobra.Command, args []string) error {
			arr := []string{"fibc", "lxc", "ovs"}
			for _, arg := range args {
				if index := indexOf(arg, arr); index < 0 {
					return fmt.Errorf("Invalid target. %s", arg)
				}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			doStatusCmd(args)
		},
	}

	return rootCmd
}
