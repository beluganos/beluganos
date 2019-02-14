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
FIBC Pakcet forwarding.
"""

import logging
from ryu.base import app_manager
from ryu.controller import handler
from ryu.controller import ofp_event
from ryu.lib.packet import ethernet
from ryu.lib.packet import vlan
from fabricflow.fibc.api import fibcapi
from fabricflow.fibc.net import ffpacket
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibcevt
from fabricflow.fibc.lib import fibclog
from fabricflow.fibc.ofc import ofc

_LOG = logging.getLogger(__name__)

def hexdump(datas):
    """
    output hexdump,
    """
    import struct
    cnt = 0
    line = ""
    for data in datas:
        line += "{0:02x} ".format(struct.unpack("B", data)[0])
        cnt += 1
        if cnt == 16:
            _LOG.info("%s", line)
            line = ""
            cnt = 0
    if line:
        _LOG.info("%s", line)


def _parse_pkt_hdr(data):
    vlan_hdr = None
    ff_hdr = None

    _, cls, data = ethernet.ethernet.parser(data)

    if cls == vlan.vlan:
        vlan_hdr, cls, data = vlan.vlan.parser(data)

    if cls == ffpacket.FFPacket:
        ff_hdr, _, _ = ffpacket.FFPacket.parser(data)

    return ff_hdr, vlan_hdr


# pylint: disable=no-self-use
class FIBCPktApp(app_manager.RyuApp):
    """
    FIBC Pakcet forwarding App
    """

    _EVENTS = [
        fibcevt.EventFIBCVsPortConfig,
    ]

    @handler.set_ev_cls(ofp_event.EventOFPPacketIn, handler.MAIN_DISPATCHER) # pylint: disable=no-member
    def on_packet_in(self, evt):
        """
        Process PacketIN event.
        """
        msg = evt.msg
        dp_id = msg.datapath.id
        port_id = get_in_port(msg)
        self._on_packet_in(msg, dp_id, port_id)


    @handler.set_ev_cls(fibcevt.EventFIBCPacketIn, handler.MAIN_DISPATCHER)
    def on_ff_packet_in(self, evt):
        """
        Process FFPacketIN event.
        """
        msg = evt.msg
        dp_id = msg.dp_id
        port_id = msg.port_no
        self._on_packet_in(msg, dp_id, port_id)


    # pylint: disable=broad-except
    # pylint: disable=no-member
    def _on_packet_in(self, msg, dp_id, port_id):
        try:
            if fibclog.dump_msg():
                _LOG.debug("packet_in(%s)", msg)

            ffpkt, vlan_hdr = _parse_pkt_hdr(msg.data)

            if ffpkt is not None:
                _LOG.debug("%s, (%d, %d)", ffpkt, dp_id, port_id)

                ffpkt_evt = fibcevt.EventFIBCVsPortConfig(ffpkt, dp_id, port_id)
                self.send_event_to_observers(ffpkt_evt)

            else:
                _LOG.debug("PacketIN (%d, %d)", dp_id, port_id)
                if fibclog.dump_pkt():
                    hexdump(msg.data)

                self.forward_pkt(vlan_hdr, msg.data, dp_id, port_id)

        except Exception as expt:
            _LOG.exception(expt)
            hexdump(msg.data)


    def forward_pkt(self, vlan_hdr, data, dp_id, port_id):
        """
        forward packet.
        """
        strip_vlan = True if vlan_hdr is not None and \
                     vlan_hdr.vid == fibcapi.OFPVID_UNTAGGED else False

        try:
            port = fibcdbm.portmap().find_by_dp(dp_id, port_id)
            vs_port = port["vs"]

            _LOG.debug("forwarding DP(%d, %d) -> VS(%d, %d)",
                       dp_id, port_id, vs_port.id, vs_port.port)

            dpath, mode = fibcdbm.dps().find_by_id(vs_port.id)
            pkt_out = ofc.pkt_out(mode)
            pkt_out(dpath, vs_port.port, strip_vlan, data)

            return

        except KeyError:
            # it may be packet from vs.
            pass

        try:
            port = fibcdbm.portmap().find_by_vs(dp_id, port_id)
            dp_port = port["dp"]

            _LOG.debug("forwarding VS(%d, %d) -> DP(%d, %d)",
                       dp_id, port_id, dp_port.id, dp_port.port)

            dpath, mode = fibcdbm.dps().find_by_id(dp_port.id)
            pkt_out = ofc.pkt_out(mode)
            pkt_out(dpath, dp_port.port, strip_vlan, data)

            return

        except KeyError:
            _LOG.warn("drop src(%d, %d)", dp_id, port_id)
            if fibclog.dump_pkt():
                hexdump(data)


def get_in_port(msg):
    """
    Get in_port from packet_in msg.
    """
    for field in msg.match.fields:
        if field.header == msg.datapath.ofproto.OXM_OF_IN_PORT:
            return field.value

    raise KeyError("in_port field not found.{0}".format(msg.match))
