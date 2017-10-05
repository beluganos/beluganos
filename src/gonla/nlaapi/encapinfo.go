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
	"gonla/nladbm"
	"gonla/nlamsg"
	"net"
)

func (e *EncapInfo) NetDst() *net.IPNet {
	return &net.IPNet{
		IP:   net.IP(e.Ip),
		Mask: net.IPMask(e.Mask),
	}
}

func (e *EncapInfo) ToNative() *nlamsg.EncapInfo {
	return &nlamsg.EncapInfo{
		Dst:  e.NetDst(),
		Vrf:  e.Vrf,
		EnId: e.EnId,
	}
}

func NewEncapInfo(dst *net.IPNet, vrf uint32) *EncapInfo {
	return &EncapInfo{
		Ip:   dst.IP,
		Mask: dst.Mask,
		Vrf:  vrf,
		EnId: 0,
	}
}

func NewEncapInfoFromNative(e *nlamsg.EncapInfo) *EncapInfo {
	return &EncapInfo{
		Ip:   e.Dst.IP,
		Mask: e.Dst.Mask,
		Vrf:  e.Vrf,
		EnId: e.EnId,
	}
}

//
// EncapInfo Key
//
func (k *EncapInfoKey) ToNative() *nladbm.EncapInfoKey {
	return &nladbm.EncapInfoKey{
		Dst: k.Dst,
		Vrf: k.Vrf,
	}
}

func NewEncapInfoKeyFromNative(k *nladbm.EncapInfoKey) *EncapInfoKey {
	return &EncapInfoKey{
		Dst: k.Dst,
		Vrf: k.Vrf,
	}
}
