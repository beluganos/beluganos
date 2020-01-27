#! /bin/bash -e
# -*- coding: utf-8 -*-

. ./install.ini

set_proxy() {
    if [ -n "${PROXY}" ]; then
        LXD_PROXY_OPT="--env http_proxy=${PROXY}"
        HTTP_PROXY_OPT="http_proxy=${PROXY} https_proxy=${PROXY}"
        export http_proxy=${PROXY}
        export https_proxy=${PROXY}
        export HTTP_PROXY=${PROXY}
        export HTTPS_PROXY=${PROXY}

        echo "using proxy. ${PROXY}"
    fi
}

download_image() {
    local IMAGE_EXIST=`lxc image list | grep ${LXD_IMAG_NAME}`
    if [ -z "${IMAGE_EXIST}" ]; then
        echo "download lxd image..."
        lxc image copy ${LXD_ORIG_NAME} local: --alias ${LXD_IMAG_NAME}
        lxc image info ${LXD_IMAG_NAME}
    fi
}

create_temp() {
    lxc launch ${LXD_IMAG_NAME} ${LXD_TEMP_NAME}

    # wait for ready conteiner.
    sleep 10

    # diable auto-update.
    lxc exec ${LXD_TEMP_NAME} systemctl -- disable unattended-upgrades || true
    lxc exec ${LXD_TEMP_NAME} systemctl -- stop unattended-upgrades || true

    lxc exec ${LXD_TEMP_NAME} ${LXD_PROXY_OPT} apt -- -y update
    lxc exec ${LXD_TEMP_NAME} ${LXD_PROXY_OPT} apt -- -y full-upgrade
    lxc exec ${LXD_TEMP_NAME} ${LXD_PROXY_OPT} apt -- -y autoremove
}

export_temp() {
    echo "Stopping ${LXD_TEMP_NAME}"
    lxc stop ${LXD_TEMP_NAME} || true

    echo "Publishing ${LXD_TEMP_NAME} as ${LXD_BASE_NAME}"
    lxc publish ${LXD_TEMP_NAME} --alias ${LXD_BASE_NAME}

    echo "Delete ${LXD_TEMP_NAME}"
    lxc delete -f ${LXD_TEMP_NAME} || true

    echo "Export ${LXD_BASE_NAME} as beluganos-lxd-${LXD_BASE_NAME}"
    lxc image export ${LXD_BASE_NAME} ${LXD_FILE_NAME}

    lxc image info ${LXD_BASE_NAME}

    if [ -d "../${FIB_DIR}" ]; then
        echo "copy lxd image to ../${FIB_DIR}/"
        install -m 644 ./${LXD_FILE_NAME}.* ../${FIB_DIR}/
    fi
}

copy_deb() {
    local DEB_FILE

    echo "copy to ${LXD_TEMP_NAME}/tmp/"
    for DEB_FILE in ${RIB_PKGS}; do
        echo "${DEB_FILE}"
        lxc file push ${DEB_FILE} ${LXD_TEMP_NAME}/tmp/ > /dev/null
    done
}

install_deb() {
    local DEB_FILE

    for DEB_FILE in ${RIB_PKGS}; do
        # echo "${DEB_FILE}"
        lxc exec ${LXD_TEMP_NAME} dpkg -- -i /tmp/${DEB_FILE}
    done
}

usage() {
    echo "- create temp container."
    echo "  $0 create-temp"
    echo "- install deb packages to temp container."
    echo "  $0 install-deb"
    echo "- export temp ccontainer as base image."
    echo "  $0 export-temp"
    echo "- download lxc image."
    echo "  $0 dl-image"
}

set_proxy

case $1 in
    create-temp)
        create_temp
        ;;
    export-temp)
        export_temp
        ;;
    install-deb)
        copy_deb
        install_deb
        ;;
    dl-image)
        download_image
        ;;
    *)
        usage
        ;;
esac
