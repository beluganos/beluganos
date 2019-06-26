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

package nlaapi

import (
	"net"

	"github.com/vishvananda/netlink"
)

func (b *BondSlaveInfo) SlaveType() string {
	return "bond"
}

func NewBondSlaveInfo() *BondSlaveInfo {
	return &BondSlaveInfo{}
}

func NewBondSlaveInfoFromNative(n *netlink.BondSlaveInfo) *BondSlaveInfo {
	return &BondSlaveInfo{
		State:                  BondState(n.State),
		MiiStatus:              BondLinkState(n.MiiStatus),
		LinkFailureCount:       n.LinkFailureCount,
		PermanentHwAddr:        net.HardwareAddr(n.PermanentHwAddr),
		QueueId:                int32(n.QueueId),
		AggregatorId:           int32(n.AggregatorId),
		ActorOperPortState:     int32(n.ActorOperPortState),
		AdPartnerOperPortState: int32(n.AdPartnerOperPortState),
	}
}

func (b *BondSlaveInfo) ToNative() *netlink.BondSlaveInfo {
	return &netlink.BondSlaveInfo{
		State:                  netlink.BondState(b.State),
		MiiStatus:              netlink.BondLinkStatus(b.MiiStatus),
		LinkFailureCount:       b.LinkFailureCount,
		PermanentHwAddr:        b.PermanentHwAddr,
		QueueId:                int(b.QueueId),
		AggregatorId:           int(b.AggregatorId),
		ActorOperPortState:     int(b.ActorOperPortState),
		AdPartnerOperPortState: int(b.AdPartnerOperPortState),
	}
}

func SlaveInfoToNative(s isLinkAttrs_SlaveInfo) netlink.LinkSlaveInfo {
	switch slaveInfo := s.(type) {
	case *LinkAttrs_BondSlaveInfo:
		return slaveInfo.BondSlaveInfo.ToNative()

	default:
		return nil
	}
}

func NewSlaveInfoFromNative(n netlink.LinkSlaveInfo) isLinkAttrs_SlaveInfo {
	switch slaveInfo := n.(type) {
	case *netlink.BondSlaveInfo:
		return &LinkAttrs_BondSlaveInfo{
			BondSlaveInfo: NewBondSlaveInfoFromNative(slaveInfo),
		}

	default:
		return nil
	}
}
