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
// NewGetL3IfaceRequest returns new instance.
//
func NewGetL3IfaceRequest(ifaceID uint32) *GetL3IfaceRequest {
	return &GetL3IfaceRequest{
		IfaceId: ifaceID,
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
func NewFindL3IfaceRequest(mac string, vid uint16) *FindL3IfaceRequest {
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
func NewGetPortInfosReply(pinfos []*PortInfo) *GetPortInfosReply {
	return &GetPortInfosReply{
		PortInfos: pinfos,
	}
}

func NewGetPortInfosRequest() *GetPortInfosRequest {
	return &GetPortInfosRequest{}
}
