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
	fflibnet "fabricflow/util/net"
	"fmt"
	"net"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//
// TestStressCmd is stress test command.
//
type TestStressCmd struct {
	Addr   string
	IfNum  uint32
	IfBase uint32

	ReID     string
	VsID     uint64
	BaseAddr string
}

func (c *TestStressCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Uint32VarP(&c.IfNum, "interface", "i", 1, "interface num.")
	cmd.Flags().Uint32VarP(&c.IfBase, "base-iinterface", "b", 1, "base interface num.")
	cmd.Flags().StringVarP(&c.Addr, "fibc-addr", "", "localhost:50061", "FIBC address.")
	return cmd
}

func (c *TestStressCmd) setRunL3Flags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.ReID, "reid", "", "", "router entity id")
	cmd.Flags().StringVarP(&c.BaseAddr, "base-addr", "", "10.0.0.0/32", "base address.")
	return c.setFlags(cmd)
}

func (c *TestStressCmd) setRunMplsFlags(cmd *cobra.Command) *cobra.Command {
	// TODO:
	return c.setFlags(cmd)
}

func (c *TestStressCmd) setRunPacketFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Uint64VarP(&c.VsID, "vsid", "", 1, "vs id.")
	return c.setFlags(cmd)
}

func (c *TestStressCmd) connectVMAPI(f func(fibcapi.FIBCVmApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	return f(fibcapi.NewFIBCVmApiClient(conn))
}

func (c *TestStressCmd) connectVSAPI(f func(fibcapi.FIBCVsApiClient) error) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	return f(fibcapi.NewFIBCVsApiClient(conn))
}

func parseFlowModCmd(cmd string) fibcapi.FlowMod_Cmd {
	switch cmd {
	case "add":
		return fibcapi.FlowMod_ADD
	case "del":
		return fibcapi.FlowMod_DELETE
	default:
		return fibcapi.FlowMod_NOP
	}
}

func parseGroupModCmd(cmd string) fibcapi.GroupMod_Cmd {
	switch cmd {
	case "add":
		return fibcapi.GroupMod_ADD
	case "del":
		return fibcapi.GroupMod_DELETE
	default:
		return fibcapi.GroupMod_NOP
	}
}

