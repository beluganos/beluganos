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

func playbookDaemonsArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("playbook name not found.")
	}

	for _, arg := range args[1:] {
		if _, ok := playbookFrrDaemons[arg]; !ok {
			return fmt.Errorf("Invalid Daemon. %s", arg)
		}
	}
	return nil
}

func playbookDaemonsCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "daemons",
		Short: "daemons file command.",
	}

	create := NewPlaybookCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name] [daemon names...]",
			Short: "Crate new daemons file.",
			Args:  playbookDaemonsArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createDaemons(args[0], args[1:])
			},
		},
	))

	return rootCmd
}

func playbookFibcYamlCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fibc",
		Short: "fibc command.",
	}

	create := NewPlaybookFibcCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new daemons file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createFibcYaml(args[0])
			},
		},
	))

	return rootCmd

}

func playbookFrrConfCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "frr",
		Short: "frr command.",
	}

	create := NewPlaybookFrrConfCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new frr.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createFrrConf(args[0])
			},
		},
	))

	return rootCmd
}

func playbookGoBGPdConfCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gobgpd",
		Short: "gobgpd command.",
	}

	create := NewPlaybookGoBGPdConfCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new gobgpd.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createGoBGPdConf(args[0])
			},
		},
	))

	return rootCmd
}

func playbookGoBGPConfCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gobgp",
		Short: "gobgo command.",
	}

	create := NewPlaybookCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new gobgp.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createGoBGPConf(args[0])
			},
		},
	))

	return rootCmd
}

func playbookLXDProfileCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "lxd-profile",
		Short: "lxd-profile command.",
	}

	create := NewPlaybookLXDProfileCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new lxd_profile.yml file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createLXDProfile(args[0])
			},
		},
	))

	return rootCmd
}

func playbookNetplanYamlCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "netplan",
		Short: "netplan command.",
	}

	create := NewPlaybookNetplanYamlCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new netplan.yml file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createNetplanYaml(args[0])
			},
		},
	))

	return rootCmd
}

func playbookRibtdConfCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ribtd",
		Short: "ribtd command.",
	}

	create := NewPlaybookCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new ribtd.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createRibtdConf(args[0])
			},
		},
	))

	return rootCmd

}

func playbookSnmpProxydConfCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "snmpproxyd-conf",
		Short: "snmpproxyd-conf command.",
	}

	create := NewPlaybookCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new snmpproxyd.conf on container file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createSnmpProxydConf(args[0])
			},
		},
	))

	return rootCmd
}

func playbookSysctlConfCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sysctl",
		Short: "sysctl command.",
	}

	create := NewPlaybookSysctlConfCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new sysctl.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createSysctlConf(args[0])
			},
		},
	))

	return rootCmd
}

func playbookRibxdConfCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ribxd",
		Short: "ribxd command.",
	}

	create := NewPlaybookRibxdConfCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new ribxd.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createRibxdConf(args[0])
			},
		},
	))

	return rootCmd
}

func playbookInventoryCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "inventory",
		Short: "inventory command.",
	}

	create := NewPlaybookCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new inventory file.",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createInventory(args...)
			},
		},
	))

	return rootCmd
}

func playbookCommonCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "common",
		Short: "common command.",
	}

	create := NewPlaybookCommonCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create",
			Short: "Crate new common files",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.createCommon()
			},
		},
	))

	return rootCmd
}

func playbookCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "playbook",
		Short: "Playbook commands.",
	}

	create := NewPlaybookCmd()
	rootCmd.AddCommand(create.setFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Create new playbook",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.create(args[0])
			},
		},
	))

	rootCmd.AddCommand(
		playbookDaemonsCmd(),
		playbookFibcYamlCmd(),
		playbookFrrConfCmd(),
		playbookGoBGPdConfCmd(),
		playbookGoBGPConfCmd(),
		playbookLXDProfileCmd(),
		playbookNetplanYamlCmd(),
		playbookRibtdConfCmd(),
		playbookSnmpProxydConfCmd(),
		playbookSysctlConfCmd(),
		playbookRibxdConfCmd(),
		playbookInventoryCmd(),
		playbookCommonCmd(),
	)

	return rootCmd
}
