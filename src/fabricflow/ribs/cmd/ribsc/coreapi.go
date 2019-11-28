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
	"strconv"

	"github.com/golang/protobuf/proto"
	gobgpapi "github.com/osrg/gobgp/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RibsCoreAPICommand struct {
	Addr string
}

func (c *RibsCoreAPICommand) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.Addr, "ribs-addr", "", "localhost:50061", "FIBC address.")
	return cmd
}

func (c *RibsCoreAPICommand) connect(f func(ribsapi.RIBSCoreApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	return f(ribsapi.NewRIBSCoreApiClient(conn))
}

func (c *RibsCoreAPICommand) monitor(nid uint8, rt string, done <-chan struct{}) error {
	return c.connect(func(client ribsapi.RIBSCoreApiClient) error {
		req := ribsapi.MonitorRibRequest{
			NId: uint32(nid),
			Rt:  rt,
		}

		stream, err := client.MonitorRib(context.Background(), &req)
		if err != nil {
			return err
		}

	FOR_LOOP:
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break FOR_LOOP
			}
			if err != nil {
				return err
			}
			if msg == nil {
				continue FOR_LOOP
			}

			path := gobgpapi.Path{}
			if err := proto.Unmarshal(msg.Path, &path); err != nil {
				return err
			}

			fmt.Printf("%s %v\n", msg.Rt, path)

		}

		return nil
	})
}

func ribsCoreAPICmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "coreapi",
		Aliases: []string{"core"},
		Short:   "RIBS Core API command.",
	}

	api := RibsCoreAPICommand{}

	rootCmd.AddCommand(api.setFlags(
		&cobra.Command{
			Use:     "monitor <nid> <RT>",
			Aliases: []string{"mon"},
			Short:   "monitor ribsd",
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				nid, err := strconv.ParseUint(args[0], 0, 8)
				if err != nil {
					return err
				}
				return api.monitor(uint8(nid), args[1], nil)
			},
		},
	))

	return rootCmd
}
