// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

package gonslapi

import (
	"github.com/beluganos/go-opennsl/opennsl"
)

//
// NewEthTypeFieldEntry returns new instance
//
func NewEthTypeFieldEntry(ethType uint16, inPort opennsl.Port) *EthTypeFieldEntry {
	return &EthTypeFieldEntry{
		EthType: uint32(ethType),
		InPort:  uint32(inPort),
	}
}

//
// NewFieldEntryEthType returns new instance
//
func NewFieldEntryEthType(ethType uint16, inPort opennsl.Port) *FieldEntry {
	return &FieldEntry{
		EntryType: FieldEntry_ETH_TYPE,
		Entry: &FieldEntry_EthType{
			EthType: NewEthTypeFieldEntry(ethType, inPort),
		},
	}
}

//
// NewDstIPFieldEntry returns new instance
//
func NewDstIPFieldEntry(ethType uint16, dstIP string, inPort opennsl.Port) *DstIpFieldEntry {
	return &DstIpFieldEntry{
		EthType: uint32(ethType),
		IpDst:   dstIP,
		InPort:  uint32(inPort),
	}
}

//
// NewFieldEntryDstIP returns new instance
//
func NewFieldEntryDstIP(ethType uint16, dstIP string, inPort opennsl.Port) *FieldEntry {
	return &FieldEntry{
		EntryType: FieldEntry_DST_IP,
		Entry: &FieldEntry_DstIp{
			DstIp: NewDstIPFieldEntry(ethType, dstIP, inPort),
		},
	}
}

//
// NewIPProtoFieldEntry returns new instance
//
func NewIPProtoFieldEntry(ethType uint16, ipProto uint8, inPort opennsl.Port) *IpProtoFieldEntry {
	return &IpProtoFieldEntry{
		EthType: uint32(ethType),
		IpProto: uint32(ipProto),
		InPort:  uint32(inPort),
	}
}

//
// NewFieldEntryIPProto returns new instance
//
func NewFieldEntryIPProto(ethType uint16, ipProto uint8, inPort opennsl.Port) *FieldEntry {
	return &FieldEntry{
		EntryType: FieldEntry_IP_PROTO,
		Entry: &FieldEntry_IpProto{
			IpProto: NewIPProtoFieldEntry(ethType, ipProto, inPort),
		},
	}
}

//
// NewGetFieldEntriesRequest returns new instance
//
func NewGetFieldEntriesRequest() *GetFieldEntriesRequest {
	return &GetFieldEntriesRequest{}
}

//
// NewGetFieldEntriesReply returns new instance
//
func NewGetFieldEntriesReply(entries []*FieldEntry) *GetFieldEntriesReply {
	return &GetFieldEntriesReply{
		Entries: entries,
	}
}

//
// NewVlanEntry returns new instance
//
func NewVlanEntry(vid opennsl.Vlan, ports []opennsl.Port, utPorts []opennsl.Port) *VlanEntry {
	convPort := func(pList []opennsl.Port) []uint32 {
		list := make([]uint32, len(pList))
		for index, p := range pList {
			list[index] = uint32(p)
		}
		return list
	}

	return &VlanEntry{
		Vid:        uint32(vid),
		Ports:      convPort(ports),
		UntagPorts: convPort(utPorts),
	}
}

//
// NewVlanEntryFromNative creates new instance from native types.
//
func NewVlanEntryFromNative(vid opennsl.Vlan, pbmp *opennsl.PBmp, upbmp *opennsl.PBmp) *VlanEntry {
	ports := []opennsl.Port{}
	pbmp.Each(func(port opennsl.Port) error {
		ports = append(ports, port)
		return nil
	})

	utPorts := []opennsl.Port{}
	upbmp.Each(func(port opennsl.Port) error {
		utPorts = append(utPorts, port)
		return nil
	})

	return NewVlanEntry(vid, ports, utPorts)
}

//
// NewGetVlansRequest returns new instance.
//
func NewGetVlansRequest() *GetVlansRequest {
	return &GetVlansRequest{}
}

//
// NewGetVlansReply returns new instance.
//
func NewGetVlansReply(vlans []*VlanEntry) *GetVlansReply {
	return &GetVlansReply{
		Vlans: vlans,
	}
}

//
// NewL2Addr returns new instance.
//
func NewL2Addr(mac string, vid opennsl.Vlan) *L2Addr {
	return &L2Addr{
		Mac: mac,
		Vid: uint32(vid),
	}
}

//
// NewL2AddrFromNative creates new instance from native types.
//
func NewL2AddrFromNative(l2addr *opennsl.L2Addr) *L2Addr {
	return &L2Addr{
		Flags: uint32(l2addr.Flags()),
		Mac:   l2addr.MAC().String(),
		Vid:   uint32(l2addr.VID()),
		Port:  uint32(l2addr.Port()),
	}
}

