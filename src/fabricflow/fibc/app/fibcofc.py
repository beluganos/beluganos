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

"""
FIBC Mod functions
"""

import logging
from ryu.base import app_manager
from ryu.controller import dpset
from ryu.controller import handler
from ryu.controller import ofp_event

from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibcevt

_LOG = logging.getLogger(__name__)

class FIBCOfcApp(app_manager.RyuApp):
    """
    FIBC OpenFlow Controller App
    """
    _EVENTS = [
        fibcevt.EventFIBCDpPortConfig,
        fibcevt.EventFIBCDpConfig,
    ]

    @handler.set_ev_cls([dpset.EventDP,
                         fibcevt.EventFIBCEnterDP], dpset.DPSET_EV_DISPATCHER)
    def on_dp(self, evt):
        """
        Process DP enter event.
        """
        dp_id = evt.dp.id

        _LOG.debug("dp_id:%s enter:%s", dp_id, evt.enter)

        if fibcdbm.dps().get_mode(dp_id) is None:
            return

        dp_status_evt = fibcevt.EventFIBCDpConfig(None, dp_id=dp_id, enter=evt.enter)
        self.send_event_to_observers(dp_status_evt)

        for port in evt.ports:
            self.send_dp_port_config(evt.dp, port, evt.enter)


    # pylint: disable=no-member
    @handler.set_ev_cls([ofp_event.EventOFPPortStatus,
                         fibcevt.EventFIBCFFPortStatus], handler.MAIN_DISPATCHER)
    def on_port_status(self, evt):
        """
        Process Port Status event.
        """
        msg = evt.msg
        dpath = msg.datapath
        port = msg.desc
        reason = msg.reason

        _LOG.debug("dp_id:%s port:%s reason: %d", dpath.id, port, reason)

        if fibcdbm.dps().get_mode(dpath.id) is None:
            return

        def _enter():
            ofp = dpath.ofproto
            return reason != ofp.OFPPR_DELETE

        self.send_dp_port_config(dpath, port, _enter())


    def send_dp_port_config(self, dpath, port, enter):
        """
        Send DpPortConfig event.
        """
        if enter:
            # if port is down, change enter to False
            ofp = dpath.ofproto
            enter = (port.state & ofp.OFPPS_LINK_DOWN) == 0

        _LOG.info("DpPortConfig: dp_id=%d, port_id=%d, enter=%s",
                  dpath.id, port.port_no, enter)

        evt = fibcevt.EventFIBCDpPortConfig(port, dpath.id, port.port_no, enter)
        self.send_event_to_observers(evt)
