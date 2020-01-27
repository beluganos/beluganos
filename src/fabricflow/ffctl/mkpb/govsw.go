// -*- coding: utf-8 -*-

// Copyright (C) 2020 Nippon Telegraph and Telephone Corporation.
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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	govswdYamlFile = "govswd.yaml"
	govswdConfFile = "govswd.conf"
)

type GovswCmd struct {
	*Command

	govswdYamlFile string
	govswdConfFile string
}

func NewGovswCmd() *GovswCmd {
	return &GovswCmd{
		Command: NewCommand(),

		govswdYamlFile: govswdYamlFile,
		govswdConfFile: govswdConfFile,
	}
}

func (c *GovswCmd) setConfFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.govswdYamlFile, "govswd-yaml-file", "", govswdYamlFile, "govswd.yaml file name")
	cmd.Flags().StringVarP(&c.govswdConfFile, "govswd-conf-file", "", govswdConfFile, "govswd.conf file name")
	return c.Command.setConfigFlags(cmd)
}

func (c *GovswCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	if err := c.createGovswdConf(playbookName); err != nil {
		return err
	}

	return c.createGovswdYaml(playbookName)
}

func (c *GovswCmd) createGovswdConf(playbookName string) error {
	opt := c.optionConfig()

	path := c.filesPath(playbookName, c.govswdConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookGovswdConf()
	t.FibcAddr = "localhost"
	t.FibcPort = opt.FibcAPIPort

	return t.Execute(f)
}

func (c *GovswCmd) createGovswdYaml(playbookName string) error {
	opt := c.optionConfig()

	path := c.filesPath(playbookName, c.govswdYamlFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookGovswdYaml()
	t.DpID = opt.GoVswdDpID
	for _, name := range c.routerNameList() {
		t.AddPatterns(fmt.Sprintf("%s\\\\.[0-9]+", name))
	}
	return t.Execute(f)
}

func NewGovswCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "govsw",
		Short: "govsw command.",
	}

	govsw := NewGovswCmd()

	rootCmd.AddCommand(govsw.setConfigFlags(
		&cobra.Command{
			Use:   "create",
			Short: "Crate new gonsl config files.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := govsw.readConfig(); err != nil {
					return err
				}
				return govsw.createConf("common")
			},
		},
	))

	return rootCmd
}
