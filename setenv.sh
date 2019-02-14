#! /bin/bash

# Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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

ARG=$1
if [ "${ARG}"x == "local"x ]; then
    echo "Use local package."
    export PYTHONPATH=$PYTHONPATH:`pwd`/src
    . ${HOME}/mypython/bin/activate
fi

export GOPATH=$HOME/go:`pwd`
export PATH=/usr/local/go/bin:$HOME/go/bin:$PATH
export LD_LIBRARY_PATH=/usr/local/lib:${HOME}/opennsl/bin/as7712
