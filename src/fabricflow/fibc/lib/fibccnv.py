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
Convert Flow/Group mod.
"""

from fabricflow.fibc.api import fibcapi_pb2 as pb


def get_dp_port(portmap, re_id, port_id):
    """
    Get port_id from portmap.
    """
    if port_id == 0:
        return 0

    port = portmap.find_by_vm(re_id=re_id, port_id=port_id)
    port = portmap.lower_port(port)
    return port["dp"].port


def conv_vlan_flow(mod, portmap):
    """
    Convert mod for VLAN Flow table.
    """
    mod.vlan.match.in_port = get_dp_port(
        portmap, mod.re_id, mod.vlan.match.in_port)


def conv_termmac_flow(mod, portmap):
    """
    Convert mod for VLAN Flow table.
    """
    mod.term_mac.match.in_port = get_dp_port(
        portmap, mod.re_id, mod.term_mac.match.in_port)


def conv_bridging_flow(mod, portmap):
    """
    Convert mod for Bridging flow table.
    """
    action = mod.bridging.action
    if action.name == pb.PolicyACLFlow.Action.OUTPUT and action.value != 0:
        mod.bridging.action.value = get_dp_port(
            portmap, mod.re_id, action.value)


def conv_policy_acl_flow(mod, portmap):
    """
    Convert mod for PolicyACL flow table.
    """
    mod.acl.match.in_port = get_dp_port(
        portmap, mod.re_id, mod.acl.match.in_port)


_FLOW_MAP = {
    pb.FlowMod.VLAN: conv_vlan_flow,
    pb.FlowMod.TERM_MAC: conv_termmac_flow,
    pb.FlowMod.BRIDGING: conv_bridging_flow,
    pb.FlowMod.POLICY_ACL: conv_policy_acl_flow,
}

def conv_flow(mod, portmap):
    """
    Convert flow mod.
    """
    func = _FLOW_MAP.get(mod.table, None)
    if func is not None:
        func(mod, portmap)


def conv_l2_interface_group(mod, portmap):
    """
    Conver Group mod for L2 Interface Group.
    """
    mod.l2_iface.port_id = get_dp_port(portmap, mod.re_id, mod.l2_iface.port_id)
    mod.l2_iface.master = get_dp_port(portmap, mod.re_id, mod.l2_iface.master)


def conv_l3_unicast_group(mod, portmap):
    """
    Conver Group mod for L3 Unicast Group.
    """
    mod.l3_unicast.port_id = get_dp_port(portmap, mod.re_id, mod.l3_unicast.port_id)
    mod.l3_unicast.phy_port_id = get_dp_port(portmap, mod.re_id, mod.l3_unicast.phy_port_id)


def conv_mpls_interface_group(mod, portmap):
    """
    Conver Group mod for MPLS Interface Group.
    """
    mod.mpls_iface.port_id = get_dp_port(portmap, mod.re_id, mod.mpls_iface.port_id)


_GROUP_MAP = {
    pb.GroupMod.L2_INTERFACE: conv_l2_interface_group,
    pb.GroupMod.L3_UNICAST: conv_l3_unicast_group,
    pb.GroupMod.MPLS_INTERFACE: conv_mpls_interface_group,
}

def conv_group(mod, portmap):
    """
    Convert Group mod.
    """
    func = _GROUP_MAP.get(mod.g_type, None)
    if func is not None:
        func(mod, portmap)
