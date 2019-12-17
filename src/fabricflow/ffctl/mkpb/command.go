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
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	IfnamePrefix    = "eth"
	CommandRootPath = "."
	CommandRole     = "lxd"
)

type Command struct {
	rootPath  string
	role      string
	overwrite bool

	configFile string
	configType string
	lxdMode    bool
	config     *Config
}

func NewCommand() *Command {
	return &Command{
		rootPath: CommandRootPath,
		role:     CommandRole,

		config: NewConfig(),
	}
}

func (c *Command) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.rootPath, "path", "", CommandRootPath, "Path to playbooks.")
	cmd.Flags().StringVarP(&c.role, "role", "", CommandRole, "Role name.")
	cmd.Flags().BoolVarP(&c.overwrite, "overwrite", "", false, "overwrite")
	return cmd
}

func (c *Command) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.configFile, "config-file", "", "", "config file path.")
	cmd.Flags().StringVarP(&c.configType, "config-type", "", "", "config file type.")
	cmd.Flags().BoolVarP(&c.lxdMode, "lxd-config-mode", "", false, "create for lxd.")
	return c.setFlags(cmd)
}

func (c *Command) setConfig(config *Config) {
	c.config = config
}

func (c *Command) readConfig() error {
	setLXDConfigMode(c.lxdMode)
	c.config.SetConfig(c.configFile, c.configType)
	return c.config.Load()
}

func (c *Command) globalConfig() *GlobalConfig {
	return c.config.Global
}

func (c *Command) optionConfig() *OptionConfig {
	return &c.config.Option
}

func (c *Command) routerNameList() []string {
	names := []string{}
	for _, router := range c.config.Router {
		names = append(names, router.Name)
	}

	return names
}

func (c *Command) routerConfig(name string) (*RouterConfig, error) {
	router := c.config.GetRouter(name)
	if router == nil {
		return nil, fmt.Errorf("router not found. %s", name)
	}
	return router, nil
}

func (c *Command) micConfig() (*RouterConfig, error) {
	mic := c.config.GetMICRouter()
	if mic == nil {
		return nil, fmt.Errorf("mic router not found.")
	}
	return mic, nil
}

func (c *Command) filesDirPath(name string) string {
	if c.config.Option.LXDConfigMode {
		return filepath.Join(c.rootPath, name)
	}
	return filepath.Join(c.rootPath, "roles", c.role, "files", name)
}

func (c *Command) filesPath(name string, filename string) string {
	return filepath.Join(c.filesDirPath(name), filename)
}

func (c *Command) mkDirAll(name string) error {
	path := c.filesDirPath(name)
	log.Debugf("%s created.", path)
	return os.MkdirAll(path, 0755)
}
