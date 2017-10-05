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

import unittest
import struct
from fabricflow.fibc.net.ffpacket import FFPacket

class TestFFPacket(unittest.TestCase):
    def setUp(self):
        pass


    def tearDown(self):
        pass


    def test_parse(self):
        data = "\x01\x02\x03"
        buf = struct.pack("!H24s24s", 0, "abc", "ABC") + data
        # exec
        pkt = FFPacket.parser(buf)
        # check
        self.assertEqual(pkt[0].re_id, "abc")
        self.assertEqual(pkt[0].ifname, "ABC")
        self.assertEqual(pkt[2], data)
        

    def test_parse_nodata(self):
        buf = struct.pack("!H24s24s", 0, "abc", "ABC")
        # exec
        pkt = FFPacket.parser(buf)
        # check
        self.assertEqual(pkt[0].re_id, "abc")
        self.assertEqual(pkt[0].ifname, "ABC")
        self.assertEqual(pkt[2], "")


    def test_serialize(self):
        pkt = FFPacket("abc", "ABC")
        buf = pkt.serialize(None, None)

        self.assertEqual(buf, struct.pack("!H24s24s", 0, "abc", "ABC"))


TESTS = [TestFFPacket]

if __name__ == "__main__":
    unittest.main()
