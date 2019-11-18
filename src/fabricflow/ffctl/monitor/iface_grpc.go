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

package monitor

import (
	"context"
	fibcapi "fabricflow/fibc/api"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
)

type IfaceStatsMessage struct {
	Time        time.Time
	Message     string
	MessageOnly bool
	PortStats   []*fibcapi.FFPortStats
}

func (m *IfaceStatsMessage) String() string {
	return fmt.Sprintf("%s %s", m.Time.Format("2006/01/02 15:04:05"), m.Message)
}

func NewIfaceStatsMessage() *IfaceStatsMessage {
	return &IfaceStatsMessage{
		Time:      time.Now(),
		PortStats: []*fibcapi.FFPortStats{},
	}
}

func (c *IfaceController) resetPortStats() error {
	if c.fibcClient == nil {
		return nil
	}

	req := fibcapi.ApModPortStatsRequest{
		DpId:   c.DpathCfg.DpID,
		PortNo: uint32(0xffffffff),
		Cmd:    fibcapi.FFPortStats_RESET,
	}

	msg := NewIfaceStatsMessage()
	msg.MessageOnly = true
	if _, err := c.fibcClient.ModPortStats(context.Background(), &req); err != nil {
		msg.Message = fmt.Sprintf("ERROR: %s", err)
	} else {
		msg.Message = "Reset stats counters."
	}

	c.statsCh <- msg

	return nil
}

func (c *IfaceController) startPortStats() error {
	fibcAddr := fmt.Sprintf("%s:%d", c.FibcAddr, c.FibcPort)
	conn, err := grpc.Dial(fibcAddr, grpc.WithInsecure())
	if err != nil {
		c.log.Errorf("servePortStats: Dial error. %s", err)
		return err
	}

	c.fibcClient = fibcapi.NewFIBCApApiClient(conn)

	c.log.Debugf("servePortStats: START. %s %s", fibcAddr, c.Interval)

	go func() {
		defer conn.Close()

		tick := time.NewTicker(c.Interval)
		defer tick.Stop()

	FOR_LOOP:
		for {
			select {
			case <-tick.C:
				portStats, err := c.getPortStats()

				msg := NewIfaceStatsMessage()
				if err != nil {
					msg.Message = fmt.Sprintf("ERROR: %s", err)
				} else {
					msg.Message = "Updated."
					msg.PortStats = portStats
				}

				c.statsCh <- msg

			case <-c.done:
				c.log.Debugf("servePortStats: EXIT.")
				break FOR_LOOP
			}
		}
	}()

	return nil
}

func (c *IfaceController) getPortStats() ([]*fibcapi.FFPortStats, error) {
	req := fibcapi.ApGetPortStatsRequest{
		DpId:   c.DpathCfg.DpID,
		PortNo: uint32(0xffffffff),
	}

	portStats := []*fibcapi.FFPortStats{}

	stream, err := c.fibcClient.GetPortStats(context.Background(), &req)
	if err != nil {
		c.log.Debugf("getPortStats: getPortStats error.%s", err)
		return portStats, err
	}

FOR_LOOP:
	for {
		stats, err := stream.Recv()
		if err == io.EOF {
			break FOR_LOOP
		}
		if err != nil {
			c.log.Debugf("getPortStats: recv error. %s", err)

			return portStats, err
		}
		if stats == nil {
			c.log.Debugf("getPortStats: invalid msg.")
			continue FOR_LOOP
		}

		portStats = append(portStats, stats)

		c.log.Debugf("getPortStats: port:%d", stats.PortNo)
	}

	return portStats, nil
}
