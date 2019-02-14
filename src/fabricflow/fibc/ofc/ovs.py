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

from fabricflow.fibc.ofc import generic

# pylint: disable=invalid-name

def _mpls_l3_vpn_group(dpath, mod, ofctl):
    """
    MPLS L3 VPN Group
    """
    return generic.mpls_l3_vpn_group(dpath, mod, ofctl, False)

setup_flow = generic.setup_flow
vlan_flow = generic.vlan_flow
termination_mac_flow = generic.termination_mac_flow
mpls1_flow = generic.mpls1_flow
unicast_routing_flow = generic.unicast_routing_flow
bridging_flow = generic.bridging_flow
policy_acl_flow = generic.policy_acl_flow
setup_group = generic.setup_group
l2_interface_group = generic.l2_interface_group
l3_unicast_group = generic.l3_unicast_group
l3_ecmp_group = generic.l3_ecmp_group
mpls_interface_group = generic.mpls_interface_group
mpls_l3_vpn_group = _mpls_l3_vpn_group
mpls_tun1_group = generic.mpls_tun1_group
mpls_swap_group = generic.mpls_swap_group
mpls_ecmp_group = generic.mpls_ecmp_group
l2_unfiltered_interface_group = generic.l2_unfiltered_interface_group
pkt_out = generic.pkt_out
get_port_stats = generic.get_port_stats
