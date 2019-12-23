<<<<<<< HEAD
<<<<<<< HEAD
#! /bin/bash
=======
#! /bin/bash -e
>>>>>>> develop
=======
#! /bin/bash -e
>>>>>>> develop
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

<<<<<<< HEAD
<<<<<<< HEAD
PIP=pip3
=======
PIP=pip
>>>>>>> develop
=======
PIP=pip
>>>>>>> develop
PATCH=patch

INST_HOME=`pwd`/etc/installer
. ${INST_HOME}/sys.sh
. ${INST_HOME}/golang.sh
. ${INST_HOME}/opennsl.sh
<<<<<<< HEAD
<<<<<<< HEAD
. ${INST_HOME}/frr.sh
. ${INST_HOME}/lxd.sh
=======
>>>>>>> develop
=======
>>>>>>> develop

apt_install_build() {
    apt install -y automake unzip gawk pkg-config git libpcap-dev
}

pip_install() {
    $PIP install -U grpcio grpcio-tools protobuf
}

beluganos_install() {
<<<<<<< HEAD
<<<<<<< HEAD
    OPTS="--with-opennsl=$BEL_ONSL_ENABLE" ./bootstrap.sh
    make install
    pushd ./etc/debian
    ./make.sh all
    ./make.sh release
    popd
}

lxd_base_build() {
    local LXD_IMAGE_TEMP="temp"
    local LXD_REL_DIR=./RELEASE

    lxc launch ${LXD_IMAGE_BARE} ${LXD_IMAGE_TEMP}
    sleep 10

    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y update
    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y full-upgrade
    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y autoremove
    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y install libc-ares2 libc6 libcap2 libjson-c3 libpam0g libreadline7 libsystemd0 logrotate iproute2

    lxc file push -r ${LXD_REL_DIR}/rib ${LXD_IMAGE_TEMP}/tmp/
    lxc exec ${LXD_IMAGE_TEMP} /tmp/rib/install.sh -- /tmp/rib

    echo "Stopping container ${LXD_IMAGE_TEMP} ..."
    lxc stop ${LXD_IMAGE_TEMP}

    echo "Publishing container ${LXD_IMAGE_TEMP} as ${LXD_IMAGE_BASE} ..."
    lxc publish ${LXD_IMAGE_TEMP} --alias ${LXD_IMAGE_BASE} || { echo "lxd_base/publish error."; exit 1; }

    echo "Deleting container ${LXD_IMAGE_TEMP} ..."
    lxc delete -f ${LXD_IMAGE_TEMP}

    lxc image info ${LXD_IMAGE_BASE}

    echo "Export ${LXD_IMAGE_BASE} as beluganos-base-lxc ..."
    lxc image export ${LXD_IMAGE_BASE} ${LXD_REL_DIR}/fib/beluganos-base-lxc

    echo "done"
}

frr_pkg_get() {
    frr_pkg_download
    mv ${LXD_WORK_DIR}/${FRR_PKG} ./RELEASE/rib/
}

do_build() {
=======
=======
>>>>>>> develop
    if [ -n "${PROXY}" ]; then
	export PROXY
    fi
    OPTS="--with-opennsl=$BEL_ONSL_ENABLE" ./bootstrap.sh
    make install
    make deb
    make rib
}

do_all() {
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
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
<<<<<<< HEAD
<<<<<<< HEAD

    lxd_init
    lxd_image
    lxd_base_build
}

case $1 in
    build) do_build;;
    frr)   frr_pkg_get;;
    test1) lxd_base_build;;
=======
=======
>>>>>>> develop
}

set_proxy_env
set_sudo

case $1 in
    all) do_all;;
<<<<<<< HEAD
>>>>>>> develop
=======
>>>>>>> develop
esac
