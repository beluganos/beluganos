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

const playbookSnmpdConf = `###############################################################################
#
# EXAMPLE.conf:
#   An example configuration file for configuring the Net-SNMP agent ('snmpd')
#   See the 'snmpd.conf(5)' man page for details
#
#  Some entries are deliberately commented out, and will need to be explicitly activated
#
###############################################################################
#
#  AGENT BEHAVIOUR
#

#  Listen for connections from the local system only
{{- if .SnmpPort }}
agentAddress udp:{{ .SnmpPort }}
{{- else }}
agentAddress udp:161
{{- end }}
#  Listen for connections on all interfaces (both IPv4 *and* IPv6)
#agentAddress udp:161,udp6:[::1]:161

###############################################################################
#
#  SNMPv3 AUTHENTICATION
#
#  Note that these particular settings don't actually belong here.
#  They should be copied to the file /var/lib/snmp/snmpd.conf
#     and the passwords changed, before being uncommented in that file *only*.
#  Then restart the agent

#  createUser authOnlyUser  MD5 "remember to change this password"
#  createUser authPrivUser  SHA "remember to change this one too"  DES
#  createUser internalUser  MD5 "this is only ever used internally, but still change the password"

#  If you also change the usernames (which might be sensible),
#  then remember to update the other occurances in this example config file to match.
###############################################################################
#
#  ACCESS CONTROL
#

                                                 #  system + hrSystem groups only
view   systemonly  included   .1.3.6.1.2.1.1
view   systemonly  included   .1.3.6.1.2.1.25.1

                                                 #  Full access from the local host
#rocommunity public  localhost
                                                 #  Default access to basic system info
 rocommunity public  default    -V systemonly
                                                 #  rocommunity6 is for IPv6
 rocommunity6 public  default   -V systemonly

                                                 #  Full access from an example network
                                                 #     Adjust this network address to match your local
                                                 #     settings, change the community string,
                                                 #     and check the 'agentAddress' setting above
#rocommunity secret  10.0.0.0/16

                                                 #  Full read-only access for SNMPv3
 rouser   authOnlyUser
                                                 #  Full write access for encrypted requests
                                                 #     Remember to activate the 'createUser' lines above
#rwuser   authPrivUser   priv

#  It's no longer typically necessary to use the full 'com2sec/group/access' configuration
#  r[ow]user and r[ow]community, together with suitable views, should cover most requirements



###############################################################################
#
#  SYSTEM INFORMATION
#

#  Note that setting these values here, results in the corresponding MIB objects being 'read-only'
#  See snmpd.conf(5) for more details
sysLocation    Sitting on the Dock of the Bay
sysContact     Me <me@example.org>
                                                 # Application + End-to-End layers
sysServices    72
#
#  Process Monitoring
#
                               # At least one  'mountd' process
proc  mountd
                               # No more than 4 'ntalkd' processes - 0 is OK
proc  ntalkd    4
                               # At least one 'sendmail' process, but no more than 10
proc  sendmail 10 1

#  Walk the UCD-SNMP-MIB::prTable to see the resulting output
#  Note that this table will be empty if there are no "proc" entries in the snmpd.conf file


#
#  Disk Monitoring
#
                               # 10MBs required on root disk, 5% free on /var, 10% free on all other disks
disk       /     10000
disk       /var  5%
includeAllDisks  10%

#  Walk the UCD-SNMP-MIB::dskTable to see the resulting output
#  Note that this table will be empty if there are no "disk" entries in the snmpd.conf file


#
#  System Load
#
                               # Unacceptable 1-, 5-, and 15-minute load averages
load   12 10 5

#  Walk the UCD-SNMP-MIB::laTable to see the resulting output
#  Note that this table *will* be populated, even without a "load" entry in the snmpd.conf file



###############################################################################
#
#  ACTIVE MONITORING
#

                                    #   send SNMPv1  traps
 trapsink     localhost public
                                    #   send SNMPv2c traps
#trap2sink    localhost public
{{- if .Trap2SinkAddr }}
trap2sink {{ .Trap2SinkAddr }} public
{{- end }}
                                    #   send SNMPv2c INFORMs
#informsink   localhost public

#  Note that you typically only want *one* of these three lines
#  Uncommenting two (or all three) will result in multiple copies of each notification.

#
#  Event MIB - automatically generate alerts
#
                                   # Remember to activate the 'createUser' lines above
iquerySecName   internalUser       
rouser          internalUser
                                   # generate traps on UCD error conditions
# defaultMonitors          yes
                                   # generate traps on linkUp/Down
# linkUpDownNotifications  yes

{{- if .LinkMonitorInterval }}
notificationEvent  linkUpTrap    linkUp   ifIndex ifAdminStatus ifOperStatus
notificationEvent  linkDownTrap  linkDown ifIndex ifAdminStatus ifOperStatus
monitor -r {{ .LinkMonitorInterval }} -e linkUpTrap   "Generate linkUp"   ifOperStatus != 2
monitor -r {{ .LinkMonitorInterval }} -e linkDownTrap "Generate linkDown" ifOperStatus == 2
{{- end }}

###############################################################################
#
#  EXTENDING THE AGENT
#

#
#  Arbitrary extension commands
#
 extend    test1   /bin/echo  Hello, world!
 extend-sh test2   echo Hello, world! ; echo Hi there ; exit 35
#extend-sh test3   /bin/sh /tmp/shtest

#  Note that this last entry requires the script '/tmp/shtest' to be created first,
#    containing the same three shell commands, before the line is uncommented

#  Walk the NET-SNMP-EXTEND-MIB tables (nsExtendConfigTable, nsExtendOutput1Table
#     and nsExtendOutput2Table) to see the resulting output

#  Note that the "extend" directive supercedes the previous "exec" and "sh" directives
#  However, walking the UCD-SNMP-MIB::extTable should still returns the same output,
#     as well as the fuller results in the above tables.


#
#  "Pass-through" MIB extension command
#
#pass .1.3.6.1.4.1.8072.2.255  /bin/sh       PREFIX/local/passtest
#pass .1.3.6.1.4.1.8072.2.255  /usr/bin/perl PREFIX/local/passtest.pl

# Note that this requires one of the two 'passtest' scripts to be installed first,
#    before the appropriate line is uncommented.
# These scripts can be found in the 'local' directory of the source distribution,
#     and are not installed automatically.

#  Walk the NET-SNMP-PASS-MIB::netSnmpPassExamples subtree to see the resulting output


#
#  AgentX Sub-agents
#
                                           #  Run as an AgentX master agent
 master          agentx
                                           #  Listen for network connections (from localhost)
                                           #    rather than the default named socket /var/agentx/master
#agentXSocket    tcp:localhost:705

{{- if .OidMap }}
# Beluganos ext-mibs.
{{- $cmdpath := .FibssnmpCmdPath }}
{{- $datapath := .FibssnmpDataPath }}
{{- range .OidMap }}
pass_persist {{ stroid .Local }} {{ $cmdpath }} --data-path {{ $datapath }}
{{- end }}
{{- end }}

{{- if .ONLOidMap }}
# Proxy to WB-SW
{{- range .ONLOidMap }}
proxy -c public -v 2c {{ .Proxy }} {{ stroid .Oid }} {{ stroid .Local }}
{{- end }}
{{- end }}
`

func NewPlaybookSnmpdConfTemplate() *template.Template {
	fm := template.FuncMap{
		"stroid": StrOid,
	}
	return template.Must(template.New("snmpd.conf").Funcs(fm).Parse(playbookSnmpdConf))
}

type PlaybookSnmpdConf struct {
	SnmpPort            uint16
	Trap2SinkAddr       string
	LinkMonitorInterval uint32
	FibssnmpCmdPath     string
	FibssnmpDataPath    string
	OidMap              []*SnmpdOidEntry
	ONLOidMap           []*SnmpdOidEntry
}

func NewPlaybookSnmpdConf() *PlaybookSnmpdConf {
	return &PlaybookSnmpdConf{
		SnmpPort:            0,
		Trap2SinkAddr:       "",
		LinkMonitorInterval: 0,
		FibssnmpCmdPath:     "/usr/bin/fibssnmp",
		FibssnmpDataPath:    "/etc/beluganos/fibc_stats.yaml",
		OidMap:              []*SnmpdOidEntry{},
		ONLOidMap:           []*SnmpdOidEntry{},
	}
}

func (p *PlaybookSnmpdConf) Execute(w io.Writer) error {
	return NewPlaybookSnmpdConfTemplate().Execute(w, p)
}
