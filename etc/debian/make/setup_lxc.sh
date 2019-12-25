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

CONF_FILE=setup_lxc.ini

LXC_WORK_DIR=/var/lib/beluganos
LXC_BASE=base
BELUGANOS_USR=beluganos

COMMON_DIR="common"

if [ -e "${LXC_WORK_DIR}/${CONF_FILE}" ]; then
    # using setup_lxc.ini on LXC
    echo "include ${LXC_WORK_DIR}/${CONF_FILE}"
    . ${LXC_WORK_DIR}/${CONF_FILE}
fi

lxc_copy_files() {
    local LXC_NAME=$1
    local FILE_NAME
    local USR_NAME
    local COPY_FILE
    local SRC_PATH

    for FILE_NAME in ${!COPY_FILES[@]}; do
        COPY_FILE=${COPY_FILES[$FILE_NAME]}
        USR_NAME=`echo "${COPY_FILE}" | sed -r 's/(.*)@(.+)/\1/'`
        LXC_PATH=`echo "${COPY_FILE}" | sed -r 's/(.+)@(.+)/\2/'`
        SRC_PATH="${LXC_WORK_DIR}/${LXC_NAME}/${FILE_NAME}"
        DST_PATH="${LXC_PREFIX}/${LXC_PATH}"

        echo "[LXC] ${SRC_PATH} -> ${DST_PATH} (${USR_NAME})"
        install -C -m 644 -o ${USR_NAME} -g ${USR_NAME} ${SRC_PATH} ${DST_PATH}
    done
}

lxc_services() {
    local CMD=$1
    local SERVICE_NAME

    for SERVICE_NAME in ${SERVICES}; do
        echo "[LXC] systemctl ${CMD} ${SERVICE_NAME}"
        systemctl ${CMD} ${SERVICE_NAME} || true
    done
}

