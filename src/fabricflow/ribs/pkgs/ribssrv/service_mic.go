// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package ribssrv

import (
	"context"
	"fabricflow/ribs/api/ribsapi"
	"fabricflow/ribs/pkgs/ribsdbm"
	"fabricflow/util/gobgp/apiutil"
	fflibnet "fabricflow/util/net"
	"fmt"
	"gonla/nlamsg"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/pkg/packet/bgp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const (
	// ServiceSyncTimeMax is max of sync time.
	ServiceSyncTimeMax = 3600 * time.Second
	// ServiceSyncTimeMin is min od sync time.
	ServiceSyncTimeMin = 100 * time.Millisecond
)

//
// MicService is mic service.
//
type MicService struct {
	RibsService
	NexthopNW string
	SyncTime  time.Duration

	syncCh     chan string
	rics       *ribsdbm.RicTable
	nexthops   *ribsdbm.NexthopTable
	nexthopMap *fflibnet.IPMap
}

func (s *MicService) initMic() error {
	_, nexthopNW, err := net.ParseCIDR(s.NexthopNW)
	if err != nil {
		s.log.Errorf("initMic: bas nexthop. %s", err)
		return err
	}
	ipgen := fflibnet.NewIPMapIPNetGenerator(nexthopNW)

	s.syncCh = make(chan string)
	s.rics = ribsdbm.NewRicTable()
	s.nexthops = ribsdbm.NewNexthopTable()
	s.nexthopMap = fflibnet.NewIPMap(ipgen)

	return nil
}

//
// Start starts main thread.
//
func (s *MicService) Start(done <-chan struct{}) error {
	s.RibsService.init("svc.mic")

	s.log.Debugf("Start:")

	if err := s.initMic(); err != nil {
		return err
	}

	if err := s.startNLACtrl(); err != nil {
		return err
	}

	if err := s.startBGPMonitor(done); err != nil {
		return err
	}

	if err := s.startCoreAPIServer(s); err != nil {
		return err
	}

	if err := s.startAPIServer(s); err != nil {
		return err
	}

	go s.serve(done)

	s.log.Debugf("Start: success.")
	return nil
}

func (s *MicService) newTicker() (*time.Ticker, bool) {
	if s.SyncTime < ServiceSyncTimeMin {
		return time.NewTicker(ServiceSyncTimeMax), true
	}

	return time.NewTicker(s.SyncTime), true
}

func (s *MicService) serve(done <-chan struct{}) {

	s.log.Infof("serve: START")

	ticker, tickActive := s.newTicker()
	defer ticker.Stop()

FOR_LOOP:
	for {
		select {
		case <-s.nla.Conn():
			s.log.Infof("serve: connected to nla.")
			if err := s.nla.Monitor(); err != nil {
				s.log.Errorf("serve: monitor nla error. %s", err)
			}

		case msg, ok := <-s.nla.Recv():
			if ok {
				s.log.Tracef("serve: nlamsg %s", msg)
				if err := nlamsg.DispatchUnion(msg, s); err != nil {
					s.log.Warnf("DispatchUnion  error. %s", err)
				}
			}

		case conn, ok := <-s.bgpConnCh:
			if ok {
				s.log.Debugf("serve: bgp connected %s", conn.RemoteAddr)
			}

		case upd, ok := <-s.bgpPathCh:
			if ok {
				s.log.Tracef("serve: from bgp RT:%s", upd.Rt)
				LogBgpPath(s.log, log.TraceLevel, upd.Path)
				s.sendToRic(upd)
			}

		case upd, ok := <-s.apiPathCh:
			if ok {
				s.log.Tracef("serve: from api RT:%s", upd.Rt)
				LogBgpPath(s.log, log.TraceLevel, upd.Path)
				s.sendToGoBGP(upd)
			}

		case rt, ok := <-s.syncCh:
			if ok {
				s.log.Tracef("serve: sync RT:%s", rt)
				s.sendAllToRic(rt)
			}

		case <-ticker.C:
			if tickActive {
				s.log.Debugf("SyncRib ALL.")
				s.sendAllToRic(ribsdbm.RTany)
			}

		case <-done:
			s.log.Infof("serve: EXIT")
			break FOR_LOOP
		}
	}
}

func (s *MicService) sendToGoBGP(upd *RibUpdate) {
	s.log.Debugf("sendToMic:")

	path := apiutil.NewNativePath(upd.Path)
	nhIP := path.GetNexthop()

	if exist := s.nexthops.Select(nhIP, ribsdbm.RTmic, func(n *ribsdbm.Nexthop) {}); exist {
		s.log.Debugf("sendToRic: Reject(RIC->MIC) via %s rt %s", nhIP, upd.Rt)
		return
	}

	sourceID := net.ParseIP(path.SourceID)

	// register nexthop as ric side.
	s.nexthops.Add(nhIP, upd.Rt, sourceID)
	s.log.Debugf("sendToRic: regiser nexthop %s rt %s src %s", nhIP, upd.Rt, sourceID)

	if err := modGoBGPPath(s.bgp.Client(), upd.Path); err != nil {
		s.log.Errorf("sendToMic: ModPath error. %s", err)
	}
}

func (s *MicService) sendToRic(upd *RibUpdate) {
	s.log.Debugf("sendToRic:")

	path := apiutil.NewNativePath(upd.Path)

	ecRT, ok := getExtendedCommunityRouteTarget(path.Attrs, ribsdbm.RTany)
	if !ok {
		s.log.Warnf("sendToRic: extCommunity(RT) not found. %s", path.Nlri)
		return
	}

	rt := ecRT.String()
	ok = s.rics.Select(rt, func(e *ribsdbm.RicEntry) {
		nlriVPN, ok := path.GetNlri().(*bgp.LabeledVPNIPAddrPrefix)
		if !ok {
			s.log.Errorf("sendToRic: bad nlri type. %s", path.GetNlri())
			return
		}

		nh := path.GetNexthop()
		if exist := s.nexthops.Select(nh, rt, func(nh *ribsdbm.Nexthop) {}); exist {
			s.log.Debugf("sendToRic: Reject(MIC->RIC) %s via %s rt:%s", nlriVPN, nh, rt)
			return
		}

		s.log.Debugf("sendToMic: Path(VPNv4): %s via %s rt %s", nlriVPN, nh, rt)

		aliasNH, err := s.nexthopMap.Value(nh)
		if err != nil {
			s.log.Errorf("sendToRic: generate alias n.h. error. %s", err)
			return
		}
		s.log.Tracef("sendToRic: Nexthop(VPNv4): %s", nh)
		s.log.Tracef("sendToRic: Nexthop(IPv4) : %s", aliasNH)

		nlriIP := apiutil.NewIPv4PrefixFromVPNv4(nlriVPN)
		s.log.Tracef("sendToRic: NLRI(VPNv4): %s", nlriVPN)
		s.log.Tracef("sendToRic: NLRI(IPv4) : %s", nlriIP)

		// register nexthop as mix side.
		s.nexthops.Add(aliasNH, ribsdbm.RTmic, nil)
		s.log.Debugf("sendToRic: register nexthop %s as MIC.", aliasNH)

		if withdraw := path.IsWithdraw; !withdraw {
			labels := apiutil.GetLabelsFromNativeAddrPrefix(path.GetNlri())
			_, dst, _ := net.ParseCIDR(nlriIP.String())

			s.log.Debugf("sendToRic: vpn nid:%d dst:%v gw:%s label:%v vpn-gw:%s",
				e.NId, dst, aliasNH, labels, nh)

			s.nla.AddVpn(e.NId, dst, aliasNH, labels, nh)
		}

		pattrsIP := []bgp.PathAttributeInterface{
			bgp.NewPathAttributeNextHop(aliasNH.String()),
		}
		for _, pattr := range path.GetPathAttrs() {
			switch pattr.GetType() {
			case bgp.BGP_ATTR_TYPE_EXTENDED_COMMUNITIES, bgp.BGP_ATTR_TYPE_MP_REACH_NLRI, bgp.BGP_ATTR_TYPE_MP_UNREACH_NLRI:
				// pass
			default:
				pattrsIP = append(pattrsIP, pattr)
			}
		}

		path.Nlri = nlriIP
		path.Attrs = pattrsIP
		path.SourceID = nh.String()
		path.Family = &api.Family{
			Afi:  api.Family_AFI_IP,
			Safi: api.Family_SAFI_UNICAST,
		}

		pathIP := path.NewAPIPath()

		s.log.Debugf("sendToRic: Path(IPv4) : %s via %s rt %s src %s", nlriIP, aliasNH, rt, nh)
		LogBgpPath(s.log, log.TraceLevel, pathIP)

		p, err := NewRibUpdateAPIFromGoBGPPath(pathIP, rt)
		if err != nil {
			s.log.Errorf("sendToRic: create RibUpdate error. %s", err)
			return
		}

		if err := e.Stream.Send(p); err != nil {
			s.log.Errorf("sendToRic: send error. %s", err)
			return
		}
	})

	if !ok {
		s.log.Debugf("sendToRic: ric not exist. rt:%s", rt)
	}
}

func (s *MicService) sendAllToRic(rt string) {
	s.log.Debugf("sendAllToRic: rt:'%s'", rt)

	go func() {
		err := listGoBGPPath(s.bgp.Client(), s.family, func(p *api.Path) error {
			path := apiutil.NewNativePath(p)
			if _, ok := getExtendedCommunityRouteTarget(path.Attrs, rt); ok {
				s.bgpPathCh <- NewRibUpdate(p, s.Family)
			}

			return nil
		})

		if err != nil {
			s.log.Errorf("sendAllToRic: ListGoBGPPath error. %s", err)
		}
	}()
}

func (s *MicService) deleteNexthopsByRT(rt string) {
	s.log.Debugf("deleteNexthopsByRT: RT:%s", rt)

	s.nexthops.DeleteByRT(rt, func(e *ribsdbm.Nexthop) bool {
		s.log.Tracef("deleteNexthopsByRT: %s", e.Key())
		return true
	})
}

func (s *MicService) deletePathsByRT(rt string) {
	s.log.Debugf("deletePathsByRT: RT:%s", rt)

	err := listGoBGPPath(s.bgp.Client(), s.family, func(p *api.Path) error {
		path := apiutil.NewNativePath(p)

		if _, ok := getExtendedCommunityRouteTarget(path.Attrs, rt); !ok {
			s.log.Debugf("deleteRicByRT: RT not match.")
			return nil
		}

		nhIP := path.GetNexthop()

		s.nexthops.Select(nhIP, rt, func(e *ribsdbm.Nexthop) {
			path.SourceID = e.SourceID.String()
			path.IsWithdraw = true
		})

		delPath := path.NewAPIPath()

		s.log.Debugf("deleteRicByRT: delete path.")
		LogBgpPath(s.log, log.TraceLevel, delPath)

		if err := deleteGoBGPPath(s.bgp.Client(), delPath); err != nil {
			s.log.Errorf("deleteRicByRT: DelPath error. %s", err)
		}

		return nil
	})

	if err != nil {
		s.log.Errorf("deleteRicByRT: ListGoBGPPath error. %s", err)
	}
}

//
// NetlinkRoute process new/del route notification.
//
func (s *MicService) NetlinkRoute(nlmsg *nlamsg.NetlinkMessage, route *nlamsg.Route) {
	if nlmsg.NId == 0 || nlmsg.Type() != unix.RTM_DELROUTE || route.Gw == nil {
		s.log.Debugf("ROUTE: ignored %d %s %s", nlmsg.NId, route.Dst, route.Gw)
		return
	}

	if exist := s.nexthops.Select(route.Gw, ribsdbm.RTmic, func(e *ribsdbm.Nexthop) {}); !exist {
		s.log.Debugf("ROUTE: not VPN route. nid:%d dst:%s gw:%s", nlmsg.NId, route.Dst, route.Gw)
		return
	}

	log.Debugf("ROUTE: DelVpn NId:%d path:%v", nlmsg.NId, route.Dst)

	if err := s.nla.DelVpn(nlmsg.NId, route.Dst); err != nil {
		s.log.Errorf("ROUTE: DelVpn error. %d %s %s", nlmsg.NId, route.Dst, err)
		return
	}
}

//
// ModRib process mod rib request.
//
func (s *MicService) ModRib(ctxt context.Context, req *ribsapi.RibUpdate) (*ribsapi.ModRibReply, error) {
	if req.Path == nil {
		return nil, fmt.Errorf("Invalid request. paths=%v", req.Path)
	}

	path := &api.Path{}
	if err := proto.Unmarshal(req.Path, path); err != nil {
		return nil, err
	}

	s.apiPathCh <- NewRibUpdate(path, req.Rt)

	return &ribsapi.ModRibReply{}, nil
}

//
// MonitorRib process monitor rib request.
//
func (s *MicService) MonitorRib(req *ribsapi.MonitorRibRequest, stream ribsapi.RIBSCoreApi_MonitorRibServer) error {
	s.log.Infof("MonitorRib: start. nid:%d RT:%s", req.NId, req.Rt)

	done := stream.Context().Done()
	if done == nil {
		return fmt.Errorf("Bad request")
	}

	e := &ribsdbm.RicEntry{
		NId:    uint8(req.NId),
		Rt:     req.Rt,
		Stream: stream,
	}
	key := e.Key()

	s.rics.Add(e)
	defer s.rics.Delete(key)

	s.syncCh <- req.Rt

	<-done

	// s.deletePathsByRT(req.Rt)
	s.deleteNexthopsByRT(req.Rt)

	s.log.Infof("MonitorRib: exit. nid:%d RT:%s", req.NId, req.Rt)
	return nil
}

//
// SyncRib process sync rib request.
//
func (s *MicService) SyncRib(ctxt context.Context, req *ribsapi.SyncRibRequest) (*ribsapi.SyncRibReply, error) {
	s.syncCh <- req.Rt
	return &ribsapi.SyncRibReply{}, nil
}

//
// GetNexthops process get nexthops request.
//
func (s *MicService) GetNexthops(req *ribsapi.GetNexthopsRequest, stream ribsapi.RIBSApi_GetNexthopsServer) error {
	s.nexthops.Range(func(key string, nh *ribsdbm.Nexthop) {
		reply := ribsapi.Nexthop{
			Key:      key,
			Rt:       nh.Rt,
			Addr:     nh.Addr.String(),
			SourceId: nh.SourceID.String(),
		}

		if err := stream.Send(&reply); err != nil {
			return
		}
	})
	return nil
}

//
// GetRics process get rics request.
//
func (s *MicService) GetRics(req *ribsapi.GetRicsRequest, stream ribsapi.RIBSApi_GetRicsServer) error {
	s.rics.Range(func(key string, e *ribsdbm.RicEntry) error {
		reply := ribsapi.RicEntry{
			Key: key,
			NId: uint32(e.NId),
			Rt:  e.Rt,
		}

		if err := stream.Send(&reply); err != nil {
			return err
		}

		return nil
	})
	return nil
}

//
// GetNexthopMap process get nexthop map request.
//
func (s *MicService) GetNexthopMap(req *ribsapi.GetIPMapRequest, stream ribsapi.RIBSApi_GetNexthopMapServer) error {
	s.nexthopMap.Walk(func(key string, val net.IP) bool {
		reply := ribsapi.IPMapEntry{
			Key:   key,
			Value: val.String(),
		}

		if err := stream.Send(&reply); err != nil {
			return false
		}

		return true
	})

	return nil
}
