// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package gonslib

import (
	api "gonsl/api"

	"github.com/beluganos/go-opennsl/opennsl"
)

//
// NewEthTypeFieldEntryAPI returns new instance
//
func NewEthTypeFieldEntryAPI(ethType uint16, inPort opennsl.Port) *api.EthTypeFieldEntry {
	return &api.EthTypeFieldEntry{
		EthType: uint32(ethType),
		InPort:  uint32(inPort),
	}
}

//
// NewFieldEntryEthTypeAPI returns new instance
//
func NewFieldEntryEthTypeAPI(ethType uint16, inPort opennsl.Port) *api.FieldEntry {
	return &api.FieldEntry{
		EntryType: api.FieldEntry_ETH_TYPE,
		Entry: &api.FieldEntry_EthType{
			EthType: NewEthTypeFieldEntryAPI(ethType, inPort),
		},
	}
}

//
// NewDstIPFieldEntryAPI returns new instance
//
func NewDstIPFieldEntryAPI(ethType uint16, dstIP string, inPort opennsl.Port) *api.DstIpFieldEntry {
	return &api.DstIpFieldEntry{
		EthType: uint32(ethType),
		IpDst:   dstIP,
		InPort:  uint32(inPort),
	}
}

//
// NewFieldEntryDstIPAPI returns new instance
//
func NewFieldEntryDstIPAPI(ethType uint16, dstIP string, inPort opennsl.Port) *api.FieldEntry {
	return &api.FieldEntry{
		EntryType: api.FieldEntry_DST_IP,
		Entry: &api.FieldEntry_DstIp{
			DstIp: NewDstIPFieldEntryAPI(ethType, dstIP, inPort),
		},
	}
}

//
// NewIPProtoFieldEntryAPI returns new instance
//
func NewIPProtoFieldEntryAPI(ethType uint16, ipProto uint8, inPort opennsl.Port) *api.IpProtoFieldEntry {
	return &api.IpProtoFieldEntry{
		EthType: uint32(ethType),
		IpProto: uint32(ipProto),
		InPort:  uint32(inPort),
	}
}

//
// NewFieldEntryIPProtoAPI returns new instance
//
func NewFieldEntryIPProtoAPI(ethType uint16, ipProto uint8, inPort opennsl.Port) *api.FieldEntry {
	return &api.FieldEntry{
		EntryType: api.FieldEntry_IP_PROTO,
		Entry: &api.FieldEntry_IpProto{
			IpProto: NewIPProtoFieldEntryAPI(ethType, ipProto, inPort),
		},
	}
}

//
// NewVlanEntryAPI returns new instance
//
func NewVlanEntryAPI(vid opennsl.Vlan, pbmp *opennsl.PBmp, upbmp *opennsl.PBmp) *api.VlanEntry {
	ports := []uint32{}
	pbmp.Each(func(port opennsl.Port) error {
		ports = append(ports, uint32(port))
		return nil
	})

	utPorts := []uint32{}
	upbmp.Each(func(port opennsl.Port) error {
		utPorts = append(utPorts, uint32(port))
		return nil
	})

	return &api.VlanEntry{
		Vid:        uint32(vid),
		Ports:      ports,
		UntagPorts: utPorts,
	}
}

//
// NewL2AddrAPI returns new instance.
//
func NewL2AddrAPI(l2addr *opennsl.L2Addr) *api.L2Addr {
	return &api.L2Addr{
		Flags: uint32(l2addr.Flags()),
		Mac:   l2addr.MAC().String(),
		Vid:   uint32(l2addr.VID()),
		Port:  uint32(l2addr.Port()),
	}
}

//
// NewL2StationAPI returns new instance.
//
func NewL2StationAPI(l2st *opennsl.L2Station) *api.L2Station {
	return &api.L2Station{
		Flags:       uint32(l2st.Flags()),
		DstMac:      l2st.DstMAC().String(),
		DstMacMask:  l2st.DstMACMask().String(),
		Vlan:        uint32(l2st.VID()),
		VlanMask:    uint32(l2st.VIDMask()),
		SrcPort:     uint32(l2st.SrcPort()),
		SrcPortMask: uint32(l2st.SrcPortMask()),
	}
}

//
// NewL3IfaceAPI returns new instance.
//
func NewL3IfaceAPI(l3iface *opennsl.L3Iface) *api.L3Iface {
	return &api.L3Iface{
		Flags:   uint32(l3iface.Flags()),
		IfaceId: uint32(l3iface.IfaceID()),
		Mac:     l3iface.MAC().String(),
		Mtu:     uint32(l3iface.MTU()),
		MtuFwd:  uint32(l3iface.MTUFwd()),
		Ttl:     uint32(l3iface.TTL()),
		Vid:     uint32(l3iface.VID()),
		Vrf:     uint32(l3iface.VRF()),
	}
}

