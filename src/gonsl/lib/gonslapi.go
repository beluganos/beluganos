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

package gonslib

import (
	api "gonsl/api"
	"net"

	"github.com/beluganos/go-opennsl/opennsl"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

const (
	l3HostTraverseMax  = 1024
	l3RouteTraverseMax = 1024
)

//
// APIServer is api server.
//
type APIServer struct {
	server *Server

	log *log.Entry
}

//
// NewAPIServer returns new instance.
//
func NewAPIServer(server *Server) *APIServer {
	return &APIServer{
		server: server,

		log: log.WithFields(log.Fields{"module": "apisrv"}),
	}
}

//
// Start starts sub modules.
//
func (s *APIServer) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	api.RegisterGoNSLApiServer(grpcServer, s)
	go grpcServer.Serve(lis)

	s.log.Infof("started.")
	return nil
}

//
// GetFieldEntries process api.GetFieldEntriesRequest.
//
func (s *APIServer) GetFieldEntries(ctxt context.Context, req *api.GetFieldEntriesRequest) (*api.GetFieldEntriesReply, error) {
	var entries []opennsl.FieldEntry
	var err error

	fields := s.server.Fields()
	result := []*api.FieldEntry{}

	if entries, err = fields.EthType.GetEntries(); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		e := FieldEntryEthType{}
		if err := fields.EthType.GetEntry(&e, entry); err != nil {
			s.log.Warnf("GetFieldEntries: GetEntry error. %d %s", entry, err)
		} else {
			result = append(result, NewFieldEntryEthTypeAPI(e.EthType, e.InPort))
		}
	}

	if entries, err = fields.DstIPv4.GetEntries(); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		e := FieldEntryDstIP{}
		if err := fields.DstIPv4.GetEntry(&e, entry); err != nil {
			s.log.Warnf("GetFieldEntries: GetEntry error. %d %s", entry, err)
		} else {
			result = append(result, NewFieldEntryDstIPAPI(uint16(e.EthType), e.Dest.String(), e.InPort))
		}
	}

	if entries, err = fields.DstIPv6.GetEntries(); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		e := FieldEntryDstIP{}
		if err := fields.DstIPv6.GetEntry(&e, entry); err != nil {
			s.log.Warnf("GetFieldEntries: GetEntry error. %d %s", entry, err)
		} else {
			result = append(result, NewFieldEntryDstIPAPI(uint16(e.EthType), e.Dest.String(), e.InPort))
		}
	}

	if entries, err = fields.IPProto.GetEntries(); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		e := FieldEntryIPProto{}
		if err := fields.IPProto.GetEntry(&e, entry); err != nil {
			s.log.Warnf("GetFieldEntries: GetEntry error. %d %s", entry, err)
		} else {
			result = append(result, NewFieldEntryIPProtoAPI(uint16(e.EthType), uint8(e.IPProto), e.InPort))
		}
	}

	return api.NewGetFieldEntriesReply(result), nil
}

//
// GetPortInfos
//