func (c *TestStressCmd) setupL2InterfaceGroup(cmd string) error {
	l2if := fibcapi.L2InterfaceGroup{
		PortId:  0,
		VlanVid: 1,
	}

	mod := fibcapi.GroupMod{
		Cmd:   parseGroupModCmd(cmd),
		GType: fibcapi.GroupMod_L2_INTERFACE,
		ReId:  c.ReID,
		Entry: &fibcapi.GroupMod_L2Iface{
			L2Iface: &l2if,
		},
	}

	return c.connectVMAPI(func(client fibcapi.FIBCVmApiClient) error {
		ctxt := context.Background()
		for index := uint32(0); index < c.IfNum; index++ {
			l2if.PortId = index + 1
			l2if.HwAddr = fmt.Sprintf("00:01:00:00:00:%02x", uint8(index%0xff))

			if _, err := client.SendGroupMod(ctxt, &mod); err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *TestStressCmd) setupL3UnicastGroup(cmd string) error {
	l3uc := fibcapi.L3UnicastGroup{
		NeId:      0,
		PortId:    0,
		VlanVid:   1,
		PhyPortId: 0,
	}

	mod := fibcapi.GroupMod{
		Cmd:   parseGroupModCmd(cmd),
		GType: fibcapi.GroupMod_L3_UNICAST,
		ReId:  c.ReID,
		Entry: &fibcapi.GroupMod_L3Unicast{
			L3Unicast: &l3uc,
		},
	}

	return c.connectVMAPI(func(client fibcapi.FIBCVmApiClient) error {
		ctxt := context.Background()
		for index := uint32(0); index < c.IfNum; index++ {
			l3uc.NeId = c.IfBase + (index % c.IfNum)
			l3uc.PortId = index + 1
			l3uc.PhyPortId = index + 1
			l3uc.EthDst = fmt.Sprintf("00:02:00:00:00:%02x", uint8(index%0xff))
			l3uc.EthSrc = fmt.Sprintf("00:01:00:00:00:%02x", uint8(index%0xff))

			if _, err := client.SendGroupMod(ctxt, &mod); err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *TestStressCmd) setup() error {
	if err := c.setupL2InterfaceGroup("add"); err != nil {
		return err
	}

	if err := c.setupL3UnicastGroup("add"); err != nil {
		return err
	}

	return nil
}

func (c *TestStressCmd) teardown() error {
	if err := c.setupL3UnicastGroup("del"); err != nil {
		return err
	}

	if err := c.setupL2InterfaceGroup("del"); err != nil {
		return err
	}

	return nil
}

func (c *TestStressCmd) runRoute(cmd string, num uint32) error {
	_, ipdst, err := net.ParseCIDR(c.BaseAddr)
	if err != nil {
		return err
	}

	ucast := fibcapi.UnicastRoutingFlow{
		Match: &fibcapi.UnicastRoutingFlow_Match{
			Origin: fibcapi.UnicastRoutingFlow_ROUTE,
		},
		Action: &fibcapi.UnicastRoutingFlow_Action{
			Name: fibcapi.UnicastRoutingFlow_Action_DEC_TTL,
		},
		GType: fibcapi.GroupMod_L3_UNICAST,
	}

	mod := fibcapi.FlowMod{
		Cmd:   parseFlowModCmd(cmd),
		Table: fibcapi.FlowMod_UNICAST_ROUTING,
		ReId:  c.ReID,
		Entry: &fibcapi.FlowMod_Unicast{
			Unicast: &ucast,
		},
	}

	return c.connectVMAPI(func(client fibcapi.FIBCVmApiClient) error {
		ctxt := context.Background()

		w := NewStopWatch()
		w.Start()

		for index := uint32(0); index < num; index++ {
			fflibnet.IncIPNet(ipdst)

			ucast.GId = c.IfBase + (index % c.IfNum)
			ucast.Match.IpDst = ipdst.String()

			if _, err := client.SendFlowMod(ctxt, &mod); err != nil {
				return err
			}

		}

		w.Point()
		w.PrintResults()

		return nil
	})
}

func (c *TestStressCmd) runHost(cmd string, num uint32) error {
	flowCmd := func() fibcapi.FlowMod_Cmd {
		switch cmd {
		case "add":
			return fibcapi.FlowMod_ADD
		case "del":
			return fibcapi.FlowMod_DELETE
		default:
			return fibcapi.FlowMod_NOP
		}
	}()

	_, ipdst, err := net.ParseCIDR(c.BaseAddr)
	if err != nil {
		return err
	}

	ucast := fibcapi.UnicastRoutingFlow{
		Match: &fibcapi.UnicastRoutingFlow_Match{
			Origin: fibcapi.UnicastRoutingFlow_NEIGH,
		},
		Action: &fibcapi.UnicastRoutingFlow_Action{
			Name: fibcapi.UnicastRoutingFlow_Action_DEC_TTL,
		},
		GType: fibcapi.GroupMod_L3_UNICAST,
	}

	mod := fibcapi.FlowMod{
		Cmd:   flowCmd,
		Table: fibcapi.FlowMod_UNICAST_ROUTING,
		ReId:  c.ReID,
		Entry: &fibcapi.FlowMod_Unicast{
			Unicast: &ucast,
		},
	}

	return c.connectVMAPI(func(client fibcapi.FIBCVmApiClient) error {
		ctxt := context.Background()

		w := NewStopWatch()
		w.Start()

		for index := uint32(0); index < num; index++ {
			fflibnet.IncIPNet(ipdst)

			ucast.GId = c.IfBase + (index % c.IfNum)
			ucast.Match.IpDst = ipdst.IP.String()

			if _, err := client.SendFlowMod(ctxt, &mod); err != nil {
				return err
			}

		}

		w.Point()
		w.PrintResults()

		return nil
	})
}

func (c *TestStressCmd) runMpls(cmd string, num uint32) error {
	// TODO
	return nil
}

func (c *TestStressCmd) runPacket(path string, num uint32) error {

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

	if log.IsLevelEnabled(log.DebugLevel) {
		hexdumpDebugLog(data)
	}

	pktin := fibcapi.FFPacketIn{
		DpId: c.VsID,
		Data: data,
	}

	return c.connectVSAPI(func(client fibcapi.FIBCVsApiClient) error {
		ctxt := context.Background()

		w := NewStopWatch()
		w.Start()

		for index := uint32(0); index < num; index++ {
			pktin.PortNo = c.IfBase + (index % c.IfNum)

			if _, err := client.SendPacketIn(ctxt, &pktin); err != nil {
				return err
			}
		}

		w.Point()
		w.PrintResults()

		return nil
	})
}

func testStressCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "stress",
		Short: "Stress test command",
	}

	stress := TestStressCmd{}

	rootCmd.AddCommand(stress.setRunL3Flags(
		&cobra.Command{
			Use:   "setup",
			Short: "setup stress test.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return stress.setup()
			},
		},
	))

	rootCmd.AddCommand(stress.setRunL3Flags(
		&cobra.Command{
			Use:     "teardown",
			Aliases: []string{"td"},
			Short:   "tear down stress test.",
			Args:    cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return stress.teardown()
			},
		},
	))

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "run command",
	}

	rootCmd.AddCommand(runCmd)

	runCmd.AddCommand(stress.setRunL3Flags(
		&cobra.Command{
			Use:   "route <add | del> <num>",
			Short: "run route stress test.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				n, err := strconv.ParseUint(args[1], 0, 32)
				if err != nil || n == 0 {
					return err
				}

				return stress.runRoute(args[0], uint32(n))
			},
		},
	))

	runCmd.AddCommand(stress.setRunL3Flags(
		&cobra.Command{
			Use:   "host <add | del> <num>",
			Short: "run hosr stress test.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				n, err := strconv.ParseUint(args[1], 0, 32)
				if err != nil || n == 0 {
					return err
				}

				return stress.runHost(args[0], uint32(n))
			},
		},
	))

	runCmd.AddCommand(stress.setRunMplsFlags(
		&cobra.Command{
			Use:   "mpls <add | del> <num>",
			Short: "run mpls stress test.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				n, err := strconv.ParseUint(args[1], 0, 32)
				if err != nil || n == 0 {
					return err
				}

				return stress.runMpls(args[0], uint32(n))
			},
		},
	))

	runCmd.AddCommand(stress.setRunPacketFlags(
		&cobra.Command{
			Use:   "packet <path | arp > <num>",
			Short: "run packet stress test.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				n, err := strconv.ParseUint(args[1], 0, 32)
				if err != nil || n == 0 {
					return err
				}

				return stress.runPacket(args[0], uint32(n))
			},
		},
	))

	return rootCmd
}
