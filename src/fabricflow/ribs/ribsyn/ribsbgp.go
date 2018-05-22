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

package ribsyn

import (
	"fabricflow/ribs/ribsmsg"
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gonla/nlalib"
	"google.golang.org/grpc"
)

func GetLabelsFromPath(p *table.Path) []uint32 {
	nlri := p.GetNlri()
	switch nlri.(type) {
	case *bgp.LabeledIPAddrPrefix:
		return nlri.(*bgp.LabeledIPAddrPrefix).Labels.Labels
	case *bgp.LabeledIPv6AddrPrefix:
		return nlri.(*bgp.LabeledIPv6AddrPrefix).Labels.Labels
	case *bgp.LabeledVPNIPAddrPrefix:
		return nlri.(*bgp.LabeledVPNIPAddrPrefix).Labels.Labels
	case *bgp.LabeledVPNIPv6AddrPrefix:
		return nlri.(*bgp.LabeledVPNIPv6AddrPrefix).Labels.Labels
	default:
		return []uint32{}
	}
}

func GetBgpPathAttribute(t bgp.BGPAttrType, p *table.Path) (bgp.PathAttributeInterface, bool) {
	for _, attr := range p.GetPathAttrs() {
		if attr.GetType() == t {
			return attr, true
		}
	}

	return nil, false
}

func GetBgpExtCommunityRouteTarget(path *table.Path) bgp.ExtendedCommunityInterface {
	if pattr, ok := GetBgpPathAttribute(bgp.BGP_ATTR_TYPE_EXTENDED_COMMUNITIES, path); ok {
		for _, extcom := range pattr.(*bgp.PathAttributeExtendedCommunities).Value {
			if _, subType := extcom.GetTypes(); subType == bgp.EC_SUBTYPE_ROUTE_TARGET {
				return extcom
			}
		}
	}

	return nil
}

func NewBgpExtCommunityRouteTarget(rt string) (bgp.PathAttributeInterface, error) {
	rtcomm, err := bgp.ParseRouteTarget(rt)
	if err != nil {
		return nil, err
	}
	excomms := []bgp.ExtendedCommunityInterface{rtcomm}
	return bgp.NewPathAttributeExtendedCommunities(excomms), nil
}

func GetBgpMpReachNlri(p *table.Path) *bgp.PathAttributeMpReachNLRI {
	if attr, ok := GetBgpPathAttribute(bgp.BGP_ATTR_TYPE_MP_REACH_NLRI, p); ok {
		return attr.(*bgp.PathAttributeMpReachNLRI)
	}

	return nil
}

func NewBgpMpReachNlri(rd, prefix string, prefixlen uint8, label uint32) (bgp.AddrPrefixInterface, error) {
	pattr, err := bgp.ParseRouteDistinguisher(rd)
	if err != nil {
		return nil, err
	}

	mpls := bgp.NewMPLSLabelStack(label)
	return bgp.NewLabeledVPNIPAddrPrefix(prefixlen, prefix, *mpls, pattr), nil
}

func GetBgpNlriPrefixList(nlri *bgp.PathAttributeMpReachNLRI) []*bgp.LabeledVPNIPAddrPrefix {
	ps := []*bgp.LabeledVPNIPAddrPrefix{}
	for _, v := range nlri.Value {
		if p, ok := v.(*bgp.LabeledVPNIPAddrPrefix); ok {
			ps = append(ps, p)
		}
	}
	return ps
}

func GetBgpAddrPrefixFromNlri(nlri bgp.AddrPrefixInterface) (string, uint8, bool) {
	if prefix, ok := nlri.(*bgp.IPAddrPrefix); ok {
		return prefix.Prefix.String(), prefix.Length, true
	} else {
		return "", 0, false
	}
}

func NewIPv4NlriFromVPNv4(prefix *bgp.LabeledVPNIPAddrPrefix) *bgp.IPAddrPrefix {
	return bgp.NewIPAddrPrefix(prefix.Length-88, prefix.Prefix.String())
}

type BgpClient struct {
	Api  api.GobgpApiClient
	Conn *grpc.ClientConn
}

func NewBgpClient(conCh chan<- *nlalib.ConnInfo, addr string) (*BgpClient, error) {
	conn, err := nlalib.NewClientConn(addr, conCh)
	if err != nil {
		return nil, err
	}

	return &BgpClient{
		Api:  api.NewGobgpApiClient(conn),
		Conn: conn,
	}, nil
}

func NewBgpMonitorRibReq(rt string) *api.MonitorRibRequest {
	family := bgp.RF_IPv4_UC
	if rt == "VPNv4" {
		family = bgp.RF_IPv4_VPN
	}

	return &api.MonitorRibRequest{
		Table: &api.Table{
			Type:   api.Resource_GLOBAL,
			Family: uint32(family),
		},
		Current: false,
	}
}

func NewBgpMonitorRib(ribCh chan<- *ribsmsg.RibUpdate, client *BgpClient, rt string) error {
	stream, err := client.Api.MonitorRib(context.Background(), NewBgpMonitorRibReq(rt))
	if err != nil {
		log.Errorf("BGPC Monitor: api error. %s", err)
		return err
	}

	go func() {
		log.Infof("BGPC Monitor: %s", rt)
		for {
			dest, err := stream.Recv()
			if err != nil {
				break
			}

			rib, err := ribsmsg.NewRibUpdateFromDestination(dest, rt)
			if err != nil {
				log.Errorf("BGPC Monitor: NewRibUpdate error. %s", err)
				continue
			}

			ribCh <- rib
		}
		log.Infof("BGPC Monitor: EXIT. %s", rt)
	}()

	return nil
}

func (c *BgpClient) AddBgpPath(path *table.Path) error {
	req := &api.AddPathRequest{
		Resource: api.Resource_GLOBAL,
		VrfId:    "",
		Path:     api.ToPathApi(path, nil),
	}
	if _, err := c.Api.AddPath(context.Background(), req); err != nil {
		return err
	}

	return nil
}

func (c *BgpClient) DelBgpPath(path *table.Path) error {
	req := &api.DeletePathRequest{
		Resource: api.Resource_GLOBAL,
		VrfId:    "",
		Path:     api.ToPathApi(path, nil),
	}
	if _, err := c.Api.DeletePath(context.Background(), req); err != nil {
		return err
	}

	return nil
}

func (c *BgpClient) GetRib(rt string, f func(*ribsmsg.RibUpdate) error) error {
	family := bgp.RF_IPv4_UC
	if rt == "VPNv4" {
		family = bgp.RF_IPv4_VPN
	}
	req := api.GetRibRequest{
		Table: &api.Table{
			Type:         api.Resource_GLOBAL,
			Family:       uint32(family),
			Name:         "",
			Destinations: []*api.Destination{},
		},
	}

	res, err := c.Api.GetRib(context.Background(), &req)
	if err != nil {
		return err
	}

	for _, d := range res.Table.Destinations {
		rib, err := ribsmsg.NewRibUpdateFromDestination(d, rt)
		if err != nil {
			log.Errorf("BGPC GetRib: Invalid Destination. %v", d)
			continue
		}

		if err := f(rib); err != nil {
			return err
		}
	}

	return nil
}
