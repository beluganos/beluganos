#! /bin/bash
# -*- coding: utf-8 -*-

. /etc/lagopus/dpdk.conf

sudo service openvswitch stop

CONF=/etc/lagopus/lagopus.conf
PID=/var/run/lagopus.pid
LOG=/var/log/lagopus.log

CORE_MASK=$CPU_MASK
PORT_MASK=$PORT_MASK

sudo lagopus -l $LOG -p $PID -C $CONF -d -- -c $CORE_MASK -n 2 -- -p $PORT_MASK  --core-assign=minimum

