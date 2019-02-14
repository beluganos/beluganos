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

opennsl_install() {
    if [ "$BEL_ONSL_ENABLE" != "yes" ]; then
        return
    fi

    $INST_HOME/opennsl_install.sh install $BEL_ONSL_PLATFORM || { echo "opennsl_install error."; exit 1; }
}

opennsl_pkg_install() {
    if [ "$BEL_ONSL_ENABLE" != "yes" ]; then
        return
    fi

    go get -u ${BEL_ONSL_PKG} || { echo "go-opennsl install error."; exit 1; }
}
