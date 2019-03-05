#! /bin/bash
#! -*- coding: utf-8 -*-

TEMPFILE=`tempfile`

do_oid_list() {
    local CONFPATH=$1

    grep oid ${CONFPATH} | sed -r 's/- oid ://g'
}

do_install() {
    local CONFPATH=$1
    local OID=$2
    local CMDS=$3

    echo "append $OID to $CONFPATH"
    echo "pass_persist ${OID} ${CMDS}" >> $TEMPFILE
}

do_uninstall() {
    local CONFPATH=$1
    local OID=$2

    echo "remove $OID from $CONFPATH"
    sed -e /"^pass_persist ${OID}\ .*"/d $CONFPATH > $TEMPFILE
}

do_set_agent_port() {
    local CONFPATH=$1
    local PORT=$2

    echo "set snmpd listen port ${PORT}"
    sed -r -e "s/^agentAddress\s+udp:(.+):161$/agentAddress udp:\1:${PORT}/" $CONFPATH > $TEMPFILE
}

do_unset_agent_port() {
    local CONFPATH=$1
    local PORT=$2

    echo "unset snmpd listen port ${PORT}"
    sed -r -e "s/^agentAddress\s+udp:(.+):${PORT}$/agentAddress udp:\1:161/" $CONFPATH > $TEMPFILE
}

do_enable_rocommunity() {
    local CONFPATH=$1

    echo "enable rocommunity public localhost"
    sed -r -e "s/^#\s*rocommunity\s+public\s+localhost/rocommunity public  localhost/" $CONFPATH > $TEMPFILE
}

do_disable_rocommunity() {
    local CONFPATH=$1

    echo "disable rocommunity public localhost"
    sed -r -e "s/^\s*rocommunity\s+public\s+localhost/\#rocommunity public  localhost/" $CONFPATH > $TEMPFILE
}

do_commit() {
    local CONFPATH=$1

    install -Tpm 644 $TEMPFILE $CONFPATH
    rm -f $TEMPFILE
}

do_show() {
    local CONFPATH=$1

    echo "--- $CONFPATH <-> $TEMPFILE ---"
    diff $CONFPATH $TEMPFILE
}

do_usage() {
    echo "$0 <install|uninstall> <conf file> <oid> \"<command>\""
    echo "$0 <set-agent-port|unset-agent-port> <conf file> <port>"
}

case $1 in
    oid-list)
        do_oid_list $2
        ;;
    install)
        do_uninstall $2 $3
        do_install $2 $3 "$4"
        do_commit $2
        ;;
    uninstall)
        do_uninstall $2 $3
        do_commit $2
        ;;
    set-agent-port)
        do_set_agent_port $2 $3
        do_commit $2
        ;;
    unset-agent-port)
        do_unset_agent_port $2 $3
        do_commit $2
        ;;
    enable-rocommunity)
        do_enable_rocommunity $2
        do_commit $2
        ;;
    disable-rocommunity)
        do_disable_rocommunity $2
        do_commit $2
        ;;
    *)
        do_usage
        ;;
esac

