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

package ribctl

import (
	fibcapi "fabricflow/fibc/api"
	"gonla/nlamsg"
)

func (r *RIBController) NewPortConfig(cmd string, reId string, link *nlamsg.Link) *fibcapi.PortConfig {
	pc := fibcapi.NewPortConfig(cmd, reId, NewLinkName(link, r.useNId), NewPortId(link), NewPortStatus(link))

	if parentIndex := link.Attrs().ParentIndex; parentIndex != 0 {
		if parent, err := r.nla.GetLink(link.NId, parentIndex); err == nil {
			pc.Link = NewLinkName(parent, r.useNId)
		}
	}

	linkType := LinkTypeFromLink(link)

	switch linkType {
	case fibcapi.LinkType_IPTUN:
		if iptun := link.Iptun(); iptun != nil {
			if tun, err := r.nla.GetIptun(link.NId, iptun.Remote); err != nil {
				r.log.Warnf("PortConfig: GetIptun error. %s", err)
			} else {
				pc.DpPort = fibcapi.NewDPPortId(uint32(tun.TnlId), linkType)
				r.log.Debugf("PortConfig: iptun.dp_port=0x%x", pc.DpPort)
			}
		}

	case fibcapi.LinkType_BRIDGE, fibcapi.LinkType_BOND:
		pc.DpPort = fibcapi.NewDPPortId(NewPortId(link), linkType)
		r.log.Debugf("PortConfig: master.dp_port=0x%x", pc.DpPort)

	default:
		// pass
	}

	return pc
}

func (r *RIBController) SendPortConfig(cmd string, link *nlamsg.Link) error {
	r.log.Debugf("PortConfig: %s %v", cmd, link)

	p := r.NewPortConfig(cmd, r.reId, link)
	return r.fib.Send(p, 0)
}
