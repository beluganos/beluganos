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

const lxdProfile = `---

{{ $name := .Name -}}
{{- $mtu := .Mtu -}}
name: "{{ .Name }}"
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

func NewLXDProfileTemplate() *template.Template {
	return template.Must(
		template.New("lxd_profile").Funcs(
			template.FuncMap{
				"HostNameIndex": func(i int) int { return i + 1 },
			},
		).Parse(lxdProfile))
}

type PlybookLXDProfile struct {
	Name        string
	MngIface    string
	BridgeIface string
	Mtu         uint16
	Ports       []uint32
}

func (p *PlybookLXDProfile) IfaceIndex(i uint) uint {
	return i + 1
}

func (p *PlybookLXDProfile) Ifaces() []string {
	ifaces := []string{}
	for _, port := range p.Ports {
		ifaces = append(ifaces, fmt.Sprintf("%s%d", IfnamePrefix, port))
	}

	return ifaces
}

func NewPlaybookLXDProfile() *PlybookLXDProfile {
	return &PlybookLXDProfile{
		Ports: []uint32{},
	}
}

func (p *PlybookLXDProfile) AddPorts(pports []uint32) {
	for _, pport := range pports {
		p.Ports = append(p.Ports, pport)
	}
}

func (p *PlybookLXDProfile) sort() {
	sort.Slice(p.Ports, func(i, j int) bool { return p.Ports[i] < p.Ports[j] })
}

func (p *PlybookLXDProfile) Execute(w io.Writer, lxdConfigMode bool) error {
	p.sort()
	if lxdConfigMode {
		return NewLXDProfileTemplate().Execute(w, p)
	}
	return NewPlaybookLXDProfileTemplate().Execute(w, p)
}
