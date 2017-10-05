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
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	"gonla/nlamsg"
	"syscall"
	"time"
)

const (
	MIC_SERVICE_INTETVAL_MAX = 3600 * time.Second
	MIC_SERVICE_INTETVAL_MIN = 100 * time.Millisecond
	MIC_SERVICE_RT_ANY       = "*"
)

//
// RIBS/MIC Service
//
type MicService struct {
	mic *MicController
	nla *NLAController
	api *CoreApiServer

	interval time.Duration
}

func NewMicService(intervalMsec int64) *MicService {
	return &MicService{
		mic: NewMicController(),
		nla: NewNLAController(),
		api: NewCoreApiServer(),

		interval: time.Duration(intervalMsec) * time.Millisecond,
	}
}

func (s *MicService) Start(nlaAddr string, micAddr string, apiAddr string) error {
	if err := s.nla.Start(nlaAddr); err != nil {
		return err
	}

	if err := s.mic.Start(micAddr); err != nil {
		return err
	}

	if err := s.api.Start(apiAddr); err != nil {
		return err
	}

	return nil
}

func (s *MicService) Serve() {
	d := func() time.Duration {
		if s.interval > MIC_SERVICE_INTETVAL_MIN {
			log.Infof("MICS: Auto Re-Sync enabled. (%s)", s.interval)
			return s.interval
		}

		return MIC_SERVICE_INTETVAL_MAX
	}()

	ticker := time.NewTicker(d)

	for {
		select {
		case <-s.mic.Conn():
			log.Infof("MICS: MIC connected.")
			if err := s.mic.Monitor(); err != nil {
				log.Errorf("MICS: mic.Monitor error. %s", err)
			}

		case <-s.nla.Conn():
			log.Infof("MICS: NLA connected.")
			if err := s.nla.Monitor(); err != nil {
				log.Errorf("MICS: nla.Monitor error. %s", err)
			}

		case conn := <-s.api.Conn():
			if conn.Leave {
				s.DeleteRic(conn.Rt)
				log.Infof("MICS: RIC disconnected. %v", conn)
			} else {
				s.AddRic(conn)
				log.Infof("MICS: RIC connected. %v", conn)
			}

		case rt := <-s.api.Sync():
			log.Infof("MICS: SyncRib(RT=%s)", rt)
			s.SendAllToRic(rt)

		case <-ticker.C:
			if s.interval > MIC_SERVICE_INTETVAL_MIN {
				log.Debugf("MICS: SyncRib(%s)", MIC_SERVICE_RT_ANY)
				s.SendAllToRic(MIC_SERVICE_RT_ANY)
			}

		case rib := <-s.mic.Recv():
			if err := s.SendToRic(rib); err != nil {
				log.Errorf("MICS: SendToRic error %s", err)
			}

		case rib := <-s.api.Recv():
			if err := s.SendToMic(rib); err != nil {
				log.Errorf("MICS: SendToMic error %s", err)
			}

		case msg := <-s.nla.Recv():
			if err := nlamsg.DispatchUnion(msg, s); err != nil {
				log.Errorf("MICS: DispatchUnion  error. %s", err)
			}
		}
	}
}

// Process Netlink event on RIC.
// if type is DELROUTE and the route is VPN, unregister VPN from NLA.
func (s *MicService) NetlinkRoute(nlmsg *nlamsg.NetlinkMessage, route *nlamsg.Route) {
	if nlmsg.NId == 0 || nlmsg.Type() != syscall.RTM_DELROUTE || route.Gw == nil {
		log.Debugf("MICS: route ignored %s %s", nlmsg, route)
		return
	}

	if nh := Tables.Nexthops.FindMic(route.Gw); nh == nil {
		log.Debugf("MICS: not VPN route %s %s", nlmsg, route)
		return
	}

	log.Debugf("MICS: DelVpn NId:%d path:%v", nlmsg.NId, route.Dst)
	if err := s.nla.DelVpn(nlmsg.NId, route.Dst); err != nil {
		log.Errorf("MICS: DelVpn error. %s %s", route, err)
	}
}

