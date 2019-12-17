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

var ribxDumpCommands = []*exec.Cmd{
	exec.Command("nlac"),
	exec.Command("ribsc", "show", "nexthop-map"),
	exec.Command("ribsc", "show", "nexthops"),
	exec.Command("ribsc", "show", "rics"),
	exec.Command("journalctl", "-t", "nlad"),
	exec.Command("journalctl", "-t", "ribcd"),
	exec.Command("journalctl", "-t", "ribpd"),
	exec.Command("journalctl", "-t", "ribtd"),
	exec.Command("journalctl", "-t", "ribsd"),
}

type RibxCmd struct {
}

func NewRibxCmd() *RibxCmd {
	return &RibxCmd{}
}

func (c *RibxCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *RibxCmd) dump(w io.Writer) error {
	for _, cmd := range ribxDumpCommands {
		fmt.Fprintf(w, "> %s\n", strings.Join(cmd.Args, " "))
		fflib.ExecCmdAndOutput(w, cmd)
	}
	return nil
}

func NewRibxCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ribx",
		Short: "maintenance ribx commands.",
	}

	ribx := NewRibxCmd()

	rootCmd.AddCommand(ribx.setFlags(
		&cobra.Command{
			Use:   "dump",
			Short: "dump ribx datas",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return ribx.dump(os.Stdout)
			},
		},
	))

	return rootCmd
}
