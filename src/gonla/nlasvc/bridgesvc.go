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

package nlasvc

import (
	"fmt"
	"gonla/nlactl"
	"gonla/nladbm"
	"gonla/nlalib"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"

	log "github.com/sirupsen/logrus"
)

type NLABridgeVlanUpdate struct {
	Link    netlink.Link
	MsgType uint16
}

func NewNLABridgeVlanUpdate(msgType uint16, link netlink.Link) *NLABridgeVlanUpdate {
	return &NLABridgeVlanUpdate{
		Link:    link,
		MsgType: msgType,
	}
}

func (n *NLABridgeVlanUpdate) String() string {
	return fmt.Sprintf("%s index=%d", nlamsg.NlMsgTypeStr(n.MsgType), n.Link.Attrs().Index)
}

type NLABridgeVlanService struct {
	nid     uint8
	service nlactl.NLAService
	table   nladbm.BridgeVlanInfoTable
	updTime time.Duration
	updChan chan *NLABridgeVlanUpdate
	chans   *nlactl.NLAChannels
	done    chan struct{}
	log     *log.Entry
}

func NewNLABridgeVlanService(service nlactl.NLAService, updTime time.Duration, updChanSize int) *NLABridgeVlanService {
	return &NLABridgeVlanService{
		nid:     0,
		service: service,
		table:   nladbm.NewBridgeVlanInfoTable(),
		updTime: updTime,
		updChan: make(chan *NLABridgeVlanUpdate, updChanSize),
		done:    make(chan struct{}),
		log:     NewLogger("NLABridgeVlanService"),
	}
}

func (s *NLABridgeVlanService) sendBrVlanMsg(msgType uint16, brvlan *nlamsg.BridgeVlanInfo) {
	s.log.Debugf("sendBrVlanMsg: %s %s", nlamsg.NlMsgTypeStr(msgType), brvlan)

	nlmsg, err := brvlan.ToNetlinkMessage(msgType)
	if err != nil {
		s.log.Errorf("sendBrVlanMsg: ToNetlinkMessage error. %s", err)
		return
	}

	s.chans.NlMsg <- nlamsg.NewNetlinkMessage(nlmsg, s.nid, nlamsg.SRC_KNL)
}

func (s *NLABridgeVlanService) sendBrVlanMsgs(msgType uint16, brvlanInfos []*nl.BridgeVlanInfo, link netlink.Link) {
	mtype := nlamsg.NlMsgTypeStr(msgType)

	s.log.Debugf("sendBrVlanMsgs: %s index=%d", mtype, link.Attrs().Index)

	for _, brvlanInfo := range brvlanInfos {
		brvlan := nlamsg.NewBridgeVlanInfoFromNetlink(s.nid, brvlanInfo, link)

		if portType := brvlan.PortType(); portType != nlamsg.BRIDGE_VLAN_PORT_ACCESS && portType != nlamsg.BRIDGE_VLAN_PORT_TRUNK {
			s.log.Debugf("sendBrVlanMsgs: SKIP %s %s", mtype, brvlan)
			continue
		}

		switch msgType {
		case nlalink.RTM_NEWBRIDGE:
			if old := s.table.Insert(brvlan); old != nil {
				if same := brvlan.Equals(old); same {
					s.log.Debugf("sendBrVlanMsgs: NOT CHANGED %s %s", mtype, brvlan)
					continue
				}
			}

			// RTM_NEWBRIDGE
			s.log.Debugf("sendBrVlanMsgs: NEW %s %s", mtype, brvlan)
			s.sendBrVlanMsg(nlalink.RTM_NEWBRIDGE, brvlan)

		case nlalink.RTM_DELBRIDGE:
			key := nladbm.BridgeVlanInfoToKey(brvlan)
			if old := s.table.Delete(key); old == nil {
				s.log.Debugf("sendBrVlanMsgs: NOT EXIST %s %s", mtype, brvlan)
				continue
			}

			// RTM_DELBRIDGE
			s.log.Debugf("sendBrVlanMsgs: DEL %s %s", mtype, brvlan)
			s.sendBrVlanMsg(nlalink.RTM_DELBRIDGE, brvlan)

		default:
			s.log.Errorf("sendBrVlanMsgs: NOP mtype:%s %s", mtype, brvlan)
		}
	}
}

func (s *NLABridgeVlanService) cleanBrVlan(index int) {
	s.log.Debugf("cleanBrVlan: index=%d", index)

	for _, key := range s.table.GCList(index) {
		if brvlan := s.table.Delete(&key); brvlan != nil {
			s.log.Debugf("cleanBrVlan: Clear index=%d", brvlan.Index)
			s.sendBrVlanMsg(nlalink.RTM_DELBRIDGE, brvlan)
		}
	}
}

