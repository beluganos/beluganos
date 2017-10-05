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
// Link
//
type Link struct {
	netlink.Link
	LnId uint16 // auto increment
	NId  uint8
}

func (ln *Link) VlanId() uint16 {
	if vlan := ln.Vlan(); vlan != nil {
		return uint16(vlan.VlanId)
	}
	return 0
}

func (ln *Link) Vlan() *netlink.Vlan {
	if vlan, ok := ln.Link.(*netlink.Vlan); ok {
		return vlan
	}
	return nil
}

func (ln *Link) Bond() *netlink.Bond {
	if bond, ok := ln.Link.(*netlink.Bond); ok {
		return bond
	}
	return nil
}

func (ln *Link) Copy() *Link {
	return &Link{
		Link: ln.Link,
		LnId: ln.LnId,
		NId:  ln.NId,
	}
}

func (ln *Link) String() string {
	a := ln.Attrs()
	return fmt.Sprintf("{Ifindex: %d %s %s %s %s Oper:%s Parent: %d Master: %d} LnId: %d NId: %d",
		a.Index, a.Name, a.HardwareAddr, ln.Type(), a.Flags, a.OperState,
		a.ParentIndex, a.MasterIndex, ln.LnId, ln.NId)
}

func NewLink(link netlink.Link, nid uint8, id uint16) *Link {
	return &Link{
		LnId: id,
		Link: link,
		NId:  nid,
	}
}

func LinkDeserialize(nlmsg *NetlinkMessage) (*Link, error) {
	link, err := netlink.LinkDeserialize(&nlmsg.Header, nlmsg.Data)
	if err != nil {
		return nil, err
	}

	return NewLink(link, nlmsg.NId, 0), nil
}
