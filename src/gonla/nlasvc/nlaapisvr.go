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

package nlasvc

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gonla/nlaapi"
	"gonla/nladbm"
	"gonla/nlalib"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"google.golang.org/grpc"
	"net"
)

//
// NLA API Server
//
type NLAApiServer struct {
	nid    uint8
	addr   string
	NlMsgs chan<- *nlamsg.NetlinkMessageUnion
}

func NewNLAApiServer(addr string, nid uint8) *NLAApiServer {
	return &NLAApiServer{
		nid:    nid,
		addr:   addr,
		NlMsgs: nil,
	}
}

func (n *NLAApiServer) Start(ch chan<- *nlamsg.NetlinkMessageUnion) error {
	listen, err := net.Listen("tcp", n.addr)
	if err != nil {
		return err
	}

	n.NlMsgs = ch

	s := grpc.NewServer()
	nlaapi.RegisterNLAApiServer(s, n)
	go s.Serve(listen)

	log.Infof("NLAApiServer: START")
	return nil
}

func (n *NLAApiServer) MonNetlink(req *nlaapi.MonNetlinkRequest, stream nlaapi.NLAApi_MonNetlinkServer) error {
	log.Infof("NLAApiServer: Monitor START. %v", req)

	client := nladbm.Clients().New()
	defer nladbm.Clients().Delete(client)

	for {
		select {
		case <-stream.Context().Done():
			log.Infof("NLAApiServer: Monitor EXIT. Closed")
			return nil

		case m := <-client:
			res := nlaapi.NewNetlinkMessageUnionFromNative(m)
			if err := stream.Send(res); err != nil {
				log.Infof("NLAApiServer: Monitor EXIT. Stream error. %s", err)
				return err
			}
		}
	}

	// return nil
}

func (n *NLAApiServer) ModVpn(ctxt context.Context, req *nlaapi.ModVpnRequest) (*nlaapi.ModVpnReply, error) {
	vpn := req.Vpn.ToNative()
	hdr := nlalib.NewNlMsghdr(uint16(req.Type), 0)
	n.NlMsgs <- nlamsg.NewNetlinkMessageUnion(&hdr, vpn, vpn.NId, nlamsg.SRC_API)
	return &nlaapi.ModVpnReply{}, nil
}

func (n *NLAApiServer) ModNetlink(ctxt context.Context, req *nlaapi.NetlinkMessageUnion) (*nlaapi.ModNetlinkReply, error) {
	nlmsg := req.ToNative()
	nlmsg.Src = nlamsg.SRC_API

	if nlmsg.Group() == nlalink.RTMGRP_VPN || nlmsg.NId == n.nid {
		// Dispatch to self.
		log.Debugf("NLAApiServer: ModNetlink to self. %s", nlmsg)
		n.NlMsgs <- nlmsg

	} else {
		// Send to slave
		node := nladbm.Nodes().Select(nladbm.NewNodeKey(nlmsg.NId))
		if node == nil {
			log.Warnf("NLAApiServer: ModNetlink Node not found. nid=%d", nlmsg.NId)
			return nil, fmt.Errorf("Node not found. nid=%d", nlmsg.NId)
		}

		log.Debugf("NLAApiServer: Send to slave. %s", nlmsg)
		if err := node.Send(nlmsg); err != nil {
			log.Errorf("NLAApiServer: ModNetlink send error. nid=%d %s", node.NId, err)
			return nil, fmt.Errorf("Failed to send request. nid=%d", node.NId)
		}
	}

	return &nlaapi.ModNetlinkReply{}, nil
}

func (n *NLAApiServer) GetLink(ctxt context.Context, req *nlaapi.LinkKey) (*nlaapi.Link, error) {
	if link := nladbm.Links().Select(req.ToNative()); link != nil {
		return nlaapi.NewLinkFromNative(link), nil
	}

	return nil, fmt.Errorf("Link not found. %v", req)
}

func (n *NLAApiServer) GetAddr(ctxt context.Context, req *nlaapi.AddrKey) (*nlaapi.Addr, error) {
	if addr := nladbm.Addrs().Select(req.ToNative()); addr != nil {
		return nlaapi.NewAddrFromNative(addr), nil
	}

	return nil, fmt.Errorf("Addr not found. %v", req)
}

func (n *NLAApiServer) GetNeigh(ctxt context.Context, req *nlaapi.NeighKey) (*nlaapi.Neigh, error) {
	if neigh := nladbm.Neighs().Select(req.ToNative()); neigh != nil {
		return nlaapi.NewNeighFromNative(neigh), nil
	}

	return nil, fmt.Errorf("Neigh not found. %v", req)
}

