#! /bin/sh
# -*- coding: utf-8 -*-

INSTALL="/usr/bin/install -v"
SYSCTL="/sbin/sysctl"

set_environment() {
    local BEL_CONF_FILE=/etc/beluganos/beluganos.conf

    if [ ! -e ${BEL_CONF_FILE} ]; then
        echo "${BEL_CONF_FILE} not exist."
        exit 1
    fi

    . ${BEL_CONF_FILE}
}

make_dirs() {
    ${INSTALL} -o ${BEL_USER} -g ${BEL_GROUP} -m 755 -d ${BEL_LOG_DIR}
    ${INSTALL} -o ${BEL_USER} -g ${BEL_GROUP} -m 755 -d ${BEL_RUN_DIR}
    ${INSTALL} -o ${BEL_USER} -g ${BEL_GROUP} -m 777 -d ${BEL_LIB_DIR}
}

set_lxd_dns() {
    local LXD_DNS_IP=`lxc network get ${LXD_BRIDGE} ipv4.address | sed -E 's#/.+##'`

    /usr/bin/systemd-resolve --interface ${LXD_BRIDGE} --set-dns ${LXD_DNS_IP} --set-domain ${LXD_DOMAIN}

    echo "LXD DNS enabled. iface:${LXD_BRIDGE} dns:${LXD_DNS_IP} domain:${LXD_DOMAIN}"
}

set_sysctl() {
    ${SYSCTL} -p ${BEL_SYSCTL_CONF}

    echo "sysctl applied. ${BEL_SYSCTL_CONF}"
}

do_all() {
    make_dirs
    set_lxd_dns
    set_sysctl
}

do_usage() {
    echo "$0 mkdirs  --  make dirs for containers."
    echo "$0 lxddns  --  enable DNS for lxd contaner name."
    echo "$0 sysctl  --  apply sysctl settings."
}

set_environment

case $1 in
    mkdirs)
        make_dirs
        ;;
    lxddns)
        set_lxd_dns
        ;;
    sysctl)
        set_sysctl
        ;;
    help)
        do_usage
        ;;
    *)
        do_all
        ;;
esac
