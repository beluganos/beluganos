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
	"fabricflow/ribt/api/ribtapi"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RibtAPICommand struct {
	Addr string
}

func (c *RibtAPICommand) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.Addr, "ribt-addr", "", "localhost:50061", "RIBT address.")
	return cmd
}

func (c *RibtAPICommand) connect(f func(ribtapi.RIBTApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	return f(ribtapi.NewRIBTApiClient(conn))
}

func (c *RibtAPICommand) dump() error {

	return c.connect(func(client ribtapi.RIBTApiClient) error {
		stream, err := client.GetTunnels(context.Background(), &ribtapi.GetTunnelsRequest{})
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

			fmt.Printf("id:%d type:%d %s->%s\n", e.Id, e.Type, e.Local, e.Remote)
			if routes := e.Routes; routes != nil {
				for key, route := range routes {
					fmt.Printf("[%s] prefix:%s nexthop:%s family:%d type:%d\n",
						key, route.Prefix, route.Nexthop, route.Family, route.TunnelType)
				}
			}
		}

		return nil
	})
}

func ribtAPICmd() *cobra.Command {
	api := RibtAPICommand{}

	rootCmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"dump"},
		Short:   "RIBT API show command.",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.dump()
		},
	}

	return rootCmd
}
