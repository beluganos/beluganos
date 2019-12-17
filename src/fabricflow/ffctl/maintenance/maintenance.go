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

package maintenance

import (
	"fabricflow/ffctl/fflib"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lxc/lxd/shared/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "maintenance",
		Short:   "maintenance commands.",
		Aliases: []string{"mn"},
	}

	rootCmd.AddCommand(
		NewMaintenanceAllCommand(),
		NewOpennNSLCommand(),
		NewLxcCommand(),
		NewFrrCommand(),
		NewSysCommand(),
		NewGoBGPCommand(),
		NewRibxCommand(),
		NewFibcCommand(),
	)

	return rootCmd
}

type MaintenanceAllCmd struct {
	outputTopDir string
	fibcAddr     string
	fibcPort     uint16
	gonslPort    uint16

	fibc *FibcCmd

	date time.Time
}

func NewMaintenanceAllCmd() *MaintenanceAllCmd {
	return &MaintenanceAllCmd{
		fibc: NewFibcCmd(),
		date: time.Now(),
	}
}

func (c *MaintenanceAllCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.outputTopDir, "output-dir", "o", ".", "output dir.")
	cmd.Flags().StringVarP(&c.fibcAddr, "fibc-addr", "", fflib.FibcHost, "fibcd addr.")
	cmd.Flags().Uint16VarP(&c.fibcPort, "fibc-port", "", fflib.FibcPort, "fibcd port.")
	cmd.Flags().Uint16VarP(&c.gonslPort, "gonsl-port", "", fflib.GonslPort, "OpenNSL agent port.")
	return cmd
}

func (c *MaintenanceAllCmd) outputDir() string {
	return filepath.Join(c.outputTopDir, c.date.Format("20060102-030405"))
}

func (c *MaintenanceAllCmd) createFile(fileName string) (*os.File, error) {
	path := filepath.Join(c.outputDir(), fileName)
	return fflib.CreateFile(path, true, nil) // overwrite, no call back for backupfile.
}

func (c *MaintenanceAllCmd) mkdir() error {
	return os.MkdirAll(c.outputDir(), 0755)
}

func (c *MaintenanceAllCmd) dumpLxc(name string) error {
	lxc := NewLxcCmd()
	f, err := c.createFile(fmt.Sprintf("%s.log", name))
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("output: %s", f.Name())

	lxc.remoteDump(f, name)
	return nil
}

func (c *MaintenanceAllCmd) dumpOpenNSL() error {
	f, err := c.createFile("gonsld.log")
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("output: %s", "gonsld.log")

	cmd := NewOpenNSLCmd()
	cmd.gonsl.Port = c.gonslPort
	cmd.fibc.Host = c.fibcAddr
	cmd.fibc.Port = c.fibcPort
	return cmd.dumpAll(f)
}

func (c *MaintenanceAllCmd) dumpFibc() error {
	f, err := c.createFile("fibcd.log")
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("output: %s", "fibcd.log")

	cmd := NewFibcCmd()
	cmd.fibc.Host = c.fibcAddr
	cmd.fibc.Port = c.fibcPort
	return cmd.dump(f)
}

func (c *MaintenanceAllCmd) dump() error {
	containers, err := fflib.LXDContainerList()
	if err != nil {
		return err
	}

	if err := c.mkdir(); err != nil {
		return err
	}

	log.Debugf("output dir: %s", c.outputDir())

	for _, container := range containers {
		if container.StatusCode == api.Running {
			log.Debugf("container: %s", container.Name)
			if err := c.dumpLxc(container.Name); err != nil {
				return err
			}
		}
	}

	return c.dumpOpenNSL()
}

func NewMaintenanceAllCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "all",
		Short: "maintenance all command.",
	}

	c := NewMaintenanceAllCmd()

	rootCmd.AddCommand(c.setFlags(
		&cobra.Command{
			Use:   "dump",
			Short: "maintenance all dump command.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.dump()
			},
		},
	))

	return rootCmd
}
