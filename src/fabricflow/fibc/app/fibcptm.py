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
FIBC Port Manager
"""

import logging
from ryu.base import app_manager
from ryu.controller import handler
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibcevt
from fabricflow.fibc.lib import fibclog

_LOG = logging.getLogger(__name__)

def get_ready_ports(portmap, port):
    """
    checi if all port associated
    """
    ready_ports = []
    if port.is_vm_ready():
        # vmの状態が有効な場合、datapath(dp/vs)の状態を見て判定する
        lw_port = portmap.lower_port(port)
        if lw_port.is_datapath_ready():
            ready_ports.append(port)
        else:
            # todo: check all slaves of lw_port
            pass

    if port.is_datapath_ready():
        # datapath(dp/vs)の状態が有効な場合、vm状態を見て判定する
        for up_port in portmap.upper_ports(port):
            if up_port.is_vm_ready():
                ready_ports.append(up_port)

    return ready_ports


class FIBCPtmApp(app_manager.RyuApp):
    """
    FIBC Pakcet forwarding App
    """

    _EVENTS = [
        fibcevt.EventFIBCDpStatus,
        fibcevt.EventFIBCPortStatus,
    ]

    def create(self, cfg):
        """
        Create Tables from config.
        """
        for dpath in cfg.dpaths:
            msg = fibcdbm.create_dp(dpath)
            evt = fibcevt.EventFIBCPortMap(msg, "add", "dp")
            self.on_port_map(evt)

        for router in cfg.routers:
            msg = fibcdbm.create_idmap(router)
            evt = fibcevt.EventFIBCPortMap(msg, "add", "idmap")
            self.on_port_map(evt)

            for port in fibcdbm.create_ports(router):
                evt = fibcevt.EventFIBCPortMap(port, "add", "port")
                self.on_port_map(evt)


    def send_port_status_event(self, port, status):
        """
        Send PortStatus event
        """
        vm_port = port["vm"]
        pts = pb.PortStatus(
            status=status,
            re_id=vm_port.id,
            port_id=vm_port.port,
            ifname=port["name"].name,
        )

        evt = fibcevt.EventFIBCPortStatus(pts)
        self.send_event_to_observers(evt)


    def send_port_status_if_ready(self, port, status):
        """
        Send PortStatus event if port is ready.
        """
        for ready_port in get_ready_ports(fibcdbm.portmap(), port):
            self.send_port_status_event(ready_port, status)


    def _send_dp_status(self, re_id, enter):
        status = pb.DpStatus.ENTER if enter else pb.DpStatus.LEAVE
        msg = pb.DpStatus(
            status=status,
            re_id=re_id,
        )
        evt = fibcevt.EventFIBCDpStatus(msg)
        self.send_event_to_observers(evt)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCVmConfig, handler.MAIN_DISPATCHER)
    def on_vm_config(self, evt):
        """
        Process VmConfig event.
        evt,msg: pb.Hello
        """
        msg = evt.msg

        _LOG.debug("VmConfig: re_id:%s enter:%s", msg.re_id, evt.enter)

        entry = fibcdbm.idmap().find_by_re_id(msg.re_id)
        if entry.update_vm_status(evt.enter):
            _LOG.debug("send DpStatus on VmConfig. %s %s", msg.re_id, evt.enter)
            self._send_dp_status(msg.re_id, evt.enter)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCDpConfig, handler.MAIN_DISPATCHER)
    def on_dp_config(self, evt):
        """
        Process DpConfig event.
        evt,msg: None
        """
        _LOG.debug("DpConfig: dp_id:%d enter:%s", evt.dp_id, evt.enter)

        entry = fibcdbm.idmap().find_by_dp_id(evt.dp_id)
        if entry.update_dp_status(evt.enter):
            _LOG.debug("send DpStatus on DpConfig. %s %s", entry["re_id"], evt.enter)
            self._send_dp_status(entry["re_id"], evt.enter)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCPortConfig, handler.MAIN_DISPATCHER)
    def on_port_config(self, evt):
        """
        Process PortConfig event
        evt.msg: instance of pb.PortConfig
        """
        msg = evt.msg
        if fibclog.dump_msg():
            _LOG.debug("%s", msg)

        try:
            port = fibcdbm.portmap().find_by_name(re_id=msg.re_id, name=msg.ifname)
            if msg.cmd == pb.PortConfig.ADD:
                port.update_vm(msg.value)
                self.send_port_status_if_ready(port, "UP")

            elif msg.cmd == pb.PortConfig.DELETE:
                self.send_port_status_if_ready(port, "DOWN")
                port.update_vm(0)

            else:
                pass

        except KeyError:
            _LOG.warn("vm port not registered. re_id:%s, ifname:%s",
                      msg.re_id, msg.ifname)

        except Exception as expt:
            _LOG.exception(expt)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCVsPortConfig, handler.MAIN_DISPATCHER)
    def on_vsport_config(self, evt):
        """
        Process VS Port Config event
        evt.msg: fibc.net.ffpacket.FFPacket
        """
        pkt = evt.msg
        if fibclog.dump_msg():
            _LOG.debug("%s", pkt)

        try:
            port = fibcdbm.portmap().find_by_name(re_id=pkt.re_id, name=pkt.ifname)
            if port.update_vs(evt.vs_id, evt.port_id):
                self.send_port_status_if_ready(port, "UP")

        except KeyError:
            _LOG.warn("vs port not registered. re_id:%s, ifname:%s",
                      pkt.re_id, pkt.ifname)

        except Exception as expt:
            _LOG.exception(expt)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCDpPortConfig, handler.MAIN_DISPATCHER)
    def on_dpport_config(self, evt):
        """
        Process DP Port Config event
        evt.msg: ryu.controller.Port
        """
        if fibclog.dump_msg():
            _LOG.debug("%s", evt.msg)

        try:
            port = fibcdbm.portmap().find_by_dp(dp_id=evt.dp_id, port_id=evt.port_id)
            if evt.enter:
                port["dpenter"] = True
                self.send_port_status_if_ready(port, "UP")

            else:
                self.send_port_status_if_ready(port, "DOWN")
                port["dpenter"] = False

        except KeyError as expt:
            _LOG.warn("dp port not registered. dpid:%d, port:%d",
                      evt.dp_id, evt.port_id)

        except Exception as expt:
            _LOG.exception(expt)


    @handler.set_ev_cls(fibcevt.EventFIBCPortMap, handler.MAIN_DISPATCHER)
    def on_port_map(self, evt):
        """
        Process PortMap event
        """
        if fibclog.dump_msg():
            _LOG.debug("%s", evt.msg)

        tbl = evt.table
        cmd = evt.cmd
        msg = evt.msg
        if tbl == "dp":
            self._on_port_map_dp(cmd, msg)

        elif tbl == "port":
            self._on_port_map_port(cmd, msg)

        elif tbl == "idmap":
            self._on_port_map_idmap(cmd, msg)

        else:
            _LOG.error("invalid portmap table. %s", tbl)


    @staticmethod
    def _on_port_map_dp(cmd, msg):
        if cmd == "add":
            fibcdbm.dps().add(msg)
        elif cmd == "delete":
            fibcdbm.dps().delete_by_dp_id(msg["dp_id"])
        else:
            _LOG.error("invalid portmap command. %s", cmd)


    @staticmethod
    def _on_port_map_port(cmd, msg):
        if cmd == "add":
            fibcdbm.portmap().add(msg)
        elif cmd == "delete":
            fibcdbm.portmap().delete_by_name_key(msg["name"])
        else:
            _LOG.error("invalid portmap command. %s", cmd)


    @staticmethod
    def _on_port_map_idmap(cmd, msg):
        if cmd == "add":
            fibcdbm.idmap().add(**msg)
        elif cmd == "delete":
            fibcdbm.idmap().delete_by_re_id(msg["re_id"])
        else:
            _LOG.error("invalid portmap command. %s", cmd)
