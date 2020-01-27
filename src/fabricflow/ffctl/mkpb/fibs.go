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
	SnmpproxydConfFile = "snmpproxyd.conf"
	SnmpproxydYamlFile = "snmpproxyd.yaml"
	FibssnmpYamlFile   = "fibssnmp.yaml"
)

type FibsCmd struct {
	*Command

	snmpConfFile       string
	snmpdConfFile      string
	snmpproxydConfFile string
	snmpproxydYamlFile string
	fibssnmpYamlFile   string
}

func NewFibsCmd() *FibsCmd {
	return &FibsCmd{
		Command: NewCommand(),

		snmpConfFile:       SnmpConfFile,
		snmpdConfFile:      SnmpdConfFile,
		snmpproxydConfFile: SnmpproxydConfFile,
		snmpproxydYamlFile: SnmpproxydYamlFile,
		fibssnmpYamlFile:   FibssnmpYamlFile,
	}
}

func (c *FibsCmd) setSnmpFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.snmpConfFile, "snmp-conf-file", "", SnmpConfFile, "snmp.conf file name")
	cmd.Flags().StringVarP(&c.snmpdConfFile, "snmpd-conf-file", "", SnmpdConfFile, "snmpd.conf file name")

	cmd.Flags().StringVarP(&c.snmpproxydConfFile, "snmpproxyd-conf-file", "", SnmpproxydConfFile, "snmpproxyd.conf file name")
	cmd.Flags().StringVarP(&c.snmpproxydYamlFile, "snmpproxyd-yaml-file", "", SnmpproxydYamlFile, "snmpproxyd.yaml file name")
	cmd.Flags().StringVarP(&c.fibssnmpYamlFile, "fibssnmp-yaml-file", "", FibssnmpYamlFile, "fibssnmp.yaml file name")
	return c.Command.setConfigFlags(cmd)
}

func (c *FibsCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	if err := c.createSnmpConf(playbookName); err != nil {
		return err
	}

	if err := c.createSnmpdConf(playbookName); err != nil {
		return err
	}

	if err := c.createSnmpproxydConf(playbookName); err != nil {
		return err
	}

	if err := c.createSnmpproxydYaml(playbookName); err != nil {
		return err
	}

	return c.createFibssnmpConf(playbookName)
}

func (c *FibsCmd) createSnmpproxydConf(playbookName string) error {
	opt := c.optionConfig()

	path := c.filesPath(playbookName, c.snmpproxydConfFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookSnmpproxydConf() // mib/trap
	t.SnmpproxydType = "mib/trap"
	t.SnmpdPort = opt.SnmpdListenPort
	t.SnmpproxydSnmpPort = opt.SnmpproxydSnmpPort
	t.SnmpproxydTrapPort = opt.SnmpproxydTrapPort
	return t.Execute(f)
}

func (c *FibsCmd) createSnmpproxydYaml(playbookName string) error {
	g := c.globalConfig()
	opt := c.optionConfig()

	dpPortMap, err := PortMap(g.DpType)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, c.snmpproxydYamlFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookSnmpproxydYaml()
	t.ONLOidMap = NewSnmpONLOidMap(g.DpAddr)
	t.Trap2Sink = opt.SnmpproxydTrap2Sink
	RangePortMap(dpPortMap, func(pport, lport uint32) {
		t.AddTrap2Map(pport, lport)
	})
	return t.Execute(f)
}

func (c *FibsCmd) createFibssnmpConf(playbookName string) error {
	path := c.filesPath(playbookName, c.fibssnmpYamlFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookFibssnmpYaml()
	return t.Execute(f)
}

func (c *FibsCmd) createSnmpConf(playbookName string) error {
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

func (c *FibsCmd) createSnmpdConf(playbookName string) error {
	opt := c.optionConfig()
	g := c.globalConfig()

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
	t.SnmpPort = opt.SnmpdListenPort
	t.OidMap = snmpOidMap
	t.ONLOidMap = NewSnmpONLOidMap(g.DpAddr)
	return t.Execute(f)
}

func NewFibsCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fibs",
		Short: "fibs command.",
	}

	fibs := NewFibsCmd()
	rootCmd.AddCommand(fibs.setSnmpFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new snmp(d).conf on container file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := fibs.readConfig(); err != nil {
					return err
				}
				return fibs.createConf("common")
			},
		},
	))

	return rootCmd

}
