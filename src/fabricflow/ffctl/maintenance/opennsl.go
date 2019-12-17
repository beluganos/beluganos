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
	"context"
	"fabricflow/ffctl/fflib"
	"fabricflow/ffctl/maintenance/opennsl"
	fibcapi "fabricflow/fibc/api"
	"fmt"
	gonslapi "gonsl/api"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var openNSLDumperList = []opennsl.Dumper{
	opennsl.NewFieldEntry(),
	opennsl.NewPort(),
	opennsl.NewVlan(),
	opennsl.NewL2Addr(),
	opennsl.NewL3Iface(),
	opennsl.NewL3Egress(),
	opennsl.NewL3Host(),
	opennsl.NewL3Route(),
	opennsl.NewTunnelInitiator(),
	opennsl.NewTunnelTerminator(),
	opennsl.NewIDMap(),
}

type OpenNSLCmd struct {
	gonsl *fflib.GonslClient
	fibc  *fflib.FibcClient
}

func NewOpenNSLCmd() *OpenNSLCmd {
	return &OpenNSLCmd{
		gonsl: fflib.NewGonslClient(),
		fibc:  fflib.NewFibcClient(),
	}
}

func (c *OpenNSLCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.gonsl.Host, "addr", "a", fflib.GonslHost, "OpenNSL agent addr.")
	cmd.Flags().Uint16VarP(&c.gonsl.Port, "port", "p", fflib.GonslPort, "OpenNSL agent port.")
	return cmd
}

func (c *OpenNSLCmd) setAllFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Uint16VarP(&c.gonsl.Port, "port", "p", fflib.GonslPort, "OpenNSL agent port.")
	cmd.Flags().StringVarP(&c.fibc.Host, "fibc-addr", "", fflib.FibcHost, "fibcd addr.")
	cmd.Flags().Uint16VarP(&c.fibc.Port, "fibc-port", "", fflib.FibcPort, "fibcd port.")
	return cmd
}

func (c *OpenNSLCmd) dpAddrList() ([]string, error) {
	addrs := []string{}

	err := c.fibc.Connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetDpEntriesRequest{
			Type: fibcapi.DbDpEntry_DPMON,
		}
		stream, err := client.GetDpEntries(context.Background(), &req)
		if err != nil {
			return err
		}

	FOR_LOOP:
		for {
			e, err := stream.Recv()
			if err == io.EOF {
				break FOR_LOOP
			}
			if err != nil {
				return err
			}
			addrs = append(addrs, e.Remote)
		}

		return nil
	})

	return addrs, err
}

func (c *OpenNSLCmd) dumpNameList(w io.Writer) {
	for _, d := range openNSLDumperList {
		fmt.Fprintf(w, "%s\n", d.Name())
	}
}

func (c *OpenNSLCmd) dumpAll(w io.Writer) error {
	hosts, err := c.dpAddrList()
	if err != nil {
		return err
	}

	for _, host := range hosts {
		c.gonsl.Host = host
		if err := c.dumpDP(w); err != nil {
			return err
		}
	}

	return nil
}

func (c *OpenNSLCmd) dumpAllByName(w io.Writer, name string) error {
	hosts, err := c.dpAddrList()
	if err != nil {
		return err
	}

	for _, host := range hosts {
		c.gonsl.Host = host
		if err := c.dumpDPByName(w, name); err != nil {
			return err
		}
	}

	return nil
}

func (c *OpenNSLCmd) dumpDP(w io.Writer) error {
	return c.gonsl.Connect(func(client gonslapi.GoNSLApiClient) error {
		for _, d := range openNSLDumperList {
			if err := d.Dump(w, client); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *OpenNSLCmd) dumpDPByName(w io.Writer, name string) error {
	return c.gonsl.Connect(func(client gonslapi.GoNSLApiClient) error {
		for _, d := range openNSLDumperList {
			if d.Name() == name {
				return d.Dump(w, client)
			}
		}
		return fmt.Errorf("Unknown name. %s", name)
	})
}

func NewOpennNSLCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "opennsl",
		Short:   "opennsl agent command.",
		Aliases: []string{"onsl", "gonsl"},
	}

	onsl := NewOpenNSLCmd()

	rootCmd.AddCommand(onsl.setFlags(
		&cobra.Command{
			Use:     "dump-dp [names ...]",
			Aliases: []string{"show-dp"},
			Short:   "dump entries of opennsl agent.",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return onsl.dumpDP(os.Stdout)
				}

				for _, arg := range args {
					if err := onsl.dumpDPByName(os.Stdout, arg); err != nil {
						return err
					}
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(onsl.setAllFlags(
		&cobra.Command{
			Use:     "dump [names ...]",
			Aliases: []string{"show"},
			Short:   "dump entries of all opennsl agents.",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return onsl.dumpAll(os.Stdout)
				}

				for _, arg := range args {
					if err := onsl.dumpAllByName(os.Stdout, arg); err != nil {
						return err
					}
				}

				return nil
			},
		},
	))

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "data-names",
			Short: "show data names.",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				onsl.dumpNameList(os.Stdout)
			},
		},
	)

	return rootCmd
}