func (n *NLAApiServer) GetRoute(ctxt context.Context, req *nlaapi.RouteKey) (*nlaapi.Route, error) {
	if route := nladbm.Routes().Select(req.ToNative()); route != nil {
		return nlaapi.NewRouteFromNative(route), nil
	}

	return nil, fmt.Errorf("Route not found. %v", req)
}

func (n *NLAApiServer) GetMpls(ctxt context.Context, req *nlaapi.MplsKey) (*nlaapi.Route, error) {
	if mpls := nladbm.Mplss().Select(req.ToNative()); mpls != nil {
		return nlaapi.NewRouteFromNative(mpls), nil
	}

	return nil, fmt.Errorf("Mpls not found. %v", req)
}

func (n *NLAApiServer) GetNode(ctxt context.Context, req *nlaapi.NodeKey) (*nlaapi.Node, error) {
	if node := nladbm.Nodes().Select(req.ToNative()); node != nil {
		return nlaapi.NewNodeFromNative(node), nil
	}

	return nil, fmt.Errorf("Node not found. %v", req)
}

func (n *NLAApiServer) GetVpn(ctxt context.Context, req *nlaapi.VpnKey) (*nlaapi.Vpn, error) {
	if vpn := nladbm.Vpns().Select(req.ToNative()); vpn != nil {
		return nlaapi.NewVpnFromNative(vpn), nil
	}

	return nil, fmt.Errorf("Vpn not found. %v", req)
}

func (n *NLAApiServer) GetEncapInfo(ctxt context.Context, req *nlaapi.EncapInfoKey) (*nlaapi.EncapInfo, error) {
	if e := nladbm.Encaps().Select(req.ToNative()); e != nil {
		return nlaapi.NewEncapInfoFromNative(e), nil
	}

	return nil, fmt.Errorf("EncapInfo not found. %v", req)
}

func (n *NLAApiServer) GetLinks(req *nlaapi.GetLinksRequest, stream nlaapi.NLAApi_GetLinksServer) error {
	nladbm.Links().Walk(func(link *nlamsg.Link) error {
		m := nlaapi.NewLinkFromNative(link)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetAddrs(req *nlaapi.GetAddrsRequest, stream nlaapi.NLAApi_GetAddrsServer) error {
	nladbm.Addrs().Walk(func(addr *nlamsg.Addr) error {
		m := nlaapi.NewAddrFromNative(addr)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetNeighs(req *nlaapi.GetNeighsRequest, stream nlaapi.NLAApi_GetNeighsServer) error {
	nladbm.Neighs().Walk(func(neigh *nlamsg.Neigh) error {
		m := nlaapi.NewNeighFromNative(neigh)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetRoutes(req *nlaapi.GetRoutesRequest, stream nlaapi.NLAApi_GetRoutesServer) error {
	nladbm.Routes().Walk(func(route *nlamsg.Route) error {
		m := nlaapi.NewRouteFromNative(route)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetMplss(req *nlaapi.GetMplssRequest, stream nlaapi.NLAApi_GetMplssServer) error {
	nladbm.Mplss().Walk(func(mpls *nlamsg.Route) error {
		m := nlaapi.NewRouteFromNative(mpls)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetNodes(req *nlaapi.GetNodesRequest, stream nlaapi.NLAApi_GetNodesServer) error {
	nladbm.Nodes().Walk(func(node *nlamsg.Node) error {
		m := nlaapi.NewNodeFromNative(node)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetVpns(req *nlaapi.GetVpnsRequest, stream nlaapi.NLAApi_GetVpnsServer) error {
	nladbm.Vpns().Walk(func(vpn *nlamsg.Vpn) error {
		m := nlaapi.NewVpnFromNative(vpn)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetStats(req *nlaapi.GetStatsRequest, stream nlaapi.NLAApi_GetStatsServer) error {
	nladbm.Stats().Walk(func(stat nladbm.Stat) error {
		m := nlaapi.NewDbStatFromNative(stat)
		return stream.Send(m)
	})
	return nil
}

func (n *NLAApiServer) GetEncapInfos(req *nlaapi.GetEncapInfosRequest, stream nlaapi.NLAApi_GetEncapInfosServer) error {
	nladbm.Encaps().Walk(func(e *nlamsg.EncapInfo) error {
		m := nlaapi.NewEncapInfoFromNative(e)
		return stream.Send(m)
	})
	return nil
}
