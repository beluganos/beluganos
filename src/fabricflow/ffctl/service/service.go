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

package service

import (
	"fabricflow/ffctl/fflib"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

const (
	ServiceFibcName  = "fibcd-go"
	ServiceFibsName  = "fibsd"
	ServiceGovswName = "govswd"
)

type ServiceCmd struct {
	FibcName   string
	FibsName   string
	GovswName  string
	AllService bool
}

func (c *ServiceCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.FibcName, "fibcd-service-name", "", ServiceFibcName, "fibcd service name.")
	cmd.Flags().StringVarP(&c.FibsName, "fibsd-service-name", "", ServiceFibsName, "fibsd service name.")
	cmd.Flags().StringVarP(&c.GovswName, "govswd-service-name", "", ServiceGovswName, "govswd service name.")
	cmd.Flags().BoolVarP(&c.AllService, "all", "", false, "execute all service.")
	return cmd
}

func (c *ServiceCmd) serviceNames(reverse bool) []string {
	names := []string{
		c.FibcName,
		c.FibsName,
		c.GovswName,
	}
	if reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(names)))
	}
	return names
}

func (c *ServiceCmd) status() error {
	for _, name := range c.serviceNames(false) {
		fflib.ExecAndOutput(os.Stdout, "systemctl", "status", "--no-pager", name)
	}

	fflib.ExecAndOutput(os.Stdout, "lxc", "list")

	return nil
}

func (c *ServiceCmd) start() error {
	for _, name := range c.serviceNames(false) {
		if err := fflib.ExecAndOutput(os.Stdout, "sudo", "systemctl", "start", name); err != nil {
			return err
		}
	}
	return nil
}

func (c *ServiceCmd) stop() error {
	for _, name := range c.serviceNames(true) {
		if err := fflib.ExecAndOutput(os.Stdout, "sudo", "systemctl", "stop", name); err != nil {
			return err
		}
	}
	return nil
}

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "service",
		Short: "service control command.",
	}

	svc := ServiceCmd{}

	rootCmd.AddCommand(svc.setFlags(
		&cobra.Command{
			Use:   "start <service name | --all>",
			Short: "start service",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if svc.AllService {
					return svc.start()
				}

				if len(args) == 0 {
					return fmt.Errorf("service not specified.")
				}

				return fflib.ExecAndOutput(os.Stdout, "sudo", "systemctl", "start", args[0])
			},
		},
	))

	rootCmd.AddCommand(svc.setFlags(
		&cobra.Command{
			Use:   "stop <service name | --all>",
			Short: "stop service",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if svc.AllService {
					return svc.stop()
				}

				if len(args) == 0 {
					return fmt.Errorf("service not specified.")
				}

				return fflib.ExecAndOutput(os.Stdout, "sudo", "systemctl", "stop", args[0])
			},
		},
	))

	rootCmd.AddCommand(svc.setFlags(
		&cobra.Command{
			Use:   "status <service name | --all>",
			Short: "show service status.",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if svc.AllService {
					svc.status()
					return nil
				}

				if len(args) == 0 {
					return fmt.Errorf("service not specified.")
				}

				fflib.ExecAndOutput(os.Stdout, "systemctl", "status", args[0])
				return nil
			},
		},
	))

	rootCmd.AddCommand(svc.setFlags(
		&cobra.Command{
			Use:   "enable <service name>",
			Short: "enable service",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return fflib.ExecAndOutput(os.Stdout, "sudo", "systemctl", "enable", args[0])
			},
		},
	))

	rootCmd.AddCommand(svc.setFlags(
		&cobra.Command{
			Use:   "disable <service name>",
			Short: "disable service",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return fflib.ExecAndOutput(os.Stdout, "sudo", "systemctl", "disable", args[0])
			},
		},
	))

	return rootCmd
}
