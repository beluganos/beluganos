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

"""
FIBC mod factory
"""

from goryu.ofproto import ofdpa_match
from goryu.ofproto import ofdpa_action
from fabricflow.fibc.ofc import generic

ofdpa_match.init()
ofdpa_action.init()

def setup_flow(dpath, mod, ofctl):
    """
    Setup flows
    """
    return generic.setup_flow(dpath, mod, ofctl, False)


def vlan_flow(dpath, mod, ofctl):
    """
    VLAN flow table.
    """
    return generic.vlan_flow(dpath, mod, ofctl, False)


def termination_mac_flow(dpath, mod, ofctl):
    """
    Termination MAC flow table.
    """
    return generic.termination_mac_flow(dpath, mod, ofctl)


def mpls1_flow(dpath, mod, ofctl):
    """
    MPLS1 flow table.
    """
    return generic.mpls1_flow(dpath, mod, ofctl, False)


def unicast_routing_flow(dpath, mod, ofctl):
    """
    Create flow_mod for Unicast Routing flow table.x
    """
    return generic.unicast_routing_flow(dpath, mod, ofctl, False)


def bridging_flow(dpath, mod, ofctl):
    """
    Bridging flow table.
    """
    return generic.bridging_flow(dpath, mod, ofctl)


def policy_acl_flow(dpath, mod, ofctl):
    """
    Policy ACL flow table.
    """
    return generic.policy_acl_flow(dpath, mod, ofctl, False)


def setup_group(dpath, mod, ofctl):
    """
    Setup Group.
    """
    return generic.setup_group(dpath, mod, ofctl)


def l2_interface_group(dpath, mod, ofctl):
    """
    L2 Interface Group
    """
    return generic.l2_interface_group(dpath, mod, ofctl)


def l3_unicast_group(dpath, mod, ofctl):
    """
    L3 Unicast Group
    """
    return generic.l3_unicast_group(dpath, mod, ofctl)


def l3_ecmp_group(dpath, mod, ofctl):
    """
    ECMP Group
    """
    return generic.l3_ecmp_group(dpath, mod, ofctl)


def mpls_interface_group(dpath, mod, ofctl):
    """
    MPLS Interface group
    """
    return generic.mpls_interface_group(dpath, mod, ofctl)


def mpls_l3_vpn_group(dpath, mod, ofctl):
    """
    MPLS L3 VPN Group
    """
    return generic.mpls_l3_vpn_group(dpath, mod, ofctl)


def mpls_tun1_group(dpath, mod, ofctl):
    """
    MPLS Tunnel1 Label Group
    """
    return generic.mpls_tun1_group(dpath, mod, ofctl)


def mpls_swap_group(dpath, mod, ofctl):
    """
    MPLS Swap Label Group
    """
    return generic.mpls_swap_group(dpath, mod, ofctl)


def mpls_ecmp_group(dpath, mod, ofctl):
    """
    MPLS ECMP Group
    """
    return generic.mpls_ecmp_group(dpath, mod, ofctl)


def l2_unfiltered_interface_group(dpath, mod, ofctl):
    """
    L2 Unfiltered Interface Group.
    """
    return generic.l2_unfiltered_interface_group(dpath, mod, ofctl)


def pkt_out(dpath, port_id, strip_vlan, data):
    """
    PacketOut
    """
    return generic.pkt_out(dpath, port_id, strip_vlan, data)


def get_port_stats(dpath, waiters, port_id, ofctl):
    """
    get port stats
    """
    return generic.get_port_stats(dpath, waiters, port_id, ofctl)


def port_mod(dpath, mod, ofctl):
    """
    PotMod
    """
    return generic.port_mod(dpath, mod, ofctl)
