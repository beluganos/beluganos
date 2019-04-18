// -*- coding: utf-8 -*-

package main

import (
	"fabricflow/util/container/interfacemap"
	"fabricflow/util/netplan"
	"fabricflow/util/sysctl"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	SYSCTL_CONF  = "/etc/sysctl.d/30-beluganos.conf"
	NETPLAN_CONF = "/etc/netplan/20-beluganos.yaml"
	STDOUT_MODE  = "-"
	TEMP_MODE    = "temp"
)

func isStdout(p string) bool {
	return (len(p) == 0) || (p == STDOUT_MODE)
}

func parseVid(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}

	if v > 0xffff {
		return 0, fmt.Errorf("Invalid VlanID. '%s'", s)
	}

	return uint16(v), nil
}

func openOutputFile(p string) (*os.File, error) {
	if isStdout(p) {
		return os.Stdout, nil
	}

	if p == TEMP_MODE {
		dir, fname := path.Split(p)
		return ioutil.TempFile(dir, fname)
	}

	return os.Create(p)
}

func sysctlMplsInputPath(ifname string, vid uint16) string {
	if vid == 0 {
		return fmt.Sprintf("net.mpls.conf.%s.input", ifname)
	}

	return fmt.Sprintf("net.mpls.conf.%s/%d.input", ifname, vid)
}

func sysctlRpFilterPath(ifname string, vid uint16) string {
	if vid == 0 {
		return fmt.Sprintf("net.ipv4.conf.%s.rp_filter", ifname)
	}

	return fmt.Sprintf("net.ipv4.conf.%s/%d.rp_filter", ifname, vid)
}

type VlanCmd struct {
	sysctlPath  string
	sysctlOut   string
	netplanPath string
	netplanOut  string
	dryRun      bool
}

func (c *VlanCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.sysctlPath, "sysctl", "s", SYSCTL_CONF, "sysctl.conf path.")
	cmd.PersistentFlags().StringVarP(&c.sysctlOut, "sysctl-out", "S", SYSCTL_CONF, "sysctl.conf output path.")
	cmd.PersistentFlags().StringVarP(&c.netplanPath, "netplan", "n", NETPLAN_CONF, "netpkan.yaml path.")
	cmd.PersistentFlags().StringVarP(&c.netplanOut, "netplan-out", "N", NETPLAN_CONF, "netpkan.yaml output path.")
	cmd.PersistentFlags().BoolVarP(&c.dryRun, "dry-run", "", false, "dry run mode.")

	return cmd
}

func (c *VlanCmd) SysctlOut() string {
	if c.dryRun {
		return STDOUT_MODE
	}
	return c.sysctlOut
}

func (c *VlanCmd) NetplanOut() string {
	if c.dryRun {
		return STDOUT_MODE
	}
	return c.netplanOut
}

func (c *VlanCmd) readSysctlConf() (*sysctl.SysctlConfig, error) {
	f, err := os.Open(c.sysctlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := sysctl.ReadConfig(f)
	return cfg, nil
}

func (c *VlanCmd) writeSysctlConf(cfg *sysctl.SysctlConfig) error {
	output := c.SysctlOut()

	temp, err := openOutputFile(output)
	if err != nil {
		return err
	}

	defer func() {
		if !isStdout(output) {
			temp.Close()
		}
	}()

	if _, err := cfg.WriteTo(temp); err != nil {
		if !isStdout(output) {
			os.Remove(temp.Name())
		}
		return err
	}

	return nil
}

func (c *VlanCmd) readNetplanConf() (map[interface{}]interface{}, error) {
	f, err := os.Open(c.netplanPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg, err := netplan.ReadConfig(f)
	return cfg, err
}

func (c *VlanCmd) writeNetplanConf(m map[interface{}]interface{}) error {
	output := c.NetplanOut()

	temp, err := openOutputFile(output)
	if err != nil {
		return err
	}

	defer func() {
		if !isStdout(output) {
			temp.Close()
		}
	}()

	if err := netplan.WriteConfig(temp, m); err != nil {
		if !isStdout(output) {
			os.Remove(temp.Name())
		}
		return err
	}

	return nil
}

func (c *VlanCmd) add(ifname string, vlanID string) error {
	vid, err := parseVid(vlanID)
	if err != nil {
		return err
	}

	sysctlCfg, err := c.readSysctlConf()
	if err != nil {
		return err
	}

	netplanCfg, err := c.readNetplanConf()
	if err != nil {
		return err
	}

	sysctlCfg.Set(sysctlMplsInputPath(ifname, vid), "1")
	sysctlCfg.Set(sysctlRpFilterPath(ifname, vid), "0")

	m, ok := interfacemap.SelectOrInsert(netplanCfg, netplan.NewVlanPath(ifname, vid)...)
	if !ok {
		return fmt.Errorf("Netplan config already exist and not map. %s %d", ifname, vid)
	}
	m["link"] = ifname
	m["id"] = vid

	if err := c.writeSysctlConf(sysctlCfg); err != nil {
		return err
	}

	if err := c.writeNetplanConf(netplanCfg); err != nil {
		return err
	}

	return nil
}

func (c *VlanCmd) del(ifname string, vlanID string) error {
	vid, err := parseVid(vlanID)
	if err != nil {
		return err
	}

	sysctlCfg, err := c.readSysctlConf()
	if err != nil {
		return err
	}

	netplanCfg, err := c.readNetplanConf()
	if err != nil {
		return err
	}

	sysctlCfg.Del(sysctlMplsInputPath(ifname, vid))
	sysctlCfg.Del(sysctlRpFilterPath(ifname, vid))

	if ok := interfacemap.Remove(netplanCfg, netplan.NewVlanPath(ifname, vid)...); !ok {
		return fmt.Errorf("Invalid ifname or vlanID. %s %d", ifname, vid)
	}

	if err := c.writeSysctlConf(sysctlCfg); err != nil {
		return err
	}

	if err := c.writeNetplanConf(netplanCfg); err != nil {
		return err
	}

	return nil
}

func vlanCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vlan",
		Short: "VLAN command.",
	}

	addCmd := &VlanCmd{}
	rootCmd.AddCommand(addCmd.setFlags(
		&cobra.Command{
			Use:   "add <ifname> <vid>",
			Short: "Add Vlan device.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return addCmd.add(args[0], args[1])
			},
		},
	))

	delCmd := &VlanCmd{}
	rootCmd.AddCommand(delCmd.setFlags(
		&cobra.Command{
			Use:   "del <ifname> <vid>",
			Short: "Delete Vlan device.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return delCmd.del(args[0], args[1])
			},
		},
	))

	return rootCmd
}