//
// NewGetL2AddrsRequest returns new instance.
//
func NewGetL2AddrsRequest() *GetL2AddrsRequest {
	return &GetL2AddrsRequest{}
}

//
// NewGetL2AddrsReply returns new instance.
//
func NewGetL2AddrsReply(addrs []*L2Addr) *GetL2AddrsReply {
	return &GetL2AddrsReply{
		Addrs: addrs,
	}
}

//
// NewL2Station returns new instance.
//
func NewL2Station(mac string, macMask string, vlan opennsl.Vlan, vlanMask opennsl.Vlan) *L2Station {
	return &L2Station{
		DstMac:     mac,
		DstMacMask: macMask,
		Vlan:       uint32(vlan),
		VlanMask:   uint32(vlanMask),
	}
}

//
// NewL2StationFromNative creates new instance from native types.
//
func NewL2StationFromNative(l2st *opennsl.L2Station) *L2Station {
	return &L2Station{
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
// NewGetL2StationsRequest returns new instance.
//
func NewGetL2StationsRequest() *GetL2StationsRequest {
	return &GetL2StationsRequest{}
}

//
// NewGetL2StationsReply returns new instance.
//
func NewGetL2StationsReply(stations []*L2Station) *GetL2StationsReply {
	return &GetL2StationsReply{
		Stations: stations,
	}
}

//
// NewL3Iface returns new instance.
//
func NewL3Iface(mac string, vid opennsl.Vlan) *L3Iface {
	return &L3Iface{
		Mac: mac,
		Vid: uint32(vid),
	}
}

//
// NewL3IfaceFromNative creates new instance from native types.
//
func NewL3IfaceFromNative(l3iface *opennsl.L3Iface) *L3Iface {
	return &L3Iface{
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
// NewGetL3IfaceRequest returns new instance.
//
func NewGetL3IfaceRequest(ifaceID opennsl.L3IfaceID) *GetL3IfaceRequest {
	return &GetL3IfaceRequest{
		IfaceId: uint32(ifaceID),
	}
}

//
// NewGetL3IfaceReply returns new instance.
//
func NewGetL3IfaceReply(iface *L3Iface) *GetL3IfaceReply {
	return &GetL3IfaceReply{
		Iface: iface,
	}
}

//
// NewGetL3IfacesRequest returns new instance.
//
func NewGetL3IfacesRequest() *GetL3IfacesRequest {
	return &GetL3IfacesRequest{}
}

//
// NewGetL3IfacesReply returns new instance.
//
func NewGetL3IfacesReply(ifaces []*L3Iface) *GetL3IfacesReply {
	return &GetL3IfacesReply{
		Ifaces: ifaces,
	}
}

//
// NewFindL3IfaceRequest returns new instance.
//
func NewFindL3IfaceRequest(mac string, vid opennsl.Vlan) *FindL3IfaceRequest {
	return &FindL3IfaceRequest{
		Mac: mac,
		Vid: uint32(vid),
	}
}

//
// NewFindL3IfaceReply returns new instance.
//
func NewFindL3IfaceReply(iface *L3Iface) *FindL3IfaceReply {
	return &FindL3IfaceReply{
		Iface: iface,
	}
}

//
// NewL3Egress returns new instance.
//
func NewL3Egress(l3egrID opennsl.L3EgressID) *L3Egress {
	return &L3Egress{
		EgressId: uint32(l3egrID),
	}
}

//
// NewL3EgressFromNative creates and returns from native types.
//
func NewL3EgressFromNative(l3egrID opennsl.L3EgressID, l3egr *opennsl.L3Egress) *L3Egress {
	return &L3Egress{
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
// NewGetL3EgressesRequest returns new instance.
//
func NewGetL3EgressesRequest() *GetL3EgressesRequest {
	return &GetL3EgressesRequest{}
}

//
// NewGetL3EgressesReply returns new instance.
//
func NewGetL3EgressesReply(l3egresses []*L3Egress) *GetL3EgressesReply {
	return &GetL3EgressesReply{
		Egresses: l3egresses,
	}
}

//
// NewL3Host returns new instance.
//
func NewL3Host(ipAddr string) *L3Host {
	return &L3Host{
		IpAddr: ipAddr,
	}
}

//
// NewL3HostFromNative creates and returns from native types.
//
func NewL3HostFromNative(host *opennsl.L3Host) *L3Host {
	return &L3Host{
		Flags:    uint32(host.Flags()),
		EgressId: uint32(host.EgressID()),
		IpAddr:   host.IPAddr().String(),
		Ip6Addr:  host.IP6Addr().String(),
		Mac:      host.NexthopMAC().String(),
		Vrf:      uint32(host.VRF()),
	}
}

//
// NewGetL3HostsRequest returns new instance.
//
func NewGetL3HostsRequest() *GetL3HostsRequest {
	return &GetL3HostsRequest{}
}

//
// NewGetL3HostsReply returns new instance.
//
func NewGetL3HostsReply(hosts []*L3Host) *GetL3HostsReply {
	return &GetL3HostsReply{
		Hosts: hosts,
	}
}

//
// NewL3Route returns new instance.
//
func NewL3Route(ipAddr string) *L3Route {
	return &L3Route{
		IpAddr: ipAddr,
	}
}

//
// NewL3RouteFromNative creates and returns from native types.
//
func NewL3RouteFromNative(route *opennsl.L3Route) *L3Route {
	return &L3Route{
		Flags:    uint32(route.Flags()),
		EgressId: uint32(route.EgressID()),
		IpAddr:   route.IP4Net().String(),
		Ip6Addr:  route.IP6Net().String(),
		Vrf:      uint32(route.VRF()),
	}
}

//
// NewGetL3RoutesRequest returns new instance.
//
func NewGetL3RoutesRequest() *GetL3RoutesRequest {
	return &GetL3RoutesRequest{}
}

//
// NewGetL3RoutesReply returns new instance.
//
func NewGetL3RoutesReply(routes []*L3Route) *GetL3RoutesReply {
	return &GetL3RoutesReply{
		Routes: routes,
	}
}

//
// IDMapName is entry kind name.
//
type IDMapName string

const (
	// IDMapNameL3Egress is L3Egress entry name.
	IDMapNameL3Egress = "L3Egress"
	// IDMapNameL3Iface is L3Iface entry name.
	IDMapNameL3Iface = "L3Iface"
	// IDMapNameTrunk is Trunk entry name.
	IDMapNameTrunk = "Trunk"
)

//
// NewIDMapEntry returns new instance.
//
func NewIDMapEntry(name IDMapName, key string, value uint32) *IDMapEntry {
	return &IDMapEntry{
		Name:  string(name),
		Key:   key,
		Value: value,
	}
}

//
// NewGetIDMapEntriesRequest returns new instance.
//
func NewGetIDMapEntriesRequest() *GetIDMapEntriesRequest {
	return &GetIDMapEntriesRequest{}
}

//
// NewGetIDMapEntriesReply returns new instance.
//
func NewGetIDMapEntriesReply(entries []*IDMapEntry) *GetIDMapEntriesReply {
	return &GetIDMapEntriesReply{
		Entries: entries,
	}
}

func NewTunnelInitiatorFromNative(tunnel *opennsl.TunnelInitiator) *TunnelInitiator {
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

	return &TunnelInitiator{
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

//
// NewGetTunnelInitiatorsRequest returns new instance.
//
func NewGetTunnelInitiatorsRequest() *GetTunnelInitiatorsRequest {
	return &GetTunnelInitiatorsRequest{}
}

//
// NewGetTunnelInitiatorsReply returns new instance.
//
func NewGetTunnelInitiatorsReply(entries []*TunnelInitiator) *GetTunnelInitiatorsReply {
	return &GetTunnelInitiatorsReply{
		Tunnels: entries,
	}
}

func NewTunnelTerminatorFromNative(tunnel *opennsl.TunnelTerminator) *TunnelTerminator {
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

	return &TunnelTerminator{
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
// NewGetTunnelTerminatorsRequest returns new instance.
//
func NewGetTunnelTerminatorsRequest() *GetTunnelTerminatorsRequest {
	return &GetTunnelTerminatorsRequest{}
}

//
// NewGetTunnelTerminatorsReply returns new instance.
//
func NewGetTunnelTerminatorsReply(entries []*TunnelTerminator) *GetTunnelTerminatorsReply {
	return &GetTunnelTerminatorsReply{
		Tunnels: entries,
	}
}

//
// PortInfo
//
func NewPortInfoFromNative(port opennsl.Port, pinfo *opennsl.PortInfo) *PortInfo {
	return &PortInfo{
		Port:         uint32(port),
		LinkStatus:   int32(pinfo.LinkStatus()),
		UntaggedVlan: uint32(pinfo.UntaggedVlan()),
	}
}

func NewGetPortInfosReply(pinfos []*PortInfo) *GetPortInfosReply {
	return &GetPortInfosReply{
		PortInfos: pinfos,
	}
}

func NewGetPortInfosRequest() *GetPortInfosRequest {
	return &GetPortInfosRequest{}
}
