// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

package gonslib

import (
	"fabricflow/fibc/api"
	"net"

	"github.com/beluganos/go-opennsl/opennsl"
)

func tunnelinitiatorSetSrcMAC(unit int, ifaceId opennsl.L3IfaceID, hwaddr net.HardwareAddr, vid opennsl.Vlan) error {
	return nil
}

func tunnelInitiatorAdd(unit int, group *fibcapi.L3UnicastGroup, ifaceId opennsl.L3IfaceID, pvid opennsl.Vlan) {

}

func tunnelInitiatorDelete(unit int, group *fibcapi.L3UnicastGroup, ifaceId opennsl.L3IfaceID) {

}

func newTunnelTerminator4(dst, src net.IP, port opennsl.Port, tunType opennsl.TunnelType) *opennsl.TunnelTerminator {
	return nil
}

func newTunnelTerminator6(dst, src net.IP, port opennsl.Port, tunType opennsl.TunnelType) *opennsl.TunnelTerminator {
	return nil
}

func newTunnelTerminators(group *fibcapi.L3UnicastGroup) (to4Tun *opennsl.TunnelTerminator, to6Tun *opennsl.TunnelTerminator) {
	return nil, nil
}

func tunnelTerminatorAdd(unit int, group *fibcapi.L3UnicastGroup) {

}

func tunnelTerminatorDelete(unit int, group *fibcapi.L3UnicastGroup) {

}
