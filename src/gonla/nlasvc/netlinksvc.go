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

package nlasvc

import (
	"gonla/nlactl"
	"gonla/nlalib"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type NLANetlinkService struct {
	log *log.Entry
}

func NewNLANetlinkService() *NLANetlinkService {
	return &NLANetlinkService{
		log: NewLogger("NLANetlinkService"),
	}
}

func (s *NLANetlinkService) Start(uint8, *nlactl.NLAChannels) error {
	s.log.Infof("Start:")
	return nil
}

func (s *NLANetlinkService) Stop() {
	s.log.Infof("Stop:")
}

func (s *NLANetlinkService) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {
	if nlmsg.Src != nlamsg.SRC_API {
		s.log.Debugf("LINK: skip. %s", nlmsg)
		return
	}

	attr := link.Attrs()
	s.log.Debugf("LINK: attrs:%v", attr)

	switch nlmsg.Type() {
	case syscall.RTM_SETLINK:
		switch attr.OperState {
		case netlink.OperUp:
			s.log.Infof("LINK: OperUp %s", attr.Name)
			if err := netlink.LinkSetUp(link.Link); err != nil {
				s.log.Errorf("NetlinkLink: LinkSetUp error. %s", err)
			}

		case netlink.OperDown:
			s.log.Infof("LINK: OperDown %s", attr.Name)
			if err := netlink.LinkSetDown(link.Link); err != nil {
				s.log.Errorf("LINK: LinkSetDown error. %s", err)
			}
		default:
			s.log.Debugf("LINK: OperState not changed.")
		}
	}
}

func (s *NLANetlinkService) getBridgeVlanVid(neigh *nlamsg.Neigh) uint16 {
	return uint16(neigh.Vlan)

	//vid := uint16(neigh.Vlan)
	//if vid == DEFAULT_VLAN_VID {
	//	nladbm.BrVlans().ListByIndex(neigh.LinkIndex, func(brvlan *nlamsg.BridgeVlanInfo) bool {
	//		if brvlan.PortType() == nlamsg.BRIDGE_VLAN_PORT_ACCESS {
	//			if brvlan.Vid != DEFAULT_VLAN_VID {
	//				vid = brvlan.Vid
	//				return false // no more callback.
	//			}
	//		}
	//
	//		return true
	//  })
	//}
	//
	// return vid
}

func (s *NLANetlinkService) NetlinkNeigh(nlmsg *nlamsg.NetlinkMessage, neigh *nlamsg.Neigh) {
	if nlmsg.Src != nlamsg.SRC_API {
		s.log.Debugf("NEIGH: skip. %s", nlmsg)
		return
	}

	if neigh.IsFdbEntry() {
		// FDB entry.
		fdb := nlalib.NewFDBEntry(
			neigh.HardwareAddr,
			neigh.LinkIndex,
			s.getBridgeVlanVid(neigh),
			0, 0,
		)

		switch nlmsg.Type() {
		case unix.RTM_NEWNEIGH:
			s.log.Debugf("NEIGH: add fdb %s if=%d vid=%d", fdb.HardwareAddr, fdb.LinkIndex, fdb.Vlan)
			if err := netlink.NeighAdd(fdb); err != nil {
				s.log.Errorf("NEIGH: add fdb error. %s", err)
			}

		case nlalink.RTM_SETNEIGH:
			s.log.Debugf("NEIGH: set fdb %s if=%d vid=%d", fdb.HardwareAddr, fdb.LinkIndex, fdb.Vlan)
			if err := netlink.NeighSet(fdb); err != nil {
				s.log.Errorf("NEIGH: set fdb error. %s", err)
			}

		case unix.RTM_DELNEIGH:
			s.log.Debugf("NEIGH: del fdb %s if=%d vid=%d", fdb.HardwareAddr, fdb.LinkIndex, fdb.Vlan)
			if err := netlink.NeighDel(fdb); err != nil {
				s.log.Errorf("NEIGH: del fdb error. %s", err)
			}

		default:
			s.log.Errorf("NEIGH: invalid fdb. msgType=%d", nlmsg.Type())
		}
	}
}
