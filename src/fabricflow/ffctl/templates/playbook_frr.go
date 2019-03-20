// -*- coding: utf-8 -*-

package templates

import (
	"fmt"
	"io"
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
{{ range .Ifnames -}}
interface {{.}}
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

type PlaybookFrrConf struct {
	RouterID string
	Ifnames  []string
}

func NewPlaybookFrrConf(routerID string) *PlaybookFrrConf {
	return &PlaybookFrrConf{
		RouterID: routerID,
		Ifnames:  []string{},
	}
}

func (p *PlaybookFrrConf) AddIface(index, vid uint) {
	if vid == 0 {
		p.Ifnames = append(p.Ifnames, fmt.Sprintf("%s%d", playbookIfnamePrefix, index))
	} else {
		p.Ifnames = append(p.Ifnames, fmt.Sprintf("%s%d.%d", playbookIfnamePrefix, index, vid))
	}
}

func (p *PlaybookFrrConf) Execute(w io.Writer) error {
	return NewPlaybookFrrConfTemplate().Execute(w, p)
}
