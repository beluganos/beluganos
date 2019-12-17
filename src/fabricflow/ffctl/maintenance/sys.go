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
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var sysDumpCommands = []*exec.Cmd{
	exec.Command("hostname"),
	exec.Command("date"),
	exec.Command("sysctl", "-a"),
	exec.Command("ip", "link"),
	exec.Command("ip", "addr"),
	exec.Command("ip", "-4", "neigh"),
	exec.Command("ip", "-6", "neigh"),
	exec.Command("ip", "-4", "route"),
	exec.Command("ip", "-6", "route"),
	exec.Command("ip", "-f", "mpls", "route"),
	exec.Command("bridge", "link", "show"),
	exec.Command("bridge", "vlan", "show"),
	exec.Command("bridge", "fdb", "show"),
}

type SysCmd struct {
}

func NewSysCmd() *SysCmd {
	return &SysCmd{}
}

func (c *SysCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *SysCmd) dump(w io.Writer) error {
	for _, cmd := range sysDumpCommands {
		fmt.Fprintf(w, "> %s\n", strings.Join(cmd.Args, " "))
		fflib.ExecCmdAndOutput(w, cmd)
	}
	return nil
}

func NewSysCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sys",
		Short: "maintenance sys commands.",
	}

	sys := NewSysCmd()

	rootCmd.AddCommand(sys.setFlags(
		&cobra.Command{
			Use:   "dump",
			Short: "dump sys datas",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return sys.dump(os.Stdout)
			},
		},
	))

	return rootCmd
}
