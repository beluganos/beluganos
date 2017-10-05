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

package nlaapi

import (
	"github.com/vishvananda/netlink"
	"gonla/nladbm"
	"gonla/nlamsg"
	"net"
)

func (n *Neigh) GetIP() net.IP {
	return net.IP(n.Ip)
}

func (n *Neigh) NetHardwareAddr() net.HardwareAddr {
	return net.HardwareAddr(n.HardwareAddr)
}

func (n *Neigh) ToNetlink() *netlink.Neigh {
	return &netlink.Neigh{
		LinkIndex:    int(n.LinkIndex),
		Family:       int(n.Family),
		State:        int(n.State),
		Type:         int(n.Type),
		Flags:        int(n.Flags),
		IP:           n.GetIP(),
		HardwareAddr: n.NetHardwareAddr(),
	}
}

func (n *Neigh) ToNative() *nlamsg.Neigh {
	return &nlamsg.Neigh{
		Neigh: n.ToNetlink(),
		NeId:  uint16(n.NeId),
		NId:   uint8(n.NId),
	}
}

func NewNeighFromNative(n *nlamsg.Neigh) *Neigh {
	return &Neigh{
		LinkIndex:    int32(n.LinkIndex),
		Family:       int32(n.Family),
		State:        int32(n.State),
		Type:         int32(n.Type),
		Flags:        int32(n.Flags),
		Ip:           n.IP,
		HardwareAddr: n.HardwareAddr,
		NId:          uint32(n.NId),
		NeId:         uint32(n.NeId),
	}
}

//
// Neigh (Key)
//

func (k *NeighKey) ToNative() *nladbm.NeighKey {
	return &nladbm.NeighKey{
		NId:  uint8(k.NId),
		Addr: k.Addr,
	}
}

func NewNeighKeyFromNative(n *nladbm.NeighKey) *NeighKey {
	return &NeighKey{
		NId:  uint32(n.NId),
		Addr: n.Addr,
	}
}
