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

from goryu.ofproto import ofmatch
from goryu.ofproto import ofaction
from goryu.ofproto import offlow
from fabricflow.fibc.api import fibcapi
from fabricflow.fibc.api import fibcapi_pb2 as pb

def _add_table_miss_flows(dpath, ofctl):
    """
    Table miss flows.
    """
    tables = [
        (pb.FlowMod.INGRESS_PORT, pb.FlowMod.VLAN, []),
        (pb.FlowMod.VLAN, pb.FlowMod.POLICY_ACL, [ofaction.clear_actions()]),
        (pb.FlowMod.TERM_MAC, pb.FlowMod.BRIDGING, []),
        (pb.FlowMod.MPLS0, pb.FlowMod.MPLS1, []),
        (pb.FlowMod.MPLS1, -1, [ofaction.clear_actions()]),
        (pb.FlowMod.MPLS2, -1, [ofaction.clear_actions()]),
        (pb.FlowMod.MPLS_L3_TYPE, -1, [ofaction.clear_actions()]),
        (pb.FlowMod.MPLS_LABEL_TRUST, pb.FlowMod.MPLS_TYPE, []),
        (pb.FlowMod.MPLS_TYPE, pb.FlowMod.POLICY_ACL, []),
        (pb.FlowMod.UNICAST_ROUTING, pb.FlowMod.POLICY_ACL, []),
        (pb.FlowMod.MULTICAST_ROUTING, pb.FlowMod.POLICY_ACL, []),
        (pb.FlowMod.BRIDGING, pb.FlowMod.POLICY_ACL, []),
        (pb.FlowMod.POLICY_ACL, -1, []),
        # (pb.FlowMod.POLICY_ACL, -1, [ofaction.output(dpath.ofproto.OFPP_CONTROLLER)]),
    ]
    for talbe_id, goto_table, actions in tables:
        if goto_table > 0:
            actions.append(ofaction.goto_table(goto_table))

        flow = offlow.flow_mod(
            match={},
            actions=actions,
            writes=[],
            table_id=talbe_id,
            priority=fibcapi.PRIORITY_DEFAULT,
        )
        ofctl.mod_flow_entry(dpath, flow, dpath.ofproto.OFPFC_ADD)

    return


def _add_mpls_l3_type_l3vpn_flows(dpath, ofctl):
    # L3 VPN Route (IPv4 Unicast)
    flow = offlow.flow_mod(
        match=ofmatch.Match().eth_type(fibcapi.ETHTYPE_MPLS),
        actions=[
            ofaction.set_mpls_type(fibcapi.MPLSTYPE_UNICAST, True),
            ofaction.pop_mpls(fibcapi.ETHTYPE_IPV4),
            ofaction.goto_table(pb.FlowMod.MPLS_LABEL_TRUST),
        ],
        writes=[],
        table_id=pb.FlowMod.MPLS_L3_TYPE,
        priority=1,
    )
    ofctl.mod_flow_entry(dpath, flow, dpath.ofproto.OFPFC_ADD)
    return


def _add_mpls_l3_type_php_flows(dpath, ofctl):
    # L3 VPN Forward (IPv4) based on this label (PHP)
    flow = offlow.flow_mod(
        match=ofmatch.Match().eth_type(fibcapi.ETHTYPE_MPLS).mpls_type(fibcapi.MPLSTYPE_PHP, True),
        actions=[
            ofaction.pop_mpls(fibcapi.ETHTYPE_IPV4),
            ofaction.goto_table(pb.FlowMod.MPLS_LABEL_TRUST),
        ],
        writes=[],
        table_id=pb.FlowMod.MPLS_L3_TYPE,
        priority=5,
    )
    ofctl.mod_flow_entry(dpath, flow, dpath.ofproto.OFPFC_ADD)
    return


def _add_mpls_l3_type_flows(dpath, ofctl):
    _add_mpls_l3_type_l3vpn_flows(dpath, ofctl)
    _add_mpls_l3_type_php_flows(dpath, ofctl)
    return


def _add_mpls_type_flows(dpath, ofctl):
    """
    MPLS Type builtin flows.
    """
    datas = [
        (fibcapi.MPLSTYPE_VPS, pb.FlowMod.POLICY_ACL),
        (fibcapi.MPLSTYPE_UNICAST, pb.FlowMod.UNICAST_ROUTING),
        (fibcapi.MPLSTYPE_MULTICAST, pb.FlowMod.MULTICAST_ROUTING),
        (fibcapi.MPLSTYPE_PHP, pb.FlowMod.POLICY_ACL),
    ]
    for mpls_type, goto_table in datas:
        flow = offlow.flow_mod(
            match=ofmatch.Match().mpls_type(mpls_type, True),
            actions=[
                ofaction.goto_table(goto_table),
            ],
            writes=[],
            table_id=pb.FlowMod.MPLS_TYPE,
            priority=1,
        )
        ofctl.mod_flow_entry(dpath, flow, dpath.ofproto.OFPFC_ADD)

    return


def setup_flows(dpath, ofctl):
    """
    OFDPA built-in flows to simulate OFDPA.
    """
    _add_table_miss_flows(dpath, ofctl)
    _add_mpls_l3_type_flows(dpath, ofctl)
    _add_mpls_type_flows(dpath, ofctl)

    return
