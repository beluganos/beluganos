#! /bin/bash
# -*- coding: utf-8 -*-

############
## Common ##
############
REID="10.0.0.1"
DPID=14
# DPTYPE=openflow
# DPTYPE=as5812
# DPTYPE=as7712x4
DPTYPE=./sample.yaml
# DAEMONS="zebra bgpd ospfd ospf6d ripd ripngd isisd pimd ldpd nhrpd"
DAEMONS="zebra ospfd ospf6d"
# VPN="--vpn"

PORTS=`tempfile`
FFCTL="../../../bin/ffctl playbook"
OVERWRITE="--overwrite"
DEBUG="-v"
FFOPT="${DEBUG} ${OVERWRITE}"


###########
## Ports ##
###########
cat <<EOF > ${PORTS}
---

ports:
  MIC:
    eth: [1, 2, 3, 4]
    vlan:
      - 1.10
      - 2.10

  RIC1:
    eth: [5, 6]

  RIC2:
    eth: [7, 8, 9, 10]
    vlan:
      - 9.10
      - 10.10
EOF


## DO NOT EDIT ##
create_playbook() {
    local NAME=$1
    local NID=$2
    local RT=$3
    local RD=$4

    if [ "$RT" != "" ]; then
        RIBXD_OPTS="--rt ${RT} --rd ${RD}"
    fi

    ${FFCTL}             create ${NAME} ${FFOPT}
    ${FFCTL} daemons     create ${NAME} ${DAEMONS} ${FFOPT}
    ${FFCTL} fibc        create ${NAME} --ports ${PORTS} --dp-type ${DPTYPE} --re-id ${REID} --dp-id ${DPID} ${FFOPT}
    ${FFCTL} frr         create ${NAME} --ports ${PORTS} --dp-type ${DPTYPE} --router-id ${REID} ${FFOPT}
    ${FFCTL} gobgp       create ${NAME} ${FFOPT}
    ${FFCTL} gobgpd      create ${NAME} --router-id ${REID} ${FFOPT}
    ${FFCTL} lxd-profile create ${NAME} --ports ${PORTS} --dp-type ${DPTYPE} ${FFOPT}
    ${FFCTL} netplan     create ${NAME} --ports ${PORTS} --dp-type ${DPTYPE} ${FFOPT}
    ${FFCTL} ribtd       create ${NAME} ${FFOPT}
    ${FFCTL} ribxd       create ${NAME} --re-id ${REID} --node-id ${NID} ${VPN} ${RIBXD_OPTS} ${FFOPT}
    ${FFCTL} snmpproxyd-conf create ${NAME} ${FFOPT}
    ${FFCTL} sysctl      create ${NAME} --ports ${PORTS} --dp-type ${DPTYPE} ${FFOPT}
}

## DO NOT EDIT ##
create_inventory() {
    ${FFCTL} inventory create $@ ${FFOPT}
}

## DO NOT EDIT ##
create_common() {
    ${FFCTL} common create --dp-type ${DPTYPE} ${FFOPT}
}

######################
## Create playbooks ##
######################
create_playbooks() {
    # create_playbook <router name> <node-id> [RT RD]
    create_playbook MIC  0
    create_playbook RIC1 10 10:10 100:10
    create_playbook RIC2 20 10:20 100:20

    create_common
    create_inventory MIC RIC1 RIC2
}

## DO NOT EDIT ##
create_playbooks
rm -f ${PORTS}
