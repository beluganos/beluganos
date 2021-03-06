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
HWADDR_NONE = "00:00:00:00:00:00"
HWADDR_DUMMY = "02:00:00:00:00:00"
HWADDR_EXACT_MASK = "ff:ff:ff:ff:ff:ff"
HWADDR_MULTICAST4 = '01:00:5e:00:00:00'
HWADDR_MULTICAST4_MASK = 'ff:ff:ff:80:00:00'
HWADDR_MULTICAST4_MATCH = HWADDR_MULTICAST4 + '/' + HWADDR_MULTICAST4_MASK
HWADDR_MULTICAST6 = '33:33:00:00:00:00'
HWADDR_MULTICAST6_MASK = 'ff:ff:00:00:00:00'
HWADDR_MULTICAST6_MATCH = HWADDR_MULTICAST6 + '/' + HWADDR_MULTICAST6_MASK
HWADDR_ISIS_LEVEL1 = "01:80:C2:00:00:14"
HWADDR_ISIS_LEVEL2 = "01:80:C2:00:00:15"

def parse_masked_mac(s): # pylint: disable=invalid-name
    """
    parse mac and mask string. (xx:xx:..:xx/yy:yy:...:yy)
    """
    items = s.split("/", 1)
    if len(items) == 1:
        return items[0], "ff:ff:ff:ff:ff:ff"
    elif len(items) == 2:
        return items[0], items[1]
    else:
        raise SyntaxError(s)

def new_masked_mac(mac, mask):
    """
    make masked mac string.
    """
    return "{0}/{1}".format(mac, mask)

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

# multicast addresses(IPv6)
MCADDR6_I_LOCAL = "ff01::/64"    # Interface Local
MCADDR6_L_LOCAL = "ff02::/64"    # Link Local
MCADDR6_S_LOCAL = "ff05::/64"    # Site Local
MCADDR6_L_ALLNODES = "ff02::1"   # All Nodes / Link Local
MCADDR6_S_ALLNODES = "ff05::1"   # All Nodes / Site Local
MCADDR6_L_ALLROUTERS = "ff02::2" # All Routers / Link Local
MCADDR6_L_ALLOSPF = "ff02::5"    # All OSPF Routers / Link Local
MCADDR6_L_ALLOSPF_DR = "ff02::6" # All OSPF Routers / Link Local
MCADDR6_L_ALLRIP = "ff02::9"     # All RIP Routers / Link Local
MCADDR6_L_ALLEIGRP = "ff02::A"   # All EIGRP Routers / Link Local
MCADDR6_L_ALLPIM = "ff02::D"     # All PIM Routers / Link Local
MCADDR6_L_ALLDHCP = "ff02::1:2"  # All DHCP Agents / Link Local
MCADDR6_S_ALLDHCP = "ff05::1:3"  # All DHCP Servers / Site Local
MCADDR6_L_ALLNTP = "ff02::101"   # All NTP Servers / Link Local
MCADDR6_S_ALLNTP = "ff05::101"   # All NTP Servers / Site Local
# unicast addresses (IPv6)
UCADDR6_L_LOCAL = "fe80::/64"

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

# DP Port Types
DPPORT_ID_MASK = 0x0000ffff
DPPORT_VRF_MASK = 0x00ff0000
DPPORT_VRF_SHiFT = 16
DPPORT_TYPE_MASK = 0xff000000
DPPORT_TYPE_SHIFT = 24

def new_dpport_id(port_id, link_type):
    return (link_type << DPPORT_TYPE_SHIFT) + (port_id & DPPORT_ID_MASK)

def parse_dpport_id(port):
    return port & DPPORT_ID_MASK, (port & DPPORT_TYPE_MASK) >> DPPORT_TYPE_SHIFT

def flow_mod_cmd(cmd, ofproto):
    """
    Convert FlowMod command (API to Ryu)
    """
    # pylint: disable=no-member
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
    # pylint: disable=no-member
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
    return ((vlan_vid << 16) & 0x0fff0000) + (port_id & DPPORT_ID_MASK)


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
    return 0xb0000000 + (port_id & DPPORT_ID_MASK)


def parse_ff_hello(data):
    """
    Parse binary data (pb.FFHello)
    """
    msg = pb.FFHello()
    msg.ParseFromString(data)
    return msg


def parse_ff_port_status(data):
    """
    Parse binary data (pb.FFPortStatus)
    """
    msg = pb.FFPortStatus()
    msg.ParseFromString(data)
    return msg


def parse_ff_multipart_request(data):
    """
    Parse binary data (pb.FFMultipart_Request)
    """
    msg = pb.FFMultipart.Request()
    msg.ParseFromString(data)
    return msg


def parse_ff_multipart_reply(data):
    """
    Parse binary data (pb.FFMultipart_Reply)
    """
    msg = pb.FFMultipart.Reply()
    msg.ParseFromString(data)
    return msg


