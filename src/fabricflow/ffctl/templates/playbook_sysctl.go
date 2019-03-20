// -*- coding: utf-8 -*-

package templates

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"
)

const playbookSysctlConf = `# -*- coding: utf-8 -*-

# enable ip forwarding
net.ipv4.ip_forward = 1
net.ipv6.conf.all.forwarding = 1

# mpls label max size.
net.mpls.platform_labels={{ .MplsLabel }}

# enable mpls protocol.
{{- range .Ifaces }}
net.mpls.conf.{{ . }}.input = 1
{{- end }}

# disable rp_filter
net.ipv4.conf.default.rp_filter = 0
net.ipv4.conf.all.rp_filter=0

{{- range .Ifaces }}
net.ipv4.conf.{{ . }}.rp_filter = 0
{{- end }}

# socket buffer size
net.core.rmem_max={{ .SockBufSize }}
net.core.wmem_max={{ .SockBufSize }}
`

func NewPlaybookSysctlConfTemplate() *template.Template {
	return template.Must(template.New("sysctl.conf").Parse(playbookSysctlConf))
}

type PlaybookSysctlConf struct {
	SockBufSize uint
	MplsLabel   uint
	ifaces      []string
}

func NewPlaybookSysctlConf() *PlaybookSysctlConf {
	return &PlaybookSysctlConf{
		ifaces: []string{},
	}
}

func (p *PlaybookSysctlConf) AddIface(index, vid uint) {
	if vid == 0 {
		p.ifaces = append(p.ifaces, fmt.Sprintf("%s%d", playbookIfnamePrefix, index))
	} else {
		p.ifaces = append(p.ifaces, fmt.Sprintf("%s%d/%d", playbookIfnamePrefix, index, vid))
	}
}

func (p *PlaybookSysctlConf) Ifaces() []string {
	sort.Slice(p.ifaces, func(i, j int) bool { return strings.Compare(p.ifaces[i], p.ifaces[j]) < 0 })
	return p.ifaces
}

func (p *PlaybookSysctlConf) Execute(w io.Writer) error {
	return NewPlaybookSysctlConfTemplate().Execute(w, p)
}
