#! /bin/bash
# -*- coding: utf-8 -*-

# Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

#
# install go-lang
#
golang_install() {
    local GO_FILE=go${GO_VER}.linux-amd64.tar.gz

    echo "Downloading ${GO_URL}/${GO_FILE}"
    wget -nc -P /tmp ${GO_URL}/${GO_FILE} || { echo "golang_install/wget error."; exit 1; }

    echo "Extracting /tmp/${GO_FILE}"
    sudo tar xf /tmp/${GO_FILE} -C /usr/local || { echo "golang_install/tar error."; exit 1; }
}

#
# install go packages
#
gopkg_install() {
    for PKG in ${GO_PKGS}; do
        echo "go get ${PKG}"
        go get -u ${PKG} || { echo "gopkg_install error."; exit 1; }
    done
}

#
# install protobuf
#
protoc_install() {
    local PROTOC_FILE=protoc-${PROTOC_VER}-linux-x86_64.zip

    echo "Downloading ${PROTOC_URL}/${PROTOC_FILE}"
    wget -nc -P /tmp ${PROTOC_URL}/${PROTOC_FILE} || { echo "protoc_install/wget error."; exit 1; }

    echo "Extracting /tmp/${PROTOC_FILE}"
    sudo unzip -o -d /usr/local/go /tmp/${PROTOC_FILE} || { echo "protoc_install/unzip error."; exit 1; }

    sudo chmod +x /usr/local/go/bin/protoc
}

#
# patch for netlink
#
netlink_patch() {
    cp ./etc/netlink/netlink_gonla.patch /tmp/
    cp ./etc/netlink/netlink_ip6tnl.patch /tmp/

    pushd ~/go/src/github.com/vishvananda/netlink/
    patch -p1 < /tmp/netlink_gonla.patch
    patch -p1 < /tmp/netlink_ip6tnl.patch
    go install || { echo "netlink_patch/install error."; exit 1; }
    popd
}

#
# Upgrade GoBGP
#
gobgp_upgrade() {
    # prepare
    cp ./etc/gobgp/gobgp-zapi-ver5-v1.33.patch /tmp/gobgp-zapi-ver5-v1.33.patch
    cp ./etc/gobgp/gobgp-encap-ipv6-v1.33.patch /tmp/gobgp-encap-ipv6-v1.33.patch
    cp ./etc/gobgp/gobgp-mp-nexthop-v1.33.patch /tmp/gobgp-mp-nexthop-v1.33.patch
    cp ./etc/gobgp/gobgp-vmx-v1.33.patch /tmp/gobgp-vmx-v1.33.patch
    cp ./etc/gobgp/gobgp-influxdata-v1.33.patch /tmp/gobgp-influxdata-v1.33.patch

    pushd ~/go/src/github.com/osrg/gobgp

    # change to specific version.
    git checkout -B ${GOBGP_VER} ${GOBGP_VER}

    # patch
    patch -p1 < /tmp/gobgp-zapi-ver5-v1.33.patch
    patch -p1 < /tmp/gobgp-encap-ipv6-v1.33.patch
    patch -p1 < /tmp/gobgp-mp-nexthop-v1.33.patch
    patch -p1 < /tmp/gobgp-vmx-v1.33.patch
    patch -p1 < /tmp/gobgp-influxdata-v1.33.patch

    # reinstall
    go install ./gobgpd || { echo "gobgp_checkout error."; exit 1; }
    go install ./gobgp || { echo "gobgp_checkout error."; exit 1; }

    popd
}

#
# patch for snmp library.
#
snmplib_patch() {
    cp ./etc/snmp/romonLogicalis-asn1.patch /tmp/

    pushd ~/go/src/github.com/PromonLogicalis/asn1
    patch -p1 < /tmp/romonLogicalis-asn1.patch
    go install || { echo "snmplib_patch/install error."; exit 1; }
    popd
}
