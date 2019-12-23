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
	fibcapi "fabricflow/fibc/api"
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
)

func (c *APAPICommand) addPortMap(reID, ifname, name, dpIDs, portIDs string) error {
	dpID, err := strconv.ParseUint(dpIDs, 0, 64)
	if err != nil {
		return err
	}

	portID, err := strconv.ParseUint(portIDs, 0, 32)
	if err != nil {
		return err
	}

	e := fibcapi.DbPortEntry{
		Key: &fibcapi.DbPortKey{
			ReId:   reID,
			Ifname: ifname,
		},
		DpPort: &fibcapi.DbPortValue{
			DpId:   dpID,
			PortId: uint32(portID),
		},
	}

	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		_, err := client.AddPortEntry(context.Background(), &e)
		return err
	})
}

func (c *APAPICommand) delPortMap(reID, ifname string) error {
	key := fibcapi.DbPortKey{
		ReId:   reID,
		Ifname: ifname,
	}

	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		_, err := client.DelPortEntry(context.Background(), &key)
		return err
	})
}

func (c *APAPICommand) addIDMap(reID, dpIDs string) error {
	dpID, err := strconv.ParseUint(dpIDs, 0, 64)
	if err != nil {
		return err
	}

	e := fibcapi.DbIdEntry{
		ReId: reID,
		DpId: dpID,
	}

	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		_, err := client.AddIDEntry(context.Background(), &e)
		return err
	})
}

func (c *APAPICommand) delIDMap(reID, dpIDs string) error {
	dpID, err := strconv.ParseUint(dpIDs, 0, 64)
	if err != nil {
		return err
	}

	e := fibcapi.DbIdEntry{
		ReId: reID,
		DpId: dpID,
	}

	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		_, err := client.DelIDEntry(context.Background(), &e)
		return err
	})
}

func (c *APAPICommand) dumpPortEntries() error {
	strbool := func(b bool) string {
		if b {
			return "+"
		}
		return "-"
	}

	strkey := func(k *fibcapi.DbPortKey) string {
		if k == nil {
			return "{}"
		}
		return fmt.Sprintf("{'%s', '%s'}", k.ReId, k.Ifname)
	}

	strport := func(p *fibcapi.DbPortValue) string {
		if p == nil {
			return "+{}"
		}
		if len(p.ReId) != 0 {
			return fmt.Sprintf("%s{'%s',0x%x}", strbool(p.Enter), p.ReId, p.PortId)
		}
		return fmt.Sprintf("%s{%d,0x%x}", strbool(p.Enter), p.DpId, p.PortId)
	}

	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetPortEntriesRequest{}

		stream, err := client.GetPortEntries(context.Background(), &req)
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

			fmt.Printf("%s VM%s DP%s VS%s P=%s M=%s\n",
				strkey(e.Key),
				strport(e.VmPort),
				strport(e.DpPort),
				strport(e.VsPort),
				strkey(e.ParentKey),
				strkey(e.MasterKey),
			)
		}

		return nil
	})
}

func (c *APAPICommand) dumpIDEntries() error {
	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetIdEntriesRequest{}

		stream, err := client.GetIDEntries(context.Background(), &req)
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

			fmt.Printf("dp_id:%d re_id:'%s'\n",
				e.DpId,
				e.ReId,
			)
		}

		return nil
	})
}

func (c *APAPICommand) dumpDpEntries(t fibcapi.DbDpEntry_Type) error {
	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetDpEntriesRequest{
			Type: t,
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

<<<<<<< HEAD
<<<<<<< HEAD
			fmt.Printf("Id:%s type:%s\n",
				e.Id,
				e.Type,
=======
=======
>>>>>>> develop
			fmt.Printf("Id:%s type:%s remote:%s\n",
				e.Id,
				e.Type,
				e.Remote,
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
			)
		}

		return nil
	})
}

func (c *APAPICommand) dumpStats() error {
	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetStatsRequest{}

		stream, err := client.GetStats(context.Background(), &req)
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

			fmt.Printf("%s:%s = %d\n",
				e.Group,
				e.Name,
				e.Value,
			)
		}

		return nil
	})
}

func dbCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:     "database",
		Aliases: []string{"db"},
		Short:   "FIBC database command.",
	}

	rootCmd.AddCommand(
		dbShowCmd(),
		dbPortMapCmd(),
		dbIDMapCmd(),
	)

	return rootCmd
}

func dbPortMapCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "portmap",
		Aliases: []string{"port", "p"},
		Short:   "FIBC portmap command.",
	}

	apapi := APAPICommand{}

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "add <re-id> <ifname> <dp-id> <dpport-id>",
			Short:   "add portmap entry",
			Aliases: []string{"a"},
			Args:    cobra.ExactArgs(5),
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.addPortMap(args[0], args[1], args[2], args[3], args[4])
			},
		},
	))

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "delete <reid> <ifname>",
			Short:   "delete portmap entry",
			Aliases: []string{"d", "del"},
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.delPortMap(args[0], args[1])
			},
		},
	))

	return rootCmd
}

func dbIDMapCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "idmap",
		Aliases: []string{"id", "i"},
		Short:   "FIBC idmap command.",
	}

	apapi := APAPICommand{}

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "add <re-id> <dp-id>",
			Short:   "add idmap entry",
			Aliases: []string{"a"},
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.addIDMap(args[0], args[1])
			},
		},
	))

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "delete <re-id> <dp-id>",
			Short:   "delete portmap entry",
			Aliases: []string{"d", "del"},
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.delIDMap(args[0], args[1])
			},
		},
	))

	return rootCmd
}

func dbShowCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"dump"},
		Short:   "FIBC databse show command.",
	}

	apapi := APAPICommand{}

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "portmap",
			Aliases: []string{"p", "port"},
			Short:   "show port entries",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.dumpPortEntries()
			},
		},
	))

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "idmap",
			Aliases: []string{"i", "id"},
			Short:   "show idmap entries",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.dumpIDEntries()
			},
		},
	))

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "dpset [type]",
			Aliases: []string{"d", "dp"},
			Short:   "show dpset entries",
			Args:    cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return apapi.dumpDpEntries(fibcapi.DbDpEntry_NOP)
				}

				t, err := fibcapi.ParseDbDpEntryType(args[0])
				if err != nil {
					return err
				}

				return apapi.dumpDpEntries(t)
			},
		},
	))

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "stats",
			Aliases: []string{"st"},
			Short:   "show stats entries",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.dumpStats()
			},
		},
	))

	return rootCmd
}
