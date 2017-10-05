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

package nlalink

import (
	"syscall"
)

const (
	// RTM_SETLINK  -> netlink.RTM_SETLINK
	RTM_SETADDR = syscall.RTM_MAX + iota
	RTM_SETNEIGH
	RTM_SETROUTE
	RTM_NEWNODE
	RTM_DELNODE
	RTM_SETNODE
	RTM_NEWVPN
	RTM_DELVPN
	RTM_SETVPN
)

const (
	RTMGRP_UNSPEC = iota
	RTMGRP_NODE
	RTMGRP_VPN
	RTMGRP_LINK
	RTMGRP_ADDR
	RTMGRP_NEIGH
	RTMGRP_ROUTE
)
