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

package main

import (
	"context"
	"fabricflow/util/gobgp/apiutil"
	"fmt"
	"net"

	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/table"
)

func AddBgpPath(client api.GobgpApiClient, path *apiutil.Path) error {
	req := &api.AddPathRequest{
		Resource: api.Resource_GLOBAL,
		VrfId:    "",
		Path:     api.ToPathApi(path.Path, nil),
	}

	_, err := client.AddPath(context.Background(), req)
	return err
}

func DelBgpPath(client api.GobgpApiClient, path *apiutil.Path) error {
	req := &api.DeletePathRequest{
		Resource: api.Resource_GLOBAL,
		VrfId:    "",
		Path:     api.ToPathApi(path.Path, nil),
	}

	_, err := client.DeletePath(context.Background(), req)
	return err
}

func NewTunnelRouteFromPath(path *apiutil.Path) (*TunnelRoute, error) {
	route := TunnelRoute{}

	switch prefix := path.GetNlri().(type) {
	case *bgp.IPAddrPrefix:
		_, route.Prefix, _ = net.ParseCIDR(prefix.String())
		route.Nexthop = path.GetNexthop()
		route.Family = prefix.AFI()

	case *bgp.IPv6AddrPrefix:
		_, route.Prefix, _ = net.ParseCIDR(prefix.String())
		route.Nexthop = path.GetNexthop()
		route.Family = prefix.AFI()

	default:
		return nil, fmt.Errorf("Unsupported Prefix. %s", prefix)
	}

	if pattr, ok := apiutil.GetPathAttribute(path, bgp.BGP_ATTR_TYPE_EXTENDED_COMMUNITIES); ok {
		extcoms := pattr.(*bgp.PathAttributeExtendedCommunities)
		if extcom, ok := apiutil.GetExtendedCommunity(extcoms, bgp.EC_SUBTYPE_ENCAPSULATION); ok {
			route.TunnelType = extcom.(*bgp.EncapExtended).TunnelType
			return &route, nil
		}
	}

	return &route, nil
}

func NewExportRouteFromPath(path *apiutil.Path) *apiutil.Path {

	attrs := []bgp.PathAttributeInterface{}
	for _, attr := range path.GetPathAttrs() {
		switch a := attr.(type) {
		case *bgp.PathAttributeCommunities:
			attr = apiutil.FilterCommunities(a, func(comm uint32) bool {
				return bgp.WellKnownCommunity(comm) != bgp.COMMUNITY_NO_EXPORT
			})

		case *bgp.PathAttributeExtendedCommunities:
			attr = apiutil.FilterExtendedCommunity(a, func(attrType bgp.ExtendedCommunityAttrType, subType bgp.ExtendedCommunityAttrSubType) bool {
				return subType != bgp.EC_SUBTYPE_ENCAPSULATION
			})
		}

		if attr != nil {
			attrs = append(attrs, attr)
		}
	}

	newPath := table.NewPath(nil, path.GetNlri(), path.IsWithdraw, attrs, path.GetTimestamp(), false)

	switch path.GetNlri().(type) {
	case *bgp.IPAddrPrefix:
		newPath.SetNexthop(net.ParseIP("0.0.0.0"))

	case *bgp.IPv6AddrPrefix:
		newPath.SetNexthop(net.ParseIP("::"))
	}

	return &apiutil.Path{Path: newPath}
}
