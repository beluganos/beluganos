// -*- coding: utf-8 -*-

package templates

import (
	"fmt"
	"io"
	"text/template"
)

const playbookLXDProfile = `---

{{ $name := .Name -}}
{{- $mtu := .Mtu -}}
- name: create profile
  lxd_profile:
    name: "{{ .Name }}"
    state: present
    config: {"security.privileged": "true"}
    devices:
      eth0: # Management LAN
        name: {{ .MngIface }}
        nictype: bridged
        parent: {{ .BridgeIface }}
        type: nic
      logdir:
        path: /var/log
        source: /var/log/beluganos/{{ .Name }}
        type: disk
      root:
        path: /
        pool: default
        type: disk
{{ range $i, $v := .Ifaces }}
      {{ $v }}:
        type: nic
        name: {{ $v }}
        host_name: "{{$name}}.{{ HostNameIndex $i}}"
        nictype: p2p
        mtu: "{{ $mtu }}"
{{- end }}
`

func NewPlaybookLXDProfileTemplate() *template.Template {
	return template.Must(
		template.New("lxd_profile").Funcs(
			template.FuncMap{
				"HostNameIndex": func(i int) int { return i + 1 },
			},
		).Parse(playbookLXDProfile))
}

type PlybookLXDProfile struct {
	Name        string
	MngIface    string
	BridgeIface string
	Mtu         uint16
	Ports       []uint
}

func (p *PlybookLXDProfile) IfaceIndex(i uint) uint {
	return i + 1
}

func (p *PlybookLXDProfile) Ifaces() []string {
	ifaces := []string{}
	for _, port := range p.Ports {
		ifaces = append(ifaces, fmt.Sprintf("%s%d", playbookIfnamePrefix, port))
	}

	return ifaces
}

func NewPlaybookLXDProfile() *PlybookLXDProfile {
	return &PlybookLXDProfile{
		Ports: []uint{},
	}
}

func (p *PlybookLXDProfile) AddPort(pport uint) {
	p.Ports = append(p.Ports, pport)
}

func (p *PlybookLXDProfile) AddPorts(pports ...uint) {
	for _, pport := range pports {
		p.AddPort(pport)
	}
}

func (p *PlybookLXDProfile) Execute(w io.Writer) error {
	return NewPlaybookLXDProfileTemplate().Execute(w, p)
}
