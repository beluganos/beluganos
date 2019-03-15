// -*- coding: utf-8 -*-

package templates

import (
	"io"
	"text/template"
)

const playbookRibtdConf = `# -*- coding: utf-8 -*-

API_LISTEN_ADDR=localhost:50051
ROUTE_FAMILY=ipv4-unicast
TUNNEL_LOCAL4=127.0.0.1/32
TUNNEL_LOCAL6=::/128
TUNNEL_PREFIX=tun

TUNNEL_TYPE_IPV6=14
TUNNEL_TYPE_FORCE=0
TUNNEL_TYPE_DEFAULT=0

# debug
# DEBUG="-v"
DUMP_TABLE_TIME=0
`

func NewPlaybookRibtdConfTemplate() *template.Template {
	return template.Must(template.New("ribtd.conf").Parse(playbookRibtdConf))
}

type PlaybookRibtdConf struct {
}

func NewPlaybookRibtdConf() *PlaybookRibtdConf {
	return &PlaybookRibtdConf{}
}

func (p *PlaybookRibtdConf) Execute(w io.Writer) error {
	return NewPlaybookRibtdConfTemplate().Execute(w, p)
}
