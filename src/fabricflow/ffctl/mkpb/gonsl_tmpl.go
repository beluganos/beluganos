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
	"text/template"
)

const playbookGonsldYaml = `---

dpaths:
  default:
    dpid: {{ .DpID }}
    addr: {{ .FibcAddr }}
    port: {{ .FibcPort }}
    # fibc_type: tcp

    l2sw:
      aging_sec: {{ .L2SWAgingSec }}
      sweep_sec: {{ .L2SWSweepSec }}
      notify_limit: {{ .L2SWNotifyLimit }}

    block_bcast:
      range: { min: {{ .L3PortStart }}, max: {{ .L3PortEnd }}, base_vid: {{ .L3VlanBase }} }
`

func NewPlaybookGonsldYamlTemplate() *template.Template {
	return template.Must(template.New("gonsld.yaml").Parse(playbookGonsldYaml))
}

type PlaybookGonsldYaml struct {
	DpID            uint64
	FibcAddr        string
	FibcPort        uint16
	L2SWAgingSec    uint32
	L2SWSweepSec    uint32
	L2SWNotifyLimit uint32
	L3PortStart     uint32
	L3PortEnd       uint32
	L3VlanBase      uint16
}

func NewPlaybookGonsldYaml() *PlaybookGonsldYaml {
	return &PlaybookGonsldYaml{
		DpID:            1234,
		FibcAddr:        "localhost",
		FibcPort:        50070,
		L2SWAgingSec:    3600,
		L2SWSweepSec:    3,
		L2SWNotifyLimit: 250,
		L3PortStart:     1,
		L3PortEnd:       190,
		L3VlanBase:      3900,
	}
}

func (p *PlaybookGonsldYaml) Execute(w io.Writer) error {
	return NewPlaybookGonsldYamlTemplate().Execute(w, p)
}

const playbookGonsldConf = `# -*- coding: utf-8 -*-

BELHOME=/etc/beluganos
CFGFILE=${BELHOME}/gonsld.yaml
LOGFILE=/var/log/gonsld.log
PIDFILE=/var/run/gonsld.pid

# Debug
# BELHOME=/tmp
# LOGFILE=/tmp/gonsld.log
# PIDFILE=/tmp/gonsld.pid
# USESIM="--use-sim"

API_ADDR="{{ .Addr }}:{{ .Port }}"
DATAPATH="default"
DAEMON_MODE="--daemon"
`

func NewPlaybookGonsldConfTemplate() *template.Template {
	return template.Must(template.New("gonsld.conf").Parse(playbookGonsldConf))
}

type PlaybookGonsldConf struct {
	Addr string
	Port uint16
}

func NewPlaybookGonsldConf() *PlaybookGonsldConf {
	return &PlaybookGonsldConf{
		Addr: "",
		Port: 50061,
	}
}

func (p *PlaybookGonsldConf) Execute(w io.Writer) error {
	return NewPlaybookGonsldConfTemplate().Execute(w, p)
}
