// -*- coding: utf-8 -*-

package templates

import (
	"io"
	"text/template"
)

const playbookInventory = `
[hosts]
localhost

[{{ .Name }}-hosts]
{{- range .Hosts }}
{{ . }}
{{- end }}
`

func NewPlaybookInventoryTemplate() *template.Template {
	return template.Must(template.New("inventory").Parse(playbookInventory))
}

type PlaybookInventory struct {
	Name  string
	Hosts []string
}

func NewPlaybookInventory() *PlaybookInventory {
	return &PlaybookInventory{
		Hosts: []string{},
	}
}

func (p *PlaybookInventory) AddHost(host string) {
	p.Hosts = append(p.Hosts, host)
}

func (p *PlaybookInventory) AddHosts(hosts ...string) {
	for _, host := range hosts {
		p.AddHost(host)
	}
}

func (p *PlaybookInventory) Execute(w io.Writer) error {
	return NewPlaybookInventoryTemplate().Execute(w, p)
}
