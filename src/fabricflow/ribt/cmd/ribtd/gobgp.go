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
	"github.com/osrg/gobgp/pkg/packet/bgp"
)

func AddBgpPath(client api.GobgpApiClient, path *apiutil.Path) error {
	req := api.AddPathRequest{
		TableType: api.TableType_GLOBAL,
		VrfId:     "",
		Path:      path.NewAPIPath(),
	}

	_, err := client.AddPath(context.Background(), &req)
	return err
}

func DelBgpPath(client api.GobgpApiClient, path *apiutil.Path) error {
	req := api.DeletePathRequest{
		TableType: api.TableType_GLOBAL,
		VrfId:     "",
		Path:      path.NewAPIPath(),
	}

	_, err := client.DeletePath(context.Background(), &req)
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

	if extcom, ok := path.GetExtCommunityPathAttr(bgp.EC_SUBTYPE_ENCAPSULATION); ok {
		route.TunnelType = extcom.(*bgp.EncapExtended).TunnelType
		return &route, nil
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

		case *bgp.PathAttributeMpReachNLRI:
			attr = nil

		case *bgp.PathAttributeMpUnreachNLRI:
			attr = nil
		}

		if attr != nil {
			attrs = append(attrs, attr)
		}
	}

	switch path.GetNlri().(type) {
	case *bgp.IPAddrPrefix:
		attrs = apiutil.AddNativePathAttributeNexthop(attrs, net.ParseIP("0.0.0.0"))

	case *bgp.IPv6AddrPrefix:
		attrs = apiutil.AddNativePathAttributeNexthop(attrs, net.ParseIP("::"))
	}

	path.Attrs = attrs

	return path
}
