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
FIBC Monitor App
"""

import time
import logging
import pymongo
from ryu.base import app_manager
from ryu.controller import ofp_event
from ryu.controller import dpset
from ryu.controller import handler
from ryu.lib import hub

_LOG = logging.getLogger(__name__)

# pylint:disable=too-few-public-methods
class DbClient(object):
    """
    FIBC Mongodb client
    """
    def __init__(self, host, port, name):
        self.client = pymongo.MongoClient(host, port)
        self.dbs = self.client[name]


    def insert(self, table, data):
        """
        Insert data
        """
        col = self.dbs[table]
        return col.insert_one(data)


# pylint: disable=no-member
class FIBCMonApp(app_manager.RyuApp):
    """
    FIBC Monitor App
    """
    _INTERVAL_SEC = 1
    _DB_ARGS = dict(
        host="localhost",
        port=27017,
        name="fflow_stats",
    )

    def __init__(self, *args, **kwargs):
        super(FIBCMonApp, self).__init__(*args, **kwargs)
        self.dbc = DbClient(**self._DB_ARGS)
        self.datapaths = {}
        self.monitor_thread = hub.spawn(self._monitor)

        _LOG.info("FIBCMonApp started.")


    def _monitor(self):
        while True:
            for datapath in self.datapaths.values():
                self._request_stats(datapath)
            hub.sleep(self._INTERVAL_SEC)


    @staticmethod
    def _request_stats(datapath):
        _LOG.debug('send stats request: %016x', datapath.id)
        ofp = datapath.ofproto
        parser = datapath.ofproto_parser

        req = parser.OFPFlowStatsRequest(datapath, 0, ofp.OFPTT_ALL)
        datapath.send_msg(req)

        req = parser.OFPGroupStatsRequest(datapath, 0, ofp.OFPG_ALL)
        datapath.send_msg(req)

        req = parser.OFPPortStatsRequest(datapath, 0, ofp.OFPP_ANY)
        datapath.send_msg(req)


    @handler.set_ev_cls(dpset.EventDP, dpset.DPSET_EV_DISPATCHER)
    def _on_dp(self, evt):
        dpath = evt.dp
        if evt.enter:
            if dpath.id not in self.datapaths:
                _LOG.info('register datapath: %016x', dpath.id)
                self.datapaths[dpath.id] = dpath

        else:
            if dpath.id in self.datapaths:
                _LOG.info('unregister datapath: %016x', dpath.id)
                del self.datapaths[dpath.id]


    @handler.set_ev_cls(ofp_event.EventOFPFlowStatsReply, handler.MAIN_DISPATCHER)
    def _flow_stats_reply_handler(self, evt):
        dpath = evt.msg.datapath
        stats = evt.msg.body
        self._insert(dpath, stats)


    @handler.set_ev_cls(ofp_event.EventOFPGroupStatsReply, handler.MAIN_DISPATCHER)
    def _group_stats_reply_handler(self, evt):
        dpath = evt.msg.datapath
        stats = evt.msg.body
        self._insert(dpath, stats)


    @handler.set_ev_cls(ofp_event.EventOFPPortStatsReply, handler.MAIN_DISPATCHER)
    def _port_stats_reply_handler(self, evt):
        dpath = evt.msg.datapath
        stats = evt.msg.body
        self._insert(dpath, stats)


    def _insert(self, dpath, stats):
        nowtime = time.time()
        ltime = time.localtime(nowtime)
        dtime = time.strftime("%Y/%m/%d %H:%M:%S", ltime)
        for stat in stats:
            for name, values in stat.to_jsondict().items():
                data = dict(
                    time=nowtime,
                    dtime=dtime,
                    dp_id=dpath.id,
                    stat=values,
                )
                self.dbc.insert(name, data)
