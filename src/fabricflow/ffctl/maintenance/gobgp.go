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

var gobgpDumpCommands = []*exec.Cmd{
	exec.Command("gobgp", "policy"),
	exec.Command("gobgp", "neighbor"),
	exec.Command("gobgp", "global"),
	exec.Command("gobgp", "global", "rib", "summary"),
	exec.Command("gobgp", "global", "rib", "-a", "ipv4"),
	exec.Command("gobgp", "global", "rib", "-a", "ipv6"),
	exec.Command("gobgp", "global", "rib", "-a", "vpnv4"),
	exec.Command("gobgp", "global", "rib", "-a", "vpnv6"),
}

type GoBGPCmd struct {
}

func NewGoBGPCmd() *GoBGPCmd {
	return &GoBGPCmd{}
}

func (c *GoBGPCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *GoBGPCmd) dump(w io.Writer) error {
	for _, cmd := range gobgpDumpCommands {
		fmt.Fprintf(w, "> %s\n", strings.Join(cmd.Args, " "))
		fflib.ExecCmdAndOutput(w, cmd)
	}
	return nil
}

func NewGoBGPCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gobgp",
		Short: "maintenance gobgp commands.",
	}

	gobgp := NewGoBGPCmd()

	rootCmd.AddCommand(gobgp.setFlags(
		&cobra.Command{
			Use:   "dump",
			Short: "dump gobgp datas",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return gobgp.dump(os.Stdout)
			},
		},
	))

	return rootCmd
}
