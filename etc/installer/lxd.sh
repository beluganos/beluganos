#! /bin/bash
# -*- coding: utf-8 -*-

# Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

lxd_init() {
    if [ "$LXD_INIT" = "yes" ]; then
        echo "lxd init..."
        lxd init --preseed < ${INST_HOME}/lxd-init.yaml
    fi
}

#
# lxdbr0 setting
#
lxd_network() {
    lxc network set ${LXD_BRIDGE} ipv4.address ${LXD_NETWORK}
    lxc network show ${LXD_BRIDGE}
}

#
# ubuntu image
#
lxd_image() {
    lxc image copy ${LXD_IMAGE_ORIG} local: --alias ${LXD_IMAGE_BARE}
    lxc image info ${LXD_IMAGE_BARE}
}

#
# base image
#
lxd_base() {
    local LXD_IMAGE_TEMP="temp"

    if [ ! -e ${LXD_WORK_DIR}/${FRR_PKG} ]; then
        echo "${LXD_WORK_DIR}/${FRR_PKG} not exist!!"
        exit -1
    fi

    lxc launch ${LXD_IMAGE_BARE} ${LXD_IMAGE_TEMP}
    sleep 10

    echo "Installing packages"
    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y update || { echo "lxd_base/update error."; exit 1; }
    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y dist-upgrade || { echo "lxd_base upgrade error"; exit 1; }
    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y install ${LXD_APT_PKGS} || { echo "lxd_base/install error."; exit 1; }
    lxc exec ${LXD_IMAGE_TEMP} apt ${APT_PROXY} -- -y autoremove

    echo "Push ${FRR_PKG} to ${LXD_IMAGE_TEMP}"
    lxc file push ${LXD_WORK_DIR}/${FRR_PKG} ${LXD_IMAGE_TEMP}/tmp/

    echo "Installing ${FRR_PKG} ..."
    lxc exec ${LXD_IMAGE_TEMP} dpkg -- -i /tmp/${FRR_PKG} || { echo "lxd_base/dpkg error."; exit 1; }

    echo "Stopping container ${LXD_IMAGE_TEMP} ..."
    lxc stop ${LXD_IMAGE_TEMP}

    echo "Publishing container ${LXD_IMAGE_TEMP} as ${LXD_IMAGE_BASE} ..."
    lxc publish ${LXD_IMAGE_TEMP} --alias ${LXD_IMAGE_BASE} || { echo "lxd_base/publish error."; exit 1; }

    echo "Deleting container ${LXD_IMAGE_TEMP} ..."
    lxc delete -f ${LXD_IMAGE_TEMP}

    lxc image info ${LXD_IMAGE_BASE}

    echo "done"
}

#
# setup lxd
#
init_lxd() {
    lxd_init
    lxd_network
    lxd_image
    lxd_base
}

