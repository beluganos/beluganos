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
	"fabricflow/ffctl/dpport"
	"fabricflow/ffctl/templates"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var playbookFrrDaemons = map[string]string{
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

func newPlaybookFrrDaemons(daemons []string) map[string]string {
	m := map[string]string{}
	for name, _ := range playbookFrrDaemons {
		m[name] = "no"
	}
	for _, daemon := range daemons {
		m[daemon] = "yes"
	}
	return m
}

func getDpPortMapAndCfg(dpType, config, playbookName string) (map[uint]uint, *dpport.PortConfig, error) {
	dpPortMap, err := dpport.PortMap(dpType)
	if err != nil {
		return nil, nil, err
	}

	dpPortCfg, _ := dpport.ReadConfig(config).PortConfig(playbookName)
	if n := len(dpPortCfg.Eth); n == 0 {
		for pport, _ := range dpPortMap {
			dpPortCfg.Eth = append(dpPortCfg.Eth, pport)
		}
	}

	dpPortMap = dpPortCfg.Filter(dpPortMap)

	return dpPortMap, dpPortCfg, nil
}

type PlaybookCmd struct {
	rootPath  string
	role      string
	overwrite bool
}

func NewPlaybookCmd() *PlaybookCmd {
	return &PlaybookCmd{}
}

func (c *PlaybookCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.rootPath, "path", "", ".", "Path to playbooks.")
	cmd.PersistentFlags().StringVarP(&c.role, "role", "", "lxd", "Role name.")
	cmd.PersistentFlags().BoolVarP(&c.overwrite, "overwrite", "", false, "overwrite")
	return cmd
}

func (c *PlaybookCmd) filesDirPath(name string) string {
	return fmt.Sprintf("%s/roles/%s/files/%s", c.rootPath, c.role, name)
}

func (c *PlaybookCmd) filesPath(name string, filename string) string {
	return fmt.Sprintf("%s/roles/%s/files/%s/%s", c.rootPath, c.role, name, filename)
}

func (c *PlaybookCmd) mkDirAll(name string) error {
	path := c.filesDirPath(name)
	log.Debugf("%s created.", path)
	return os.MkdirAll(path, 0755)
}

func (c *PlaybookCmd) create(name string) error {
	path := fmt.Sprintf("%s/lxd-%s.yaml", c.rootPath, name)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	p := templates.NewPlaybook(name)
	if err := p.Execute(f); err != nil {
		return err
	}

	if err := c.mkDirAll(name); err != nil {
		return err
	}

	return nil
}

func (c *PlaybookCmd) createInventory(hosts ...string) error {
	name := hosts[0]
	path := fmt.Sprintf("%s/lxd-%s.inv", c.rootPath, name)
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})

	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookInventory()
	t.Name = name
	t.AddHosts(hosts...)
	return t.Execute(f)
}

func (c *PlaybookCmd) createDaemons(playbookName string, daemons []string) error {
	path := c.filesPath(playbookName, "daemons")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})

	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookDaemons()
	t.SetMap(newPlaybookFrrDaemons(daemons))
	return t.Execute(f)
}

