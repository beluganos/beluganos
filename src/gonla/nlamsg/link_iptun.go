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

type Iptun struct {
	*Link
	TnlId    uint16
	LocalMAC net.HardwareAddr
}

func NewIptun(link *Link) *Iptun {
	e := &Iptun{
		Link:     link.Copy(),
		LocalMAC: nil,
	}

	return e
}

func (e *Iptun) Copy() *Iptun {
	return &Iptun{
		Link:     e.Link.Copy(),
		TnlId:    e.TnlId,
		LocalMAC: e.LocalMAC,
	}
}

func (e *Iptun) Remote() net.IP {
	return e.Iptun().Remote
}

func (e *Iptun) Local() net.IP {
	return e.Iptun().Local
}

func (v *Iptun) String() string {
	return fmt.Sprintf("%s tnlId %d MAC %s", v.Link, v.TnlId, v.LocalMAC)
}
