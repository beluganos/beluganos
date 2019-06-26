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

const playbookBrVlanYaml = `---

network:
  vlans:
{{- range .Vlans }}
    {{ .Name }};
      id: {{ .Id }}
      ids: {{ .Ids }}
{{- end }}

  bridges:
    l2swbr0:
      vlan_filtering: 1
      interfaces:
{{- range .Vlans }}
        - {{ .Name }}
{{- end}}
`

func NewPlaybookBrVlanYamlTemplate() *template.Template {
	return template.Must(template.New("bridge_vlan.yaml").Parse(playbookBrVlanYaml))
}

type PlaybookBrVlanVlan struct {
	Name string
	Id   uint
	Ids  []uint
}

func NewPlaybookBrVlanVlan(ifname string) *PlaybookBrVlanVlan {
	return &PlaybookBrVlanVlan{
		Name: ifname,
		Ids:  []uint{},
	}
}

type PlaybookBrVlanYaml struct {
	Vlans []*PlaybookBrVlanVlan
}

func NewPlaybookBrVlanYaml() *PlaybookBrVlanYaml {
	return &PlaybookBrVlanYaml{
		Vlans: []*PlaybookBrVlanVlan{},
	}
}

func (p *PlaybookBrVlanYaml) find(ifname string) *PlaybookBrVlanVlan {
	for _, vlan := range p.Vlans {
		if vlan.Name == ifname {
			return vlan
		}
	}

	return nil
}

func (p *PlaybookBrVlanYaml) AddAccessPort(index uint, vid uint) {
	name := fmt.Sprintf("%s%d", playbookIfnamePrefix, index)
	if p.find(name) == nil {
		vlan := NewPlaybookBrVlanVlan(name)
		vlan.Id = vid
		p.Vlans = append(p.Vlans, vlan)
	}
}

func (p *PlaybookBrVlanYaml) AddTrunkPort(index uint, vids ...uint) {
	name := fmt.Sprintf("%s%d", playbookIfnamePrefix, index)
	vlan := p.find(name)
	if vlan == nil {
		vlan = NewPlaybookBrVlanVlan(name)
		p.Vlans = append(p.Vlans, vlan)
	}

	vlan.Ids = append(vlan.Ids, vids...)
}

func (p *PlaybookBrVlanYaml) Execute(w io.Writer) error {
	return NewPlaybookBrVlanYamlTemplate().Execute(w, p)
}
