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
	"gonla/nladbm"
	"gonla/nlamsg"
	"net"
)

//
// Iptun
//
func NewIptun(nid uint8, mac net.HardwareAddr, link *nlamsg.Link) *Iptun {
	if mac == nil {
		mac = []byte{}
	}

	return &Iptun{
		NId:      uint32(nid),
		Link:     NewLinkFromNative(link),
		LocalMac: mac,
	}
}

func NewIptunFromNative(v *nlamsg.Iptun) *Iptun {
	mac := v.LocalMAC
	if mac == nil {
		mac = []byte{}
	}

	return &Iptun{
		NId:      uint32(v.NId),
		Link:     NewLinkFromNative(v.Link),
		LocalMac: mac,
		TnlId:    uint32(v.TnlId),
	}
}

func (v *Iptun) ToNative() *nlamsg.Iptun {
	var mac net.HardwareAddr
	if len(v.LocalMac) != 0 {
		mac = v.LocalMac
	}
	return &nlamsg.Iptun{
		Link:     v.Link.ToNative(),
		LocalMAC: mac,
		TnlId:    uint16(v.TnlId),
	}
}

func (v *Iptun) GetRemoteIP() net.IP {
	return net.IP(v.Link.GetIptun().Remote)
}

func (v *Iptun) GetLocalIP() net.IP {
	return net.IP(v.Link.GetIptun().Local)
}

func (v *Iptun) GetLocalMACAddr() net.HardwareAddr {
	return net.HardwareAddr(v.LocalMac)
}

//
// Iptun (Key)
//
func (k *IptunKey) ToNative() *nladbm.IptunKey {
	return &nladbm.IptunKey{
		NId:    uint8(k.NId),
		Remote: k.GetRemoteIP().String(),
	}
}

func NewIptunKeyFromNative(n *nladbm.IptunKey) *IptunKey {
	return &IptunKey{
		NId:    uint32(n.NId),
		Remote: net.ParseIP(n.Remote),
	}
}

func (k *IptunKey) GetRemoteIP() net.IP {
	return net.IP(k.Remote)
}
