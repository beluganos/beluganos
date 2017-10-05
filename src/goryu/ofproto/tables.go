// -*- coding; utf-8 -*-

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

package ofproto

import (
	"fmt"
)

var tableNames_names = map[uint8]string{
	0:  "Ingress",
	5:  "Port_DSCP",
	6:  "Port_PCP",
	7:  "Tun_DSCP",
	8:  "Tun_PCP",
	9:  "Injected_OAM",
	10: "VLAN",
	11: "VLAN1",
	12: "Ingress_M.P.",
	13: "MPLS_L2_Port",
	15: "MPLS_L2_DSCP",
	16: "MPLS_L2_PCP",
	20: "TermMAC",
	21: "L3_Type",
	22: "MPLS0",
	24: "MPLS1",
	25: "MPLS2",
	26: "MPLS_TP_M.P.",
	27: "MPLS_L3_Type",
	28: "MPLS_Trust",
	29: "MPLS_Type",
	30: "Unicast",
	40: "Multicast",
	50: "Bridging",
	55: "L2_Policer",
	56: "L2_PolicerAt",
	60: "ACL",
	65: "Color",
}

func StrTable(tableNo uint8) string {
	if name, ok := tableNames_names[tableNo]; ok {
		return name
	}
	return fmt.Sprintf("(%d)", tableNo)
}

func TableValue(name string) uint8 {
	for v, n := range tableNames_names {
		if n == name {
			return v
		}
	}
	return 255
}
