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
	BridgeVlanFileName = "bridge_vlan.yaml"
)

type BridgeVlanCmd struct {
	*Command
	fileName string
}

func NewBridgeVlanCmd() *BridgeVlanCmd {
	return &BridgeVlanCmd{
		Command:  NewCommand(),
		fileName: BridgeVlanFileName,
	}
}

func (c *BridgeVlanCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "file-name", "", BridgeVlanFileName, "output file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c *BridgeVlanCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createBrVlanConf(playbookName)
}

func (c *BridgeVlanCmd) createBrVlanConf(playbookName string) error {
	opt := c.optionConfig()
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

	t := NewPlaybookBrVlanYaml()
	t.Bridge = opt.L2SWBridge
	if l2sw := r.L2SW; l2sw != nil {
		t.AddAccessPorts(l2sw.Access)
		t.AddTrunkPorts(l2sw.Trunk)
	}

	return t.Execute(f)
}

func NewBridgeVlanCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "brvlan",
		Short: "bridge vlan command.",
	}

	brvlan := NewBridgeVlanCmd()
	rootCmd.AddCommand(brvlan.setConfigFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Create new bridge vlan config file",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := brvlan.readConfig(); err != nil {
					return err
				}
				return brvlan.createConf(args[0])
			},
		},
	))

	return rootCmd
}
