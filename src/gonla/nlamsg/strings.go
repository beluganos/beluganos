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

package nlamsg

import (
	"fmt"
	"gonla/nlamsg/nlalink"
	"syscall"
)

var RTMStrings = map[uint16]string{
	nlalink.RTM_NEWNODE:  "RTM_NEWNODE",
	nlalink.RTM_DELNODE:  "RTM_DELNODE",
	nlalink.RTM_SETNODE:  "RTM_SETNODE",
	nlalink.RTM_NEWVPN:   "RTM_NEWVPN",
	nlalink.RTM_DELVPN:   "RTM_DELVPN",
	nlalink.RTM_SETVPN:   "RTM_SETVPN",
	syscall.RTM_NEWLINK:  "RTM_NEWLINK",
	syscall.RTM_DELLINK:  "RTM_DELLINK",
	syscall.RTM_SETLINK:  "RTM_SETLINK",
	syscall.RTM_NEWADDR:  "RTM_NEWADDR",
	syscall.RTM_DELADDR:  "RTM_DELADDR",
	nlalink.RTM_SETADDR:  "RTM_SETADDR",
	syscall.RTM_NEWNEIGH: "RTM_NEWNEIGH",
	syscall.RTM_DELNEIGH: "RTM_DELNEIGH",
	nlalink.RTM_SETNEIGH: "RTM_SETNEIGH",
	syscall.RTM_NEWROUTE: "RTM_NEWROUTE",
	syscall.RTM_DELROUTE: "RTM_DELROUTE",
	nlalink.RTM_SETROUTE: "RTM_SETROUTE",
}

func NlMsgTypeStr(t uint16) string {
	if s, ok := RTMStrings[t]; ok {
		return s
	}
	return fmt.Sprintf("RTM_UNSPEC(%d)", t)
}

var RTMGRPStrings = map[uint16]string{
	nlalink.RTMGRP_NODE:  "RTMGRP_NODE",
	nlalink.RTMGRP_VPN:   "RTMGRP_VPN",
	nlalink.RTMGRP_LINK:  "RTMGRP_LINK",
	nlalink.RTMGRP_ADDR:  "RTMGRP_ADDR",
	nlalink.RTMGRP_NEIGH: "RTMGRP_NEIGH",
	nlalink.RTMGRP_ROUTE: "RTMGRP_ROUTE",
}

func NlMsgGroupStr(g uint16) string {
	if s, ok := RTMGRPStrings[g]; ok {
		return s
	}
	return fmt.Sprintf("RTMGRP_UNKNOWN(%d)", g)
}

func NlMsgHdrStr(hdr *syscall.NlMsghdr) string {
	return fmt.Sprintf("%d %s", hdr.Len, NlMsgTypeStr(hdr.Type))
}

func NlMsgStr(nlmsg *syscall.NetlinkMessage) string {
	return fmt.Sprintf("%s %d", NlMsgHdrStr(&nlmsg.Header), len(nlmsg.Data))
}

var RT_SCOPE_strings = map[uint8]string{
	syscall.RT_SCOPE_UNIVERSE: "UNIVERSE",
	syscall.RT_SCOPE_SITE:     "SITE",
	syscall.RT_SCOPE_LINK:     "LINK",
	syscall.RT_SCOPE_HOST:     "HOST",
	syscall.RT_SCOPE_NOWHERE:  "NOWHERE",
}

func ScopeStr(scope uint8) string {
	if s, ok := RT_SCOPE_strings[scope]; ok {
		return s
	}
	return fmt.Sprintf("UNSPEC(%d)", scope)
}
