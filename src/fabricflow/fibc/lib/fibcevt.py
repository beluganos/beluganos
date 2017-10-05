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
FIBC Event
"""

from ryu.controller import event

# pylint: disable=too-few-public-methods
class EventFIBCBase(event.EventBase):
    """
    FIBC Base event
    """
    def __init__(self, msg):
        super(EventFIBCBase, self).__init__()
        self.msg = msg


# pylint: disable=too-few-public-methods
class EventFIBCPortConfig(EventFIBCBase):
    """
    FIBC PortConfig event
    """
    pass


# pylint: disable=too-few-public-methods
class EventFIBCDpPortConfig(EventFIBCBase):
    """
    FIBC DP PortConfig event.
    """
    def __init__(self, msg, dp_id, port_id, enter):
        super(EventFIBCDpPortConfig, self).__init__(msg)
        self.dp_id = dp_id
        self.port_id = port_id
        self.enter = enter


# pylint: disable=too-few-public-methods
class EventFIBCVsPortConfig(EventFIBCBase):
    """
    FIBC VS PortConfig event
    """
    def __init__(self, msg, vs_id, port_id):
        super(EventFIBCVsPortConfig, self).__init__(msg)
        self.vs_id = vs_id
        self.port_id = port_id


# pylint: disable=too-few-public-methods
class EventFIBCPortStatus(EventFIBCBase):
    """
    FIBC VM PortConfig event
    """
    pass


# pylint: disable=too-few-public-methods
class EventFIBCVmConfig(EventFIBCBase):
    """
    FIBC VM Config event
    """
    def __init__(self, msg, enter):
        super(EventFIBCVmConfig, self).__init__(msg)
        self.enter = enter


# pylint: disable=too-few-public-methods
class EventFIBCDpConfig(EventFIBCBase):
    """
    FIBC Dp config event
    """
    def __init__(self, msg, dp_id, enter):
        super(EventFIBCDpConfig, self).__init__(msg)
        self.dp_id = dp_id
        self.enter = enter


# pylint: disable=too-few-public-methods
class EventFIBCDpStatus(EventFIBCBase):
    """
    FIBC Dp status event
    """
    pass


# pylint: disable=too-few-public-methods
class EventFIBCFlowMod(EventFIBCBase):
    """
    FIBC FlowMod event
    """
    pass


# pylint: disable=too-few-public-methods
class EventFIBCGroupMod(EventFIBCBase):
    """
    FIBC GroupMod event
    """
    pass


class EventFIBCPortMap(EventFIBCBase):
    """
    FIBC PortMap event
    cmd: "ADD" "DELETE"

    table: "dp"
    msg: fibcdbm.create_dp()

    table: "idmap"
    msg: fibcdbm.create_idmap()

    table: "port"
    msg: fibcdbm.create_ports()
    """
    def __init__(self, msg, cmd, table):
        super(EventFIBCPortMap, self).__init__(msg)
        self.cmd = cmd
        self.table = table
