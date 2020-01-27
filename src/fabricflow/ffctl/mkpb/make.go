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

package mkpb

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type makeCmfInterface interface {
	setConfig(config *Config)
	createConf(string) error
	copyFlags(*Command)
}

type MakeCmd struct {
	*Command

	withOptions bool
	lxcName     string
}

func NewMakeCmd() *MakeCmd {
	return &MakeCmd{
		Command: NewCommand(),
	}
}

func (c *MakeCmd) setSampleFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolVarP(&c.withOptions, "with-options", "", false, "show options.")
	cmd.Flags().StringVarP(&c.lxcName, "lxc-name", "n", "mic", "container name.")
	return cmd
}

func (c *MakeCmd) run(name string) error {
	cmdList := []makeCmfInterface{
		NewBridgeVlanCmd(),
		NewFibcCmd(),
		NewFrrCmd(),
		NewGobgpCmd(),
		NewLXDCmd(),
		NewNetplanCmd(),
		NewRibtCmd(),
		NewRibxCmd(),
		NewSnmpCmd(),
		NewSysctlCmd(),
		NewNFTableCmd(),
	}

	for _, cmd := range cmdList {
		log.Debugf("%T.run(%s)", cmd, name)
		cmd.copyFlags(c.Command)
		cmd.setConfig(c.config)
		if err := cmd.createConf(name); err != nil {
			return err
		}
	}

	return nil

}

func (c *MakeCmd) all() error {
	names := c.routerNameList()
	if len(names) == 0 {
		return fmt.Errorf("routers not exist.")
	}

	pb := NewPlaybookCmd()
	pb.copyFlags(c.Command)
	pb.setConfig(c.config)
	if err := pb.createConf(names[0]); err != nil {
		return err
	}

	fibs := NewFibsCmd()
	fibs.copyFlags(c.Command)
	fibs.setConfig(c.config)
	if err := fibs.createConf("common"); err != nil {
		return err
	}

	gonsl := NewGonslCmd()
	gonsl.copyFlags(c.Command)
	gonsl.setConfig(c.config)
	if err := gonsl.createConf("common"); err != nil {
		return err
	}

	govsw := NewGovswCmd()
	govsw.copyFlags(c.Command)
	govsw.setConfig(c.config)
	if err := govsw.createConf("common"); err != nil {
		return err
	}

	for _, name := range names {
		if err := c.run(name); err != nil {
			return err
		}
	}

	return pb.createInventory(names...)
}

func (c *MakeCmd) showSample(sampleName string) error {
	config, err := MakeSampleConfig(sampleName, c.lxcName)
	if err != nil {
		return err
	}

	if err := ExecPlaybookSampleMakeYamlTempl(config, os.Stdout); err != nil {
		return err
	}

	if c.withOptions {
		return ExecPlaybookSampleMakeOptionYamlTempl(&config.Option, os.Stdout)
	}

	return nil
}

func NewMakeCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "make",
		Short: "make command.",
	}

	mk := NewMakeCmd()
	rootCmd.AddCommand(mk.setConfigFlags(
		&cobra.Command{
			Use:   "all",
			Short: "make all command",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := mk.readConfig(); err != nil {
					return err
				}
				return mk.all()
			},
		},
	))

	rootCmd.AddCommand(mk.setSampleFlags(
		&cobra.Command{
			Use:   "sample <l3 | l3-vlan | l2sw | iptun>",
			Short: "sample config comamnd",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return mk.showSample(args[0])
			},
		},
	))

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "dp-types",
			Short: "show dp type list",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				for _, s := range PortMapKeys() {
					fmt.Printf("%s\n", s)
				}
			},
		},
	)
	return rootCmd
}
