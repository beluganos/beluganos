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
from goryu.ofproto import ofmatch
from goryu.ofproto import ofaction
from goryu.ofproto import offlow
from goryu.ofproto import ofgroup
from fabricflow.fibc.api import fibcapi
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.ofc import ofdpa2_builtin

_LOG = logging.getLogger("fibc.gen.generic")

def _lagopus_bugfix(dpath, ofctl):
    """
    To avoid lagopus bug.
    """
    flow = offlow.flow_mod(
        match=ofmatch.Match().eth_type(fibcapi.ETHTYPE_MPLS),
        actions=[],
        writes=[],
        table_id=pb.FlowMod.POLICY_ACL,
        priority=fibcapi.PRIORITY_HIGHEST,
    )
    ofctl.mod_flow_entry(dpath, flow, dpath.ofproto.OFPFC_ADD)


def setup_term_mac_flow(dpath, ofctl):
    """
    Termination MAC flow (for setup)
    """
    matches = [
        ofmatch.Match().eth_dst(fibcapi.HWADDR_MULTICAST4_MATCH).eth_type(fibcapi.ETHTYPE_IPV4),
        ofmatch.Match().eth_dst(fibcapi.HWADDR_MULTICAST6_MATCH).eth_type(fibcapi.ETHTYPE_IPV6),
    ]
    actions = [ofaction.goto_table(pb.FlowMod.MULTICAST_ROUTING)]
    for match in matches:
        flow = offlow.flow_mod(
            match=match,
            actions=actions,
            writes=[],
            table_id=pb.FlowMod.TERM_MAC,
            priority=2,
        )

        ofctl.mod_flow_entry(dpath, flow, dpath.ofproto.OFPFC_ADD)


# pylint: disable=line-too-long
def setup_policy_acl_flow(dpath, ofctl):
    """
    Policy ACL flows when dp enter.
    """
    matches = [
        ofmatch.Match().eth_type(fibcapi.ETHTYPE_LACP),
        ofmatch.Match().eth_type(fibcapi.ETHTYPE_ARP),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV4).ip_proto(fibcapi.IPPROTO_ICMP4),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV6).ip_proto(fibcapi.IPPROTO_ICMP6),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV4).ip_proto(fibcapi.IPPROTO_OSPF),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV6).ip_proto(fibcapi.IPPROTO_OSPF),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV4).ip_proto(fibcapi.IPPROTO_TCP).tcp_src(fibcapi.TCPPORT_BGP),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV4).ip_proto(fibcapi.IPPROTO_TCP).tcp_dst(fibcapi.TCPPORT_BGP),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV6).ip_proto(fibcapi.IPPROTO_TCP).tcp_src(fibcapi.TCPPORT_BGP),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV6).ip_proto(fibcapi.IPPROTO_TCP).tcp_dst(fibcapi.TCPPORT_BGP),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV4).ip_proto(fibcapi.IPPROTO_TCP).tcp_src(fibcapi.TCPPORT_LDP),
        # ofmatch.Match().eth_type(fibcapi.ETHTYPE_IPV4).ip_proto(fibcapi.IPPROTO_TCP).tcp_dst(fibcapi.TCPPORT_LDP),
        ofmatch.Match().ip_dst(fibcapi.MCADDR_ALLROUTERS),
        ofmatch.Match().ip_dst(fibcapi.MCADDR_OSPF_HELLO),
        ofmatch.Match().ip_dst(fibcapi.MCADDR_OSPF_ALLDR),
    ]
    actions = [ofaction.output(dpath.ofproto.OFPP_CONTROLLER)]

    for match in matches:
        flow = offlow.flow_mod(
            match=match,
            actions=actions,
            writes=[],
            table_id=pb.FlowMod.POLICY_ACL,
            priority=fibcapi.PRIORITY_NORMAL,
        )

        ofctl.mod_flow_entry(dpath, flow, dpath.ofproto.OFPFC_ADD)


def setup_flow(dpath, mod, ofctl, ofdpa_sim=True):
    """
    Setup flows.
    """
    _LOG.debug("Default FLow: %d %s", dpath.id, mod)

    dpath.send_msg(offlow.clear_all(dpath))
    for msg in ofgroup.clear_all(dpath):
        dpath.send_msg(msg)

    if ofdpa_sim:
        ofdpa2_builtin.setup_flows(dpath, ofctl)
        _lagopus_bugfix(dpath, ofctl)

    setup_term_mac_flow(dpath, ofctl)
    setup_policy_acl_flow(dpath, ofctl)


