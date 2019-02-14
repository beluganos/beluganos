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

package nlaapi

import (
	"gonla/nlamsg"
	"net"
)

//
// NeighNotun
//
func (n *NeighNotun) ToNative() nlamsg.NeighTunnel {
	return nil
}

func NewNeighNotunFromNative(n nlamsg.NeighTunnel) *NeighNotun {
	return &NeighNotun{}
}

//
// NeighIptun
//
func (n *NeighIptun) GetSrcIP() net.IP {
	return net.IP(n.SrcIp)
}

func (n *NeighIptun) ToNative() *nlamsg.NeighIptun {
	return &nlamsg.NeighIptun{
		TunType: n.TunType,
		SrcIP:   n.SrcIp,
	}
}

func NewNeighIptunFromNative(n *nlamsg.NeighIptun) *NeighIptun {
	if n == nil {
		return nil
	}

	return &NeighIptun{
		TunType: n.TunType,
		SrcIp:   n.SrcIP,
	}
}
