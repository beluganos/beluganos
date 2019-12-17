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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	NetplanFileName = "netplan.yaml"
)

type NetplanCmd struct {
	*Command

	fileName string
}

func NewNetplanCmd() *NetplanCmd {
	return &NetplanCmd{
		Command: NewCommand(),

		fileName: NetplanFileName,
	}
}

func (c *NetplanCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "file-name", "", NetplanFileName, "output file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c *NetplanCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createNetplanConf(playbookName)
}

func (c *NetplanCmd) createNetplanConf(playbookName string) error {
	opt := c.optionConfig()
	g := c.globalConfig()
	r, err := c.routerConfig(playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, c.fileName)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	portmap, err := PortMap(g.DpType)
	if err != nil {
		return err
	}

	log.Debugf("%s created.", path)

	t := NewPlaybookNetplanYaml()
	for _, eth := range ConvToPPortList(r.Eth, portmap) {
		t.AddEth(eth, opt.LXDMtu)
	}
	for eth, vlans := range r.Vlan {
		if HasPPort(eth, portmap) {
			for _, vlan := range vlans {
				t.AddVlan(eth, vlan)
			}
		}
	}

	return t.Execute(f)
}

func NewNetplanCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "netplan",
		Short: "netplan command.",
	}

	netplan := NewNetplanCmd()
	rootCmd.AddCommand(netplan.setConfigFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new netplan config file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := netplan.readConfig(); err != nil {
					return err
				}
				return netplan.createConf(args[0])
			},
		},
	))

	return rootCmd
}