def vlan_flow(dpath, mod, ofctl, use_metadata=True):
    """
    VLAN flow table.
    """
    _LOG.debug("VLAN FLow: %d %s", dpath.id, mod)

    cmd = fibcapi.flow_mod_cmd(mod.cmd, dpath.ofproto)
    entry = mod.vlan

    def _match():
        match = ofmatch.Match()
        match.in_port(entry.match.in_port)
        match.vlan_vid(entry.match.vid, entry.match.vid_mask)
        return match

    def _actions():
        if not offlow.is_action_needed(dpath, cmd):
            return []

        actions = []
        for action in entry.actions:
            if action.name == pb.VLANFlow.Action.PUSH_VLAN:
                vlan_vid = action.value | fibcapi.OFPVID_PRESENT
                actions.append(ofaction.push_vlan(fibcapi.ETHTYPE_VLAN_Q))
                actions.append(ofaction.set_field("vlan_vid", vlan_vid))

            elif action.name == pb.VLANFlow.Action.SET_VRF:
                actions.append(ofaction.set_vrf(action.value, use_metadata))

        actions.append(ofaction.goto_table(entry.goto_table))
        return actions

    flow = offlow.flow_mod(
        match=_match, actions=_actions, writes=[],
        table_id=pb.FlowMod.VLAN, priority=3)

    ofctl.mod_flow_entry(dpath, flow, cmd)


def termination_mac_flow(dpath, mod, ofctl):
    """
    Termination MAC flow table.
    """
    _LOG.debug("TERM MAC FLow: %d %s %s", dpath.id, mod, ofctl)

    cmd = fibcapi.flow_mod_cmd(mod.cmd, dpath.ofproto)
    entry = mod.term_mac

    match = ofmatch.Match().eth_type(entry.match.eth_type).eth_dst(entry.match.eth_dst)

    def _actions():
        if not offlow.is_action_needed(dpath, cmd):
            return []

        return [ofaction.goto_table(entry.goto_table)]

    flow = offlow.flow_mod(
        match=match, actions=_actions, writes=[],
        table_id=pb.FlowMod.TERM_MAC, priority=fibcapi.PRIORITY_LOW)

    ofctl.mod_flow_entry(dpath, flow, cmd)


def mpls1_flow(dpath, mod, ofctl, use_metadata=True):
    """
    MPLS1 flow table.
    """
    _LOG.debug("MPLS1 FLow: %d %s %s", dpath.id, mod, ofctl)

    cmd = fibcapi.flow_mod_cmd(mod.cmd, dpath.ofproto)
    entry = mod.mpls1

    def _match():
        match = ofmatch.Match()
        match.eth_type(fibcapi.ETHTYPE_MPLS)
        match.mpls_bos(entry.match.bos)
        match.mpls_label(entry.match.label)
        return match

    def _actions():
        if not offlow.is_action_needed(dpath, cmd):
            return []

        actions = [
            ofaction.goto_table(entry.goto_table),
            ofaction.dec_mpls_ttl(),
        ]

        if entry.goto_table == pb.FlowMod.MPLS_TYPE:
            actions.append(ofaction.set_mpls_type(fibcapi.MPLSTYPE_PHP, use_metadata))

        for action in entry.actions:
            if action.name == pb.MPLSFlow.Action.POP_LABEL:
                actions.append(ofaction.pop_mpls(action.value))

            elif action.name == pb.MPLSFlow.Action.SET_VRF:
                actions.append(ofaction.set_vrf(action.value, use_metadata))

        return actions

    def _writes():
        if not offlow.is_action_needed(dpath, cmd):
            return []
        if entry.g_type == pb.GroupMod.MPLS_INTERFACE:
            return [ofaction.group(fibcapi.mpls_interface_group_id(entry.g_id))]
        if entry.g_type == pb.GroupMod.MPLS_SWAP:
            return [ofaction.group(fibcapi.mpls_label_group_id(5, entry.g_id))]
        if entry.g_type == pb.GroupMod.MPLS_FF:
            return [ofaction.group(fibcapi.mpls_ff_group_id(entry.g_id))]
        if entry.g_type == pb.GroupMod.MPLS_ECMP:
            return [ofaction.group(fibcapi.mpls_ecmp_group_id(entry.g_id))]
        return []

    flow = offlow.flow_mod(
        match=_match, actions=_actions, writes=_writes,
        table_id=pb.FlowMod.MPLS1, priority=1)

    ofctl.mod_flow_entry(dpath, flow, cmd)


