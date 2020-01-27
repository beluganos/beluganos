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

const playbookBrVlanYaml = `---

network:
  vlans:
  {{- range .Vlans }}
    {{ .Name }};
      id: {{ .Id }}
      ids: {{ stru16 .Ids }}
  {{- else }}
    {{- " {}" }}
  {{- end }}

  bridges:
    {{ .Bridge }}:
      vlan_filtering: 1
      interfaces:
      {{- range .Vlans }}
        - {{ .Name }}
      {{- else }}
        {{- " []" }}
      {{- end }}
`

func NewPlaybookBrVlanYamlTemplate() *template.Template {
	fm := template.FuncMap{
		"stru16": StrUint16Slice,
		"stru32": StrUint32Slice,
	}

	return template.Must(template.New("bridge_vlan.yaml").Funcs(fm).Parse(playbookBrVlanYaml))
}

type PlaybookBrVlanVlan struct {
	Index uint32
	Name  string
	Id    uint16
	Ids   []uint16
}

func NewPlaybookBrVlanVlan(ifname string, index uint32) *PlaybookBrVlanVlan {
	return &PlaybookBrVlanVlan{
		Index: index,
		Name:  ifname,
		Ids:   []uint16{},
	}
}

func (p *PlaybookBrVlanVlan) sort() {
	sort.Slice(p.Ids, func(i, j int) bool { return p.Ids[i] < p.Ids[j] })
}

type PlaybookBrVlanYaml struct {
	Bridge string
	Vlans  []*PlaybookBrVlanVlan
}

func NewPlaybookBrVlanYaml() *PlaybookBrVlanYaml {
	return &PlaybookBrVlanYaml{
		Bridge: "l2swbr0",
		Vlans:  []*PlaybookBrVlanVlan{},
	}
}

func (p *PlaybookBrVlanYaml) AddAccessPorts(ports map[uint32]uint16) {
	if ports == nil {
		return
	}

	for index, vlan := range ports {
		vlan := &PlaybookBrVlanVlan{
			Index: index,
			Name:  NewPhyIfname(index),
			Id:    vlan,
			Ids:   []uint16{},
		}

		p.Vlans = append(p.Vlans, vlan)
	}
}

func (p *PlaybookBrVlanYaml) AddTrunkPorts(ports map[uint32][]uint16) {
	if ports == nil {
		return
	}

	for index, vlans := range ports {
		vlan := &PlaybookBrVlanVlan{
			Index: index,
			Name:  NewPhyIfname(index),
			Id:    0,
			Ids:   vlans,
		}

		p.Vlans = append(p.Vlans, vlan)
	}
}

func (p *PlaybookBrVlanYaml) sort() {
	for _, vlan := range p.Vlans {
		vlan.sort()
	}
	sort.Slice(p.Vlans, func(i, j int) bool { return p.Vlans[i].Index < p.Vlans[j].Index })
}

func (p *PlaybookBrVlanYaml) Execute(w io.Writer) error {
	p.sort()
	return NewPlaybookBrVlanYamlTemplate().Execute(w, p)
}
