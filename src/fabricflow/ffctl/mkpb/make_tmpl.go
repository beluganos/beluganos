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

const playbookSampleMakeYaml = `---

global:
  re-id: {{ .Global.ReID }}
  dp-id: {{ .Global.DpID }}
  dp-type: {{ .Global.DpType }}
  dp-mode: {{ .Global.DpMode }}
  dp-addr: {{ .Global.DpAddr }}
  vpn: {{ .Global.Vpn }}

router:
{{- range .Router }}
  - name: {{ .Name }}
    nid:  {{ .NodeID }}  # 0: MIC, >0: RIC
    eth:  {{ stru32 .Eth }}  # if empty, all ports.
    daemons: {{ strstr .Daemons }}
  {{- if len .RtRd }}
    rt-rd: {{ strstr .RtRd }} # [<RT>, <RD>]
  {{- end }}

  {{- if .Vlan }}
    Vlan:
    {{- range $index, $vlan := .Vlan }}
      {{ $index }}: {{ stru16 $vlan }}
    {{- else }}
      {{- " {}" }}
    {{- end }}
  {{- end }}

  {{- if .L2SW }}
    l2sw:
    {{- if .L2SW.Access }}
      access:
      {{- range $index, $vlan := .L2SW.Access }}
        {{ $index }}: {{ $vlan }}
      {{- else }}
        {{- " {}" }}
      {{- end }}
    {{- end }}
    {{- if .L2SW.Trunk }}
      trunk:
      {{- range $index, $vlans := .L2SW.Trunk }}
        {{ $index }}: {{ stru16 $vlans }}
      {{- else }}
        {{- " {}" }}
      {{- end }}
    {{- end }}
  {{- end }}

  {{- if .IPTun }}
    iptun:
      bgp-route-family: {{ .IPTun.BgpRouteFamily }}
      local-addr-range4: {{ .IPTun.LocalAddrRange4 }}
      local-addr-range6: {{ .IPTun.LocalAddrRange6 }}
      remote-routes:
      {{- range .IPTun.RemoteRoutes }}
        - {{ . }}
      {{- else }}
        {{- " []"}}
      {{- end }}
  {{- end }}
{{- end }}
`

func NewPlaybookSampleMakeYamlTempl() *template.Template {
	fm := template.FuncMap{
		"stru16": StrUint16Slice,
		"stru32": StrUint32Slice,
		"strstr": StrStringSlice,
	}
	return template.Must(template.New("make-sample.yaml").Funcs(fm).Parse(playbookSampleMakeYaml))
}

func ExecPlaybookSampleMakeYamlTempl(c *Config, w io.Writer) error {
	return NewPlaybookSampleMakeYamlTempl().Execute(w, c)
}

const playbookSampleMakeOptionYaml = `
# to use options, uncomment 'option:' and entries.
# option:
{{- range .List }}
#  {{ .Key }}: {{ .Val }}     # {{ printf "%T" .Val }}
{{- end }}
`

func NewPlaybookSampleMakeOptionYamlTempl() *template.Template {
	return template.Must(template.New("make-option.yaml").Parse(playbookSampleMakeOptionYaml))
}

func ExecPlaybookSampleMakeOptionYamlTempl(c *OptionConfig, w io.Writer) error {
	return NewPlaybookSampleMakeOptionYamlTempl().Execute(w, c)
}
