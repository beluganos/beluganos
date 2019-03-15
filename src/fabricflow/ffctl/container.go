// -*- coding: utf-8 -*-

package main

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

const (
	containerImageName = "base"
)

func doDeviceContainer(name string, excludeIfaces []string, bridge string, cmd string, force bool) error {
	devices, err := containerDevices(name, excludeIfaces)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "doDeviceContainer",
			"call": "containerDevices",
		}).Errorf("%s", err)
		return err
	}

	for _, device := range devices {
		nictype, err := containerDeviceProperty(name, device, "nictype")
		if err != nil {
			log.WithFields(log.Fields{
				"func": "doDeviceContainer",
				"call": "containerDeviceProperty",
			}).Errorf("%s", err)
			return err
		}

		if nictype != "p2p" {
			continue
		}

		ifname, err := containerDeviceProperty(name, device, "host_name")
		if err != nil {
			log.WithFields(log.Fields{
				"func": "doDeviceContainer",
				"call": "containerDeviceProperty",
			}).Errorf("%s", err)
			return err
		}

		if err := execAndOutput("sudo", "ovs-vsctl", cmd, bridge, ifname); err != nil {
			log.WithFields(log.Fields{
				"func": "doDeviceContainer",
				"call": "execAndOutput",
			}).Errorf("%s", err)

			if force {
				continue
			}

			return err
		}

		log.WithFields(log.Fields{
			"container": name,
			"device":    ifname,
		}).Debugf("%s success.", cmd)
	}

	return nil
}

func doAddContainer(name string, excludeIfaces []string, bridge string) error {
	if err := doDeviceContainer(name, excludeIfaces, bridge, "add-port", false); err != nil {
		log.WithFields(log.Fields{
			"container":  name,
			"ovs-bridge": bridge,
		}).Errorf("add container error. %s", err)
		return err
	}

	log.WithFields(log.Fields{
		"container":  name,
		"ovs-bridge": bridge,
	}).Infof("add container ok.")

	return nil
}

func doDeleteContainer(name string, excludeIfaces []string, bridge string) error {
	if err := doDeviceContainer(name, excludeIfaces, bridge, "del-port", true); err != nil {
		log.WithFields(log.Fields{
			"container":  name,
			"ovs-bridge": bridge,
		}).Errorf("delete container error. %s", err)
		return err
	}

	log.WithFields(log.Fields{
		"container":  name,
		"ovs-bridge": bridge,
	}).Infof("delete container ok.")

	return nil
}

func doShowContainer(name string) error {
	return execAndOutput("lxc", "info", name)
}

func doListContainer() error {
	return execAndOutput("lxc", "list")
}

func doStartContainer(name string) error {
	return execAndOutput("lxc", "start", name)
}

func doStopContainer(name string) error {
	return execAndOutput("lxc", "stop", name)
}

func doConContainer(name string) error {
	return execAndWait("lxc", "exec", name, "bash")
}

func containerCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "container",
		Short: "Container commands.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doListContainer()
		},
	}

	var bridge string
	var excludeIfnames []string
	rootCmd.PersistentFlags().StringArrayVarP(&excludeIfnames, "exclude", "", []string{"eth0", "root"}, "Exclude devices.")
	rootCmd.PersistentFlags().StringVarP(&bridge, "bridge", "", ovsBridgeDefault, "ovs-bridge name.")

	addCmd := &cobra.Command{
		Use:   "add <container name>",
		Short: "Add container",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doAddContainer(args[0], excludeIfnames, bridge)
		},
	}
	rootCmd.AddCommand(addCmd)

	delCmd := &cobra.Command{
		Use:     "delete <container name>",
		Aliases: []string{"del"},
		Short:   "Delete container",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doDeleteContainer(args[0], excludeIfnames, bridge)
		},
	}
	rootCmd.AddCommand(delCmd)

	showCmd := &cobra.Command{
		Use:   "show <container name>",
		Short: "Show container",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doShowContainer(args[0])
		},
	}
	rootCmd.AddCommand(showCmd)

	startCmd := &cobra.Command{
		Use:     "start <container name>",
		Aliases: []string{"sta"},
		Short:   "Start container",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doStartContainer(args[0])
		},
	}
	rootCmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop <container name>",
		Short: "Stop container",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doStopContainer(args[0])
		},
	}
	rootCmd.AddCommand(stopCmd)

	conCmd := &cobra.Command{
		Use:     "console <container name>",
		Aliases: []string{"con"},
		Short:   "Run  container console",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doConContainer(args[0])
		},
	}
	rootCmd.AddCommand(conCmd)

	return rootCmd
}
