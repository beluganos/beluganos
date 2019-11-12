// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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
	"google.golang.org/grpc"
)

type FIBGrpcController struct {
	addr string

	log *log.Entry
}

func NewFIBGrpcController(addr string) *FIBGrpcController {
	return &FIBGrpcController{
		addr: addr,

		log: log.WithFields(log.Fields{"module": "FIBGrpcController"}),
	}
}

func (c *FIBGrpcController) connect(f func(client fibcapi.FIBCApApiClient) error) error {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	return f(fibcapi.NewFIBCApApiClient(conn))
}

func (c *FIBGrpcController) Dps() ([]uint64, error) {
	dpIds := []uint64{}
	err := c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetDpEntriesRequest{
			Type: fibcapi.DbDpEntry_DPMON,
		}

		stream, err := client.GetDpEntries(context.Background(), &req)
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
			dpId, _ := strconv.ParseUint(e.Id, 0, 64)
			dpIds = append(dpIds, dpId)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dpIds, nil
}

func (c *FIBGrpcController) PortStats(dpId uint64, statsNames []string) (PortStats, error) {
	portStats := NewPortStats()
	err := c.connect(func(client fibcapi.FIBCApApiClient) error {
		req := fibcapi.ApGetPortStatsRequest{
			DpId:   dpId,
			PortNo: 0xffffffff,
			Names:  statsNames,
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

			m := map[string]interface{}{}
			m["port_no"] = stats.PortNo
			for key, val := range stats.Values {
				m[key] = val
			}
			for key, val := range stats.SValues {
				m[key] = val
			}

			portStats = append(portStats, m)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return portStats, nil
}
