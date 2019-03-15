// -*- coding: utf-8 -*-

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
