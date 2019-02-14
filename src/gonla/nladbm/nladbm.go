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

package nladbm

var (
	clients ClientTable
	nodes   NodeTable
	links   LinkTable
	addrs   AddrTable
	neighs  NeighTable
	routes  RouteTable
	mplss   MplsTable
	vpns    VpnTable
	encaps  EncapInfoTable
	stats   StatTable
)

func Create() {
	clients = NewClientTable()
	nodes = NewNodeTable()
	links = NewLinkTable()
	addrs = NewAddrTable()
	neighs = NewNeighTable()
	routes = NewRouteTable()
	mplss = NewMplsTable()
	vpns = NewVpnTable()
	encaps = NewEncapInfoTable()
	stats = NewStatTable()
}

func Clients() ClientTable {
	return clients
}

func Nodes() NodeTable {
	return nodes
}

func Links() LinkTable {
	return links
}

func Addrs() AddrTable {
	return addrs
}

func Neighs() NeighTable {
	return neighs
}

func Routes() RouteTable {
	return routes
}

func Mplss() MplsTable {
	return mplss
}

func Vpns() VpnTable {
	return vpns
}

func Encaps() EncapInfoTable {
	return encaps
}

func Stats() StatTable {
	return stats
}
