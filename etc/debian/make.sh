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

# Proxy
# PROXY=http://192.168.1.100:8080

# Frr version
# FRRVER="frr-stable"
# FRRVER="frr-7"
FRRVER="frr-6"

# DO NOT EDIT
FIB_MODULES="fib-common fib-fibc fib-fibs govsw"
RIB_MODULES="rib-common rib-ribc rib-ribs rib-ribp rib-ribt rib-snmp rib-ribn gonla gobgp"
WBS_MODULES="gonsl"
MODULES="${FIB_MODULES} ${RIB_MODULES} ${WBS_MODULES}"
MAKE=`pwd`/make/make.sh
RELDIR=`pwd`/../../RELEASE
DEBDIR=`pwd`/deb-cache

set_proxy() {
    if [ -n "${PROXY}" ]; then
        HTTP_PROXY_OPT="http_proxy=${PROXY} https_proxy=${PROXY}"
        export http_proxy=${PROXY}
        export https_proxy=${PROXY}
        export HTTP_PROXY=${PROXY}
        export HTTPS_PROXY=${PROXY}

        echo "using proxy. ${PROXY}"
    fi
}

do_make() {
    local OPTS=$1
    local MODULE

    for MODULE in ${MODULES}; do
        pushd ./${MODULE}
        ${MAKE} ${OPTS}
        popd
    done
}

release_bel() {
    local INSTDIR=${RELDIR}/$1
    local DIRLIST=$2

    install -d ${INSTDIR}

    for DIRNAME in ${DIRLIST}; do
        echo "${DIRNAME} -> ${INSTDIR}"
        install -C -m 644 ./${DIRNAME}/*.deb ${INSTDIR}/
    done

    pushd ${INSTDIR}
    md5sum *.deb > md5sum.txt
    popd
}

release_fib() {
    local INSTDIR=${RELDIR}/$1

    install -d ${INSTDIR}
    install -C -m 644 ./make/install.ini    ${INSTDIR}/install.ini
    install -C -m 755 ./make/install_fib.sh ${INSTDIR}/install.sh
    install -C -m 644 ./make/lxd-init.yaml  ${INSTDIR}/lxd-init.yaml
    install -C -m 644 ${DEBDIR}/* ${INSTDIR}/
    rm -f ${INSTDIR}/frr*.deb ${INSTDIR}/libc-ares*.deb ${INSTDIR}/libyang0.16*.deb
}

release_rib() {
    local INSTDIR=${RELDIR}/$1

    install -d ${INSTDIR}
    install -C -m 644 ./make/install.ini    ${INSTDIR}/install.ini
    install -C -m 755 ./make/install_rib.sh ${INSTDIR}/install.sh
    install -C -m 644 ./make/Makefile_rib   ${INSTDIR}/Makefile
    install -C -m 644 ${DEBDIR}/* ${INSTDIR}/
}

download_deb() {
    install -d ${DEBDIR}
    pushd ${DEBDIR}

    local PKG_LIST
    local PKG_NAME
    PKG_LIST="snmpd snmp snmp-mibs-downloader libsnmp-base"
    PKG_LIST="${PKG_LIST} libc6 libsnmp30 libssl1.1 libsensors4 libc-ares2 libyang0.16"
    PKG_LIST="${PKG_LIST} adduser debconf lsb-base smistrip libjansson4 libnftnl7 nftables"

    for PKG_NAME in ${PKG_LIST}; do
        apt-get download ${PKG_NAME} || echo "download error. ${PKG_NAME}"
    done

    popd
}

download_frr() {
    install -d ${DEBDIR}
    pushd ${DEBDIR}

    curl -s https://deb.frrouting.org/frr/keys.asc | sudo apt-key add -
    sudo rm -f /etc/apt/sources.list.d/frr.list
    echo deb https://deb.frrouting.org/frr $(lsb_release -s -c) $FRRVER | sudo tee -a /etc/apt/sources.list.d/frr.list
    sudo ${HTTP_PROXY_OPT} apt update
    apt-get download frr frr-pythontools

    popd
}

do_release() {
    release_fib fib
    release_bel fib  "${FIB_MODULES}"

    release_rib rib
    release_bel rib  "${RIB_MODULES}"

    release_bel wbsw "${WBS_MODULES}"
}

do_usage() {
    echo "$0 all"
    echo "$0 release"
    echo "$0 clean"
}

set_proxy

case $1 in
    all)
        do_make all
        ;;
    release)
        do_release
        ;;
    clean)
        rm -fr ${DEBDIR}
        do_make clean
        ;;
    dl-deb)
        download_deb
        ;;
    dl-frr)
        download_frr
        ;;
    deb)
        pushd ./$2
        ${MAKE} $3
        popd
        ;;
    *)
        do_usage
        ;;
esac