//
// NewL3EgressAPI returns new instance.
//
func NewL3EgressAPI(l3egrID opennsl.L3EgressID, l3egr *opennsl.L3Egress) *api.L3Egress {
	return &api.L3Egress{
		Flags:    uint32(l3egr.Flags()),
		Flags2:   uint32(l3egr.Flags2()),
		EgressId: uint32(l3egrID),
		IfaceId:  uint32(l3egr.IfaceID()),
		Mac:      l3egr.MAC().String(),
		Vid:      uint32(l3egr.VID()),
		Port:     uint32(l3egr.Port()),
	}
}

//
// NewL3HostAPI returns new instance.
//
func NewL3HostAPI(host *opennsl.L3Host) *api.L3Host {
	return &api.L3Host{
		Flags:    uint32(host.Flags()),
		EgressId: uint32(host.EgressID()),
		IpAddr:   host.IPAddr().String(),
		Ip6Addr:  host.IP6Addr().String(),
		Mac:      host.NexthopMAC().String(),
		Vrf:      uint32(host.VRF()),
	}
}

//
// NewL3RouteAPI returns new instance.
//
func NewL3RouteAPI(route *opennsl.L3Route) *api.L3Route {
	return &api.L3Route{
		Flags:    uint32(route.Flags()),
		EgressId: uint32(route.EgressID()),
		IpAddr:   route.IP4Net().String(),
		Ip6Addr:  route.IP6Net().String(),
		Vrf:      uint32(route.VRF()),
	}
}

func NewTunnelInitiatorAPI(tunnel *opennsl.TunnelInitiator) *api.TunnelInitiator {
	dstIp, srcIp := func() (string, string) {
		switch tunnel.Type() {
		case opennsl.TunnelTypeIPIP4encap:
			return tunnel.DstIP4().String(), tunnel.SrcIP4().String()
		case opennsl.TunnelTypeIPIP6encap:
			return tunnel.DstIP6().String(), tunnel.SrcIP6().String()
		default:
			return tunnel.Type().String(), tunnel.Type().String()
		}
	}()

	return &api.TunnelInitiator{
		Flags:      uint32(tunnel.Flags()),
		TunnelId:   uint32(tunnel.TunnelID()),
		TunnelType: tunnel.Type().String(),
		L3IfaceId:  uint32(tunnel.L3IfaceID()),
		DstMac:     tunnel.DstMAC().String(),
		SrcMac:     tunnel.SrcMAC().String(),
		DstIp:      dstIp,
		SrcIp:      srcIp,
		DstPort:    uint32(tunnel.UdpDstPort()),
		SrcPort:    uint32(tunnel.UdpSrcPort()),
		Ttl:        uint32(tunnel.TTL()),
		Mtu:        uint32(tunnel.MTU()),
		Vlan:       uint32(tunnel.VID()),
	}
}

func NewTunnelTerminatorAPI(tunnel *opennsl.TunnelTerminator) *api.TunnelTerminator {
	dstIp, srcIp := func() (string, string) {
		switch tunnel.Type() {
		case opennsl.TunnelTypeIPIP4toIP4, opennsl.TunnelTypeIPIP4toIP6:
			return tunnel.DstIPNet4().String(), tunnel.SrcIPNet4().String()
		case opennsl.TunnelTypeIPIP6toIP4, opennsl.TunnelTypeIPIP6toIP6:
			return tunnel.DstIPNet6().String(), tunnel.SrcIPNet6().String()
		default:
			return tunnel.Type().String(), tunnel.Type().String()
		}
	}()

	return &api.TunnelTerminator{
		Flags:      uint32(tunnel.Flags()),
		TunnelId:   uint32(tunnel.TunnelID()),
		TunnelType: tunnel.Type().String(),
		RemotePort: uint32(tunnel.RemotePort()),
		DstIp:      dstIp,
		SrcIp:      srcIp,
		DstPort:    uint32(tunnel.UdpDstPort()),
		SrcPort:    uint32(tunnel.UdpSrcPort()),
		Vlan:       uint32(tunnel.VID()),
		Vrf:        uint32(tunnel.VRF()),
	}
}

//
// PortInfoAPI
//
func NewPortInfoAPI(port opennsl.Port, pinfo *opennsl.PortInfo) *api.PortInfo {
	return &api.PortInfo{
		Port:         uint32(port),
		LinkStatus:   int32(pinfo.LinkStatus()),
		UntaggedVlan: uint32(pinfo.UntaggedVlan()),
	}
}
