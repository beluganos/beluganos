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

beluganos_install() {
    OPTS="--with-opennsl=$BEL_ONSL_ENABLE" ./bootstrap.sh
    if [ "${ENABLE_VIRTUALENV}" = "yes" ]; then
        make release
    else
        make install
        make fflow-install
        sudo make fibc-install
    fi
}

netconf_install() {
    if [ "${BEL_NC_ENABLE}" != "yes" ]; then
        return
    fi

    local NC_DIR=${LXD_WORK_DIR}/netconf

    if [ -e $NC_DIR ]; then
        pushd $NC_DIR
    else
        git clone $BEL_NC_URL $NC_DIR || { echo "beluganos_netconf/clone error."; exit 1; }
        pushd $NC_DIR
    fi

    PROXY=${PROXY} ./create.sh beluganos-netconf

    popd
}
