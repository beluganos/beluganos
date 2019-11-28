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

LXC_LIST="test"
ONSLD_ADDR="172.16.0.58:50061"

MK_TAR="yes"
PREFIX="test-lxc"

ech() {
    echo ""
    echo "$LXC> $*"
}

do_cmds_vtysh() {
    ech "hostname"
    hostname

    ech "date"
    date

    # config
    ech 'vtysh -c "show running-config"'
    vtysh -c "show running-config"

    # interface
    ech 'vtysh -c "show interface"'
    vtysh -c "show interface"

    # route
    ech 'vtysh -c "show ip route"'
    vtysh -c "show ip route"

    ech 'vtysh -c "show ipv6 route"'
    vtysh -c "show ipv6 route"

    # OSPF
    ech 'vtysh -c "show ip ospf neighbor"'
    vtysh -c "show ip ospf neighbor"

    ech 'vtysh -c "show ipv6 ospf6 neighbor"'
    vtysh -c "show ipv6 ospf6 neighbor"

    ech 'vtysh -c "show ip ospf route"'
    vtysh -c "show ip ospf route"

    ech 'vtysh -c "show ipv6 ospf6 route"'
    vtysh -c "show ipv6 ospf6 route"

    ech 'vtysh -c "show ip ospf database"'
    vtysh -c "show ip ospf database"

    ech 'vtysh -c "show ip6 ospf6 database"'
    vtysh -c "show ip6 ospf6 database"

    # LDP
    ech 'vtysh -c "show mpls ldp binding"'
    vtysh -c "show mpls ldp binding"

    ech 'vtysh -c "show mpls ldp discovery"'
    vtysh -c "show mpls ldp discovery"

    ech 'vtysh -c "show mpls ldp interface"'
    vtysh -c "show mpls ldp interface"

    ech 'vtysh -c "show mpls ldp neighbor"'
    vtysh -c "show mpls ldp neighbor"
}

do_cmds_ip() {
    # Netlink
    ech 'ip link'
    ip link

    ech 'ip addr'
    ip addr

    ech 'ip -4 neigh'
    ip -4 neigh

    ech 'ip -6 neigh'
    ip -6 neigh

    ech 'ip -4 route'
    ip -4 route

    ech 'ip -6 route'
    ip -6 route

    ech 'ip -f mpls route'
    ip -f mpls route

    # Bridge(L2SW)
    ech 'bridge link show'
    bridge link show

    ech 'bridge vlan show'
    bridge vlan show

    ech 'bridge fdb show'
    bridge fdb show
}

do_cmds_gobgp() {
    # GoBGP
    ech 'gobgp policy'
    gobgp policy

    ech 'gobgp neighbor'
    gobgp neighbor

    ech 'gobgp global rib -a ipv4'
    gobgp global rib -a ipv4

    ech 'gobgp global rib -a ipv6'
    gobgp global rib -a ipv6

    ech 'gobgp global rib -a vpnv4'
    gobgp global rib -a vpnv4
}

do_cmds_ribx() {
    # FFlow
    ech "nlac"
    nlac

    ech "ribsc -a 127.0.0.1:50073"
    ribsc

    ech "journalctl -t ribcd"
    journalctl -t ribcd
}

do_cmds() {

    do_cmds_ip

    do_cmds_vtysh

    do_cmds_gobgp

    do_cmds_ribx

    # sysctl
    ech "sysctl"
    sysctl -a
}

dump_netconf() {
    local LOGTOP=$1
    local LOGDIR=$LOGTOP/netconf
    mkdir -p ${LOGDIR}

    local MODULES="beluganos-interfaces beluganos-routing-policy beluganos-network-instance"
    local MODULE
    for MODULE in ${MODULES}; do
        echo $MODULE
        sysrepocfg -d running -x - ${MODULE} > ${LOGDIR}/${MODULE}.running.xml 2>&1
        sysrepocfg -d startup -x - ${MODULE} > ${LOGDIR}/${MODULE}.startup.xml 2>&1
    done
}

dump_flows() {
    local LOGDIR=$1

    ofdump > ${LOGDIR}/ofdump.log 2>&1
    onsldump -a ${ONSLD_ADDR} > ${LOGDIR}/onsldump.log 2>&1
}

copy_fib_log() {
    local LOGDIR=$1

    fibcdmp >> ${LOGDIR}/fibcdmp.log 2>&1
    cp /tmp/fibc.log* ${LOGDIR}/
    cp /tmp/*.pcap    ${LOGDIR}/
}

copy_lxc_log() {
    local LXCNAME=$1
    local LOGTOP=$2
    local LOGDIR="${LOGTOP}/${LXCNAME}"

    mkdir -p ${LOGDIR}

    lxc file pull             ${LXCNAME}/etc/vrf.conf        ${LOGDIR}/
    lxc file pull             ${LXCNAME}/etc/vrf.conf.backup ${LOGDIR}/
    lxc file pull --recursive ${LXCNAME}/etc/frr             ${LOGDIR}/
    lxc file pull --recursive ${LXCNAME}/etc/sysctl.d        ${LOGDIR}/
    lxc file pull --recursive ${LXCNAME}/etc/netplan         ${LOGDIR}/
    lxc file pull --recursive ${LXCNAME}/etc/beluganos       ${LOGDIR}/
}

do_lxc_log() {
    local CMD=$0
    local LOGDIR=$1
    local LXCNAME=$2
    local LOGFILE="${LOGDIR}/${LXCNAME}.log"

    if [ "${LOGDIR}" = "" ]; then
        echo "log dir is empty."
        exit 1
    fi
    if [ "${LXCNAME}" = "" ]; then
        echo "container name is empty."
        exit 1
    fi

    lxc file push ${CMD} ${LXCNAME}/tmp/_remote_cmd.sh
    lxc exec ${LXCNAME} /tmp/_remote_cmd.sh cmds 2>&1 | tee ${LOGFILE}
    copy_lxc_log ${LXCNAME} ${LOGDIR}
}

do_log() {
    local CMD=$0
    local LOGTOP=$1
    local LOGNAME="${PREFIX}-`date +%Y%m%d%H%M%S`"
    local LOGDIR="${LOGTOP}/${LOGNAME}"

    if [ "${LOGTOP}" = "" ]; then
        echo "log dir is empty."
        exit 1
    fi

    mkdir -p     ${LOGDIR}
    copy_fib_log ${LOGDIR}
    dump_flows   ${LOGDIR}
    dump_netconf ${LOGDIR}

    local LXCNAME
    for LXCNAME in ${LXC_LIST}; do
        do_lxc_log ${LOGDIR} ${LXCNAME}
    done

    if [ "${MK_TAR}" = "yes" ]; then
        tar Jcf ${LOGNAME}.tar.xz ${LOGDIR}
    fi
}

do_usage() {
    echo "$0 - log collect tools for beluganos."
    echo ""
    echo "> $0 log <log dir>"
    echo "    Collect logs from LXCs('${LXC_LIST}') and fibc, netconf..."
    echo "    Edit LXC_LIST value to change lxc name list."
    echo ""
    echo "> $0 lxc-log <log dir> <lxc name>"
    echo "    Collect logs from <lxc name>."
    echo ""
    echo "  Settings:"
    echo "   LXC_LIST   : lxc name list."
    echo "   ONSLD_ADDR : OpenNSL agent address."
}

case $1 in
    log)
        do_log $2
        ;;
    lxc-log)
        do_lxc_log $2 $3
        ;;
    cmds)
        do_cmds
        ;;
    *) do_usage;;
esac

