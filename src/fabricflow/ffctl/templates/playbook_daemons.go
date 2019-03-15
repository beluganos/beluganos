// -*- coding: utf-8 -*-

package templates

import (
	"io"
	"text/template"
)

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

{{ range . -}}
{{ .Name }}={{ .Arg }}
{{  end -}}

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

func (p *PlaybookDaemons) Set(daemon, arg string) {
	p.daemons = append(p.daemons, NewPlaybookDaemon(daemon, arg))
}

func (p *PlaybookDaemons) SetMap(m map[string]string) {
	for daemon, yesno := range m {
		p.Set(daemon, yesno)
	}
}

func (p *PlaybookDaemons) SetYes(daemons ...string) {
	for _, daemon := range daemons {
		p.Set(daemon, "yes")
	}
}

func (p *PlaybookDaemons) SetNo(daemons ...string) {
	for _, daemon := range daemons {
		p.Set(daemon, "no")
	}
}

func (p *PlaybookDaemons) Execute(w io.Writer) error {
	return NewPlaybookDaemonsTemplate().Execute(w, p.daemons)
}
