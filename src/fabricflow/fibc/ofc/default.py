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

import logging
from goryu.ofproto import offlow
from goryu.ofproto import ofgroup
from goryu.ofproto import ofaction

_LOG = logging.getLogger(__name__)

def setup_flow(dpath, mod, ofctl):
    """
    Setup flows.
    """
    _LOG.debug("Default FLow: %d %s", dpath.id, mod)

    ofp = dpath.ofproto

    # send whole packet to controller
    dpath.send_msg(offlow.set_sw_config(dpath))

    # clear all flows/groups.
    dpath.send_msg(offlow.clear_all(dpath))
    for msg in ofgroup.clear_all(dpath):
        dpath.send_msg(msg)

    # send all packet to controller
    flow = offlow.flow_mod(
        table_id=0,
        priority=0,
        match={},
        actions=[
            ofaction.output(ofp.OFPP_CONTROLLER),
        ],
        writes=[],
    )
    ofctl.mod_flow_entry(dpath, flow, ofp.OFPFC_ADD)


def vlan_flow(dpath, mod, ofctl):
    """
    VLAN flow table.
    """
    _LOG.debug("VLAN FLow: %d %s %s", dpath.id, mod, ofctl)


def termination_mac_flow(dpath, mod, ofctl):
    """
    Termination MAC flow table.
    """
    _LOG.debug("TERM MAC FLow: %d %s %s", dpath.id, mod, ofctl)


def mpls1_flow(dpath, mod, ofctl):
    """
    MPLS1 flow table.
    """
    _LOG.debug("MPLS1 FLow: %d %s %s", dpath.id, mod, ofctl)


def unicast_routing_flow(dpath, mod, ofctl):
    """
    Create flow_mod for Unicast Routing flow table.x
    """
    _LOG.debug("Unicast Routing FLow: %d %s %s", dpath.id, mod, ofctl)


def bridging_flow(dpath, mod, ofctl):
    """
    Bridging flow table.
    """
    _LOG.debug("Bridging FLow: %d %s %s", dpath.id, mod, ofctl)


def policy_acl_flow(dpath, mod, ofctl):
    """
    Policy ACL flow table.
    """
    _LOG.debug("ACL FLow: %d %s %s", dpath.id, mod, ofctl)


def setup_group(dpath, mod, ofctl):
    """
    Default Group.
    """
    _LOG.debug("Setup Group: %d %s %s", dpath.id, mod, ofctl)


def l2_interface_group(dpath, mod, ofctl):
    """
    L2 Interface Group
    """
    _LOG.debug("L2 Interface Group: %d %s %s", dpath.id, mod, ofctl)


def l3_unicast_group(dpath, mod, ofctl):
    """
    L3 Unicast Group
    """
    _LOG.debug("L3 Unicast Group: %d %s %s", dpath.id, mod, ofctl)


def l3_ecmp_group(dpath, mod, ofctl):
    """
    ECMP Group
    """
    _LOG.debug("L3 ECMP Group: %d %s %s", dpath.id, mod, ofctl)


def mpls_interface_group(dpath, mod, ofctl):
    """
    MPLS Interface group
    """
    _LOG.debug("MPLS Interface Group: %d %s %s", dpath.id, mod, ofctl)


def mpls_l3_vpn_group(dpath, mod, ofctl):
    """
    MPLS L3 VPN Group
    """
    _LOG.debug("MPLS L3 CPN Group: %d %s %s", dpath.id, mod, ofctl)


def mpls_tun1_group(dpath, mod, ofctl):
    """
    MPLS Tunnel1 Label Group
    """
    _LOG.debug("MPLS Tunnel1 Group: %d %s %s", dpath.id, mod, ofctl)


def mpls_swap_group(dpath, mod, ofctl):
    """
    MPLS Swap Label Group
    """
    _LOG.debug("MPLS Swap Group: %d %s %s", dpath.id, mod, ofctl)


def mpls_ecmp_group(dpath, mod, ofctl):
    """
    MPLS ECMP Group
    """
    _LOG.debug("MPLS ECMP Group: %d %s %s", dpath.id, mod, ofctl)


def l2_unfiltered_interface_group(dpath, mod, ofctl):
    """
    L2 Unfiltered Interface Group.
    """
    _LOG.debug("L2 Unfiltered Interface Group: %d %s %s", dpath.id, mod, ofctl)
