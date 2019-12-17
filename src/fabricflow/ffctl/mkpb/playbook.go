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

type PlaybookCmd struct {
	*Command
}

func NewPlaybookCmd() *PlaybookCmd {
	return &PlaybookCmd{
		Command: NewCommand(),
	}
}

func (c *PlaybookCmd) createConf(name string) error {
	return c.createPlaybook(name)
}

func (c *PlaybookCmd) createPlaybook(name string) error {
	if len(name) == 0 {
		if routers := c.config.Router; len(routers) == 0 {
			name = "default"
		} else {
			name = routers[0].Name
		}
	}

	path := fmt.Sprintf("%s/lxd-%s.yaml", c.rootPath, name)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	p := NewPlaybook(name)
	return p.Execute(f)
}

func (c *PlaybookCmd) createInventory(hosts ...string) error {
	var name string
	if len(hosts) == 0 {
		if routers := c.config.Router; len(routers) == 0 {
			name = "default"
		} else {
			for _, router := range routers {
				hosts = append(hosts, router.Name)
			}
			name = hosts[0]
		}
	} else {
		name = hosts[0]
	}

	path := fmt.Sprintf("%s/lxd-%s.inv", c.rootPath, name)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})

	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := NewPlaybookInventory()
	t.Name = name
	t.AddHosts(hosts...)
	return t.Execute(f)
}

func NewPlaybookCommand() *cobra.Command {
	pb := NewPlaybookCmd()
	rootCmd := pb.setConfigFlags(
		&cobra.Command{
			Use:   "init [playbook name]",
			Short: "initialize repository.",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var playbookName string
				if len(args) == 0 {
					if err := pb.readConfig(); err != nil {
						return err
					}
				} else {
					playbookName = args[0]
				}
				return pb.createConf(playbookName)
			},
		},
	)

	return rootCmd
}

func NewInventoryCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "inventory",
		Short: "inventory command.",
	}

	pb := NewPlaybookCmd()
	rootCmd.AddCommand(pb.setConfigFlags(
		&cobra.Command{
			Use:   "create <host name ...>",
			Short: "Crate new inventory file.",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := pb.readConfig(); err != nil {
					return err
				}
				return pb.createInventory(args...)
			},
		},
	))

	return rootCmd
}
