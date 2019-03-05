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
# install pip by get_pip.py
#
get_pip() {
    if [ "${ENABLE_VIRTUALENV}" != "yes" ]; then
        wget -nc -P /tmp ${GET_PIP_URL}/${GET_PIP_FILE} || { echo "pip_install/wget error."; exit 1; }
        sudo python /tmp/${GET_PIP_FILE} --proxy="${PROXY}" || { echo "pip_install/python error."; exit 1; }
    fi
}

#
# python virtualenv
#
make_virtenv() {
    if [ "${ENABLE_VIRTUALENV}" = "yes" ]; then
        if [ -d ${VIRTUALENV} ]; then
            echo "${VIRTUALENV} already exist."
        else
            virtualenv ${VIRTUALENV}
        fi
    fi
}


#
# install python packages
#
pip_install() {
    get_pip
    $PIP install -U -r ${INST_HOME}/${PIP_PKG_LIST} || { echo "pip_install/pip error."; exit 1; }
}

#
# Ryu ofdpa patch
#
ryu_patch() {
    cp ./etc/ryu/ryu_ofdpa2.patch /tmp/

    if [ "${ENABLE_VIRTUALENV}" = "yes" ]; then
        pushd ${VIRTUALENV}/lib/python2.7/site-packages
    else
        pushd /usr/local/lib/python2.7/dist-packages
    fi

    $PATCH -b -p1 < /tmp/ryu_ofdpa2.patch

    popd
}
