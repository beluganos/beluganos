// -*- coding: utf-8 -*-
// +build !gobgpv2

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
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/table"
)

type Path struct {
	*table.Path
}

func NewNativePath(p *api.Path) (*Path, error) {
	path, err := p.ToNativePath()
	if err != nil {
		return nil, err
	}
	return &Path{path}, nil
}

func NewApiPath(p *Path) *api.Path {
	return api.ToPathApi(p.Path, nil)
}

func GetNativePathAttributes(p *api.Path) ([]bgp.PathAttributeInterface, error) {
	return p.GetNativePathAttributes()
}

func GetPathAttribute(path *Path, typ bgp.BGPAttrType) (bgp.PathAttributeInterface, bool) {
	for _, attr := range path.GetPathAttrs() {
		if attr.GetType() == typ {
			return attr, true
		}
	}

	return nil, false
}
