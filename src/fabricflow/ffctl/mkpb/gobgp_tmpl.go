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

package mkpb

import (
	"io"
	"text/template"
)

const playbookGoBGPdConf = `# -*- coding: utf-8; mode:toml -*-

[global.config]
    as = {{.AS}}
    router-id = {{"\""}}{{.RouterID}}{{"\""}}

[zebra]
  [zebra.config]
    enabled = {{ .ZAPIEnable }}
    version = {{ .ZAPIVersion }}
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
	RouterID    string
	AS          uint32
	ZAPIVersion uint16
	ZAPIEnable  bool
}

func NewPlaybookGoBGPdConf() *PlaybookGoBGPdConf {
	return &PlaybookGoBGPdConf{
		RouterID:    "0.0.0.0",
		ZAPIVersion: 5,
	}
}

func (p *PlaybookGoBGPdConf) Execute(w io.Writer) error {
	return NewPlaybookGoBGPdConfTemplate().Execute(w, p)
}

const playbookGoBGPConf = `# -*- coding: utf-8 -*-

CONF_PATH = /etc/frr/gobgpd.conf
CONF_TYPE = toml
LOG_LEVEL = debug
PPROF_OPT = --pprof-disable
API_HOSTS = {{ .APIAddr }}:{{ .APIPort }}
`

func NewPlaybookGoBGPConfTemplate() *template.Template {
	return template.Must(template.New("gobgpd.conf").Parse(playbookGoBGPConf))
}

type PlaybookGoBGPConf struct {
	APIAddr string
	APIPort uint16
}

func NewPlaybookGoBGPConf() *PlaybookGoBGPConf {
	return &PlaybookGoBGPConf{
		APIAddr: "localhost",
		APIPort: 50051,
	}
}

func (p *PlaybookGoBGPConf) Execute(w io.Writer) error {
	return NewPlaybookGoBGPConfTemplate().Execute(w, p)
}
