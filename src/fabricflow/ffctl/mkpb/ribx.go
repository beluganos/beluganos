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
	RibxFileName = "ribxd.conf"
)

type RibxCmd struct {
	*Command

	fileName string
}

func NewRibxCmd() *RibxCmd {
	return &RibxCmd{
		Command: NewCommand(),

		fileName: RibxFileName,
	}
}

func (c *RibxCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "file-name", "", RibxFileName, "outpur file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c *RibxCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createRibxConf(playbookName)
}

func (c *RibxCmd) createRibxConf(playbookName string) error {
	opt := c.optionConfig()
	g := c.globalConfig()
	r, err := c.routerConfig(playbookName)
	if err != nil {
		return err
	}
	mic, err := c.micConfig()
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

	t := NewPlaybookRibxdConf()
	t.NID = r.NodeID
	t.ReID = g.ReID
	t.Name = playbookName
	t.RT = r.RT()
	t.RD = r.RD()
	t.Vpn = g.Vpn
	t.Mic = mic.Name

	t.NLARecvChanSize = opt.NLARecvChannelSize
	t.NLARecvSockBufSize = opt.NetlinkSocketBufSize
	t.NLABrVlanUpdateSec = opt.NLABrVlanUpdateSec
	t.NLABrVlanChanSize = opt.NLABrVlanChanSize
	t.VpnNexthop = opt.VPMNexthopNetwork
	t.VpnNexhopBridge = opt.VPNPseudoBridge
	t.NLACorePort = opt.NLACorePort
	t.NLAAPIPort = opt.NLAAPIPort
	t.FibcAPIAddr = opt.FibcAPIAddr
	t.FibcAPIPort = opt.FibcAPIPort
	t.RibsCorePort = opt.RibsCorePort
	t.RibsAPIPort = opt.RibsAPIPort
	t.RibpAPIPort = opt.RibpAPIPort

	t.LogLevel = opt.RibxLogLevel
	t.LogDump = opt.RibxLogDump

	return t.Execute(f)
}

func NewRibxCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ribX",
		Aliases: []string{"ribx", "rib"},
		Short:   "ribX command.",
	}

	ribx := NewRibxCmd()
	rootCmd.AddCommand(ribx.setConfigFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new ribxd.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := ribx.readConfig(); err != nil {
					return err
				}
				return ribx.createConf(args[0])
			},
		},
	))

	return rootCmd
}
