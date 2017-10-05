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
	"fabricflow/ribs/ribsapi"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func dumpRics(c ribsapi.RIBSApiClient) {
	stream, err := c.GetRics(context.Background(), &ribsapi.GetRicsRequest{})
	if err != nil {
		log.Printf("GetRics error. %s", err)
		return
	}
	for {
		ric, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("RIC: %v", ric)
	}
}

func dumpNexthops(c ribsapi.RIBSApiClient) {
	stream, err := c.GetNexthops(context.Background(), &ribsapi.GetNexthopsRequest{})
	if err != nil {
		log.Printf("GetNexthops error. %s", err)
		return
	}

	for {
		nh, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("NH : %v", nh)
	}
}

func dumpNexthopMap(c ribsapi.RIBSApiClient) {
	stream, err := c.GetNexthopMap(context.Background(), &ribsapi.GetNexthopMapRequest{})
	if err != nil {
		log.Printf("GetNexthopMap error. %s", err)
		return
	}

	for {
		nh, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("MAP: %v", nh)
	}
}

func main() {

	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:50072", "RIBS API address.")
	flag.Parse()

	log.SetFlags(0)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Printf("grpc.Dial error. %v", err)
		return
	}
	defer conn.Close()

	c := ribsapi.NewRIBSApiClient(conn)

	dumpRics(c)
	dumpNexthops(c)
	dumpNexthopMap(c)
}
