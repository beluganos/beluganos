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
	"gonla/nlamsg/nlalink"
	"net"
)

func (v *Vpn) GetIPNet() *net.IPNet {
	return &net.IPNet{
		IP:   net.IP(v.Ip),
		Mask: net.IPMask(v.Mask),
	}
}

func (v *Vpn) NetGw() net.IP {
	return net.IP(v.Gw)
}

func (v *Vpn) NetVpnGw() net.IP {
	return net.IP(v.VpnGw)
}

func (v *Vpn) ToNetlink() *nlalink.Vpn {
	return &nlalink.Vpn{
		Ip:    v.Ip,
		Mask:  v.Mask,
		Gw:    v.Gw,
		Label: v.Label,
		VpnGw: v.VpnGw,
	}
}

func (v *Vpn) ToNative() *nlamsg.Vpn {
	return &nlamsg.Vpn{
		Vpn:   v.ToNetlink(),
		VpnId: v.VpnId,
		NId:   uint8(v.NId),
	}
}

func NewVpn(nid uint8, dst *net.IPNet, gw net.IP, label uint32, vpnGw net.IP) *Vpn {
	if vpnGw == nil {
		vpnGw = gw
	}

	return &Vpn{
		NId:   uint32(nid),
		Ip:    dst.IP,
		Mask:  dst.Mask,
		Gw:    gw,
		Label: label,
		VpnGw: vpnGw,
		VpnId: 0,
	}
}

func NewVpnFromNative(v *nlamsg.Vpn) *Vpn {
	return &Vpn{
		NId:   uint32(v.NId),
		Ip:    v.Ip,
		Mask:  v.Mask,
		Gw:    v.Gw,
		Label: v.Label,
		VpnGw: v.VpnGw,
		VpnId: v.VpnId,
	}
}

//
// Vpn (Key)
//
func (k *VpnKey) ToNative() *nladbm.VpnKey {
	return &nladbm.VpnKey{
		NId: uint8(k.NId),
		Dst: k.Dst,
	}
}

func NewVpnKeyFromNative(n *nladbm.VpnKey) *VpnKey {
	return &VpnKey{
		NId: uint32(n.NId),
		Dst: n.Dst,
	}
}
