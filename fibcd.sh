#! /bin/bash
# -*- coding: utf-8 -*-

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

OPTS+=" --config-file /etc/fabricflow/fibc.conf"
OPTS+=" --log-config-file /etc/fabricflow/fibc.log.conf"
OPTS+=" --verbose"

ryu-manager ryu.app.ofctl_rest fabricflow.fibc.app.fibcapp $OPTS
