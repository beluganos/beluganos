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
	"fabricflow/fibc/api"
	log "github.com/sirupsen/logrus"
	"gonla/nlamsg"
)

func NewPortConfig(cmd string, reId string, link *nlamsg.Link) *fibcapi.PortConfig {
	return fibcapi.NewPortConfig(cmd, reId, NewLinkLinkName(link), NewPortId(link))
}

func (r *RIBController) SendPortConfig(cmd string, link *nlamsg.Link) error {
	p := NewPortConfig(cmd, r.reId, link)
	return r.fib.Send(p, 0)
}

func (r *RIBController) SendPortConfigs() {
	err := r.nla.GetLinks(func(link *nlamsg.Link) error {
		return r.SendPortConfig("ADD", link)
	})

	if err != nil {
		log.Errorf("RIBController: SendPortConfig error. %s", err)
		return
	}

	log.Infof("RIBController: PortConfigs sent.")
}
