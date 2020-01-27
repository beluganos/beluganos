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
	GobgpConfFile  = "gobgp.conf"
	GobgpdConfFile = "gobgpd.conf"
)

type GobgpCmd struct {
	*Command

	gobgpConfFile  string
	gobgpdConfFile string
}

func NewGobgpCmd() *GobgpCmd {
	return &GobgpCmd{
		Command: NewCommand(),

		gobgpConfFile:  GobgpConfFile,
		gobgpdConfFile: GobgpdConfFile,
	}
}

func (c *GobgpCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.gobgpConfFile, "gobgp-conf", "", GobgpConfFile, "gobgpd service config")
	cmd.Flags().StringVarP(&c.gobgpdConfFile, "gobgpd-conf", "", GobgpdConfFile, "gobgpd config")
	return c.Command.setConfigFlags(cmd)
}

func (c *GobgpCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	if err := c.createGoBGPdConf(playbookName); err != nil {
		return err
	}

	return c.createGoBGPConf(playbookName)
}

func (c *GobgpCmd) createGoBGPConf(playbookName string) error {
	opt := c.optionConfig()

	path := c.filesPath(playbookName, c.gobgpConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookGoBGPConf()
	t.APIAddr = opt.GoBGPAPIAddr
	t.APIPort = opt.GoBGPAPIPort
	return t.Execute(f)
}

func (c *GobgpCmd) createGoBGPdConf(playbookName string) error {
	opt := c.optionConfig()
	g := c.globalConfig()

	r, err := c.routerConfig(playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, c.gobgpdConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	zapiEnabled := func() bool {
		if g.Vpn {
			return r.NodeID != 0 // VPN-MIC => false, VPN-RIC => true
		}
		return opt.GoBGPZAPIEnable
	}()

	t := NewPlaybookGoBGPdConf()
	t.RouterID = g.ReID
	t.AS = opt.GoBGPAs
	t.ZAPIVersion = opt.GoBGPZAPIVersion
	t.ZAPIEnable = zapiEnabled
	return t.Execute(f)
}

func NewGoBGPCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gobgp",
		Short: "gobgp command.",
	}

	gobgp := NewGobgpCmd()
	rootCmd.AddCommand(gobgp.setConfigFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new gobgp config file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := gobgp.readConfig(); err != nil {
					return err
				}
				return gobgp.createConf(args[0])
			},
		},
	))

	return rootCmd
}
