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
// netlink.Addr
//
func CopyAddr(src *netlink.Addr) *netlink.Addr {
	dst := *src
	return &dst
}

//
// Addr
//
type Addr struct {
	*netlink.Addr
	AdId   uint32 // auto increment
	Family int32
	Index  int32
	NId    uint8
}

func (a *Addr) Copy() *Addr {
	return &Addr{
		Addr:   CopyAddr(a.Addr),
		AdId:   a.AdId,
		Family: a.Family,
		Index:  a.Index,
		NId:    a.NId,
	}
}

func (a *Addr) String() string {
	return fmt.Sprintf("{Ifindex: %d %s} AdId: %d NId: %d", a.Index, a.Addr, a.AdId, a.NId)
}

func NewAddr(addr *netlink.Addr, family int32, index int32, nid uint8, id uint32) *Addr {
	return &Addr{
		AdId:   id,
		Addr:   addr,
		Family: family,
		Index:  index,
		NId:    nid,
	}
}

func AddrDeserialize(nlmsg *NetlinkMessage) (*Addr, error) {
	addr, family, index, err := netlink.AddrDeserialize(nlmsg.Data)
	if err != nil {
		return nil, err
	}

	return NewAddr(&addr, int32(family), int32(index), nlmsg.NId, 0), nil
}
