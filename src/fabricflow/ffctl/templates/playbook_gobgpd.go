// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package templates

import (
	"io"
	"text/template"
)

const playbookGoBGPdConf = `# -*- coding: utf-8; mode: toml -*-

[global.config]
    as = {{.AS}}
    router-id = {{"\""}}{{.RouterID}}{{"\""}}

[zebra]
  [zebra.config]
    enabled = false
    version = 5
    url = "unix:/var/run/frr/zserv.api"
    # redistribute-route-type-list = ["connect"]

# [[neighbors]]
#   [neighbors.config]
#     neighbor-address = "10.0.0.2"
#     peer-as = 65000
#   [[neighbors.afi-safis]]
#     [neighbors.afi-safis.config]
#       afi-safi-name = "ipv4-unicast"
#   [neighbors.apply-policy.config]
#     export-policy-list = ["policy-nexthop-self"]
#     default-export-policy = "accept-route"
#
# [[policy-definitions]]
#   name = "policy-nexthop-self"
#   [[policy-definitions.statements]]
#     [policy-definitions.statements.actions.bgp-actions]
#       set-next-hop = "self"
#     [policy-definitions.statements.actions]
#       route-disposition = "accept-route"
`

func NewPlaybookGoBGPdConfTemplate() *template.Template {
	return template.Must(template.New("gobgpd.conf").Parse(playbookGoBGPdConf))
}

type PlaybookGoBGPdConf struct {
	RouterID string
	AS       uint
}

func NewPlaybookGoBGPdConf() *PlaybookGoBGPdConf {
	return &PlaybookGoBGPdConf{
		RouterID: "0.0.0.0",
	}
}

func (p *PlaybookGoBGPdConf) Execute(w io.Writer) error {
	return NewPlaybookGoBGPdConfTemplate().Execute(w, p)
}
