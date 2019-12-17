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
	"fmt"
	"io"
	"sort"
	"text/template"
)

const playbookFrrConf = `! -*- coding: utf-8 -*-
!
frr defaults datacenter
username cumulus nopassword
service integrated-vtysh-config
!
ipv6 forwarding
ip forwarding
!
no log monitor
log file /var/log/frr/frr.log
log syslog informational
log timestamp precision 6
!
{{ range .Ifaces -}}
interface {{ .Ifname }}
{{ end -}}
!
interface lo
  ip address {{.RouterID}}/32
!
router-id {{.RouterID}}
!
router ospf
  router-id {{.RouterID}}
!
router ospf6
  router-id {{.RouterID}}
!
mpls ldp
  router-id {{.RouterID}}
!
line vty
!
`

func NewPlaybookFrrConfTemplate() *template.Template {
	return template.Must(template.New("frr.conf").Parse(playbookFrrConf))
}

type PlaybookFrrConfIface struct {
	Index uint32
	Vlan  uint16
}

func (p *PlaybookFrrConfIface) Ifname() string {
	if p.Vlan == 0 {
		return fmt.Sprintf("%s%d", IfnamePrefix, p.Index)
	}

	return fmt.Sprintf("%s%d.%d", IfnamePrefix, p.Index, p.Vlan)
}

type PlaybookFrrConf struct {
	RouterID string
	Ifaces   []*PlaybookFrrConfIface
}

func NewPlaybookFrrConf(routerID string) *PlaybookFrrConf {
	return &PlaybookFrrConf{
		RouterID: routerID,
		Ifaces:   []*PlaybookFrrConfIface{},
	}
}

func (p *PlaybookFrrConf) AddIface(index uint32, vlan uint16) {
	iface := &PlaybookFrrConfIface{
		Index: index,
		Vlan:  vlan,
	}
	p.Ifaces = append(p.Ifaces, iface)
}

func (p *PlaybookFrrConf) sort() {
	sort.Slice(p.Ifaces, func(i, j int) bool {
		if p.Ifaces[i].Index == p.Ifaces[j].Index {
			return p.Ifaces[i].Vlan < p.Ifaces[j].Vlan
		}

		return p.Ifaces[i].Index < p.Ifaces[j].Index
	})
}

func (p *PlaybookFrrConf) Execute(w io.Writer) error {
	p.sort()
	return NewPlaybookFrrConfTemplate().Execute(w, p)
}

const playbookDaemons = `
 This file tells the frr package which daemons to start.
#
# Entries are in the format: <daemon>=(yes|no|priority)
#   0, "no"  = disabled
#   1, "yes" = highest priority
#   2 .. 10  = lower priorities
# Read /usr/share/doc/frr/README.Debian for details.
#
# Sample configurations for these daemons can be found in
# /usr/share/doc/frr/examples/.
#
# ATTENTION:
#
# When activation a daemon at the first time, a config file, even if it is
# empty, has to be present *and* be owned by the user and group "frr", else
# the daemon will not be started by /etc/init.d/frr. The permissions should
# be u=rw,g=r,o=.
# When using "vtysh" such a config file is also needed. It should be owned by
# group "frrvty" and set to ug=rw,o= though. Check /etc/pam.d/frr, too.
#
# The watchfrr daemon is always started. Per default in monitoring-only but
# that can be changed via /etc/frr/debian.conf.
#

{{- range . }}
{{ .Name }}={{ .Arg }}
{{-  end }}
`

func NewPlaybookDaemonsTemplate() *template.Template {
	return template.Must(template.New("daemons").Parse(playbookDaemons))
}

type PlaybookDaemon struct {
	Name string
	Arg  string
}

func NewPlaybookDaemon(name, arg string) *PlaybookDaemon {
	return &PlaybookDaemon{
		Name: name,
		Arg:  arg,
	}
}

type PlaybookDaemons struct {
	daemons []*PlaybookDaemon
}

func NewPlaybookDaemons() *PlaybookDaemons {
	return &PlaybookDaemons{
		daemons: []*PlaybookDaemon{},
	}
}

func (p *PlaybookDaemons) SetMap(m map[string]string) {
	for name, arg := range m {
		daemon := &PlaybookDaemon{
			Name: name,
			Arg:  arg,
		}
		p.daemons = append(p.daemons, daemon)
	}
}

func (p *PlaybookDaemons) Execute(w io.Writer) error {
	return NewPlaybookDaemonsTemplate().Execute(w, p.daemons)
}
