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

package fibccmd

import (
	"fmt"
)

const (
	MOD_NAME_VM_MODIFY = "portcfg/vm"
	MOD_NAME_VS_MODIFY = "portcfg/vs"
	MOD_NAME_DP_MODIFY = "portcfg/dp"
	MOD_NAME_ADD       = "portmap/port/add"
	MOD_NAME_DELETE    = "portmap/port/delete"
)

type ModEntry struct {
	Name string                   `yaml:"name" json:"name"`
	Cmd  string                   `yaml:"cmd" json:"cmd"`
	ReId string                   `yaml:"re_id" json:"re_id"`
	VsId uint64                   `yaml:"vs_id" json:"vs_id"`
	DpId uint64                   `yaml:"dp_id" json:"dp_id"`
	Args []map[string]interface{} `yaml:"args" json:"args"`
}

type ModClient struct {
	Addr string
}

func NewModClient(addr string) *ModClient {
	return &ModClient{
		Addr: addr,
	}
}

func (m *ModClient) Send(e *ModEntry) error {
	url := fmt.Sprintf("http://%s/fib/%s", m.Addr, e.Name)
	return httpJson(url, e)
}

func (m *ModClient) Sends(mods []*ModEntry) error {
	for _, e := range mods {
		if err := m.Send(e); err != nil {
			return err
		}
	}
	return nil
}
