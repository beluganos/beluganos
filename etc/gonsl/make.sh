#! /bin/bash
# -*- coding: utf-8 -*-

PREFIX=`pwd`

PKGDIR=${PREFIX}/work
SRCDIR=`pwd`

make_dirs() {
    install -pd ${PKGDIR}/usr/bin
    install -pd ${PKGDIR}/usr/lib
    install -pd ${PKGDIR}/etc/opennsl/drivers
    install -pd ${PKGDIR}/etc/init.d
    install -pd ${PKGDIR}/DEBIAN
}

copy_files() {
    install -pm  755 ${SRCDIR}/gonsld              ${PKGDIR}/usr/bin/
    install -pm  644 ${SRCDIR}/libopennsl.so.1     ${PKGDIR}/usr/lib/
    install -pm  755 ${SRCDIR}/opennsl.conf        ${PKGDIR}/etc/opennsl/
    install -pm  644 ${SRCDIR}/linux-kernel-bde.ko ${PKGDIR}/etc/opennsl/drivers/
    install -pm  644 ${SRCDIR}/linux-user-bde.ko   ${PKGDIR}/etc/opennsl/drivers/
    install -pm  644 ${SRCDIR}/linux-bcm-knet.ko   ${PKGDIR}/etc/opennsl/drivers/

    install -pm  644 ${SRCDIR}/files/gonsld.yaml   ${PKGDIR}/etc/opennsl/
    install -pm  644 ${SRCDIR}/files/gonsld.conf   ${PKGDIR}/etc/opennsl/
    install -Tpm 755 ${SRCDIR}/files/gonsld.initd  ${PKGDIR}/etc/init.d/gonsld

    install -pm  755 ${SRCDIR}/debian/control      ${PKGDIR}/DEBIAN/
    install -pm  755 ${SRCDIR}/debian/postinst     ${PKGDIR}/DEBIAN/
    install -pm  755 ${SRCDIR}/debian/prerm        ${PKGDIR}/DEBIAN/
}

make_deb() {
    fakeroot dpkg-deb --build ${PKGDIR} .
}

do_debpkg() {
    make_dirs
    copy_files
    make_deb
}

file_exist() {
    local FILENAME=${SRCDIR}/$1
    if [ ! -e $FILENAME ]; then
        echo "[*NG*] $1 not found!!"
    else
        echo "[ OK ] $1"
    fi
}


do_check() {
    file_exist gonsld
    file_exist libopennsl.so.1
    file_exist opennsl.conf
    file_exist linux-kernel-bde.ko
    file_exist linux-user-bde.ko
    file_exist linux-bcm-knet.ko

    file_exist files/gonsld.yaml
    file_exist files/gonsld.conf
    file_exist files/gonsld.initd

    file_exist debian/control
    file_exist debian/postinst
    file_exist debian/prerm
}


do_clean() {
    rm -fr ${PKGDIR}
}

do_distclean() {
    rm -f *.ko *.deb
    rm -f gonsld opennsl.conf libopennsl.so.1 libopennsl.so
}

do_usage() {
    echo "$0 <check | deb | clean>"
}


case $1 in
    make-dirs)
        make_dirs
        ;;

    copy-files)
        copy_files
        ;;

    deb)
        do_check
        do_debpkg
        ;;

    check)
        do_check
        ;;

    clean)
        do_clean
        ;;
    distclean)
        do_clean
        do_distclean
        ;;

    *)
        do_usage
        ;;
esac
