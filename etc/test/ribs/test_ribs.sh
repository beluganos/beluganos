#! /bin/bash
# -*- coding: utf-8 -*-

GOBGP="$HOME/go/bin/gobgp"
MICOPT="-p 10001"
RICOPT="-p 10002"
RIBSBIN="../../../bin/ribs2c"
RIBSOPT="--ribs-addr localhost:50072"
MICDEV=micbr1
RICDEV=ricbr1
GWADDR=10.0.0.1
GWMAC=11:22:33:44:55:66
SLEEP=1

VPN1_PREFIX=100.0.10.0/24
VPN1_NEXTHOP=10.0.1.1
VPN1_RD=10:10
VPN1_RT=10:5
VPN1_LABEL=10101

VPN2_PREFIX=100.0.20.0/24
VPN2_NEXTHOP=10.0.2.1
VPN2_RD=20:10
VPN2_RT=20:5
VPN2_LABEL=10201

IP1_PREFIX=20.0.10.0/24
IP1_NEXTHOP=20.0.1.1

IP2_PREFIX=20.0.20.0/24
IP2_NEXTHOP=20.0.2.1

do_init() {
    echo "### init ###"

    # sudo ip link add $MICDEV type bridge
    sudo ip link add $RICDEV type bridge
}

do_clean() {
    echo "### clean ###"

    #sudo ip link del $MICDEV
    sudo ip link del $RICDEV
}

do_mic_rib() {
    CMD=$1

    $GOBGP $MICOPT global rib $CMD -a vpnv4 $VPN1_PREFIX label $VPN1_LABEL rd $VPN1_RD rt $VPN1_RT \
           nexthop $VPN1_NEXTHOP origin igp med 10 local-pref 110
    $GOBGP $MICOPT global rib $CMD -a vpnv4 $VPN2_PREFIX label $VPN2_LABEL rd $VPN2_RD rt $VPN2_RT \
           nexthop $VPN2_NEXTHOP origin igp med 10 local-pref 110
}

do_ric_rib() {
    CMD=$1

    $GOBGP $RICOPT global rib $CMD -a ipv4 $IP1_PREFIX nexthop $IP1_NEXTHOP origin egp med 10 local-pref 120
    $GOBGP $RICOPT global rib $CMD -a ipv4 $IP2_PREFIX nexthop $IP2_NEXTHOP origin egp med 10 local-pref 120
}

do_show() {
    echo ""
    echo "[MIC] IPv4 ----------"
    $GOBGP $MICOPT global rib

    echo ""
    echo "[MIC] VPNv4 ---------"
    $GOBGP $MICOPT global rib -a vpnv4

    echo ""
    echo "[RIC] IPv4 ----------"
    $GOBGP $RICOPT global rib

    echo ""
    echo "[RIC] VPNv4 ---------"
    $GOBGP $RICOPT global rib -a vpnv4

    echo ""
    echo "[RIBS] RICS -----"
    $RIBSBIN $RIBSOPT dump rics

    echo ""
    echo "[RIBS] Nexthops -----"
    $RIBSBIN $RIBSOPT dump nexthops

    echo ""
    echo "[RIBS] NexthopMap -----"
    $RIBSBIN $RIBSOPT dump nexthop-map
}

do_usage() {
    echo "$0 <init | clean | show>"
    echo "$0 mic <add-rib | del-rib>"
    echo "$0 ric <add-rib | del-rib>"
}

case $1 in
    init)  do_init;;
    clean) do_clean;;
    show)  do_show;;
    mic)
        case $2 in
            add-rib) do_mic_rib add;;
            del-rib) do_mic_rib del;;
            *)       do_usage;;
        esac
        ;;

    ric)
        case $2 in
            add-rib) do_ric_rib add;;
            del-rib) do_ric_rib del;;
	    *)       do_usage;;
        esac
        ;;

    *) do_usage;;
esac
