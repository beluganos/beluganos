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
// APAPICommand is APAPI command
//
type APAPICommand struct {
	Addr string
}

func (c *APAPICommand) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.Addr, "fibc-addr", "", "localhost:50061", "FIBC address.")
	return cmd
}

func (c *APAPICommand) connect(f func(fibcapi.FIBCApApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	return f(fibcapi.NewFIBCApApiClient(conn))
}

func (c *APAPICommand) monitor(done <-chan struct{}) error {
	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		req := fibcapi.ApMonitorRequest{}
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

			fibcapi.LogApMonitorReply(
				log.StandardLogger(),
				log.DebugLevel,
				reply,
			)
		}

		return nil
	})
}

func (c *APAPICommand) getPortStats(dpID uint64, portID uint32) error {
	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetPortStatsRequest{
			DpId:   dpID,
			PortNo: portID,
		}
		stream, err := client.GetPortStats(context.Background(), &req)
		if err != nil {
			return err
		}

	FOR_LOOP:
		for {
			stats, err := stream.Recv()
			if err == io.EOF {
				break FOR_LOOP
			}
			if err != nil {
				return err
			}
			if stats == nil {
				continue FOR_LOOP
			}

			fmt.Printf("Port: %d\n", stats.PortNo)
			for k, v := range stats.Values {
				fmt.Printf("%s: %v\n", k, v)
			}
			for k, v := range stats.SValues {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
		return nil
	})
}

func (c *APAPICommand) modPortStats(dpID uint64, portID uint32, cmd fibcapi.FFPortStats_Cmd) error {
	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApModPortStatsRequest{
			DpId:   dpID,
			PortNo: portID,
			Cmd:    cmd,
		}
		if _, err := client.ModPortStats(context.Background(), &req); err != nil {
			return err
		}

		return nil
	})
}

func (c *APAPICommand) oamAuditRouteCnt() error {
	return c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.OAM_Request{
			DpId:    0,
			OamType: fibcapi.OAM_AUDIT_ROUTE_CNT,
			Body: &fibcapi.OAM_Request_AuditRouteCnt{
				AuditRouteCnt: &fibcapi.OAM_AuditRouteCntRequest{},
			},
		}
		if _, err := client.RunOAM(context.Background(), &req); err != nil {
			return err
		}

		return nil
	})
}

func apAPICmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "apapi",
		Aliases: []string{"ap"},
		Short:   "FIBC AP API command.",
	}

	apapi := APAPICommand{}

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "monitor",
			Aliases: []string{"mon"},
			Short:   "monitor fibcd",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.monitor(nil)
			},
		},
	))

	rootCmd.AddCommand(
		apAPIPortStatsCmd(),
		apAPIOAMCmd(),
	)

	return rootCmd
}

func apAPIPortStatsArgs(args []string) (uint64, uint32, error) {
	dpID, err := strconv.ParseUint(args[0], 0, 64)
	if err != nil {
		return 0, 0, err
	}

	portID := uint32(0xffffffff)
	if len(args) == 2 {
		v, err := strconv.ParseUint(args[1], 0, 32)
		if err != nil {
			return 0, 0, err
		}

		portID = uint32(v)
	}

	return dpID, portID, nil
}

func apAPIPortStatsCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "port-stats",
		Aliases: []string{"ps"},
		Short:   "port stats command.",
	}

	apapi := APAPICommand{}

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "show <dp-id> [port-id]",
			Aliases: []string{"get", "dump"},
			Short:   "show port stats",
			Args:    cobra.RangeArgs(1, 2),
			RunE: func(cmd *cobra.Command, args []string) error {
				dpID, portID, err := apAPIPortStatsArgs(args)
				if err != nil {
					return err
				}

				return apapi.getPortStats(dpID, portID)
			},
		},
	))

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:     "clear <dp-id> [port-id]",
			Aliases: []string{"reset", "cls"},
			Short:   "clear port stats",
			Args:    cobra.RangeArgs(1, 2),
			RunE: func(cmd *cobra.Command, args []string) error {
				dpID, portID, err := apAPIPortStatsArgs(args)
				if err != nil {
					return err
				}

				return apapi.modPortStats(dpID, portID, fibcapi.FFPortStats_RESET)
			},
		},
	))

	return rootCmd
}

func apAPIOAMCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "oam",
		Short: "oam command.",
	}

	apapi := APAPICommand{}

	rootCmd.AddCommand(apapi.setFlags(
		&cobra.Command{
			Use:   "audit-route-cnt",
			Short: "Audit Route(Count)",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return apapi.oamAuditRouteCnt()
			},
		},
	))

	return rootCmd
}
