#! /bin/bash
# -*- coding: utf-8 -*-

PREFIX=`pwd`
SRCDIR=`pwd`
BINDIR=`pwd`/../../bin
WORKDIR=${PREFIX}/work

file_exist() {
    local FILENAME=${SRCDIR}/$1
    if [ ! -e $FILENAME ]; then
        echo "[*NG*] $1 not found!!"
    else
        echo "[ OK ] $1"
    fi
}

make_dirs() {
    install -pd ${WORKDIR}/usr/bin
    install -pd ${WORKDIR}/etc
    # install -pd ${WORKDIR}/etc/network/interfaces.d
    # install -pd ${WORKDIR}/etc/systemd/system
    install -pd ${WORKDIR}/DEBIAN
}

copy_files() {
    local MODE=$1

    if [ -z "${MODE}" ]; then
        echo "mode not found."
        exit 1
    fi

    install -pm 755 ${BINDIR}/ffctl                  ${WORKDIR}/usr/bin/ffctl-ztp
    install -pm 755 ${SRCDIR}/files/rc.local         ${WORKDIR}/etc/rc.local
    # install -pm 644 ${SRCDIR}/files/ma1              ${WORKDIR}/etc/network/interfaces.d/ma1
    # install -pm 644 ${SRCDIR}/files/ztp-init.service ${WORKDIR}/etc/systemd/system/
    install -pm 755 ${SRCDIR}/files/ztp-init         ${WORKDIR}/usr/bin/
    install -pm 644 ${SRCDIR}/files/ztp-init-${MODE}.conf ${WORKDIR}/etc/ztp-init.conf

    install -pm 755 ${SRCDIR}/debian/control-${MODE}   ${WORKDIR}/DEBIAN/control
    install -pm 755 ${SRCDIR}/debian/preinst   ${WORKDIR}/DEBIAN
    install -pm 755 ${SRCDIR}/debian/postinst  ${WORKDIR}/DEBIAN
    install -pm 755 ${SRCDIR}/debian/prerm     ${WORKDIR}/DEBIAN
    install -pm 755 ${SRCDIR}/debian/postrm    ${WORKDIR}/DEBIAN
}

check_files() {
    file_exist ${BINDIR}/ffctl
}

clean_all() {
    rm -fr ${WORKDIR}
}

make_deb() {
    fakeroot dpkg-deb --build ${WORKDIR} .
}

make_pkg() {
    local MODE=$1
    make_dirs
    copy_files ${MODE}
    make_deb
}

do_all() {
    clean_all
    make_pkg kvm
    clean_all
    make_pkg ztp
}

do_usage() {
    echo "$0 <command>"
    echo "command:"
    echo "  make-dirs         :  make work dirs."
    echo "  copy-files <mode> : copy files to work."
    echo "  check             : check files."
    echo "  deb <mode>        : create deb."
    echo "  clean             : clean work files."
    echo "  <mode>            : 'kvm' or 'ztp'"
    exit 1
}


case $1 in
    make-dirs)  make_dirs;;
    copy-files) copy_files $2;;
    check)      check_files;;
    clean)      clean_all;;
    deb)        make_pkg $2;;
    all)        do_all;;
    *)          do_usage;;
esac
