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
FIBC Flow/GroupMod
"""

import logging
from time import sleep
from ryu.base import app_manager
from ryu.lib import ofctl_v1_3 as ofctl
from ryu.controller import dpset
from ryu.controller import handler
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibccnv
from fabricflow.fibc.lib import fibclog
from fabricflow.fibc.lib import fibcevt
from fabricflow.fibc.ofc import ofc

_LOG = logging.getLogger(__name__)
_SEND_MOD_WAIT_SEC = 0.025

def _find_dp_by_re_id(re_id):
    dp_id = fibcdbm.idmap().find_by_re_id(re_id)["dp_id"]
    return fibcdbm.dps().find_by_id(dp_id)


# pylint: disable=no-self-use
class FIBCModApp(app_manager.RyuApp):
    """
    FIBC Flow/GroupMod App
    """
    # pylint: disable=no-self-use
    # pylint: disable=broad-except
    @handler.set_ev_cls([dpset.EventDP,
                         fibcevt.EventFIBCEnterDP], dpset.DPSET_EV_DISPATCHER)
    def on_dp(self, evt):
        """
        Process Dp Entre event
        """
        if not evt.enter:
            return

        try:
            mode = fibcdbm.dps().get_mode(evt.dp.id, "default")
            flow = ofc.flow(mode, -1)
            group = ofc.group(mode, -1)

            flow(evt.dp, None, ofctl)
            sleep(_SEND_MOD_WAIT_SEC)
            group(evt.dp, None, ofctl)

        except Exception as expt:
            _LOG.exception(expt)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCFlowMod, handler.MAIN_DISPATCHER)
    def on_flow_mod(self, evt):
        """
        Process FlowMod event
        """
        mod = evt.msg
        if fibclog.dump_msg():
            _LOG.debug("%s", mod)

        try:
            dpath, mode = _find_dp_by_re_id(mod.re_id)
            if dpath is not None:
                fibccnv.conv_flow(mod, fibcdbm.portmap())
                func = ofc.flow(mode, mod.table)
                func(dpath, mod, ofctl)
                sleep(_SEND_MOD_WAIT_SEC)

        except Exception as expt:
            _LOG.exception(expt)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCGroupMod, handler.MAIN_DISPATCHER)
    def on_group_mod(self, evt):
        """
        Process GroupMod event
        """
        mod = evt.msg
        if fibclog.dump_msg():
            _LOG.debug(mod)

        try:
            dpath, mode = _find_dp_by_re_id(mod.re_id)
            if dpath is not None:
                fibccnv.conv_group(mod, fibcdbm.portmap())
                func = ofc.group(mode, mod.g_type)
                func(dpath, mod, ofctl)
                sleep(_SEND_MOD_WAIT_SEC)

        except Exception as expt:
            _LOG.exception(expt)


    # pylint: disable=broad-except
    @handler.set_ev_cls(fibcevt.EventFIBCFFPortMod, handler.MAIN_DISPATCHER)
    def on_ff_port_mod(self, evt):
        """
        Process PortMod event
        """
        mod = evt.msg
        if fibclog.dump_msg():
            _LOG.debug(mod)

        try:
            dpath, mode = fibcdbm.dps().find_by_id(mod.dp_id)
            func = ofc.port_mod(mode)
            func(dpath, mod, ofctl)

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

        dp_id = evt.dp_id
        port_id = evt.port_id
        link_up = (evt.state & 0x01) == 0 # 0x01:ofp.OFPPS_LINK_DOWN

        try:
            port = fibcdbm.portmap().find_by_dp(dp_id, port_id)
            vs_port = port["vs"]
            vs_dp, mode = fibcdbm.dps().find_by_id(vs_port.id)
            func = ofc.port_mod(mode)

            status = pb.PortStatus.UP if link_up else pb.PortStatus.DOWN
            mod = pb.FFPortMod(
                dp_id=vs_port.id,
                port_no=vs_port.port,
                hw_addr=port["vs_hw_addr"],
                status=status,
            )

            func(vs_dp, mod, ofctl)

            return

        except KeyError as expt:
            _LOG.warn("dp port not registered. dpid:%d, port:%d", dp_id, port_id)

        except Exception as expt:
            _LOG.exception(expt)
