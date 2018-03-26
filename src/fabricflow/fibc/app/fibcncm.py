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

"""
FIBC nc-module functions
"""

# pylint: disable=bare-except
# pylint: disable=invalid-name

import logging
import yaml

from ryu.base import app_manager
from ryu.controller import dpset
from ryu.controller import handler

from fabricflow.fibc.dbm import fibcdbm
_LOG = logging.getLogger(__name__)

def _load_config(config):
    try:
        with open(config, "r+") as f:
            return yaml.load(f)
    except:
        return dict(dps=dict())

def _write_config(config, data):
    with open(config, "w+") as f:
        f.write(yaml.dump(data))

def _get_ports(evt):
    def _port_entry(port):
        return dict(
            port_no=port.port_no,
            hw_addr=port.hw_addr,
            name=port.name,
        )

    def _port_list():
        if not evt.enter:
            return []

        ofp = evt.dp.ofproto
        return [_port_entry(port) for port in evt.ports if port.port_no <= ofp.OFPP_MAX]

    return dict(ports=_port_list())


class FIBCNcmApp(app_manager.RyuApp):
    """
    FIBC nc-module config App
    """

    def __init__(self, *args, **kwargs):
        super(FIBCNcmApp, self).__init__(*args, **kwargs)
        self.config = None


    def init(self, config):
        """
        Initialize
        """
        self.config = config


    @handler.set_ev_cls(dpset.EventDP, dpset.DPSET_EV_DISPATCHER)
    def on_dp(self, evt):
        """
        Process DP enter event.
        """

        dp_id = evt.dp.id

        mode = fibcdbm.dps().get_mode(dp_id)
        if mode is None or mode == "default":
            _LOG.debug("Ignore %s. mode=%s", dp_id, mode)
            return

        datas = _load_config(self.config)
        datas["dps"][dp_id] = _get_ports(evt)
        _write_config(self.config, datas)

        _LOG.debug("new config:%s %s", dp_id, datas)