func (s *MicService) SendToMic(rib *ribsmsg.RibUpdate) error {

	for _, path := range rib.Paths {
		nexthop := path.GetNexthop()
		if nh := Tables.Nexthops.FindMic(nexthop); nh != nil {
			log.Debugf("MICS: Reject(RIC->MIC) %v", path)
			return nil
		}

		srcId := path.GetSource().ID
		log.Debugf("MICS: Nexthop RIC:%s RT:%s Src:%s", nexthop, rib.Rt, srcId)
		Tables.Nexthops.AddRic(nexthop, rib.Rt, srcId)

		log.Debugf("MICS: SendToMic %v RT:%s", path, rib.Rt)
		if err := s.mic.SendBgpPath(path); err != nil {
			log.Errorf("MICS: mic.SendBgpPath error. %s", err)
		}
	}

	return nil
}

func (s *MicService) SendToRic(rib *ribsmsg.RibUpdate) error {

	for _, vpnPath := range rib.Paths {
		ecRt := GetBgpExtCommunityRouteTarget(vpnPath)
		if ecRt == nil {
			log.Warnf("MICS: ExtCom(RT) not found. %v", vpnPath)
			continue
		}

		rt := ecRt.String()
		ric := Tables.Rics.FindByRt(rt)
		if ric == nil {
			log.Warnf("MICS: RIC entry not found. %s", rt)
			continue
		}

		nexthopVpn := vpnPath.GetNexthop()
		if nh := Tables.Nexthops.FindRic(nexthopVpn, rt); nh != nil {
			log.Debugf("MICS: Reject(MIC->RIC) %v", vpnPath)
			continue
		}

		nexthopIP, err := Tables.NexthopMap.Value(nexthopVpn)
		if err != nil {
			log.Errorf("MICS: Nexthop overflow. %s", nexthopVpn)
			continue
		}

		Tables.Nexthops.AddMic(nexthopIP)
		log.Debugf("MICS: Nexthop MIC:%s(%s) RT:%s", nexthopVpn, nexthopIP, rt)

		ipPaths, err := s.mic.Translate(vpnPath, nexthopIP)
		if err != nil {
			log.Warnf("MICS: Translate failed. %s", err)
			continue
		}

		labels := GetLabelsFromPath(vpnPath)

		for _, ipPath := range ipPaths {
			if !ipPath.IsWithdraw {
				log.Debugf("MICS: AddVpn NId:%d Path:%v Label:%v vpnGw:%s",
					ric.NId, ipPath, labels, nexthopVpn)
				s.nla.AddVpn(ric.NId, ipPath, labels, nexthopVpn)
			}
		}

		r := ribsmsg.NewRibUpdate("", ipPaths, rt)
		log.Debugf("MICS: send %v", r)
		if err := ric.SendRib(r); err != nil {
			return err
		}
	}

	return nil
}

func (s *MicService) SendAllToRic(rt string) {
	filter := func(ecRt string) bool {
		if rt == MIC_SERVICE_RT_ANY {
			return true
		}
		return ecRt == rt
	}

	go s.mic.GetRibs(filter)
}

func (s *MicService) AddRic(ric *RicEntry) {
	Tables.Rics.Add(ric)
}

func (s *MicService) DeleteRic(rt string) {
	s.mic.UpdateRibs(func(path *table.Path) bool {
		ecRt := GetBgpExtCommunityRouteTarget(path).String()
		if rt != MIC_SERVICE_RT_ANY && rt != ecRt {
			return false
		}

		entry := Tables.Nexthops.FindRic(path.GetNexthop(), ecRt)
		if entry == nil {
			return false
		}

		path.GetSource().ID = entry.SrcId
		path.IsWithdraw = true
		return true
	})
	Tables.Rics.Delete(rt)
}
