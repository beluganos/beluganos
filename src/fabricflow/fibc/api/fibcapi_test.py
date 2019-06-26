#! /usr/bin/env python
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
from goryu.ofproto import offlow
from fabricflow.fibc.api import fibcapi as api
import fabricflow.fibc.api.fibcapi_pb2 as pb


class TestVLANFlow(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_1(self):
        pass


class TestApi(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_priority(self):
        datas = [
            ("1.0.0.0/24", 16400+offlow.PRIORITY_BAND*24, 16400),
            ("1.0.0.0/24", 16405+offlow.PRIORITY_BAND*24, 16405),
            ("1.0.0.0/16", 16405+offlow.PRIORITY_BAND*16, 16405),
        ]
        for ip, pri, base in datas:
            self.assertEqual(offlow.priority_for_ipaddr(ip, base), pri)

    def test_l2addr_status(self):
        addrs = [
            api.new_l2addr("11:22:33:44:55:66", 10, 100, "ADD"),
            api.new_l2addr("11:22:33:44:55:77", 11, 101, "DELETE", "eth1"),
        ]
        msg = api.new_l2addr_status("1.1.1.1", addrs)
        data = msg.SerializeToString()
        msg = api.parse_l2addr_status(data)

        self.assertEqual(msg.re_id, "1.1.1.1")

        addr = msg.addrs[0]
        self.assertEqual(addr.hw_addr, "11:22:33:44:55:66")
        self.assertEqual(addr.vlan_vid, 10)
        self.assertEqual(addr.port_id, 100)
        self.assertEqual(addr.reason, pb.L2Addr.ADD)
        self.assertEqual(addr.ifname, "")

        addr = msg.addrs[1]
        self.assertEqual(addr.hw_addr, "11:22:33:44:55:77")
        self.assertEqual(addr.vlan_vid, 11)
        self.assertEqual(addr.port_id, 101)
        self.assertEqual(addr.reason, pb.L2Addr.DELETE)
        self.assertEqual(addr.ifname, "eth1")


    def test_ff_l2addr_status(self):
        addrs = [
            api.new_l2addr("11:22:33:44:55:66", 10, 100, "ADD"),
            api.new_l2addr("11:22:33:44:55:77", 11, 101, "DELETE", "eth1"),
        ]
        msg = api.new_ff_l2addr_status(1234, addrs)
        data = msg.SerializeToString()
        msg = api.parse_ff_l2addr_status(data)

        self.assertEqual(msg.dp_id, 1234)

        addr = msg.addrs[0]
        self.assertEqual(addr.hw_addr, "11:22:33:44:55:66")
        self.assertEqual(addr.vlan_vid, 10)
        self.assertEqual(addr.port_id, 100)
        self.assertEqual(addr.reason, pb.L2Addr.ADD)
        self.assertEqual(addr.ifname, "")

        addr = msg.addrs[1]
        self.assertEqual(addr.hw_addr, "11:22:33:44:55:77")
        self.assertEqual(addr.vlan_vid, 11)
        self.assertEqual(addr.port_id, 101)
        self.assertEqual(addr.reason, pb.L2Addr.DELETE)
        self.assertEqual(addr.ifname, "eth1")


class TestDpMultipartPort(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new(self):
        pass


TESTS = [TestApi, TestVLANFlow, TestDpMultipartPort]

if __name__ == "__main__":
    unittest.main()
