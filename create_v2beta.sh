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

PIP=pip
PATCH=patch

INST_HOME=`pwd`/etc/installer
. ${INST_HOME}/sys.sh
. ${INST_HOME}/golang.sh
. ${INST_HOME}/opennsl.sh

apt_install_build() {
    apt install -y automake unzip gawk pkg-config git libpcap-dev
}

pip_install() {
    $PIP install -U grpcio grpcio-tools protobuf
}

beluganos_install() {
    if [ -n "${PROXY}" ]; then
	export PROXY
    fi
    OPTS="--with-opennsl=$BEL_ONSL_ENABLE" ./bootstrap.sh
    make install
    make deb
    make rib
}

do_all() {
    confirm "Install ALL" || exit 1

    apt_install_build
    golang_install
    protoc_install
    opennsl_install

    . ./setenv.sh ${ENABLE_VIRTUALENV}

    pip_install
    gopkg_install
    snmplib_patch
    netlink_patch
    gobgp_upgrade

    beluganos_install
}

set_proxy_env
set_sudo

case $1 in
    all) do_all;;
esac
