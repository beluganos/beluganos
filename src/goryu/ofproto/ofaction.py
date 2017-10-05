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
FIBC Ryu message helper functions
"""

from goryu.ofproto import ofproto

def _action(action_type, **kwargs):
    dic = dict(type=action_type)
    dic.update(**kwargs)
    return dic

def output(port):
    """
    OUTPUT Action
    """
    return _action("OUTPUT", port=port)

def copy_ttl_out():
    """
    COPY_TTL_OUT Action
    """
    return _action("COPY_TTL_OUT")

def copy_ttl_in():
    """
    COPY_TTL_IN Action
    """
    return _action("COPY_TTL_IN")

def set_mpls_ttl(ttl):
    """
    SET_MPLS_TTL Action
    """
    return _action("SET_MPLS_TTL", mpls_ttl=ttl)

def dec_mpls_ttl():
    """
    DEC_MPLS_TTL Action
    """
    return _action("DEC_MPLS_TTL")

def push_vlan(ethertype):
    """
    PUSH_VLAN
    """
    return _action("PUSH_VLAN", ethertype=ethertype)

def pop_vlan():
    """
    POP_VLAN
    """
    return _action("POP_VLAN")

def push_mpls(ethertype):
    """
    PUSH_MPLS
    """
    return _action("PUSH_MPLS", ethertype=ethertype)

def pop_mpls(ethertype):
    """
    POP_MPLS
    """
    return _action("POP_MPLS", ethertype=ethertype)

def set_queue(queue_id):
    """
    SET_QUEUE
    """
    return _action("SET_QUEUE", queue_id=queue_id)

def group(group_id):
    """
    GROUP Action
    """
    return _action("GROUP", group_id=group_id)

def set_nw_ttl(nw_ttl):
    """
    SET_NW_TTL Action
    """
    return _action("SET_NW_TTL", nw_ttl=nw_ttl)

def dec_nw_ttl():
    """
    DEC_NW_TTL Action
    """
    return _action("DEC_NW_TTL")

def set_field(field, value):
    """
    SET_FIELD Action
    """
    return _action("SET_FIELD", field=field, value=value)

def push_pbb(ethertype):
    """
    PUSH_PBB Action
    """
    return _action("PUSH_PBB", ethertype=ethertype)

def pop_pbb():
    """
    POP_PBB Action
    """
    return _action("POP_PBB")

def write_actions(actions):
    """
    WRITE_ACTIONS Action
    """
    return _action("WRITE_ACTIONS", actions=actions)

def clear_actions():
    """
    CLEAR_ACTIONS Action
    """
    return _action("CLEAR_ACTIONS")

def goto_table(table_id):
    """
    GOTO_TABLE Action
    """
    return _action("GOTO_TABLE", table_id=table_id)

def write_metadata(metadata, metadata_mask):
    """
    WRITE_METADATA Action
    """
    return _action("WRITE_METADATA", metadata=metadata, metadata_mask=metadata_mask)

def meter(meter_id):
    """
    METER Action
    """
    return _action("METER", meter_id=meter_id)

def set_vrf(vrf, use_metadata):
    """
    SET_FIELD(vrf) Action
    """
    if use_metadata:
        return write_metadata(
            vrf << ofproto.VRF_METADATA_SHIFT, ofproto.VRF_METADATA_MASK)

    return set_field("vrf", vrf)

def set_mpls_type(mpls_type, use_metadata):
    """
    SET_FIELD(mpls_type) Action
    """
    if use_metadata:
        return write_metadata(
            mpls_type << ofproto.MPLSTYPE_METADATA_SHIFT, ofproto.MPLSTYPE_METADATA_MASK)

    return set_field("mpls_type", mpls_type)