func (c *PlaybookCmd) createGoBGPConf(playbookName string) error {
	path := c.filesPath(playbookName, "gobgp.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookGoBGPConf()
	return t.Execute(f)
}

func (c *PlaybookCmd) createRibtdConf(playbookName string) error {
	path := c.filesPath(playbookName, "ribtd.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookRibtdConf()
	return t.Execute(f)
}

func (c *PlaybookCmd) createSnmpProxydConf(playbookName string) error {
	path := c.filesPath(playbookName, "snmpproxyd.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookSnmpProxydConf(false)
	return t.Execute(f)
}

func (c *PlaybookCmd) createSnmpProxydYaml(playbookName string) error {
	path := c.filesPath(playbookName, "snmpproxyd.yaml")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookSnmpProxydYaml()
	return t.Execute(f)
}

type PlaybookCommdnCmd struct {
	*PlaybookCmd
	dpType string
}

func NewPlaybookCommonCmd() *PlaybookCommdnCmd {
	return &PlaybookCommdnCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookCommdnCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().StringVarP(&c.dpType, "dp-type", "", "as5812", "datapath type.")

	return cmd
}

func (c *PlaybookCommdnCmd) createSnmpProxyConf() error {
	path := c.filesPath("common", "snmpproxyd.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookSnmpProxydConf(true)
	return t.Execute(f)
}

func (c *PlaybookCommdnCmd) createSnmpProxyYaml() error {
	dpPortMap, err := dpport.PortMap(c.dpType)
	if err != nil {
		return err
	}

	path := c.filesPath("common", "snmpproxyd.yaml")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookSnmpProxydYaml()
	for pport, lport := range dpPortMap {
		t.AddTrap2Map(pport, lport)
	}
	return t.Execute(f)
}

func (c *PlaybookCommdnCmd) createCommon() error {

	if err := c.mkDirAll("common"); err != nil {
		return err
	}

	if err := c.createSnmpProxyConf(); err != nil {
		return err
	}

	if err := c.createSnmpProxyYaml(); err != nil {
		return err
	}

	return nil
}

type PlaybookFibcCmd struct {
	*PlaybookCmd
	reID   string
	dpID   uint64
	dpName string
	dpMode string
	dpType string
	config string
}

func NewPlaybookFibcCmd() *PlaybookFibcCmd {
	return &PlaybookFibcCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookFibcCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().StringVarP(&c.reID, "re-id", "", "0.0.0.0", "router entity id.")
	cmd.PersistentFlags().Uint64VarP(&c.dpID, "dp-id", "", 0, "datapath id. (default: random)")
	cmd.PersistentFlags().StringVarP(&c.dpName, "dp-name", "", "", "datapath name. (default: random)")
	cmd.PersistentFlags().StringVarP(&c.dpMode, "dp-mode", "", "onsl", "datapath mode.")
	cmd.PersistentFlags().StringVarP(&c.dpType, "dp-type", "", "as5812", "datapath type.")
	cmd.PersistentFlags().StringVarP(&c.config, "ports", "", "./ports.yaml", "Port filepath.")

	return cmd
}

func (c *PlaybookFibcCmd) createFibcYaml(playbookName string) error {
	dpPortMap, _, err := getDpPortMapAndCfg(c.dpType, c.config, playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, "fibc.yml")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	dpID := func() uint64 {
		if c.dpID != 0 {
			return c.dpID
		}
		r := rand.New(rand.NewSource(time.Now().Unix()))
		return r.Uint64()
	}()

	dpName := func() string {
		if len(c.dpName) != 0 {
			return c.dpName
		}
		return fmt.Sprintf("dp_%d", dpID)
	}()

	t := templates.NewPlaybookFibcYaml(c.reID, c.dpName)
	t.ReID = c.reID
	t.DpID = dpID
	t.DpName = dpName
	t.DpMode = c.dpMode
	t.AddPorts(dpPortMap)
	return t.Execute(f)
}

type PlaybookFrrConfCmd struct {
	*PlaybookCmd
	routerID string
	dpType   string
	config   string
}

func NewPlaybookFrrConfCmd() *PlaybookFrrConfCmd {
	return &PlaybookFrrConfCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookFrrConfCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().StringVarP(&c.routerID, "router-id", "", "0.0.0.0", "router id.")
	cmd.PersistentFlags().StringVarP(&c.dpType, "dp-type", "", "as5812", "datapath type.")
	cmd.PersistentFlags().StringVarP(&c.config, "ports", "", "./ports.yaml", "Port filepath.")

	return cmd
}

func (c *PlaybookFrrConfCmd) createFrrConf(playbookName string) error {
	_, dpPortCfg, err := getDpPortMapAndCfg(c.dpType, c.config, playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, "frr.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookFrrConf(c.routerID)
	for _, dev := range dpPortCfg.Devices() {
		t.AddIface(dev.Eth, dev.Vid)
	}
	return t.Execute(f)
}

type PlaybookGoBGPdConfCmd struct {
	*PlaybookCmd
	routerID string
	AS       uint
}

func NewPlaybookGoBGPdConfCmd() *PlaybookGoBGPdConfCmd {
	return &PlaybookGoBGPdConfCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookGoBGPdConfCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().StringVarP(&c.routerID, "router-id", "", "0.0.0.0", "router id.")
	cmd.PersistentFlags().UintVarP(&c.AS, "as", "", 65001, "AS number.")
	return cmd
}

func (c *PlaybookGoBGPdConfCmd) createGoBGPdConf(playbookName string) error {
	path := c.filesPath(playbookName, "gobgpd.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookGoBGPdConf()
	t.RouterID = c.routerID
	t.AS = c.AS
	return t.Execute(f)
}

type PlaybookLXDProfileCmd struct {
	*PlaybookCmd
	mngIface    string
	bridgeIface string
	mtu         uint16
	dpType      string
	config      string
}

func NewPlaybookLXDProfileCmd() *PlaybookLXDProfileCmd {
	return &PlaybookLXDProfileCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookLXDProfileCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().StringVarP(&c.mngIface, "mng-iface", "", "eth0", "managment interface.")
	cmd.PersistentFlags().StringVarP(&c.bridgeIface, "bridge-iface", "", "lxdbr0", "lxd bridge interface.")
	cmd.PersistentFlags().Uint16VarP(&c.mtu, "mtu", "", 9000, "mtu")
	cmd.PersistentFlags().StringVarP(&c.dpType, "dp-type", "", "as5812", "datapath type.")
	cmd.PersistentFlags().StringVarP(&c.config, "ports", "", "./ports.yaml", "Port filepath.")

	return cmd
}

func (c *PlaybookLXDProfileCmd) createLXDProfile(playbookName string) error {
	_, dpPortCfg, err := getDpPortMapAndCfg(c.dpType, c.config, playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, "lxd_profile.yml")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookLXDProfile()
	t.Name = playbookName
	t.MngIface = c.mngIface
	t.BridgeIface = c.bridgeIface
	t.Mtu = c.mtu
	t.AddPorts(dpPortCfg.Eth...)
	return t.Execute(f)
}

type PlaybookNetplanYamlCmd struct {
	*PlaybookCmd
	mtu    uint16
	dpType string
	config string
}

func NewPlaybookNetplanYamlCmd() *PlaybookNetplanYamlCmd {
	return &PlaybookNetplanYamlCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookNetplanYamlCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().Uint16VarP(&c.mtu, "mtu", "", 9000, "mtu")
	cmd.PersistentFlags().StringVarP(&c.dpType, "dp-type", "", "as5812", "datapath type.")
	cmd.PersistentFlags().StringVarP(&c.config, "ports", "", "./ports.yaml", "Port filepath.")

	return cmd
}

func (c *PlaybookNetplanYamlCmd) createNetplanYaml(playbookName string) error {
	dpPortMap, dpPortCfg, err := getDpPortMapAndCfg(c.dpType, c.config, playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, "netplan.yml")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookNetplanYaml()
	for pport, _ := range dpPortMap {
		t.AddEth(pport, c.mtu)
	}
	for _, vlan := range dpPortCfg.Vlans() {
		t.AddVlan(vlan.Eth, vlan.Vid)
	}

	return t.Execute(f)
}

type PlaybookSysctlConfCmd struct {
	*PlaybookCmd
	mplsLabel   uint
	sockBufSize uint
	dpType      string
	config      string
}

func NewPlaybookSysctlConfCmd() *PlaybookSysctlConfCmd {
	return &PlaybookSysctlConfCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookSysctlConfCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().UintVarP(&c.mplsLabel, "mpls-label", "", 10240, "mpls label max size")
	cmd.PersistentFlags().UintVarP(&c.sockBufSize, "sockbuf-size", "", 8388608, "socket buffer size.")
	cmd.PersistentFlags().StringVarP(&c.dpType, "dp-type", "", "as5812", "datapath type.")
	cmd.PersistentFlags().StringVarP(&c.config, "ports", "", "./ports.yaml", "Port filepath.")

	return cmd
}

func (c *PlaybookSysctlConfCmd) createSysctlConf(playbookName string) error {
	dpPortMap, dpPortCfg, err := getDpPortMapAndCfg(c.dpType, c.config, playbookName)
	if err != nil {
		return err
	}

	path := c.filesPath(playbookName, "sysctl.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookSysctlConf()
	t.SockBufSize = c.sockBufSize
	t.MplsLabel = c.mplsLabel
	for pport, _ := range dpPortMap {
		t.AddIface(pport, 0)
	}
	for _, vlan := range dpPortCfg.Vlans() {
		t.AddIface(vlan.Eth, vlan.Vid)
	}

	return t.Execute(f)
}

type PlaybookRibxdConfCmd struct {
	*PlaybookCmd
	nid  uint8
	reID string
	rt   string
	rd   string
	vpn  bool
}

func NewPlaybookRibxdConfCmd() *PlaybookRibxdConfCmd {
	return &PlaybookRibxdConfCmd{
		PlaybookCmd: NewPlaybookCmd(),
	}
}

func (c *PlaybookRibxdConfCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	c.PlaybookCmd.setFlags(cmd)
	cmd.PersistentFlags().Uint8VarP(&c.nid, "node-id", "", 0, "node id.")
	cmd.PersistentFlags().StringVarP(&c.reID, "re-id", "", "0.0.0.0", "router entity id.")
	cmd.PersistentFlags().StringVarP(&c.rt, "rt", "", "", "route target.")
	cmd.PersistentFlags().StringVarP(&c.rd, "rd", "", "", "route dist.")
	cmd.PersistentFlags().BoolVarP(&c.vpn, "vpn", "", false, "vpn or not.")
	return cmd
}

func (c *PlaybookRibxdConfCmd) createRibxdConf(playbookName string) error {
	path := c.filesPath(playbookName, "ribxd.conf")
	f, err := createFile(path, c.overwrite, func(backup string) {
		log.Debugf("%s backup", backup)
	})
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("%s created.", path)

	t := templates.NewPlaybookRibxdConf()
	t.NID = c.nid
	t.ReID = c.reID
	t.RT = c.rt
	t.RD = c.rd
	t.Vpn = c.vpn
	t.Name = playbookName

	return t.Execute(f)
}
