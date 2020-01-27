// -*- coding: utf-8 -*-

// Copyright (C) 2020 Nippon Telegraph and Telephone Corporation.
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
	"text/template"
)

const playbookGovswdYaml = `---

dpaths:
  default:
    dpid: {{ .DpID }}
    ifaces:
      names:
      {{- range .IfNames }}
        - "{{ . }}"
      {{- else }}
        {{- " []" }}
      {{- end }}
      patterns:
      {{- range .Patterns }}
        - "{{ . }}"
      {{- else }}
        {{- " []" }}
      {{- end }}
      blacklist:
      {{- range .Blacklist }}
        - "{{ . }}"
      {{- else }}
        {{- " []" }}
      {{- end }}
`

func NewPlaybookGovswdYamlTemplate() *template.Template {
	return template.Must(template.New("govswd.yaml").Parse(playbookGovswdYaml))
}

type PlaybookGovswdYaml struct {
	DpID      uint64
	IfNames   []string
	Patterns  []string
	Blacklist []string
}

func NewPlaybookGovswdYaml() *PlaybookGovswdYaml {
	return &PlaybookGovswdYaml{
		DpID:      1111,
		IfNames:   []string{},
		Patterns:  []string{},
		Blacklist: []string{},
	}
}

func (p *PlaybookGovswdYaml) AddIfNames(ifnames ...string) {
	p.IfNames = append(p.IfNames, ifnames...)
}

func (p *PlaybookGovswdYaml) AddPatterns(patterns ...string) {
	p.Patterns = append(p.Patterns, patterns...)
}

func (p *PlaybookGovswdYaml) AddBlacklist(ifnames ...string) {
	p.Blacklist = append(p.Blacklist, ifnames...)
}

func (p *PlaybookGovswdYaml) Execute(w io.Writer) error {
	return NewPlaybookGovswdYamlTemplate().Execute(w, p)
}

const playbookGovswdConf = `# -*- coding: utf-8 -*-

CONFIG_FILE={{ .ConfigFile }}
FIBC_ADDR={{ .FibcAddr }}
FIBC_PORT={{ .FibcPort }}
# OPTS="-v"
# OPTS="--trace"
`

func NewPlaybookGovswdConfTemplate() *template.Template {
	return template.Must(template.New("govswd.conf").Parse(playbookGovswdConf))
}

type PlaybookGovswdConf struct {
	ConfigFile string
	FibcAddr   string
	FibcPort   uint16
}

func NewPlaybookGovswdConf() *PlaybookGovswdConf {
	return &PlaybookGovswdConf{
		ConfigFile: "/etc/beluganos/govswd.yaml",
		FibcAddr:   "localhost",
		FibcPort:   50070,
	}
}

func (p *PlaybookGovswdConf) Execute(w io.Writer) error {
	return NewPlaybookGovswdConfTemplate().Execute(w, p)
}
