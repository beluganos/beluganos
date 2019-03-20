// -*- coding: utf-8 -*-

package templates

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
	PPort uint
	LPort uint
}

func (p *PlaybookFibcPort) Ifname() string {
	return fmt.Sprintf("%s%d", playbookIfnamePrefix, p.PPort)
}

func NewPlaybookFibcPort(pport, lport uint) *PlaybookFibcPort {
	return &PlaybookFibcPort{
		PPort: pport,
		LPort: lport,
	}
}

type PlaybookFibcYaml struct {
	Desc   string
	ReID   string
	DpName string
	DpID   uint64
	DpMode string
	Ports  []*PlaybookFibcPort
}

func NewPlaybookFibcYaml(reID, dpName string) *PlaybookFibcYaml {
	return &PlaybookFibcYaml{
		Desc:   fmt.Sprintf("router(%s)", reID),
		ReID:   reID,
		DpName: dpName,
		DpID:   0,
		DpMode: "onsl",
		Ports:  []*PlaybookFibcPort{},
	}
}

func (p *PlaybookFibcYaml) AddPort(pport, lport uint) {
	p.Ports = append(p.Ports, NewPlaybookFibcPort(pport, lport))
}

func (p *PlaybookFibcYaml) AddPorts(m map[uint]uint) {
	var keys []uint
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, pport := range keys {
		p.AddPort(pport, m[pport])
	}
}

func (p *PlaybookFibcYaml) Execute(w io.Writer) error {
	return NewPlaybookFibcYamlTemplate().Execute(w, p)
}
