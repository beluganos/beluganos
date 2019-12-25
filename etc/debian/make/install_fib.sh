#! /bin/bash -e
# -*- coding: utf-8 -*-

. ./install.ini

init_lxd() {
    lxd init --preseed < ./lxd-init.yaml
    lxc network set ${LXD_BRIDGE} ipv4.address ${LXD_NETWORK}
    lxc network show ${LXD_BRIDGE}
}

import_image() {
    lxc image import ./${LXD_FILE_NAME}.tar.gz --alias ${LXD_BASE_NAME}
}

install_deb() {
    local DEB_FILE

    for DEB_FILE in ${FIB_PKGS}; do
        # echo "${DEB_FILE}"
        dpkg -i ${DEB_FILE}
    done
}

do_all() {
    init_lxd
    import_image
    install_deb
}

usage() {
    echo "- run all"
    echo "  $0 all"
    echo "- initialize lxd and network,"
    echo "  $0 init-lxd"
    echo "- import lxd image."
    echo "  $0 import"
    echo "- install deb packages"
    echo "  $0 install-deb"
}

case $1 in
    all)
        do_all
        ;;
    init-lxd)
        init_lxd
        ;;

    import)
        import_image
        ;;

    install-deb)
        install_deb
        ;;

    *)
        usage
        ;;
esac
