// -*- coding: utf-8 -*-

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
SNMPPROXYD_ADDR=192.169.1.1:161
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
      - name:  ifOperStatus
        oid:   .1.3.6.1.2.1.2.2.1.8
        local: .1.3.6.1.4.99999.2.2.1.8
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
      - addr: 192.169.122.1:161
`

func NewPlaybookSnmpProxydConfTemplate() *template.Template {
	return template.Must(template.New("snmpproxyd.conf").Parse(playbookSnmpProxydConf))
}

type PlaybookSnmpProxydConf struct {
	Host bool
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
