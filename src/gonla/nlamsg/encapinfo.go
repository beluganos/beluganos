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
	"net"
)

//
// EncapInfo
//
type EncapInfo struct {
	Dst  *net.IPNet
	Vrf  uint32
	EnId uint32
}

func (e *EncapInfo) Copy() *EncapInfo {
	d := *e.Dst
	return NewEncapInfo(&d, e.Vrf, e.EnId)
}

func (e *EncapInfo) String() string {
	return fmt.Sprintf("Dst: %s Vrf: %d EnId: %d", e.Dst, e.Vrf, e.EnId)
}

func NewEncapInfo(dst *net.IPNet, vrf uint32, enId uint32) *EncapInfo {
	return &EncapInfo{
		Dst:  dst,
		Vrf:  vrf,
		EnId: enId,
	}
}
