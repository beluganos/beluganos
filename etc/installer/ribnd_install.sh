#! /bin/bash
# -*- coding: utf-8 -*-

SRCDIR=`pwd`
BINDIR=../../bin

do_install() {
    local LXCNAME=$1

    if [ "${LXCNAME}" = "" ]; then
        echo "container not specfied."
        exit 1
    fi

    echo "copy files..."
    lxc file push ${BINDIR}/ribnd         ${LXCNAME}/usr/bin/
    lxc file push ${BINDIR}/ribnc         ${LXCNAME}/usr/bin/
    lxc file push ${SRCDIR}/ribnd.conf    ${LXCNAME}/etc/beluganos/
    lxc file push ${SRCDIR}/ribnd.yaml    ${LXCNAME}/etc/beluganos/
    lxc file push ${SRCDIR}/ribnd.service ${LXCNAME}/etc/systemd/system/

    echo "register as systemd service."
    lxc exec ${LXCNAME} systemctl -- daemon-reload

    echo "install success."
}

do_change_enable() {
    local LXCNAME=$1
    local MODE=$2

    
    if [ "${MODE}" = "" ]; then
        echo "invalid argument"
        exit 1
    fi

    lxc exec ${LXCNAME} systemctl -- ${MODE} ribnd

    echo "${MODE} ribnd success."
}

usage() {
    echo "$0 install <container>"
    echo "$0 enable  <container>"
    echo "$0 disable <container>"
}

case $1 in
    install)
        do_install $2
        ;;
    enable)
        do_change_enable $2 enable
        ;;

    disable)
        do_change_enable $2 disable
        ;;
    *)
        usage
        ;;
esac
