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
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

const (
	containerImageName = "base"
)

type containerCommand struct {
	excludeIfnames []string
	bridge         string

	log *log.Entry
}

func (c *containerCommand) setFlags(cmd *cobra.Command) *cobra.Command {
	c.log = log.WithFields(log.Fields{
		"module": "container",
	})
	return cmd
}

func (c *containerCommand) setIfaceFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringArrayVarP(&c.excludeIfnames, "exclude", "", []string{"eth0", "root", "logdir"}, "Exclude devices.")
	cmd.Flags().StringVarP(&c.bridge, "bridge", "", ovsBridgeDefault, "ovs-bridge name.")
	return c.setFlags(cmd)
}

func (c *containerCommand) modOvsPort(name, cmd string, force bool) error {
	c.log.Debugf("modOvsPort: name:%s cmd:%s", name, cmd)

	devices, err := containerDevices(name, c.excludeIfnames)
	if err != nil {
		c.log.Errorf("modOvsPort: name:%s cmd:%s %s", name, cmd, err)
		return err
	}

	for _, device := range devices {
		c.log.Debugf("modOvsPort: name:%s cmd:%s device:%s", name, cmd, device)

		nictype, err := containerDeviceProperty(name, device, "nictype")
		if err != nil {
			c.log.Errorf("modOvsPort: name:%s cmd:%s device:%s %s", name, cmd, device, err)
			return err
		}

		if nictype != "p2p" {
			continue
		}

		ifname, err := containerDeviceProperty(name, device, "host_name")
		if err != nil {
			c.log.Errorf("modOvsPort: name:%s cmd:%s device:%s %s", name, cmd, device, err)
			return err
		}

		if err := execAndOutput("sudo", "ovs-vsctl", cmd, c.bridge, ifname); err != nil {
			if force {
				c.log.Warnf("modOvsPort: name:%s cmd:%s device:%s %s", name, cmd, device, err)
				continue
			}

			c.log.Errorf("modOvsPort: name:%s cmd:%s device:%s %s", name, cmd, device, err)
			return err
		}
	}

	return nil
}

func (c *containerCommand) addOvsPort(name string) error {
	c.log.Debugf("addOvsPort: name:%s", name)

	if err := c.modOvsPort(name, "add-port", false); err != nil {
		c.log.Errorf("addOvsPort: %s", err)
		return err
	}

	return nil
}

func (c *containerCommand) deleteOvsPort(name string) error {
	c.log.Debugf("deleteOvsPort: name:%s", name)

	if err := c.modOvsPort(name, "del-port", true); err != nil {
		c.log.Errorf("deleteOvsPort: %s", err)
		return err
	}

	return nil
}

func (c *containerCommand) showInfo(name string) error {
	return execAndOutput("lxc", "info", name)
}

func (c *containerCommand) showList() error {
	return execAndOutput("lxc", "list")
}

func (c *containerCommand) start(name string) error {
	return execAndOutput("lxc", "start", name)
}

func (c *containerCommand) stop(name string) error {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.showList()
		},
	}

	rootCmd.AddCommand(c.setIfaceFlags(
		&cobra.Command{
			Use:   "add <container name>",
			Short: "Add container",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.addOvsPort(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setIfaceFlags(
		&cobra.Command{
			Use:     "delete <container name>",
			Aliases: []string{"del"},
			Short:   "Delete container",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.deleteOvsPort(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setFlags(
		&cobra.Command{
			Use:   "show <container name>",
			Short: "Show container",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.showInfo(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setFlags(
		&cobra.Command{
			Use:     "start <container name>",
			Aliases: []string{"sta"},
			Short:   "Start container",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.start(args[0])
			},
		},
	))

	rootCmd.AddCommand(c.setFlags(
		&cobra.Command{
			Use:   "stop <container name>",
			Short: "Stop container",
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
			Short:   "Run container console",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.console(args[0])
			},
		},
	))

	return rootCmd
}
