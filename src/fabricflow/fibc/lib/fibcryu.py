# -*- coding: utf-8 -*-

# Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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
FIBC Puseudo ryu objects.
"""

import random
import ryu.ofproto.ofproto_v1_3 as ofproto

class FFDatapath(object):
    """
    ryu.controller.controller.Datapath
    """

    _XID_MASK = 0xffffffff

    def __init__(self, send_q, dp_id):
        self.send_q = send_q
        self.id = dp_id # pylint: disable=invalid-name
        self.xid = random.randint(0, self._XID_MASK)
        self.ofproto = ofproto


    def get_xid(self):
        """
        Get xid
        """
        def _xid():
            xid = (self.xid + 1) & self._XID_MASK
            if xid == 0:
                return 1
            return xid

        self.xid = _xid()
        return self.xid


    def send_msg(self, mtype, msg, xid=0):
        """
        Serialize to binary data and put into send queue.
        """
        data = msg.SerializeToString()

        if xid == 0:
            xid = self.get_xid()

        self.send_q.put((mtype, data, xid))

        return xid


class FFPortStatus(object):
    """
    ryu.ofproto.ofproto_vX_X_parser.OFPortStatus
    """

    # pylint: disable=too-few-public-methods

    def __init__(self, datapath, reason, desc):
        self.datapath = datapath
        self.reason = reason
        self.desc = desc
