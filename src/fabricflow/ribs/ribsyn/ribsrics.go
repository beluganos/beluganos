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
	"fabricflow/ribs/ribsapi"
	"fabricflow/ribs/ribsmsg"
	"fmt"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	"gonla/nlalib"
)

//
// RIC Entry
//
type RicEntry struct {
	ribsmsg.RicEntry
	Client *BgpClient
	nlaCtl *NLAController
	Ch     chan *ribsapi.RibUpdate
}

func (r *RicEntry) Active() bool {
	return (r != nil) && (r.Client != nil)
}

func (r *RicEntry) Start(conCh chan *nlalib.ConnInfo) error {
	if r == nil {
		return fmt.Errorf("RICE Start: Invalid entry. %v", r)
	}

	client, err := NewBgpClient(conCh, r.GetAddr())
	if err != nil {
		return err
	}

	r.Client = client
	log.Infof("RICE Start: %s", r.GetAddr())
	return nil
}

func (r *RicEntry) MonitorRib(ch chan<- *ribsmsg.RibUpdate) error {
	if !r.Active() {
		return fmt.Errorf("RICE Monitor: Invalid status. %v", r)
	}

	log.Infof("RICE Monitor: %s", r.GetAddr())
	return NewBgpMonitorRib(ch, r.Client, r.Rt)
}

func (r *RicEntry) Stop() {
	if r.Active() {
		r.Client.Conn.Close()
		r.Client = nil
		log.Infof("RICE Stop: %s", r.GetAddr())
	}
}

func (r *RicEntry) GetRibs() {
	err := r.Client.GetRib(r.Rt, func(rib *ribsmsg.RibUpdate) error {
		nrib, err := ribsapi.NewRibUpdateFromNative(rib)
		if err != nil {
			return err
		}

		r.Ch <- nrib
		return nil
	})

	if err != nil {
		log.Errorf("RICE GetRibs: GetRib error. %s", err)
	}
}

func (r *RicEntry) Translate(rib *ribsmsg.RibUpdate, f func(*table.Path, *table.Path, *RicEntry) error) error {
	if !r.Active() {
		return fmt.Errorf("RICE Translate: Invalid status. %v", r)
	}

	for _, ipPath := range rib.Paths {
		vpnPath, err := r.TranslatePath(ipPath)

		if err != nil {
			return err
		}

		if err := f(vpnPath, ipPath, r); err != nil {
			return err
		}
	}

	return nil
}

func (r *RicEntry) TranslatePath(path *table.Path) (*table.Path, error) {
	prefix, prefixlen, ok := GetBgpAddrPrefixFromNlri(path.GetNlri())
	if !ok {
		return nil, fmt.Errorf("RICE TransPath: Get prefix error. %v", path)
	}

	newNlri, err := NewBgpMpReachNlri(r.Rd, prefix, prefixlen, r.Label)
	if err != nil {
		return nil, err
	}

	ecRt, err := NewBgpExtCommunityRouteTarget(r.Rt)
	if err != nil {
		return nil, err
	}

	pattrs := r.TranslatePathAttrs(path.GetPathAttrs())
	pattrs = append(pattrs, bgp.NewPathAttributeNextHop(path.GetNexthop().String()))
	pattrs = append(pattrs, ecRt)

	newPath := table.NewPath(
		path.GetSource(),
		newNlri,
		path.IsWithdraw,
		pattrs,
		path.GetTimestamp(),
		path.NoImplicitWithdraw(),
	)

	return newPath, nil
}

func (r *RicEntry) TranslatePathAttrs(pattrs []bgp.PathAttributeInterface) []bgp.PathAttributeInterface {
	newPattrs := []bgp.PathAttributeInterface{}
	for _, pattr := range pattrs {
		switch pattr.GetType() {
		case bgp.BGP_ATTR_TYPE_NEXT_HOP:
			// nothing to do.
		default:
			newPattrs = append(newPattrs, pattr)
		}
	}
	return newPattrs
}

func (r *RicEntry) SendBgpPath(path *table.Path) error {
	if !r.Active() {
		return fmt.Errorf("RICE Send: Invalid status. %v", r)
	}

	if path.IsWithdraw {
		return r.Client.DelBgpPath(path)
	} else {
		return r.Client.AddBgpPath(path)
	}
}

func NewRicEntryOnRic(nid uint8, addr string, port uint16, rt, rd string, label uint32, nlaCtl *NLAController) *RicEntry {
	return &RicEntry{
		Client: nil,
		nlaCtl: nlaCtl,
		Ch:     nil,
		RicEntry: ribsmsg.RicEntry{
			NId:   nid,
			Addr:  addr,
			Port:  port,
			Rt:    rt,
			Rd:    rd,
			Label: label,
			Leave: false,
		},
	}
}

func NewRicEntryOnMic(nid uint8, rt string) *RicEntry {
	return &RicEntry{
		Client: nil,
		nlaCtl: nil,
		Ch:     make(chan *ribsapi.RibUpdate),
		RicEntry: ribsmsg.RicEntry{
			NId:   nid,
			Rt:    rt,
			Leave: false,
		},
	}
}

func (r *RicEntry) SendRib(rib *ribsmsg.RibUpdate) error {
	nrib, err := ribsapi.NewRibUpdateFromNative(rib)
	if err != nil {
		return err
	}

	r.Ch <- nrib
	return nil
}

func (r *RicEntry) Close() {
	if r.Ch != nil {
		close(r.Ch)
		r.Ch = nil
	}
}
