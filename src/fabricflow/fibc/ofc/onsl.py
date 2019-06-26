# -*- coding: utf-8 -*-

# Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

"""
FIBC mod factory
"""

import logging
from ryu.lib import hub
from fabricflow.fibc.api import fibcapi
from fabricflow.fibc.api import fibcapi_pb2 as pb

_LOG = logging.getLogger("fibc.ofc.onsl")

 # pylint: disable=unused-argument

def setup_flow(dpath, mod, ofctl):
    """
    Setup flows.
    """
    _LOG.debug("Default FLow: %d %s", dpath.id, mod)

    """
    matches = [
        fibcapi.new_policy_acl_match(eth_type=fibcapi.ETHTYPE_LACP),
        fibcapi.new_policy_acl_match(eth_type=fibcapi.ETHTYPE_ARP),
        fibcapi.new_policy_acl_match(eth_type=0x0800, ip_dst=fibcapi.MCADDR_ALLROUTERS),
        fibcapi.new_policy_acl_match(eth_type=0x0800, ip_dst=fibcapi.MCADDR_OSPF_HELLO),
        fibcapi.new_policy_acl_match(eth_type=0x0800, ip_dst=fibcapi.MCADDR_OSPF_ALLDR),

        fibcapi.new_policy_acl_match(eth_type=0x86dd, ip_dst=fibcapi.MCADDR6_I_LOCAL),
        fibcapi.new_policy_acl_match(eth_type=0x86dd, ip_dst=fibcapi.MCADDR6_L_LOCAL),
        fibcapi.new_policy_acl_match(eth_type=0x86dd, ip_dst=fibcapi.MCADDR6_S_LOCAL),
        fibcapi.new_policy_acl_match(eth_type=0x86dd, ip_dst=fibcapi.UCADDR6_L_LOCAL),

        fibcapi.new_policy_acl_match(eth_dst=fibcapi.HWADDR_ISIS_LEVEL1),
        fibcapi.new_policy_acl_match(eth_dst=fibcapi.HWADDR_ISIS_LEVEL2),
    ]
    action = fibcapi.new_policy_acl_action("OUTPUT")

    for match in matches:
        acl = pb.PolicyACLFlow(match=match, action=action)
        mod = pb.FlowMod(
            cmd="ADD",
            table="POLICY_ACL",
            re_id="",
            acl=acl)

        dpath.send_msg(pb.FLOW_MOD, mod)
    """


def vlan_flow(dpath, mod, ofctl):
    """
    VLAN flow table.
    """
    _LOG.debug("VLAN FLow: %d %s", dpath.id, mod)

    dpath.send_msg(pb.FLOW_MOD, mod)


def termination_mac_flow(dpath, mod, ofctl):
    """
    Termination MAC flow table.
    """
    _LOG.debug("TERM MAC FLow: %d %s", dpath.id, mod)

    dpath.send_msg(pb.FLOW_MOD, mod)


def mpls1_flow(dpath, mod, ofctl):
    """
    MPLS1 flow table.
    """
    _LOG.debug("MPLS1 FLow: %d %s", dpath.id, mod)

    dpath.send_msg(pb.FLOW_MOD, mod)

def unicast_routing_flow(dpath, mod, ofctl):
    """
    Create flow_mod for Unicast Routing flow table.
    """
    _LOG.debug("Unicast Routing FLow: %d %s", dpath.id, mod)

    dpath.send_msg(pb.FLOW_MOD, mod)


def bridging_flow(dpath, mod, ofctl):
    """
    Bridging flow table.
    """
    _LOG.debug("Bridging FLow: %d %s", dpath.id, mod)

    dpath.send_msg(pb.FLOW_MOD, mod)


def policy_acl_flow(dpath, mod, ofctl):
    """
    Policy ACL flow table.
    """
    _LOG.debug("ACL FLow: %d %s", dpath.id, mod)

    if mod.acl.match.in_port == 0:
        # ignore flow mod that match.in_port not specified.
        return

    dpath.send_msg(pb.FLOW_MOD, mod)


