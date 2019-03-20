// -*- coding: utf-8 -*-

package templates

import (
	"io"
	"text/template"
)

const playbookRibxdConf = `# -*- coding: utf-8; mode: toml -*-

[node]
nid   = {{ .NID }}
reid  = "{{ .ReID }}"
label = 100000
allow_duplicate_ifname = false

[log]
level = 5
dump  = 0

[nla]
core  = "{{ .Name }}:50061"
api   = "127.0.0.1:50062"
recv_chan_size = 65536
recv_sock_buf = 8388608

[ribc]
{{- if eq .NID 0 }}
fibc  = "192.169.1.1:50070"
{{- else }}
disable = true
{{- end }}

[ribs]
{{- if .Vpn }}
core = "{{ .Name }}:50071"
api  = "127.0.0.1:50072"

[ribs.bgpd]
addr = "127.0.0.1"
{{ if eq .NID 0 }}
[ribs.nexthops]
mode = "translate"
args = "1.1.0.0/24"
{{- else }}
[ribs.vrf]
rt = "{{ .RT }}"
rd = "{{ .RD }}"
iface = "ffbr0"
{{- end }}
{{- else }}
disable = true
{{- end }}

[ribp]
api = "127.0.0.1:50091"
interval = 5000
`

func NewPlaybookRibxdConfTemplate() *template.Template {
	return template.Must(template.New("ribxd.conf").Parse(playbookRibxdConf))
}

type PlaybookRibxdConf struct {
	NID  uint8
	ReID string
	Name string
	RT   string
	RD   string
	Vpn  bool
}

func NewPlaybookRibxdConf() *PlaybookRibxdConf {
	return &PlaybookRibxdConf{}
}

func (p *PlaybookRibxdConf) Execute(w io.Writer) error {
	return NewPlaybookRibxdConfTemplate().Execute(w, p)
}
