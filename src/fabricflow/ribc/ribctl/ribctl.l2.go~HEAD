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
	"fmt"
	"net"

	fibcapi "fabricflow/fibc/api"

	"golang.org/x/sys/unix"
)

func L2AddrReasonToMsgType(reason fibcapi.L2Addr_Reason) (uint16, error) {
	switch reason {
	case fibcapi.L2Addr_ADD:
		return unix.RTM_NEWNEIGH, nil

	case fibcapi.L2Addr_DELETE:
		return unix.RTM_DELNEIGH, nil

	default:
		return 0, fmt.Errorf("Bad reason. %d", reason)
	}
}

func (r *RIBController) SetFdb(addr *fibcapi.L2Addr, ifindex int) error {
	mtype, err := L2AddrReasonToMsgType(addr.Reason)
	if err != nil {
		return err
	}

	hwaddr, err := net.ParseMAC(addr.HwAddr)
	if err != nil {
		return err
	}

	nid, _ := ParsePortId(addr.PortId)
	if err := r.nla.ModFdb(nid, ifindex, hwaddr, uint16(addr.VlanVid), mtype); err != nil {
		return err
	}

	return nil
}
