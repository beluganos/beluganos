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
FFPacket for ryu.
"""

import struct
from ryu.lib.packet import packet_base
from ryu.lib.packet.ethernet import ethernet

FFPACKET_ETHTYPE = 0x0a0a

class FFPacket(packet_base.PacketBase):
    """
    FFPacket class.
    """

    _PACK_STR = "!H24s24s"
    _MIN_LEN = struct.calcsize(_PACK_STR)

    def __init__(self, re_id, ifname):
        super(FFPacket, self).__init__()
        self.re_id = re_id
        self.ifname = ifname

    @classmethod
    def parser(cls, buf):
        """
        Parse FFPacket from binary.
        """
        _, re_id, ifname = struct.unpack_from(cls._PACK_STR, buf)
        return cls(re_id.rstrip("\0"), ifname.rstrip("\0")), None, buf[cls._MIN_LEN:]


    def serialize(self, payload, prev):
        """
        Serialize FFPacket to binary.
        """
        return struct.pack(self._PACK_STR, 0, self.re_id, self.ifname)


ethernet.register_packet_type(FFPacket, FFPACKET_ETHTYPE)
