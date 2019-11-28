#! /bin/bash
# -*- coding: utf-8 -*-

# Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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

. ./create.ini

PIP=pip
PATCH=patch

INST_HOME=`pwd`/etc/installer
. $INST_HOME/sys.sh
. $INST_HOME/frr.sh
. $INST_HOME/lxd.sh
. $INST_HOME/ovs.sh
. $INST_HOME/golang.sh
. $INST_HOME/python.sh
. $INST_HOME/opennsl.sh
. $INST_HOME/beluganos.sh

do_all() {
    confirm "Install ALL" || exit 1

    # install packages and tool
    apt_install
    golang_install
    protoc_install
    opennsl_install

    # create virtual env
    make_virtenv

    # enable virtual-env and go-env
    . ./setenv.sh ${ENABLE_VIRTUALENV}

    # install packages
    pip_install
    gopkg_install
    snmplib_patch
    netlink_patch
    gobgp_upgrade
    ryu_patch

    # create frr deb package
    frr_pkg

    # initailize systems
    init_lxd
    init_sys
    init_ovs ${BELUG_OVS_BRIDGE} 127.0.0.1

    # beluganos-netconf
    netconf_install

    # beluganos-beluganos
    beluganos_install
}

do_minimal() {
    confirm "Install minimal" || exit 1

    sudo ${HTTP_PROXY} ${APT_OPTION} apt -y install ${APT_MINS}
    make_virtenv
    . ./setenv.sh ${ENABLE_VIRTUALENV}
    get_pip
    $PIP install -U ${PIP_PROXY} ansible
    init_lxd
    init_sys
    init_ovs ${SAMPLE_OVS_BRIDGE} ${BELUG_OFC_ADDR} ${SAMPLE_OVS_DPID}
}

do_opennsl() {
    opennsl_install
    . ./setenv.sh ${ENABLE_VIRTUALENV}
    beluganos_install
}

do_usage() {
    echo "Usage $0 [OPTIONS]"
    echo "Options:"
    echo "  ''    : run all"
    echo "  pkg   : update apt-packages and re-install golang and protoc"
    echo "  pip   : update pip-packages"
    echo "  gopkg : update go-packages"
    echo "  min   : minimal install for frr container."
    echo "  help  : show this message"
}

set_proxy_env
set_proxy_lxd
set_sudo
case $1 in
    pkg)
        apt_install
        golang_install
        protoc_install
        ;;
    pip)
        . ./setenv.sh ${ENABLE_VIRTUALENV}
        pip_install
        ryu_patch
        ;;
    gopkg)
        gopkg_install
        snmplib_patch
        netlink_patch
        gobgp_upgrade
        ;;
    min)
        do_minimal
        ;;
    netconf)
        netconf_install
        ;;
    opennsl)
        do_opennsl
        ;;
    help)
        do_usage
        ;;
    *)
        do_all
        ;;
esac
