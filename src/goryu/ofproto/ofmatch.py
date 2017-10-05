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

# pylint: disable=too-many-public-methods
class Match(dict):
    """
    Match
    """
    def _put(self, name, value):
        """
        Set match field
        """
        self[name] = value
        return self

    def in_port(self, port):
        """
        in_port=<num>
        """
        return self._put("in_port", port)

    def in_phy_port(self, port):
        """
        in_phy_port=<num>
        """
        return self._put("in_phy_port", port)

    def metadata(self, metadata):
        """
        metadata=<num> or "metadata/mask"
        """
        return self._put("metadata", metadata)

    def eth_dst(self, eth_dst):
        """
        eth_dst="hwaddr" or "hwaddr/mask"
        """
        return self._put("eth_dst", eth_dst)

    def eth_src(self, eth_src):
        """
        eth_src="hwaddr" or "hwaddr/mask"
        """
        return self._put("eth_src", eth_src)

    def eth_type(self, eth_type):
        """
        eth_type=<num>
        """
        return self._put("eth_type", eth_type)

    def vlan_vid(self, vlan_vid, mask=0):
        """
        vlan_vid=<decimal>  (PRESENT BIT added by ofctl)
        vlan_vid="decimal"  (PRESENT BIT added by ofctl)
        vlan_vid="0x..."    (PRESENT BIT not added)
        vlan_vid="vid/mask" (PRESENT BIT not added)
        """
        if not mask:
            if vlan_vid == 0:
                return self._put("vlan_vid", "0x0")

            return self._put("vlan_vid", vlan_vid)

        return self._put("vlan_vid", "{0}/{1}".format(vlan_vid, mask))

    def vlan_pcp(self, vlan_pcp):
        """
        vlan_pcp=<num>
        """
        return self._put("vlan_pcp", vlan_pcp)

    def ip_dscp(self, ip_dscp):
        """
        ip_dscp=<num>
        """
        return self._put("ip_dscp", ip_dscp)

    def ip_ecn(self, ip_ecn):
        """
        ip_ecn=<num>
        """
        return self._put("ip_ecn", ip_ecn)

    def ip_proto(self, ip_proto):
        """
        ip_proto=<num>
        """
        return self._put("ip_proto", ip_proto)

    def ipv4_src(self, ipv4_src):
        """
        ipv4_src="ip" or "ip/<num>" or "ip/mask"
        """
        return self._put("ipv4_src", ipv4_src)

    def ipv4_dst(self, ipv4_dst):
        """
        ipv4_dst="ip" or "ip/<num>" or "ip/mask"
        """
        return self._put("ipv4_dst", ipv4_dst)

    def tcp_src(self, tcp_src):
        """
        tcp_src=<num>
        """
        return self._put("tcp_src", tcp_src)

    def tcp_dst(self, tcp_dst):
        """
        tcp_dst=<num>
        """
        return self._put("tcp_dst", tcp_dst)

    def udp_src(self, udp_src):
        """
        udp_src=<num>
        """
        return self._put("udp_src", udp_src)

    def udp_dst(self, udp_dst):
        """
        udp_dst=<num>
        """
        return self._put("udp_dst", udp_dst)

    def sctp_src(self, sctp_src):
        """
        sctp_src=<num>
        """
        return self._put("sctp_src", sctp_src)

    def sctp_dst(self, sctp_dst):
        """
        sctp_dst=<num>
        """
        return self._put("sctp_dst", sctp_dst)

    def icmpv4_type(self, icmpv4_type):
        """
        icmpv4_type=<num>
        """
        return self._put("icmpv4_type", icmpv4_type)

    def icmpv4_code(self, icmpv4_code):
        """
        icmpv4_code=<num>
        """
        return self._put("icmpv4_code", icmpv4_code)

    def arp_op(self, arp_op):
        """
        arp_op=<num>
        """
        return self._put("arp_op", arp_op)

    def arp_spa(self, arp_spa):
        """
        arp_spa="ip" or "ip/<num>" or "ip/mask"
        """
        self._put("arp_spa", arp_spa)

    def arp_tpa(self, arp_tpa):
        """
        arp_tpa="ip" or "ip/<num>" or "ip/mask"
        """
        self._put("arp_tpa", arp_tpa)

    def arp_sha(self, arp_sha):
        """
        arp_sha="hwaddr" or "hwaddr/mask"
        """
        self._put("arp_sha", arp_sha)

    def arp_tha(self, arp_tha):
        """
        arp_tha="hwaddr" or "hwaddr/mask"
        """
        self._put("arp_tha", arp_tha)

    def ipv6_src(self, ipv6_src):
        """
        ipv6_src="ip" or "ip/<num>" or "ip/mask"
        """
        return self._put("ipv6_src", ipv6_src)

    def ipv6_dst(self, ipv6_dst):
        """
        ipv6_dst="ip" or "ip/<num>" or "ip/mask"
        """
        return self._put("ipv6_dst", ipv6_dst)

    def ipv6_flabel(self, ipv6_flabel):
        """
        ipv6_flabel=<num>
        """
        return self._put("ipv6_flabel", ipv6_flabel)

    def icmpv6_type(self, icmpv6_type):
        """
        icmpv6_type=<num>
        """
        return self._put("icmpv6_type", icmpv6_type)

    def icmpv6_code(self, icmpv6_code):
        """
        icmpv6_code=<num>
        """
        return self._put("icmpv6_code", icmpv6_code)

    def ipv6_nd_target(self, ipv6_nd_target):
        """
        ipv6_nd_target="ip" or "ip/<num>" or "ip/mask"
        """
        return self._put("ipv6_nd_target", ipv6_nd_target)

    def ipv6_nd_sll(self, ipv6_nd_sll):
        """
        ipv6_nd_sll="hwaddr" or "hwaddr/mask"
        """
        return self._put("ipv6_nd_sll", ipv6_nd_sll)

    def ipv6_nd_tll(self, ipv6_nd_tll):
        """
        ipv6_nd_tll="hwaddr" or "hwaddr/mask"
        """
        return self._put("ipv6_nd_tll", ipv6_nd_tll)

    def mpls_label(self, mpls_label):
        """
        mpls_label=<num>
        """
        return self._put("mpls_label", mpls_label)

    def mpls_tc(self, mpls_tc):
        """
        mpls_tc=<num>
        """
        return self._put("mpls_tc", mpls_tc)

    def mpls_bos(self, mpls_bos):
        """
        mpls_bos=<num>
        """
        def _fix_bos():
            return 1 if mpls_bos else 0
        return self._put("mpls_bos", _fix_bos())

    def pbb_isid(self, pbb_isid):
        """
        pbb_isid=<num> or "num" or "num/mask"
        """
        return self._put("pbb_isid", pbb_isid)

    def tunnel_id(self, tunnel_id):
        """
        tunnel_id=<num> or "num" or "num/mask"
        """
        return self._put("tunnel_id", tunnel_id)

    def ipv6_exthdr(self, ipv6_exthdr):
        """
        ipv6_exthdr=<num> or "num" or "num/mask"
        """
        return self._put("ipv6_exthdr", ipv6_exthdr)

    def vrf(self, vrf, use_metadata):
        """
        vrf=<num>
        """
        if use_metadata:
            return self.metadata("{0}/{1}".format(
                vrf << ofproto.VRF_METADATA_SHIFT, ofproto.VRF_METADATA_MASK))

        if vrf:
            return self._put("vrf", vrf)

        return self

    def mpls_type(self, mpls_type, use_metadata):
        """
        mpls_type=<num>
        """
        if use_metadata:
            return self.metadata("{0}/{1}".format(
                mpls_type << ofproto.MPLSTYPE_METADATA_SHIFT, ofproto.MPLSTYPE_METADATA_MASK))

        return self._put("mpls_type", mpls_type)

    def ip_dst(self, ip_dst):
        """
        ip_dst=<ipv4>/<mask> or <ipv6>/<mask>
        """
        if ":" in str(ip_dst):
            return self.eth_type(0x86dd).ipv6_dst(ip_dst)

        return self.eth_type(0x0800).ipv4_dst(ip_dst)

    def ip_src(self, ip_src):
        """
        ip_dst=<ipv4>/<mask> or <ipv6>/<mask>
        """
        if ":" in str(ip_src):
            return self.eth_type(0x86dd).ipv6_src(ip_src)

        return self.eth_type(0x0800).ipv4_src(ip_src)
