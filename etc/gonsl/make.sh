#! /bin/bash
# -*- coding: utf-8 -*-

PREFIX=`pwd`
SRCDIR=`pwd`
BINDIR=`pwd`/../../bin
OPENNSL_HWTYPE=as7712
OPENNSL_BINDIR=`pwd`/../../../opennsl/bin/${OPENNSL_HWTYPE}

WORKDIR=${PREFIX}/work
GONSLD_DIR=${WORKDIR}/gonsl
OPENNSL_DIR=${WORKDIR}/opennsl

file_exist() {
    local FILENAME=${SRCDIR}/$1
    if [ ! -e $FILENAME ]; then
        echo "[*NG*] $1 not found!!"
    else
        echo "[ OK ] $1"
    fi
}

gonsld_make_dirs() {
    install -pd ${GONSLD_DIR}/usr/bin
    install -pd ${GONSLD_DIR}/etc/beluganos
    install -pd ${GONSLD_DIR}/etc/init.d
    install -pd	${GONSLD_DIR}/DEBIAN
}

gonsld_prepare_files() {
    install -pm 755 ${BINDIR}/gonsld  ${SRCDIR}/
}

gonsld_copy_files() {
    install -pm  755 ${SRCDIR}/gonsld             ${GONSLD_DIR}/usr/bin/
    install -pm  644 ${SRCDIR}/files/gonsld.yaml  ${GONSLD_DIR}/etc/beluganos/
    install -pm  644 ${SRCDIR}/files/gonsld.conf  ${GONSLD_DIR}/etc/beluganos/
    install -Tpm 755 ${SRCDIR}/files/gonsld.initd ${GONSLD_DIR}/etc/init.d/gonsld

    install -pm  755 ${SRCDIR}/debian/gonsld_control  ${GONSLD_DIR}/DEBIAN/control
    install -pm  755 ${SRCDIR}/debian/gonsld_postinst ${GONSLD_DIR}/DEBIAN/postinst
    install -pm  755 ${SRCDIR}/debian/gonsld_prerm    ${GONSLD_DIR}/DEBIAN/prerm
}

gonsld_check() {
    file_exist gonsld
    file_exist files/gonsld.yaml
    file_exist files/gonsld.conf
    file_exist files/gonsld.initd

    file_exist debian/gonsld_control
    file_exist debian/gonsld_postinst
    file_exist debian/gonsld_prerm
}

gonsld_clean() {
    rm -vfr ${GONSLD_DIR}
}

gonsld_make_deb() {
    fakeroot dpkg-deb --build ${GONSLD_DIR} .
}

gonsld_make_pkg() {
    gonsld_make_dirs
    gonsld_copy_files
    gonsld_make_deb
}

opennsl_make_dirs() {
    install -pd ${OPENNSL_DIR}/usr/bin/
    install -pd ${OPENNSL_DIR}/usr/lib/
    install -pd ${OPENNSL_DIR}/etc/opennsl/drivers/
    install -pd ${OPENNSL_DIR}/DEBIAN
}

opennsl_prepare_files() {
    install -pm 644 ${OPENNSL_BINDIR}/libopennsl.so.1     ${SRCDIR}/
    install -pm 644 ${OPENNSL_BINDIR}/linux-kernel-bde.ko ${SRCDIR}/
    install -pm 644 ${OPENNSL_BINDIR}/linux-user-bde.ko   ${SRCDIR}/
    install -pm 644 ${OPENNSL_BINDIR}/linux-bcm-knet.ko   ${SRCDIR}/
    if [ -e ${SRCDIR}/files/config.${OPENNSL_HWTYPE} ]; then
	install -pm 644  ${SRCDIR}/files/config.${OPENNSL_HWTYPE}  ${SRCDIR}/opennsl.conf
    else
	install -pm 644 ${OPENNSL_BINDIR}/config.${OPENNSL_HWTYPE} ${SRCDIR}/opennsl.conf
    fi
}

opennsl_copy_files() {
    install -pm  755 ${SRCDIR}/files/opennsl_setup  ${OPENNSL_DIR}/usr/bin/
    install -pm  644 ${SRCDIR}/libopennsl.so.1      ${OPENNSL_DIR}/usr/lib/
    install -pm  755 ${SRCDIR}/opennsl.conf         ${OPENNSL_DIR}/etc/opennsl/
    install -pm  644 ${SRCDIR}/linux-kernel-bde.ko  ${OPENNSL_DIR}/etc/opennsl/drivers/
    install -pm  644 ${SRCDIR}/linux-user-bde.ko    ${OPENNSL_DIR}/etc/opennsl/drivers/
    install -pm  644 ${SRCDIR}/linux-bcm-knet.ko    ${OPENNSL_DIR}/etc/opennsl/drivers/

    install -pm  755 ${SRCDIR}/debian/opennsl_control  ${OPENNSL_DIR}/DEBIAN/control
    install -pm  755 ${SRCDIR}/debian/opennsl_postinst ${OPENNSL_DIR}/DEBIAN/postinst
    install -pm  755 ${SRCDIR}/debian/opennsl_prerm    ${OPENNSL_DIR}/DEBIAN/prerm
}

opennsl_clean() {
    rm -vfr ${OPENNSL_DIR}
}

opennsl_make_deb() {
    fakeroot dpkg-deb --build ${OPENNSL_DIR} .
}

opennsl_check() {
    file_exist files/opennsl_setup
    file_exist libopennsl.so.1
    file_exist opennsl.conf
    file_exist linux-kernel-bde.ko
    file_exist linux-user-bde.ko
    file_exist linux-bcm-knet.ko

    file_exist debian/opennsl_control
    file_exist debian/opennsl_postinst
    file_exist debian/opennsl_prerm
}

opennsl_make_pkg() {
    opennsl_make_dirs
    opennsl_copy_files
    opennsl_make_deb
}

do_distclean() {
    rm -vf *.ko *.deb
    rm -vf gonsld opennsl.conf libopennsl.so.1 libopennsl.so
    rm -fr ${WORKDIR}
}

do_usage() {
    echo "$0 <gonsld | opennsl> ><command>"
    echo "command:"
    echo "  prepare : prepare files."
    echo "  check   : check files."
    echo "  deb     : create deb."
    echo "  clean   : clean files."

    exit 1
}


case $1 in
    gonsld)
	case $2 in
	    make-dirs)  gonsld_make_dirs;;
	    prepare)    gonsld_prepare_files;;
	    copy-files) gonsld_copy_files;;
	    check)      gonsld_check;;
	    clean)      gonsld_clean;;
	    deb)        gonsld_make_pkg;;
	    *)          do_usage;;
	esac
	;;

    opennsl)
	case $2 in
	    make-dirs)  opennsl_make_dirs;;
	    prepare)    opennsl_prepare_files;;
	    copy-files) opennsl_copy_files;;
	    check)      opennsl_check;;
	    clean)      opennsl_clean;;
	    deb)        opennsl_make_pkg;;
	    *)          do_usage;;
	esac
	;;

    distclean)
	do_distclean
	;;

    *)
        do_usage
        ;;
esac
