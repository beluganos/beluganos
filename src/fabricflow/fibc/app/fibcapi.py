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
from ryu.controller import handler
from fabricflow.fibc.net import fibcnet
from fabricflow.fibc.api import fibcapi
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.lib import fibcevt
from fabricflow.fibc.lib import fibclog

_LOG = logging.getLogger(__name__)

_RE_ID_ANY = "0.0.0.0"

class FIBCApiClients(object):
    """
    FIBCApi Client Manager
    """
    def __init__(self):
        self.clients = dict()
        self.monitors = None


    def register(self, re_id, soc, addr):
        """
        Register client
        """
        if re_id in self.clients:
            return False

        self.clients[re_id] = dict(
            re_id=re_id,
            soc=soc,
            addr=addr,
        )
        _LOG.info("client %s registerd.", re_id)
        return True


    def unregister(self, re_id):
        """
        Unregister client and close socket
        """
        client = self.clients.pop(re_id)
        if client is not None:
            client["soc"].close()
            _LOG.info("client %s unregisterd.", re_id)


    def get_sock(self, re_id):
        """
        Get socket object.
        """
        return self.clients[re_id]["soc"]


    def get_mon(self):
        """
        Get socket for monitor
        """
        client = self.clients.get(_RE_ID_ANY, None)
        if client is not None:
            return client["soc"]


class FIBCApiApp(app_manager.RyuApp):
    """
    FIBC Api Server
    """

    _EVENTS = [
        fibcevt.EventFIBCVmConfig,
        fibcevt.EventFIBCPortConfig,
        fibcevt.EventFIBCFlowMod,
        fibcevt.EventFIBCGroupMod,
    ]

    def __init__(self, *args, **kwargs):
        super(FIBCApiApp, self).__init__(*args, **kwargs)
        self.clients = FIBCApiClients()

    # pylint: disable=broad-except
    def on_connect(self, soc, addr):
        """
        Receive and process message from clinet.
        - addr: (ip, port)
        """
        _LOG.info("NewConnection %s", addr)

        hdr, data = fibcnet.read_fib_msg(soc)
        if hdr is None or hdr[0] != pb.HELLO:
            _LOG.error("Invalid message. %s", hdr)
            soc.close()
            return

        hello = fibcapi.parse_hello(data)
        if fibclog.dump_msg():
            _LOG.debug("%s", hello)

        if not self.clients.register(hello.re_id, soc, addr):
            _LOG.error("re_id:%s already exist", hello.re_id)
            soc.close()
            return

        vm_evt = fibcevt.EventFIBCVmConfig(hello, True)
        self.send_evt(vm_evt)

        while True:
            try:
                hdr, data = fibcnet.read_fib_msg(soc)
                if hdr is None:
                    _LOG.info("Disconnected %s", addr)
                    break

                self.dispatch(hdr, data)

            except Exception as ex:
                _LOG.exception("%s", ex)

        self.clients.unregister(hello.re_id)

        vm_evt = fibcevt.EventFIBCVmConfig(hello, False)
        self.send_evt(vm_evt)


    def send_evt(self, evt):
        """
        Send Event
        """
        _LOG.info("%s", evt.msg)
        self.send_event_to_observers(evt)


    def dispatch(self, hdr, data):
        """
        Dispatch mmessage.
        """
        mtype = hdr[0]
        if mtype == pb.PORT_CONFIG:
            port_conf = fibcapi.parse_port_config(data)
            port_conf_evt = fibcevt.EventFIBCPortConfig(port_conf)
            self.send_evt(port_conf_evt)

        elif mtype == pb.FLOW_MOD:
            flow_mod = fibcapi.parse_flow_mod(data)
            flow_mod_evt = fibcevt.EventFIBCFlowMod(flow_mod)
            self.send_evt(flow_mod_evt)

        elif mtype == pb.GROUP_MOD:
            group_mod = fibcapi.parse_group_mod(data)
            group_mod_evt = fibcevt.EventFIBCGroupMod(group_mod)
            self.send_evt(group_mod_evt)

        else:
            _LOG.warn("Unknown message %s", hdr)
            if fibclog.dump_msg():
                _LOG.debug("%s", data)


    # pylint: disable=no-self-use
    @handler.set_ev_cls(fibcevt.EventFIBCPortStatus, handler.MAIN_DISPATCHER)
    def on_port_status(self, evt):
        """
        Process PortStatus event
        """
        msg = evt.msg
        _LOG.info("PortStatus: %s", msg)

        try:
            mon = self.clients.get_mon()
            if mon is not None:
                fibcnet.write_fib_msg(mon, pb.PORT_STATUS, 0, msg.SerializeToString())

            client = self.clients.get_sock(msg.re_id)
            fibcnet.write_fib_msg(client, pb.PORT_STATUS, 0, msg.SerializeToString())

        except KeyError as err:
            _LOG.warn("client not exist. re_id:%s, %s", msg.re_id, err)


    # pylint: disable=no-self-use
    @handler.set_ev_cls(fibcevt.EventFIBCDpStatus, handler.MAIN_DISPATCHER)
    def on_dp_status(self, evt):
        """
        Process DpStatus event
        msg: pb.DpStatus
        """
        msg = evt.msg
        _LOG.info("DpStatus: %s", msg)

        try:
            mon = self.clients.get_mon()
            if mon is not None:
                fibcnet.write_fib_msg(mon, pb.DP_STATUS, 0, msg.SerializeToString())

            client = self.clients.get_sock(msg.re_id)
            fibcnet.write_fib_msg(client, pb.DP_STATUS, 0, msg.SerializeToString())

        except KeyError as err:
            _LOG.warn("client not exist. re_id:%s, %s", msg.re_id, err)
