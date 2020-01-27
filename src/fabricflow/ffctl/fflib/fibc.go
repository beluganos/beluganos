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

package fflib

import (
	fibcapi "fabricflow/fibc/api"
	"fmt"

	"google.golang.org/grpc"
)

const (
	FibcHost = "localhost"
	FibcPort = uint16(50070)
)

type FibcClient struct {
	Host string
	Port uint16
}

func NewFibcClient() *FibcClient {
	return &FibcClient{
		Host: FibcHost,
		Port: FibcPort,
	}
}

func (c *FibcClient) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *FibcClient) Connect(f func(fibcapi.FIBCApApiClient) error) error {
	conn, err := grpc.Dial(c.Addr(), grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	return f(fibcapi.NewFIBCApApiClient(conn))
}
