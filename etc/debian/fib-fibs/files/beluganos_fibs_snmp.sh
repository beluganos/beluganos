#! /bin/bash
# -*- coding: utf-8 -*-

# Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FIBSSNMP_CONF="/etc/beluganos/fibssnmp.yaml"
SNMPD_CONF="/etc/snmp/snmpd.conf"

# for debug
# FIBSSNMP_CONF="./fibssnmp.yaml"
# SNMPD_CONF="./snmpd.conf"

if [ -e "/etc/beluganos/snmpproxyd.conf" ]; then
    . /etc/beluganos/snmpproxyd.conf
else
    SNMPD_ADDR="127.0.0.1:1161"
fi
    
FIBS_AGENT_BIN="/usr/bin/fibssnmp -v"

oid_list() {
    grep oid ${FIBSSNMP_CONF} | sed -r /'[ ]*#'/d | sed -r 's/[ ]*- oid[\ ]*:[ ]*//g'
}

oid_register() {
    local OIDLIST=`oid_list`
    local OID

    for OID in ${OIDLIST}; do
        echo "pass_persist ${OID} ${FIBS_AGENT_BIN}" >> ${SNMPD_CONF}
    done
}

oid_unregister() {
    local OIDLIST=`oid_list`
    local OID

    local TEMP_FILE=`tempfile`

    for OID in ${OIDLIST}; do
        sed -e /"^pass_persist ${OID}\ .*"/d ${SNMPD_CONF} > ${TEMP_FILE}
        mv -f ${TEMP_FILE} ${SNMPD_CONF}
    done
}

oid_backup() {
    cp ${SNMPD_CONF} ${SNMPD_CONF}.bak
}

agentaddr_update() {
    local TEMP_FILE=`tempfile`
    local UDP_ADDR=$1

    sed -r "s/^agentAddress\s+udp:.+/agentAddress udp:${UDP_ADDR}/g" ${SNMPD_CONF} > ${TEMP_FILE}
    mv ${TEMP_FILE} ${SNMPD_CONF}
}

snmpd_conf_recover() {
    cp ${SNMPD_CONF}.bak ${SNMPD_CONF}
}

usage() {
    echo "$0 <param>"
    echo " param:"
    echo "  - register:   register ex-mib."
    echo "  - unregister: unregister ex-mib."
    echo "  - replace:    replace ex-mib."
}

case $1 in
    register)
        oid_backup
        oid_register
	agentaddr_update ${SNMPD_ADDR}
        ;;
    unregister)
        oid_backup
        oid_unregister
	agentaddr_update "localhost:161"
        ;;
    replace)
        oid_backup
        oid_unregister
        oid_register
	agentaddr_update ${SNMPD_ADDR}
        ;;
    oid-list)
	oid_list
	;;
    agentaddr)
	oid_backup
	agentaddr_update ${SNMPD_ADDR}
	;;
    recover)
	snmpd_conf_recover
        ;;
    *) usage;;
esac
