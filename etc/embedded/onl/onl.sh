#! /bin/bash
# -*- coding: utf-8 -*-

# PROXY=http://172.16.0.1:8080
ONSL_BIN=as7712

ONL_VER="04257be"
KNL_VER="3.16.53"
ONSL_DIR=`pwd`/OpenNSL-master
ONSL_GIT=https://github.com/Broadcom-Switch/OpenNSL.git

RELEASE_DIR=`pwd`/RELEASE

set_proxy() {
  if [ -n "$PROXY" ]; then
    export http_proxy=$PROXY
    export https_proxy=$PROXY
  fi
}

inst_docker() {
  set_proxy

  apt -y update
  apt -y dist-upgrade
  apt -y install linux-image-extra-$(uname -r) linux-image-extra-virtual
  apt -y install apt-transport-https ca-certificates curl software-properties-common
  apt -y install python-minimal apt-cacher-ng
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
  add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  apt -y update
  apt -y install docker-ce
}

init_docker() {
  set_proxy

  if [ -n "$PROXY" ]; then
    echo "Add these line in [Service] section."
    echo "Environment=\"http_proxy=$PROXY\""
    echo "Environment=\"https_proxy=$PROXY\""
    echo "hit any key..."
    read wait
    vi /lib/systemd/system/docker.service

    systemctl daemon-reload
    systemctl restart docker.service
  fi

  systemctl stop apt-cacher-ng
  systemctl disable apt-cacher-ng

}

clone_onl() {
  set_proxy
  local ONL_DIR=OpenNetworkLinux

  if [ -e ${ONL_DIR} ]; then
    echo "${ONL_DIR} already exist."
  else
    git clone https://github.com/opencomputeproject/OpenNetworkLinux
    pushd ${ONL_DIR}/
    git checkout -b onl_kvm ${ONL_VER}
    patch -p1 < ../onl.patch
    popd
  fi

  install -pm 755 -t ${ONL_DIR}/ $0
  install -pm 644 -t ${ONL_DIR}/ *.patch
  install -pm 644 -t ${ONL_DIR}/ opennsl_*

  echo "clone onl ok."
}

start_docker() {
  ./docker/tools/onlbuilder -8
}

setup_onl() {
  set_proxy

  if [ -n "$PROXY" ]; then
    echo "Proxy: $PROXY" >> /etc/apt-cacher-ng/acng.conf
    echo "Acquire::http::Proxy \"$PROXY\";" > /etc/apt/apt.conf
    echo "Acquire::http::proxy::127.0.0.1 \"DIRECT\";" >> /etc/apt/apt.conf
    echo "Acquire::https::Proxy \"$PROXY\";" >> /etc/apt/apt.conf
    echo "Acquire::https::proxy::127.0.0.1 \"DIRECT\";" >> /etc/apt/apt.conf
    git config --global url."https://".insteadOf git://
  fi

  apt-cacher-ng -c /etc/apt-cacher-ng/
}

build_onl() {
  set_proxy
  source setup.env
  make amd64
}

copy_ko() {
  local KO_RELEASE_DIR=${RELEASE_DIR}/ko

  pushd packages/base/amd64/kernels/kernel-3.16-lts-x86-64-all/builds/linux-${KNL_VER}/
  make arch/x86/kvm
  make drivers/vhost
  make fs/autofs4
  popd

  mkdir -p ${KO_RELEASE_DIR}
  for path in $(find packages/base/amd64 -name *.ko); do
    install -pm 644  $path ${KO_RELEASE_DIR}/
  done
}

clone_opennsl() {
  set_proxy

  if [ ! -e ${ONSL_DIR} ]; then
    mkdir -p ${ONSL_DIR}
    pushd ${ONSL_DIR}
    git init
    git remote add origin $ONSL_GIT
    git config core.sparseCheckout true
    echo /bin/${ONSL_BIN}        >> .git/info/sparse-checkout
    echo /include                >> .git/info/sparse-checkout
    echo /sdk-6.5.12-gpl-modules >> .git/info/sparse-checkout
    git fetch --depth 1 origin master
    git pull  --depth 1 origin master
    popd
  fi
}

build_opennsl() {
  local ONSL_MODULES=${ONSL_DIR}/sdk-6.5.12-gpl-modules/systems/linux/user/x86-smp_generic_64-2_6
  export KERNDIR=`pwd`/packages/base/amd64/kernels/kernel-3.16-lts-x86-64-all/builds/linux-${KNL_VER}

  pushd ${ONSL_MODULES}
  make
  popd
}

copy_opennsl() {
  local ONSL_BUILD_DIR=${ONSL_DIR}//sdk-6.5.12-gpl-modules/build
  local ONSL_RELEASE_DIR=${RELEASE_DIR}/opennsl

  mkdir -p ${ONSL_RELEASE_DIR}
  for path in $(find ${ONSL_BUILD_DIR} -name *.ko); do
    install -pm 644 $path ${ONSL_RELEASE_DIR}/
  done
}

dpkg_opennsl() {
  local ONSL_RELEASE_DIR=${RELEASE_DIR}/opennsl
  local WORK_DIR=${RELEASE_DIR}/opennsl/work

  mkdir -p ${WORK_DIR}

  install -pd ${WORK_DIR}/usr/bin/
  install -pd ${WORK_DIR}/usr/lib/
  install -pd ${WORK_DIR}/etc/opennsl/drivers/
  install -pd ${WORK_DIR}/DEBIAN/

  install -pm 755 opennsl_control   ${WORK_DIR}/DEBIAN/control
  install -pm 755 opennsl_postinst  ${WORK_DIR}/DEBIAN/postinst
  install -pm 755 opennsl_prerm     ${WORK_DIR}/DEBIAN/prerm

  install -pm 755 opennsl_setup                                ${WORK_DIR}/usr/bin/
  install -pm 644 ${ONSL_DIR}/bin/${ONSL_BIN}/libopennsl.so.1  ${WORK_DIR}/usr/lib/
  install -pm 644 ${ONSL_RELEASE_DIR}/linux-kernel-bde.ko      ${WORK_DIR}/etc/opennsl/drivers/
  install -pm 644 ${ONSL_RELEASE_DIR}/linux-user-bde.ko        ${WORK_DIR}/etc/opennsl/drivers/
  install -pm 644 ${ONSL_RELEASE_DIR}/linux-bcm-knet.ko        ${WORK_DIR}/etc/opennsl/drivers/

  pushd ${WORK_DIR}

  fakeroot dpkg-deb --build . ${ONSL_RELEASE_DIR}

  popd

  rm -fr ${WORK_DIR}
}

usage() {
  echo "0. edit PROXY value in $0"
  echo "1. sudo $0 install"
  echo "2. $0 clone"
  echo "3. cd OpenNetworkLinux"
  echo "--- to build onl ---"
  echo "4. sudo $0 docker"
  echo "5. $0 build"
  echo "--- to build opennsl library ---"
  echo "# after build onl"
  echo "4. $0 clone-opennsl"
  echo "5. sudo $0 docker"
  echo "6. $0 build-opennsl"
}

case $1 in
  install)
    inst_docker
    init_docker
    ;;

  clone)
    clone_onl
    ;;

  docker)
    start_docker
    ;;

  build)
    setup_onl
    build_onl
    copy_ko
    ;;

  clone-opennsl)
    clone_opennsl
    ;;

  build-opennsl)
    build_opennsl
    copy_opennsl
    dpkg_opennsl
    ;;

  dpkg-opennsl)
    dpkg_opennsl
    ;;

  *)
    usage
    ;;

esac

