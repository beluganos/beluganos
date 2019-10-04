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

func NewBondSlaveInfoFromNative(n *netlink.BondSlave) *BondSlaveInfo {
	return &BondSlaveInfo{
		State:                  BondState(n.State),
		MiiStatus:              BondLinkState(n.MiiStatus),
		LinkFailureCount:       n.LinkFailureCount,
		PermanentHwAddr:        net.HardwareAddr(n.PermHardwareAddr),
		QueueId:                int32(n.QueueId),
		AggregatorId:           int32(n.AggregatorId),
		ActorOperPortState:     int32(n.AdActorOperPortState),
		AdPartnerOperPortState: int32(n.AdPartnerOperPortState),
	}
}

func (b *BondSlaveInfo) ToNative() *netlink.BondSlave {
	return &netlink.BondSlave{
		State:                  netlink.BondSlaveState(b.State),
		MiiStatus:              netlink.BondSlaveMiiStatus(b.MiiStatus),
		LinkFailureCount:       b.LinkFailureCount,
		PermHardwareAddr:       b.PermanentHwAddr,
		QueueId:                uint16(b.QueueId),
		AggregatorId:           uint16(b.AggregatorId),
		AdActorOperPortState:   uint8(b.ActorOperPortState),
		AdPartnerOperPortState: uint16(b.AdPartnerOperPortState),
	}
}

func SlaveInfoToNative(s isLinkAttrs_SlaveInfo) netlink.LinkSlave {
	switch slaveInfo := s.(type) {
	case *LinkAttrs_BondSlaveInfo:
		return slaveInfo.BondSlaveInfo.ToNative()

	default:
		return nil
	}
}

func NewSlaveInfoFromNative(n netlink.LinkSlave) isLinkAttrs_SlaveInfo {
	switch slaveInfo := n.(type) {
	case *netlink.BondSlave:
		return &LinkAttrs_BondSlaveInfo{
			BondSlaveInfo: NewBondSlaveInfoFromNative(slaveInfo),
		}

	default:
		return nil
	}
}
