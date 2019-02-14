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

package apiutil

import (
	"github.com/osrg/gobgp/packet/bgp"
)

func GetExtendedCommunity(communities *bgp.PathAttributeExtendedCommunities, typ bgp.ExtendedCommunityAttrSubType) (bgp.ExtendedCommunityInterface, bool) {
	for _, community := range communities.Value {
		if _, subType := community.GetTypes(); subType == typ {
			return community, true
		}
	}
	return nil, false
}

func FilterExtendedCommunity(comms *bgp.PathAttributeExtendedCommunities, filter func(bgp.ExtendedCommunityAttrType, bgp.ExtendedCommunityAttrSubType) bool) *bgp.PathAttributeExtendedCommunities {
	newComms := []bgp.ExtendedCommunityInterface{}
	for _, comm := range comms.Value {
		if ok := filter(comm.GetTypes()); ok {
			newComms = append(newComms, comm)
		}
	}

	if len(newComms) == 0 {
		return nil
	}

	comms.Value = newComms
	return comms
}

func FilterCommunities(attr *bgp.PathAttributeCommunities, filter func(uint32) bool) *bgp.PathAttributeCommunities {
	newComms := []uint32{}
	for _, comm := range attr.Value {
		if ok := filter(comm); ok {
			newComms = append(newComms, comm)
		}
	}

	if len(newComms) == 0 {
		return nil
	}

	attr.Value = newComms
	return attr
}