def unicast_routing_flow(dpath, mod, ofctl, use_metadata=True):
    """
    Create flow_mod for Unicast Routing flow table.
    """
    _LOG.debug("Unicast Routing FLow: %d %s", dpath.id, mod)

    cmd = fibcapi.flow_mod_cmd(mod.cmd, dpath.ofproto)
    entry = mod.unicast

    match = ofmatch.Match().ip_dst(entry.match.ip_dst).vrf(entry.match.vrf, use_metadata)
    def _actions():
        if not offlow.is_action_needed(dpath, cmd):
            return []
        return [ofaction.goto_table(pb.FlowMod.POLICY_ACL)]

    def _writes():
        if not offlow.is_action_needed(dpath, cmd):
            return []

        writes = [ofaction.dec_nw_ttl()]
        if entry.g_type == pb.GroupMod.L3_UNICAST:
            writes.append(ofaction.group(fibcapi.l3_unicast_group_id(entry.g_id)))
        elif entry.g_type == pb.GroupMod.L3_ECMP:
            writes.append(ofaction.group(fibcapi.l3_ecmp_group_id(entry.g_id)))
        elif entry.g_type == pb.GroupMod.MPLS_L3_VPN:
            writes.append(ofaction.group(fibcapi.mpls_label_group_id(2, entry.g_id)))
        else:
            pass

        return writes

    def _priority_base():
        if entry.g_type == pb.GroupMod.MPLS_L3_VPN:
            return fibcapi.PRIORITY_BASE_VPN

        return fibcapi.PRIORITY_BASE_UC

    priority = offlow.priority_for_ipaddr(entry.match.ip_dst, _priority_base())

    flow = offlow.flow_mod(
        match=match, actions=_actions, writes=_writes,
        table_id=pb.FlowMod.UNICAST_ROUTING, priority=priority)

    ofctl.mod_flow_entry(dpath, flow, cmd)


def bridging_flow(dpath, mod, ofctl):
    """
    Bridging flow table.
    """
    _LOG.debug("Bridging FLow: %d %s %s", dpath.id, mod, ofctl)


def policy_acl_flow(dpath, mod, ofctl, use_metadata=True):
    """
    Policy ACL flow table.
    """
    _LOG.debug("ACL FLow: %d %s", dpath.id, mod)

    cmd = fibcapi.flow_mod_cmd(mod.cmd, dpath.ofproto)
    entry = mod.acl
    match = ofmatch.Match().ip_dst(entry.match.ip_dst).vrf(entry.match.vrf, use_metadata)
    def _actions():
        if not offlow.is_action_needed(dpath, cmd):
            return []
        return [ofaction.output(dpath.ofproto.OFPP_CONTROLLER)]

    flow = offlow.flow_mod(
        match=match, actions=_actions, writes=[],
        table_id=pb.FlowMod.POLICY_ACL, priority=fibcapi.PRIORITY_HIGH)

    ofctl.mod_flow_entry(dpath, flow, cmd)


def setup_group(dpath, mod, ofctl):
    """
    Setup Group.
    """
    _LOG.debug("Default Group: %d %s %s", dpath.id, mod, ofctl)


def l2_interface_group(dpath, mod, ofctl):
    """
    L2 Interface Group
    """
    _LOG.debug("L2 Interface Group: %d %s", dpath.id, mod)

    ofproto = dpath.ofproto
    entry = mod.l2_iface
    cmd = fibcapi.group_mod_cmd(mod.cmd, dpath.ofproto)
    gid = fibcapi.l2_interface_group_id(entry.port_id, entry.vlan_vid)
    def _buckets():
        if not ofgroup.is_bucket_needed(dpath, cmd):
            return []

        actions = [ofaction.output(entry.port_id)]
        if entry.vlan_vid == ofproto.OFPVID_NONE:
            actions.insert(0, ofaction.pop_vlan())
        return [dict(actions=actions)]

    group = ofgroup.group_mod(gid, "INDIRECT", _buckets)
    ofctl.mod_group_entry(dpath, group, cmd)


