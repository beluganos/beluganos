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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	FibcFileName = "fibc.yml"
)

type FibcCmd struct {
	*Command

	fileName string
}

func NewFibcCmd() *FibcCmd {
	return &FibcCmd{
		Command: NewCommand(),

		fileName: FibcFileName,
	}
}

func (c *FibcCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "file-name", "", FibcFileName, "output file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c *FibcCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createFibcConf(playbookName)
}

func (c *FibcCmd) createFibcConf(playbookName string) error {
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

	log.Debugf("%s created.", path)

	portmap, err := PortMap(g.DpType)
	if err != nil {
		return err
	}

	dpName := fmt.Sprintf("dp_%d", g.DpID)

	t := NewPlaybookFibcYaml(g.ReID)
	t.DpName = dpName
	t.DpID = g.DpID
	t.DpMode = g.DpMode
	t.AddPorts(ConvToPortMap(r.Eth, portmap))
	return t.Execute(f)
}

func NewFibcCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fibc",
		Short: "fibc command.",
	}

	fibc := NewFibcCmd()
	rootCmd.AddCommand(fibc.setConfigFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new fibc config file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := fibc.readConfig(); err != nil {
					return err
				}
				return fibc.createConf(args[0])
			},
		},
	))

	return rootCmd

}
