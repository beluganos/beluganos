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
FIB packet capture
"""

import logging
from dpkt import pcap
from ryu.base import app_manager
from ryu.controller import handler
from ryu.controller import ofp_event
from fabricflow.fibc.lib import fibcevt

_LOG = logging.getLogger(__name__)

class FIBCPcapApp(app_manager.RyuApp):
    """
    FIB packet capture App
    """

    def create(self, path=None):
        """
        Initialize App
        """
        # pylint: disable=attribute-defined-outside-init

        if path:
            try:
                self.pcap = pcap.Writer(open(path, "wb"))
                _LOG.info("Pcap Start. %s", path)

            except OSError as ex:
                _LOG.exception(ex)

        else:
            self.pcap = None
            _LOG.info("Pcap diabled.")


    @handler.set_ev_cls(ofp_event.EventOFPPacketIn, handler.MAIN_DISPATCHER) # pylint: disable=no-member
    def on_packet_in(self, evt):
        """
        Process PacketIN event.
        """
        msg = evt.msg
        self._write_to_pcap(msg.data)


    @handler.set_ev_cls(fibcevt.EventFIBCPacketIn, handler.MAIN_DISPATCHER)
    def on_ff_packet_in(self, evt):
        """
        Process FFPacketIN event.
        """
        msg = evt.msg
        self._write_to_pcap(msg.data)


    def _write_to_pcap(self, pkt):
        if not self.pcap:
            return

        self.pcap.writepkt(pkt)
