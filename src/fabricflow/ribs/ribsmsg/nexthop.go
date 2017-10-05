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

package ribsmsg

import (
	"fmt"
	"net"
)

type Nexthop struct {
	Rt    string
	Addr  net.IP
	SrcId net.IP
}

func NewNexthop(rt string, addr net.IP, srcId net.IP) *Nexthop {
	return &Nexthop{
		Rt:    rt,
		Addr:  addr,
		SrcId: srcId,
	}
}

func (e *Nexthop) IsMic() bool {
	return e.Rt == ""
}

func (e *Nexthop) String() string {
	if e.IsMic() {
		return fmt.Sprintf("%s RT:- Src:%s", e.Addr, e.SrcId)
	} else {
		return fmt.Sprintf("%s RT:%s Src:%s", e.Addr, e.Rt, e.SrcId)
	}
}