def l3_unicast_group(dpath, mod, ofctl):
    """
    L3 Unicast Group
    """
    _LOG.debug("L3 Unicast Group: %d %s", dpath.id, mod)

    entry = mod.l3_unicast
    cmd = fibcapi.group_mod_cmd(mod.cmd, dpath.ofproto)
    gid = fibcapi.l3_unicast_group_id(entry.ne_id)
    def _buckets():
        if not ofgroup.is_bucket_needed(dpath, cmd):
            return []

        next_gid = fibcapi.l2_interface_group_id(entry.port_id, entry.vlan_vid)
        vlan_vid = fibcapi.adjust_vlan_vid(entry.vlan_vid) | fibcapi.OFPVID_PRESENT
        actions = [
            ofaction.set_field("eth_src", entry.eth_src),
            ofaction.set_field("eth_dst", entry.eth_dst),
            ofaction.set_field("vlan_vid", vlan_vid),
            ofaction.group(next_gid),
        ]
        return [dict(actions=actions)]

    group = ofgroup.group_mod(gid, "INDIRECT", _buckets)
    ofctl.mod_group_entry(dpath, group, cmd)


def l3_ecmp_group(dpath, mod, ofctl):
    """
    ECMP Group
    """
    _LOG.debug("L3 ECMP Group: %d %s %s", dpath.id, mod, ofctl)


def mpls_interface_group(dpath, mod, ofctl):
    """
    MPLS Interface group
    """
    _LOG.debug("MPLS Interface Group: %d %s", dpath.id, mod)

    entry = mod.mpls_iface
    cmd = fibcapi.group_mod_cmd(mod.cmd, dpath.ofproto)
    gid = fibcapi.mpls_interface_group_id(entry.ne_id)
    def _buckets():
        if not ofgroup.is_bucket_needed(dpath, cmd):
            return []

        next_gid = fibcapi.l2_interface_group_id(entry.port_id, entry.vlan_vid)
        vlan_vid = fibcapi.adjust_vlan_vid(entry.vlan_vid) | fibcapi.OFPVID_PRESENT
        actions = [
            ofaction.set_field("eth_src", entry.eth_src),
            ofaction.set_field("eth_dst", entry.eth_dst),
            ofaction.set_field("vlan_vid", vlan_vid),
            ofaction.group(next_gid),
        ]
        return [dict(actions=actions)]

    group = ofgroup.group_mod(gid, "INDIRECT", _buckets)
    ofctl.mod_group_entry(dpath, group, cmd)


def mpls_l3_vpn_group(dpath, mod, ofctl, mpls_bos=True):
    """
    MPLS L3 VPN Group
    """
    _LOG.debug("MPLS L3 VPN Group: %d %s %s", dpath.id, mod, ofctl)

    def get_next_gid(entry):
        """
        Get Next Group Id
        """
        if entry.ne_id != 0:
            return fibcapi.mpls_interface_group_id(entry.ne_id)
        elif entry.new_dst_id != 0:
            return fibcapi.mpls_label_group_id(3, entry.new_dst_id)

        return None

    entry = mod.mpls_label
    cmd = fibcapi.group_mod_cmd(mod.cmd, dpath.ofproto)
    gid = fibcapi.mpls_label_group_id(2, entry.dst_id)
    def _buckets():
        if not ofgroup.is_bucket_needed(dpath, cmd):
            return []

        next_gid = get_next_gid(entry)
        actions = [
            ofaction.push_mpls(fibcapi.ETHTYPE_MPLS),
            ofaction.set_field("mpls_label", entry.new_label),
        ]

        if mpls_bos:
            actions.append(ofaction.set_field("mpls_bos", 1))

        actions.append(ofaction.set_mpls_ttl(64))
        actions.append(ofaction.group(next_gid))

        return [dict(actions=actions)]

    group = ofgroup.group_mod(gid, "INDIRECT", _buckets)
    ofctl.mod_group_entry(dpath, group, cmd)


