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
	"fabricflow/ribs/api/ribsapi"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RibsAPICommand struct {
	Addr string
}

func (c *RibsAPICommand) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.Addr, "ribs-addr", "", "localhost:50061", "FIBC address.")
	return cmd
}

func (c *RibsAPICommand) connect(f func(ribsapi.RIBSApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	return f(ribsapi.NewRIBSApiClient(conn))
}

func (c *RibsAPICommand) dumpRics() error {
	return c.connect(func(client ribsapi.RIBSApiClient) error {
		stream, err := client.GetRics(context.Background(), &ribsapi.GetRicsRequest{})
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
			if e == nil {
				continue FOR_LOOP
			}

			fmt.Printf("RIC[%s]: nid:%d RT:%s\n", e.Key, e.NId, e.Rt)
		}

		return nil
	})
}

func (c *RibsAPICommand) dumpNexthops() error {
	return c.connect(func(client ribsapi.RIBSApiClient) error {
		stream, err := client.GetNexthops(context.Background(), &ribsapi.GetNexthopsRequest{})
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
			if e == nil {
				continue FOR_LOOP
			}

			fmt.Printf("Nexthop[%s] %s RT:%s src:%s\n", e.Key, e.Addr, e.Rt, e.SourceId)
		}

		return nil
	})
}

func (c *RibsAPICommand) dumpNexthopMap() error {
	return c.connect(func(client ribsapi.RIBSApiClient) error {
		stream, err := client.GetNexthopMap(context.Background(), &ribsapi.GetIPMapRequest{})
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
			if e == nil {
				continue FOR_LOOP
			}

			fmt.Printf("NexthopMap: %s = %s\n", e.Key, e.Value)
		}

		return nil
	})
}

func ribsAPICmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "appapi",
		Aliases: []string{"api", "app"},
		Short:   "RIBS API command.",
	}

	return rootCmd
}

func ribsAPIDumpCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"dump"},
		Short:   "RIBS API show command.",
	}

	api := RibsAPICommand{}

	rootCmd.AddCommand(api.setFlags(
		&cobra.Command{
			Use:     "nexthops",
			Aliases: []string{"nh"},
			Short:   "show nexthops.",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return api.dumpNexthops()
			},
		},
	))

	rootCmd.AddCommand(api.setFlags(
		&cobra.Command{
			Use:     "nexthop-map",
			Aliases: []string{"nm"},
			Short:   "show nexthop map.",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return api.dumpNexthopMap()
			},
		},
	))

	rootCmd.AddCommand(api.setFlags(
		&cobra.Command{
			Use:     "rics",
			Aliases: []string{"r", "rs", "ric"},
			Short:   "show rics.",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return api.dumpRics()
			},
		},
	))

	return rootCmd
}
