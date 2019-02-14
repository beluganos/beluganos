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
# frr deb package
#
frr_pkg() {
    if [ "$FRR_BRANCH" = "3.0" ]; then
        frr_pkg_build
    else
        frr_pkg_download
    fi
}

frr_pkg_build() {
    local FRR_DIR=${LXD_WORK_DIR}/frr

    if [ -e $FRR_DIR ]; then
        pushd $FRR_DIR
    else
        git clone $FRR_URL $FRR_DIR || { echo "frr_pkg/clone error."; exit 1; }
        cp etc/frr/frr.patch /tmp/
        cp etc/frr/frr-stable-3.0-for-ubuntu-1804.patch /tmp/

        pushd $FRR_DIR
        git checkout -b $FRR_BRANCH origin/stable/$FRR_BRANCH
        patch -p1 < /tmp/frr.patch
        patch -p1 < /tmp/frr-stable-3.0-for-ubuntu-1804.patch
        ln -s debianpkg debian
    fi

    ./bootstrap.sh
    ./configure
    make dist
    make -f debian/rules backports

    cd ${LXD_WORK_DIR}
    tar xvf ${FRR_DIR}/frr_*.orig.tar.gz
    cd frr-*
    . /etc/os-release
    tar xvf ${FRR_DIR}/frr_*${ID}${VERSION_ID}*.debian.tar.xz

    fakeroot ./debian/rules binary

    popd
}

frr_pkg_download() {
    pushd ${LXD_WORK_DIR}
    wget -O ${FRR_PKG} ${FRR_DOWNLOAD}
    popd
}
