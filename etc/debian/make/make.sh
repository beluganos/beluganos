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

. ./make.conf

if [ "${VERSION}" = "" ]; then
    VERSION=1.0.0-2
fi

PREFIX=`pwd`
WORKDIR=${PREFIX}/work
BINDIR=${PREFIX}/../../../bin
GO_BINDIR=${HOME}/go/bin
USR_BIN=usr/bin
DEBIANS="control postinst preinst postrm prerm"

check_exist() {
    local FILENAME=$1

    if [ ! -e ${FILENAME} ]; then
        echo "[*NG*] $1 not found!!"
    else
        echo "[ OK ] $1"
    fi
}

build_deb() {
    local PKGDIR=$1

    fakeroot dpkg-deb --build ${PKGDIR} .
}

make_debian_dir() {
    local DIRNAME=$1

    install -pd ${DIRNAME}/DEBIAN
}

copy_debians() {
    local SRCDIR=$1
    local DSTDIR=$2
    local FILENAME

    for FILENAME in ${DEBIANS}; do
        install -pm 755 ${SRCDIR}/debian/${FILENAME} ${DSTDIR}/DEBIAN/
    done
}

check_debians() {
    local DIRNAME=$1
    local FILENAME

    for FILENAME in ${DEBIANS}; do
        check_exist ${DIRNAME}/DEBIAN/${FILENAME}
    done
}

set_version() {
    local FILENAME=$1/DEBIAN/control
    local VER=$2

    sed -e "s/^Version: .*/Version: ${VER}/g" ${FILENAME} > control.temp
    mv control.temp ${FILENAME}
}

make_dirs() {
    local DIRNAME
    for DIRNAME in $DIRS; do
        install -pd ${WORKDIR}/${DIRNAME}
    done

    make_debian_dir ${WORKDIR}
}

copy_files() {
    local FILENAME

    if [ -n "${GO_BINS}" ]; then
        for FILENAME in ${GO_BINS}; do
            install -pm 755 ${GO_BINDIR}/${FILENAME} ${WORKDIR}/${USR_BIN}/
        done
    fi

    if [ -n "${BINS}" ]; then
        for FILENAME in ${BINS}; do
            install -pm 755 ${BINDIR}/${FILENAME} ${WORKDIR}/${USR_BIN}/
        done
    fi

    for FILENAME in "${!COPY_FILES[@]}"; do
        DIRNAME=${COPY_FILES[$FILENAME]}
        install -pm 644 ${PREFIX}/files/${FILENAME} ${WORKDIR}/${DIRNAME}
    done

    copy_debians ${PREFIX}  ${WORKDIR}
    set_version  ${WORKDIR} ${VERSION}
}

check_files() {
    local FILENAME

    if [ -n "${GO_BINS}" ]; then
        for FILENAME in ${GO_BINS}; do
            check_exist ${WORKDIR}/${USR_BIN}/${FILENAME}
        done
    fi

    if [ -n "${BINS}" ]; then
        for FILENAME in ${BINS}; do
            check_exist ${WORKDIR}/${USR_BIN}/${FILENAME}
        done
    fi

    for FILENAME in "${!COPY_FILES[@]}"; do
        DIRNAME=${COPY_FILES[$FILENAME]}
        check_exist ${WORKDIR}/${DIRNAME}/${FILENAME}
    done

    check_debians ${WORKDIR}
}

do_clean() {
    rm -vfr ${WORKDIR}
    rm -vf *.deb
}

do_build() {
    build_deb ${WORKDIR}
}

do_all() {
    make_dirs
    copy_files
    check_files
    do_build
}

usage() {
    echo "$0 <option>"
    echo "option:"
    echo "  dirs   - make work dirs."
    echo "  files  - copy files to work dir."
    echo "  check  - check files."
    echo "  deb    - build deb package."
    echo "  all    - do all."
    echo "  clean  - clean all work files."
}

case $1 in
    dirs)   make_dirs;;
    files)  copy_files;;
    check)  check_files;;
    deb)    do_build;;
    all)    do_all;;
    clean)  do_clean;;
    *)      usage;;
esac
