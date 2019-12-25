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

const playbookRibxdConf = `# -*- coding: utf-8; mode: toml -*-

[node]
nid   = {{ .NID }}
reid  = "{{ .ReID }}"
label = 100000
allow_duplicate_ifname = false

[log]
level = {{ .LogLevel }}
dump  = {{ .LogDump }}

[nla]
core  = "{{ .Mic }}:{{ .NLACorePort }}"
api   = "127.0.0.1:{{ .NLAAPIPort }}"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
recv_chan_size = {{ .RecvChanSize }}
recv_sock_buf = {{ .RecvSockBufSize }}
=======
recv_chan_size = {{ .NLARecvChanSize }}
recv_sock_buf = {{ .NLARecvSockBufSize }}
>>>>>>> develop
=======
recv_chan_size = {{ .NLARecvChanSize }}
recv_sock_buf = {{ .NLARecvSockBufSize }}
>>>>>>> develop
=======
recv_chan_size = {{ .NLARecvChanSize }}
recv_sock_buf = {{ .NLARecvSockBufSize }}
>>>>>>> develop

  [[nla.iptun]]
  nid = 0
  remotes = []
  # remotes = ["10.100.1.0/24", "10.100.2.0/24"]
  # remotes = ["2001:db8:1::/64", "2001:db8:2::/64"]

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> develop
=======
>>>>>>> develop
  [nla.bridge_vlan]
  update_sec = {{ .NLABrVlanUpdateSec }}
  chan_size = {{ .NLABrVlanChanSize }}

<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop

[ribc]
{{- if eq .NID 0 }}
fibc  = "{{ .FibcAPIAddr }}:{{ .FibcAPIPort }}"
# fibc_type = "tcp"
{{- else }}
disable = true
{{- end }}

{{- if .Vpn }}
{{/* MPLS-VPN */}}
[ribs]
core = "{{ .Mic }}:{{ .RibsCorePort }}"
api  = "127.0.0.1:{{ .RibsAPIPort }}"

  [ribs.bgpd]
  addr = "127.0.0.1"

  {{- if eq .NID 0 }}
  {{/* MPLS-VPN (MIC) */}}
  [ribs.nexthops]
  mode = "translate"
  args = "{{ .VpnNexthop }}"
  {{- else }}
  {{/* MPLS-VPN (RIB) */}}
  [ribs.vrf]
  rt = "{{ .RT }}"
  rd = "{{ .RD }}"
  iface = "{{ .VpnNexhopBridge }}"
  {{- end }}
{{- else }}
{{/* NOT VPN */}}
[ribs]
disable = true
{{- end }}

[ribp]
api = "127.0.0.1:{{ .RibpAPIPort }}"
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
	Mic  string

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	RecvChanSize    uint64
	RecvSockBufSize uint64
	VpnNexthop      string // x.x.x.x/y
	VpnNexhopBridge string
=======
=======
>>>>>>> develop
=======
>>>>>>> develop
	NLARecvChanSize    uint64
	NLARecvSockBufSize uint64
	NLABrVlanUpdateSec uint32
	NLABrVlanChanSize  uint32
	VpnNexthop         string // x.x.x.x/y
	VpnNexhopBridge    string
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop

	NLACorePort  uint16
	NLAAPIPort   uint16
	FibcAPIAddr  string
	FibcAPIPort  uint16
	RibsCorePort uint16
	RibsAPIPort  uint16
	RibpAPIPort  uint16

	LogLevel uint8
	LogDump  uint8
}

func NewPlaybookRibxdConf() *PlaybookRibxdConf {
	return &PlaybookRibxdConf{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		RecvChanSize:    65536,
		RecvSockBufSize: 8388608,
=======
=======
>>>>>>> develop
=======
>>>>>>> develop
		NLARecvChanSize:    65536,
		NLARecvSockBufSize: 1024 * 1024 * 8,
		NLABrVlanUpdateSec: 60 * 30,
		NLABrVlanChanSize:  4096 * 4,

<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
		VpnNexthop:      "1.1.0.0/24",
		VpnNexhopBridge: "ffbr0",

		NLACorePort:  50061,
		NLAAPIPort:   50062,
		FibcAPIAddr:  "192.169.1.1",
		FibcAPIPort:  50070,
		RibsCorePort: 50071,
		RibsAPIPort:  50072,
		RibpAPIPort:  50091,
	}
}

func (p *PlaybookRibxdConf) Execute(w io.Writer) error {
	return NewPlaybookRibxdConfTemplate().Execute(w, p)
}
