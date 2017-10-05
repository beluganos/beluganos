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

package ribsyn

import (
	"fabricflow/ribs/ribsmsg"
	"github.com/osrg/gobgp/table"
	"gonla/nlalib"
	"net"
)

//
// RIC Controller
//
type RicController struct {
	RibCh chan *ribsmsg.RibUpdate
	ConCh chan *nlalib.ConnInfo
}

func NewRicController() *RicController {
	return &RicController{
		RibCh: make(chan *ribsmsg.RibUpdate),
		ConCh: make(chan *nlalib.ConnInfo),
	}
}

func (c *RicController) Conn() <-chan *nlalib.ConnInfo {
	return c.ConCh
}

func (c *RicController) Recv() <-chan *ribsmsg.RibUpdate {
	return c.RibCh
}

func (c *RicController) Start(addr string) error {
	return Tables.Rics.FindByAddr(addr).Start(c.ConCh)
}

func (c *RicController) Stop(addr string) {
	Tables.Rics.FindByAddr(addr).Stop()
}

func (c *RicController) Monitor(addr net.IP) error {
	return Tables.Rics.FindByAddr(addr.String()).MonitorRib(c.RibCh)
}

func (c *RicController) Translate(rib *ribsmsg.RibUpdate, f func(*table.Path, *table.Path, *RicEntry) error) error {
	return Tables.Rics.FindByRt(rib.Rt).Translate(rib, f)
}

func (c *RicController) GetRibs(rt string) {
	Tables.Rics.FindByRt(rt).GetRibs()
}
