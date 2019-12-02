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

package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	containerImageName = "base"
	containerGovswCmd  = "govswc"
)

type containerCommand struct {
	excludeDevices []string
	vswCmd         string
	withIfaces     bool
}

func (c *containerCommand) setFlags(cmd *cobra.Command) *cobra.Command {
	return cmd
}

func (c *containerCommand) setContainerFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolVarP(&c.withIfaces, "with-interfaces", "I", false, "exec (un)register interfaces.")
	return c.setIfaceFlags(cmd)
}

func (c *containerCommand) setIfaceFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringArrayVarP(&c.excludeDevices, "exclude", "", []string{"eth0", "root", "logdir"}, "Exclude devices.")
	cmd.Flags().StringVarP(&c.vswCmd, "vsw-command", "", containerGovswCmd, "vsw command.")
	return c.setFlags(cmd)
}

func (c *containerCommand) modPort(name, cmd string, force bool) error {
	log.Debugf("modPort: name:%s cmd:%s", name, cmd)

	ifnames, err := containerHostIfnames(name, c.excludeDevices)
	if err != nil {
		return err
	}

	for _, ifname := range ifnames {
		if err := execAndOutput(c.vswCmd, "interface", cmd, ifname); err != nil {
			if force {
				log.Warnf("%s error. %s %s %s", c.vswCmd, cmd, ifname, err)
				continue
			}

			log.Errorf("%s error. %s %s %s", c.vswCmd, cmd, ifname, err)
			return err
		}

		log.Debugf("%s %s %s success.", c.vswCmd, cmd, ifname)
	}

	return nil
}

func (c *containerCommand) registerPort(name string) error {
	log.Debugf("addPort: name:%s", name)

	if err := c.modPort(name, "add", false); err != nil {
		log.Errorf("addPort: %s", err)
		return err
	}

	return nil
}

func (c *containerCommand) unregisterPort(name string) error {
	log.Debugf("deletePort: name:%s", name)

	if err := c.modPort(name, "delete", true); err != nil {
		log.Errorf("deletePort: %s", err)
		return err
	}

	return nil
}

func (c *containerCommand) showInfo(name string) {
	execAndOutput("lxc", "info", name)
	execAndOutput("lxc", "config", "show", name)
}

func (c *containerCommand) showList() error {
	return execAndOutput("lxc", "list")
}

func (c *containerCommand) start(name string) error {
	if err := execAndOutput("lxc", "start", name); err != nil {
		return err
	}

	if c.withIfaces {
		c.registerPort(name)
	}

	return nil
}

func (c *containerCommand) stop(name string) error {
	if c.withIfaces {
		c.unregisterPort(name)
	}

	return execAndOutput("lxc", "stop", name)
}

func (c *containerCommand) console(name string) error {
	return execAndWait("lxc", "exec", name, "bash")
}

func containerCmd() *cobra.Command {
	c := containerCommand{}

	rootCmd := &cobra.Command{
		Use:   "container",
		Short: "Container commands.",
	}

	rootCmd.AddCommand(c.setIfaceFlags(
		&cobra.Command{
			Use:     "register <container name>",
			Short:   "register container interfaces.",
			Aliases: []string{"r", "reg"},
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.registerPort(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setIfaceFlags(
		&cobra.Command{
			Use:     "unregister <container name>",
			Aliases: []string{"u", "unreg"},
			Short:   "unregister container interfaces.",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.unregisterPort(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setFlags(
		&cobra.Command{
			Use:   "status <container name>",
			Short: "status container status",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				c.showInfo(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setFlags(
		&cobra.Command{
			Use:   "list",
			Short: "show container list",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.showList()
			},
		},
	))

	rootCmd.AddCommand(c.setContainerFlags(
		&cobra.Command{
			Use:     "start <container name>",
			Aliases: []string{"sta"},
			Short:   "start container",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.start(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setContainerFlags(
		&cobra.Command{
			Use:   "stop <container name>",
			Short: "stop container",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.stop(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setFlags(
		&cobra.Command{
			Use:     "console <container name>",
			Aliases: []string{"con"},
			Short:   "run container console",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.console(args[0])
			},
		},
	))

	return rootCmd
}
