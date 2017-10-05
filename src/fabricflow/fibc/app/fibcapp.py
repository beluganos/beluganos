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
FIB Controller
"""

import logging
from ryu.base import app_manager
from ryu.ofproto import ofproto_v1_3 as ofproto
from ryu.lib import hub
from ryu.controller import dpset
from ryu.app.wsgi import WSGIApplication
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibccfg
from fabricflow.fibc.app import fibcapi
from fabricflow.fibc.app import fibcmod
from fabricflow.fibc.app import fibcofc
from fabricflow.fibc.app import fibcptm
from fabricflow.fibc.app import fibcpkt
from fabricflow.fibc.app import fibcwap

_LOG = logging.getLogger(__name__)

def stream_server(host, func):
    """
    spawn service on StreamServer.
    """
    sserver = hub.StreamServer(host, func)
    sserver.serve_forever()


def get_config():
    """
    Load configuration.
    """
    import sys
    from ryu import cfg

    conf = cfg.CONF
    conf.register_opts([
        cfg.StrOpt("cfg_path", default="/etc/fabricflow/fibc.d"),
        cfg.StrOpt("api_addr", default="127.0.0.1"),
        cfg.IntOpt("api_port", default=50051),
    ])
    conf(sys.argv[1:])
    return conf


class FIBCApp(app_manager.RyuApp):
    """
    FIB Controller Main App
    """
    OFP_VERSIONS = [ofproto.OFP_VERSION]

    _CONTEXTS = {
        "dpset"  : dpset.DPSet,
        'wsgi'   : WSGIApplication,
        "apiapp" : fibcapi.FIBCApiApp,
        "modapp" : fibcmod.FIBCModApp,
        "ofcapp" : fibcofc.FIBCOfcApp,
        "ptmapp" : fibcptm.FIBCPtmApp,
        "pktapp" : fibcpkt.FIBCPktApp,
        "webapp" : fibcwap.FIBCRestApp,
    }

    def __init__(self, *args, **kwargs):
        super(FIBCApp, self).__init__(*args, **kwargs)

        config = get_config()

        dps = kwargs["dpset"]
        fibcdbm.create(dps)

        ptm = kwargs["ptmapp"]
        ptm.create(fibccfg.load_dir(config.cfg_path))

        wsgi = kwargs["wsgi"]
        webapp = kwargs["webapp"]
        webapp.create(wsgi)

        apiapp = kwargs["apiapp"]
        hub.spawn(stream_server, (config.api_addr, config.api_port), apiapp.on_connect)
