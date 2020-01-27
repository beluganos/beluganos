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
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//
// VMAPICommand is VMAPI command
//
type VMAPICommand struct {
	Addr string
	ReID string

	// port-config
	ParentLink string
	MasterLink string
	PortStatus string
	DpPort     uint32
}

//
// GetPortStatus returns port status.
//
func (c *VMAPICommand) GetPortStatus() fibcapi.PortStatus_Status {
	v, _ := fibcapi.ParsePortStatus(c.PortStatus)
	return v
}

func (c *VMAPICommand) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.ReID, "reid", "", "", "router entity id")
	cmd.Flags().StringVarP(&c.Addr, "fibc-addr", "", FibcAddr, "FIBC address.")
	return cmd
}

func (c *VMAPICommand) setPortConfigFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.ParentLink, "parent", "", "", "parent ifname.")
	cmd.Flags().StringVarP(&c.MasterLink, "master", "", "", "master ifname.")
	cmd.Flags().StringVarP(&c.PortStatus, "port-status", "", "UP", "port status. <UP | DOWN>")
	cmd.Flags().Uint32VarP(&c.DpPort, "dp-port", "", 0, "dp port id.")
	return c.setFlags(cmd)
}

func (c *VMAPICommand) connect(f func(fibcapi.FIBCVmApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	return f(fibcapi.NewFIBCVmApiClient(conn))
}

func (c *VMAPICommand) hello() error {
	req := fibcapi.Hello{
		ReId: c.ReID,
	}

	return c.connect(func(client fibcapi.FIBCVmApiClient) error {
		if _, err := client.SendHello(context.Background(), &req); err != nil {
			return err
		}

		return nil
	})
}

func (c *VMAPICommand) monitor(done <-chan struct{}) error {
	req := fibcapi.VmMonitorRequest{
		ReId: c.ReID,
	}

	return c.connect(func(client fibcapi.FIBCVmApiClient) error {
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

			fibcapi.LogVmMonitorReply(
				log.StandardLogger(),
				log.DebugLevel,
				reply,
			)
		}

		return nil
	})
}

func (c *VMAPICommand) portConfig(cmd fibcapi.PortConfig_Cmd, ifname string, portStr string) error {
	portID, err := strconv.ParseUint(portStr, 10, 32)
	if err != nil {
		return err
	}

	req := fibcapi.PortConfig{
		Cmd:    cmd,
		ReId:   c.ReID,
		Ifname: ifname,
		PortId: uint32(portID),
		Link:   c.ParentLink,
		Master: c.MasterLink,
		Status: c.GetPortStatus(),
		DpPort: c.DpPort,
	}

	return c.connect(func(client fibcapi.FIBCVmApiClient) error {
		if _, err := client.SendPortConfig(context.Background(), &req); err != nil {
			return err
		}

		return nil
	})
}

func vmAPICmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "vmapi",
		Aliases: []string{"vm"},
		Short:   "FIBC VM API command.",
	}

	vmapi := VMAPICommand{}

	rootCmd.AddCommand(vmapi.setFlags(
		&cobra.Command{
			Use:   "hello",
			Short: "send hello message.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vmapi.hello()
			},
		},
	))

	rootCmd.AddCommand(vmapi.setFlags(
		&cobra.Command{
			Use:     "monitor",
			Aliases: []string{"mon"},
			Short:   "monitor fibc.",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vmapi.monitor(nil)
			},
		},
	))

	rootCmd.AddCommand(vmAPIPortConfigCmd())

	return rootCmd
}

func vmAPIPortConfigCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "port-config",
		Aliases: []string{"pc"},
		Short:   "send port config.",
	}

	vmapi := VMAPICommand{}

	rootCmd.AddCommand(vmapi.setPortConfigFlags(
		&cobra.Command{
			Use:   "add <ifname> <port-id>",
			Short: "send add port config.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vmapi.portConfig(fibcapi.PortConfig_ADD, args[0], args[1])
			},
		},
	))

	rootCmd.AddCommand(vmapi.setPortConfigFlags(
		&cobra.Command{
			Use:   "del <ifname> <port-id>",
			Short: "send del port config.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vmapi.portConfig(fibcapi.PortConfig_DELETE, args[0], args[1])
			},
		},
	))

	rootCmd.AddCommand(vmapi.setPortConfigFlags(
		&cobra.Command{
			Use:     "modify <ifname> <port-id>",
			Short:   "send modify port config.",
			Aliases: []string{"mod", "update"},
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vmapi.portConfig(fibcapi.PortConfig_DELETE, args[0], args[1])
			},
		},
	))

	return rootCmd
}
