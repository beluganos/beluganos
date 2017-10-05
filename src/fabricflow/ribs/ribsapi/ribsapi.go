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

package ribsapi

import (
	"fabricflow/ribs/ribsmsg"
	"github.com/golang/protobuf/proto"
	"github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/table"
)

func NewRicEntryFromNative(key string, r *ribsmsg.RicEntry) *RicEntry {
	return &RicEntry{
		Key:   key,
		NId:   uint32(r.NId),
		Addr:  r.Addr,
		Port:  uint32(r.Port),
		Rt:    r.Rt,
		Rd:    r.Rd,
		Label: r.Label,
	}
}

func NewPathFromNative(p *table.Path) *Path {
	return &Path{
		Prefix:  p.GetNlri().String(),
		Nexthop: p.GetNexthop().String(),
	}
}

func NewNexthopFromNative(nh *ribsmsg.Nexthop) *Nexthop {
	return &Nexthop{
		Addr:  nh.Addr.String(),
		Rt:    nh.Rt,
		SrcId: nh.SrcId.String(),
	}
}

func NewNexthopMapFromNative(key string, val string) *NexthopMap {
	return &NexthopMap{
		Key: key,
		Val: val,
	}
}

func DeserializePath(b []byte) (*table.Path, error) {
	p := &gobgpapi.Path{}
	if err := proto.Unmarshal(b, p); err != nil {
		return nil, err
	}
	return p.ToNativePath()
}

func SerializePath(path *table.Path) ([]byte, error) {
	p := gobgpapi.ToPathApi(path)
	return proto.Marshal(p)
}

func (r *RibUpdate) ToNative() (*ribsmsg.RibUpdate, error) {
	paths := make([]*table.Path, len(r.Paths))
	for i, b := range r.Paths {
		p, err := DeserializePath(b)
		if err != nil {
			return nil, err
		}
		paths[i] = p
	}
	return &ribsmsg.RibUpdate{
		Rt:     r.Rt,
		Prefix: r.Prefix,
		Paths:  paths,
	}, nil
}

func NewRibUpdateFromNative(rib *ribsmsg.RibUpdate) (*RibUpdate, error) {
	paths := make([][]byte, len(rib.Paths))
	for i, p := range rib.Paths {
		b, err := SerializePath(p)
		if err != nil {
			return nil, err
		}
		paths[i] = b
	}
	return &RibUpdate{
		Rt:     rib.Rt,
		Prefix: rib.Prefix,
		Paths:  paths,
	}, nil
}