def mpls_tun1_group(dpath, mod, ofctl):
    """
    MPLS Tunnel1 Label Group
    """
    _LOG.debug("MPLS Tunnel1 Group: %d %s %s", dpath.id, mod, ofctl)

    entry = mod.mpls_label
    cmd = fibcapi.group_mod_cmd(mod.cmd, dpath.ofproto)
    gid = fibcapi.mpls_label_group_id(3, entry.dst_id)
    def _buckets():
        if not ofgroup.is_bucket_needed(dpath, cmd):
            return []

        next_gid = fibcapi.mpls_interface_group_id(entry.ne_id)
        actions = [
            ofaction.push_mpls(fibcapi.ETHTYPE_MPLS),
            ofaction.set_field("mpls_label", entry.new_label),
            ofaction.set_mpls_ttl(64),
            ofaction.group(next_gid),
        ]
        return [dict(actions=actions)]

    group = ofgroup.group_mod(gid, "INDIRECT", _buckets)
    ofctl.mod_group_entry(dpath, group, cmd)


def mpls_swap_group(dpath, mod, ofctl):
    """
    MPLS Swap Label Group
    """
    _LOG.debug("MPLS Swap Group: %d %s %s", dpath.id, mod, ofctl)

    entry = mod.mpls_label
    cmd = fibcapi.group_mod_cmd(mod.cmd, dpath.ofproto)
    gid = fibcapi.mpls_label_group_id(5, entry.dst_id)
    def _buckets():
        if not ofgroup.is_bucket_needed(dpath, cmd):
            return []

        next_gid = fibcapi.mpls_interface_group_id(entry.ne_id)
        actions = [
            ofaction.set_field("mpls_label", entry.new_label),
            ofaction.group(next_gid),
        ]
        return [dict(actions=actions)]

    group = ofgroup.group_mod(gid, "INDIRECT", _buckets)
    ofctl.mod_group_entry(dpath, group, cmd)


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


def pkt_out(dpath, port_id, strip_vlan, data):
    """
    PacketOut
    """
    _LOG.debug("PacketOUT: %s %s", dpath.id, len(data))

    parser = dpath.ofproto_parser
    ofp = dpath.ofproto

    actions = [parser.OFPActionOutput(port_id)]
    if strip_vlan:
        actions.insert(0, parser.OFPActionPopVlan())

    msg = parser.OFPPacketOut(datapath=dpath,
                              buffer_id=ofp.OFP_NO_BUFFER,
                              in_port=ofp.OFPP_ANY,
                              actions=actions,
                              data=data)

    dpath.send_msg(msg)


def get_port_stats(dpath, waiters, port_id, ofctl):
    """
    get port stats
    """
    # pylint: disable=unused-argument
    _LOG.debug("get_port_stats: %d %s", dpath.id, port_id)

    stats = ofctl.get_port_stats(dpath, waiters, port_id)
    return _fix_port_stats_names(stats)


def _fix_port_stats_names(stats):
    key_map = {
        'rx_packets': "ifInUcastPkts",
        'tx_packets': "ifOutUcastPkts",
        'rx_bytes': "ifInOctets",
        'tx_bytes': "ifOutOctets",
        'rx_dropped':"ifInDiscards",
        'tx_dropped': "ifOutDiscards",
        'rx_errors': "ifInErrors",
        'tx_errors': "ifOutErrors",
        'rx_frame_err': "etherStatsJabbers",
        'rx_over_err': "etherStatsOversizePkts",
        'rx_crc_err': "etherStatsCRCAlignErrors",
        'collisions': "etherStatsCollisions",
    }

    def _new_port_stats(port_stats):
        port_stats["ifInNUcastPkts"] = 0
        port_stats["ifOutNUcastPkts"] = 0
        return {key_map.get(key, key): val for key, val in port_stats.items()}

    def _new_port_stats_list(port_stats_list):
        return [_new_port_stats(port_stats) for port_stats in port_stats_list]

    return {dpid:_new_port_stats_list(port_stats_list) for dpid, port_stats_list in stats.items()}


def port_mod(dpath, mod, ofctl):
    """
    PotMod
    mod: api.FFPortMod
    """
    _LOG.debug("port_mod: %s %s", dpath, mod)

    parser = dpath.ofproto_parser
    ofp = dpath.ofproto

    config = 0 if port.state == pb.PortStatus.UP else ofp.OFPPC_PORT_DOWN
    msg = parser.OFPPortMod(
        port_no=mod.port_no,
        hw_addr=mod.hw_addr,
        config=config,
        mask=ofp.OFPPC_PORT_DOWN,
        advertise=0,
    )

    dpath.send_msg(msg)
