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

const playbookSnmpProxydConf = `# -*- coding: utf-8 -*-

{{ if .Host -}}
# snmpproxyd-mib/trap
CONF=/etc/beluganos/snmpproxyd.yaml
LISTEN_MIB=:161
LISTEN_TRAP=:162
SNMPD_ADDR=localhost:8161
{{- else -}}
# snmproxyd-ifmond
SNMPPROXYD_ADDR={{ .SnmpProxydAddr }}:161
TRAP_RESEND_TIME=10s
{{- end }}

# debug
# DEBUG="-v"
DUMP_TABLE_TIME=0
DUMP_TABLE_MIB=/tmp/snmpproxyd_mib_tables
DUMP_TABLE_TRAP=/tmp/snmpproxyd_trap_tables
`

const playbookSnmpProxydYaml = `---

snmpproxy:
  default:
    oidmap:
      - name:  ifIndex
        oid:   .1.3.6.1.2.1.2.2.1.1
        local: .1.3.6.1.4.99999.2.2.1.1
      - name:  ifDescr
        oid:   .1.3.6.1.2.1.2.2.1.2
        local: .1.3.6.1.4.99999.2.2.1.2
      - name:  ifType
        oid:   .1.3.6.1.2.1.2.2.1.3
        local: .1.3.6.1.4.99999.2.2.1.3
      - name:  ifMtu
        oid:   .1.3.6.1.2.1.2.2.1.4
        local: .1.3.6.1.4.99999.2.2.1.4
      - name:  ifSpeed
        oid:   .1.3.6.1.2.1.2.2.1.5
        local: .1.3.6.1.4.99999.2.2.1.5
      - name:  ifPhysAddress
        oid:   .1.3.6.1.2.1.2.2.1.6
        local: .1.3.6.1.4.99999.2.2.1.6
      - name:  ifAdminStatus
        oid:   .1.3.6.1.2.1.2.2.1.7
        local: .1.3.6.1.4.99999.2.2.1.7
      - name:  ifOperStatus
        oid:   .1.3.6.1.2.1.2.2.1.8
        local: .1.3.6.1.4.99999.2.2.1.8
      - name:  ifLastChange
        oid:   .1.3.6.1.2.1.2.2.1.9
        local: .1.3.6.1.4.99999.2.2.1.9
      - name:  ifInOctets
        oid:   .1.3.6.1.2.1.2.2.1.10
        local: .1.3.6.1.4.99999.2.2.1.10
      - name:  ifInUcastPkts
        oid:   .1.3.6.1.2.1.2.2.1.11
        local: .1.3.6.1.4.99999.2.2.1.11
      - name:  ifInNUcastPkts
        oid:   .1.3.6.1.2.1.2.2.1.12
        local: .1.3.6.1.4.99999.2.2.1.12
      - name:  ifInDiscards
        oid:   .1.3.6.1.2.1.2.2.1.13
        local: .1.3.6.1.4.99999.2.2.1.13
      - name:  ifInErrors
        oid:   .1.3.6.1.2.1.2.2.1.14
        local: .1.3.6.1.4.99999.2.2.1.14
      - name:  ifInUnknownProtos
        oid:   .1.3.6.1.2.1.2.2.1.15
        local: .1.3.6.1.4.99999.2.2.1.15
      - name:  ifOutOctets
        oid:   .1.3.6.1.2.1.2.2.1.16
        local: .1.3.6.1.4.99999.2.2.1.16
      - name:  ifOutUcastPkts
        oid:   .1.3.6.1.2.1.2.2.1.17
        local: .1.3.6.1.4.99999.2.2.1.17
      - name:  ifOutNUcastPkts
        oid:   .1.3.6.1.2.1.2.2.1.18
        local: .1.3.6.1.4.99999.2.2.1.18
      - name:  ifOutDiscards
        oid:   .1.3.6.1.2.1.2.2.1.19
        local: .1.3.6.1.4.99999.2.2.1.19
      - name:  ifOutErrors
        oid:   .1.3.6.1.2.1.2.2.1.20
        local: .1.3.6.1.4.99999.2.2.1.20
      - name:  ifOutQLen
        oid:   .1.3.6.1.2.1.2.2.1.21
        local: .1.3.6.1.4.99999.2.2.1.21
      - name:  ifSpecific
        oid:   .1.3.6.1.2.1.2.2.1.22
        local: .1.3.6.1.4.99999.2.2.1.22
      - name:  ifName
        oid:   .1.3.6.1.2.1.31.1.1.1.1
        local: .1.3.6.1.4.99999.31.1.1.1.1

    ifmap:
      oidmap:
        min: 0
        max: 1023
      shift:
        min: 1024
        max: 2147483647

    trap2map:
{{- range $name, $port := .Trap2Map }}
      {{ $name }}: {{ $port }}
{{- else }}
      {{- "{}" }}
{{- end }}

    trap2sink:
      - addr: 192.168.122.1:161
`

func NewPlaybookSnmpProxydConfTemplate() *template.Template {
	return template.Must(template.New("snmpproxyd.conf").Parse(playbookSnmpProxydConf))
}

type PlaybookSnmpProxydConf struct {
	Host           bool
	SnmpProxydAddr string
}

func NewPlaybookSnmpProxydConf(host bool) *PlaybookSnmpProxydConf {
	return &PlaybookSnmpProxydConf{
		Host: host,
	}
}

func (p *PlaybookSnmpProxydConf) Execute(w io.Writer) error {
	return NewPlaybookSnmpProxydConfTemplate().Execute(w, p)
}

func NewPlaybookSnmpProxydYamlTemplate() *template.Template {
	return template.Must(template.New("snmpproxyd.yaml").Parse(playbookSnmpProxydYaml))
}

type PlaybookSnmpProxydYaml struct {
	Trap2Map map[string]uint
}

func NewPlaybookSnmpProxydYaml() *PlaybookSnmpProxydYaml {
	return &PlaybookSnmpProxydYaml{
		Trap2Map: map[string]uint{},
	}
}

func (p *PlaybookSnmpProxydYaml) AddTrap2Map(pport, lport uint) {
	ifname := fmt.Sprintf("%s%d", playbookIfnamePrefix, pport)
	p.Trap2Map[ifname] = lport
}

func (p *PlaybookSnmpProxydYaml) Execute(w io.Writer) error {
	return NewPlaybookSnmpProxydYamlTemplate().Execute(w, p)
}