func (s *APIServer) GetPortInfos(ctxt context.Context, req *api.GetPortInfosRequest) (*api.GetPortInfosReply, error) {
	pbmp, err := PortBmpGet(s.server.Unit())
	if err != nil {
		return nil, err
	}

	portInfos := []*api.PortInfo{}
	err = pbmp.Each(func(port opennsl.Port) error {
		pinfo := opennsl.NewPortInfo()
		if err := pinfo.PortSelectiveGet(s.server.Unit(), port); err != nil {
			return err
		}

		portInfos = append(portInfos, NewPortInfoAPI(port, pinfo))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return api.NewGetPortInfosReply(portInfos), nil
}

//
// GetVlans process api.GetVlansRequest.
//
func (s *APIServer) GetVlans(ctxt context.Context, req *api.GetVlansRequest) (*api.GetVlansReply, error) {
	entries := []*api.VlanEntry{}
	err := opennsl.VlanTraverse(s.server.Unit(), func(unit int, vlan opennsl.Vlan, pbmp *opennsl.PBmp, ubmp *opennsl.PBmp) opennsl.OpenNSLError {
		entries = append(entries, NewVlanEntryAPI(vlan, pbmp, ubmp))
		return opennsl.E_NONE
	})

	if err != nil {
		return nil, err
	}

	return api.NewGetVlansReply(entries), nil
}

//
// GetL2Addrs process api.GetL2AddrsRequest.
//
func (s *APIServer) GetL2Addrs(ctxt context.Context, req *api.GetL2AddrsRequest) (*api.GetL2AddrsReply, error) {
	l2addrs := []*api.L2Addr{}
	err := opennsl.L2Traverse(s.server.Unit(), func(unit int, l2addr *opennsl.L2Addr) opennsl.OpenNSLError {
		l2addrs = append(l2addrs, NewL2AddrAPI(l2addr))
		return opennsl.E_NONE
	})

	if err != nil {
		return nil, err
	}

	return api.NewGetL2AddrsReply(l2addrs), nil
}

//
// GetL2Stations process api.GetL2StationsRequest.
//
func (s *APIServer) GetL2Stations(ctxt context.Context, req *api.GetL2StationsRequest) (*api.GetL2StationsReply, error) {
	l2stations := []*api.L2Station{}
	s.server.idmaps.L2Stations.Traverse(func(key L2StationIDKey, l2stationId opennsl.L2StationID) bool {
		l2station, err := l2stationId.Get(s.server.Unit())
		if err != nil {
			s.log.Errorf("L2StationGet error. %d %s", l2stationId, err)
		} else {
			l2stations = append(l2stations, NewL2StationAPI(l2station))
		}
		return true
	})

	return api.NewGetL2StationsReply(l2stations), nil
}

//
// FindL3Iface process api.FindL3IfaceRequest.
//
func (s *APIServer) FindL3Iface(ctxt context.Context, req *api.FindL3IfaceRequest) (*api.FindL3IfaceReply, error) {
	mac, err := net.ParseMAC(req.Mac)
	if err != nil {
		return nil, err
	}

	l3iface, err := opennsl.L3IfaceFind(s.server.Unit(), mac, opennsl.Vlan(req.Vid))
	if err != nil {
		return nil, err
	}

	return api.NewFindL3IfaceReply(NewL3IfaceAPI(l3iface)), nil
}

//
// GetL3Iface process api.GetL3IfaceRequest.
//
func (s *APIServer) GetL3Iface(ctxt context.Context, req *api.GetL3IfaceRequest) (*api.GetL3IfaceReply, error) {
	l3iface, err := opennsl.L3IfaceGet(s.server.Unit(), opennsl.L3IfaceID(req.IfaceId))
	if err != nil {
		return nil, err
	}

	return api.NewGetL3IfaceReply(NewL3IfaceAPI(l3iface)), nil
}

//
// GetL3Ifaces process api.GetL3IfacesRequest.
//
func (s *APIServer) GetL3Ifaces(ctxt context.Context, req *api.GetL3IfacesRequest) (*api.GetL3IfacesReply, error) {

	l3Ifaces := []*api.L3Iface{}
	s.server.idmaps.L3Ifaces.Traverse(func(key L3IfaceIDKey, ifaceId opennsl.L3IfaceID) bool {
		l3iface, err := opennsl.L3IfaceGet(s.server.Unit(), ifaceId)
		if err != nil {
			s.log.Errorf("GetL3Ifaces: L3IfaceGet error. %d %s", ifaceId, err)
		} else {
			l3Ifaces = append(l3Ifaces, NewL3IfaceAPI(l3iface))
		}

		return true
	})

	return api.NewGetL3IfacesReply(l3Ifaces), nil
}

//
// GetL3Egresses process api.GetL3EgressesRequest.
//
func (s *APIServer) GetL3Egresses(ctxt context.Context, req *api.GetL3EgressesRequest) (*api.GetL3EgressesReply, error) {
	l3egrs := []*api.L3Egress{}
	err := opennsl.L3EgressTraverse(s.server.Unit(), func(unit int, l3egrId opennsl.L3EgressID, l3egr *opennsl.L3Egress) opennsl.OpenNSLError {
		l3egrs = append(l3egrs, NewL3EgressAPI(l3egrId, l3egr))
		return opennsl.E_NONE
	})

	if err != nil {
		return nil, err
	}

	return api.NewGetL3EgressesReply(l3egrs), nil
}

//
// GetL3Hosts process api.GetL3HostsRequest.
//
func (s *APIServer) GetL3Hosts(ctxt context.Context, req *api.GetL3HostsRequest) (*api.GetL3HostsReply, error) {
	hosts := []*api.L3Host{}
	err := opennsl.L3HostTraverse(s.server.Unit(), 0, 0, l3HostTraverseMax, func(unit int, index int, host *opennsl.L3Host) opennsl.OpenNSLError {
		hosts = append(hosts, NewL3HostAPI(host))
		return opennsl.E_NONE
	})

	if err != nil {
		return nil, err
	}

	err = opennsl.L3HostTraverse(s.server.Unit(), uint32(opennsl.L3_IP6), 0, l3HostTraverseMax, func(unit int, index int, host *opennsl.L3Host) opennsl.OpenNSLError {
		hosts = append(hosts, NewL3HostAPI(host))
		return opennsl.E_NONE
	})

	if err != nil {
		return nil, err
	}

	return api.NewGetL3HostsReply(hosts), nil
}

//
// GetL3Routes process api.GetL3RoutesRequest.
//
func (s *APIServer) GetL3Routes(ctxt context.Context, req *api.GetL3RoutesRequest) (*api.GetL3RoutesReply, error) {
	routes := []*api.L3Route{}
	err := opennsl.L3RouteTraverse(s.server.Unit(), 0, 0, l3RouteTraverseMax, func(unit int, index int, route *opennsl.L3Route) opennsl.OpenNSLError {
		routes = append(routes, NewL3RouteAPI(route))
		return opennsl.E_NONE
	})

	if err != nil {
		return nil, err
	}

	err = opennsl.L3RouteTraverse(s.server.Unit(), uint32(opennsl.L3_IP6), 0, l3RouteTraverseMax, func(unit int, index int, route *opennsl.L3Route) opennsl.OpenNSLError {
		routes = append(routes, NewL3RouteAPI(route))
		return opennsl.E_NONE
	})

	if err != nil {
		return nil, err
	}

	return api.NewGetL3RoutesReply(routes), nil
}

//
// GetIDMapEntries process api.GetIDMapEntriesRequest.
//
func (s *APIServer) GetIDMapEntries(ctxt context.Context, req *api.GetIDMapEntriesRequest) (*api.GetIDMapEntriesReply, error) {
	entries := []*api.IDMapEntry{}
	s.server.idmaps.L3Ifaces.Traverse(func(key L3IfaceIDKey, value opennsl.L3IfaceID) bool {
		entry := api.NewIDMapEntry(api.IDMapNameL3Iface, key.String(), uint32(value))
		entries = append(entries, entry)
		return true
	})

	s.server.idmaps.L3Egress.Traverse(func(key L3EgressIDKey, value opennsl.L3EgressID) bool {
		entry := api.NewIDMapEntry(api.IDMapNameL3Egress, key.String(), uint32(value))
		entries = append(entries, entry)
		return true
	})

	s.server.idmaps.Trunks.Traverse(func(key TrunkIDKey, value opennsl.Trunk) bool {
		entry := api.NewIDMapEntry(api.IDMapNameTrunk, key.String(), uint32(value))
		entries = append(entries, entry)
		return true
	})

	return api.NewGetIDMapEntriesReply(entries), nil
}

//
// GetTunnelInitiators process api.GetTunnelInitiatorsRequest.
//
func (s *APIServer) GetTunnelInitiators(ctxt context.Context, req *api.GetTunnelInitiatorsRequest) (*api.GetTunnelInitiatorsReply, error) {
	tunnels := []*api.TunnelInitiator{}
	opennsl.TunnelInitiatorTraverse(s.server.Unit(), func(unit int, initiator *opennsl.TunnelInitiator) opennsl.OpenNSLError {
		tunnels = append(tunnels, NewTunnelInitiatorAPI(initiator))
		return opennsl.E_NONE
	})
	return api.NewGetTunnelInitiatorsReply(tunnels), nil
}

//
// GetTunnelTerminators process api.GetTunnelTerminatorsRequest.
//
func (s *APIServer) GetTunnelTerminators(ctxt context.Context, req *api.GetTunnelTerminatorsRequest) (*api.GetTunnelTerminatorsReply, error) {
	tunnels := []*api.TunnelTerminator{}
	opennsl.TunnelTerminatorTraverse(s.server.Unit(), func(unit int, terminator *opennsl.TunnelTerminator) opennsl.OpenNSLError {
		tunnels = append(tunnels, NewTunnelTerminatorAPI(terminator))
		return opennsl.E_NONE
	})
	return api.NewGetTunnelTerminatorsReply(tunnels), nil
}
