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

#
# ask confirm
#
confirm() {
    MSG=$1

    echo "$MSG [y/N]"
    read ans
    case $ans in
        [yY]) return 0;;
        *) return 1;;
    esac
}

set_proxy_env() {
    if [ "${PROXY}"x != ""x ]; then
        APT_PROXY="--env http_proxy=${PROXY}"
        HTTP_PROXY="http_proxy=${PROXY} https_proxy=${PROXY}"
        export http_proxy=${PROXY}
        export https_proxy=${PROXY}
        export HTTP_PROXY=${PROXY}
        export HTTPS_PROXY=${PROXY}

        echo "--- Proxy settings ---"
        echo "APT_PROXY=${APT_PROXY}"
        echo "HTTP_PROXY=${HTTP_PROXY}"
    fi
}

set_proxy_lxd() {
    if [ "${PROXY}"x != ""x ]; then
        lxc config set core.proxy_http ${PROXY}
        lxc config set core.proxy_https ${PROXY}

        echo "--- Proxy settings ---"
	echo "lxc proxy ${PROXY}"
    fi
}

set_sudo() {
    if [ "${ENABLE_VIRTUALENV}" != "yes" ]; then
        PIP="sudo -E $PIP"
        PATCH="sudo $PATCH"
    fi
}


#
# install deb packages
#
apt_install() {
    sudo ${HTTP_PROXY} ${APT_OPTION} apt -y install ${APT_PKGS} || { echo "apt_install error."; exit 1; }
    sudo apt -y autoremove
}

#
# modules
#
init_module() {
    sudo cp -v etc/modules/modules.conf  /etc/modules-load.d/beluganos.conf
    sudo cp -v etc/modules/modprobe.conf /etc/modprobe.d/beluganos.conf
    sudo modprobe -a belbonding mpls_router mpls_iptunnel ip_tunnel ip6_tunnel
    sudo netplan apply
}

#
# system
#
init_sys() {
    sudo useradd -s /sbin/nologin -r ${BELUG_USER}
    sudo mkdir -v -p ${BELUG_HOME}
    sudo mkdir -v -p ${BELUG_DIR}

    local IFACE_TEMP=/tmp/interfaces_temp
    cat >  ${IFACE_TEMP} <<EOF
# -*- coding: utf-8 -*-
network:
  version: 2
  renderer: networkd
  ethernets:
    ${BELUG_OFC_IFACE}:
      addresses:
        - ${BELUG_OFC_ADDR}
EOF
    sudo cp ${IFACE_TEMP} /etc/netplan/02-beluganos.yaml

    init_module
}
