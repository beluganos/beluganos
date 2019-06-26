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
FIBC Api Server
"""

import logging
from ryu.base import app_manager
from ryu.lib import hub
from ryu.controller import handler
from fabricflow.fibc.net import fibcnet
from fabricflow.fibc.api import fibcapi
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.lib import fibcevt
from fabricflow.fibc.lib import fibclog
from fabricflow.fibc.lib import fibcryu
from fabricflow.fibc.dbm import fibcdbm

_LOG = logging.getLogger(__name__)

_RE_ID_ANY = "0.0.0.0"

class FIBCApiVmController(object):
    """
    FIBCApi Vm Controller
    """
    def __init__(self, soc, data, app):
        self.soc = soc
        self.app = app
        self.hello = fibcapi.parse_hello(data)

    def _send_evt(self, evt):
        _LOG.info("%s %s %s", evt.mtype, evt.msg, evt)
        self.app.send_event_to_observers(evt)

    def get_id(self):
        """
        get client id.
        """
        return self.hello.re_id # pylint: disable=no-member

    def initialize(self):
        """
        send vm-config(ente) event
        """
        cid = self.get_id()
        if cid in self.app.clients:
            raise KeyError()

        self.app.clients[cid] = self
        _LOG.info("client(VM) %s registerd.", cid)

        if cid != _RE_ID_ANY:
            evt = fibcevt.EventFIBCVmConfig(self.hello, True)
            self._send_evt(evt)

    def finalize(self):
        """
        send vm-config(leave) event
        """
        cid = self.get_id()
        self.app.clients.pop(cid)
        _LOG.info("client(VM) %s unregistered.", cid)

        if cid != _RE_ID_ANY:
            evt = fibcevt.EventFIBCVmConfig(self.hello, False)
            self._send_evt(evt)

    def send_data(self, mtype, data, xid=0):
        """
        send data to vm
        """
        fibcnet.write_fib_msg(self.soc, mtype, xid, data)

    def dispatch(self, hdr, data):
        """
        Dispatch mmessage.
        """
        mtype = hdr[0]
        if mtype == pb.PORT_CONFIG:
            port_conf = fibcapi.parse_port_config(data)
            port_conf_evt = fibcevt.EventFIBCPortConfig(port_conf)
            self._send_evt(port_conf_evt)

        elif mtype == pb.FLOW_MOD:
            flow_mod = fibcapi.parse_flow_mod(data)
            flow_mod_evt = fibcevt.EventFIBCFlowMod(flow_mod)
            self._send_evt(flow_mod_evt)

        elif mtype == pb.GROUP_MOD:
            group_mod = fibcapi.parse_group_mod(data)
            group_mod_evt = fibcevt.EventFIBCGroupMod(group_mod)
            self._send_evt(group_mod_evt)

        else:
            _LOG.warn("Unknown message %s", hdr)
            if fibclog.dump_msg():
                _LOG.debug("%s", data)


class FIBCApiDpController(object):
    """
    FIBCApi DP Controller
    """

    def __init__(self, soc, data, app):
        self.soc = soc
        self.app = app
        self.que = hub.Queue() # (mtype, data, xid)
        self.hello = fibcapi.parse_ff_hello(data)
        self.dpath = fibcryu.FFDatapath(self.que, self.get_id())
        self.ports = list()

    def _send_evt(self, evt):
        _LOG.info("%s %s", evt.mtype, evt)
        self.app.send_event_to_observers(evt)

    def _process_que(self):
        _LOG.debug("_process_que %s started.", self.get_id())

        while True:
            msg = self.que.get() # msg is None or (mtype, data, xid)
            if msg is None:
                break

            self.send_data(*msg)

        _LOG.debug("_process_que %s exit.", self.get_id())

    def get_id(self):
        """
        get client id
        """
        return self.hello.dp_id # pylint: disable=no-member

    def initialize(self):
        """
        send MultiPartRequet(PoerDesc)
        """
        fibcdbm.dps().add_dp(self.dpath)
        _LOG.debug("FFDatapath registered. %s", self.dpath)

        msg = fibcapi.new_ff_multipart_request_portdesc(self.get_id(), internal=True)
        self.send_msg(pb.FF_MULTIPART_REQUEST, msg)

        hub.spawn(self._process_que)

    def finalize(self):
        """
        send event (leave)
        """
        self.que.put(None)

        evt = fibcevt.EventFIBCEnterDP(self.dpath, False, self.ports)
        self._send_evt(evt)

        fibcdbm.dps().del_dp(self.get_id())
        _LOG.debug("FFDatapath unregistered. %s", self.dpath)

    def send_msg(self, mtype, msg, xid=0):
        """
        put to send queue.
        """
        data = msg.SerializeToString()
        self.que.put((mtype, data, xid))

    def send_data(self, mtype, data, xid=0):
        """
        write msgs to fibc sock.
        """
        fibcnet.write_fib_msg(self.soc, mtype, xid, data)

    def dispatch(self, hdr, data):
        """
        Dispatch mmessage.
        """
        mtype = hdr[0]
        if mtype == pb.FF_MULTIPART_REPLY:
            msg = fibcapi.parse_ff_multipart_reply(data)
            if msg.mp_type == pb.FFMultipart.PORT_DESC and msg.port_desc.internal: # pylint: disable=no-member

                self.ports = msg.port_desc.port # pylint: disable=no-member

                evt = fibcevt.EventFIBCEnterDP(self.dpath, True, msg.port_desc.port) # pylint: disable=no-member
                self._send_evt(evt)

            else:
                evt = fibcevt.EventFIBCMultipartReply(self.dpath,
                                                      msg,
                                                      fibcnet.get_fib_header_xid(hdr))
                self._send_evt(evt)

        elif mtype == pb.FF_PACKET_IN:
            msg = fibcapi.parse_ff_packet_in(data)
            evt = fibcevt.EventFIBCPacketIn(self.dpath, msg, fibcnet.get_fib_header_xid(hdr))
            self._send_evt(evt)

        elif mtype == pb.FF_PORT_STATUS:
            msg = fibcapi.parse_ff_port_status(data)
            msg = fibcryu.FFPortStatus(self.dpath, msg.reason, msg.port) # pylint: disable=no-member
            evt = fibcevt.EventFIBCFFPortStatus(msg)
            self._send_evt(evt)

        elif mtype == pb.FF_L2ADDR_STATUS:
            msg = fibcapi.parse_ff_l2addr_status(data)
            evt = fibcevt.EventFIBCFFL2AddrStatus(self.dpath, msg)
            self._send_evt(evt)

        else:
            _LOG.warn("Unknown message dp_id:%d %s", self.dpath.id, hdr)
            if fibclog.dump_msg():
                _LOG.debug("%s", data)


class FIBCApiApp(app_manager.RyuApp):
    """
    FIBC Api Server
    """

    _EVENTS = [
        fibcevt.EventFIBCVmConfig,
        fibcevt.EventFIBCPortConfig,
        fibcevt.EventFIBCFlowMod,
        fibcevt.EventFIBCGroupMod,
        fibcevt.EventFIBCEnterDP,
        fibcevt.EventFIBCMultipartReply,
        fibcevt.EventFIBCPacketIn,
        fibcevt.EventFIBCFFPortStatus,
        fibcevt.EventFIBCFFL2AddrStatus,
    ]

    def __init__(self, *args, **kwargs):
        super(FIBCApiApp, self).__init__(*args, **kwargs)
        self.clients = dict()


    def _stream_server(self, host):
        sserver = hub.StreamServer(host, self.on_connect)
        sserver.serve_forever()


    def start_server(self, host):
        """
        Start server
        host: (addr, port)
        """
        hub.spawn(self._stream_server, host)
        _LOG.info("Server started.")


    def send_to_monitor(self, mtype, data, xid):
        """
        send msgs to monitor.
        """
        monitor = self.clients.get(_RE_ID_ANY, None)
        if monitor is not None:
            monitor.send_data(mtype, xid, data)

    def on_connect(self, soc, addr):
        """
        Receive and process message from clinet.
        - addr: (ip, port)
        """
        _LOG.info("NewConnection %s", addr)

        hdr, data = fibcnet.read_fib_msg(soc)

        try:
            if hdr[0] == pb.HELLO:
                ctl = FIBCApiVmController(soc, data, self)
            elif hdr[0] == pb.FF_HELLO:
                ctl = FIBCApiDpController(soc, data, self)
            else:
                raise TypeError()

            if fibclog.dump_msg():
                _LOG.debug("%s", ctl.hello)

            ctl.initialize()

            while True:
                try:
                    hdr, data = fibcnet.read_fib_msg(soc)
                    if hdr is None:
                        _LOG.info("Disconnected %s", addr)
                        break

                    _LOG.debug("Recv %s", hdr)

                    ctl.dispatch(hdr, data)

                except Exception as ex: # pylint: disable=broad-except
                    _LOG.exception("%s", ex)

            ctl.finalize()

        except Exception as ex: # pylint: disable=broad-except
            _LOG.exception("Invalid message. %s %s", hdr, ex)

        finally:
            soc.close()
            _LOG.debug("Connection closed %s", addr)


    def send_to_vm(self, mtype, msg, xid=0):
        """
        Send message to ribc.
        """

        cid = msg.re_id
        if fibclog.dump_msg():
            _LOG.debug("%s %s %s %s", cid, mtype, msg, xid)

        try:
            data = msg.SerializeToString()
            client = self.clients[cid]
            self.send_to_monitor(mtype, data, xid)
            client.send_data(mtype, data, xid)

        except KeyError as err:
            _LOG.warn("client not exist. id:%s, %s", cid, err)


    @handler.set_ev_cls([fibcevt.EventFIBCPortStatus,
                         fibcevt.EventFIBCL2AddrStatus,
                         fibcevt.EventFIBCDpStatus], handler.MAIN_DISPATCHER)
    def _send_msg_to_vm_handler(self, evt):
        """
        process event to send to ribc.
        """
        msg = evt.msg
        mtype = evt.mtype
        _LOG.info("send_to_vm: %s %s", mtype, msg)
        self.send_to_vm(mtype, msg)
