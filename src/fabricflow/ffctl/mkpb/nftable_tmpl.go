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

const playbookNFTableConf = `
flush ruleset

define ic_ssh = {{ .PortSSH }}
define ic_snmp = {{ .PortSNMP }}
define ic_snmptrap = {{ .PortSNMPTrap }}

table ip nat {
    chain prerouting {
        type nat hook prerouting priority -100; policy accept;
        {{- range .NATPreRoutingRules }}
        {{ . }}
        {{- end }}
    }

    chain postrouting {
        type nat hook postrouting priority 100; policy accept;
        {{- range .NATPostRoutingRules }}
        {{ . }}
        {{- end }}
    }
}

table inet filter {
    chain input {
        type filter hook input priority 0; policy accept;
    }

    chain forward {
        type filter hook forward priority 0; policy accept;
    }

    chain output {
        type filter hook output priority 0; policy accept;
    }
}
`

func NewPlaybookNFTablesConfTemplate() *template.Template {
	return template.Must(template.New("nftables.conf").Parse(playbookNFTableConf))
}

type PlaybookNFTableConf struct {
	PortSNMP     uint16
	PortSNMPTrap uint16
	PortSSH      uint16

	NATPreRoutingRules  []string
	NATPostRoutingRules []string
}

func NewPlaybookNFTableConf() *PlaybookNFTableConf {
	return &PlaybookNFTableConf{
		PortSNMP:     1161,
		PortSNMPTrap: 1162,
		PortSSH:      122,

		NATPreRoutingRules:  []string{},
		NATPostRoutingRules: []string{},
	}
}

func (c *PlaybookNFTableConf) AddNATPreRoutingRule(s string) {
	c.NATPreRoutingRules = append(c.NATPreRoutingRules, s)
}

func (c *PlaybookNFTableConf) AddNATPostRoutingRule(s string) {
	c.NATPostRoutingRules = append(c.NATPostRoutingRules, s)
}

func (p *PlaybookNFTableConf) sort() {
}

func (p *PlaybookNFTableConf) Execute(w io.Writer) error {
	p.sort()
	return NewPlaybookNFTablesConfTemplate().Execute(w, p)
}
