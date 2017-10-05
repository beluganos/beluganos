// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ofproto

import (
	"fmt"
	"strings"
)

type Match map[string]interface{}

func (m Match) String() string {
	slist := []string{}

	for name, value := range m {
		switch value.(type) {
		case float64:
			slist = append(slist, fmt.Sprintf("%s=0x%x", name, int64(value.(float64))))
		case string:
			slist = append(slist, fmt.Sprintf("%s=\"%s\"", name, value))
		default:
			slist = append(slist, fmt.Sprintf("%s='%v'", name, value))
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(slist, ","))
}

func (m *Match) Set(name string, value interface{}) *Match {
	(*m)[name] = value
	return m
}

func (m *Match) InPort(port uint32) *Match {
	return m.Set("in_port", port)
}

func (m *Match) InPhyPort(port uint32) *Match {
	return m.Set("in_phy_port", port)
}

func (m *Match) MetadataAndMask(metadata uint64, metadata_mask uint64) *Match {
	return m.Set("metadata", fmt.Sprintf("%d/%d", metadata, metadata_mask))
}

func (m *Match) Metadata(metadata uint64) *Match {
	return m.Set("metadata", metadata)
}

func (m *Match) EthDst(ethDst string) *Match {
	return m.Set("eth_dst", ethDst)
}

func (m *Match) EthSrc(ethSrc string) *Match {
	return m.Set("eth_src", ethSrc)
}

func (m *Match) EthType(ethType uint16) *Match {
	return m.Set("eth_type", ethType)
}

func (m *Match) VlanVidAndMask(vlanVid uint16, mask uint16) *Match {
	return m.Set("vlan_vid", fmt.Sprintf("%d/%d", vlanVid, mask))
}

func (m *Match) VlanVid(vlanVid uint16) *Match {
	if vlanVid == 0 {
		return m.Set("vlan_vid", "0x0")
	}

	return m.Set("vlan_vid", vlanVid)
}

func (m *Match) VlanPCP(pcp uint8) *Match {
	return m.Set("vlan_pcp", pcp)
}

func (m *Match) IPDscp(dscp uint8) *Match {
	return m.Set("ip_dscp", dscp)
}

func (m *Match) IPEcn(ecn uint8) *Match {
	return m.Set("ip_ecn", ecn)
}

func (m *Match) IPProto(proto uint8) *Match {
	return m.Set("ip_proto", proto)
}

func (m *Match) IPv4Src(ipv4Src string) *Match {
	return m.Set("ipv4_src", ipv4Src)
}

func (m *Match) IPv4Dst(ipv4Dst string) *Match {
	return m.Set("ipv4_dst", ipv4Dst)
}

func (m *Match) TCPSrc(tcpSrc uint16) *Match {
	return m.Set("tcp_src", tcpSrc)
}

func (m *Match) TCPDst(tcpDst uint16) *Match {
	return m.Set("tcp_dst", tcpDst)
}

func (m *Match) UDPSrc(udpSrc uint16) *Match {
	return m.Set("udp_src", udpSrc)
}

func (m *Match) UDPDst(udpDst uint16) *Match {
	return m.Set("udp_dst", udpDst)
}

func (m *Match) SctpSrc(sctpSrc uint16) *Match {
	return m.Set("sctp_src", sctpSrc)
}

func (m *Match) SctpDst(sctpDst uint16) *Match {
	return m.Set("sctp_dst", sctpDst)
}

func (m *Match) ICMPv4Code(code uint8) *Match {
	return m.Set("icmpv4_code", code)
}

func (m *Match) ICMPv4Type(t uint8) *Match {
	return m.Set("icmpv4_type", t)
}

func (m *Match) ARPOp(op uint16) *Match {
	return m.Set("arp_op", op)
}

func (m *Match) ARPSpa(spa string) *Match {
	return m.Set("arp_spa", spa)
}

func (m *Match) ARPTpa(tpa string) *Match {
	return m.Set("arp_tpa", tpa)
}

func (m *Match) ARPSha(sha string) *Match {
	return m.Set("arp_sha", sha)
}

func (m *Match) ARPTha(tha string) *Match {
	return m.Set("arp_tha", tha)
}

func (m *Match) IPv6Src(ipv6Src string) *Match {
	return m.Set("ipv6_src", ipv6Src)
}

func (m *Match) IPv6Dst(ipv6Dst string) *Match {
	return m.Set("ipv6_dst", ipv6Dst)
}

func (m *Match) IPv6Flabel(flabel uint32) *Match {
	return m.Set("ipv6_flabel", flabel)
}

func (m *Match) ICMPv6Code(code uint8) *Match {
	return m.Set("icmpv6_code", code)
}

func (m *Match) ICMPv6Type(t uint8) *Match {
	return m.Set("icmpv6_type", t)
}

func (m *Match) IPv6NDTarger(ipv6NDTarget string) *Match {
	return m.Set("ipv6_nd_target", ipv6NDTarget)
}

func (m *Match) IPv6NDSll(ipv6NDSll string) *Match {
	return m.Set("ipv6_nd_sll", ipv6NDSll)
}

func (m *Match) IPv6NDTll(ipv6NDTll string) *Match {
	return m.Set("ipv6_nd_tll", ipv6NDTll)
}

func (m *Match) MPLSLabel(label uint32) *Match {
	return m.Set("mpls_label", label)
}

func (m *Match) MPLSTc(tc uint8) *Match {
	return m.Set("mpls_tc", tc)
}

func (m *Match) MPLSBos(bos uint8) *Match {
	return m.Set("mpls_bos", bos)
}

func (m *Match) PBBIsid(isid uint32) *Match {
	return m.Set("pbb_isid", isid)
}

func (m *Match) PBBIsidAndMask(isid, mask uint32) *Match {
	return m.Set("pbb_isid", fmt.Sprintf("%d/%d", isid, mask))
}

func (m *Match) TunnelId(tunnelId uint64) *Match {
	return m.Set("tunnel_id", tunnelId)
}

func (m *Match) TunnelIdAndMask(tunnelId, mask uint64) *Match {
	return m.Set("tunnel_id", fmt.Sprintf("%d/%d", tunnelId, mask))
}

func (m *Match) IPv6ExtHdr(exthdr uint16) *Match {
	return m.Set("ipv6_exthdr", exthdr)
}

func (m *Match) IPv6ExtHdrAndMask(exthdr, mask uint16) *Match {
	return m.Set("ipv6_exthdr", fmt.Sprintf("%d/%d", exthdr, mask))
}

func (m *Match) Vrf(vrf uint8, useMetadata bool) *Match {
	if useMetadata {
		return m.MetadataAndMask(uint64(vrf)<<VRF_METADATA_SHIFT, VRF_METADATA_MASK)
	} else {
		return m.Set("vrf", vrf)
	}
}

func (m *Match) MPLSType(mplsType uint8, useMetadata bool) *Match {
	if useMetadata {
		return m.MetadataAndMask(uint64(mplsType)<<MPLSTYPE_METADATA_SHIFT, MPLSTYPE_METADATA_MASK)
	} else {
		return m.Set("mpls_type", mplsType)
	}
}

func (m *Match) IPDst(ipDst string) *Match {
	if strings.Contains(ipDst, ":") {
		return m.EthType(0x86dd).IPv6Dst(ipDst)
	} else {
		return m.EthType(0x0800).IPv4Dst(ipDst)
	}
}

func (m *Match) IPSrc(ipSrc string) *Match {
	if strings.Contains(ipSrc, ":") {
		return m.EthType(0x86dd).IPv6Src(ipSrc)
	} else {
		return m.EthType(0x0800).IPv4Src(ipSrc)
	}
}
