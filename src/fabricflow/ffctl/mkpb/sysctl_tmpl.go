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
	"fmt"
	"io"
	"sort"
	"text/template"
)

const playbookSysctlConf = `# -*- coding: utf-8 -*-

# enable ip forwarding
net.ipv4.ip_forward = 1
net.ipv6.conf.all.forwarding = 1

# mpls label max size.
net.mpls.platform_labels = {{ .MplsLabel }}

# enable mpls protocol.
{{- range .Ifaces }}
net.mpls.conf.{{ .Ifname }}.input = 1
{{- end }}

# disable rp_filter
net.ipv4.conf.default.rp_filter = 0
net.ipv4.conf.all.rp_filter = 0

{{- range .Ifaces }}
net.ipv4.conf.{{ .Ifname }}.rp_filter = 0
{{- end }}

# socket buffer size
net.core.rmem_max={{ .SockBufSize }}
net.core.wmem_max={{ .SockBufSize }}
`

func NewPlaybookSysctlConfTemplate() *template.Template {
	return template.Must(template.New("sysctl.conf").Parse(playbookSysctlConf))
}

type PlaybookSysctlIface struct {
	Index uint32
	Vlan  uint16
}

func (p *PlaybookSysctlIface) Ifname() string {
	if p.Vlan == 0 {
		return fmt.Sprintf("%s%d", IfnamePrefix, p.Index)
	}
	return fmt.Sprintf("%s%d/%d", IfnamePrefix, p.Index, p.Vlan)
}

func NewPlaybookSysctlIface(index uint32, vlan uint16) *PlaybookSysctlIface {
	return &PlaybookSysctlIface{
		Index: index,
		Vlan:  vlan,
	}
}

type PlaybookSysctlConf struct {
	SockBufSize uint64
	MplsLabel   uint32
	Ifaces      []*PlaybookSysctlIface
}

func NewPlaybookSysctlConf() *PlaybookSysctlConf {
	return &PlaybookSysctlConf{
		Ifaces: []*PlaybookSysctlIface{},
	}
}

func (p *PlaybookSysctlConf) AddIface(index uint32, vlan uint16) {
	p.Ifaces = append(p.Ifaces, NewPlaybookSysctlIface(index, vlan))
}

func (p *PlaybookSysctlConf) sort() {
	sort.Slice(p.Ifaces, func(i, j int) bool {
		if p.Ifaces[i].Index == p.Ifaces[j].Index {
			return p.Ifaces[i].Vlan < p.Ifaces[j].Vlan
		}
		return p.Ifaces[i].Index < p.Ifaces[j].Index
	})
}

func (p *PlaybookSysctlConf) Execute(w io.Writer) error {
	p.sort()
	return NewPlaybookSysctlConfTemplate().Execute(w, p)
}
