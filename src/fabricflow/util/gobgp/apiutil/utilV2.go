// -*- coding: utf-8 -*-
// +build !gobgpv1

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

	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/pkg/packet/bgp"
)

type Path struct {
	Nlri       bgp.AddrPrefixInterface
	Attrs      []bgp.PathAttributeInterface
	SourceID   string
	NeighborIP string
	IsWithdraw bool
	Family     *api.Family

	apiPath *api.Path
}

func NewNativePath(p *api.Path) *Path {
	nlri, _ := GetNativeNlri(p)
	attrs, _ := GetNativePathAttributes(p)
	return &Path{
		Nlri:       nlri,
		Attrs:      attrs,
		SourceID:   p.SourceId,
		NeighborIP: p.NeighborIp,
		IsWithdraw: p.IsWithdraw,
		Family:     p.Family,

		apiPath: p,
	}
}

func (p *Path) NewAPIPath() *api.Path {
	p.apiPath.Nlri = MarshalNLRI(p.Nlri)
	p.apiPath.Pattrs = MarshalPathAttributes(p.Attrs)
	p.apiPath.SourceId = p.SourceID
	p.apiPath.NeighborIp = p.NeighborIP
	p.apiPath.IsWithdraw = p.IsWithdraw
	p.apiPath.Family = p.Family

	return p.apiPath
}

func (path *Path) GetNlri() bgp.AddrPrefixInterface {
	return path.Nlri
}

func (path *Path) GetPathAttrs() []bgp.PathAttributeInterface {
	return path.Attrs
}

func (path *Path) GetPathAttr(typ bgp.BGPAttrType) (bgp.PathAttributeInterface, bool) {
	return GetNativePathAttribute(typ, path.GetPathAttrs())
}

func (path *Path) SetPathAttr(pattr bgp.PathAttributeInterface) {
	path.Attrs = SetNativePathAttribute(path.GetPathAttrs(), pattr)
}

func (path *Path) GetExtCommunityPathAttr(subType bgp.ExtendedCommunityAttrSubType) (bgp.ExtendedCommunityInterface, bool) {
	return GetNativeExtCommunityAttribute(subType, path.GetPathAttrs())
}

func (path *Path) GetNexthop() net.IP {
	if nh, ok := GetNexthopIPFromNativePathAttributes(path.GetPathAttrs()); ok {
		return nh
	}
	return net.IP{}
}

func getNLRI(family bgp.RouteFamily, buf []byte) (bgp.AddrPrefixInterface, error) {
	afi, safi := bgp.RouteFamilyToAfiSafi(family)
	nlri, err := bgp.NewPrefixFromRouteFamily(afi, safi)
	if err != nil {
		return nil, err
	}
	if err := nlri.DecodeFromBytes(buf); err != nil {
		return nil, err
	}
	return nlri, nil
}

func GetNativeNlri(p *api.Path) (bgp.AddrPrefixInterface, error) {
	if len(p.NlriBinary) > 0 {
		return getNLRI(ToRouteFamily(p.Family), p.NlriBinary)
	}
	return UnmarshalNLRI(ToRouteFamily(p.Family), p.Nlri)
}

func GetNativePathAttribute(attrType bgp.BGPAttrType, attrs []bgp.PathAttributeInterface) (bgp.PathAttributeInterface, bool) {
	for _, attr := range attrs {
		if attr.GetType() == attrType {
			return attr, true
		}
	}
	return nil, false
}

func SetNativePathAttribute(attrs []bgp.PathAttributeInterface, newAttr bgp.PathAttributeInterface) []bgp.PathAttributeInterface {
	for index, attr := range attrs {
		if attr.GetType() == newAttr.GetType() {
			attrs[index] = newAttr
			return attrs
		}
	}

	return append(attrs, newAttr)
}

