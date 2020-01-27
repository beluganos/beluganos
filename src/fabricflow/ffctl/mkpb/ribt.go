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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	RibtFileName = "ribtd.conf"
)

type RibtCmd struct {
	*Command

	fileName string
}

func NewRibtCmd() *RibtCmd {
	return &RibtCmd{
		Command: NewCommand(),

		fileName: RibtFileName,
	}
}

func (c *RibtCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "file-name", "", RibtFileName, "output file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c *RibtCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createRibtConf(playbookName)
}

func (c *RibtCmd) createRibtConf(playbookName string) error {
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

	t := NewPlaybookRibtdConf()
	t.GoBGPAPIAddr = opt.GoBGPAPIAddr
	t.GoBGPAPIPort = opt.GoBGPAPIPort
	if iptun := r.IPTun; iptun != nil {
		t.BgpRouteFamily = iptun.BgpRouteFamily
		t.TunnelLocalRange4 = iptun.LocalAddrRange4
		t.TunnelLocalRange6 = iptun.LocalAddrRange6
	}
	t.TunnelIFPrefix = opt.IPTunnelIFPrefix
	t.TunnelTypeIPv6 = opt.IPTunnelTypeIPv6
	t.DumpTableDuration = time.Duration(opt.RibtDumpSec) * time.Second

	return t.Execute(f)
}

func NewRibtCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ribt",
		Short: "ribt command.",
	}

	ribt := NewRibtCmd()
	rootCmd.AddCommand(ribt.setConfigFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new ribtd.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := ribt.readConfig(); err != nil {
					return err
				}
				return ribt.createConf(args[0])
			},
		},
	))

	return rootCmd

}
