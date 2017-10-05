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

package nlalink

import (
	"github.com/golang/protobuf/proto"
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

func NewVpn(ipNet *net.IPNet, gw net.IP, label uint32, vpnGw net.IP) *Vpn {
	if vpnGw == nil {
		vpnGw = gw
	}

	return &Vpn{
		Ip:    ipNet.IP,
		Mask:  ipNet.Mask,
		Gw:    gw,
		Label: label,
		VpnGw: vpnGw,
	}
}

func VpnDeserialize(b []byte) (*Vpn, error) {
	vpn := &Vpn{}
	if err := proto.Unmarshal(b, vpn); err != nil {
		return nil, err
	}
	return vpn, nil
}

func VpnSerialize(vpn *Vpn) ([]byte, error) {
	return proto.Marshal(vpn)
}
