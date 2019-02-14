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
# Open vSwitch
#
init_ovs() {
    local BRIDGE=$1
    local OFCADDR=$2
    local DPID=$3

    sudo ovs-vsctl add-br ${BRIDGE}
    sudo ovs-vsctl set-controller ${BRIDGE} tcp:${OFCADDR}
    if [ "$DPID"x != ""x ]; then
        sudo ovs-vsctl set bridge ${BRIDGE} other-config:datapath-id=${DPID}
    fi
    sudo ovs-vsctl show
    sudo ovs-ofctl show ${BRIDGE}
}
