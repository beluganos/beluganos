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

package maintenance

import (
	"fabricflow/ffctl/fflib"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lxc/lxd/shared/api"
	"github.com/spf13/cobra"
)

type DumpCmdInterface interface {
	dump(io.Writer) error
}

var lxcDumpCmds = []DumpCmdInterface{
	NewFrrCmd(),
	NewSysCmd(),
	NewGoBGPCmd(),
	NewRibxCmd(),
}

type LxcCmd struct {
	nameSep string
	nameAll bool
}

func NewLxcCmd() *LxcCmd {
	return &LxcCmd{}
}

func (c *LxcCmd) setListFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.nameSep, "separator", "s", " ", "separator.")
	cmd.Flags().BoolVarP(&c.nameAll, "all", "a", false, "show all names.")
	return cmd
}

func (c *LxcCmd) dump(w io.Writer) {
	for _, cmd := range lxcDumpCmds {
		cmd.dump(w)
	}
}

func (c *LxcCmd) remoteDump(w io.Writer, name string) {
	fflib.ExecAndOutput(w, "lxc", "exec", name, "--", "ffctl", "maintenance", "lxc", "dump-local")
}

func (c *LxcCmd) names(w io.Writer) error {
	containers, err := fflib.LXDContainerList()
	if err != nil {
		return err
	}

	names := []string{}
	for _, container := range containers {
		if c.nameAll || container.StatusCode == api.Running {
			names = append(names, container.Name)
		}
	}

	fmt.Fprintf(w, "%s", strings.Join(names, c.nameSep))
	return nil
}

func NewLxcCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "lxc",
		Short: "maintenance lxc commands.",
	}

	c := NewLxcCmd()

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "dump <lxc names...>",
			Short: "dump lxc datas",
			Args:  cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				for _, arg := range args {
					c.remoteDump(os.Stdout, arg)
				}
			},
		},
	)

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "dump-local",
			Short: "dump lxc local datas",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				c.dump(os.Stdout)
			},
		},
	)

	rootCmd.AddCommand(c.setListFlags(
		&cobra.Command{
			Use:   "names",
			Short: "shoe lxc names.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.names(os.Stdout)
			},
		},
	))

	return rootCmd
}
