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
	fibcapi "fabricflow/fibc/api"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type FibcCmd struct {
	fibc *fflib.FibcClient
}

func NewFibcCmd() *FibcCmd {
	return &FibcCmd{
		fibc: fflib.NewFibcClient(),
	}
}

func (c *FibcCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fibc.Host, "fibc-addr", "", fflib.FibcHost, "fibcd addr.")
	cmd.Flags().Uint16VarP(&c.fibc.Port, "fibc-port", "", fflib.FibcPort, "fibcd port.")
	return cmd
}

func (c *FibcCmd) dump(w io.Writer) error {
	return c.fibc.Connect(func(client fibcapi.FIBCApApiClient) error {
		if err := c.writePortEntries(w, client); err != nil {
			return err
		}

		if err := c.writeIDEntries(w, client); err != nil {
			return err
		}

		if err := c.writeDpEntries(w, client, fibcapi.DbDpEntry_NOP); err != nil {
			return err
		}

		return c.dumpStats(w)
	})
}

func (c *FibcCmd) dumpPortEntries(w io.Writer) error {
	return c.fibc.Connect(func(client fibcapi.FIBCApApiClient) error {
		return c.writePortEntries(w, client)
	})
}

func (c *FibcCmd) writePortEntries(w io.Writer, client fibcapi.FIBCApApiClient) error {
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

		fmt.Fprintf(w, "%s VM%s DP%s VS%s P=%s M=%s\n",
			strkey(e.Key),
			strport(e.VmPort),
			strport(e.DpPort),
			strport(e.VsPort),
			strkey(e.ParentKey),
			strkey(e.MasterKey),
		)
	}

	return nil
}

func (c *FibcCmd) dumpIDEntries(w io.Writer) error {
	return c.fibc.Connect(func(client fibcapi.FIBCApApiClient) error {
		return c.writeIDEntries(w, client)
	})
}

func (c *FibcCmd) writeIDEntries(w io.Writer, client fibcapi.FIBCApApiClient) error {
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

		fmt.Fprintf(w, "dp_id:%d re_id:'%s'\n",
			e.DpId,
			e.ReId,
		)
	}

	return nil
}

func (c *FibcCmd) dumpDpEntries(w io.Writer, t fibcapi.DbDpEntry_Type) error {
	return c.fibc.Connect(func(client fibcapi.FIBCApApiClient) error {
		return c.writeDpEntries(w, client, t)
	})
}

func (c *FibcCmd) writeDpEntries(w io.Writer, client fibcapi.FIBCApApiClient, t fibcapi.DbDpEntry_Type) error {
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

		fmt.Fprintf(w, "Id:%s type:%s remote:%s\n",
			e.Id,
			e.Type,
			e.Remote,
		)
	}

	return nil
}

func (c *FibcCmd) dumpStats(w io.Writer) error {
	return c.fibc.Connect(func(client fibcapi.FIBCApApiClient) error {
		return c.writeStats(w, client)
	})
}

func (c *FibcCmd) writeStats(w io.Writer, client fibcapi.FIBCApApiClient) error {
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

		fmt.Fprintf(w, "%s:%s = %d\n",
			e.Group,
			e.Name,
			e.Value,
		)
	}

	return nil
}

func NewFibcCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fibc",
		Short: "maintenance fibc commands.",
	}

	fibc := NewFibcCmd()

	rootCmd.AddCommand(fibc.setFlags(
		&cobra.Command{
			Use:   "dump",
			Short: "dump fibc datas",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return fibc.dump(os.Stdout)
			},
		},
	))

	rootCmd.AddCommand(fibc.setFlags(
		&cobra.Command{
			Use:   "dump-portmap",
			Short: "dump fibc portmap datas.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return fibc.dumpPortEntries(os.Stdout)
			},
		},
	))

	rootCmd.AddCommand(fibc.setFlags(
		&cobra.Command{
			Use:   "dump-idmap",
			Short: "dump fibc idmap datas.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return fibc.dumpIDEntries(os.Stdout)
			},
		},
	))

	rootCmd.AddCommand(fibc.setFlags(
		&cobra.Command{
			Use:   "dump-dpset [type]",
			Short: "dump fibc dpset datas.",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return fibc.dumpDpEntries(os.Stdout, fibcapi.DbDpEntry_NOP)
				}

				t, err := fibcapi.ParseDbDpEntryType(args[0])
				if err != nil {
					return err
				}

				return fibc.dumpDpEntries(os.Stdout, t)
			},
		},
	))

	rootCmd.AddCommand(fibc.setFlags(
		&cobra.Command{
			Use:   "dump-stats",
			Short: "show stats entries",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return fibc.dumpStats(os.Stdout)
			},
		},
	))

	return rootCmd
}