func (s *NLABridgeVlanService) updateBrVlans() {
	s.log.Debugf("updateBrVlans:")

	brvlanMap, err := netlink.BridgeVlanList()
	if err != nil {
		s.log.Errorf("updateBrVlans: BridgeVlanList error. %s", err)
		return
	}

	s.table.GCInit()

	for ifindex, brvlans := range brvlanMap {
		link, err := netlink.LinkByIndex(int(ifindex))
		if err != nil {
			s.log.Errorf("updateBrVlans: Link get error. index=%d %s", ifindex, err)
			continue
		}

		s.log.Debugf("updateBrVlans: Update index=%d", ifindex)

		s.sendBrVlanMsgs(nlalink.RTM_NEWBRIDGE, brvlans, link)
	}

	s.cleanBrVlan(-1)
}

func (s *NLABridgeVlanService) updateBrVlan(msgType uint16, link netlink.Link) {
	s.log.Debugf("updateBrVlan: index=%d", link.Attrs().Index)

	brvlanMap, err := netlink.BridgeVlanList()
	if err != nil {
		s.log.Errorf("updateBrVlan: BridgeVlanList error. %s", err)
		return
	}

	ifindex := link.Attrs().Index

	s.log.Debugf("updateBrVlan: Update index=%d", ifindex)

	brvlans, ok := brvlanMap[int32(ifindex)]
	if !ok {
		brvlans = []*nl.BridgeVlanInfo{}
	}

	s.table.GCInit()
	s.sendBrVlanMsgs(msgType, brvlans, link)
	s.cleanBrVlan(ifindex)
}

func (s *NLABridgeVlanService) NetlinkLink(nlmsg *nlamsg.NetlinkMessage, link *nlamsg.Link) {

	if nlmsg.Src != nlamsg.SRC_KNL {
		s.log.Debugf("LINK: skip. %s", nlmsg)
		return
	}

	if nid := link.NId; nid != s.nid {
		s.log.Debugf("LINK: skip. nid %d", nid)
		return
	}

	s.log.Debugf("LINK: %d index=%d", nlmsg.Type(), link.Attrs().Index)

	msgType := func() uint16 {
		switch nlmsg.Type() {
		case syscall.RTM_NEWLINK, syscall.RTM_SETLINK:
			return nlalink.RTM_NEWBRIDGE

		case syscall.RTM_DELLINK:
			return nlalink.RTM_DELBRIDGE

		default:
			return 0
		}
	}()

	s.log.Debugf("NetlinkLink: %d -> %d", nlmsg.Type(), msgType)
	s.updChan <- NewNLABridgeVlanUpdate(msgType, link.Link)
}

func (n *NLABridgeVlanService) NetlinkNeigh(nlmsg *nlamsg.NetlinkMessage, neigh *nlamsg.Neigh) {
	if nlmsg.Src != nlamsg.SRC_KNL {
		n.log.Debugf("NEIG skip. %s", nlmsg)
		return
	}

	if n.nid != neigh.NId {
		n.log.Debugf("NEIGH: skip. nid %d", neigh.NId)
	}

	if !neigh.IsFdbEntry() {
		n.log.Debugf("NEIG skip. not fdb. %s", nlmsg)
		return
	}

	if nlalib.IsInvalidHardwareAddr(neigh.HardwareAddr) {
		n.log.Debugf("NEIG skip. mac '%s'", neigh.HardwareAddr)
		return
	}

	n.log.Debugf("NEIG")

	switch nlmsg.Type() {
	case syscall.RTM_NEWNEIGH:
		n.log.Debugf("NEIGH: NEW %v", neigh)
		if old := nladbm.Neighs().Insert(neigh); old != nil {
			nlmsg.Header.Type = nlalink.RTM_SETNEIGH
		}

		nlamsg.DispatchNeigh(nlmsg, neigh, n.service)

	case syscall.RTM_DELNEIGH:
		n.log.Debugf("NEIGH: DEL %v", neigh)
		if old := nladbm.Neighs().Delete(nladbm.NeighToKey(neigh)); old != nil {
			neigh.NeId = old.NeId
		}

		nlamsg.DispatchNeigh(nlmsg, neigh, n.service)

	default:
		n.log.Errorf("NEIGH Invalid message. %v", nlmsg)
	}
}

func (s *NLABridgeVlanService) Serve(done chan struct{}) {
	s.log.Infof("Serve: START")

	tick := time.NewTicker(s.updTime)
	defer tick.Stop()

FOR_LABEL:
	for {
		select {
		case <-tick.C:
			s.log.Debugf("Serve: Auto Update")
			s.updateBrVlans()

		case update := <-s.updChan:
			s.log.Debugf("Serve: %s", update)
			s.updateBrVlan(update.MsgType, update.Link)

		case <-s.done:
			s.log.Infof("Serve: EXIT")
			break FOR_LABEL
		}
	}
}

func (s *NLABridgeVlanService) Start(nid uint8, chans *nlactl.NLAChannels) error {
	s.nid = nid
	s.chans = chans
	go s.Serve(s.done)
	return nil
}

func (s *NLABridgeVlanService) Stop() {
	close(s.done)
}
