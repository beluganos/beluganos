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

const playbookSnmpConf = `# As the snmp packages come without MIB files due to license reasons, loading
# of MIBs is disabled by default. If you added the MIBs you can reenable
# loading them by commenting out the following line.
mibs +ALL
`

func NewPlaybookSnmpConfTemplate() *template.Template {
	return template.Must(template.New("snmp.conf").Parse(playbookSnmpConf))
}

type PlaybookSnmpConf struct {
}

func NewPlaybookSnmpConf() *PlaybookSnmpConf {
	return &PlaybookSnmpConf{}
}

func (p *PlaybookSnmpConf) Execute(w io.Writer) error {
	return NewPlaybookSnmpConfTemplate().Execute(w, p)
}