def parse_ff_packet_in(data):
    """
    Parse binary data (pb.FFPacketn)
    """
    msg = pb.FFPacketIn()
    msg.ParseFromString(data)
    return msg


def new_ff_packet_out(dp_id, port_no, data):
    """
    Create FFPacketOut messasge
    """
    return pb.FFPacketOut(
        dp_id=dp_id,
        port_no=port_no,
        data=data,
    )

def new_ff_multipart_request_port(dp_id, port_no, names):
    """
    Create Multipart Request (Port)
    """
    return pb.FFMultipart.Request(
        mp_type=pb.FFMultipart.PORT, # pylint: disable=no-member
        dp_id=dp_id,
        port=pb.FFMultipart.PortRequest(port_no=port_no, names=names),
    )

def new_ff_multipart_reply_port(dp_id, stats):
    """
    Create Multopart Reply (Port)
    """
    return pb.FFMultipart.Reply(
        mp_type=pb.FFMultipart.PORT, # pylint: disable=no-member
        dp_id=dp_id,
        port=pb.FFMultipart.PortReply(stats=stats)
    )

def new_ff_multipart_request_portdesc(dp_id, internal=False): # pylint: disable=invalid-name
    """
    Create Multipart Request (PortDesc)
    """
    return pb.FFMultipart.Request(
        mp_type=pb.FFMultipart.PORT_DESC, # pylint: disable=no-member
        dp_id=dp_id,
        port_desc=pb.FFMultipart.PortDescRequest(internal=internal)
    )

def new_ff_multipart_reply_portdesc(dp_id, ports, internal=False):
    """
    Create Multipart Reply (PortDesc)
    """
    return pb.FFMultipart.Reply(
        mp_type=pb.FFMultipart.PORT_DESC, # pylint: disable=no-member
        dp_id=dp_id,
        port_desc=pb.FFMultipart.PortDescReply(port=ports, internal=internal)
    )

def new_policy_acl_match(**kwargs):
    """
    Create PolicyACLFlow.Match
    """
    return pb.PolicyACLFlow.Match(
        in_port=kwargs.get("in_port", 0),
        ip_dst=kwargs.get("ip_dst", ""),
        vrf=kwargs.get("vrf", 0),
        eth_type=kwargs.get("eth_type", 0),
        ip_proto=kwargs.get("ip_proto", 0),
        tp_src=kwargs.get("tp_src", 0),
        tp_dst=kwargs.get("tp_dst", 0),
        eth_dst=kwargs.get("eth_dst", ""))

def new_policy_acl_action(name, value=0):
    """
    Create PolicyACLFlow.Action
    """
    return pb.PolicyACLFlow.Action(
        name=name,
        value=value)

def new_ff_port_mod(dp_id, port_no, status):
    """
    Create FFPortMod message
    """
    return pb.FFPortMod(
        dp_id=dp_id,
        port_no=port_no,
        status=status)

def bridgevlaninfo_type_from_flags(flags):
    """
    converts BridgeVlanInfo Flags to PortType.
    """
    if (flags & pb.BridgeVlanInfo.PVID) != 0 and (flags & pb.BridgeVlanInfo.UNTAGGED) != 0:
        return pb.BridgeVlanInfo.ACCESS_PORT

    if (flags & pb.BridgeVlanInfo.PVID) == 0 and (flags & pb.BridgeVlanInfo.UNTAGGED) == 0:
        return pb.BridgeVlanInfo.TRUNK_PORT

    return pb.BridgeVlanInfo.NONE_PORT


def parse_l2addr(data):
    """
    Parse binary data (pb.L2Addr)
    """
    msg = pb.L2Addr()
    msg.ParseFromString(data)
    return msg


def new_l2addr(hw_addr, vlan_vid, port_id, reason, ifname=""):
    """
    Create L2Addr message
    """
    return pb.L2Addr(
        hw_addr=hw_addr,
        vlan_vid=vlan_vid,
        port_id=(port_id & DPPORT_ID_MASK),
        ifname=ifname,
        reason=reason,
    )


def parse_ff_l2addr_status(data):
    """
    Parse binary data (pb.FFL2AddrStatus)
    """
    msg = pb.FFL2AddrStatus()
    msg.ParseFromString(data)
    return msg


def new_ff_l2addr_status(dp_id, addrs):
    """
    Create FFL2AddrStatus message
    """
    return pb.FFL2AddrStatus(
        dp_id=dp_id,
        addrs=addrs,
    )


def parse_l2addr_status(data):
    """
    Parse binary data (pb.L2AddrStatus)
    """
    msg = pb.L2AddrStatus()
    msg.ParseFromString(data)
    return msg


def new_l2addr_status(re_id, addrs):
    """
    Create L2AddrStatus message
    """
    return pb.L2AddrStatus(
        re_id=re_id,
        addrs=addrs,
    )
