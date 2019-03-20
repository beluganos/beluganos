// -*- coding: utf-8 -*-

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
