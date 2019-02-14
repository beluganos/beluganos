// -*- coding: utf-8 -*-
// +build gobgpv2

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

package apiutil

import (
	"net"
	"time"

	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/pkg/packet/bgp"
)

func NewNativePath(p *api.Path) (*Path, error) {
	nlri, _ := GetNativeNlri(p)
	attrs, _ := GetNativePathAttributes(p)
	return &Path{
		Nlri:       nlri,
		Age:        p.Age,
		Best:       p.Best,
		Attrs:      attrs,
		Stale:      p.Stale,
		Withdrawal: p.IsWithdraw,
		SourceID:   net.ParseIP(p.SourceId),
		NeighborIP: net.ParseIP(p.NeighborIp),
	}, nil
}

func NewApiPath(p *Path) *api.Path {
	return NewPath(
		p.Nlri,
		p.Withdrawal,
		p.Attrs,
		time.Unix(p.Age, 0),
	)
}

func (path *Path) GetPathAttrs() []bgp.PathAttributeInterface {
	return path.Attrs
}

func (path *Path) getPathAttr(typ bgp.BGPAttrType) (bgp.PathAttributeInterface, bool) {
	for _, attr := range path.GetPathAttrs() {
		if attr.GetType() == typ {
			return attr, true
		}
	}
	return nil, false
}

func GetPathAttribute(path *Path, typ bgp.BGPAttrType) (bgp.PathAttributeInterface, bool) {
	return path.getPathAttr(typ)
}

func (path *Path) GetNexthop() net.IP {
	if attr, ok := path.getPathAttr(bgp.BGP_ATTR_TYPE_NEXT_HOP); ok {
		return attr.(*bgp.PathAttributeNextHop).Value
	}
	return net.IP{}
}

func (path *Path) GetNlri() bgp.AddrPrefixInterface {
	return path.Nlri
}
