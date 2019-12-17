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
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	LXDProfileFileName = "lxd_profile.yml"
)

type LXDCmd struct {
	*Command

	fileName  string
	outputDir string
}

func NewLXDCmd() *LXDCmd {
	return &LXDCmd{
		Command: NewCommand(),

		fileName: LXDProfileFileName,
	}
}

func (c *LXDCmd) setConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "profile-file", "", LXDProfileFileName, "lxd profile file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c *LXDCmd) setConvFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "profile-file", "", LXDProfileFileName, "lxd profile file name.")
	cmd.Flags().StringVarP(&c.outputDir, "output", "o", ".", "Path to output")
	return c.Command.setFlags(cmd)
}

func (c *LXDCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createProfile(playbookName)
}

func (c *LXDCmd) createProfile(playbookName string) error {
	opt := c.optionConfig()
	g := c.globalConfig()
	r, err := c.routerConfig(playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, c.fileName)
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

	t := NewPlaybookLXDProfile()
	t.Name = playbookName
	t.MngIface = opt.LXDMngInterface
	t.BridgeIface = opt.LXDBridge
	t.Mtu = opt.LXDMtu
	t.AddPorts(ConvToPPortList(r.Eth, portmap))
	return t.Execute(f, c.config.Option.LXDConfigMode)
}

func (c *LXDCmd) convProfile(playbookName string) error {
	path := c.filesPath(playbookName, c.fileName)
	log.Debugf("convert '%s'", path)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	funcs := template.FuncMap{
		"lxcname": func() string { return playbookName },
	}

	t, err := template.New("lxd_profile").Funcs(funcs).Parse(string(b))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, nil); err != nil {
		return err
	}

	data := []map[string]interface{}{}
	if err := yaml.Unmarshal(buf.Bytes(), &data); err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("file is empty,")
	}
	profileMap, ok := data[0]["lxd_profile"]
	if !ok {
		return fmt.Errorf("'lxd_profile' not exist.'")
	}

	if m, ok := profileMap.(map[interface{}]interface{}); ok {
		delete(m, "state")
	}

	profile, err := yaml.Marshal(profileMap)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(c.outputDir, c.fileName))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(profile); err != nil {
		return err
	}
	f.Sync()

	return nil
}

func NewLXDCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "lxd",
		Short: "lxd command.",
	}

	lxd := NewLXDCmd()
	rootCmd.AddCommand(lxd.setConfigFlags(
		&cobra.Command{
			Use:   "create [playbook name]",
			Short: "Crate new lxd_profile.yml file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := lxd.readConfig(); err != nil {
					return err
				}
				return lxd.createConf(args[0])
			},
		},
	))

	rootCmd.AddCommand(lxd.setConvFlags(
		&cobra.Command{
			Use:     "convert-to-lxd [playbook name]",
			Short:   "Convert lxd_profile.yaml to lxd profile.",
			Aliases: []string{"conv", "convert"},
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return lxd.convProfile(args[0])
			},
		},
	))

	return rootCmd
}
