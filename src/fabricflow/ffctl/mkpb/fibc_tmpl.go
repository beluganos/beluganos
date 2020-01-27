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

const playbookFibcYaml = `---

routers:
  - desc: {{.Desc}}
    re_id: {{.ReID}}  # Router entity id.
    datapath: {{.DpName}}
    ports:
      {{- range .Ports }}
      - { name: {{ .Ifname }},  port: {{ .LPort }} }
      {{- end }}

#
datapaths:
  - name: {{.DpName}}
    dp_id: {{.DpID}}
    mode: {{.DpMode}}
`

func NewPlaybookFibcYamlTemplate() *template.Template {
	return template.Must(template.New("daemons").Parse(playbookFibcYaml))
}

type PlaybookFibcPort struct {
	PPort uint32
	LPort uint32
}

func (p *PlaybookFibcPort) Ifname() string {
	return NewPhyIfname(p.PPort)
}

type PlaybookFibcYaml struct {
	Desc   string
	ReID   string
	DpName string
	DpID   uint64
	DpMode string
	Ports  []*PlaybookFibcPort
}

func NewPlaybookFibcYaml(reID string) *PlaybookFibcYaml {
	return &PlaybookFibcYaml{
		Desc:  fmt.Sprintf("router(%s)", reID),
		ReID:  reID,
		Ports: []*PlaybookFibcPort{},
	}
}

func (p *PlaybookFibcYaml) AddPorts(ports map[uint32]uint32) {
	for pport, lport := range ports {
		port := &PlaybookFibcPort{
			PPort: pport,
			LPort: lport,
		}
		p.Ports = append(p.Ports, port)
	}
}

func (p *PlaybookFibcYaml) sort() {
	sort.Slice(p.Ports, func(i, j int) bool {
		return p.Ports[i].PPort < p.Ports[j].PPort
	})
}

func (p *PlaybookFibcYaml) Execute(w io.Writer) error {
	p.sort()
	return NewPlaybookFibcYamlTemplate().Execute(w, p)
}
