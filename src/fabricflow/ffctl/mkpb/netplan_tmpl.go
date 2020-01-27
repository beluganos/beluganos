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
	"sort"
	"text/template"
)

const playbookNetplanYaml = `---

network:
  version: 2
  ethernets:
  {{- range .Eths }}
    {{ .Name }}:
      mtu: {{ .Mtu }}
  {{- else }}
    {{- " {}" }}
  {{- end }}

  {{- if .Vlans }}
  vlans:
  {{- range .Vlans }}
    {{ .Link }}.{{ .Vid }}:
      link: {{ .Link }}
      id: {{ .Vid  }}
  {{- else }}
    {{- " {}" }}
  {{- end }}
{{- end }}
`

func NewPlaybookNetplanYamlTemplate() *template.Template {
	return template.Must(template.New("netplan.yaml").Parse(playbookNetplanYaml))
}

type PlaybookNetplanEth struct {
	Index uint32
	Name  string
	Mtu   uint16
}

func NewPlaybookNetplanEth(index uint32, mtu uint16) *PlaybookNetplanEth {
	return &PlaybookNetplanEth{
		Index: index,
		Name:  NewPhyIfname(index),
		Mtu:   mtu,
	}
}

type PlaybookNetplanVlan struct {
	Index uint32
	Link  string
	Vid   uint16
}

func NewPlaybookNetplanVlan(index uint32, vid uint16) *PlaybookNetplanVlan {
	return &PlaybookNetplanVlan{
		Index: index,
		Link:  NewPhyIfname(index),
		Vid:   vid,
	}
}

type PlaybookNetplanYaml struct {
	Eths  []*PlaybookNetplanEth
	Vlans []*PlaybookNetplanVlan
}

func NewPlaybookNetplanYaml() *PlaybookNetplanYaml {
	return &PlaybookNetplanYaml{
		Eths:  []*PlaybookNetplanEth{},
		Vlans: []*PlaybookNetplanVlan{},
	}
}

func (p *PlaybookNetplanYaml) AddEth(index uint32, mtu uint16) {
	p.Eths = append(p.Eths, NewPlaybookNetplanEth(index, mtu))
}

func (p *PlaybookNetplanYaml) AddVlan(index uint32, vid uint16) {
	p.Vlans = append(p.Vlans, NewPlaybookNetplanVlan(index, vid))
}

func (p *PlaybookNetplanYaml) sort() {
	sort.Slice(p.Eths, func(i, j int) bool {
		return p.Eths[i].Index < p.Eths[j].Index
	})
	sort.Slice(p.Vlans, func(i, j int) bool {
		if p.Vlans[i].Index == p.Vlans[j].Index {
			return p.Vlans[i].Vid < p.Vlans[j].Vid
		}
		return p.Vlans[i].Index < p.Vlans[j].Index
	})
}

func (p *PlaybookNetplanYaml) Execute(w io.Writer) error {
	p.sort()
	return NewPlaybookNetplanYamlTemplate().Execute(w, p)
}
