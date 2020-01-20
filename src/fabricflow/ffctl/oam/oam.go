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

package oam

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "oam",
		Short: "oam command.",
	}

	rootCmd.AddCommand(
		NewAuditCommand(),
	)

	return rootCmd
}

func NewAuditCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "audit",
		Short: "audit command",
	}

	audit := NewAuditCmd()

	rootCmd.AddCommand(audit.setFlags(
		&cobra.Command{
			Use:     "route-cnt",
			Short:   "audit route count,",
			Aliases: []string{"rc"},
			RunE: func(cmd *cobra.Command, args []string) error {
				return audit.routeCnt()
			},
		},
	))

	return rootCmd
}
