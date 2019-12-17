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

const playbookRibtdConf = `# -*- coding: utf-8 -*-

ROUTE_FAMILY=ipv4-unicast
# ROUTE_FAMILY=ipv6-unicast

TUNNEL_LOCAL4=127.0.0.1/32
# TUNNEL_LOCAL4=10.100.0.0/24

TUNNEL_LOCAL6=::/128
# TUNNEL_LOCAL6=2001:db8::/64

$ GoBGP API
API_LISTEN_ADDR=localhost:50051

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
