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

const (
	NFTableConfFileName = "nftables.conf"
)

func newNATPreRoutingRuleSSH(ifname, newDst string) string {
	return fmt.Sprintf("iif \"%s\" tcp dport $ic_ssh dnat to %s:ssh", ifname, newDst)
}

func newNATPreRoutingRuleSNMP(ifname, newDst string) string {
	return fmt.Sprintf("iif \"%s\" udp dport $ic_snmp dnat to %s:snmp", ifname, newDst)
}

func newNATPreRoutingRuleSNMPTrap(ifname, newDst string) string {
	return fmt.Sprintf("iif \"%s\" udp dport $ic_snmptrap dnat to %s:snmp-trap", ifname, newDst)
}

func newNATPostRoutingRuleOIFMasq(ifname string) string {
	return fmt.Sprintf("oif \"%s\" masquerade", ifname)
}

func newNATPostRoutingRuleSnmpTrapMasq() string {
	return "udp dport snmp-trap masquerade"
}

type NFTableCmd struct {
	*Command

	fileName string
}

func NewNFTableCmd() *NFTableCmd {
	return &NFTableCmd{
		Command: NewCommand(),

		fileName: NFTableConfFileName,
	}
}

func (c NFTableCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fileName, "file-name", "", NFTableConfFileName, "output file name.")
	return c.Command.setConfigFlags(cmd)
}

func (c NFTableCmd) createConf(playbookName string) error {
	if err := c.mkDirAll(playbookName); err != nil {
		return err
	}

	return c.createNFTableConf(playbookName)
}

func (c NFTableCmd) createNFTableConf(playbookName string) error {
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

	t := NewPlaybookNFTableConf()
	t.PortSSH = opt.InChOpeSSHPort
	t.PortSNMP = opt.InChOpeSNMPPort
	t.PortSNMPTrap = opt.InChOpeSNMPTrapPort
	for _, eth := range ConvToPPortList(r.Eth, portmap) {
		ifname := NewPhyIfname(eth)
		t.AddNATPreRoutingRule(newNATPreRoutingRuleSSH(ifname, opt.SnmpproxydAddr()))
		t.AddNATPreRoutingRule(newNATPreRoutingRuleSNMP(ifname, opt.SnmpproxydAddr()))
	}
	t.AddNATPreRoutingRule(newNATPreRoutingRuleSNMPTrap(opt.LXDMngInterface, opt.InChOpeSNMPTrapSink))
	t.AddNATPostRoutingRule(newNATPostRoutingRuleOIFMasq(opt.LXDMngInterface))
	t.AddNATPostRoutingRule(newNATPostRoutingRuleSnmpTrapMasq())

	return t.Execute(f)
}

func NewNFTableCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nft",
		Short: "nftable command.",
	}

	nft := NewNFTableCmd()
	rootCmd.AddCommand(nft.setFlags(
		&cobra.Command{
			Use:   "create <playbook name>",
			Short: "Create new nftable.conf",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := nft.readConfig(); err != nil {
					return err
				}
				return nft.createConf(args[0])
			},
		},
	))
	return rootCmd
}
