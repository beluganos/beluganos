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

const playbookNetplanYaml = `---

network:
  version: 2
  ethernets:
{{- range .Eths }}
    {{ .Name }}:
      mtu: {{ .Mtu }}
{{- else }}
  {{- "{}" }}
{{- end }}

  vlans:
{{- range .Vlans }}
    {{ .Link }}.{{ .Vid }}:
      link: {{ .Link }}
      id: {{ .Vid  }}
{{- else }}
  {{- "{}" }}
{{- end }}
`

func NewPlaybookNetplanYamlTemplate() *template.Template {
	return template.Must(template.New("netplan.yaml").Parse(playbookNetplanYaml))
}

type PlaybookNetplanEth struct {
	Name string
	Mtu  uint16
}

func NewPlaybookNetplanEth(name string, mtu uint16) *PlaybookNetplanEth {
	return &PlaybookNetplanEth{
		Name: name,
		Mtu:  mtu,
	}
}

type PlaybookNetplanVlan struct {
	Link string
	Vid  uint
}

func NewPlaybookNetplanVlan(link string, vid uint) *PlaybookNetplanVlan {
	return &PlaybookNetplanVlan{
		Link: link,
		Vid:  vid,
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

func (p *PlaybookNetplanYaml) AddEth(index uint, mtu uint16) {
	name := fmt.Sprintf("%s%d", playbookIfnamePrefix, index)
	p.Eths = append(p.Eths, NewPlaybookNetplanEth(name, mtu))
}

func (p *PlaybookNetplanYaml) AddVlan(index uint, vid uint) {
	link := fmt.Sprintf("%s%d", playbookIfnamePrefix, index)
	p.Vlans = append(p.Vlans, NewPlaybookNetplanVlan(link, vid))
}

func (p *PlaybookNetplanYaml) Execute(w io.Writer) error {
	return NewPlaybookNetplanYamlTemplate().Execute(w, p)
}
