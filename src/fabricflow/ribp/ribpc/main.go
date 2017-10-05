// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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
	"fabricflow/ribp/api"
	"flag"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Args struct {
	Addr string
	Name string
}

func getargs() *Args {
	args := &Args{}
	flag.StringVar(&args.Addr, "a", "127.0.0.1:50053", "RIBP API address")
	flag.StringVar(&args.Name, "n", "any", "ifname to send ffpacket.")
	flag.Parse()

	return args
}

func main() {

	args := getargs()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(args.Addr, opts...)
	if err != nil {
		log.Errorf("grpc.Dial error. %v", err)
		return
	}
	defer conn.Close()

	c := ribpapi.NewRIBPApiClient(conn)
	req := &ribpapi.FFPacketRequest{
		Ifname: args.Name,
	}
	reply, err := c.SendFFPacket(context.Background(), req)

	log.Debug("reply: %v, err: %s", reply, err)
}
