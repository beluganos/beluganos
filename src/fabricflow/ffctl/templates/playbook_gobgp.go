// -*- coding: utf-8 -*-

package templates

import (
	"io"
	"text/template"
)

const playbookGoBGPConf = `# -*- coding: utf-8 -*-

CONF_PATH = /etc/frr/gobgpd.conf
CONF_TYPE = toml
LOG_LEVEL = debug
PPROF_OPT = --pprof-disable
API_HOSTS = 127.0.0.1:50051
`

func NewPlaybookGoBGPConfTemplate() *template.Template {
	return template.Must(template.New("gobgpd.conf").Parse(playbookGoBGPConf))
}

type PlaybookGoBGPConf struct {
}

func NewPlaybookGoBGPConf() *PlaybookGoBGPConf {
	return &PlaybookGoBGPConf{}
}

func (p *PlaybookGoBGPConf) Execute(w io.Writer) error {
	return NewPlaybookGoBGPConfTemplate().Execute(w, p)
}
