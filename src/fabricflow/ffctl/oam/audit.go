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

package oam

import (
	"context"
	"fabricflow/ffctl/fflib"
	fibcapi "fabricflow/fibc/api"

	"github.com/spf13/cobra"
)

type AuditCmd struct {
	fibc *fflib.FibcClient
}

func NewAuditCmd() *AuditCmd {
	return &AuditCmd{
		fibc: fflib.NewFibcClient(),
	}
}

func (c *AuditCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.fibc.Host, "fibc-addr", "", fflib.FibcHost, "fibcd addr.")
	cmd.Flags().Uint16VarP(&c.fibc.Port, "fibc-port", "", fflib.FibcPort, "fibcd port.")
	return cmd
}

func (c *AuditCmd) routeCnt() error {
	return c.fibc.Connect(func(client fibcapi.FIBCApApiClient) error {
		audit := fibcapi.NewOAMAuditRouteCntRequest()
		req := fibcapi.NewOAMRequest(0).SetAuditRouteCnt(audit)

		_, err := client.RunOAM(context.Background(), req)
		return err
	})
}