def setup_group(dpath, mod, ofctl):
    """
    Setup Group.
    """
    _LOG.debug("Default Group: %d %s", dpath.id, mod)


def l2_interface_group(dpath, mod, ofctl):
    """
    L2 Interface Group.
    """
    _LOG.debug("L2 Interface Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def l3_unicast_group(dpath, mod, ofctl):
    """
    L3 Unicast Group.
    """
    _LOG.debug("L3 Unicast Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def l3_ecmp_group(dpath, mod, ofctl):
    """
    ECMP Group.
    """
    _LOG.debug("L3 ECMP Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def mpls_interface_group(dpath, mod, ofctl):
    """
    MPLS Interface group.
    """
    _LOG.debug("MPLS Interface Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def mpls_l3_vpn_group(dpath, mod, ofctl):
    """
    MPLS L3 VPN Group.
    """
    _LOG.debug("MPLS L3 VPN Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def mpls_tun1_group(dpath, mod, ofctl):
    """
    MPLS Tunnel1 Label Group
    """
    _LOG.debug("MPLS Tunnel1 Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def mpls_swap_group(dpath, mod, ofctl):
    """
    MPLS Swap Label Group.
    """
    _LOG.debug("MPLS Swap Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def mpls_ecmp_group(dpath, mod, ofctl):
    """
    MPLS ECMP Group
    """
    _LOG.debug("MPLS ECMP Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def l2_unfiltered_interface_group(dpath, mod, ofctl):
    """
    L2 Unfiltered Interface Group.
    """
    _LOG.debug("L2 Unfiltered Interface Group: %d %s", dpath.id, mod)

    dpath.send_msg(pb.GROUP_MOD, mod)


def pkt_out(dpath, port_id, strip_vlan, data):
    """
    PacketOut
    """
    _LOG.debug("PacketOUT: %s %d %d", dpath.id, port_id, len(data))

    msg = fibcapi.new_ff_packet_out(dpath.id, port_id, data)
    dpath.send_msg(pb.FF_PACKET_OUT, msg)


_PORT_STATS_NAMES = (
    "ifInOctets",
    "ifInUcastPkts",
    "ifInNUcastPkts",
    "ifInDiscards",
    "ifInErrors",
    "ifOutOctets",
    "ifOutUcastPkts",
    "ifOutNUcastPkts",
    "ifOutDiscards",
    "ifOutErrors",
)

def get_port_stats(dpath, waiters, port_id, ofctl):
    """
    get port stats
    """
    _LOG.debug("get_port_stats: %d %s", dpath.id, port_id)

    xid = dpath.get_xid()
    lock = hub.Event()
    msgs = list()
    prev_msg_num = len(msgs)
    if port_id is None:
        port_id = dpath.ofproto.OFPP_ANY

    dp_waiter = waiters.setdefault(dpath.id, {})
    dp_waiter[xid] = (lock, msgs)

    msg = fibcapi.new_ff_multipart_request_port(dpath.id, port_id, _PORT_STATS_NAMES)
    dpath.send_msg(pb.FF_MULTIPART_REQUEST, msg, xid)

    while True:
        lock.wait(timeout=0.5)

        curr_msg_num = len(msgs)
        if curr_msg_num == prev_msg_num:
            break

        prev_msg_num = curr_msg_num

    if not lock.is_set():
        # Timeout
        del dp_waiter[xid]
        if not dp_waiter:
            del waiters[dpath.id]
        return None

    stats_list = list()
    for msg in msgs:
        for stats in msg.port.stats:
            stats_entry = {k:v for k, v in stats.values.items()}
            stats_entry["port_no"] = stats.port_no
            stats_list.append(stats_entry)

    return {dpath.id: stats_list}


def port_mod(dpath, mod, ofctl):
    """
    PotMod
    """
    _LOG.debug("port_mod: %s %s", dpath, mod)

    dpath.send_msg(pb.FF_PORT_MOD, mod)
