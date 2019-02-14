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

package nlamsg

import (
	"fmt"
	"net"
)

type NeighIptun struct {
	TunType string
	SrcIP   net.IP
}

func (n *NeighIptun) Copy() NeighTunnel {
	dst := *n
	return &dst
}

func (n *NeighIptun) String() string {
	return fmt.Sprintf("type: '%s', src: %s", n.TunType, n.SrcIP)
}

func NewNeighIptun(tunType string, srcIp net.IP) *NeighIptun {
	return &NeighIptun{
		TunType: tunType,
		SrcIP:   srcIp,
	}
}
