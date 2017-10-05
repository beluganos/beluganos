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
from ryu.lib.packet import packet
from ryu.lib.packet import vlan
from fabricflow.fibc.api import fibcapi
from fabricflow.fibc.net import ffpacket
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibcevt
from fabricflow.fibc.lib import fibclog

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


# pylint: disable=no-self-use
class FIBCPktApp(app_manager.RyuApp):
    """
    FIBC Pakcet forwarding App
    """

    _EVENTS = [
        fibcevt.EventFIBCVsPortConfig,
    ]

    # pylint: disable=broad-except
    # pylint: disable=no-member
    @handler.set_ev_cls(ofp_event.EventOFPPacketIn, handler.MAIN_DISPATCHER)
    def on_packet_in(self, evt):
        """
        Process PacketIN event.
        """
        try:
            if fibclog.dump_msg():
                _LOG.debug("packet_in(%s)", evt.msg)

            msg = evt.msg
            dp_id = msg.datapath.id
            port_id = get_in_port(msg)

            pkt = packet.Packet(msg.data)
            ffpkt = pkt.get_protocol(ffpacket.FFPacket)

            if ffpkt is not None:
                _LOG.debug("%s, (%d, %d)", ffpkt, dp_id, port_id)

                ffpkt_evt = fibcevt.EventFIBCVsPortConfig(ffpkt, dp_id, port_id)
                self.send_event_to_observers(ffpkt_evt)

            else:
                _LOG.debug("PacketIN (%d, %d)", dp_id, port_id)
                if fibclog.dump_pkt():
                    _LOG.debug("%s", pkt)

                self.forward_pkt(pkt, dp_id, port_id)

        except Exception as expt:
            _LOG.exception(expt)
            hexdump(msg.data)



    def forward_pkt(self, pkt, dp_id, port_id):
        """
        forward packet.
        """
        try:
            port = fibcdbm.portmap().find_by_dp(dp_id, port_id)
            vs_port = port["vs"]
            vlan_hdr = pkt.get_protocol(vlan.vlan)
            strip_vlan = True if vlan_hdr is not None and \
                         vlan_hdr.vid == fibcapi.OFPVID_UNTAGGED else False

            _LOG.debug("forwarding DP(%d, %d) -> VS(%d, %d)",
                       dp_id, port_id, vs_port.id, vs_port.port)

            self.packetout(pkt, vs_port.id, vs_port.port, strip_vlan)

            return

        except KeyError:
            # it may be packet from vs.
            pass

        try:
            port = fibcdbm.portmap().find_by_vs(dp_id, port_id)
            dp_port = port["dp"]

            _LOG.debug("forwarding VS(%d, %d) -> DP(%d, %d)",
                       dp_id, port_id, dp_port.id, dp_port.port)

            self.packetout(pkt, dp_port.id, dp_port.port)

            return

        except KeyError:
            _LOG.warn("drop src(%d, %d)", dp_id, port_id)
            if fibclog.dump_pkt():
                _LOG.debug("drop %s", pkt)


    def packetout(self, pkt, dp_id, port_id, strip_vlan=False):
        """
        Send Packetout message.
        """
        dpath = fibcdbm.dps().find_by_id(dp_id)

        actions = [dpath.ofproto_parser.OFPActionOutput(port_id)]
        if strip_vlan:
            actions.insert(0, dpath.ofproto_parser.OFPActionPopVlan())

        msg = dpath.ofproto_parser.OFPPacketOut(datapath=dpath,
                                                buffer_id=dpath.ofproto.OFP_NO_BUFFER,
                                                in_port=dpath.ofproto.OFPP_ANY,
                                                actions=actions,
                                                data=pkt.data)
        if fibclog.dump_msg():
            _LOG.debug("PacketOUT(%s)", msg)

        dpath.send_msg(msg)


def get_in_port(msg):
    """
    Get in_port from packet_in msg.
    """
    for field in msg.match.fields:
        if field.header == msg.datapath.ofproto.OXM_OF_IN_PORT:
            return field.value

    raise KeyError("in_port field not found.{0}".format(msg.match))
