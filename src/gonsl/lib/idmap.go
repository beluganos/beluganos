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
	"fabricflow/fibc/api"

	"github.com/beluganos/go-opennsl/opennsl"
)

//
// IDMaps has sub-maps.
//
type IDMaps struct {
	L2Stations *opennsl.L2StationIDMap
	L3Ifaces   *opennsl.L3IfaceIDMap
	L3Egress   *opennsl.L3EgressIDMap
}

//
// NewIDMaps returns new instance.
//
func NewIDMaps() *IDMaps {
	return &IDMaps{
		L2Stations: opennsl.NewL2StationIDMap(nil),
		L3Ifaces:   opennsl.NewL3IfaceIDMap(fibcapi.NewL2InterfaceGroupID),
		L3Egress:   opennsl.NewL3EgressIDMap(fibcapi.NewL3UnicastGroupID),
	}
}
