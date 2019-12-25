<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
#! /bin/sh
=======
#! /bin/bash -e
>>>>>>> develop
=======
#! /bin/bash -e
>>>>>>> develop
=======
#! /bin/bash -e
>>>>>>> develop
=======
#! /bin/bash -e
>>>>>>> develop
# -*- coding: utf-8 -*-

. ./install.ini

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
BEL_PKGS="} beluganos-fib-*_amd64.deb"
ALL_PKGS="${DEB_PKGS} ${BEL_PKGS}"

=======
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
=======
>>>>>>> develop
init_lxd() {
    lxd init --preseed < ./lxd-init.yaml
    lxc network set ${LXD_BRIDGE} ipv4.address ${LXD_NETWORK}
    lxc network show ${LXD_BRIDGE}
}

import_image() {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
    lxc image import ./beluganos-lxd-*.tar.gz --alias ${LXD_BASE_NAME}
=======
    lxc image import ./${LXD_FILE_NAME}.tar.gz --alias ${LXD_BASE_NAME}
>>>>>>> develop
=======
    lxc image import ./${LXD_FILE_NAME}.tar.gz --alias ${LXD_BASE_NAME}
>>>>>>> develop
=======
    lxc image import ./${LXD_FILE_NAME}.tar.gz --alias ${LXD_BASE_NAME}
>>>>>>> develop
=======
    lxc image import ./${LXD_FILE_NAME}.tar.gz --alias ${LXD_BASE_NAME}
>>>>>>> develop
}

install_deb() {
    local DEB_FILE

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
    for DEB_FILE in ${ALL_PKGS}; do
=======
    for DEB_FILE in ${FIB_PKGS}; do
>>>>>>> develop
=======
    for DEB_FILE in ${FIB_PKGS}; do
>>>>>>> develop
=======
    for DEB_FILE in ${FIB_PKGS}; do
>>>>>>> develop
=======
    for DEB_FILE in ${FIB_PKGS}; do
>>>>>>> develop
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
