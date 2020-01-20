#! /bin/bash -e
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

. ./create.ini

PATCH=patch

INST_HOME=`pwd`/etc/installer
. ${INST_HOME}/sys.sh
. ${INST_HOME}/golang.sh
. ${INST_HOME}/opennsl.sh

apt_install() {
    sudo -E apt update
    sudo -E apt install -y automake unzip gawk pkg-config git libpcap-dev
}

beluganos_install() {
    if [ -n "${PROXY}" ]; then
        export PROXY
    fi
    OPTS="--with-opennsl=$BEL_ONSL_ENABLE --without-python-grpc" ./bootstrap.sh
    make install
    make deb
    make rib
}

do_all() {
    confirm "Install ALL" || exit 1

    apt_install
    golang_install
    protoc_install
    opennsl_install

    . ./setenv.sh

    gopkg_install
    snmplib_patch
    netlink_patch
    gobgp_upgrade

    beluganos_install
}

do_usage() {
    echo "$0 all       - install and setup."
    echo "$0 aptpkg    - install apt packages"
    echo "$0 gopkg     - install go packages."
    echo "$0 golang    - install golang and protoc"
    echo "$0 beluganos - build and create beluganos packages.."
}

set_proxy_env
set_sudo

case $1 in
    all)
        do_all
        ;;
    aptpkg)
        apt_install
        ;;
    golang)
        golang_install
        protoc_install
        ;;
    gopkg)
        gopkg_install
        ;;
    beluganos)
        beluganos_install
        ;;
    *)
        do_usage
        ;;
esac
