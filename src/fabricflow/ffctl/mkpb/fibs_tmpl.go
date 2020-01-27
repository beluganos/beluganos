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
	"time"
)

const playbookSnmpproxydConf = `# -*- coding: utf-8 -*-

{{- if eq .SnmpproxydType "ifmon" }}
# snmproxyd-ifmond
SNMPPROXYD_ADDR={{ .SnmpproxydAddr }}:{{ .SnmpproxydTrapPort }}
IFNOTIFY_RESEND={{ .SnmpproxydIfResend }}
TRAP_RESEND_TIME={{ .SnmpproxydIfResend }}
{{- else if eq .SnmpproxydType "mib/trap" }}
# snmpproxyd-mib/trap
CONF=/etc/beluganos/snmpproxyd.yaml
LISTEN_MIB=:{{ .SnmpproxydSnmpPort }}
LISTEN_TRAP=:{{ .SnmpproxydTrapPort }}
SNMPD_ADDR=localhost:{{ .SnmpdPort }}
{{- else }}
{{- end }}

# debug
# DEBUG="-v"
DUMP_TABLE_TIME=0
DUMP_TABLE_MIB=/tmp/snmpproxyd_mib_tables
DUMP_TABLE_TRAP=/tmp/snmpproxyd_trap_tables
`

func NewPlaybookSnmpproxydConfTemplate() *template.Template {
	return template.Must(template.New("snmpproxyd.conf").Parse(playbookSnmpproxydConf))
}

type PlaybookSnmpproxydConf struct {
	SnmpproxydType     string
	SnmpproxydAddr     string
	SnmpproxydSnmpPort uint16
	SnmpproxydTrapPort uint16
	SnmpdPort          uint16
	SnmpproxydIfResend time.Duration
}

func NewPlaybookSnmpproxydConf() *PlaybookSnmpproxydConf {
	return &PlaybookSnmpproxydConf{
		SnmpproxydType:     "",
		SnmpproxydAddr:     "192.169.1.1",
		SnmpproxydSnmpPort: 161,
		SnmpproxydTrapPort: 162,
		SnmpdPort:          1161,
		SnmpproxydIfResend: 15 * time.Minute,
	}
}

func (p *PlaybookSnmpproxydConf) Execute(w io.Writer) error {
	return NewPlaybookSnmpproxydConfTemplate().Execute(w, p)
}

const playbookSnmpproxydYaml = `---

snmpproxy:
  default:
    oidmap:
    {{- range .OidMap }}
      - name:  {{ .Name }}
        oid:   {{ stroid .Oid }}
        local: {{ stroid .Local }}
    {{- end }}
    {{- range .ONLOidMap }}
      - name:  {{ .Name }}
        oid:   {{ stroid .Oid }}
        local: {{ stroid .Local }}
        proxy: {{ .Proxy }}
    {{- end }}

    ifmap:
      oidmap:
        min: 0
        max: 1023
      shift:
        min: 1024
        max: 2147483647

    trap2map:
{{- range .Trap2Map }}
      {{ .Ifname }}: {{ .Port }}
{{- else }}
      {{- " {}" }}
{{- end }}

    trap2sink:
    {{- range .Trap2Sink }}
      - addr: {{ . }}
    {{- else }}
      {{- " []" }}
    {{- end }}
`

func NewPlaybookSnmpproxydYamlTemplate() *template.Template {
	fm := template.FuncMap{
		"stroid": StrOid,
	}
	return template.Must(template.New("snmpproxyd.yaml").Funcs(fm).Parse(playbookSnmpproxydYaml))
}

type PlaybookSnmpproxydTrap2MapEntry struct {
	Index  uint32
	Ifname string
	Port   uint32
}

type PlaybookSnmpproxydYaml struct {
	Trap2Map  []*PlaybookSnmpproxydTrap2MapEntry
	OidMap    []*SnmpdOidEntry
	ONLOidMap []*SnmpdOidEntry
	Trap2Sink []string
}

func NewPlaybookSnmpproxydYaml() *PlaybookSnmpproxydYaml {
	return &PlaybookSnmpproxydYaml{
		Trap2Map:  []*PlaybookSnmpproxydTrap2MapEntry{},
		OidMap:    snmpOidMap,
		ONLOidMap: []*SnmpdOidEntry{},
		Trap2Sink: []string{},
	}
}

func (p *PlaybookSnmpproxydYaml) AddTrap2Map(pport, lport uint32) {
	e := &PlaybookSnmpproxydTrap2MapEntry{
		Index:  pport,
		Ifname: NewPhyIfname(pport),
		Port:   lport,
	}
	p.Trap2Map = append(p.Trap2Map, e)
}

func (p *PlaybookSnmpproxydYaml) sort() {
	sort.Slice(p.Trap2Map, func(i, j int) bool { return p.Trap2Map[i].Index < p.Trap2Map[j].Index })
}

func (p *PlaybookSnmpproxydYaml) Execute(w io.Writer) error {
	p.sort()
	return NewPlaybookSnmpproxydYamlTemplate().Execute(w, p)
}

const playbookFibssnmpConf = `---

handlers:
{{- range .OidMap }}
  - oid:  {{ stroid .Local }}
    name: {{ .Name }}
    type: {{ .SnmpType }}
{{- end }}
`

func NewPlaybookFibssnmpYamlTemplate() *template.Template {
	fm := template.FuncMap{
		"stroid": StrOid,
	}
	return template.Must(template.New("fibssnmp.yaml").Funcs(fm).Parse(playbookFibssnmpConf))
}

type PlaybookFibssnmpYaml struct {
	OidMap []*SnmpdOidEntry
}

func NewPlaybookFibssnmpYaml() *PlaybookFibssnmpYaml {
	return &PlaybookFibssnmpYaml{
		OidMap: snmpOidMap,
	}
}

func (p *PlaybookFibssnmpYaml) Execute(w io.Writer) error {
	return NewPlaybookFibssnmpYamlTemplate().Execute(w, p)
}
