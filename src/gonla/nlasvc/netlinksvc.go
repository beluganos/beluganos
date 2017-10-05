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
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"gonla/nlactl"
	"gonla/nlamsg"
	"syscall"
)

type NLANetlinkService struct {
}

func NewNLANetlinkService() *NLANetlinkService {
	return &NLANetlinkService{}
}

func (s *NLANetlinkService) Start(uint8, *nlactl.NLAChannels) error {
	log.Infof("NetlinkService: START")
	return nil
}

func (s *NLANetlinkService) Stop() {
	log.Infof("NetlinkService: STOP")
}

func (s *NLANetlinkService) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {
	if nlmsg.Src != nlamsg.SRC_API {
		log.Debugf("NetlinkService: Link skip. %s", nlmsg)
		return
	}

	attr := link.Attrs()
	log.Debugf("NetlinkService: LinkAttr %v", attr)

	switch nlmsg.Type() {
	case syscall.RTM_SETLINK:
		switch attr.OperState {
		case netlink.OperUp:
			log.Infof("NetlinkService: OperUp %s", attr.Name)
			if err := netlink.LinkSetUp(link.Link); err != nil {
				log.Errorf("NetlinkService: LinkSetUp error. %s", err)
			}

		case netlink.OperDown:
			log.Infof("NetlinkService: OperDown %s", attr.Name)
			if err := netlink.LinkSetDown(link.Link); err != nil {
				log.Errorf("NetlinkService: LinkSetDown error. %s", err)
			}
		default:
			log.Debugf("NetlinkService: OperState not changed.")
		}
	}
}