func GetNativePathAttributes(p *api.Path) ([]bgp.PathAttributeInterface, error) {
	pattrsLen := 0
	if p.PattrsBinary != nil {
		pattrsLen = len(p.PattrsBinary)
	}

	if pattrsLen > 0 {
		pattrs := make([]bgp.PathAttributeInterface, 0, pattrsLen)
		for _, attr := range p.PattrsBinary {
			a, err := bgp.GetPathAttribute(attr)
			if err != nil {
				return nil, err
			}
			err = a.DecodeFromBytes(attr)
			if err != nil {
				return nil, err
			}
			pattrs = append(pattrs, a)
		}
		return pattrs, nil
	}
	return UnmarshalPathAttributes(p.Pattrs)
}

func GetNativeExtCommunityAttribute(subType bgp.ExtendedCommunityAttrSubType, pattrs []bgp.PathAttributeInterface) (bgp.ExtendedCommunityInterface, bool) {
	pattr, ok := GetNativePathAttribute(bgp.BGP_ATTR_TYPE_EXTENDED_COMMUNITIES, pattrs)
	if !ok {
		return nil, false
	}

	for _, extcom := range pattr.(*bgp.PathAttributeExtendedCommunities).Value {
		if _, st := extcom.GetTypes(); st == subType {
			return extcom, true
		}
	}

	return nil, false
}

func GetNexthopIPFromNativePathAttributes(attrs []bgp.PathAttributeInterface) (net.IP, bool) {
	if attr, ok := GetNativePathAttribute(bgp.BGP_ATTR_TYPE_NEXT_HOP, attrs); ok {
		return attr.(*bgp.PathAttributeNextHop).Value, true
	}

	if attr, ok := GetNativePathAttribute(bgp.BGP_ATTR_TYPE_MP_REACH_NLRI, attrs); ok {
		return attr.(*bgp.PathAttributeMpReachNLRI).Nexthop, true
	}

	return nil, false
}

func AddNativePathAttributeNexthop(attrs []bgp.PathAttributeInterface, nh net.IP) []bgp.PathAttributeInterface {
	attr := bgp.NewPathAttributeNextHop(nh.String())
	return SetNativePathAttribute(attrs, attr)
}

func AddNativePathAttributeNLRI(attrs []bgp.PathAttributeInterface, nh net.IP, nlri []bgp.AddrPrefixInterface) []bgp.PathAttributeInterface {
	attr := bgp.NewPathAttributeMpReachNLRI(nh.String(), nlri)
	return SetNativePathAttribute(attrs, attr)
}

func SetNexthopIPToNativePathAttributes(attrs []bgp.PathAttributeInterface, nh net.IP) []bgp.PathAttributeInterface {
	if _, ok := GetNativePathAttribute(bgp.BGP_ATTR_TYPE_NEXT_HOP, attrs); ok {
		attrs = AddNativePathAttributeNexthop(attrs, nh)
	}

	if oldAttr, ok := GetNativePathAttribute(bgp.BGP_ATTR_TYPE_MP_REACH_NLRI, attrs); ok {
		oldNlri := oldAttr.(*bgp.PathAttributeMpReachNLRI)
		attrs = AddNativePathAttributeNLRI(attrs, nh, oldNlri.Value)
	}

	return attrs
}

func GetLabelsFromNativeAddrPrefix(nlri bgp.AddrPrefixInterface) []uint32 {
	switch prefix := nlri.(type) {
	case *bgp.LabeledIPAddrPrefix:
		return prefix.Labels.Labels
	case *bgp.LabeledIPv6AddrPrefix:
		return prefix.Labels.Labels
	case *bgp.LabeledVPNIPAddrPrefix:
		return prefix.Labels.Labels
	case *bgp.LabeledVPNIPv6AddrPrefix:
		return prefix.Labels.Labels
	default:
		return []uint32{}
	}
}

func NewIPv4PrefixFromVPNv4(prefix *bgp.LabeledVPNIPAddrPrefix) *bgp.IPAddrPrefix {
	return bgp.NewIPAddrPrefix(prefix.Length-88, prefix.Prefix.String())
}

func ToRouteFamily(f *api.Family) bgp.RouteFamily {
	return bgp.AfiSafiToRouteFamily(uint16(f.Afi), uint8(f.Safi))
}

func ToApiFamily(afi uint16, safi uint8) *api.Family {
	return &api.Family{
		Afi:  api.Family_Afi(afi),
		Safi: api.Family_Safi(safi),
	}
}
