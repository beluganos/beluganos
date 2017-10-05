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

package nlamsg

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

//
// netlink.Neigh
//
func CopyNeigh(src *netlink.Neigh) *netlink.Neigh {
	dst := *src
	return &dst
}

//
// Neigh
//
type Neigh struct {
	*netlink.Neigh
	NeId uint16
	NId  uint8
}

func (n *Neigh) Copy() *Neigh {
	return &Neigh{
		Neigh: CopyNeigh(n.Neigh),
		NeId:  n.NeId,
		NId:   n.NId,
	}
}

func (n *Neigh) String() string {
	return fmt.Sprintf("{Ifindex: %d %s} NeId: %d, NId: %d", n.LinkIndex, n.Neigh, n.NeId, n.NId)
}

func NewNeigh(neigh *netlink.Neigh, nid uint8, id uint16) *Neigh {
	return &Neigh{
		NeId:  id,
		Neigh: neigh,
		NId:   nid,
	}
}

func NeighDeserialize(nlmsg *NetlinkMessage) (*Neigh, error) {
	neigh, err := netlink.NeighDeserialize(nlmsg.Data)
	if err != nil {
		return nil, err
	}

	return NewNeigh(neigh, nlmsg.NId, 0), nil
}
