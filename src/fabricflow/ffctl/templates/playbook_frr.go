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
	"fmt"
	"io"
	"text/template"
)

const playbookFrrConf = `! -*- coding: utf-8 -*-
!
frr defaults datacenter
username cumulus nopassword
service integrated-vtysh-config
!
ipv6 forwarding
ip forwarding
!
no log monitor
log file /var/log/frr/frr.log
log syslog informational
log timestamp precision 6
!
{{ range .Ifnames -}}
interface {{.}}
{{ end -}}
!
interface lo
  ip address {{.RouterID}}/32
!
router-id {{.RouterID}}
!
router ospf
  router-id {{.RouterID}}
!
router ospf6
  router-id {{.RouterID}}
!
mpls ldp
  router-id {{.RouterID}}
!
line vty
!
`

func NewPlaybookFrrConfTemplate() *template.Template {
	return template.Must(template.New("frr.conf").Parse(playbookFrrConf))
}

type PlaybookFrrConf struct {
	RouterID string
	Ifnames  []string
}

func NewPlaybookFrrConf(routerID string) *PlaybookFrrConf {
	return &PlaybookFrrConf{
		RouterID: routerID,
		Ifnames:  []string{},
	}
}

func (p *PlaybookFrrConf) AddIface(index, vid uint) {
	if vid == 0 {
		p.Ifnames = append(p.Ifnames, fmt.Sprintf("%s%d", playbookIfnamePrefix, index))
	} else {
		p.Ifnames = append(p.Ifnames, fmt.Sprintf("%s%d.%d", playbookIfnamePrefix, index, vid))
	}
}

func (p *PlaybookFrrConf) Execute(w io.Writer) error {
	return NewPlaybookFrrConfTemplate().Execute(w, p)
}
