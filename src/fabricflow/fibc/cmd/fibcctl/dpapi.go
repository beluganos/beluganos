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
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//
// DPAPICommand is DPAPI command
//
type DPAPICommand struct {
	Addr string
	DpID uint64

	// hello
	DpType string

	// multipart
	Xid uint32

	// multipart.PortStats
	portStatsType string
}

func (c *DPAPICommand) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Uint64VarP(&c.DpID, "dpid", "i", 0, "daapath id.")
	cmd.Flags().StringVarP(&c.DpType, "dp-type", "", "OPENNSL", "datapath type.")
	cmd.Flags().StringVarP(&c.Addr, "fibc-addr", "a", "localhost:50061", "FIBC address.")
	return cmd
}

func (c *DPAPICommand) setMultipartFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Uint32VarP(&c.Xid, "xid", "x", 0, "xid.")
	return c.setFlags(cmd)
}

func (c *DPAPICommand) setMultipartPortFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.portStatsType, "file-type", "t", "yaml", "port status file type.")
	return c.setMultipartFlags(cmd)
}

func (c *DPAPICommand) connect(f func(fibcapi.FIBCDpApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fibcapi.NewFIBCDpApiClient(conn)
	return f(client)
}

func (c *DPAPICommand) hello() error {
	dpType, err := fibcapi.ParseDpType(c.DpType)
	if err != nil {
		return err
	}

	hello := fibcapi.FFHello{
		DpId:   c.DpID,
		DpType: dpType,
	}

	return c.connect(func(client fibcapi.FIBCDpApiClient) error {
		if _, err := client.SendHello(context.Background(), &hello); err != nil {
			return err
		}

		return nil
	})
}

func (c *DPAPICommand) packetIn(portID uint32, path string) error {
	data, err := func() ([]byte, error) {
		switch path {
		case "arp":
			return newPacketARP(), nil

		default:
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			return parseHexDump(f)
		}
	}()

	if err != nil {
		return err
	}

	msg := fibcapi.FFPacketIn{
		DpId:   c.DpID,
		PortNo: portID,
		Data:   data,
	}

	return c.connect(func(client fibcapi.FIBCDpApiClient) error {
		if _, err := client.SendPacketIn(context.Background(), &msg); err != nil {
			return err
		}

		return nil
	})
}

func (c *DPAPICommand) monitor(done <-chan struct{}) error {
	req := fibcapi.DpMonitorRequest{
		DpId: c.DpID,
	}

	return c.connect(func(client fibcapi.FIBCDpApiClient) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stream, err := client.Monitor(ctx, &req)
		if err != nil {
			return err
		}

		if done != nil {
			go func() {
				<-done
				cancel()
			}()
		}

	FOR_LOOP:
		for {
			reply, err := stream.Recv()
			if err == io.EOF {
				break FOR_LOOP
			}
			if err != nil {
				return err
			}

			fibcapi.LogDpMonitorReply(
				log.StandardLogger(),
				log.DebugLevel,
				reply,
				true,
			)
		}

		return nil
	})
}

func (c *DPAPICommand) multipartPortDesc(args []string) error {
	ffports := []*fibcapi.FFPort{}

	for _, arg := range args {
		portID, err := strconv.ParseUint(arg, 0, 32)
		if err != nil {
			return err
		}

		ffport := &fibcapi.FFPort{
			PortNo: uint32(portID),
		}

		ffports = append(ffports, ffport)
	}

	reply := fibcapi.DpMultipartReply{
		Xid: c.Xid,
		Reply: &fibcapi.FFMultipart_Reply{
			DpId:   c.DpID,
			MpType: fibcapi.FFMultipart_PORT_DESC,
			Body: &fibcapi.FFMultipart_Reply_PortDesc{
				PortDesc: &fibcapi.FFMultipart_PortDescReply{
					Internal: true,
					Port:     ffports,
				},
			},
		},
	}

	return c.connect(func(client fibcapi.FIBCDpApiClient) error {
		if _, err := client.SendMultipartReply(context.Background(), &reply); err != nil {
			return err
		}

		return nil
	})
}

func (c *DPAPICommand) multipartPortStat(path string) error {
	cfg := DPAPIConfig{}
	if err := cfg.ReadFile(path, c.portStatsType); err != nil {
		return err
	}

	stats := []*fibcapi.FFPortStats{}
	for _, ps := range cfg.Multipart.GetPortStats() {
		log.Debugf("ps %v", ps)
		stat := &fibcapi.FFPortStats{
			PortNo: ps.PortID(),
			Values: ps.GetValues(),
		}
		stats = append(stats, stat)
	}

	reply := fibcapi.DpMultipartReply{
		Xid: c.Xid,
		Reply: &fibcapi.FFMultipart_Reply{
			DpId:   c.DpID,
			MpType: fibcapi.FFMultipart_PORT,
			Body: &fibcapi.FFMultipart_Reply_Port{
				Port: &fibcapi.FFMultipart_PortReply{
					Stats: stats,
				},
			},
		},
	}

	return c.connect(func(client fibcapi.FIBCDpApiClient) error {
		if _, err := client.SendMultipartReply(context.Background(), &reply); err != nil {
			return err
		}

		return nil
	})
}

func (c *DPAPICommand) multipartPortStats(paths []string) error {
	for _, path := range paths {
		if err := c.multipartPortStat(path); err != nil {
			return err
		}
	}

	return nil
}

func dpAPICmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "dpapi",
		Aliases: []string{"dp"},
		Short:   "FIBC DP API command.",
	}

	dpapi := DPAPICommand{}

	rootCmd.AddCommand(dpapi.setFlags(
		&cobra.Command{
			Use:   "hello",
			Short: "send hello message.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return dpapi.hello()
			},
		},
	))

	rootCmd.AddCommand(dpapi.setFlags(
		&cobra.Command{
			Use:     "packetin <port id> <filename | arp>",
			Aliases: []string{"pkt", "pktin"},
			Short:   "send packet in message,",
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				portID, err := strconv.ParseUint(args[0], 0, 32)
				if err != nil {
					return err
				}
				return dpapi.packetIn(uint32(portID), args[1])
			},
		},
	))

	rootCmd.AddCommand(dpapi.setFlags(
		&cobra.Command{
			Use:     "monitor",
			Aliases: []string{"mon"},
			Short:   "monitor fibc.",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return dpapi.monitor(nil)
			},
		},
	))

	rootCmd.AddCommand(dpMultipartCmd())

	return rootCmd
}

func dpMultipartCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "multipart",
		Aliases: []string{"mp"},
		Short:   "multipart command",
	}

	dpapi := DPAPICommand{}

	rootCmd.AddCommand(dpapi.setMultipartFlags(
		&cobra.Command{
			Use:     "portdesc <portID, ...>",
			Aliases: []string{"port-desc", "pd"},
			Short:   "send port-desc message.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return dpapi.multipartPortDesc(args)
			},
		},
	))

	rootCmd.AddCommand(dpapi.setMultipartPortFlags(
		&cobra.Command{
			Use:     "portstats <file1, ...>",
			Aliases: []string{"port-stats", "ps"},
			Short:   "send port-stats message.",
			Args:    cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return dpapi.multipartPortStats(args)
			},
		},
	))

	return rootCmd
}
