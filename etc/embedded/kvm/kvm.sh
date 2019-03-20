#! /bin/bash
# -*- coding: utf-8 -*-

NW_NAME="default"
DOM_FILE="domain.xml"

add_network() {
    local NAME=$1
    virsh net-define    ./networks.xml
    virsh net-autostart $NAME
    virsh net-start     $NAME
}

del_network() {
    local NAME=$1
    virsh net-destroy  $NAME
    virsh net-undefine $NAME
}

list_network() {
    virsh net-list --all
}

add_domain() {
    local XMLFILE=$1
    virsh create $XMLFILE
}

del_domain() {
    local DOMAIN=$1
    virsh shutdown $DOMAIN
}

list_domain() {
    virsh list --all
    cat /var/lib/libvirt/dnsmasq/default.leases
}

do_usage() {
    echo "$0 network <add | del | list>"
    echo "$0 domain  add"
    echo "$0 domain  del <name>"
    echo "$0 domain  list"
}

case $1 in
    network)
        case $2 in
            add) add_network $NW_NAME;;
            del) del_network $NW_NAME;;
            list) list_network;;
            *) do_usage;;
        esac
        ;;
    domain)
        case $2 in
            add) add_domain $DOM_FILE;;
            del) del_domain $3;;
            list) list_domain ;;
            *) do_usae;;
        esac
        ;;
    *) do_usage;;
esac
