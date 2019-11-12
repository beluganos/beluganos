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
	"context"
	"fmt"
	"govsw/api/vswapi"
	"govsw/pkgs/govsw"
	"io"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type GovswCmd struct {
	Addr string
	Port uint16
}

func (c *GovswCmd) getGovswdAddr() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

func (c *GovswCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.Addr, "govswd-addr", "a", "localhost", "govswd address.")
	cmd.Flags().Uint16VarP(&c.Port, "govswd-port", "", govsw.VSWAPI_PORT, "govswd port")
	return cmd
}

func (c *GovswCmd) connect(f func(vswapi.VswApiClient) error) error {
	conn, err := grpc.Dial(c.getGovswdAddr(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := vswapi.NewVswApiClient(conn)
	return f(client)
}

func (c *GovswCmd) dumpLink(client vswapi.VswApiClient) error {
	stream, err := client.GetLinks(context.Background(), &vswapi.GetLinksRequest{})
	if err != nil {
		log.Errorf("GetLinks error. %s", err)
		return err
	}

FOR_LOOP:
	for {
		link, err := stream.Recv()
		if err == io.EOF {
			break FOR_LOOP
		}
		if err != nil {
			return err
		}

		fmt.Printf("link : %s\n", link)
	}

	return nil
}

func (c *GovswCmd) dumpIfname(client vswapi.VswApiClient) error {
	stream, err := client.GetIfnames(context.Background(), &vswapi.GetIfnamesRequest{})
	if err != nil {
		log.Errorf("GetLinks error. %s", err)
		return err
	}

FOR_LOOP:
	for {
		ifname, err := stream.Recv()
		if err == io.EOF {
			break FOR_LOOP
		}
		if err != nil {
			return err
		}

		fmt.Printf("iface: %s\n", ifname)
	}

	return nil
}

func (c *GovswCmd) dumpStats(client vswapi.VswApiClient) error {
	stream, err := client.GetStats(context.Background(), &vswapi.GetStatsRequest{})
	if err != nil {
		log.Errorf("GetStats error. %s", err)
		return err
	}

FOR_LOOP:
	for {
		stat, err := stream.Recv()
		if err == io.EOF {
			break FOR_LOOP
		}
		if err != nil {
			return err
		}

		keys := []string{}
		for key := range stat.Values {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			fmt.Printf("stats: %s/%s: %d\n", stat.Group, key, stat.Values[key])
		}
	}

	return nil
}

func (c *GovswCmd) dump(targets []string) error {
	if len(targets) == 0 {
		targets = []string{"link", "ifname"}
	}

	return c.connect(func(client vswapi.VswApiClient) error {
		for _, target := range targets {
			switch target {
			case "link":
				return c.dumpLink(client)

			case "ifname":
				return c.dumpIfname(client)

			case "stats":
				return c.dumpStats(client)

			default:
				log.Errorf("bad target '%s'", target)
				return fmt.Errorf("unknown target. %s", target)
			}
		}

		return nil
	})
}

var modLinkCmd_values = map[string]vswapi.ModLinkRequest_Cmd{
	"up":   vswapi.ModLinkRequest_UP,
	"down": vswapi.ModLinkRequest_DOWN,
}

func parseLinkCmd(cmd string) (vswapi.ModLinkRequest_Cmd, error) {
	if v, ok := modLinkCmd_values[cmd]; ok {
		return v, nil
	}

	return vswapi.ModLinkRequest_NOP, fmt.Errorf("invalid command. %s", cmd)
}

func (c *GovswCmd) modIface(cmd vswapi.ModIfnameRequest_Cmd, ifname string) error {
	req := vswapi.ModIfnameRequest{
		Cmd:    cmd,
		Ifname: ifname,
	}

	return c.connect(func(client vswapi.VswApiClient) error {
		if _, err := client.ModIfname(context.Background(), &req); err != nil {
			log.Errorf("ModIfname error. %s", err)
			return err
		}

		log.Debugf("ModIfname success. %s %s", cmd, ifname)
		return nil
	})
}

func (c *GovswCmd) modLink(cmd vswapi.ModLinkRequest_Cmd, ifname string) error {
	req := vswapi.ModLinkRequest{
		Cmd:    cmd,
		Ifname: ifname,
	}

	return c.connect(func(client vswapi.VswApiClient) error {
		if _, err := client.ModLink(context.Background(), &req); err != nil {
			log.Errorf("ModLink error. %s", err)
			return err
		}

		log.Debugf("ModLink success. %s %s", cmd, ifname)
		return nil
	})
}

func (c *GovswCmd) saveConfig() error {
	req := vswapi.SaveConfigRequest{}

	return c.connect(func(client vswapi.VswApiClient) error {
		if _, err := client.SaveConfig(context.Background(), &req); err != nil {
			log.Errorf("SaveConfig error. %s", err)
			return err
		}

		log.Debugf("SaveConfig success.")
		return nil
	})
}

func dumpCmd() *cobra.Command {

	vswcmd := GovswCmd{}

	rootCmd := vswcmd.setFlags(
		&cobra.Command{
			Use:   "dump <link | ifname | stats>",
			Short: "dump command.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return vswcmd.dump(args)
			},
		},
	)

	return rootCmd
}

func ifaceCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "interface",
		Aliases: []string{"i", "iface", "intf"},
		Short:   "interface command",
	}

	vswcmd := GovswCmd{}

	rootCmd.AddCommand(vswcmd.setFlags(
		&cobra.Command{
			Use:   "add <ifname>",
			Short: "add interface name command.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vswcmd.modIface(vswapi.ModIfnameRequest_ADD, args[0])
			},
		},
	))

	rootCmd.AddCommand(vswcmd.setFlags(
		&cobra.Command{
			Use:   "add-regex <regex>",
			Short: "add interface regex command.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vswcmd.modIface(vswapi.ModIfnameRequest_REG, args[0])
			},
		},
	))

	rootCmd.AddCommand(vswcmd.setFlags(
		&cobra.Command{
			Use:     "delete <ifname | regex>",
			Aliases: []string{"del"},
			Short:   "delete interface name or regex command.",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vswcmd.modIface(vswapi.ModIfnameRequest_DELETE, args[0])
			},
		},
	))

	rootCmd.AddCommand(vswcmd.setFlags(
		&cobra.Command{
			Use:   "sync [ifname]",
			Short: "sync command.",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ifname := func() string {
					if len(args) == 1 {
						return args[0]
					}
					return "*"
				}()
				return vswcmd.modIface(vswapi.ModIfnameRequest_SYNC, ifname)
			},
		},
	))

	return rootCmd
}

func linkCmd() *cobra.Command {
	vswcmd := GovswCmd{}

	rootCmd := vswcmd.setFlags(
		&cobra.Command{
			Use:   "link <up | down> <ifname>",
			Short: "link command",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := parseLinkCmd(args[0])
				if err != nil {
					return err
				}
				return vswcmd.modLink(c, args[1])
			},
		},
	)

	return rootCmd
}

func configCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "config",
		Short:   "config command.",
		Aliases: []string{"cfg", "conf"},
	}

	vswcmd := GovswCmd{}

	rootCmd.AddCommand(vswcmd.setFlags(
		&cobra.Command{
			Use:     "save",
			Short:   "save current setting to file.",
			Aliases: []string{"write"},
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vswcmd.saveConfig()
			},
		},
	))

	return rootCmd
}
