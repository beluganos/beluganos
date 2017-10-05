#! /usr/bin/env python
# -*- coding: utf-8 -*-

# Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import grpc
from nlaapi import nlaapi_pb2 as api

def dump_nodes(stub):
    print("# Node")

    cnt = 0
    for node in stub.GetNodes(api.GetNodesRequest()):
        cnt += 1
        print(node)

    print("# {0} nodes".format(cnt))

def dump_links(stub):
    print("# Link")

    cnt = 0
    for link in stub.GetLinks(api.GetLinksRequest()):
        cnt += 1
        print(link)

    print("# {0} links".format(cnt))

def dump_addrs(stub):
    print("# Addr")

    cnt = 0
    for addr in stub.GetAddrs(api.GetAddrsRequest()):
        cnt += 1
        print(addr)

    print("# {0} addrs".format(cnt))

def dump_neighs(stub):
    print("# Neigh")

    cnt = 0
    for neigh in stub.GetNeighs(api.GetNeighsRequest()):
        cnt += 1
        print(neigh)

    print("# {0} neighs".format(cnt))

def dump_routes(stub):
    print("# Route")

    cnt = 0
    for route in stub.GetRoutes(api.GetRoutesRequest()):
        cnt += 1
        print(route)

    print("# {0} routes".format(cnt))

def dump_mplss(stub):
    print("# MPLS")

    cnt = 0
    for mpls in stub.GetMplss(api.GetMplssRequest()):
        cnt += 1
        print(mpls)

    print("# {0} mplss".format(cnt))

def dump_vpns(stub):
    print("# VPN")

    cnt = 0
    for vpn in stub.GetVpns(api.GetVpnsRequest()):
        cnt += 1
        print(vpn)

    print("# {0} vpns".format(cnt))

def mon_netlink(stub):
    for nlmsg in stub.MonNetlink(api.MonNetlinkRequest()):
        print(nlmsg)

def dump_all(stub):
    dump_nodes(stub)
    dump_links(stub)
    dump_addrs(stub)
    dump_neighs(stub)
    dump_routes(stub)
    dump_mplss(stub)
    dump_vpns(stub)

dump_cmd = dict(
    node = dump_nodes,
    link = dump_links,
    addr = dump_addrs,
    neigh= dump_neighs,
    route= dump_routes,
    mpls = dump_mplss,
    vpn  = dump_vpns,
    all  = dump_all,
)


def _getopts():
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--addr", default="127.0.0.1:50052")
    parser.add_argument("cmd")
    parser.add_argument("table", nargs="?", default="all")
    return parser.parse_args(), parser


def _help():
    _, p = _getopts
    p.print_help()

def _main():
    opts, _ = _getopts()

    channel = grpc.insecure_channel(opts.addr)
    stub = api.NLAApiStub(channel)

    if opts.cmd == "dump":
        if opts.table not in dump_cmd:
            print("{0} not found".format(opts.table))
            return

        dump_cmd[opts.table](stub)

    elif opts.cmd == "mon":
        mon_netlink(stub)

    else:
        _help()

if __name__ == "__main__":
    _main()
        
