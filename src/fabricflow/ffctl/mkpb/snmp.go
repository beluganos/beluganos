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
	SnmpConfFile  = "snmp.conf"
	SnmpdConfFile = "snmpd.conf"
)

type SnmpCmd struct {
	*Command

	snmpConfFile       string
	snmpdConfFile      string
	snmpproxydConfFile string
}

func NewSnmpCmd() *SnmpCmd {
	return &SnmpCmd{
		Command: NewCommand(),

		snmpConfFile:       SnmpConfFile,
		snmpdConfFile:      SnmpdConfFile,
		snmpproxydConfFile: SnmpproxydConfFile,
	}
}

func (c *SnmpCmd) setSnmpFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.snmpproxydConfFile, "snmpproxyd-conf-file", "", SnmpproxydConfFile, "snmpproxyd.conf file name")
	cmd.Flags().StringVarP(&c.snmpConfFile, "snmp-conf-file", "", SnmpConfFile, "snmp.conf file name")
	cmd.Flags().StringVarP(&c.snmpdConfFile, "snmpd-conf-file", "", SnmpdConfFile, "snmpd.conf file name")
	return c.Command.setConfigFlags(cmd)
}

func (c *SnmpCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	if err := c.createSnmpConf(playbookName); err != nil {
		return err
	}

	if err := c.createSnmpdConf(playbookName); err != nil {
		return err
	}

	return c.createSnmpproxydConf(playbookName)
}

func (c *SnmpCmd) createSnmpConf(playbookName string) error {
	path := c.filesPath(playbookName, c.snmpConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookSnmpConf()

	return t.Execute(f)
}

func (c *SnmpCmd) createSnmpdConf(playbookName string) error {
	opt := c.optionConfig()

	path := c.filesPath(playbookName, c.snmpdConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookSnmpdConf()
	t.Trap2SinkAddr = opt.SnmpproxydAddr
	t.LinkMonitorInterval = opt.SnmpdLinkmonInterval

	return t.Execute(f)
}

func (c *SnmpCmd) createSnmpproxydConf(playbookName string) error {
	path := c.filesPath(playbookName, c.snmpproxydConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookSnmpproxydConf() // proxyd
	t.SnmpproxydAddr = c.optionConfig().SnmpproxydAddr
	return t.Execute(f)
}

func NewSnmpCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "snmp",
		Short: "snmp command.",
	}

	snmp := NewSnmpCmd()
	rootCmd.AddCommand(snmp.setSnmpFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new snmp(d).conf on container file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := snmp.readConfig(); err != nil {
					return err
				}
				return snmp.createConf(args[0])
			},
		},
	))

	return rootCmd
}
