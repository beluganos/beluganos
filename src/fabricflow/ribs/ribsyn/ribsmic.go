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

package ribsyn

import (
	"fabricflow/ribs/ribsmsg"
	"fmt"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	"gonla/nlalib"
	"net"
)

type MicController struct {
	Client *BgpClient
	RibCh  chan *ribsmsg.RibUpdate
	ConCh  chan *nlalib.ConnInfo
}

func NewMicController() *MicController {
	return &MicController{
		Client: nil,
		RibCh:  make(chan *ribsmsg.RibUpdate),
		ConCh:  make(chan *nlalib.ConnInfo),
	}
}

func (c *MicController) Conn() <-chan *nlalib.ConnInfo {
	return c.ConCh
}

func (c *MicController) Recv() <-chan *ribsmsg.RibUpdate {
	return c.RibCh
}

func (c *MicController) GetRibs(filter func(string) bool) {
	log.Debugf("MICC GetRibs: START")

	c.Client.GetRib("VPNv4", func(rib *ribsmsg.RibUpdate) error {
		ecRt := GetBgpExtCommunityRouteTarget(rib.Paths[0])
		if filter(ecRt.String()) {
			c.RibCh <- rib
		}
		return nil
	})

	log.Debugf("MICC GetRibs: END")
}

func (c *MicController) UpdateRibs(filter func(*table.Path) bool) {
	log.Debugf("MICC UpdRibs: START")

	c.Client.GetRib("VPNv4", func(rib *ribsmsg.RibUpdate) error {
		for _, path := range rib.Paths {
			if filter(path) {
				log.Debugf("MICC UpdRibs: %v", path)
				if err := c.SendBgpPath(path); err != nil {
					log.Errorf("MICC UpdRibs: SendBgpPath error. %s", err)
				}
			}
		}
		return nil
	})

	log.Debugf("MICC UpdRibs: END")
}

func (c *MicController) Start(addr string) error {
	client, err := NewBgpClient(c.ConCh, addr)
	if err != nil {
		return err
	}

	c.Client = client
	log.Infof("MICC Start:")
	return nil
}

func (c *MicController) Monitor() error {
	log.Infof("MICC Monitor:")
	return NewBgpMonitorRib(c.RibCh, c.Client, "VPNv4")
}

func (c *MicController) Translate(vpnPath *table.Path, nexthop net.IP) ([]*table.Path, error) {
	nlri := vpnPath.GetNlri()
	if nlri == nil {
		return nil, fmt.Errorf("MICC Translate: NLRI not found. %v", vpnPath)
	}

	prefix := nlri.(*bgp.LabeledVPNIPAddrPrefix)
	if prefix == nil {
		return nil, fmt.Errorf("MICC Translate: Bad NLRI. %v %v", vpnPath, nlri)
	}

	pattrs := c.TranslatePathAttrs(vpnPath.GetPathAttrs())
	pattrs = append(pattrs, bgp.NewPathAttributeNextHop(nexthop.String()))

	source := *vpnPath.GetSource()
	source.ID = vpnPath.GetNexthop()

	newPath := table.NewPath(
		&source,
		NewIPv4NlriFromVPNv4(prefix),
		vpnPath.IsWithdraw,
		pattrs,
		vpnPath.GetTimestamp(),
		vpnPath.NoImplicitWithdraw(),
	)

	return []*table.Path{newPath}, nil
}

func (c *MicController) TranslatePathAttrs(pattrs []bgp.PathAttributeInterface) []bgp.PathAttributeInterface {
	newPattrs := []bgp.PathAttributeInterface{}
	for _, pattr := range pattrs {
		switch pattr.GetType() {
		case bgp.BGP_ATTR_TYPE_EXTENDED_COMMUNITIES, bgp.BGP_ATTR_TYPE_MP_REACH_NLRI:
			// nothing to do.
		default:
			newPattrs = append(newPattrs, pattr)
		}
	}

	return newPattrs
}

func (c *MicController) SendBgpPath(path *table.Path) error {
	if path.IsWithdraw {
		return c.Client.DelBgpPath(path)
	} else {
		return c.Client.AddBgpPath(path)
	}
}
