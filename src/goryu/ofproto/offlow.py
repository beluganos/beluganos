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
Flow Mod functions
"""

from goryu.ofproto import ofaction

PRIORITY_BAND = 16
PRIORITY_BASE = 16400

def priority_for_ipaddr(ipaddr, base=PRIORITY_BASE):
    """
    priority of flow entry.
    """
    from netaddr.ip import IPNetwork
    ipnw = IPNetwork(ipaddr)
    return ipnw.prefixlen * PRIORITY_BAND + base


def is_action_needed(dpath, cmd):
    """
    judge actions are needed.
    """
    ofp = dpath.ofproto
    return cmd not in (ofp.OFPFC_DELETE, ofp.OFPFC_DELETE_STRICT)


def flow_mod(match, actions, writes, table_id, priority):
    """
    flow_mod
    """
    mtc = match() if callable(match) else match
    act = actions() if callable(actions) else actions
    wri = writes() if callable(writes) else writes
    return {
        "table_id": table_id,
        "priority": priority,
        "match"   : mtc,
        "actions" : act + [ofaction.write_actions(wri)],
    }


def clear_all(dpath):
    """
    Clear all flow tables.
    """
    # return {"table_id": dpath.ofproto.OFPTT_ALL}
    ofp = dpath.ofproto
    return dpath.ofproto_parser.OFPFlowMod(
        dpath, table_id=ofp.OFPTT_ALL, command=ofp.OFPFC_DELETE)


def set_sw_config(dpath):
    """
    Set Switch Config
    """
    ofp = dpath.ofproto
    return dpath.ofproto_parser.OFPSetConfig(
        dpath, ofp.OFPC_FRAG_NORMAL, ofp.OFPCML_MAX)
