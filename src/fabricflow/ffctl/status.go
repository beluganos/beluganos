// -*- coding: utf-8 -*-

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
