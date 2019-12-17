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

var frrDaemons = map[string]string{
	"zebra":  "no",
	"bgpd":   "no",
	"ospfd":  "no",
	"ospf6d": "no",
	"ripd":   "no",
	"ripngd": "no",
	"isisd":  "no",
	"pimd":   "no",
	"ldpd":   "no",
	"nhrpd":  "no",
}

func newFrrDaemons(daemons []string) map[string]string {
	m := map[string]string{}
	for name, _ := range frrDaemons {
		m[name] = "no"
	}
	for _, daemon := range daemons {
		m[daemon] = "yes"
	}
	return m
}

const (
	FrrConfFile    = "frr.conf"
	FrrDaemonsFile = "daemons"
)

type FrrCmd struct {
	*Command
}

func NewFrrCmd() *FrrCmd {
	return &FrrCmd{
		Command: NewCommand(),
	}
}

func (c *FrrCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	if err := c.createDaemons(playbookName); err != nil {
		return err
	}

	return c.createFrrConf(playbookName)
}

func (c *FrrCmd) createDaemons(playbookName string) error {
	r, err := c.routerConfig(playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, FrrDaemonsFile)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})

	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookDaemons()
	t.SetMap(newFrrDaemons(r.Daemons))
	return t.Execute(f)
}

func (c *FrrCmd) createFrrConf(playbookName string) error {
	g := c.globalConfig()
	r, err := c.routerConfig(playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, FrrConfFile)
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

	t := NewPlaybookFrrConf(g.ReID)
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

func NewFrrCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "frr",
		Short: "frr config file command.",
	}

	frr := NewFrrCmd()
	rootCmd.AddCommand(frr.setConfigFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Crate new frr.conf file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := frr.readConfig(); err != nil {
					return err
				}
				return frr.createConf(args[0])
			},
		},
	))

	return rootCmd
}
