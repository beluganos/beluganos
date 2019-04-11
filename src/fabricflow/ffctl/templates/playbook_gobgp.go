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
