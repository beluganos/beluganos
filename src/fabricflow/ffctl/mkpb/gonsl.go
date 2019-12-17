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
	gonsldYamlFile = "gonsld.yaml"
	gonsldConfFile = "gonsld.conf"
)

type GonslCmd struct {
	*Command

	gonsldYamlFile string
	gonsldConfFile string
}

func NewGonslCmd() *GonslCmd {
	return &GonslCmd{
		Command: NewCommand(),

		gonsldYamlFile: gonsldYamlFile,
		gonsldConfFile: gonsldConfFile,
	}
}

func (c *GonslCmd) setConfFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.gonsldYamlFile, "gonsld-yaml-file", "", gonsldYamlFile, "gonsld.yaml file name")
	cmd.Flags().StringVarP(&c.gonsldConfFile, "gonsld-conf-file", "", gonsldConfFile, "gonsld.conf file name")
	return c.Command.setConfigFlags(cmd)
}

func (c *GonslCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createGonsldYaml(playbookName)
}

func (c *GonslCmd) createGonsldYaml(playbookName string) error {
	g := c.globalConfig()
	opt := c.optionConfig()

	path := c.filesPath(playbookName, c.gonsldYamlFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookGonsldYaml()
	t.DpID = g.DpID
	t.FibcAddr = opt.FibcAPIAddr
	t.FibcPort = opt.FibcAPIPort
	t.L2SWAgingSec = opt.L2SWAgingSec
	t.L2SWSweepSec = opt.L2SWSweepSec
	t.L2SWNotifyLimit = opt.L2SWNotifyLimit
	t.L3PortStart = opt.L3PortStart
	t.L3PortEnd = opt.L3PortEnd
	t.L3VlanBase = opt.L3VlanBase

	return t.Execute(f)
}

func (c *GonslCmd) createGonsldConf(playbookName string) error {
	opt := c.optionConfig()

	path := c.filesPath(playbookName, c.gonsldConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookGonsldConf()
	t.Addr = opt.GonsldListenAddr
	t.Port = opt.GonsldListenPort

	return t.Execute(f)
}

func NewGonslCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gonsl",
		Short: "gonsl command.",
	}

	gonsl := NewGonslCmd()

	rootCmd.AddCommand(gonsl.setConfigFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new gonsl config files.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := gonsl.readConfig(); err != nil {
					return err
				}
				return gonsl.createConf("common")
			},
		},
	))

	return rootCmd
}
