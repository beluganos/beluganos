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
fabricflow.fibc.ofc init
"""

from fabricflow.fibc.ofc import generic
from fabricflow.fibc.ofc import ovs
from fabricflow.fibc.ofc import ofdpa2
from fabricflow.fibc.ofc import default
from fabricflow.fibc.api import fibcapi_pb2 as pb

def _flow_dict(mod):
    return {
        -1                        : mod.setup_flow,
        pb.FlowMod.VLAN           : mod.vlan_flow,
        pb.FlowMod.TERM_MAC       : mod.termination_mac_flow,
        pb.FlowMod.MPLS1          : mod.mpls1_flow,
        pb.FlowMod.UNICAST_ROUTING: mod.unicast_routing_flow,
        pb.FlowMod.BRIDGING       : mod.bridging_flow,
        pb.FlowMod.POLICY_ACL     : mod.policy_acl_flow,
    }

def _group_dict(mod):
    return {
        -1                         : mod.setup_group,
        pb.GroupMod.L2_INTERFACE   : mod.l2_interface_group,
        pb.GroupMod.L3_UNICAST     : mod.l3_unicast_group,
        pb.GroupMod.L3_ECMP        : mod.l3_ecmp_group,
        pb.GroupMod.MPLS_INTERFACE : mod.mpls_interface_group,
        pb.GroupMod.MPLS_L3_VPN    : mod.mpls_l3_vpn_group,
        pb.GroupMod.MPLS_TUNNEL1   : mod.mpls_tun1_group,
        pb.GroupMod.MPLS_SWAP      : mod.mpls_swap_group,
        pb.GroupMod.MPLS_ECMP      : mod.mpls_ecmp_group,
        pb.GroupMod.L2_UF_INTERFACE: mod.l2_unfiltered_interface_group,
    }


def _mod_entry(mod):
    return (_flow_dict(mod), _group_dict(mod))


def _get_mod(mode, index, table):
    if mode not in _MODS:
        mode = "default"

    return _MODS[mode][index][table]


_MODS = {
    "default": _mod_entry(default),
    "generic": _mod_entry(generic),
    "ofdpa2" : _mod_entry(ofdpa2),
    "ovs"    : _mod_entry(ovs),
}


def flow(mode, table):
    """
    Get Generator of flow table.
    """
    return _get_mod(mode, 0, table)


def group(mode, g_type):
    """
    Get Generator of group talbe.
    """
    return _get_mod(mode, 1, g_type)
