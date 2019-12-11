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

package dhcplib

import (
	"fmt"
	"io"
	"net"
	"text/template"
)

const iscDhcpConfigTemplate = `
# /etc/dhcp/dhcpd.conf
{{range .ExtOptions -}}
option {{.Name}} code {{.Code}} = {{.Type}}
{{end}}
subnet {{.Subnet}} netmask {{.Mask}} {
# range {{.RangeBegin}} {{.RangeEnd}}
{{range .Options -}}
option {{.Value}};
{{end}}
}
`

func NewISCDhcpdConfigTemplate() *template.Template {
	return template.Must(template.New("isc-dhcpd.conf").Parse(iscDhcpConfigTemplate))
}

type ISCDhcpdOption struct {
	Name  string
	Value string
	Code  uint8
	Type  string
}

func NewISCDhcpdOption(name string, code uint8, value string, typ string) *ISCDhcpdOption {
	return &ISCDhcpdOption{
		Name:  name,
		Value: value,
		Code:  code,
		Type:  typ,
	}
}

type ISCDhcpdConfig struct {
	Options []*ISCDhcpdOption
	IPNet   *net.IPNet
}

func NewISCDhcpdConfig(addr string) (*ISCDhcpdConfig, error) {
	ip, ipnet, err := net.ParseCIDR(addr)
	if err != nil {
		return nil, err
	}

	c := &ISCDhcpdConfig{
		IPNet: ipnet,
		Options: []*ISCDhcpdOption{
			NewISCDhcpdOption(
				OptionNameRouter,
				OptionCodeRouter,
				fmt.Sprintf("%s %s", OptionNameRouter, ip),
				"",
			),

			NewISCDhcpdOption(
				OptionNameDefaultURL,
				OptionCodeDefaultURL,
				fmt.Sprintf("%s = \"http://%s/onie-installer\"", OptionNameDefaultURL, ip),
				"",
			),
			NewISCDhcpdOption(
				OptionNameBeluganosKVMURL,
				OptionCodeBeluganosKVMURL,
				fmt.Sprintf("%s = \"http://%s/beluganos-kvm-installer\"", OptionNameBeluganosKVMURL, ip),
				"text",
			),
			NewISCDhcpdOption(
				OptionNameBeluganosZTPURL,
				OptionCodeBeluganosZTPURL,
				fmt.Sprintf("%s = \"http://%s/beluganos-ztp-installer\"", OptionNameBeluganosZTPURL, ip),
				"text",
			),
		},
	}

	return c, nil
}

func (c *ISCDhcpdConfig) Execute(w io.Writer) error {
	return NewISCDhcpdConfigTemplate().Execute(w, c)
}

func (c *ISCDhcpdConfig) Subnet() string {
	return c.IPNet.IP.String()
}

func (c *ISCDhcpdConfig) Mask() string {
	return net.IP(c.IPNet.Mask).String()
}

func (c *ISCDhcpdConfig) RangeBegin() string {
	return c.IPNet.IP.String()
}

func (c *ISCDhcpdConfig) RangeEnd() string {
	return c.IPNet.IP.String()
}

func (c *ISCDhcpdConfig) ExtOptions() []*ISCDhcpdOption {
	opts := []*ISCDhcpdOption{}
	for _, opt := range c.Options {
		if len(opt.Type) != 0 {
			opts = append(opts, opt)
		}
	}

	return opts
}
