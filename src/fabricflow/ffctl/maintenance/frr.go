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

var frrDumpCommands = []*exec.Cmd{
	exec.Command("vtysh", "-c", "show running-config"),
	exec.Command("vtysh", "-c", "show interface"),
	exec.Command("vtysh", "-c", "show ip route"),
	exec.Command("vtysh", "-c", "show ipv6 route"),
	exec.Command("vtysh", "-c", "show ip ospf neighbor"),
	exec.Command("vtysh", "-c", "show ipv6 ospf6 neighbor"),
	exec.Command("vtysh", "-c", "show ip ospf route"),
	exec.Command("vtysh", "-c", "show ipv6 ospf6 route"),
	exec.Command("vtysh", "-c", "show ip ospf database"),
	exec.Command("vtysh", "-c", "show ip6 ospf6 database"),
	exec.Command("vtysh", "-c", "show mpls ldp binding"),
	exec.Command("vtysh", "-c", "show mpls ldp discovery"),
	exec.Command("vtysh", "-c", "show mpls ldp interface"),
	exec.Command("vtysh", "-c", "show mpls ldp neighbor"),
}

type FrrCmd struct {
}

func NewFrrCmd() *FrrCmd {
	return &FrrCmd{}
}

func (c *FrrCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *FrrCmd) dump(w io.Writer) error {
	for _, cmd := range frrDumpCommands {
		fmt.Fprintf(w, "> %s\n", strings.Join(cmd.Args, " "))
		fflib.ExecCmdAndOutput(w, cmd)
	}
	return nil
}

func NewFrrCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "frr",
		Short: "maintenance frr commands.",
	}

	frr := NewFrrCmd()

	rootCmd.AddCommand(frr.setFlags(
		&cobra.Command{
			Use:   "dump",
			Short: "dump frr datas",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return frr.dump(os.Stdout)
			},
		},
	))

	return rootCmd

}
