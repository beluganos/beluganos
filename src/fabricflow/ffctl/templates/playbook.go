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

const playbookIfnamePrefix = "eth"

const playbook = `---

- hosts: hosts
  connection: local
  roles:
    - { role: lxd, mode: host }
  tags:
    - host

- hosts: hosts
  connection: local
  tasks:
    - include_role:
        name: lxd
      vars:
        mode: create
      with_items:
        - {{"\"{{ groups['"}}{{.Name}}-hosts{{"'] }}\""}}
      loop_control:
        loop_var: lxcname

- hosts: {{.Name}}-hosts
  connection: lxd
  roles:
    - { role: lxd, lxcname: {{"\"{{ inventory_hostname }}\""}}, mode: setup }
  tags:
    - setup
    - lxd

`

func NewPlaybookTemplate() *template.Template {
	return template.Must(template.New("playbook").Parse(playbook))
}

type Playbook struct {
	Name string
}

func NewPlaybook(name string) *Playbook {
	return &Playbook{
		Name: name,
	}
}

func (p *Playbook) Execute(w io.Writer) error {
	return NewPlaybookTemplate().Execute(w, p)
}
