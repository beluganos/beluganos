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
	"net"
	"text/template"
	"time"
)

const playbookRibtdConf = `# -*- coding: utf-8 -*-

# monitor gobgp rib of this route family.
ROUTE_FAMILY={{ .BgpRouteFamily }}

# one of lo interface addresses contained in TUNNEL_LOCAL[4/6]
# is used as tunnel local address.
TUNNEL_LOCAL4={{ .TunnelLocal4 }}
TUNNEL_LOCAL6={{ .TunnelLocal6 }}

$ GoBGP API
API_LISTEN_ADDR={{ .GoBGPAPIAddr }}:{{ .GoBGPAPIPort }}

# DO NOT EDIT.
TUNNEL_PREFIX={{ .TunnelIFPrefix }}
TUNNEL_TYPE_IPV6={{ .TunnelTypeIPv6 }}
TUNNEL_TYPE_FORCE={{ .TunnelTypeForce }}
TUNNEL_TYPE_DEFAULT={{ .TunnelTypeDefault }}

# debug
# DEBUG="-v"
DUMP_TABLE_TIME={{ .DumpTableDuration }}
`

func NewPlaybookRibtdConfTemplate() *template.Template {
	return template.Must(template.New("ribtd.conf").Parse(playbookRibtdConf))
}

type PlaybookRibtdConf struct {
	GoBGPAPIAddr string
	GoBGPAPIPort uint16

	BgpRouteFamily string

	TunnelLocalRange4 string
	TunnelLocalRange6 string
	TunnelIFPrefix    string
	TunnelTypeIPv6    uint16
	TunnelTypeForce   uint16
	TunnelTypeDefault uint16

	DumpTableDuration time.Duration
}

func NewPlaybookRibtdConf() *PlaybookRibtdConf {
	return &PlaybookRibtdConf{
		GoBGPAPIAddr: "localhost",
		GoBGPAPIPort: 50051,

		BgpRouteFamily:    "ipv4-unicast",
		TunnelLocalRange4: "127.0.0.1/32",
		TunnelLocalRange6: "::1/128",
		TunnelIFPrefix:    "tun",
		TunnelTypeIPv6:    14,

		DumpTableDuration: 0,
	}
}

func (p *PlaybookRibtdConf) TunnelLocal4() string {
	if _, ipnet, err := net.ParseCIDR(p.TunnelLocalRange4); err == nil {
		if ipnet.IP.To4() != nil {
			return ipnet.String()
		}
	}

	return "127.0.0.1/32"
}

func (p *PlaybookRibtdConf) TunnelLocal6() string {
	if _, ipnet, err := net.ParseCIDR(p.TunnelLocalRange6); err == nil {
		if ipnet.IP.To4() == nil {
			return ipnet.String()
		}
	}

	return "::1/128"
}

func (p *PlaybookRibtdConf) Execute(w io.Writer) error {
	return NewPlaybookRibtdConfTemplate().Execute(w, p)
}
