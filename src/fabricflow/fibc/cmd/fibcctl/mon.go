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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//
// MonitorCmd is monitor command.
//
type MonitorCmd struct {
	Addr string
	ReID string
	DpID uint64
	VsID uint64

	noAPAPI bool
	noDPAPI bool
	noVSAPI bool
	noVMAPI bool
}

func (c *MonitorCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.Addr, "fibc-addr", "", "localhost:50061", "fibc address.")
	cmd.Flags().Uint64VarP(&c.DpID, "dpid", "", 0, "daapath id.")
	cmd.Flags().StringVarP(&c.ReID, "reid", "", "", "router entity id")
	cmd.Flags().Uint64VarP(&c.VsID, "vsid", "", 0, "vswitch id.")

	cmd.Flags().BoolVarP(&c.noAPAPI, "no-ap", "", false, "disable ap-api monitor.")
	cmd.Flags().BoolVarP(&c.noDPAPI, "no-dp", "", false, "disable dp-api monitor.")
	cmd.Flags().BoolVarP(&c.noVSAPI, "no-vs", "", false, "disable vs-api monitor.")
	cmd.Flags().BoolVarP(&c.noVMAPI, "no-vm", "", false, "disable vm-api monitor.")

	return cmd
}

func (c *MonitorCmd) monitor(done chan struct{}) error {
	log.Infof("monitor %s", c.Addr)

	if !c.noAPAPI {
		log.Infof("AP API monitor enabled.")

		apcmd := APAPICommand{
			Addr: c.Addr,
		}

		go func() {
			if err := apcmd.monitor(done); err != nil {
				log.Errorf("AP API monitor error. %s", err)
				close(done)
			}
		}()
	}

	if !c.noDPAPI {
		log.Infof("DP API monitor enabled. dp-id:%d", c.DpID)

		dpcmd := DPAPICommand{
			Addr: c.Addr,
			DpID: c.DpID,
		}

		go func() {
			if err := dpcmd.monitor(done); err != nil {
				log.Errorf("DP API monitor error. %s", err)
				close(done)
			}
		}()
	}

	if !c.noVMAPI {
		log.Infof("VM API monitor enabled. re-id:%s", c.ReID)

		vmcmd := VMAPICommand{
			Addr: c.Addr,
			ReID: c.ReID,
		}

		go func() {
			if err := vmcmd.monitor(done); err != nil {
				log.Errorf("VM API monitor error. %s", err)
				close(done)
			}
		}()
	}

	if !c.noVSAPI {
		log.Infof("VS API monitor enabled. vs-id:%d", c.VsID)

		vscmd := VSAPICommand{
			Addr: c.Addr,
			VsID: c.VsID,
		}

		go func() {
			if err := vscmd.monitor(done); err != nil {
				log.Errorf("VS API monitor error. %s", err)
				close(done)
			}
		}()
	}

	<-done

	return nil
}

func monCmd() *cobra.Command {
	mon := MonitorCmd{}
	rootCmd := &cobra.Command{
		Use:   "mon <name, ...>",
		Short: "monitor fibc messages",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			done := make(chan struct{})
			return mon.monitor(done)
		},
	}

	return mon.setFlags(rootCmd)
}
