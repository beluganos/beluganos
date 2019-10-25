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
from fabricflow.fibc.api import fibcapi_pb2 as pb

class EventFIBCBase(event.EventBase):
    """
    FIBC Base event
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg, mtype=pb.UNSPEC):
        super(EventFIBCBase, self).__init__()
        self.msg = msg
        self.mtype = mtype


class EventFIBCPortConfig(EventFIBCBase):
    """
    FIBC PortConfig event
    msg: ryu,ofproto.OFPPort
    """
    # pylint: disable=too-few-public-methods
    pass


# pylint: disable=too-few-public-methods
class EventFIBCDpPortConfig(EventFIBCBase):
    """
    FIBC DP PortConfig event.
    msg: ryu,ofproto.OFPPort
    """
    def __init__(self, msg, dp_id, port_id, enter, state):
        super(EventFIBCDpPortConfig, self).__init__(msg)
        self.dp_id = dp_id
        self.port_id = port_id
        self.enter = enter
        self.state = state


class EventFIBCVsPortConfig(EventFIBCBase):
    """
    FIBC VS PortConfig event
    msg: ffpacket
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg, vs_id, port_id):
        super(EventFIBCVsPortConfig, self).__init__(msg)
        self.vs_id = vs_id
        self.port_id = port_id


class EventFIBCPortStatus(EventFIBCBase):
    """
    FIBC VM PortConfig event
    msg: pb.PortStatis
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg):
        super(EventFIBCPortStatus, self).__init__(msg, pb.PORT_STATUS)


class EventFIBCVmConfig(EventFIBCBase):
    """
    FIBC VM Config event
    msg; pb.Hello
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg, enter):
        super(EventFIBCVmConfig, self).__init__(msg, pb.HELLO)
        self.enter = enter


class EventFIBCDpConfig(EventFIBCBase):
    """
    FIBC Dp config event
    msg: None
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg, dp_id, enter):
        super(EventFIBCDpConfig, self).__init__(msg)
        self.dp_id = dp_id
        self.enter = enter


class EventFIBCDpStatus(EventFIBCBase):
    """
    FIBC Dp status event
    msg: pb.DpStatus
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg):
        super(EventFIBCDpStatus, self).__init__(msg, pb.DP_STATUS)


class EventFIBCFFPortMod(EventFIBCBase):
    """
    FIBC FFPortMod event
    msg; pb.FFPortMod
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg):
        super(EventFIBCFFPortMod, self).__init__(msg, pb.FF_PORT_MOD)


class EventFIBCFlowMod(EventFIBCBase):
    """
    FIBC FlowMod event
    msg: pb.FlowMod
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg):
        super(EventFIBCFlowMod, self).__init__(msg, pb.FLOW_MOD)


class EventFIBCGroupMod(EventFIBCBase):
    """
    FIBC GroupMod event
    msg: pb.GroupMod
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, msg):
        super(EventFIBCGroupMod, self).__init__(msg, pb.GROUP_MOD)


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
    # pylint: disable=too-few-public-methods

    def __init__(self, msg, cmd, table):
        super(EventFIBCPortMap, self).__init__(msg)
        self.cmd = cmd
        self.table = table


class EventFIBCEnterDP(EventFIBCBase):
    """
    FIBC EnterDP event

    msg: None
    dp : FFDatapath
    enter: True or False
    ports: list of fibcapi.FFPort
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, dpath, enter, ports):
        super(EventFIBCEnterDP, self).__init__(None)
        self.dp = dpath # pylint: disable=invalid-name
        self.enter = enter
        self.ports = ports


class EventFIBCFFPortStatus(EventFIBCBase):
    """
    FIBC FFPortStatus event
    msg: lib.fibcryu.FFPortStatus
    """
    # pylint: disable=too-few-public-methods
    pass


class EventFIBCMultipartRequest(EventFIBCBase):
    """
    FIBC FFMultipart Request event

    msg: pb.FFMultipart.Request
    dpath: FFDatapath
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, dpath, msg):
        super(EventFIBCMultipartRequest, self).__init__(msg. pb.FF_MULTIPART_REQUEST)
        self.dp = dpath # pylint: disable=invalid-name


class EventFIBCMultipartReply(EventFIBCBase):
    """
    FIBC FFMultipart Reply event

    msg: pb.FFMultipart.Reply
    dpath: FFDatapath
    """
    # pylint: disable=too-few-public-methods

    def __init__(self, dpath, msg, xid=0):
        super(EventFIBCMultipartReply, self).__init__(msg, pb.FF_MULTIPART_REPLY)
        self.dp = dpath # pylint: disable=invalid-name
        self.xid = xid


class EventFIBCPacketIn(EventFIBCBase):
    """
    FIBC FFPacketIn event

    msg: pb.FFPacketIn
    dpath: FFDatapath
    """
    def __init__(self, dpath, msg, xid=0):
        super(EventFIBCPacketIn, self).__init__(msg, pb.FF_PACKET_IN)
        self.dp = dpath # pylint: disable=invalid-name
        self.xid = xid


class EventFIBCPacketOut(EventFIBCBase):
    """
    FIBC FFPacketOut event

    msg: pb.FFPacketOut
    dpath: FFDatapath
    """
    def __init__(self, dpath, msg, xid=0):
        super(EventFIBCPacketOut, self).__init__(msg, pb.FF_PACKET_OUT)
        self.datapath = dpath
        self.xid = xid


class EventFIBCL2AddrStatus(EventFIBCBase):
    """
    FIBC L2AddrStatus event

    msg: pb.L2AddrStatus
    dpath: FFDatapath
    """
    def __init__(self, msg):
        super(EventFIBCL2AddrStatus, self).__init__(msg, pb.L2ADDR_STATUS)


class EventFIBCFFL2AddrStatus(EventFIBCBase):
    """
    FIBC FFL2AddrStatus event

    msg: pb.FFL2AddrStatus
    dpath: FFDatapath
    """
    def __init__(self, dpath, msg, xid=0):
        super(EventFIBCFFL2AddrStatus, self).__init__(msg, pb.FF_L2ADDR_STATUS)
        self.datapath = dpath
        self.xid = xid