copy_to_lxc() {
    local CONF_DIR=$1
    local LXC_NAME=$2

    echo "push ${CONF_DIR}/${LXC_NAME} -> ${LXC_NAME}${LXC_WORK_DIR}"

    lxc file push -p -r ${CONF_DIR}/*.sh        ${LXC_NAME}${LXC_WORK_DIR}/
    lxc file push -p -r ${CONF_DIR}/*.ini       ${LXC_NAME}${LXC_WORK_DIR}/
    lxc file push -p -r ${CONF_DIR}/${LXC_NAME} ${LXC_NAME}${LXC_WORK_DIR}/
}

install_to_fib() {
    local FILE_MODE=$1
    local SRC_PATH=$2
    local DST_PATH=$3

    echo "install ${FILE_MODE} ${SRC_PATH} -> ${DST_PATH}"
    install -C -m ${FILE_MODE} -g ${BELUGANOS_USR} -o ${BELUGANOS_USR} ${SRC_PATH} ${DST_PATH}
}

do_install_fib() {
    local CONF_DIR=$1
    local LXC_NAME=$2

    install_to_fib 644 ${CONF_DIR}/${LXC_NAME}/fibc.yml /etc/beluganos/fibc.d/fibc-lxc-${CONF_DIR}.yml
    install_to_fib 644 ${CONF_DIR}/${COMMON_DIR}/snmpproxyd.yaml /etc/beluganos/
    install_to_fib 644 ${CONF_DIR}/${COMMON_DIR}/fibssnmp.yaml   /etc/beluganos/
    install_to_fib 644 ${CONF_DIR}/${COMMON_DIR}/snmp.conf       /etc/snmp/
    install_to_fib 644 ${CONF_DIR}/${COMMON_DIR}/snmpd.conf      /etc/snmp/
}

create_lxc() {
    local CONF_DIR=$1
    local LXC_NAME=$2
    local LXC_EXIST=`lxc list | grep " ${LXC_NAME} "`
    local LXC_RUNNING=`lxc list | grep " ${LXC_NAME} " | grep RUNNING`

    if [ -z "${LXC_EXIST}" ]; then
        sudo mkdir -p /var/log/beluganos/${LXC_NAME}
        lxc profile edit ${LXC_NAME} < ${CONF_DIR}/${LXC_NAME}/lxd_profile.yml
        lxc launch ${LXC_BASE} ${LXC_NAME} -p ${LXC_NAME}
    else
        echo "INFO: ${LXC_NAME} already exist."
        if [ -z "${LXC_RUNNING}" ]; then
            lxc start ${LXC_NAME} || true
        else
            echo "INFO: ${LXC_NAME} already running."
        fi
    fi
}

exec_on_lxc() {
    local LXC_NAME=$1

    lxc exec ${LXC_NAME} ${LXC_WORK_DIR}/$0 -v -- lxc $2 $3 $4 $5
}

usage() {
    echo "beluganos tool."
    echo "  $0 install-fib <conf dir>"
    echo "    install deb packages and initialize for beluganos server."
    echo ""
    echo "  $0 install-lxc <conf dir> <lxc name>"
    echo "    create and initialize container."
    echo ""
    echo "  $0 update-lxc <conf dir> <lxc name>"
    echo "    replace config files and restart services on container."
}

usage_detail() {
    echo "internal commands (on host)"
    echo "  $0 init <conf dir>"
    echo "    install deb packages and import lxc image."
    echo ""
    echo "  $0 setup <conf dir>"
    echo "    copy config files for fib."
    echo ""
    echo "  $0 create-lxc <conf dir> <lxc name>"
    echo "    create lxc profile and launch container."
    echo ""
    echo "  $0 push-to-lxc <conf dir> <lxc name>"
    echo "    copy config files to container (${LXC_WORK_DIR})."
    echo ""
    echo "  $0 exec-on-lxc <lxc name> <lxc subcommand> <params ...>"
    echo "    run 'lxc exec <lxc name> ${LXC_WORK_DIR}/$0 -- lxc <command> <params...>'"
    echo ""
    echo "internal commands (on container)"
    echo "  $0 lxc install <lxc name>"
    echo "    run 'lxc mkdir', 'lxc setup', 'lxc services enable' on container."
    echo ""
    echo "  $0 lxc update <lxc name>"
    echo "    run 'lxc setup', 'lxc services restart' on container."
    echo ""
    echo "  $0 lxc setup <lxc name>"
    echo "    copy config files for rib on container."
    echo ""
    echo "  $0 lxc services <command>"
    echo "    run systemctl <command> on container."
}

case $1 in
    install-fib) # $2:<conf dir> $3:<lxc name>
        do_install_fib $2 $3
        ;;

    install-lxc) # $2:<conf dir> $3:<lxc name>
        create_lxc $2 $3
        copy_to_lxc $2 $3
        exec_on_lxc $3 install $3
        ;;

    update-lxc) # $2:<conf dir> $3:<lxc name>
        copy_files_to_lxc $2 $3
        exec_on_lxc $3 update $3
        ;;

    # internal command
    create-lxc) # $2:<conf dir> $3:<lxc name>
        create_lxc $2 $3
        ;;
    copy-to-lxc) # $2:<conf dir> $3:<lxc name>
        copy_to_lxc $2 $3
        ;;
    exec-on-lxc) # $2:<lxc name> $3:<lxc subcommand> $4,$5:<params>
        exec_on_lxc $2 $3 $4 $5
        ;;
    lxc) # run on LXC.
        case $2 in
            install) # $3:<lxc name>
                lxc_copy_files $3
                lxc_services enable
                ;;
            update) # $3:<lxc name>
                lxc_copy_files $3
                lxc_services restart
                ;;
            setup)
                lxc_copy_files $3
                ;;

            services) # $3:<systemctl command>
                lxc_services $3
                ;;
            *)
                usage
                ;;
        esac
        ;;
    usage)
        usage
        echo ""
        usage_detail
        ;;
    *)
        usage
        ;;
esac
