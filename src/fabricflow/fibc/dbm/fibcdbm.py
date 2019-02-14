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
FIB Controller database
"""

import logging
from fabricflow.fibc.dbm.portmap import FIBCDbPortMapTable
from fabricflow.fibc.dbm.portmap import FIBCPortEntry
from fabricflow.fibc.dbm.dps import FIBCDbDpTable
from fabricflow.fibc.dbm.idmap import FIBCDbIdMapTable

_LOG = logging.getLogger(__name__)

# pylint: disable=too-few-public-methods
class FIBCDbm(object):
    """
    FIB Controller database manager
    """
    def __init__(self):
        self.portmap = None
        self.idmap = None
        self.dps = None

    def create(self, dpset, cfg):
        """
        Create Table instances
        cfg: fibccfg.Config
        """
        self.dps = FIBCDbDpTable(dpset)
        self.idmap = FIBCDbIdMapTable()
        self.portmap = FIBCDbPortMapTable()

        for dpath in cfg.dpaths:
            self.dps.add_entry(dpath)

        for router in cfg.get_ext_routers():
            self.idmap.add(re_id=router["re_id"], dp_id=router["dp_id"])
            for port in router["ports"]:
                self.portmap.add(FIBCPortEntry.new(**port))


_INSTANCE = FIBCDbm()

def portmap():
    """
    Get Port map table
    """
    return _INSTANCE.portmap

def idmap():
    """
    Get ID map table
    """
    return _INSTANCE.idmap

def dps():
    """
    Get DP ports table
    """
    return _INSTANCE.dps

def show():
    """
    Show datas.
    """
    portmap().show()
    idmap().show()

def dump(writer):
    """
    Dump datas.
    """
    portmap().dump(writer)
    idmap().dump(writer)

def create(dpset, cfg):
    """
    Create and initialize FIBC DBM instance.
    """
    _LOG.info("creating FIBCDbm")

    _INSTANCE.create(dpset, cfg)
