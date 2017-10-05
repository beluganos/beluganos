# -*- coding: utf-8 -*-

"""
fibapi helper module.
"""

import fabricflow.fibc.api.fibcapi_pb2 as pb

# ethertypes
ETHTYPE_IPV4 = 0x0800
ETHTYPE_IPV6 = 0x86dd
ETHTYPE_MPLS = 0x8847
ETHTYPE_LACP = 0x8809
ETHTYPE_ARP = 0x0806
ETHTYPE_VLAN_Q = 0x8100
ETHTYPE_VLAN_AD = 0x88a8

# hardware addresses
HWADDR_MULTICAST4 = '01:00:5e:00:00:00'
HWADDR_MULTICAST4_MASK = 'ff:ff:ff:80:00:00'
HWADDR_MULTICAST4_MATCH = HWADDR_MULTICAST4 + '/' + HWADDR_MULTICAST4_MASK
HWADDR_MULTICAST6 = '33:33:00:00:00:00'
HWADDR_MULTICAST6_MASK = 'ff:ff:00:00:00:00'
HWADDR_MULTICAST6_MATCH = HWADDR_MULTICAST6 + '/' + HWADDR_MULTICAST6_MASK

# vlan
OFPVID_UNTAGGED = 0x0001
OFPVID_PRESENT = 0x1000
OFPVID_NONE = 0x0000
OFPVID_ABSENT = 0x0000

def adjust_vlan_vid(vid):
    """
    Fix VLAN ID to UNTAGGED VID if vid is 0
    """
    return vid if vid != OFPVID_NONE else OFPVID_UNTAGGED


# protocol numbers
IPPROTO_ICMP4 = 1
IPPROTO_TCP = 6
IPPROTO_UDP = 17
IPPROTO_ICMP6 = 58
IPPROTO_OSPF = 89

# port numbers
TCPPORT_BGP = 179
TCPPORT_LDP = 646

# multicast addresses
MCADDR_ALLHOSTS = "224.0.0.1"
MCADDR_ALLROUTERS = "224.0.0.2"
MCADDR_OSPF_HELLO = "224.0.0.5"
MCADDR_OSPF_ALLDR = "224.0.0.6"

# Flow priority
PRIORITY_DEFAULT = 0
PRIORITY_LOW = 16400
PRIORITY_NORMAL = PRIORITY_LOW*2
PRIORITY_HIGH = PRIORITY_LOW*3
PRIORITY_HIGHEST = 65530
PRIORITY_BASE_VPN = PRIORITY_LOW
PRIORITY_BASE_UC = PRIORITY_LOW+1

# OFDPA Extensions
MPLSTYPE_NONE = 0x00
MPLSTYPE_VPS = 0x01
MPLSTYPE_UNICAST = 0x08
MPLSTYPE_MULTICAST = 0x10
MPLSTYPE_PHP = 0x20

def flow_mod_cmd(cmd, ofproto):
    """
    Convert FlowMod command (API to Ryu)
    """
    cmd_map = {
        pb.FlowMod.ADD          : ofproto.OFPFC_ADD,
        pb.FlowMod.MODIFY       : ofproto.OFPFC_MODIFY,
        pb.FlowMod.MODIFY_STRICT: ofproto.OFPFC_MODIFY_STRICT,
        pb.FlowMod.DELETE       : ofproto.OFPFC_DELETE,
        pb.FlowMod.DELETE_STRICT: ofproto.OFPFC_DELETE_STRICT,
    }
    return cmd_map.get(cmd, None)


def group_mod_cmd(cmd, ofproto):
    """
    Convert GroupMod command (API to Ryu)
    """
    cmd_map = {
        pb.GroupMod.ADD          : ofproto.OFPGC_ADD,
        pb.GroupMod.MODIFY       : ofproto.OFPGC_MODIFY,
        pb.GroupMod.DELETE       : ofproto.OFPGC_DELETE,
    }
    return cmd_map.get(cmd, None)


def parse_hello(data):
    """
    Parse binary data (pb.Hello)
    """
    msg = pb.Hello()
    msg.ParseFromString(data)
    return msg


def parse_port_config(data):
    """
    Parse binary data (pb.PortConfig)
    """
    msg = pb.PortConfig()
    msg.ParseFromString(data)
    return msg


def parse_flow_mod(data):
    """
    Parse binary data (pb.FlowMod)
    """
    msg = pb.FlowMod()
    msg.ParseFromString(data)
    return msg


def parse_group_mod(data):
    """
    Parse binary data (pb.GroupMod)
    """
    msg = pb.GroupMod()
    msg.ParseFromString(data)
    return msg


def l2_interface_group_id(port_id, vlan_vid):
    """
    L2 Interface Group Id
    """
    vlan_vid = adjust_vlan_vid(vlan_vid)
    return ((vlan_vid << 16) & 0x0fff0000) + (port_id & 0xffff)


def l2_rewrite_group_id(ne_id):
    """
    L2 Rewrite Group Id
    """
    return 0x10000000 + (ne_id & 0x0fffffff)


def l3_unicast_group_id(ne_id):
    """
    L3 Unicast Group Id
    """
    return 0x20000000 + (ne_id & 0x0fffffff)


def l2_multicast_group_id(mc_id, vlan_vid):
    """
    L2 Multicast Group Id
    """
    vlan_vid = adjust_vlan_vid(vlan_vid)
    return 0x30000000 + ((vlan_vid << 16) & 0x0fff0000) + (mc_id & 0xffff)


def l2_flood_group_id(fd_id, vlan_vid):
    """
    L2 Flood Group Id
    """
    vlan_vid = adjust_vlan_vid(vlan_vid)
    return 0x40000000 + ((vlan_vid << 16) & 0x0fff0000) + (fd_id & 0xffff)


def l3_interface_group_id(ne_id):
    """
    L3 Interface Group Id
    """
    return 0x50000000 + (ne_id & 0x0fffffff)


def l3_multicast_group_id(mc_id, vlan_vid):
    """
    L3 Multicast Group Id
    """
    vlan_vid = adjust_vlan_vid(vlan_vid)
    return 0x60000000 + ((vlan_vid << 16) & 0x0fff0000) + (mc_id & 0xffff)


def l3_ecmp_group_id(ecmp_id):
    """
    L3 ECMP Group Id
    """
    return 0x70000000 + (ecmp_id & 0x0fffffff)


def l2_overlay_group_id(tun_id, sub_type, index):
    """
    L2 Overlay Group Id
    """
    return 0x80000000 + ((tun_id << 12) & 0x0ffff000) + \
        ((sub_type << 10) & 0x0800) + (index & 0x07ff)


def mpls_interface_group_id(ne_id):
    """
    MPLS Interface Group Id
    """
    return 0x90000000 + (ne_id & 0x00ffffff)


def mpls_label_group_id(sub_type, label):
    """
    MPLS Label Group Id
    sub_type:
    - 1: L2 VPN Label
    - 2: L3 VPN Label
    - 3: Tunnel Label 1
    - 4: Tunnel Label 2
    - 5: Swap Label
    """
    return 0x90000000 + ((sub_type << 24) & 0x0f000000) + (label & 0x00ffffff)


def mpls_ff_group_id(index):
    """
    MPLS Fast Failover Group Id
    """
    return 0xa6000000 + (index & 0x00ffffff)


def mpls_ecmp_group_id(index):
    """
    MPLS ECMP Group Id
    """
    return 0xa8000000 + (index & 0x00ffffff)


def l2_unfiltered_iface_group_id(port_id):
    """
    L2 Unfiltered Interface Group
    """
    return 0xb0000000 + (port_id & 0xffff)
