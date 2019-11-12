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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//
// VSAPICommand is VSAPI command
//
type VSAPICommand struct {
	Addr string
	VsID uint64
}

func (c *VSAPICommand) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Uint64VarP(&c.VsID, "vsid", "", 0, "vswitch id.")
	cmd.Flags().StringVarP(&c.Addr, "fibc-addr", "", "localhost:50061", "FIBC address.")
	return cmd
}

func (c *VSAPICommand) connect(f func(fibcapi.FIBCVsApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	return f(fibcapi.NewFIBCVsApiClient(conn))
}

func (c *VSAPICommand) hello() error {
	return c.connect(func(client fibcapi.FIBCVsApiClient) error {
		hello := fibcapi.FFHello{
			DpId: c.VsID,
		}

		if _, err := client.SendHello(context.Background(), &hello); err != nil {
			return err
		}

		return nil
	})
}

func (c *VSAPICommand) monitor(done <-chan struct{}) error {
	return c.connect(func(client fibcapi.FIBCVsApiClient) error {
		req := fibcapi.VsMonitorRequest{
			VsId: c.VsID,
		}

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

			fibcapi.LogVsMonitorReply(
				log.StandardLogger(),
				log.DebugLevel,
				reply,
				true,
			)
		}

		return nil
	})
}

func (c *VSAPICommand) pktin(client fibcapi.FIBCVsApiClient, portID uint32, path string) error {
	data, err := parseHexDumpFile(path)
	if err != nil {
		return err
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		hexdumpDebugLog(data)
	}
	pktin := fibcapi.FFPacketIn{
		DpId:   c.VsID,
		PortNo: portID,
		Data:   data,
	}

	if _, err := client.SendPacketIn(context.Background(), &pktin); err != nil {
		return err
	}

	return nil
}

func (c *VSAPICommand) pktins(port string, paths []string) error {
	portID, err := strconv.ParseUint(port, 0, 32)
	if err != nil {
		return err
	}

	return c.connect(func(client fibcapi.FIBCVsApiClient) error {
		for _, path := range paths {
			if err := c.pktin(client, uint32(portID), path); err != nil {
				fmt.Printf("packet out error. %s", err)
			}
		}

		return nil
	})
}

func (c *VSAPICommand) ffpacket(reID, ifname, port string) error {
	portID, err := strconv.ParseUint(port, 0, 32)
	if err != nil {
		return err
	}

	return c.connect(func(client fibcapi.FIBCVsApiClient) error {
		req := fibcapi.FFPacket{
			DpId:   c.VsID,
			PortNo: uint32(portID),
			ReId:   reID,
			Ifname: ifname,
		}

		if _, err := client.SendFFPacket(context.Background(), &req); err != nil {
			return err
		}

		return nil
	})
}

func vsAPICmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "vsapi",
		Aliases: []string{"vs"},
		Short:   "FIBC VS API command.",
	}

	vsapi := VSAPICommand{}

	rootCmd.AddCommand(vsapi.setFlags(
		&cobra.Command{
			Use:   "hello",
			Short: "send hello message.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vsapi.hello()
			},
		},
	))

	rootCmd.AddCommand(vsapi.setFlags(
		&cobra.Command{
			Use:     "monitor",
			Aliases: []string{"mon"},
			Short:   "monitor fibcd.",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vsapi.monitor(nil)
			},
		},
	))

	rootCmd.AddCommand(vsapi.setFlags(
		&cobra.Command{
			Use:   "pktin <port-id> <files, ...>",
			Short: "send packet in.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vsapi.pktins(args[0], args[1:])
			},
		},
	))

	rootCmd.AddCommand(vsapi.setFlags(
		&cobra.Command{
			Use:   "ffpkt <re-id> <ifname> <port-id>",
			Short: "send ffpacket.",
			Args:  cobra.ExactArgs(3),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vsapi.ffpacket(args[0], args[1], args[2])
			},
		},
	))

	return rootCmd
}
