// -*- coding: utf-8 -*-

package templates

import (
	"io"
	"text/template"
)

const playbookGoBGPdConf = `# -*- coding: utf-8; mode: toml -*-

[global.config]
    as = {{.AS}}
    router-id = {{"\""}}{{.RouterID}}{{"\""}}

[zebra]
  [zebra.config]
    enabled = false
    version = 5
    url = "unix:/var/run/frr/zserv.api"
    # redistribute-route-type-list = ["connect"]

# [[neighbors]]
#   [neighbors.config]
#     neighbor-address = "10.0.0.2"
#     peer-as = 65000
#   [[neighbors.afi-safis]]
#     [neighbors.afi-safis.config]
#       afi-safi-name = "ipv4-unicast"
#   [neighbors.apply-policy.config]
#     export-policy-list = ["policy-nexthop-self"]
#     default-export-policy = "accept-route"
#
# [[policy-definitions]]
#   name = "policy-nexthop-self"
#   [[policy-definitions.statements]]
#     [policy-definitions.statements.actions.bgp-actions]
#       set-next-hop = "self"
#     [policy-definitions.statements.actions]
#       route-disposition = "accept-route"
`

func NewPlaybookGoBGPdConfTemplate() *template.Template {
	return template.Must(template.New("gobgpd.conf").Parse(playbookGoBGPdConf))
}

type PlaybookGoBGPdConf struct {
	RouterID string
	AS       uint
}

func NewPlaybookGoBGPdConf() *PlaybookGoBGPdConf {
	return &PlaybookGoBGPdConf{
		RouterID: "0.0.0.0",
	}
}

func (p *PlaybookGoBGPdConf) Execute(w io.Writer) error {
	return NewPlaybookGoBGPdConfTemplate().Execute(w, p)
}
