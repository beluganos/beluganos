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
	SysctlFileName = "sysctl.conf"
)

type SysctlCmd struct {
	*Command

	fileName string
}

func NewSysctlCmd() *SysctlCmd {
	return &SysctlCmd{
		Command: NewCommand(),

		fileName: SysctlFileName,
	}
}

func (c *SysctlCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "file-name", "", SysctlFileName, "output file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c *SysctlCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createSysctlConf(playbookName)
}

func (c *SysctlCmd) createSysctlConf(playbookName string) error {
	g := c.globalConfig()
	r, err := c.routerConfig(playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, "sysctl.conf")
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

	t := NewPlaybookSysctlConf()
	t.SockBufSize = c.optionConfig().NetlinkSocketBufSize
	t.MplsLabel = c.optionConfig().SysctlMPLSLabelMax

	for _, eth := range ConvToPPortList(r.Eth, portmap) {
		t.AddIface(eth, 0)
	}
	for eth, vlans := range r.Vlan {
		if HasPPort(eth, portmap) {
			for _, vlan := range vlans {
				t.AddIface(eth, vlan)
			}
		}
	}

	return t.Execute(f)
}

func NewSysctlCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sysctl",
		Short: "sysctl command.",
	}

	sysctl := NewSysctlCmd()
	rootCmd.AddCommand(sysctl.setConfigFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new sysctl.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := sysctl.readConfig(); err != nil {
					return err
				}
				return sysctl.createConf(args[0])
			},
		},
	))

	return rootCmd
}
